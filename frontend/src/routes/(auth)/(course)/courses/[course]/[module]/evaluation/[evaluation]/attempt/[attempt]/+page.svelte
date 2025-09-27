<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import LoadingSpinner from '$lib/components/ui/loading-spinner.svelte';
	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';
	import type { PageProps } from './$types';
	import { QuestionType } from '$lib/types';
	import {
		Clock,
		ChevronLeft,
		ChevronRight,
		Send,
		AlertTriangle,
		Wifi,
		WifiOff
	} from '@lucide/svelte';
	import { EvaluationAttemptController } from '$lib/controllers/evaluationAttempt';
	import { onMount, onDestroy } from 'svelte';
	import { toast } from 'svelte-sonner';
	import { goto } from '$app/navigation';
	import { validateAnswers, calculateRemainingTime, isAttemptActive } from '$lib/utils/examHelpers';
	import { browser } from '$app/environment';
	import { page } from '$app/state';

	let { data }: PageProps = $props();

	const evaluationAttemptController = new EvaluationAttemptController();

	const attempt = data.attempt;
	const questions = $derived(attempt.questions.toSorted((a, b) => a.id - b.id));

	let currentQuestionIndex = $state(0);
	let answers = $state<Record<string, number[]>>({});
	let timeLeft = $state<number>(0);
	let timer: any; // NodeJS.Timeout equivalent
	let submitting = $state(false);
	let showConfirmSubmit = $state(false);
	let isOnline = $state(true);
	let autoSaveTimer: any;
	let lastSaved = $state<Date | null>(null);

	const currentQuestion = $derived(questions[currentQuestionIndex]);
	const progress = $derived(
		Number((((currentQuestionIndex + 1) / questions.length) * 100).toFixed(2))
	);
	const answeredCount = $derived(Object.keys(answers).length);
	const allAnswered = $derived(answeredCount === questions.length);

	// Check if attempt is still valid
	const attemptIsActive = $derived(
		isAttemptActive(data.attempt, data.attempt.evaluation?.time_limit)
	);

	// Initialize timer and check attempt validity
	onMount(() => {
		if (browser) {
			// Check online status
			isOnline = navigator.onLine;
			window.addEventListener('online', () => (isOnline = true));
			window.addEventListener('offline', () => (isOnline = false));

			// Check if attempt is already submitted
			if (data.attempt.submitted_at) {
				console.log(data.attempt.submitted_at);

				toast.error('Este intento ya ha sido enviado');
				goto(
					`/courses/${page.params.course}/${page.params.module}/evaluation/${page.params.evaluation}/attempt/${page.params.attempt}/results`
				);
				return;
			}

			// Initialize timer if evaluation has time limit
			if (data.attempt.evaluation?.time_limit) {
				timeLeft = calculateRemainingTime(data.attempt, data.attempt.evaluation.time_limit);
				console.log(timeLeft);

				if (timeLeft <= 0) {
					toast.warning('¡Tiempo agotado! El examen se enviará automáticamente.');
					submitExam();
					return;
				}

				timer = setInterval(() => {
					timeLeft--;
					if (timeLeft <= 0) {
						toast.warning('¡Tiempo agotado! El examen se enviará automáticamente.');
						submitExam();
					}
				}, 1000);
			}

			// Auto-save functionality (every 30 seconds)
			autoSaveTimer = setInterval(() => {
				if (Object.keys(answers).length > 0 && isOnline && !submitting) {
					saveProgress();
				}
			}, 30000);

			// Load any existing answers (if continuing an attempt)
			if (data.attempt.answers && data.attempt.answers.length > 0) {
				const existingAnswers = new Map();
				data.attempt.answers.forEach((answer) => {
					existingAnswers.set(answer.attempt_question_id, answer.selected_option_ids);
				});
				answers = existingAnswers;
			}
		}
	});

	onDestroy(() => {
		if (timer) clearInterval(timer);
		if (autoSaveTimer) clearInterval(autoSaveTimer);
		if (browser) {
			window.removeEventListener('online', () => (isOnline = true));
			window.removeEventListener('offline', () => (isOnline = false));
		}
	});

	function formatTime(seconds: number): string {
		const mins = Math.floor(seconds / 60);
		const secs = seconds % 60;
		return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
	}

	function handleAnswer(questionId: number, optionIds: number[]) {
		answers[questionId.toString()] = optionIds;
	}

	function handleSingleChoice(questionId: number, optionId: number) {
		handleAnswer(questionId, [optionId]);
	}

	function handleMultipleChoice(questionId: number, optionId: number, checked: boolean) {
		const currentAnswers = answers[questionId.toString()] || [];
		if (checked) {
			handleAnswer(questionId, [...currentAnswers, optionId]);
		} else {
			handleAnswer(
				questionId,
				currentAnswers.filter((id) => id !== optionId)
			);
		}
	}

	function goToQuestion(index: number) {
		if (index >= 0 && index < questions.length) {
			currentQuestionIndex = index;
		}
	}

	function nextQuestion() {
		goToQuestion(currentQuestionIndex + 1);
	}

	function prevQuestion() {
		goToQuestion(currentQuestionIndex - 1);
	}

	async function saveProgress() {
		try {
			// This would be an API call to save partial progress
			// For now, we'll just update the lastSaved timestamp
			lastSaved = new Date();
		} catch (error) {
			console.warn('Failed to save progress:', error);
		}
	}

	async function submitExam() {
		if (submitting || !attemptIsActive) return;

		// Validate answers before submitting
		const validation = validateAnswers(questions, answers);
		if (!validation.isValid) {
			toast.error(`Error en las respuestas: ${validation.errors.join(', ')}`);
			return;
		}

		if (!isOnline) {
			toast.error('Sin conexión a internet. Por favor verifica tu conexión e intenta de nuevo.');
			return;
		}

		submitting = true;
		try {
			// Convert answers to the required format
			const submissionAnswers = Array.from(answers.entries()).map(
				([questionId, selectedOptionIds]) => ({
					attempt_question_id: questionId,
					selected_option_ids: selectedOptionIds,
					is_correct: false, // This will be calculated by the backend
					points: 0 // This will be calculated by the backend
				})
			);

			// Add empty answers for unanswered questions
			questions.forEach((question) => {
				if (!answers.has(question.id)) {
					submissionAnswers.push({
						attempt_question_id: question.id,
						selected_option_ids: [],
						is_correct: false,
						points: 0
					});
				}
			});

			const result = await evaluationAttemptController.submitAttempt(page.params.attempt, {
				answers: submissionAnswers
			});

			toast.success('¡Examen enviado exitosamente!');

			// Redirect to results page
			goto(
				`/courses/${page.params.course}/module/${page.params.module}/evaluation/${page.params.evaluation}/attempt/${page.params.attempt}/results`
			);
		} catch (error) {
			console.error('Error submitting exam:', error);
			toast.error('Error al enviar el examen. Por favor intenta de nuevo.');
		} finally {
			submitting = false;
			showConfirmSubmit = false;
		}
	}

	function confirmSubmit() {
		const validation = validateAnswers(questions, answers);

		if (validation.warnings.length > 0) {
			showConfirmSubmit = true;
		} else {
			submitExam();
		}
	}

	// Get current answers for the current question
	const currentAnswers = $derived(answers[currentQuestion?.id.toString()] || []);

	// Handle beforeunload to warn user about unsaved changes
	onMount(() => {
		const handleBeforeUnload = (e: BeforeUnloadEvent) => {
			if (Object.keys(answers).length > 0 && !submitting && attemptIsActive) {
				e.preventDefault();
				e.returnValue = 'Tienes un examen en progreso. ¿Estás seguro de que quieres salir?';
				return e.returnValue;
			}
		};

		if (browser) {
			window.addEventListener('beforeunload', handleBeforeUnload);
		}

		return () => {
			if (browser) {
				window.removeEventListener('beforeunload', handleBeforeUnload);
			}
		};
	});
</script>

<div class="mx-auto max-w-4xl p-6">
	<!-- Header with timer and progress -->
	<div class="mb-6 flex items-center justify-between">
		<div>
			<h1 class="text-2xl font-bold">{data.attempt.evaluation?.title}</h1>
			<p class="text-muted-foreground">
				Pregunta {currentQuestionIndex + 1} de {questions.length}
			</p>
		</div>

		<div class="flex items-center gap-4">
			<!-- Connection status -->
			<div class="flex items-center gap-2">
				{#if !isOnline}
					<WifiOff class="h-4 w-4 text-red-600" />
					<span class="text-xs text-red-600"> Sin conexión </span>
				{/if}
			</div>

			<!-- Timer -->
			{#if timeLeft > 0}
				<div
					class="flex items-center gap-2 font-mono text-lg {timeLeft < 300 ? 'text-red-600' : ''}"
				>
					<Clock class="h-5 w-5" />
					{formatTime(timeLeft)}
				</div>
			{/if}
		</div>
	</div>

	<!-- Progress bar -->
	<div class="mb-6">
		<div class="mb-2 flex items-center justify-between">
			<span class="text-muted-foreground text-sm">Progreso del examen</span>
			<div class="flex items-center gap-4">
				<span class="text-muted-foreground text-sm"
					>{answeredCount}/{questions.length} respondidas</span
				>
				{#if lastSaved}
					<span class="text-xs text-green-600">
						Guardado a las {lastSaved.toLocaleTimeString()}
					</span>
				{/if}
			</div>
		</div>
		<div class="h-2 w-full rounded-full bg-gray-200">
			<div
				class="bg-primary h-2 rounded-full transition-all duration-300"
				style="width: {progress}%"
			></div>
		</div>
	</div>

	<!-- Question navigation -->
	<div class="mb-6">
		<div class="flex flex-wrap gap-2">
			{#each questions as question, index}
				<button
					onclick={() => goToQuestion(index)}
					class="h-10 w-10 rounded-lg border-2 text-sm font-medium transition-colors
						{index === currentQuestionIndex
						? 'border-primary bg-primary text-primary-foreground'
						: answers[question.id.toString()] != undefined &&
							  answers[question.id.toString()].length > 0
							? 'border-green-500 bg-green-100 text-green-800'
							: 'border-gray-300 bg-white hover:bg-gray-50'}"
				>
					{index + 1}
				</button>
			{/each}
		</div>
	</div>

	{#if currentQuestion}
		{#key currentQuestion}
			<!-- Question Card -->
			<Card.Root class="mb-6">
				<Card.Header>
					<div class="flex items-start justify-between">
						<div class="flex-1">
							<Card.Title class="text-lg">
								{currentQuestion.text}
							</Card.Title>
							{#if currentQuestion.explanation}
								<Card.Description class="mt-2">
									{currentQuestion.explanation}
								</Card.Description>
							{/if}
						</div>
						<div class="flex items-center gap-2">
							<Badge variant="secondary">
								{currentQuestion.points}
								{currentQuestion.points === 1 ? 'punto' : 'puntos'}
							</Badge>
							<Badge variant={currentQuestion.type === QuestionType.SINGLE ? 'default' : 'outline'}>
								{currentQuestion.type === QuestionType.SINGLE ? 'Una opción' : 'Múltiple opción'}
							</Badge>
						</div>
					</div>
				</Card.Header>

				<Card.Content>
					{#if currentQuestion.type === QuestionType.SINGLE}
						<!-- Single choice question -->
						<div class="space-y-3">
							{#each currentQuestion.answer_options as option}
								<div class="hover:bg-muted/50 flex items-center space-x-3 rounded-lg border p-3">
									<input
										type="radio"
										name="question-{currentQuestion.id}"
										id="option-{option.id}"
										value={option.id}
										checked={currentAnswers.includes(option.id)}
										onchange={() => handleSingleChoice(currentQuestion.id, option.id)}
										class="text-primary h-4 w-4"
									/>
									<Label for="option-{option.id}" class="flex-1 cursor-pointer">
										{option.text}
									</Label>
								</div>
							{/each}
						</div>
					{:else}
						<!-- Multiple choice question -->
						<div class="space-y-3">
							{#each currentQuestion.answer_options as option}
								<div class="hover:bg-muted/50 flex items-center space-x-3 rounded-lg border p-3">
									<input
										type="checkbox"
										id="option-{option.id}"
										checked={currentAnswers.includes(option.id)}
										onchange={(e) =>
											handleMultipleChoice(
												currentQuestion.id,
												option.id,
												(e.target as HTMLInputElement).checked
											)}
										class="text-primary h-4 w-4"
									/>
									<Label for="option-{option.id}" class="flex-1 cursor-pointer">
										{option.text}
									</Label>
								</div>
							{/each}
						</div>
					{/if}
				</Card.Content>
			</Card.Root>
		{/key}
	{/if}

	<!-- Navigation and submit -->
	<div class="flex items-center justify-between">
		<Button variant="outline" onclick={prevQuestion} disabled={currentQuestionIndex === 0}>
			<ChevronLeft class="mr-2 h-4 w-4" />
			Anterior
		</Button>

		<div class="flex gap-2">
			{#if currentQuestionIndex === questions.length - 1}
				<Button
					onclick={confirmSubmit}
					disabled={submitting || !isOnline}
					class="bg-green-600 hover:bg-green-700"
				>
					{#if submitting}
						<LoadingSpinner size="sm" class="mr-2" />
						Enviando...
					{:else}
						<Send class="mr-2 h-4 w-4" />
						Enviar Examen
					{/if}
				</Button>
			{:else}
				<Button onclick={nextQuestion} disabled={currentQuestionIndex === questions.length - 1}>
					Siguiente
					<ChevronRight class="ml-2 h-4 w-4" />
				</Button>
			{/if}
		</div>
	</div>
</div>

<!-- Confirmation dialog -->
{#if showConfirmSubmit}
	<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
		<Card.Root class="w-full max-w-md">
			<Card.Header>
				<Card.Title class="flex items-center gap-2">
					<AlertTriangle class="h-5 w-5 text-orange-500" />
					Confirmar envío
				</Card.Title>
			</Card.Header>
			<Card.Content>
				<div class="mb-4 rounded-lg border border-orange-200 bg-orange-50 p-4">
					<div class="flex items-center">
						<AlertTriangle class="mr-2 h-4 w-4 text-orange-600" />
						<span class="text-sm font-medium text-orange-800">Preguntas sin responder</span>
					</div>
					<p class="mt-1 text-sm text-orange-700">
						Tienes {questions.length - answeredCount} preguntas sin responder. ¿Estás seguro de que quieres
						enviar el examen?
					</p>
				</div>
			</Card.Content>
			<Card.Footer class="flex gap-2">
				<Button variant="outline" onclick={() => (showConfirmSubmit = false)} class="flex-1">
					Cancelar
				</Button>
				<Button
					onclick={submitExam}
					disabled={submitting}
					class="flex-1 bg-green-600 hover:bg-green-700"
				>
					{#if submitting}
						<LoadingSpinner size="sm" class="mr-2" />
						Enviando...
					{:else}
						Enviar de todas formas
					{/if}
				</Button>
			</Card.Footer>
		</Card.Root>
	</div>
{/if}

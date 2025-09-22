<script lang="ts">
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$components/ui/card';
	import { Badge } from '$components/ui/badge';
	import { Button } from '$components/ui/button';
	// Remove Progress import since it doesn't exist
	// Remove Checkbox import since it doesn't exist
	// Remove RadioGroupItem import since it doesn't exist
	import { Label } from '$components/ui/label';
	import { 
		Clock, 
		ChevronLeft, 
		ChevronRight, 
		CheckCircle,
		AlertCircle,
		Send
	} from '@lucide/svelte';
	import { EvaluationAttemptController } from '$lib';
	import { toast } from 'svelte-sonner';
	import { goto } from '$app/navigation';
	import { onMount, onDestroy } from 'svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	const evaluationAttemptController = new EvaluationAttemptController();

	// Exam state
	let currentQuestionIndex = $state(0);
	let answers = $state<Record<number, number[]>>({});
	let timeRemaining = $state<number | null>(null);
	let isSubmitting = $state(false);
	let intervalId: NodeJS.Timeout | null = null;

	// Computed properties
	const currentQuestion = $derived(data.attempt.questions[currentQuestionIndex]);
	const totalQuestions = $derived(data.attempt.questions.length);
	const progressPercentage = $derived(((currentQuestionIndex + 1) / totalQuestions) * 100);
	const canGoNext = $derived(currentQuestionIndex < totalQuestions - 1);
	const canGoPrev = $derived(currentQuestionIndex > 0);
	const hasAnsweredCurrent = $derived(answers[currentQuestion.id] && answers[currentQuestion.id].length > 0);
	const answeredCount = $derived(Object.keys(answers).filter(qId => answers[parseInt(qId)]?.length > 0).length);

	onMount(() => {
		// Initialize time remaining if there's a time limit
		const evaluation = data.attempt.evaluation;
		if (evaluation?.time_limit) {
			const startTime = new Date(data.attempt.started_at).getTime();
			const timeLimitMs = evaluation.time_limit * 60 * 1000; // Convert to milliseconds
			const elapsed = Date.now() - startTime;
			timeRemaining = Math.max(0, Math.floor((timeLimitMs - elapsed) / 1000));

			// Start timer countdown
			if (timeRemaining > 0) {
				intervalId = setInterval(() => {
					if (timeRemaining !== null && timeRemaining > 0) {
						timeRemaining--;
					} else {
						// Auto submit when time runs out
						autoSubmit();
					}
				}, 1000);
			}
		}

		// Initialize answers from existing attempt answers if any
		data.attempt.answers.forEach(answer => {
			answers[answer.attempt_question_id] = answer.selected_option_ids;
		});
	});

	onDestroy(() => {
		if (intervalId) {
			clearInterval(intervalId);
		}
	});

	function selectAnswer(optionId: number) {
		const questionId = currentQuestion.id;
		
		if (currentQuestion.type === 'single_choice') {
			// Single choice: replace the selection
			answers[questionId] = [optionId];
		} else {
			// Multiple choice: toggle the selection
			if (!answers[questionId]) {
				answers[questionId] = [];
			}
			
			const currentAnswers = answers[questionId];
			if (currentAnswers.includes(optionId)) {
				answers[questionId] = currentAnswers.filter(id => id !== optionId);
			} else {
				answers[questionId] = [...currentAnswers, optionId];
			}
		}
	}

	function goToQuestion(index: number) {
		if (index >= 0 && index < totalQuestions) {
			currentQuestionIndex = index;
		}
	}

	function nextQuestion() {
		if (canGoNext) {
			currentQuestionIndex++;
		}
	}

	function prevQuestion() {
		if (canGoPrev) {
			currentQuestionIndex--;
		}
	}

	async function submitAttempt() {
		if (isSubmitting) return;

		// Confirm submission
		if (!confirm('¿Estás seguro de que quieres enviar tu evaluación? Esta acción no se puede deshacer.')) {
			return;
		}

		isSubmitting = true;

		try {
			// Prepare answers in the format expected by the API
			const submissionAnswers = data.attempt.questions.map(question => ({
				attempt_question_id: question.id,
				selected_option_ids: answers[question.id] || []
			}));

			await evaluationAttemptController.submitAttempt(data.attempt.id, {
				answers: submissionAnswers
			});

			toast.success('Evaluación enviada correctamente');
			
			// Redirect to results
			goto(`/courses/${data.courseId}/${data.moduleId}/${data.evaluationId}/results/${data.attempt.id}`);
		} catch (error) {
			console.error('Error submitting attempt:', error);
			toast.error('Error al enviar la evaluación. Intenta nuevamente.');
		} finally {
			isSubmitting = false;
		}
	}

	async function autoSubmit() {
		toast.warning('Se acabó el tiempo. Enviando evaluación automáticamente...');
		await submitAttempt();
	}

	function formatTime(seconds: number): string {
		const minutes = Math.floor(seconds / 60);
		const secs = seconds % 60;
		return `${minutes}:${secs.toString().padStart(2, '0')}`;
	}

	function getTimeColor(): string {
		if (timeRemaining === null) return 'text-muted-foreground';
		if (timeRemaining < 300) return 'text-red-600'; // < 5 minutes
		if (timeRemaining < 600) return 'text-orange-600'; // < 10 minutes
		return 'text-muted-foreground';
	}
</script>

<div class="min-h-screen bg-background">
	<!-- Header with progress and timer -->
	<div class="border-b bg-card sticky top-0 z-10">
		<div class="container mx-auto px-4 py-4">
			<div class="flex items-center justify-between">
				<div class="flex items-center gap-4">
					<h1 class="text-lg font-semibold">{data.attempt.evaluation?.title}</h1>
					<Badge variant="outline">
						Pregunta {currentQuestionIndex + 1} de {totalQuestions}
					</Badge>
				</div>

				{#if timeRemaining !== null}
					<div class="flex items-center gap-2">
						<Clock class="h-4 w-4 {getTimeColor()}" />
						<span class="font-mono font-medium {getTimeColor()}">
							{formatTime(timeRemaining)}
						</span>
					</div>
				{/if}
			</div>

			<div class="mt-4">
				<!-- Simple progress bar -->
				<div class="w-full bg-gray-200 rounded-full h-2">
					<div class="bg-primary h-2 rounded-full transition-all duration-300" style="width: {progressPercentage}%"></div>
				</div>
				<div class="flex justify-between text-sm text-muted-foreground mt-1">
					<span>Respondidas: {answeredCount}/{totalQuestions}</span>
					<span>{Math.round(progressPercentage)}% completado</span>
				</div>
			</div>
		</div>
	</div>

	<div class="container mx-auto px-4 py-8">
		<div class="grid gap-8 lg:grid-cols-4">
			<!-- Question Content -->
			<div class="lg:col-span-3">
				<Card>
					<CardHeader>
						<div class="flex items-start justify-between">
							<div class="flex-1">
								<CardTitle class="text-xl">
									Pregunta {currentQuestionIndex + 1}
								</CardTitle>
								<CardDescription class="mt-2 text-base leading-relaxed">
									{currentQuestion.text}
								</CardDescription>
							</div>
							<Badge variant="secondary">
								{currentQuestion.points} {currentQuestion.points === 1 ? 'punto' : 'puntos'}
							</Badge>
						</div>
					</CardHeader>
					<CardContent>
						<div class="space-y-4">
							{#if currentQuestion.type === 'single_choice'}
								<div class="text-sm text-muted-foreground mb-4">
									Selecciona una respuesta:
								</div>
								<div class="space-y-4">
									{#each currentQuestion.answer_options as option}
										<div class="flex items-center space-x-2">
											<input
												type="radio"
												name="question-{currentQuestion.id}"
												value={option.id}
												id="option-{option.id}"
												checked={answers[currentQuestion.id]?.includes(option.id) || false}
												onchange={() => selectAnswer(option.id)}
												class="h-4 w-4 text-primary focus:ring-primary border-gray-300"
											/>
											<Label
												for="option-{option.id}"
												class="text-base leading-relaxed cursor-pointer flex-1"
											>
												{option.text}
											</Label>
										</div>
									{/each}
								</div>
							{:else}
								<div class="text-sm text-muted-foreground mb-4">
									Selecciona todas las respuestas correctas:
								</div>
								<div class="space-y-4">
									{#each currentQuestion.answer_options as option}
										<div class="flex items-start space-x-3">
											<input
												type="checkbox"
												id="option-{option.id}"
												checked={answers[currentQuestion.id]?.includes(option.id) || false}
												onchange={() => selectAnswer(option.id)}
												class="mt-1 h-4 w-4 rounded border-gray-300 text-primary focus:ring-primary"
											/>
											<Label
												for="option-{option.id}"
												class="text-base leading-relaxed cursor-pointer flex-1"
											>
												{option.text}
											</Label>
										</div>
									{/each}
								</div>
							{/if}
						</div>

						<!-- Navigation -->
						<div class="flex items-center justify-between mt-8 pt-6 border-t">
							<Button
								variant="outline"
								onclick={prevQuestion}
								disabled={!canGoPrev}
							>
								<ChevronLeft class="mr-2 h-4 w-4" />
								Anterior
							</Button>

							<div class="text-sm text-muted-foreground">
								{#if hasAnsweredCurrent}
									<CheckCircle class="inline h-4 w-4 text-green-600 mr-1" />
									Respondida
								{:else}
									<AlertCircle class="inline h-4 w-4 text-orange-600 mr-1" />
									Sin responder
								{/if}
							</div>

							{#if canGoNext}
								<Button onclick={nextQuestion}>
									Siguiente
									<ChevronRight class="ml-2 h-4 w-4" />
								</Button>
							{:else}
								<Button onclick={submitAttempt} disabled={isSubmitting} class="bg-green-600 hover:bg-green-700">
									{#if isSubmitting}
										Enviando...
									{:else}
										<Send class="mr-2 h-4 w-4" />
										Enviar Evaluación
									{/if}
								</Button>
							{/if}
						</div>
					</CardContent>
				</Card>
			</div>

			<!-- Question Navigator -->
			<div class="lg:col-span-1">
				<Card class="sticky top-32">
					<CardHeader>
						<CardTitle class="text-base">Navegación</CardTitle>
						<CardDescription>
							Haz clic en cualquier número para ir a esa pregunta
						</CardDescription>
					</CardHeader>
					<CardContent>
						<div class="grid grid-cols-5 gap-2">
							{#each data.attempt.questions as question, index}
								<Button
									variant={currentQuestionIndex === index ? "default" : "outline"}
									size="sm"
									class="relative {answers[question.id]?.length > 0 ? 'ring-2 ring-green-200' : ''}"
									onclick={() => goToQuestion(index)}
								>
									{index + 1}
									{#if answers[question.id]?.length > 0}
										<CheckCircle class="absolute -top-1 -right-1 h-3 w-3 text-green-600 fill-current" />
									{/if}
								</Button>
							{/each}
						</div>

						<div class="mt-4 pt-4 border-t space-y-2 text-xs text-muted-foreground">
							<div class="flex items-center gap-2">
								<div class="w-3 h-3 border border-green-200 rounded ring-1 ring-green-200"></div>
								<span>Respondida</span>
							</div>
							<div class="flex items-center gap-2">
								<div class="w-3 h-3 border rounded bg-primary text-primary-foreground"></div>
								<span>Actual</span>
							</div>
							<div class="flex items-center gap-2">
								<div class="w-3 h-3 border rounded"></div>
								<span>Sin responder</span>
							</div>
						</div>

						{#if totalQuestions - answeredCount > 0}
							<div class="mt-4 p-3 bg-orange-50 border border-orange-200 rounded-md">
								<p class="text-sm text-orange-800">
									<AlertCircle class="inline h-4 w-4 mr-1" />
									Te faltan {totalQuestions - answeredCount} pregunta{totalQuestions - answeredCount === 1 ? '' : 's'} por responder
								</p>
							</div>
						{:else}
							<div class="mt-4 p-3 bg-green-50 border border-green-200 rounded-md">
								<p class="text-sm text-green-800">
									<CheckCircle class="inline h-4 w-4 mr-1" />
									¡Has respondido todas las preguntas!
								</p>
							</div>
						{/if}
					</CardContent>
				</Card>
			</div>
		</div>
	</div>
</div>
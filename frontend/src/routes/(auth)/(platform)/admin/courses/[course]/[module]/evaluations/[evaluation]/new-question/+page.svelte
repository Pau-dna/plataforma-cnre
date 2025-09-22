<script lang="ts">
	import type { PageProps } from './$types';
	import type { CreateQuestionDTO, CreateAnswerDTO } from '$lib/types';
	import { QuestionType } from '$lib/types';
	import { QuestionController, AnswerController } from '$lib/controllers';
	import { goto } from '$app/navigation';
	import { toast } from 'svelte-sonner';
	import { Plus, Trash2, Check } from '@lucide/svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import Input from '$lib/components/ui/input/input.svelte';
	import Label from '$lib/components/ui/label/label.svelte';
	import Textarea from '$lib/components/ui/textarea/textarea.svelte';
	import Card from '$lib/components/ui/card/card.svelte';
	import CardHeader from '$lib/components/ui/card/card-header.svelte';
	import CardTitle from '$lib/components/ui/card/card-title.svelte';
	import CardContent from '$lib/components/ui/card/card-content.svelte';
	import Back from '$lib/components/kit/Back.svelte';

	let { data }: PageProps = $props();

	const evaluation = data.evaluation;
	const courseId = data.courseId;
	const moduleId = data.moduleId;
	const questionController = new QuestionController();
	const answerController = new AnswerController();

	let formData = $state({
		text: '',
		type: QuestionType.SINGLE as QuestionType,
		explanation: '',
		points: 1
	});

	interface AnswerForm {
		text: string;
		is_correct: boolean;
		id?: number; // For temporary IDs
	}

	let answers = $state<AnswerForm[]>([
		{ text: '', is_correct: false },
		{ text: '', is_correct: false }
	]);

	let isSubmitting = $state(false);

	function addAnswer() {
		answers.push({ text: '', is_correct: false });
	}

	function removeAnswer(index: number) {
		if (answers.length <= 2) {
			toast.error('Debe tener al menos 2 respuestas');
			return;
		}
		answers = answers.filter((_, i) => i !== index);
	}

	function toggleCorrectAnswer(index: number) {
		if (formData.type === QuestionType.SINGLE) {
			// For single choice, only one answer can be correct
			answers.forEach((answer, i) => {
				answer.is_correct = i === index;
			});
		} else {
			// For multiple choice, toggle the selected answer
			answers[index].is_correct = !answers[index].is_correct;
		}
	}

	async function handleSubmit() {
		// Basic validation
		if (!formData.text.trim()) {
			toast.error('El texto de la pregunta es requerido');
			return;
		}

		if (formData.points < 1) {
			toast.error('Los puntos deben ser mayor a 0');
			return;
		}

		// Validate answers
		const validAnswers = answers.filter(a => a.text.trim());
		if (validAnswers.length < 2) {
			toast.error('Debe tener al menos 2 respuestas válidas');
			return;
		}

		const correctAnswers = validAnswers.filter(a => a.is_correct);
		if (correctAnswers.length === 0) {
			toast.error('Debe marcar al menos una respuesta como correcta');
			return;
		}

		if (formData.type === QuestionType.SINGLE && correctAnswers.length > 1) {
			toast.error('Para preguntas de opción única, solo puede haber una respuesta correcta');
			return;
		}

		isSubmitting = true;

		try {
			// Create question
			const questionData: CreateQuestionDTO = {
				text: formData.text.trim(),
				type: formData.type,
				explanation: formData.explanation.trim() || undefined,
				points: formData.points,
				evaluation_id: evaluation.id
			};

			const newQuestion = await questionController.createQuestion(questionData);

			// Create answers
			const answerPromises = validAnswers.map(async (answer, index) => {
				const answerData: CreateAnswerDTO = {
					text: answer.text.trim(),
					is_correct: answer.is_correct,
					order: index + 1, // Simple ordering, as requested to ignore but required for model
					question_id: newQuestion.id
				};
				return answerController.createAnswer(answerData);
			});

			await Promise.all(answerPromises);

			toast.success('Pregunta creada exitosamente');
			goto(`/admin/courses/${courseId}/${moduleId}/evaluations/${evaluation.id}`);
		} catch (error) {
			console.error('Error creating question:', error);
			toast.error('Error al crear la pregunta', {
				description: error instanceof Error ? error.message : String(error)
			});
		} finally {
			isSubmitting = false;
		}
	}

	function handleCancel() {
		goto(`/admin/courses/${courseId}/${moduleId}/evaluations/${evaluation.id}`);
	}
</script>

<Back href="/admin/courses/{courseId}/{moduleId}/evaluations/{evaluation.id}" />

<div class="flex flex-col gap-6">
	<div>
		<h1 class="text-h1">Crear Nueva Pregunta</h1>
		<p class="text-subtitle">Evaluación: {evaluation.title}</p>
	</div>

	<div class="space-y-6">
		<!-- Question Details -->
		<Card>
			<CardHeader>
				<CardTitle>Información de la Pregunta</CardTitle>
			</CardHeader>
			<CardContent class="space-y-4">
				<div class="space-y-2">
					<Label for="text">Pregunta *</Label>
					<Textarea
						id="text"
						bind:value={formData.text}
						placeholder="Escribe tu pregunta aquí..."
						required
						disabled={isSubmitting}
						rows={3}
					/>
				</div>

				<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
					<div class="space-y-2">
						<Label for="type">Tipo de Pregunta *</Label>
						<select
							id="type"
							bind:value={formData.type}
							required
							disabled={isSubmitting}
							class="w-full px-3 py-2 text-sm rounded-md border border-input bg-background ring-offset-background placeholder:text-muted-foreground focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
						>
							<option value={QuestionType.SINGLE}>Opción única</option>
							<option value={QuestionType.MULTIPLE}>Opción múltiple</option>
						</select>
						<p class="text-xs text-muted-foreground">
							{#if formData.type === QuestionType.SINGLE}
								Solo una respuesta puede ser correcta
							{:else}
								Múltiples respuestas pueden ser correctas
							{/if}
						</p>
					</div>

					<div class="space-y-2">
						<Label for="points">Puntos *</Label>
						<Input
							id="points"
							type="number"
							bind:value={formData.points}
							min="1"
							max="100"
							required
							disabled={isSubmitting}
						/>
					</div>
				</div>

				<div class="space-y-2">
					<Label for="explanation">Explicación (opcional)</Label>
					<Textarea
						id="explanation"
						bind:value={formData.explanation}
						placeholder="Explicación que se mostrará después de responder..."
						disabled={isSubmitting}
						rows={2}
					/>
				</div>
			</CardContent>
		</Card>

		<!-- Answers -->
		<Card>
			<CardHeader class="flex flex-row items-center justify-between space-y-0">
				<CardTitle>Respuestas</CardTitle>
				<Button type="button" variant="outline" size="sm" onclick={addAnswer} disabled={isSubmitting}>
					<Plus class="h-4 w-4" />
					<span>Agregar Respuesta</span>
				</Button>
			</CardHeader>
			<CardContent class="space-y-4">
				{#each answers as answer, index (index)}
					<div class="flex items-start gap-3 p-4 border border-gray-200 rounded-lg">
						<div class="flex items-center gap-2 pt-2">
							<input
								type="checkbox"
								checked={answer.is_correct}
								onchange={() => toggleCorrectAnswer(index)}
								disabled={isSubmitting}
								class="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500 focus:ring-2"
							/>
							{#if answer.is_correct}
								<Check class="h-4 w-4 text-green-600" />
							{/if}
						</div>
						<div class="flex-1 space-y-2">
							<Label for="answer-{index}">Respuesta {index + 1}</Label>
							<Textarea
								id="answer-{index}"
								bind:value={answer.text}
								placeholder="Escribe la respuesta aquí..."
								disabled={isSubmitting}
								rows={2}
							/>
						</div>
						{#if answers.length > 2}
							<Button
								type="button"
								variant="ghost"
								size="sm"
								onclick={() => removeAnswer(index)}
								disabled={isSubmitting}
								class="mt-6 text-red-600 hover:text-red-800 hover:bg-red-50"
							>
								<Trash2 class="h-4 w-4" />
							</Button>
						{/if}
					</div>
				{/each}

				<div class="text-sm text-muted-foreground p-3 bg-blue-50 border border-blue-200 rounded-lg">
					<p class="font-medium mb-1">Instrucciones:</p>
					<ul class="list-disc list-inside space-y-1">
						<li>Marca las respuestas correctas con el checkbox</li>
						<li>Puedes agregar tantas respuestas como necesites (mínimo 2)</li>
						{#if formData.type === QuestionType.SINGLE}
							<li class="text-orange-700">Solo una respuesta puede ser correcta para este tipo</li>
						{:else}
							<li class="text-green-700">Múltiples respuestas pueden ser correctas para este tipo</li>
						{/if}
					</ul>
				</div>
			</CardContent>
		</Card>

		<!-- Actions -->
		<div class="flex justify-end gap-4">
			<Button variant="outline" type="button" onclick={handleCancel} disabled={isSubmitting}>
				Cancelar
			</Button>
			<Button
				type="submit"
				disabled={isSubmitting}
				onclick={handleSubmit}
				class="bg-green-500 hover:bg-green-600"
			>
				{#if isSubmitting}
					Creando...
				{:else}
					Crear Pregunta
				{/if}
			</Button>
		</div>
	</form>
</div>
<script lang="ts">
	import type { PageProps } from './$types';
	import type { UpdateQuestionDTO, UpdateAnswerDTO, CreateAnswerDTO } from '$lib/types';
	import { QuestionType } from '$lib/types';
	import { QuestionController, AnswerController } from '$lib/controllers';
	import { goto } from '$app/navigation';
	import { toast } from 'svelte-sonner';
	import { Plus, Trash2, Check, Save, ArrowLeft } from '@lucide/svelte';
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

	const question = data.question;
	const courseID = data.courseId;
	const moduleID = data.moduleId;
	const evaluationID = data.evaluationId;
	const questionController = new QuestionController();
	const answerController = new AnswerController();

	let formData = $state({
		text: question.text,
		type: question.type as QuestionType,
		explanation: question.explanation || '',
		points: question.points
	});

	interface AnswerForm {
		id?: number;
		text: string;
		is_correct: boolean;
		order: number;
		isNew?: boolean;
		isDeleted?: boolean;
	}

	// Initialize answers from existing question data
	let answers = $state<AnswerForm[]>(
		question.answers?.map((answer) => ({
			id: answer.id,
			text: answer.text,
			is_correct: answer.is_correct,
			order: answer.order,
			isNew: false,
			isDeleted: false
		})) || [
			{ text: '', is_correct: false, order: 1, isNew: true },
			{ text: '', is_correct: false, order: 2, isNew: true }
		]
	);

	let isSubmitting = $state(false);

	function addAnswer() {
		const maxOrder = Math.max(...answers.map((a) => a.order), 0);
		answers.push({
			text: '',
			is_correct: false,
			order: maxOrder + 1,
			isNew: true
		});
	}

	function removeAnswer(index: number) {
		const activeAnswers = answers.filter((a) => !a.isDeleted);
		if (activeAnswers.length <= 2) {
			toast.error('Debe tener al menos 2 respuestas');
			return;
		}

		const answer = answers[index];
		if (answer.isNew) {
			// Remove new answers completely
			answers = answers.filter((_, i) => i !== index);
		} else {
			// Mark existing answers for deletion
			answer.isDeleted = true;
		}
	}

	function toggleCorrectAnswer(index: number) {
		const answer = answers[index];
		if (answer.isDeleted) return;

		if (formData.type === QuestionType.SINGLE) {
			// For single choice, only one answer can be correct
			answers.forEach((ans, i) => {
				if (!ans.isDeleted) {
					ans.is_correct = i === index;
				}
			});
		} else {
			// For multiple choice, toggle the selected answer
			answer.is_correct = !answer.is_correct;
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
		const validAnswers = answers.filter((a) => !a.isDeleted && a.text.trim());
		if (validAnswers.length < 2) {
			toast.error('Debe tener al menos 2 respuestas válidas');
			return;
		}

		const correctAnswers = validAnswers.filter((a) => a.is_correct);
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
			// Update question
			const questionData: UpdateQuestionDTO = {
				text: formData.text.trim(),
				type: formData.type,
				explanation: formData.explanation.trim() || undefined,
				points: formData.points
			};

			await questionController.updateQuestion(question.id, questionData);

			// Handle answers
			const answerPromises: Promise<unknown>[] = [];

			// Delete answers marked for deletion
			const answersToDelete = answers.filter((a) => a.isDeleted && !a.isNew && a.id);
			answerPromises.push(
				...answersToDelete.map((answer) => answerController.deleteAnswer(answer.id!))
			);

			// Update existing answers
			const answersToUpdate = answers.filter((a) => !a.isDeleted && !a.isNew && a.id);
			answerPromises.push(
				...answersToUpdate.map((answer) => {
					const answerData: UpdateAnswerDTO = {
						text: answer.text.trim(),
						is_correct: answer.is_correct,
						order: answer.order
					};
					return answerController.updateAnswer(answer.id!, answerData);
				})
			);

			// Create new answers
			const answersToCreate = answers.filter((a) => !a.isDeleted && a.isNew && a.text.trim());
			answerPromises.push(
				...answersToCreate.map((answer) => {
					const answerData: CreateAnswerDTO = {
						text: answer.text.trim(),
						is_correct: answer.is_correct,
						order: answer.order,
						question_id: question.id
					};
					return answerController.createAnswer(answerData);
				})
			);

			await Promise.all(answerPromises);

			toast.success('Pregunta actualizada exitosamente');
			// eslint-disable-next-line svelte/no-navigation-without-resolve
			goto(`/admin/courses/${courseID}/${moduleID}/evaluations/${evaluationID}`);
		} catch (error) {
			console.error('Error updating question:', error);
			toast.error('Error al actualizar la pregunta', {
				description: error instanceof Error ? error.message : String(error)
			});
		} finally {
			isSubmitting = false;
		}
	}

	function handleCancel() {
		// eslint-disable-next-line svelte/no-navigation-without-resolve
		goto(`/admin/courses/${courseID}/${moduleID}/evaluations/${evaluationID}`);
	}

	// Filter out deleted answers for display
	let visibleAnswers = $derived(answers.filter((a) => !a.isDeleted));
</script>

<Back href="/admin/courses/{courseID}/{moduleID}/evaluations/{evaluationID}" />

<div class="flex flex-col gap-6">
	<div>
		<h1 class="text-h1">Editar Pregunta</h1>
		<p class="text-subtitle">Pregunta ID: {question.id}</p>
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

				<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
					<div class="space-y-2">
						<Label for="type">Tipo de Pregunta *</Label>
						<select
							id="type"
							bind:value={formData.type}
							required
							disabled={isSubmitting}
							class="border-input bg-background ring-offset-background placeholder:text-muted-foreground focus:ring-ring w-full rounded-md border px-3 py-2 text-sm focus:ring-2 focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
						>
							<option value={QuestionType.SINGLE}>Opción única</option>
							<option value={QuestionType.MULTIPLE}>Opción múltiple</option>
						</select>
						<p class="text-muted-foreground text-xs">
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
				<Button
					type="button"
					variant="outline"
					size="sm"
					onclick={addAnswer}
					disabled={isSubmitting}
				>
					<Plus class="h-4 w-4" />
					<span>Agregar Respuesta</span>
				</Button>
			</CardHeader>
			<CardContent class="space-y-4">
				{#each visibleAnswers as answer, index (answer.id || `new-${index}`)}
					<div
						class="flex items-start gap-3 rounded-lg border border-gray-200 p-4 {answer.isNew
							? 'bg-blue-50'
							: ''}"
					>
						<div class="flex items-center gap-2 pt-2">
							<input
								type="checkbox"
								checked={answer.is_correct}
								onchange={() => toggleCorrectAnswer(answers.indexOf(answer))}
								disabled={isSubmitting}
								class="h-4 w-4 rounded border-gray-300 bg-gray-100 text-blue-600 focus:ring-2 focus:ring-blue-500"
							/>
							{#if answer.is_correct}
								<Check class="h-4 w-4 text-green-600" />
							{/if}
						</div>
						<div class="flex-1 space-y-2">
							<Label for="answer-{index}">
								Respuesta {index + 1}
								{#if answer.isNew}
									<span class="text-xs font-medium text-blue-600">(Nueva)</span>
								{/if}
							</Label>
							<Textarea
								id="answer-{index}"
								bind:value={answer.text}
								placeholder="Escribe la respuesta aquí..."
								disabled={isSubmitting}
								rows={2}
							/>
						</div>
						{#if visibleAnswers.length > 2}
							<Button
								type="button"
								variant="ghost"
								size="sm"
								onclick={() => removeAnswer(answers.indexOf(answer))}
								disabled={isSubmitting}
								class="mt-6 text-red-600 hover:bg-red-50 hover:text-red-800"
							>
								<Trash2 class="h-4 w-4" />
							</Button>
						{/if}
					</div>
				{/each}

				<div class="text-muted-foreground rounded-lg border border-blue-200 bg-blue-50 p-3 text-sm">
					<p class="mb-1 font-medium">Instrucciones:</p>
					<ul class="list-inside list-disc space-y-1">
						<li>Marca las respuestas correctas con el checkbox</li>
						<li>Puedes agregar tantas respuestas como necesites (mínimo 2)</li>
						{#if formData.type === QuestionType.SINGLE}
							<li class="text-orange-700">Solo una respuesta puede ser correcta para este tipo</li>
						{:else}
							<li class="text-green-700">
								Múltiples respuestas pueden ser correctas para este tipo
							</li>
						{/if}
						<li class="text-blue-700">Las respuestas nuevas aparecen con fondo azul claro</li>
					</ul>
				</div>
			</CardContent>
		</Card>

		<!-- Actions -->
		<div class="flex justify-end gap-4">
			<Button variant="outline" type="button" onclick={handleCancel} disabled={isSubmitting}>
				<ArrowLeft class="h-4 w-4" />
				Cancelar
			</Button>
			<Button
				type="submit"
				disabled={isSubmitting}
				onclick={handleSubmit}
				class="bg-blue-500 hover:bg-blue-600"
			>
				<Save class="h-4 w-4" />
				{#if isSubmitting}
					Actualizando...
				{:else}
					Guardar Cambios
				{/if}
			</Button>
		</div>
	</div>
</div>

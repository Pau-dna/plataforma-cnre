<script lang="ts">
	import type { PageProps } from './$types';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Textarea } from '$lib/components/ui/textarea/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import type { UpdateEvaluationDTO } from '$lib/types';
	import { EvaluationController } from '$lib/controllers';
	import { ApiError } from '$lib/utils/error';
	import { toast } from 'svelte-sonner';
	import { goto } from '$app/navigation';
	import Back from '$lib/components/kit/Back.svelte';

	let { data }: PageProps = $props();

	let submitting = $state(false);
	const formdata = $state<UpdateEvaluationDTO>({
		title: data.evaluation.title,
		description: data.evaluation.description || '',
		question_count: data.evaluation.question_count,
		answer_options_count: data.evaluation.answer_options_count,
		passing_score: data.evaluation.passing_score,
		max_attempts: data.evaluation.max_attempts || undefined,
		time_limit: data.evaluation.time_limit || undefined
	});

	const evaluationController = new EvaluationController();

	async function handleSubmit() {
		// Validate data
		if (!formdata.title || !formdata.question_count || formdata.question_count < 1 || !formdata.answer_options_count || formdata.answer_options_count < 2 || !formdata.passing_score || formdata.passing_score < 0 || formdata.passing_score > 100) {
			toast.error('Por favor, complete correctamente todos los campos obligatorios.');
			return;
		}

		try {
			submitting = true;
			await evaluationController.updateEvaluation(data.evaluationId, formdata);
			toast.success('Configuración de evaluación actualizada con éxito.');
			// Redirect back to evaluation management page
			goto(`/admin/courses/${data.courseId}/${data.moduleId}/evaluations/${data.evaluationId}`);
		} catch (error) {
			const apiError = ApiError.from(error);

			// Handle specific error types
			if (apiError.isValidationError()) {
				toast.error('Datos inválidos', {
					description: apiError.message
				});
			} else if (apiError.isAuthError()) {
				toast.error('Sin autorización', {
					description: 'No tienes permisos para editar evaluaciones.'
				});
			} else {
				toast.error('Error al actualizar la configuración', {
					description: apiError.getUserMessage()
				});
			}
		} finally {
			submitting = false;
		}
	}
</script>

<Back href="/admin/courses/{data.courseId}/{data.moduleId}/evaluations/{data.evaluationId}" />

<div class="flex items-center justify-center">
	<Card.Root class="form-card w-full max-w-2xl">
		<Card.Header class="flex flex-col gap-3">
			<Card.Title class="text-h2">Configuración de Evaluación</Card.Title>
			<Card.Description>
				Modifique la configuración de la evaluación "{data.evaluation.title}" en el módulo "{data.module.title}".
			</Card.Description>
		</Card.Header>
		<Card.Content class="flex flex-col gap-6">
			<div class="flex flex-col gap-2">
				<Label for="title">Título de la evaluación *</Label>
				<Input
					id="title"
					placeholder="Ingrese el título de la evaluación"
					bind:value={formdata.title}
				/>
			</div>

			<div class="flex flex-col gap-2">
				<Label for="description">Descripción</Label>
				<Textarea
					id="description"
					placeholder="Ingrese una descripción opcional de la evaluación"
					bind:value={formdata.description}
				/>
			</div>

			<div class="grid grid-cols-1 gap-6 sm:grid-cols-2">
				<div class="flex flex-col gap-2">
					<Label for="question_count">Número de preguntas *</Label>
					<Input
						id="question_count"
						type="number"
						min="1"
						placeholder="5"
						bind:value={formdata.question_count}
					/>
				</div>

				<div class="flex flex-col gap-2">
					<Label for="answer_options_count">Opciones por pregunta *</Label>
					<Input
						id="answer_options_count"
						type="number"
						min="2"
						max="6"
						placeholder="4"
						bind:value={formdata.answer_options_count}
					/>
				</div>
			</div>

			<div class="grid grid-cols-1 gap-6 sm:grid-cols-2">
				<div class="flex flex-col gap-2">
					<Label for="passing_score">Puntaje mínimo para aprobar (%) *</Label>
					<Input
						id="passing_score"
						type="number"
						min="0"
						max="100"
						placeholder="70"
						bind:value={formdata.passing_score}
					/>
				</div>

				<div class="flex flex-col gap-2">
					<Label for="max_attempts">Intentos máximos (opcional)</Label>
					<Input
						id="max_attempts"
						type="number"
						min="1"
						placeholder="3"
						bind:value={formdata.max_attempts}
					/>
				</div>
			</div>

			<div class="flex flex-col gap-2">
				<Label for="time_limit">Tiempo límite en minutos (opcional)</Label>
				<Input
					id="time_limit"
					type="number"
					min="1"
					placeholder="30"
					bind:value={formdata.time_limit}
				/>
			</div>
		</Card.Content>
		<Card.Footer>
			<Button
				onclick={handleSubmit}
				disabled={submitting}
				class="w-full bg-pink-500 hover:bg-pink-900"
			>
				{submitting ? 'Guardando...' : 'Guardar configuración'}
			</Button>
		</Card.Footer>
	</Card.Root>
</div>
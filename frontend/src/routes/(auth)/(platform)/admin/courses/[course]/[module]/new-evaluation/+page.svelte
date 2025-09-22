<script lang="ts">
	import type { PageProps } from './$types';
	import type { CreateEvaluationDTO } from '$lib/types';
	import { ContentType } from '$lib/types';
	import { EvaluationController } from '$lib/controllers';
	import { goto } from '$app/navigation';
	import { toast } from 'svelte-sonner';
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

	const module = data.module;
	const courseId = data.courseId;
	const moduleId = data.moduleId;
	const evaluationController = new EvaluationController();

	let formData = $state({
		title: '',
		description: '',
		question_count: 1,
		passing_score: 60,
		max_attempts: undefined as number | undefined,
		time_limit: undefined as number | undefined
	});

	let isSubmitting = $state(false);

	async function handleSubmit() {
		// Basic validation
		if (!formData.title.trim()) {
			toast.error('El título es requerido');
			return;
		}

		if (formData.question_count < 1) {
			toast.error('Debe tener al menos 1 pregunta');
			return;
		}

		if (formData.passing_score < 0 || formData.passing_score > 100) {
			toast.error('El puntaje mínimo debe estar entre 0 y 100');
			return;
		}

		isSubmitting = true;

		try {
			// Get the next order number (simplified, could be improved)
			const nextOrder = Date.now() % 1000000; // Simple ordering

			const evaluationData: CreateEvaluationDTO = {
				title: formData.title.trim(),
				description: formData.description.trim() || undefined,
				type: ContentType.EVALUATION,
				question_count: formData.question_count,
				passing_score: formData.passing_score,
				max_attempts: formData.max_attempts || undefined,
				time_limit: formData.time_limit || undefined,
				module_id: moduleId,
				order: nextOrder
			};

			const newEvaluation = await evaluationController.createEvaluation(evaluationData);
			toast.success('Evaluación creada exitosamente');
			goto(`/admin/courses/${courseId}/${moduleId}`);
		} catch (error) {
			console.error('Error creating evaluation:', error);
			toast.error('Error al crear la evaluación', {
				description: error instanceof Error ? error.message : String(error)
			});
		} finally {
			isSubmitting = false;
		}
	}

	function handleCancel() {
		goto(`/admin/courses/${courseId}/${moduleId}`);
	}
</script>

<Back href="/admin/courses/{courseId}/{moduleId}" />

<div class="flex flex-col gap-6">
	<div>
		<h1 class="text-h1">Crear Nueva Evaluación</h1>
		<p class="text-subtitle">Módulo: {module.title}</p>
	</div>

	<Card class="w-full max-w-4xl">
		<CardHeader>
			<CardTitle>Información de la Evaluación</CardTitle>
		</CardHeader>
		<CardContent>
			<form onsubmit|preventDefault={handleSubmit} class="space-y-6">
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
					<div class="space-y-2 md:col-span-2">
						<Label for="title">Título *</Label>
						<Input
							id="title"
							bind:value={formData.title}
							placeholder="Título de la evaluación"
							required
							disabled={isSubmitting}
						/>
					</div>

					<div class="space-y-2 md:col-span-2">
						<Label for="description">Descripción</Label>
						<Textarea
							id="description"
							bind:value={formData.description}
							placeholder="Descripción opcional de la evaluación"
							disabled={isSubmitting}
							rows={3}
						/>
					</div>

					<div class="space-y-2">
						<Label for="question_count">Número de Preguntas *</Label>
						<Input
							id="question_count"
							type="number"
							bind:value={formData.question_count}
							min="1"
							max="100"
							required
							disabled={isSubmitting}
						/>
					</div>

					<div class="space-y-2">
						<Label for="passing_score">Puntaje Mínimo (%) *</Label>
						<Input
							id="passing_score"
							type="number"
							bind:value={formData.passing_score}
							min="0"
							max="100"
							required
							disabled={isSubmitting}
						/>
					</div>

					<div class="space-y-2">
						<Label for="max_attempts">Intentos Máximos</Label>
						<Input
							id="max_attempts"
							type="number"
							bind:value={formData.max_attempts}
							min="1"
							placeholder="Sin límite"
							disabled={isSubmitting}
						/>
						<p class="text-sm text-muted-foreground">
							Dejar vacío para permitir intentos ilimitados
						</p>
					</div>

					<div class="space-y-2">
						<Label for="time_limit">Tiempo Límite (minutos)</Label>
						<Input
							id="time_limit"
							type="number"
							bind:value={formData.time_limit}
							min="1"
							placeholder="Sin límite de tiempo"
							disabled={isSubmitting}
						/>
						<p class="text-sm text-muted-foreground">
							Dejar vacío para permitir tiempo ilimitado
						</p>
					</div>
				</div>

				<div class="flex justify-end gap-4">
					<Button variant="outline" type="button" onclick={handleCancel} disabled={isSubmitting}>
						Cancelar
					</Button>
					<Button
						type="submit"
						disabled={isSubmitting}
						class="bg-purple-500 hover:bg-purple-600"
					>
						{#if isSubmitting}
							Creando...
						{:else}
							Crear Evaluación
						{/if}
					</Button>
				</div>
			</form>
		</CardContent>
	</Card>
</div>
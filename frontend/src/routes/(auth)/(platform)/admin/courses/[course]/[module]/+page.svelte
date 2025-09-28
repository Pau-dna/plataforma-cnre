<script lang="ts">
	import type { PageProps } from './$types';
	import type { Content, Evaluation, ReorderItemDTO } from '$lib/types';
	import { Plus } from '@lucide/svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import Back from '$lib/components/kit/Back.svelte';
	import { ContentController, EvaluationController } from '$lib/controllers';
	import ContentCard from '$lib/components/content/ContentCard.svelte';
	import EvaluationCard from '$lib/components/evaluation/EvaluationCard.svelte';
	import { toast } from 'svelte-sonner';

	let { data }: PageProps = $props();

	const module = $state(data.module);
	let contents = $state(data.contents);
	let evaluations = $state(data.evaluations);
	const courseId = data.courseId;
	const contentController = new ContentController();
	const evaluationController = new EvaluationController();

	function handleContentUpdate(updated: Content) {
		// Find and update the content in the array
		const index = contents.findIndex((c) => c.id === updated.id);
		if (index !== -1) {
			contents[index] = updated;
		}
	}

	function handleContentDelete(deleted: Content) {
		// Remove the content from the array
		contents = contents.filter((c) => c.id !== deleted.id);
	}

	function handleEvaluationUpdate(updated: Evaluation) {
		// Find and update the evaluation in the array
		const index = evaluations.findIndex((e) => e.id === updated.id);
		if (index !== -1) {
			evaluations[index] = updated;
		}
	}

	function handleEvaluationDelete(deleted: Evaluation) {
		// Remove the evaluation from the array
		evaluations = evaluations.filter((e) => e.id !== deleted.id);
	}

	async function handleMoveUp(content: Content) {
		const currentIndex = contents.findIndex((c) => c.id === content.id);
		if (currentIndex <= 0) return; // Can't move up if it's the first content

		const previousContent = contents[currentIndex - 1];

		try {
			// Create reorder data - swap the orders
			await contentController.reorderContent(module.id, [
				{ id: content.id, order: previousContent.order },
				{ id: previousContent.id, order: content.order }
			]);

			// Update local state - swap the contents and their orders
			const updatedContent = { ...content, order: previousContent.order };
			const updatedPreviousContent = { ...previousContent, order: content.order };

			contents[currentIndex] = updatedPreviousContent;
			contents[currentIndex - 1] = updatedContent;

			// Re-sort contents by order to ensure consistency
			contents = contents.toSorted((a, b) => a.order - b.order);

			toast.success('Contenido movido hacia arriba correctamente.');
		} catch (error) {
			console.error('Error moving content up:', error);
			toast.error('Error al mover el contenido.', {
				description: error instanceof Error ? error.message : String(error)
			});
		}
	}

	async function handleMoveDown(content: Content) {
		const currentIndex = contents.findIndex((c) => c.id === content.id);
		if (currentIndex >= contents.length - 1) return; // Can't move down if it's the last content

		const nextContent = contents[currentIndex + 1];

		try {
			// Create reorder data - swap the orders
			await contentController.reorderContent(module.id, [
				{ id: content.id, order: nextContent.order },
				{ id: nextContent.id, order: content.order }
			]);

			// Update local state - swap the contents and their orders
			const updatedContent = { ...content, order: nextContent.order };
			const updatedNextContent = { ...nextContent, order: content.order };

			contents[currentIndex] = updatedNextContent;
			contents[currentIndex + 1] = updatedContent;

			// Re-sort contents by order to ensure consistency
			contents = contents.toSorted((a, b) => a.order - b.order);

			toast.success('Contenido movido hacia abajo correctamente.');
		} catch (error) {
			console.error('Error moving content down:', error);
			toast.error('Error al mover el contenido.', {
				description: error instanceof Error ? error.message : String(error)
			});
		}
	}

	async function handleEvaluationMoveUp(evaluation: Evaluation) {
		const currentIndex = evaluations.findIndex((e) => e.id === evaluation.id);
		if (currentIndex <= 0) return; // Can't move up if it's the first evaluation

		const previousEvaluation = evaluations[currentIndex - 1];

		try {
			// Update the orders
			await evaluationController.updateEvaluationPatch(evaluation.id, {
				order: previousEvaluation.order
			});
			await evaluationController.updateEvaluationPatch(previousEvaluation.id, {
				order: evaluation.order
			});

			// Update local state - swap the evaluations and their orders
			const updatedEvaluation = { ...evaluation, order: previousEvaluation.order };
			const updatedPreviousEvaluation = { ...previousEvaluation, order: evaluation.order };

			evaluations[currentIndex] = updatedPreviousEvaluation;
			evaluations[currentIndex - 1] = updatedEvaluation;

			// Re-sort evaluations by order to ensure consistency
			evaluations = evaluations.toSorted((a, b) => a.order - b.order);

			toast.success('Evaluación movida hacia arriba correctamente.');
		} catch (error) {
			console.error('Error moving evaluation up:', error);
			toast.error('Error al mover la evaluación.', {
				description: error instanceof Error ? error.message : String(error)
			});
		}
	}

	async function handleEvaluationMoveDown(evaluation: Evaluation) {
		const currentIndex = evaluations.findIndex((e) => e.id === evaluation.id);
		if (currentIndex >= evaluations.length - 1) return; // Can't move down if it's the last evaluation

		const nextEvaluation = evaluations[currentIndex + 1];

		try {
			// Update the orders
			await evaluationController.updateEvaluationPatch(evaluation.id, {
				order: nextEvaluation.order
			});
			await evaluationController.updateEvaluationPatch(nextEvaluation.id, {
				order: evaluation.order
			});

			// Update local state - swap the evaluations and their orders
			const updatedEvaluation = { ...evaluation, order: nextEvaluation.order };
			const updatedNextEvaluation = { ...nextEvaluation, order: evaluation.order };

			evaluations[currentIndex] = updatedNextEvaluation;
			evaluations[currentIndex + 1] = updatedEvaluation;

			// Re-sort evaluations by order to ensure consistency
			evaluations = evaluations.toSorted((a, b) => a.order - b.order);

			toast.success('Evaluación movida hacia abajo correctamente.');
		} catch (error) {
			console.error('Error moving evaluation down:', error);
			toast.error('Error al mover la evaluación.', {
				description: error instanceof Error ? error.message : String(error)
			});
		}
	}
</script>

<Back href="/admin/courses/{courseId}" />

<div class="flex flex-col gap-6">
	<div class="flex flex-col gap-2">
		<div class="flex items-center justify-between">
			<h1 class="text-h1">Contenidos del Módulo</h1>
			<div class="flex gap-2">
				<Button
					href="/admin/courses/{courseId}/{module.id}/new-content"
					class="bg-pink-500 hover:bg-pink-900"
				>
					<Plus class="h-4 w-4 leading-none" />
					<span class="leading-none">Crear Contenido</span>
				</Button>
				<Button
					href="/admin/courses/{courseId}/{module.id}/new-evaluation"
					class="bg-purple-500 hover:bg-purple-900"
				>
					<Plus class="h-4 w-4 leading-none" />
					<span class="leading-none">Crear Evaluación</span>
				</Button>
			</div>
		</div>
		<p class="text-subtitle">{module.title}</p>
		<p class="text-muted-foreground text-sm">{module.description}</p>
	</div>

	<!-- Contenidos -->
	<div class="flex flex-col gap-4">
		<h2 class="text-xl font-semibold">Contenidos</h2>
		<div class="grid grid-cols-1 gap-4">
			{#each contents as content, index (content.id)}
				<ContentCard
					{content}
					{courseId}
					actDate={content.updated_at}
					onupdate={handleContentUpdate}
					ondelete={handleContentDelete}
					onmoveup={handleMoveUp}
					onmovedown={handleMoveDown}
					canMoveUp={index > 0}
					canMoveDown={index < contents.length - 1}
				/>
			{/each}

			{#if contents.length === 0}
				<div class="flex flex-col items-center justify-center py-12 text-center">
					<p class="text-muted-foreground mb-4 text-lg">Este módulo no tiene contenidos aún</p>
					<Button
						href="/admin/courses/{courseId}/{module.id}/new-content"
						class="bg-pink-500 hover:bg-pink-900"
					>
						<Plus class="h-4 w-4 leading-none" />
						<span class="leading-none">Crear Primer Contenido</span>
					</Button>
				</div>
			{/if}
		</div>
	</div>

	<!-- Evaluaciones -->
	<div class="flex flex-col gap-4">
		<h2 class="text-xl font-semibold">Evaluaciones</h2>
		<div class="grid grid-cols-1 gap-4">
			{#each evaluations as evaluation, index (evaluation.id)}
				<EvaluationCard
					{evaluation}
					{courseId}
					onupdate={handleEvaluationUpdate}
					ondelete={handleEvaluationDelete}
					onmoveup={handleEvaluationMoveUp}
					onmovedown={handleEvaluationMoveDown}
					canMoveUp={index > 0}
					canMoveDown={index < evaluations.length - 1}
				/>
			{/each}

			{#if evaluations.length === 0}
				<div class="flex flex-col items-center justify-center py-12 text-center">
					<p class="text-muted-foreground mb-4 text-lg">Este módulo no tiene evaluaciones aún</p>
					<Button
						href="/admin/courses/{courseId}/{module.id}/new-evaluation"
						class="bg-purple-500 hover:bg-purple-900"
					>
						<Plus class="h-4 w-4 leading-none" />
						<span class="leading-none">Crear Primera Evaluación</span>
					</Button>
				</div>
			{/if}
		</div>
	</div>
</div>

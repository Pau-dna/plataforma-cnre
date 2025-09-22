<script lang="ts">
	import type { PageProps } from './$types';
	import type { Content, ReorderItemDTO } from '$lib/types';
	import { Plus } from '@lucide/svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import Back from '$lib/components/kit/Back.svelte';
	import { ContentController } from '$lib/controllers';
	import ContentCard from '$lib/components/content/ContentCard.svelte';
	import { toast } from 'svelte-sonner';

	let { data }: PageProps = $props();

	const module = $state(data.module);
	let contents = $state(data.contents);
	const courseId = data.courseId;
	const contentController = new ContentController();

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
</script>

<Back href="/admin/courses/{courseId}" />

<div class="flex flex-col gap-4">
	<div class="flex flex-col gap-2">
		<div class="flex items-center justify-between">
			<h1 class="text-h1">Contenidos del Módulo</h1>
			<Button
				href="/admin/courses/{courseId}/{module.id}/new-content"
				class="bg-pink-500 hover:bg-pink-900"
			>
				<Plus class="h-4 w-4 leading-none" />
				<span class="leading-none">Crear Contenido</span>
			</Button>
		</div>
		<p class="text-subtitle">{module.title}</p>
		<p class="text-muted-foreground text-sm">{module.description}</p>
	</div>

	<div class="grid grid-cols-1 gap-4">
		{#each contents as content, index (content.id)}
			<ContentCard
				{content}
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

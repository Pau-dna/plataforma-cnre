<script lang="ts">
	import ModuleCard from '$lib/components/module/ModuleCard.svelte';
	import type { PageProps } from './$types';
	import type { Module, ReorderItemDTO } from '$lib/types';
	import { Plus } from '@lucide/svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import Back from '$lib/components/kit/Back.svelte';
	import { ModuleController } from '$lib/controllers/module';

	let { data }: PageProps = $props();

	const course = $state(data.course);
	const modules = $state(data.modules);
	const moduleController = new ModuleController();

	function handleModuleUpdate(updated: Module) {
		// Find and update the module in the array
		const index = modules.findIndex(m => m.id === updated.id);
		if (index !== -1) {
			modules[index] = updated;
		}
	}

	async function handleMoveUp(module: Module) {
		const currentIndex = modules.findIndex(m => m.id === module.id);
		if (currentIndex <= 0) return; // Can't move up if it's the first module
		
		const previousModule = modules[currentIndex - 1];
		
		// Create reorder data - swap the orders
		const reorderData: ReorderItemDTO[] = [
			{ id: module.id, order: previousModule.order },
			{ id: previousModule.id, order: module.order }
		];
		
		try {
			await moduleController.reorderModules(course.id, reorderData);
			
			// Update local state - swap the modules and their orders
			const updatedModule = { ...module, order: previousModule.order };
			const updatedPreviousModule = { ...previousModule, order: module.order };
			
			modules[currentIndex] = updatedModule;
			modules[currentIndex - 1] = updatedPreviousModule;
			
			// Re-sort modules by order to ensure consistency
			modules.sort((a, b) => a.order - b.order);
		} catch (error) {
			console.error('Error moving module up:', error);
		}
	}

	async function handleMoveDown(module: Module) {
		const currentIndex = modules.findIndex(m => m.id === module.id);
		if (currentIndex >= modules.length - 1) return; // Can't move down if it's the last module
		
		const nextModule = modules[currentIndex + 1];
		
		// Create reorder data - swap the orders
		const reorderData: ReorderItemDTO[] = [
			{ id: module.id, order: nextModule.order },
			{ id: nextModule.id, order: module.order }
		];
		
		try {
			await moduleController.reorderModules(course.id, reorderData);
			
			// Update local state - swap the modules and their orders
			const updatedModule = { ...module, order: nextModule.order };
			const updatedNextModule = { ...nextModule, order: module.order };
			
			modules[currentIndex] = updatedModule;
			modules[currentIndex + 1] = updatedNextModule;
			
			// Re-sort modules by order to ensure consistency
			modules.sort((a, b) => a.order - b.order);
		} catch (error) {
			console.error('Error moving module down:', error);
		}
	}
</script>

<Back href="/admin/courses" />

<div class="flex flex-col gap-4">
	<div class="flex flex-col gap-2">
		<div class="flex items-center justify-between">
			<h1 class="text-h1">Módulos del Curso</h1>
			<Button href="/admin/courses/{course.id}/new-module" class="bg-pink-500 hover:bg-pink-900">
				<Plus class="h-4 w-4 leading-none" />
				<span class="leading-none">Crear Módulo</span>
			</Button>
		</div>
		<p class="text-subtitle">Inducción al CNRE</p>
	</div>

	<div class="grid grid-cols-1 gap-4">
		{#each modules as modulo, index}
			<ModuleCard
				module={modulo}
				actDate={modulo.updated_at}
				onupdate={handleModuleUpdate}
				onmoveup={handleMoveUp}
				onmovedown={handleMoveDown}
				canMoveUp={index > 0}
				canMoveDown={index < modules.length - 1}
			/>
		{/each}
	</div>
</div>

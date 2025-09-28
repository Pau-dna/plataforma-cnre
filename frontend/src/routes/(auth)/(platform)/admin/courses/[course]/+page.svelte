<script lang="ts">
	import ModuleCard from '$lib/components/module/ModuleCard.svelte';
	import type { PageProps } from './$types';
	import type { Module, ReorderItemDTO } from '$lib/types';
	import { Plus } from '@lucide/svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import Back from '$lib/components/kit/Back.svelte';
	import { ModuleController } from '$lib/controllers/module';
	import { toast } from 'svelte-sonner';

	let { data }: PageProps = $props();

	const course = $state(data.course);
	let modules = $state(data.modules);
	const moduleController = new ModuleController();

	function handleModuleUpdate(updated: Module) {
		// Find and update the module in the array
		const index = modules.findIndex((m) => m.id === updated.id);
		if (index !== -1) {
			modules[index] = updated;
		}
	}

	function handleModuleDelete(deleted: Module) {
		// Remove the module from the array
		modules = modules.filter((m) => m.id !== deleted.id);
	}

	async function handleMoveUp(modulo: Module) {
		const currentIndex = modules.findIndex((m) => m.id === modulo.id);
		if (currentIndex <= 0) return; // Can't move up if it's the first module

		const previousModule = modules[currentIndex - 1];

		try {
			// Create reorder data - swap the orders
			moduleController.reorderModules(course.id, [
				{ id: modulo.id, order: previousModule.order },
				{ id: previousModule.id, order: modulo.order }
			]);

			// Update local state - swap the modules and their orders
			const updatedModule = { ...modulo, order: previousModule.order };
			const updatedPreviousModule = { ...previousModule, order: modulo.order };

			modules[currentIndex] = updatedPreviousModule;
			modules[currentIndex - 1] = updatedModule;

			// Re-sort modules by order to ensure consistency
			modules = modules.toSorted((a, b) => a.order - b.order);

			toast.success('Módulo movido hacia arriba correctamente.');
		} catch (error) {
			console.error('Error al mover módulo hacia arriba:', error);
			toast.error('Error al mover el módulo.', {
				description: error instanceof Error ? error.message : String(error)
			});
		}
	}

	async function handleMoveDown(modulo: Module) {
		const currentIndex = modules.findIndex((m) => m.id === modulo.id);
		if (currentIndex >= modules.length - 1) return; // Can't move down if it's the last module

		const nextModule = modules[currentIndex + 1];

		try {
			// Create reorder data - swap the orders
			moduleController.reorderModules(course.id, [
				{ id: modulo.id, order: nextModule.order },
				{ id: nextModule.id, order: modulo.order }
			]);

			// Update local state - swap the modules and their orders
			const updatedModule = { ...modulo, order: nextModule.order };
			const updatedNextModule = { ...nextModule, order: modulo.order };

			modules[currentIndex] = updatedNextModule;
			modules[currentIndex + 1] = updatedModule;

			// Re-sort modules by order to ensure consistency
			modules = modules.toSorted((a, b) => a.order - b.order);

			toast.success('Módulo movido hacia abajo correctamente.');
		} catch (error) {
			console.error('Error al mover módulo hacia abajo:', error);
			toast.error('Error al mover el módulo.', {
				description: error instanceof Error ? error.message : String(error)
			});
		}
	}

	$inspect(modules);
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
		{#each modules as modulo, index (modulo.id)}
			<ModuleCard
				module={modulo}
				actDate={modulo.updated_at}
				onupdate={handleModuleUpdate}
				ondelete={handleModuleDelete}
				onmoveup={handleMoveUp}
				onmovedown={handleMoveDown}
				canMoveUp={index > 0}
				canMoveDown={index < modules.length - 1}
			/>
		{/each}

		{#if modules.length === 0}
			<div class="flex flex-col items-center justify-center py-12 text-center">
				<p class="text-muted-foreground mb-4 text-lg">Este curso no tiene módulos aún</p>
				<Button href="/admin/courses/{course.id}/new-module" class="bg-pink-500 hover:bg-pink-900">
					<Plus class="h-4 w-4 leading-none" />
					<span class="leading-none">Crear Primer Módulo</span>
				</Button>
			</div>
		{/if}
	</div>
</div>

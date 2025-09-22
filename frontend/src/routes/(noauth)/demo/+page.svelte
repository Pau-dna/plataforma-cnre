<script lang="ts">
	import ModuleCard from '$lib/components/module/ModuleCard.svelte';
	import type { Module } from '$lib/types';

	// Mock data for demonstration
	let modules: Module[] = $state([
		{
			id: 1,
			title: "Módulo 1: Introducción",
			description: "Introducción al sistema CNRE",
			order: 0,
			course_id: 1,
			created_at: "2024-01-01",
			updated_at: "2024-01-01"
		},
		{
			id: 2,
			title: "Módulo 2: Fundamentos",
			description: "Fundamentos básicos de la plataforma",
			order: 1,
			course_id: 1,
			created_at: "2024-01-02",
			updated_at: "2024-01-02"
		},
		{
			id: 3,
			title: "Módulo 3: Avanzado",
			description: "Temas avanzados del sistema",
			order: 2,
			course_id: 1,
			created_at: "2024-01-03",
			updated_at: "2024-01-03"
		}
	]);

	function handleModuleUpdate(updated: Module) {
		const index = modules.findIndex(m => m.id === updated.id);
		if (index !== -1) {
			modules[index] = updated;
		}
	}

	function handleMoveUp(module: Module) {
		console.log('Move up:', module.title);
		const currentIndex = modules.findIndex(m => m.id === module.id);
		if (currentIndex <= 0) return;
		
		const previousModule = modules[currentIndex - 1];
		
		// Simulate API call success - swap the orders
		const updatedModule = { ...module, order: previousModule.order };
		const updatedPreviousModule = { ...previousModule, order: module.order };
		
		modules[currentIndex] = updatedModule;
		modules[currentIndex - 1] = updatedPreviousModule;
		
		// Re-sort modules by order to ensure consistency
		modules.sort((a, b) => a.order - b.order);
		
		console.log('Modules after move up:', modules.map(m => `${m.title} (order: ${m.order})`));
	}

	function handleMoveDown(module: Module) {
		console.log('Move down:', module.title);
		const currentIndex = modules.findIndex(m => m.id === module.id);
		if (currentIndex >= modules.length - 1) return;
		
		const nextModule = modules[currentIndex + 1];
		
		// Simulate API call success - swap the orders
		const updatedModule = { ...module, order: nextModule.order };
		const updatedNextModule = { ...nextModule, order: module.order };
		
		modules[currentIndex] = updatedModule;
		modules[currentIndex + 1] = updatedNextModule;
		
		// Re-sort modules by order to ensure consistency
		modules.sort((a, b) => a.order - b.order);
		
		console.log('Modules after move down:', modules.map(m => `${m.title} (order: ${m.order})`));
	}
</script>

<div class="p-8 max-w-4xl mx-auto">
	<h1 class="text-3xl font-bold mb-6">Demo: Reordenamiento de Módulos</h1>
	<p class="mb-4 text-gray-600">
		Use las flechas para cambiar el orden de los módulos. Los cambios se reflejan inmediatamente en la interfaz.
		(En la versión completa, los cambios se envían al backend mediante API)
	</p>
	
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
	
	<div class="mt-8 p-4 bg-gray-100 rounded">
		<h3 class="font-semibold mb-2">Estado actual del orden:</h3>
		<ul class="list-decimal list-inside">
			{#each modules as module}
				<li>{module.title} (Order: {module.order})</li>
			{/each}
		</ul>
	</div>
</div>
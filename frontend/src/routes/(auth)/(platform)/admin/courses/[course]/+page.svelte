<script lang="ts">
	import ModuleCard from '$lib/components/module/ModuleCard.svelte';
	import type { PageProps } from './$types';
	import type { Module } from '$lib/types';
	import { Plus } from '@lucide/svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import Back from '$lib/components/kit/Back.svelte';

	let { data }: PageProps = $props();

	const course = $state(data.course);
	const modules = $state(data.modules);

	function handleModuleUpdate(updated: Module) {
		// Find and update the module in the array
		const index = modules.findIndex(m => m.id === updated.id);
		if (index !== -1) {
			modules[index] = updated;
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
		{#each modules as modulo}
			<ModuleCard
				module={modulo}
				actDate={modulo.updated_at}
				onupdate={handleModuleUpdate}
			/>
		{/each}
	</div>
</div>

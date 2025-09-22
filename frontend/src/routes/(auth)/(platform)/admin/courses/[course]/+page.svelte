<script lang="ts">
	import ModuleCard from '$lib/components/module/ModuleCard.svelte';
	import type { PageProps } from './$types';
	import { Plus } from '@lucide/svelte';
	import Button from '$lib/components/ui/button/button.svelte';

	let { data }: PageProps = $props();

	const modules = [
		{
			id: 1,
			title: 'Introducción al CNRE',
			description: 'Descripción del módulo 1',
			actDate: '07/09/2025'
		},
		{
			id: 2,
			title: 'Estructura del CNRE',
			description: 'Descripción del módulo 2',
			actDate: '07/09/2025'
		},
		{
			id: 3,
			title: 'Funciones y Responsabilidades',
			description: 'Descripción del módulo 3',
			actDate: '08/09/2025'
		}
	];
</script>

<div class="flex flex-col gap-4">
	<div class="flex flex-col gap-2">
		<div class="flex items-center justify-between">
			<h1 class="text-h1">Módulos del Curso</h1>
			<Button href="/admin/courses/create" class="bg-pink-500 hover:bg-pink-900">
				<Plus class="h-4 w-4 leading-none" />
				<span class="leading-none">Crear Módulo</span>
			</Button>
		</div>
		<p class="text-subtitle">Inducción al CNRE</p>
	</div>

	<div class="grid grid-cols-1 gap-4">
		{#each modules as module}
			<ModuleCard
				id={module.id}
				title={module.title}
				description={module.description}
				actDate={module.actDate}
			/>
		{/each}
	</div>
</div>

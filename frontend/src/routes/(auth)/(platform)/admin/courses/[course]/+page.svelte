<script lang="ts">
	import ModuleCard from '$lib/components/module/ModuleCard.svelte';
	import type { PageProps } from './$types';
	import { Plus } from '@lucide/svelte';
	import Button from '$lib/components/ui/button/button.svelte';

	let { data }: PageProps = $props();

	const course = $state(data.course);
	const modules = $state(data.modules);
</script>

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
			/>
		{/each}
	</div>
</div>

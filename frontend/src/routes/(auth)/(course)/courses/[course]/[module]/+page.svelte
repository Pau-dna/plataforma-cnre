<script lang="ts">
	import type { Module } from '$lib';
	import Content from '$lib/components/course/Content.svelte';
	import StudentEvaluationCard from '$lib/components/evaluation/StudentEvaluationCard.svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	const modulo = $derived(data.modules.find((m) => m.id === data.moduleID) as Module);
	const contentsCompleted = $derived<Record<string, boolean>>(
		Object.fromEntries(data.progress.map((p) => [p.id.toString(), p.completed]))
	);
</script>

<div class="flex flex-col gap-y-4 px-4 md:px-0">
	<div class="flex flex-col gap-y-2">
		<h1 class="text-2xl font-bold md:text-3xl">{modulo.title}</h1>
		<p class="text-muted-foreground">{modulo.description}</p>
	</div>

	<!-- Contents Section -->
	<div class="flex flex-col gap-y-4">
		<h2 class="text-muted-foreground font-medium">Contenidos</h2>

		{#if data.contents.length === 0}
			<div class="flex flex-col items-center justify-center py-12 text-center">
				<p class="text-muted-foreground text-lg">Este módulo no tiene contenidos disponibles aún</p>
			</div>
		{:else}
			<div class="flex flex-col gap-y-2">
				{#each data.contents as content (content.id)}
					<Content completed={contentsCompleted[content.id.toString()]} {content} active={false} />
				{/each}
			</div>
		{/if}
	</div>

	<!-- Evaluations Section -->
	<div class="flex flex-col gap-y-4">
		<h2 class="text-muted-foreground font-medium">Evaluaciones</h2>

		{#if data.evaluations.length === 0}
			<div class="flex flex-col items-center justify-center py-12 text-center">
				<p class="text-muted-foreground text-lg">
					Este módulo no tiene evaluaciones disponibles aún
				</p>
			</div>
		{:else}
			<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
				{#each data.evaluations as evaluation (evaluation.id)}
					<StudentEvaluationCard {evaluation} courseId={data.course.id} />
				{/each}
			</div>
		{/if}
	</div>
</div>

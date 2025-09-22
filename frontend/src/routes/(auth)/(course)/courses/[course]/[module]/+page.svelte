<script lang="ts">
	import Content from '$lib/components/course/Content.svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	const modulo = $derived(data.modules.find((m) => m.id === data.moduleID));
</script>

<div class="flex flex-col gap-y-8">
	<div class="flex flex-col gap-y-2">
		<h1 class="text-3xl font-bold">{modulo.title}</h1>
		<p class="text-muted-foreground">{modulo.description}</p>
	</div>

	<div class="flex flex-col gap-y-4">
		<h2 class="text-muted-foreground font-medium">Contenidos</h2>

		{#if data.contents.length === 0}
			<div class="flex flex-col items-center justify-center py-12 text-center">
				<p class="text-muted-foreground text-lg">Este módulo no tiene contenidos disponibles aún</p>
			</div>
		{:else}
			<div class="flex flex-col gap-y-2">
				{#each data.contents as content (content.id)}
					<Content {content} active={false} />
				{/each}
			</div>
		{/if}
	</div>
</div>

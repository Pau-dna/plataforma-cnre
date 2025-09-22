<script lang="ts">
	import Content from '$lib/components/course/Content.svelte';
	import EmptyState from '$lib/components/ui/empty-state';
	import { BookOpen } from '@lucide/svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	const modulo = $derived(data.modules.find((m) => m.id === data.moduleID));
</script>

<div class="flex flex-col gap-y-8">
	{#if modulo}
		<div class="flex flex-col gap-y-2">
			<h1 class="text-3xl font-bold">{modulo.title}</h1>
			<p class="text-muted-foreground">{modulo.description}</p>
		</div>

		<div class="flex flex-col gap-y-4">
			<h2 class="text-muted-foreground font-medium">Contenidos</h2>

			{#if data.contents && data.contents.length > 0}
				<div class="flex flex-col gap-y-2">
					{#each data.contents as content (content.id)}
						<Content {content} active={false} />
					{/each}
				</div>
			{:else}
				<EmptyState
					icon={BookOpen}
					title="No hay contenidos disponibles"
					description="Este módulo aún no tiene contenidos publicados. Los contenidos aparecerán aquí una vez que sean agregados por el instructor."
				/>
			{/if}
		</div>
	{/if}
</div>

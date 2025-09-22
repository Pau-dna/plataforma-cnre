<script lang="ts">
	import Content from '$lib/components/course/Content.svelte';
	import { BookOpen } from '@lucide/svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	const modulo = $derived(data.modules.find((m) => m.id === data.moduleID));
</script>

{#if modulo}
	<div class="flex flex-col gap-y-8">
		<div class="flex flex-col gap-y-2">
			<h1 class="text-3xl font-bold">{modulo.title}</h1>
			<p class="text-muted-foreground">{modulo.description}</p>
		</div>

		<div class="flex flex-col gap-y-4">
			<h2 class="text-muted-foreground font-medium">Contenidos</h2>

			{#if data.contents.length > 0}
				<div class="flex flex-col gap-y-2">
					{#each data.contents as content (content.id)}
						<Content {content} active={false} />
					{/each}
				</div>
			{:else}
				<div
					class="flex flex-col items-center justify-center rounded-lg border border-dashed border-gray-300 p-12 text-center"
				>
					<BookOpen class="text-muted-foreground mb-4 size-12" />
					<h3 class="mb-2 text-lg font-medium text-gray-900">
						Este módulo aún no tiene contenidos
					</h3>
					<p class="text-muted-foreground">Los contenidos se agregarán próximamente.</p>
				</div>
			{/if}
		</div>
	</div>
{/if}

<script lang="ts">
	import { CheckCircle } from '@lucide/svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
</script>

<div class="flex flex-col gap-y-4 px-4 md:px-0">
	<h1 class="text-h1">{data.course.title}</h1>
	<p class="text-muted-foreground max-w-prose text-pretty">{data.course.description}</p>

	<div class="grid grid-cols-1 gap-8 lg:grid-cols-12">
		<div class="lg:col-span-8">
			<div class="aspect-[16/6] w-full rounded-lg bg-indigo-50"></div>
		</div>

		<div class="flex flex-col gap-y-2 lg:col-span-4">
			<h2 class="text-h2">Tu progreso</h2>

			{#each data.modules as modulo}
				<div class="bg-muted/5 flex items-center gap-x-2 rounded-lg border p-3">
					<CheckCircle class="size-4 text-lime-500" />
					<p class="text-muted-foreground text-sm font-medium">{modulo.title}</p>
				</div>
			{/each}
		</div>
	</div>
</div>

<script lang="ts">
	import { CheckCircle } from '@lucide/svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
</script>

<div class="flex flex-col gap-y-4">
	<h1 class="text-h1">{data.course.title}</h1>
	<p class="text-muted-foreground max-w-prose text-pretty">{data.course.description}</p>

	<div class="grid grid-cols-12 gap-8">
		<div class="col-span-8">
			<!-- svelte-ignore a11y_media_has_caption -->
			<video class="aspect-video w-full rounded-lg" src="https://cnre-storage.imlargo.dev/WhatsApp%20Video%202025-09-22%20at%2010.26.20%20AM.mp4" controls></video>
		</div>

		<div class="col-span-4 flex flex-col gap-y-2">
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

<style>
	.bg {
	}
</style>

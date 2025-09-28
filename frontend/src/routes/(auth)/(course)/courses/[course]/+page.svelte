<script lang="ts">
	import { CheckCircle } from '@lucide/svelte';
	import type { PageProps } from './$types';
	import * as Card from '$lib/components/ui/card/index.js';

	let { data }: PageProps = $props();
</script>

<div class="flex flex-col gap-y-4 px-4 md:px-0">
	<h1 class="text-h1">{data.course.title}</h1>
	<p class="text-muted-foreground max-w-prose text-pretty">{data.course.description}</p>

	<div class="grid grid-cols-1 lg:grid-cols-12 gap-8">
		<div class="col-span-1 lg:col-span-8">
			<!-- svelte-ignore a11y_media_has_caption -->
			<video
				class="aspect-video w-full rounded-lg"
				src="https://cnre-storage.imlargo.dev/WhatsApp%20Video%202025-09-22%20at%2010.26.20%20AM.mp4"
				controls
			></video>
		</div>

		<div class="col-span-1 lg:col-span-4 flex flex-col gap-y-2">
			<h2 class="text-h2">Módulos</h2>

			{#each data.modules as modulo}
				<a href="/courses/{data.course.id}/{modulo.id}">
					<Card.Root>
						<Card.Header>
							<Card.Title>{modulo.title}</Card.Title>
							<Card.Description class="line-clamp-1">{modulo.description}</Card.Description>
						</Card.Header>
					</Card.Root>
				</a>
			{/each}
		</div>
	</div>
</div>

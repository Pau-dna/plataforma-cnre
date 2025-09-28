<script lang="ts">
	import { CheckCircle, Circle } from '@lucide/svelte';
	import ProgressBar from '$lib/components/course/ProgressBar.svelte';
	import CourseCompletion from '$lib/components/course/CourseCompletion.svelte';
	import EvaluationStatus from '$lib/components/course/EvaluationStatus.svelte';
	import CompletionRequirements from '$lib/components/course/CompletionRequirements.svelte';
	import type { PageProps } from './$types';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Badge } from '$lib/components/ui/badge';

	let { data }: PageProps = $props();

	const modulesProgress = $derived(
		Object.fromEntries(data.progress.modules_progress.map((mod) => [mod.module_id.toString(), mod]))
	);
</script>

<div class="flex flex-col gap-y-4 px-4 md:px-0">
	<h1 class="text-h1">{data.course.title}</h1>
	<p class="text-muted-foreground max-w-prose text-pretty">{data.course.description}</p>

	<div class="grid grid-cols-1 gap-8 lg:grid-cols-12">
		<div class="col-span-1 lg:col-span-8">
			<!-- svelte-ignore a11y_media_has_caption -->
			<video
				class="aspect-video w-full rounded-lg"
				src="https://cnre-storage.imlargo.dev/WhatsApp%20Video%202025-09-22%20at%2010.26.20%20AM.mp4"
				controls
			></video>
		</div>

		<div class="col-span-1 flex flex-col gap-y-2 lg:col-span-4">
			<h2 class="text-h2">Módulos</h2>

			{#each data.modules as modulo}
				{@const progress = modulesProgress[modulo.id.toString()]}
				<a href="/courses/{data.course.id}/{modulo.id}">
					<Card.Root>
						<Card.Header>
							<div class="flex items-center gap-x-3">
								<Card.Title>{modulo.title}</Card.Title>

								{#if progress.is_completed}
									<Badge
										class="flex max-w-max items-center justify-center gap-2 bg-green-100 text-xs"
									>
										<div class="h-2 w-2 rounded-full bg-green-600"></div>
										<span class="font-medium text-green-800">Completado</span>
									</Badge>
								{:else}
									<Badge
										class="flex max-w-max items-center justify-center gap-2 bg-amber-100 text-xs text-amber-800"
									>
										{progress.percentage.toFixed(2)}%
									</Badge>
								{/if}
							</div>
							<Card.Description class="line-clamp-1">{modulo.description}</Card.Description>
						</Card.Header>
					</Card.Root>
				</a>
			{/each}
		</div>
	</div>
</div>

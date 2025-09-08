<script lang="ts">
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
</script>

<div class="flex flex-col gap-y-4">
	<h1 class="text-h1">{data.course.title}</h1>
	<p class="text-muted-foreground max-w-prose text-pretty">{data.course.description}</p>

	<div class="grid grid-cols-12 gap-8">
		<div class="col-span-8">
			<div class="bg-muted aspect-[16/6] w-full rounded-lg border"></div>
		</div>

		<div class="col-span-4">
			<h2 class="text-h2">Tu progreso</h2>
		</div>
	</div>
</div>

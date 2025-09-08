<script lang="ts">
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
</script>

<div class="flex flex-col gap-y-4">
	<h1 class="text-h1">{data.course.title}</h1>
	<p class="text-pretty max-w-prose text-muted-foreground">{data.course.description}</p>

	<div class="aspect-[16/6] w-full bg-muted border rounded-lg">

	</div>
</div>
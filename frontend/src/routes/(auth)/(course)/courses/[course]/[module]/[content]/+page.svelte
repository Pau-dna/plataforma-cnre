<script lang="ts">
	import VideoPlayer from '$lib/components/course/VideoPlayer.svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	let content = $state(data.content);
</script>

<div class="flex flex-col gap-y-8">
	<div class="flex flex-col gap-y-2">
		<h1 class="text-3xl font-bold">{content.title}</h1>
		<p class="text-muted-foreground">{content.description}</p>
	</div>

	<div class="aspect-video w-full">
		<VideoPlayer url={content?.media_url || ""} />
	</div>
</div>

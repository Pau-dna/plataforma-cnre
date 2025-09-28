<script lang="ts">
	type Props = {
		url: string;
	};

	const { url }: Props = $props();

	function getYouTubeVideoId(url: string): string | null {
		const regex =
			/(?:youtube\.com\/(?:[^/]+\/.+\/|(?:v|e(?:mbed)?)\/|.*[?&]v=)|youtu\.be\/)([^"&?/\s]{11})/;
		const match = url.match(regex);
		return match ? match[1] : null;
	}

	let isLoading = $state(false);
	const videoID = getYouTubeVideoId(url);
	const embedUrl = `https://www.youtube.com/embed/${videoID}?rel=0&modestbranding=1`;
</script>

<div class="aspect-video w-8/12 overflow-hidden rounded-lg bg-black shadow-lg">
	{#if isLoading}
		<div class="bg-muted flex h-full w-8/12 animate-pulse items-center justify-center">
			<div class="text-muted-foreground">{'Cargando video...'}</div>
		</div>
	{/if}
	<!-- svelte-ignore element_invalid_self_closing_tag -->
	<iframe
		src={embedUrl}
		title="Reproductor de video de YouTube"
		class={`h-full w-full ${isLoading ? 'hidden' : 'block'}`}
		allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
		allowfullscreen
		onload={() => (isLoading = false)}
	></iframe>
</div>

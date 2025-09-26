<script lang="ts">
	import { detectUrls, normalizeUrl } from '$lib/utils/url';

	type Props = {
		text: string;
		class?: string;
	};

	const { text, class: className = '' }: Props = $props();

	// Process the text to detect URLs
	const textSegments = $derived(detectUrls(text));
</script>

<div class={className}>
	{#each textSegments as segment}
		{#if segment.isUrl}
			<a
				href={normalizeUrl(segment.text)}
				target="_blank"
				rel="noopener noreferrer"
				class="break-all text-blue-600 underline hover:text-blue-800"
			>
				{segment.text}
			</a>
		{:else}
			{segment.text}
		{/if}
	{/each}
</div>

<script lang="ts">
	import type { Component } from 'svelte';

	type Props = {
		icon: Component;
		title: string;
		description?: string;
		class?: string;
	};

	let { icon: IconComponent, title, description, class: className = '', ...restProps }: Props = $props();
</script>

<div class="flex flex-col items-center justify-center py-12 px-4 text-center {className}" {...restProps}>
	<div class="text-muted-foreground mb-4">
		<IconComponent class="h-12 w-12 mx-auto" />
	</div>
	<h3 class="text-lg font-medium mb-2">{title}</h3>
	{#if description}
		<p class="text-muted-foreground text-sm max-w-sm">{description}</p>
	{/if}
</div>
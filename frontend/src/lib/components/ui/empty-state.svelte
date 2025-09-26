<script lang="ts">
	type Props = {
		icon?: any;
		title: string;
		description?: string;
		action?: {
			text: string;
			href?: string;
			onclick?: () => void;
		};
	};

	const { icon: Icon, title, description, action }: Props = $props();
</script>

<div class="flex flex-col items-center justify-center py-12 text-center">
	{#if Icon}
		<Icon class="h-12 w-12 text-muted-foreground mx-auto mb-4" />
	{/if}
	
	<h3 class="text-lg font-semibold mb-2">{title}</h3>
	
	{#if description}
		<p class="text-muted-foreground mb-4 max-w-sm">{description}</p>
	{/if}
	
	{#if action}
		{#if action.href}
			<a href={action.href} class="inline-flex items-center justify-center rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 bg-primary text-primary-foreground hover:bg-primary/90 h-10 px-4 py-2">
				{action.text}
			</a>
		{:else if action.onclick}
			<button 
				onclick={action.onclick}
				class="inline-flex items-center justify-center rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 bg-primary text-primary-foreground hover:bg-primary/90 h-10 px-4 py-2"
			>
				{action.text}
			</button>
		{/if}
	{/if}
</div>
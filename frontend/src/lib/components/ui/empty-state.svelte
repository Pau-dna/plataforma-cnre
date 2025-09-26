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
		<Icon class="text-muted-foreground mx-auto mb-4 h-12 w-12" />
	{/if}

	<h3 class="mb-2 text-lg font-semibold">{title}</h3>

	{#if description}
		<p class="text-muted-foreground mb-4 max-w-sm">{description}</p>
	{/if}

	{#if action}
		{#if action.href}
			<a
				href={action.href}
				class="ring-offset-background focus-visible:ring-ring bg-primary text-primary-foreground hover:bg-primary/90 inline-flex h-10 items-center justify-center rounded-md px-4 py-2 text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50"
			>
				{action.text}
			</a>
		{:else if action.onclick}
			<button
				onclick={action.onclick}
				class="ring-offset-background focus-visible:ring-ring bg-primary text-primary-foreground hover:bg-primary/90 inline-flex h-10 items-center justify-center rounded-md px-4 py-2 text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50"
			>
				{action.text}
			</button>
		{/if}
	{/if}
</div>

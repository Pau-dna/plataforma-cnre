<script lang="ts">
	type Props = {
		size?: 'sm' | 'md' | 'lg';
		class?: string;
	};

	const { size = 'md', class: className }: Props = $props();

	const sizeClasses = {
		sm: 'h-4 w-4',
		md: 'h-6 w-6',
		lg: 'h-8 w-8'
	};
</script>

<div class="inline-block animate-spin {sizeClasses[size]} {className}">
	<svg class="h-full w-full" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
		<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"
		></circle>
		<path
			class="opacity-75"
			fill="currentColor"
			d="m4 12a8 8 0 0 1 8-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 0 1 4 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
		></path>
	</svg>
</div>

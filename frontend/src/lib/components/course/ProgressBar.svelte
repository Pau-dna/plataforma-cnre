<script lang="ts">
	interface Props {
		progress: number;
		title: string;
		subtitle?: string;
		size?: 'sm' | 'md' | 'lg';
	}

	let { progress, title, subtitle, size = 'md' }: Props = $props();

	const sizeClasses = {
		sm: 'h-2',
		md: 'h-3',
		lg: 'h-4'
	};

	const textSizeClasses = {
		sm: 'text-sm',
		md: 'text-base',
		lg: 'text-lg'
	};
</script>

<div class="flex flex-col gap-2">
	<div class="flex items-center justify-between">
		<div class="flex flex-col">
			<h3 class="{textSizeClasses[size]} font-medium">{title}</h3>
			{#if subtitle}
				<p class="text-muted-foreground text-sm">{subtitle}</p>
			{/if}
		</div>
		<span class="text-muted-foreground text-sm font-medium">{progress}%</span>
	</div>

	<div class="w-full rounded-full bg-gray-200 {sizeClasses[size]}">
		<div
			class="bg-gradient-to-r from-blue-500 to-purple-600 {sizeClasses[
				size
			]} rounded-full transition-all duration-300 ease-out"
			style="width: {progress}%"
		></div>
	</div>
</div>

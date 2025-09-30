<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import type { Icon as IconType } from '@lucide/svelte';
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';

	type Colors = 'sky' | 'pink' | 'teal' | 'yellow';

	type Props = {
		title: string;
		value: string | number;
		icon: typeof IconType;
		color?: Colors;
		subtitle?: string;
		isLoading?: boolean;
	};

	const { title, value, icon: Icon, color = 'sky', subtitle, isLoading = false }: Props = $props();

	const colorClasses: Record<Colors, string> = {
		sky: 'text-sky-500 bg-sky-100',
		pink: 'text-pink-500 bg-pink-100',
		teal: 'text-teal-500 bg-teal-100',
		yellow: 'text-yellow-500 bg-yellow-100'
	};
</script>

<Card.Root class="col-span-4 md:col-span-2 lg:col-span-1">
	<Card.Content>
		<div class="flex items-center justify-between">
			<div class="flex flex-col gap-4">
				<Card.Title>
					{#if isLoading}
						<Skeleton class="h-5 w-20" />
					{:else}
						{title}
					{/if}
				</Card.Title>
				<div class="flex flex-col gap-1">
					<span class="text-primary text-2xl font-bold">
						{#if isLoading}
							<Skeleton class="h-7 w-24" />
						{:else}
							{value}
						{/if}
					</span>
					<span class="text-muted-foreground text-xs">
						{#if subtitle}
							<p class="text-muted-foreground text-xs">
								{#if isLoading}
									<Skeleton class="h-4 w-16" />
								{:else}
									{subtitle}
								{/if}
							</p>
						{/if}
					</span>
				</div>
			</div>
			<div class={`rounded-xl p-3 ${colorClasses[color]}`}>
				<Icon class="h-6 w-6" />
			</div>
		</div>
	</Card.Content>
</Card.Root>

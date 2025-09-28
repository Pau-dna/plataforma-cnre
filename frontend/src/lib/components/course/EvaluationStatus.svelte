<script lang="ts">
	import { Badge } from '$lib/components/ui/badge';
	import { CheckCircle, Clock, XCircle, AlertCircle } from '@lucide/svelte';

	interface Props {
		evaluation: {
			id: number;
			title: string;
			description?: string;
			passing_score: number;
			max_attempts?: number;
			hasPassed?: boolean;
		};
		size?: 'sm' | 'md' | 'lg';
	}

	let { evaluation, size = 'md' }: Props = $props();

	const sizeClasses = {
		sm: 'text-xs',
		md: 'text-sm', 
		lg: 'text-base'
	};

	const iconSizes = {
		sm: 'h-3 w-3',
		md: 'h-4 w-4',
		lg: 'h-5 w-5'
	};
</script>

<div class="flex items-center gap-2">
	{#if evaluation.hasPassed}
		<CheckCircle class="{iconSizes[size]} text-green-600" />
		<span class="{sizeClasses[size]} font-medium text-green-700">{evaluation.title}</span>
		<Badge variant="secondary" class="bg-green-100 text-green-800 text-xs">
			Aprobada
		</Badge>
	{:else}
		<AlertCircle class="{iconSizes[size]} text-amber-500" />
		<span class="{sizeClasses[size]} font-medium text-gray-700">{evaluation.title}</span>
		<Badge variant="secondary" class="bg-amber-100 text-amber-800 text-xs">
			Pendiente
		</Badge>
	{/if}
</div>

{#if evaluation.description && size !== 'sm'}
	<p class="text-xs text-gray-500 ml-6 mt-1">{evaluation.description}</p>
{/if}
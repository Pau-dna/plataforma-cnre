<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import { ArrowDown, ArrowUp, FileText, Ellipsis, GripVertical } from '@lucide/svelte';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import type { Content } from '$lib/types';

	type Props = {
		content: Content;
		actDate?: string;
		onupdate?: (content: Content) => void;
		onmoveup?: (content: Content) => void;
		onmovedown?: (content: Content) => void;
		canMoveUp?: boolean;
		canMoveDown?: boolean;
	};

	const {
		content,
		actDate,
		onupdate,
		onmoveup,
		onmovedown,
		canMoveUp = true,
		canMoveDown = true
	}: Props = $props();

	function handleMoveUp() {
		onmoveup?.(content);
	}

	function handleMoveDown() {
		onmovedown?.(content);
	}

	// Format the date for display
	const formattedDate = actDate ? new Date(actDate).toLocaleDateString('es-ES') : 'N/A';
</script>

<Card.Root>
	<Card.Header>
		<div class="flex items-center justify-between">
			<div class="flex items-center gap-2">
				<Button size="sm" variant="ghost" class="bg-muted">
					<GripVertical class="text-muted-foreground h-4 w-4" />
				</Button>
				<span class="text-muted-foreground font-semibold">Contenido #{content.order}</span>
			</div>
			<DropdownMenu.Root>
				<DropdownMenu.Trigger>
					<Button size="sm" variant="ghost">
						<Ellipsis class="h-4 w-4 leading-none" />
					</Button>
				</DropdownMenu.Trigger>
				<DropdownMenu.Content>
					<DropdownMenu.Group>
						<DropdownMenu.Item>Ver Detalles</DropdownMenu.Item>
						<DropdownMenu.Item>Editar</DropdownMenu.Item>
						<DropdownMenu.Separator />
						<DropdownMenu.Item class="text-destructive">Eliminar</DropdownMenu.Item>
					</DropdownMenu.Group>
				</DropdownMenu.Content>
			</DropdownMenu.Root>
		</div>
		<div class="flex items-center gap-4">
			<div class="flex flex-col gap-0">
				<Button size="sm" variant="ghost" disabled={!canMoveUp} onclick={handleMoveUp}>
					<ArrowUp class="h-4 w-4" />
				</Button>
				<Button size="sm" variant="ghost" disabled={!canMoveDown} onclick={handleMoveDown}>
					<ArrowDown class="h-4 w-4" />
				</Button>
			</div>
			<div class="flex flex-1 flex-col gap-1">
				<div class="flex items-center gap-2">
					<Card.Title class="text-lg">{content.title}</Card.Title>
					<Badge variant="secondary" class="flex items-center gap-1">
						<FileText class="h-3 w-3" />
						<span>Contenido</span>
					</Badge>
				</div>
				{#if content.description}
					<Card.Description>{content.description}</Card.Description>
				{/if}
			</div>
		</div>
	</Card.Header>
	<Card.Content class="flex w-full justify-between">
		<div class="flex items-center gap-4">
			{#if content.media_url}
				<div class="flex items-center gap-2 text-blue-800">
					<FileText class="h-4 w-4 leading-none" />
					<span class="text-sm leading-none">Con multimedia</span>
				</div>
			{/if}
			{#if content.body}
				<div class="flex items-center gap-2 text-green-800">
					<FileText class="h-4 w-4 leading-none" />
					<span class="text-sm leading-none">{content.body.length} caracteres</span>
				</div>
			{/if}
		</div>
		<span class="text-muted-foreground text-sm leading-none">Actualizado el {formattedDate}</span>
	</Card.Content>
</Card.Root>

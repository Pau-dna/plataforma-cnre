<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import { ArrowDown, ArrowUp, FileText, Ellipsis, GripVertical } from '@lucide/svelte';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import DeleteConfirmDialog from '$lib/components/ui/DeleteConfirmDialog.svelte';
	import type { Content } from '$lib/types';
	import { ContentController } from '$lib/controllers/content';
	import { toast } from 'svelte-sonner';

	type Props = {
		content: Content;
		actDate?: string;
		courseId?: number;
		onupdate?: (content: Content) => void;
		ondelete?: (content: Content) => void;
		onmoveup?: (content: Content) => void;
		onmovedown?: (content: Content) => void;
		canMoveUp?: boolean;
		canMoveDown?: boolean;
	};

	const {
		content,
		actDate,
		courseId,
		onupdate,
		ondelete,
		onmoveup,
		onmovedown,
		canMoveUp = true,
		canMoveDown = true
	}: Props = $props();

	let openDelete = $state(false);
	const contentController = new ContentController();

	function handleMoveUp() {
		onmoveup?.(content);
	}

	function handleMoveDown() {
		onmovedown?.(content);
	}

	async function handleDelete() {
		try {
			await contentController.deleteContent(content.id);
			toast.success('Contenido eliminado con éxito.');
			ondelete?.(content);
		} catch (error) {
			console.error('Error deleting content:', error);
			toast.error('Error al eliminar el contenido.', {
				description: error instanceof Error ? error.message : String(error)
			});
		}
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
						{#if courseId}
							<DropdownMenu.Item>
								<a
									href="/admin/courses/{courseId}/{content.module_id}/content/{content.id}/edit"
									class="flex w-full"
								>
									Editar
								</a>
							</DropdownMenu.Item>
						{:else}
							<DropdownMenu.Item>Editar</DropdownMenu.Item>
						{/if}
						<DropdownMenu.Separator />
						<DropdownMenu.Item class="text-destructive" onclick={() => (openDelete = true)}>
							Eliminar
						</DropdownMenu.Item>
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

<DeleteConfirmDialog
	bind:open={openDelete}
	title="¿Eliminar contenido?"
	description="Esta acción eliminará permanentemente el contenido '{content.title}'. No se puede deshacer."
	onConfirm={handleDelete}
/>

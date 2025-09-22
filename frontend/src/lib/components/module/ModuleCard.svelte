<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import { ArrowDown, ArrowUp, BookOpen, CirclePlay, Ellipsis, GripVertical } from '@lucide/svelte';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import EditModule from './EditModule.svelte';
	import DeleteConfirmDialog from '$lib/components/ui/DeleteConfirmDialog.svelte';
	import type { Module } from '$lib/types';
	import { ModuleController } from '$lib/controllers/module';

	type Props = {
		module: Module;
		actDate?: string;
		onupdate?: (module: Module) => void;
		ondelete?: (module: Module) => void;
		onmoveup?: (module: Module) => void;
		onmovedown?: (module: Module) => void;
		canMoveUp?: boolean;
		canMoveDown?: boolean;
	};

	const {
		module: modulo,
		actDate,
		onupdate,
		ondelete,
		onmoveup,
		onmovedown,
		canMoveUp = true,
		canMoveDown = true
	}: Props = $props();
	
	let openEdit = $state(false);
	let openDelete = $state(false);
	const moduleController = new ModuleController();

	function handleModuleUpdate(updated: Module) {
		onupdate?.(updated);
	}

	function handleMoveUp() {
		onmoveup?.(modulo);
	}

	function handleMoveDown() {
		onmovedown?.(modulo);
	}

	async function handleDelete() {
		try {
			await moduleController.deleteModule(modulo.id);
			ondelete?.(modulo);
		} catch (error) {
			console.error('Error deleting module:', error);
			// TODO: Show error toast
		}
	}
</script>

<Card.Root>
	<Card.Header>
		<div class="flex items-center justify-between">
			<div class="flex items-center gap-2">
				<Button size="sm" variant="ghost" class="bg-muted">
					<GripVertical class="text-muted-foreground h-4 w-4" />
				</Button>

				<span class="text-muted-foreground font-semibold">Módulo {modulo.id}</span>
			</div>
			<DropdownMenu.Root>
				<DropdownMenu.Trigger>
					<Button size="sm" variant="ghost">
						<Ellipsis class="h-4 w-4 leading-none" />
					</Button>
				</DropdownMenu.Trigger>
				<DropdownMenu.Content>
					<DropdownMenu.Group>
						<DropdownMenu.Item>
							<a href={`/admin/courses/${modulo.course_id}/${modulo.id}`}> Ver Detalles </a>
						</DropdownMenu.Item>
						<DropdownMenu.Item onclick={() => (openEdit = true)}>Editar</DropdownMenu.Item>
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
			<div class="flex flex-col gap-1">
				<Card.Title class="text-lg">{modulo.title}</Card.Title>
				<Card.Description>{modulo.description}</Card.Description>
			</div>
		</div>
	</Card.Header>
	<Card.Content class="flex w-full justify-between">
		<div class="flex items-center gap-4">
			<Button
				href="/admin/courses/{modulo.course_id}/{modulo.id}"
				variant="ghost"
				size="sm"
				class="flex items-center gap-2 text-blue-800 hover:bg-blue-50 hover:text-blue-900"
			>
				<BookOpen class="h-4 w-4 leading-none" />
				<span class="text-sm leading-none">Gestionar Contenidos</span>
			</Button>
		</div>
		<span class="text-muted-foreground text-sm leading-none">Actualizado el {actDate}</span>
	</Card.Content>
</Card.Root>

<EditModule module={modulo} bind:openEdit onupdate={handleModuleUpdate} />
<DeleteConfirmDialog
	bind:open={openDelete}
	title="¿Eliminar módulo?"
	description="Esta acción eliminará permanentemente el módulo '{modulo.title}' y todo su contenido asociado. No se puede deshacer."
	onConfirm={handleDelete}
/>

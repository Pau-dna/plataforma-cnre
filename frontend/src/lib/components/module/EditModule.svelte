<script lang="ts">
	import type { Snippet } from 'svelte';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Textarea } from '$lib/components/ui/textarea/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { ModuleController } from '$lib';
	import type { Module, UpdateModuleDTO } from '$lib/types';

	type Props = {
		module: Module;
		openEdit?: boolean;
		children?: Snippet;
		onUpdated?: (updatedModule: Module) => void;
	};

	let { module, children, openEdit = $bindable(), onUpdated }: Props = $props();

	let isLoading = $state(false);
	let error = $state('');

	const formdata = $state<UpdateModuleDTO>({
		title: module.title,
		description: module.description || '',
		order: module.order
	});

	// Reset form when module changes or dialog opens
	$effect(() => {
		if (openEdit) {
			formdata.title = module.title;
			formdata.description = module.description || '';
			formdata.order = module.order;
			error = '';
		}
	});

	async function handleSubmit(event: Event) {
		event.preventDefault();
		
		if (!formdata.title?.trim()) {
			error = 'El nombre del módulo es requerido';
			return;
		}

		try {
			isLoading = true;
			error = '';

			const moduleController = new ModuleController();
			const updatedModule = await moduleController.updateModule(module.id, formdata);
			
			onUpdated?.(updatedModule);
			openEdit = false;
		} catch (err) {
			error = err instanceof Error ? err.message : 'Error al actualizar el módulo';
		} finally {
			isLoading = false;
		}
	}
</script>

<Dialog.Root bind:open={openEdit}>
	<Dialog.Content class="form-card flex flex-col gap-6">
		<Dialog.Header class="flex flex-col gap-3">
			<Dialog.Title class="text-h2">Editar módulo</Dialog.Title>
			<Dialog.Description>Complete los detalles para editar el módulo.</Dialog.Description>
		</Dialog.Header>
		<form onsubmit={handleSubmit}>
			<div class="flex flex-col gap-6">
				{#if error}
					<div class="text-sm text-red-600 bg-red-50 border border-red-200 p-3 rounded-md">
						{error}
					</div>
				{/if}
				<div class="flex flex-col gap-2">
					<Label for="edit-module-title">Nombre del módulo</Label>
					<Input 
						bind:value={formdata.title} 
						id="edit-module-title" 
						placeholder="Ingrese el nombre del módulo"
						required
						disabled={isLoading}
					/>
				</div>
				<div class="flex flex-col gap-2">
					<Label for="edit-module-description">Descripción</Label>
					<Textarea 
						bind:value={formdata.description} 
						id="edit-module-description" 
						placeholder="Ingrese la descripción del módulo"
						disabled={isLoading}
					/>
				</div>
				<div class="flex flex-col gap-2">
					<Label for="edit-module-order">Orden</Label>
					<Input 
						bind:value={formdata.order} 
						id="edit-module-order" 
						type="number"
						placeholder="Orden del módulo en el curso"
						min="0"
						disabled={isLoading}
					/>
				</div>
			</div>
			<Dialog.Footer>
				<div class="flex gap-2 w-full">
					<Button 
						type="button" 
						variant="outline" 
						class="w-full"
						onclick={() => openEdit = false}
						disabled={isLoading}
					>
						Cancelar
					</Button>
					<Button 
						type="submit"
						class="w-full bg-pink-500 hover:bg-pink-900"
						disabled={isLoading}
					>
						{#if isLoading}
							Guardando...
						{:else}
							Editar módulo
						{/if}
					</Button>
				</div>
			</Dialog.Footer>
		</form>
	</Dialog.Content>
</Dialog.Root>
<script lang="ts">
	import type { Snippet } from 'svelte';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Textarea } from '$lib/components/ui/textarea/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import type { Module } from '$lib/types';
	import { toast } from 'svelte-sonner';
	import { ModuleController } from '$lib/controllers';

	type Props = {
		module: Module;
		openEdit?: boolean;
		children?: Snippet;
		onupdate?: (module: Module) => void;
	};

	let { module, children, onupdate, openEdit = $bindable() }: Props = $props();

	let submitting = $state(false);
	const formData = $state<Partial<Module>>({
		title: module.title,
		description: module.description
	});

	const moduleController = new ModuleController();
	async function handleUpdate() {
		// Validate data
		if (!formData.title || !formData.description) {
			toast.error('Por favor, complete todos los campos.');
			return;
		}

		try {
			submitting = true;
			const updated = await moduleController.updateModulePatch(module.id, formData);
			toast.success('Módulo actualizado con éxito.');
			openEdit = false;
			onupdate?.(updated);
		} catch (error) {
			toast.error('Error al actualizar el módulo.', {
				description: error instanceof Error ? error.message : String(error)
			});
		} finally {
			submitting = false;
		}
	}

	$effect(() => {
		if (!openEdit) {
			formData.title = module.title;
			formData.description = module.description;
		}
	});
</script>

<Dialog.Root bind:open={openEdit}>
	<Dialog.Content class="form-card flex flex-col gap-6">
		<Dialog.Header class="flex flex-col gap-3">
			<Dialog.Title class="text-h2">Editar módulo</Dialog.Title>
			<Dialog.Description>Complete los detalles para editar el módulo.</Dialog.Description>
		</Dialog.Header>
		<div class="flex flex-col gap-6">
			<div class="flex flex-col gap-2">
				<Label for="name">Nombre del módulo</Label>
				<Input bind:value={formData.title} id="name" placeholder="Ingrese el nombre del módulo" />
			</div>
			<div class="flex flex-col gap-2">
				<Label for="description">Descripción</Label>
				<Textarea
					bind:value={formData.description}
					id="description"
					placeholder="Ingrese la descripción del módulo"
				/>
			</div>
		</div>
		<Dialog.Footer>
			<Button
				onclick={handleUpdate}
				disabled={submitting}
				class="w-full bg-pink-500 hover:bg-pink-900">Editar módulo</Button
			>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>

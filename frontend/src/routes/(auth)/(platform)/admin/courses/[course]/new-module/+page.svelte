<script lang="ts">
	import type { PageProps } from './$types';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Textarea } from '$lib/components/ui/textarea/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import type { Module, CreateModuleDTO } from '$lib';
	import { ModuleController } from '$lib/controllers';
	import { page } from '$app/state';
	import { toast } from 'svelte-sonner';
	import { goto } from '$app/navigation';
	import Back from '$lib/components/kit/Back.svelte';

	let { data }: PageProps = $props();

	let submitting = $state(false);
	const formdata = $state<CreateModuleDTO>({
		title: '',
		description: '',
		course_id: parseInt(page.params.course as string)
	});

	const moduleController = new ModuleController();

	async function handleSubmit() {
		// Validate data
		if (!formdata.title || !formdata.description) {
			toast.error('Por favor, complete todos los campos.');
			return;
		}

		try {
			submitting = true;
			const newModule = await moduleController.createModule(formdata);
			toast.success('Módulo creado con éxito.');
			// Redirect back to the course page
			goto(`/admin/courses/${page.params.course}`);
		} catch (error) {
			toast.error('Error al crear el módulo.', {
				description: error instanceof Error ? error.message : String(error)
			});
		} finally {
			submitting = false;
		}
	}
</script>

<Back href="/admin/courses/{page.params.course}" />

<div class="flex items-center justify-center">
	<Card.Root class="form-card">
		<Card.Header class="flex flex-col gap-3">
			<Card.Title class="text-h2">Crear módulo</Card.Title>
			<Card.Description>Complete los detalles para crear un nuevo módulo.</Card.Description>
		</Card.Header>
		<Card.Content class="flex flex-col gap-6">
			<div class="flex flex-col gap-2">
				<Label for="name">Nombre del módulo</Label>
				<Input id="name" placeholder="Ingrese el nombre del módulo" bind:value={formdata.title} />
			</div>
			<div class="flex flex-col gap-2">
				<Label for="description">Descripción</Label>
				<Textarea
					id="description"
					placeholder="Ingrese la descripción del módulo"
					bind:value={formdata.description}
				/>
			</div>
		</Card.Content>
		<Card.Footer>
			<Button
				onclick={handleSubmit}
				disabled={submitting}
				class="w-full bg-pink-500 hover:bg-pink-900"
			>
				{submitting ? 'Creando...' : 'Crear módulo'}
			</Button>
		</Card.Footer>
	</Card.Root>
</div>

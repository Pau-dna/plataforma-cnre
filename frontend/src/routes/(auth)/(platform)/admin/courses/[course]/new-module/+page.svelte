<script lang="ts">
	import type { PageProps } from './$types';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Textarea } from '$lib/components/ui/textarea/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { ModuleController } from '$lib';
	import type { CreateModuleDTO } from '$lib/types';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';

	let { data }: PageProps = $props();

	let isLoading = $state(false);
	let error = $state('');

	const courseId = parseInt($page.params.course as string);

	const formdata = $state<CreateModuleDTO>({
		title: '',
		description: '',
		order: 0,
		course_id: courseId
	});

	async function handleSubmit(event: Event) {
		event.preventDefault();
		
		if (!formdata.title.trim()) {
			error = 'El nombre del módulo es requerido';
			return;
		}

		try {
			isLoading = true;
			error = '';

			const moduleController = new ModuleController();
			const newModule = await moduleController.createModule(formdata);
			
			// Redirect to course modules list after successful creation
			goto(`/admin/courses/${courseId}`);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Error al crear el módulo';
		} finally {
			isLoading = false;
		}
	}
</script>

<div class="flex items-center justify-center">
	<Card.Root class="form-card">
		<Card.Header class="flex flex-col gap-3">
			<Card.Title class="text-h2">Crea un módulo</Card.Title>
			<Card.Description>Complete los detalles para crear un nuevo módulo.</Card.Description>
		</Card.Header>
		<form onsubmit={handleSubmit}>
			<Card.Content class="flex flex-col gap-6">
				{#if error}
					<div class="text-sm text-red-600 bg-red-50 border border-red-200 p-3 rounded-md">
						{error}
					</div>
				{/if}
				<div class="flex flex-col gap-2">
					<Label for="title">Nombre del módulo</Label>
					<Input 
						id="title" 
						placeholder="Ingrese el nombre del módulo" 
						bind:value={formdata.title}
						required
						disabled={isLoading}
					/>
				</div>
				<div class="flex flex-col gap-2">
					<Label for="description">Descripción</Label>
					<Textarea 
						id="description" 
						placeholder="Ingrese la descripción del módulo" 
						bind:value={formdata.description}
						disabled={isLoading}
					/>
				</div>
				<div class="flex flex-col gap-2">
					<Label for="order">Orden</Label>
					<Input 
						id="order" 
						type="number"
						placeholder="Orden del módulo en el curso" 
						bind:value={formdata.order}
						min="0"
						disabled={isLoading}
					/>
				</div>
			</Card.Content>
			<Card.Footer>
				<Button 
					type="submit"
					class="w-full bg-pink-500 hover:bg-pink-900"
					disabled={isLoading}
				>
					{#if isLoading}
						Creando...
					{:else}
						Crear módulo
					{/if}
				</Button>
			</Card.Footer>
		</form>
	</Card.Root>
</div>

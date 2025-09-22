<script lang="ts">
	import type { PageProps } from './$types';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Textarea } from '$lib/components/ui/textarea/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import type { CreateContentDTO } from '$lib/types';
	import { ContentType } from '$lib/types';
	import { ContentController } from '$lib/controllers';
	import { page } from '$app/state';
	import { toast } from 'svelte-sonner';
	import { goto } from '$app/navigation';
	import Back from '$lib/components/kit/Back.svelte';

	let { data }: PageProps = $props();

	let submitting = $state(false);
	const formdata = $state<CreateContentDTO>({
		title: '',
		description: '',
		body: '',
		media_url: '',
		type: ContentType.CONTENT,
		module_id: data.moduleId
	});

	const contentController = new ContentController();

	async function handleSubmit() {
		// Validate data
		if (!formdata.title || !formdata.body) {
			toast.error('Por favor, complete los campos obligatorios (Título y Contenido).');
			return;
		}

		try {
			submitting = true;
			const newContent = await contentController.createContent(formdata);
			toast.success('Contenido creado con éxito.');
			// Redirect back to the module contents page
			goto(`/admin/courses/${data.courseId}/${data.moduleId}`);
		} catch (error) {
			toast.error('Error al crear el contenido.', {
				description: error instanceof Error ? error.message : String(error)
			});
		} finally {
			submitting = false;
		}
	}
</script>

<Back href="/admin/courses/{data.courseId}/{data.moduleId}" />

<div class="flex items-center justify-center">
	<Card.Root class="form-card w-full max-w-2xl">
		<Card.Header class="flex flex-col gap-3">
			<Card.Title class="text-h2">Crear contenido</Card.Title>
			<Card.Description>
				Complete los detalles para crear un nuevo contenido en el módulo "{data.module.title}".
			</Card.Description>
		</Card.Header>
		<Card.Content class="flex flex-col gap-6">
			<div class="flex flex-col gap-2">
				<Label for="title">Título del contenido *</Label>
				<Input
					id="title"
					placeholder="Ingrese el título del contenido"
					bind:value={formdata.title}
				/>
			</div>

			<div class="flex flex-col gap-2">
				<Label for="description">Descripción</Label>
				<Textarea
					id="description"
					placeholder="Ingrese una descripción opcional del contenido"
					bind:value={formdata.description}
				/>
			</div>

			<div class="flex flex-col gap-2">
				<Label for="body">Contenido *</Label>
				<Textarea
					id="body"
					placeholder="Ingrese el contenido principal (texto, HTML, etc.)"
					bind:value={formdata.body}
					rows={8}
				/>
			</div>

			<div class="flex flex-col gap-2">
				<Label for="media_url">URL de multimedia (opcional)</Label>
				<Input
					id="media_url"
					placeholder="https://ejemplo.com/video.mp4 o https://ejemplo.com/imagen.jpg"
					bind:value={formdata.media_url}
				/>
			</div>

			<div class="flex flex-col gap-2">
				<Label for="order">Orden</Label>
				<Input
					id="order"
					type="number"
					min="1"
					placeholder="Orden de aparición del contenido"
					bind:value={formdata.order}
				/>
			</div>
		</Card.Content>
		<Card.Footer>
			<Button
				onclick={handleSubmit}
				disabled={submitting}
				class="w-full bg-pink-500 hover:bg-pink-900"
			>
				{submitting ? 'Creando...' : 'Crear contenido'}
			</Button>
		</Card.Footer>
	</Card.Root>
</div>

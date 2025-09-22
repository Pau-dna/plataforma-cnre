<script lang="ts">
	import type { PageProps } from './$types';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Textarea } from '$lib/components/ui/textarea/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { CourseController } from '$lib';
	import type { CreateCourseDTO } from '$lib/types';
	import { goto } from '$app/navigation';

	let { data }: PageProps = $props();

	let isLoading = $state(false);
	let error = $state('');

	const formdata = $state<CreateCourseDTO>({
		title: '',
		description: '',
		short_description: '',
		image_url: ''
	});

	async function handleSubmit(event: Event) {
		event.preventDefault();
		
		if (!formdata.title.trim()) {
			error = 'El nombre del curso es requerido';
			return;
		}

		if (!formdata.description.trim()) {
			error = 'La descripción del curso es requerida';
			return;
		}

		try {
			isLoading = true;
			error = '';

			const courseController = new CourseController();
			const newCourse = await courseController.createCourse(formdata);
			
			// Redirect to admin courses list after successful creation
			goto('/admin/courses');
		} catch (err) {
			error = err instanceof Error ? err.message : 'Error al crear el curso';
		} finally {
			isLoading = false;
		}
	}
</script>

<div class="flex items-center justify-center">
	<Card.Root class="form-card">
		<Card.Header class="flex flex-col gap-3">
			<Card.Title class="text-h2">Crea un curso</Card.Title>
			<Card.Description>Complete los detalles para crear un nuevo curso.</Card.Description>
		</Card.Header>
		<form onsubmit={handleSubmit}>
			<Card.Content class="flex flex-col gap-6">
				{#if error}
					<div class="text-sm text-red-600 bg-red-50 border border-red-200 p-3 rounded-md">
						{error}
					</div>
				{/if}
				<div class="flex flex-col gap-2">
					<Label for="title">Nombre del curso</Label>
					<Input 
						id="title" 
						placeholder="Ingrese el nombre del curso" 
						bind:value={formdata.title}
						required
						disabled={isLoading}
					/>
				</div>
				<div class="flex flex-col gap-2">
					<Label for="description">Descripción</Label>
					<Textarea 
						id="description" 
						placeholder="Ingrese la descripción del curso" 
						bind:value={formdata.description}
						required
						disabled={isLoading}
					/>
				</div>
				<div class="flex flex-col gap-2">
					<Label for="short_description">Descripción corta (opcional)</Label>
					<Input 
						id="short_description" 
						placeholder="Ingrese una descripción corta" 
						bind:value={formdata.short_description}
						disabled={isLoading}
					/>
				</div>
				<div class="flex flex-col gap-2">
					<Label for="image_url">URL de imagen (opcional)</Label>
					<Input 
						id="image_url" 
						placeholder="Ingrese la URL de la imagen" 
						bind:value={formdata.image_url}
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
						Crear curso
					{/if}
				</Button>
			</Card.Footer>
		</form>
	</Card.Root>
</div>

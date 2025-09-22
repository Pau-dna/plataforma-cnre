<script lang="ts">
	import type { Snippet } from 'svelte';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Textarea } from '$lib/components/ui/textarea/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { CourseController } from '$lib';
	import type { Course, UpdateCourseDTO } from '$lib/types';

	type Props = {
		course: Course;
		openEdit?: boolean;
		children?: Snippet;
		onUpdated?: (updatedCourse: Course) => void;
	};

	let { course, children, openEdit = $bindable(), onUpdated }: Props = $props();

	let isLoading = $state(false);
	let error = $state('');

	const formdata = $state<UpdateCourseDTO>({
		title: course.title,
		description: course.description,
		short_description: course.short_description || '',
		image_url: course.image_url || ''
	});

	// Reset form when course changes or dialog opens
	$effect(() => {
		if (openEdit) {
			formdata.title = course.title;
			formdata.description = course.description;
			formdata.short_description = course.short_description || '';
			formdata.image_url = course.image_url || '';
			error = '';
		}
	});

	async function handleSubmit(event: Event) {
		event.preventDefault();
		
		if (!formdata.title?.trim()) {
			error = 'El nombre del curso es requerido';
			return;
		}

		if (!formdata.description?.trim()) {
			error = 'La descripción del curso es requerida';
			return;
		}

		try {
			isLoading = true;
			error = '';

			const courseController = new CourseController();
			const updatedCourse = await courseController.updateCourse(course.id, formdata);
			
			onUpdated?.(updatedCourse);
			openEdit = false;
		} catch (err) {
			error = err instanceof Error ? err.message : 'Error al actualizar el curso';
		} finally {
			isLoading = false;
		}
	}
</script>

<Dialog.Root bind:open={openEdit}>
	<Dialog.Content class="form-card flex flex-col gap-6">
		<Dialog.Header class="flex flex-col gap-3">
			<Dialog.Title class="text-h2">Edita el curso</Dialog.Title>
			<Dialog.Description>Complete los detalles para editar el curso.</Dialog.Description>
		</Dialog.Header>
		<form onsubmit={handleSubmit}>
			<div class="flex flex-col gap-6">
				{#if error}
					<div class="text-sm text-red-600 bg-red-50 border border-red-200 p-3 rounded-md">
						{error}
					</div>
				{/if}
				<div class="flex flex-col gap-2">
					<Label for="edit-title">Nombre del curso</Label>
					<Input 
						bind:value={formdata.title} 
						id="edit-title" 
						placeholder="Ingrese el nombre del curso"
						required
						disabled={isLoading}
					/>
				</div>
				<div class="flex flex-col gap-2">
					<Label for="edit-description">Descripción</Label>
					<Textarea 
						bind:value={formdata.description} 
						id="edit-description" 
						placeholder="Ingrese la descripción del curso"
						required
						disabled={isLoading}
					/>
				</div>
				<div class="flex flex-col gap-2">
					<Label for="edit-short-description">Descripción corta (opcional)</Label>
					<Input 
						bind:value={formdata.short_description} 
						id="edit-short-description" 
						placeholder="Ingrese una descripción corta"
						disabled={isLoading}
					/>
				</div>
				<div class="flex flex-col gap-2">
					<Label for="edit-image-url">URL de imagen (opcional)</Label>
					<Input 
						bind:value={formdata.image_url} 
						id="edit-image-url" 
						placeholder="Ingrese la URL de la imagen"
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
							Editar curso
						{/if}
					</Button>
				</div>
			</Dialog.Footer>
		</form>
	</Dialog.Content>
</Dialog.Root>

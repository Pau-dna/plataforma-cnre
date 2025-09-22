<script lang="ts">
	import type { Snippet } from 'svelte';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Textarea } from '$lib/components/ui/textarea/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import type { Course } from '$lib/types';

	type Props = {
		course: Course;
		openEdit?: boolean;
		children?: Snippet;
	};

	let { course, children, openEdit = $bindable() }: Props = $props();
</script>

<Dialog.Root bind:open={openEdit}>
	<Dialog.Content class="form-card flex flex-col gap-6">
		<Dialog.Header class="flex flex-col gap-3">
			<Dialog.Title class="text-h2">Edita el curso</Dialog.Title>
			<Dialog.Description>Complete los detalles para editar el curso.</Dialog.Description>
		</Dialog.Header>
		<div class="flex flex-col gap-6">
			<div class="flex flex-col gap-2">
				<Label for="name">Nombre del curso</Label>
				<Input value={course.title} id="name" placeholder="Ingrese el nombre del curso" />
			</div>
			<div class="flex flex-col gap-2">
				<Label for="description">Descripción</Label>
				<Textarea value={course.description} id="description" placeholder="Ingrese la descripción del curso" />
			</div>
		</div>
		<Dialog.Footer>
			<Button class="w-full bg-pink-500 hover:bg-pink-900">Editar curso</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>

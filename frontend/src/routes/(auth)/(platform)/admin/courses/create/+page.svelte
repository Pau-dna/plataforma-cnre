<script lang="ts">
	import type { PageProps } from './$types';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Textarea } from '$lib/components/ui/textarea/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Label } from '$lib/components/ui/label/index.js';

	let { data }: PageProps = $props();
</script>

<div class="flex items-center justify-center">
	<Card.Root class="form-card">
		<Card.Header class="flex flex-col gap-3">
			<Card.Title class="text-h2">Crea un curso</Card.Title>
			<Card.Description>Complete los detalles para crear un nuevo curso.</Card.Description>
		</Card.Header>
		<Card.Content class="flex flex-col gap-6">
			<div class="flex flex-col gap-2">
				<Label for="name">Nombre del curso</Label>
				<Input id="name" placeholder="Ingrese el nombre del curso" />
			</div>
			<div class="flex flex-col gap-2">
				<Label for="description">Descripción</Label>
				<Textarea id="description" placeholder="Ingrese la descripción del curso" />
			</div>
		</Card.Content>
		<Card.Footer>
			<Button class="w-full bg-pink-500 hover:bg-pink-900">Crear curso</Button>
		</Card.Footer>
	</Card.Root>
</div>

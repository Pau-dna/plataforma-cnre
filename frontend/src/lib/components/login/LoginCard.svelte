<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { ArrowRight } from '@lucide/svelte';
</script>

<Card.Root class="w-md border-sky-500/25 shadow-lg shadow-sky-500/25">
	<img src="/images/logo.png" alt="CNRE Logo" class="mx-auto h-40 w-auto" />
	<Card.Content class="flex flex-col items-center gap-4 p-6">
		<Card.Title class="text-center text-2xl font-bold"
			>Bienvenido(a) a la Plataforma de Cursos del CNRE</Card.Title
		>
		<Card.Description class="text-center"
			>Ingrese con su correo institucional para acceder a todos los recursos disponibles.</Card.Description
		>
	</Card.Content>
	<Card.Footer>
		<Button class="w-full bg-pink-500 hover:bg-pink-900">
			<span>Continuar con Google</span>
			<ArrowRight />
		</Button>
	</Card.Footer>
</Card.Root>

<script lang="ts">
	import { Badge } from '$lib/components/ui/badge/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Progress } from '$lib/components/ui/progress/index.js';
	import * as Table from '$lib/components/ui/table/index.js';
	import { BookOpen, CircleCheckBig, TrendingUp, Users } from '@lucide/svelte';
	type Props = {};

	const {}: Props = $props();

	const estudiantes = [
		{
			nombre: 'Juan Perez',
			curso: 'Capacitación CNRE',
			progreso: '50%',
			estado: 'En progreso',
			inscripcion: '2024-01-15'
		},
		{
			nombre: 'Maria Gomez',
			curso: 'Capacitación CNRE',
			progreso: '100%',
			estado: 'Completado',
			inscripcion: '2024-01-10'
		},
		{
			nombre: 'Luis Rodriguez',
			curso: 'Capacitación CNRE',
			progreso: '20%',
			estado: 'En progreso',
			inscripcion: '2024-01-20'
		},
		{
			nombre: 'Ana Martinez',
			curso: 'Capacitación CNRE',
			progreso: '80%',
			estado: 'En progreso',
			inscripcion: '2024-01-05'
		}
	];
</script>

<div class="grid grid-cols-4 gap-4">
	<Card.Root class="col-span-4 md:col-span-2 lg:col-span-1">
		<Card.Content>
			<div class="flex items-center justify-between">
				<div class="flex flex-col gap-4">
					<Card.Title>Total estudiantes</Card.Title>
					<div class="flex flex-col gap-1">
						<span class="text-primary text-2xl font-bold"> 20 </span>
						<span class="text-muted-foreground text-xs"> Estudiantes activos </span>
					</div>
				</div>
				<div class={`rounded-xl bg-sky-100 p-3 text-sky-500`}>
					<Users class="h-6 w-6" />
				</div>
			</div>
		</Card.Content>
	</Card.Root>

	<Card.Root class="col-span-4 md:col-span-2 lg:col-span-1">
		<Card.Content>
			<div class="flex items-center justify-between">
				<div class="flex flex-col gap-4">
					<Card.Title>Cursos activos</Card.Title>
					<div class="flex flex-col gap-1">
						<span class="text-primary text-2xl font-bold"> 1 </span>
						<span class="text-muted-foreground text-xs"> Cursos disponibles </span>
					</div>
				</div>
				<div class={`rounded-xl bg-pink-100 p-3 text-pink-500`}>
					<BookOpen class="h-6 w-6" />
				</div>
			</div>
		</Card.Content>
	</Card.Root>

	<Card.Root class="col-span-4 md:col-span-2 lg:col-span-1">
		<Card.Content>
			<div class="flex items-center justify-between">
				<div class="flex flex-col gap-4">
					<Card.Title>Tasa de finalización</Card.Title>
					<div class="flex flex-col gap-1">
						<span class="text-primary text-2xl font-bold"> 30% </span>
						<span class="text-muted-foreground text-xs"> Estudiantes que completaron </span>
					</div>
				</div>
				<div class={`rounded-xl bg-teal-100 p-3 text-teal-500`}>
					<CircleCheckBig class="h-6 w-6" />
				</div>
			</div>
		</Card.Content>
	</Card.Root>

	<Card.Root class="col-span-4 md:col-span-2 lg:col-span-1">
		<Card.Content>
			<div class="flex items-center justify-between">
				<div class="flex flex-col gap-4">
					<Card.Title>Progreso promedio</Card.Title>
					<div class="flex flex-col gap-1">
						<span class="text-primary text-2xl font-bold"> 60% </span>
						<span class="text-muted-foreground text-xs"> Estudiantes que completaron </span>
					</div>
				</div>
				<div class={`rounded-xl bg-yellow-100 p-3 text-yellow-500`}>
					<TrendingUp class="h-6 w-6" />
				</div>
			</div>
		</Card.Content>
	</Card.Root>
</div>

<Card.Root>
	<Card.Header>
		<Card.Title>Rendimiento por curso</Card.Title>
		<Card.Description>Estadísticas de inscripción y finalización por curso</Card.Description>
	</Card.Header>
	<Card.Content>
		<Card.Root>
			<Card.Header class="flex justify-between">
				<div class="flex flex-col gap-2">
					<Card.Title>Capacitación CNRE</Card.Title>
					<Card.Description class="line-clamp-1">Capacitación de ciudad y sede</Card.Description>
				</div>
				<Badge>50% completado</Badge>
			</Card.Header>
			<Card.Content class="flex flex-col gap-2">
				<div class="flex items-center justify-between">
					<span>Estudiantes inscritos:</span>
					<span>10</span>
				</div>
				<div class="flex items-center justify-between">
					<span>Han completado:</span>
					<span>5</span>
				</div>
				<Progress value={50} class="h-2 w-full" />
			</Card.Content>
		</Card.Root>
	</Card.Content>
</Card.Root>

<Card.Root>
	<Card.Header>
		<Card.Title>Progreso de Estudiantes</Card.Title>
		<Card.Description>Seguimiento detallado del progreso de cada estudiante</Card.Description>
	</Card.Header>
	<Card.Content>
		<Table.Root>
			<Table.Header>
				<Table.Row>
					<Table.Head>Estudiante</Table.Head>
					<Table.Head>Curso</Table.Head>
					<Table.Head>Progreso</Table.Head>
					<Table.Head>Estado</Table.Head>
					<Table.Head>Inscripción</Table.Head>
				</Table.Row>
			</Table.Header>
			<Table.Body>
				{#each estudiantes as estudiante (estudiante.nombre)}
					<Table.Row>
						<Table.Cell class="font-medium">{estudiante.nombre}</Table.Cell>
						<Table.Cell>{estudiante.curso}</Table.Cell>
						<Table.Cell>{estudiante.progreso}</Table.Cell>
						<Table.Cell>{estudiante.estado}</Table.Cell>
						<Table.Cell>{estudiante.inscripcion}</Table.Cell>
					</Table.Row>
				{/each}
			</Table.Body>
		</Table.Root>
	</Card.Content>
</Card.Root>

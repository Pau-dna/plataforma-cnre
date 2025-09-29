<script lang="ts">
    import * as Avatar from "$lib/components/ui/avatar/index.js";
	import { Badge } from '$lib/components/ui/badge/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Progress } from '$lib/components/ui/progress/index.js';
	import * as Table from '$lib/components/ui/table/index.js';
	import { EnrollmentController } from '$lib/controllers';
	import type { CourseKPIResponse, Enrollment } from '$lib/types';
	import { BookOpen, CircleCheckBig, TrendingUp, Users } from '@lucide/svelte';
	import { onMount } from 'svelte';
	import KPICard from './KPICard.svelte';
	import { Root } from "../ui/alert";
	type Props = {};

	const {}: Props = $props();

	function formatDate(date: string | Date | undefined): string {
		if (!date) return 'N/A';
		const d = typeof date === 'string' ? new Date(date) : date;
		return d.toLocaleDateString('es-ES');
	}

	const enrollmentController = new EnrollmentController();
	let estudiantes = $state<Enrollment[]>();
	let courseKPIs = $state<CourseKPIResponse>();
	let loading = $state(true);

	onMount(async () => {
		estudiantes = await enrollmentController.getCourseEnrollments(1);
		courseKPIs = await enrollmentController.getCourseKPIs(1);
		loading = false;
	});
</script>

<div class="grid grid-cols-4 gap-4">
	<KPICard
		title="Total estudiantes"
		value={loading ? '' : (courseKPIs?.student_count ?? '0')}
		icon={Users}
		color="sky"
		subtitle="Estudiantes activos"
		isLoading={loading}
	/>
	<KPICard
		title="Cursos activos"
		value={loading ? '' : (courseKPIs?.active_courses ?? '0')}
		icon={BookOpen}
		color="pink"
		subtitle="Cursos disponibles"
		isLoading={loading}
	/>
	<KPICard
		title="Tasa de finalización"
		value={loading ? '' : courseKPIs?.completion_rate ? `${courseKPIs.completion_rate}%` : '0%'}
		icon={CircleCheckBig}
		color="teal"
		subtitle="Estudiantes que completaron"
		isLoading={loading}
	/>
	<KPICard
		title="Progreso promedio"
		value={loading ? '' : courseKPIs?.average_progress ? `${courseKPIs.average_progress}%` : '0%'}
		icon={TrendingUp}
		color="yellow"
		subtitle="Estudiantes que completaron"
		isLoading={loading}
	/>
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
				{#if estudiantes && estudiantes.length > 0}
					{#each estudiantes as estudiante (estudiante.id)}
						<Table.Row>
							<Table.Cell class="flex gap-2 items-center">
                                <Avatar.Root class="h-8 w-8 mr-2">
                                    <Avatar.Image src={estudiante.user?.avatar_url} alt={estudiante.user?.fullname || 'Avatar'} />
                                    <Avatar.Fallback>{estudiante.user?.fullname ? estudiante.user.fullname.charAt(0) : 'U'}</Avatar.Fallback>
                                </Avatar.Root>
                                <div class="flex flex-col gap-0">
                                    <span class="font-medium">{estudiante.user?.fullname}</span>
                                    <span class="text-muted-foreground">{estudiante.user?.email}</span>
                                </div>
                            </Table.Cell>
							<Table.Cell>{estudiante.course?.title}</Table.Cell>
							<Table.Cell>
								<Progress value={estudiante.progress} max={100} class="h-2 w-24" />
							</Table.Cell>
							<Table.Cell>
								<Badge
									class={estudiante.progress === 100
										? 'bg-teal-100 text-teal-500'
										: 'bg-pink-100 text-pink-500'}
								>
									{estudiante.progress === 100 ? 'Completado' : 'En progreso'}
								</Badge>
							</Table.Cell>
							<Table.Cell>{formatDate(estudiante.enrolled_at)}</Table.Cell>
						</Table.Row>
					{/each}
				{/if}
			</Table.Body>
		</Table.Root>
	</Card.Content>
</Card.Root>

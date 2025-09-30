<script lang="ts">
	import * as Avatar from '$lib/components/ui/avatar/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Progress } from '$lib/components/ui/progress/index.js';
	import * as Table from '$lib/components/ui/table/index.js';
	import { EnrollmentController } from '$lib/controllers';
	import type { CourseKPIResponse, Enrollment } from '$lib/types';
	import { BookOpen, CircleCheckBig, TrendingUp, Users } from '@lucide/svelte';
	import { onMount } from 'svelte';
	import KPICard from './KPICard.svelte';
	import { Root } from '../ui/alert';
	type Props = {};

	const {}: Props = $props();

	function formatDate(date: string | Date | undefined): string {
		if (!date) return 'N/A';
		const d = typeof date === 'string' ? new Date(date) : date;
		return d.toLocaleDateString('es-ES', {
			year: 'numeric',
			month: 'short',
			day: 'numeric'
		});
	}

	const enrollmentController = new EnrollmentController();
	let students = $state<Enrollment[]>();
	let courseKPIs = $state<CourseKPIResponse>();
	let loading = $state(true);

	let studentCount = $derived.by(() => (students ? students.length : 0));
	let completedCount = $derived.by(() =>
		students ? students.filter((student) => student.progress === 100).length : 0
	);
	let percentageCompleted = $derived.by(() =>
		studentCount > 0 ? Math.round((completedCount / studentCount) * 100) : 0
	);

	onMount(async () => {
		students = await enrollmentController.getCourseEnrollments(1);
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
		value={loading
			? ''
			: courseKPIs?.completion_rate
				? `${courseKPIs.completion_rate.toFixed(2)}%`
				: '0%'}
		icon={CircleCheckBig}
		color="teal"
		subtitle="Estudiantes que completaron"
		isLoading={loading}
	/>
	<KPICard
		title="Progreso promedio"
		value={loading
			? ''
			: courseKPIs?.average_progress
				? `${courseKPIs.average_progress.toFixed(2)}%`
				: '0%'}
		icon={TrendingUp}
		color="yellow"
		subtitle="Progreso de los estudiantes"
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
					<Card.Description class="line-clamp-1"
						>Capacitación de Ciudad y Sede, Antiracismo, VBG, Anticapacitismo y Funcionamiento del
						CNRE</Card.Description
					>
				</div>
				<Badge class="bg-sky-100 text-sky-500">{percentageCompleted}% completado</Badge>
			</Card.Header>
			<Card.Content class="flex flex-col gap-2">
				<div class="flex items-center justify-between">
					<span>Estudiantes inscritos:</span>
					<span>{studentCount}</span>
				</div>
				<div class="flex items-center justify-between">
					<span>Han completado:</span>
					<span>{completedCount}</span>
				</div>
				<Progress value={percentageCompleted} class="h-2 w-full" />
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
				{#if students && students.length > 0}
					{#each students as student (student.id)}
						<Table.Row>
							<Table.Cell class="flex items-center gap-2">
								<Avatar.Root class="mr-2 h-8 w-8">
									<Avatar.Image
										src={student.user?.avatar_url}
										alt={student.user?.fullname || 'Avatar'}
									/>
									<Avatar.Fallback
										>{student.user?.fullname
											? student.user.fullname.charAt(0)
											: 'U'}</Avatar.Fallback
									>
								</Avatar.Root>
								<div class="flex flex-col gap-0">
									<span class="font-medium">{student.user?.fullname}</span>
									<span class="text-muted-foreground">{student.user?.email}</span>
								</div>
							</Table.Cell>
							<Table.Cell>{student.course?.title}</Table.Cell>
							<Table.Cell>
								<Progress value={student.progress} max={100} class="h-2 w-24" />
							</Table.Cell>
							<Table.Cell>
								<Badge
									class={student.progress === 100
										? 'bg-teal-100 text-teal-500'
										: 'bg-pink-100 text-pink-500'}
								>
									{student.progress === 100 ? 'Completado' : 'En progreso'}
								</Badge>
							</Table.Cell>
							<Table.Cell>{formatDate(student.enrolled_at)}</Table.Cell>
						</Table.Row>
					{/each}
				{/if}
			</Table.Body>
		</Table.Root>
	</Card.Content>
</Card.Root>

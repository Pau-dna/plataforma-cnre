<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import { onMount } from 'svelte';
	import KPICard from './KPICard.svelte';
	import { BookOpen, TrendingUp, Trophy } from '@lucide/svelte';
	import { EnrollmentController } from '$lib/controllers';
	import type { Enrollment } from '$lib/types';
	import { authStore } from '$lib/stores/auth.svelte';
	import { Progress } from '$lib/components/ui/progress/index.js';
	import Badge from '../ui/badge/badge.svelte';

	let isLoading = $state(true);

	const enrollmentController = new EnrollmentController();
	let enrollments = $state<Enrollment[]>();

	let completedCourses = $derived.by(() =>
		enrollments ? enrollments.filter((enrollment) => enrollment.progress === 100).length : 0
	);

	let mediaProgress = $derived.by(() => {
		if (!enrollments || enrollments.length === 0) return 0;
		const totalProgress = enrollments.reduce(
			(sum, enrollment) => sum + (enrollment.progress || 0),
			0
		);
		return Math.round(totalProgress / enrollments.length);
	});

	let estimatedTimeToComplete = $derived.by(() => {
		if (!enrollments || enrollments.length === 0) return '0h';

		const totalRemainingHours = enrollments.reduce((sum, enrollment) => {
			// Asumiendo que cada curso tiene una duración estimada (en horas)
			// Si no existe, puedes usar un valor por defecto
			const courseDuration = 10; // 10 horas por defecto
			const progressDecimal = (enrollment.progress || 0) / 100;
			const remainingProgress = 1 - progressDecimal;
			const remainingHours = courseDuration * remainingProgress;

			return sum + remainingHours;
		}, 0);

		// Formatear el resultado
		if (totalRemainingHours < 1) {
			return `${Math.round(totalRemainingHours * 60)}min`;
		} else {
			return `${Math.round(totalRemainingHours)}h`;
		}
	});

	onMount(async () => {
		enrollments = await enrollmentController.getUserEnrollments(authStore?.user?.id as number);
		isLoading = false;
	});
</script>

<div class="grid grid-cols-4 gap-4">
	<KPICard
		title="Cursos activos"
		value={enrollments?.length || 0}
		icon={BookOpen}
		color="sky"
		subtitle="Cursos inscritos"
		{isLoading}
	/>
	<KPICard
		title="Completados"
		value={completedCourses}
		icon={Trophy}
		color="pink"
		subtitle="Cursos terminados"
		{isLoading}
	/>
	<KPICard
		title="Progreso promedio"
		value="{mediaProgress}%"
		icon={TrendingUp}
		color="teal"
		subtitle="en todos los cursos"
		{isLoading}
	/>
	<KPICard
		title="Tiempo estimado"
		value={estimatedTimeToComplete}
		icon={BookOpen}
		color="yellow"
		subtitle="para completar"
		{isLoading}
	/>
</div>

<div class="grid grid-cols-2 gap-4">
	<Card.Root>
		<Card.Header>
			<Card.Title>Mis Cursos</Card.Title>
			<Card.Description>Tu progreso en cada curso inscrito</Card.Description>
		</Card.Header>
		<Card.Content>
			{#if enrollments && enrollments.length > 0}
				{#each enrollments as enrollment (enrollment.id)}
					<div class="flex items-start gap-3">
						<div class="size-12 rounded-xl bg-gradient-to-r from-sky-500 to-pink-500">
							{#if enrollment.course?.image_url}
								<img src={enrollment.course?.image_url} alt="" class="object-cover" />
							{/if}
						</div>
						<div class="min-w-0 flex-1">
							<div class="mb-1 flex items-center justify-between">
								<span class="line-clamp-1 font-semibold">{enrollment.course?.title}</span>
								<Badge
									class={enrollment.progress === 100
										? 'bg-teal-100 text-teal-500'
										: 'bg-pink-100 text-pink-500'}
									>{enrollment.progress === 100 ? 'completado' : enrollment.progress}</Badge
								>
							</div>
							<span class="text-muted-foreground mb-2 line-clamp-1 text-sm"
								>{enrollment.course?.description}</span
							>

							<Progress value={enrollment.progress} class="h-2 w-full" />
						</div>
					</div>
				{/each}
			{/if}
		</Card.Content>
	</Card.Root>
	<Card.Root>
		<Card.Header>
			<Card.Title>Actividad Reciente</Card.Title>
			<Card.Description>Tus últimas acciones en la plataforma</Card.Description>
		</Card.Header>
        <Card.Content>
            
        </Card.Content>
	</Card.Root>
</div>

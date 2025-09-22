<script lang="ts">
	import CourseCard from '$lib/components/course/CourseCard.svelte';
	import EmptyState from '$lib/components/ui/empty-state';
	import { BookMarked } from '@lucide/svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	let courses = $state((data?.enrollments || []).map((enrollment) => enrollment.course).filter(Boolean));
</script>

<div class="flex flex-col gap-6">
	<h1 class="text-h1">Mis Cursos</h1>

	{#if courses.length === 0}
		<EmptyState
			icon={BookMarked}
			title="No tienes cursos inscritos"
			description="Aún no estás inscrito en ningún curso. Explora los cursos disponibles y comienza tu aprendizaje."
		/>
	{:else}
		<div class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
			{#each courses as course}
				<CourseCard {course} />
			{/each}
		</div>
	{/if}
</div>

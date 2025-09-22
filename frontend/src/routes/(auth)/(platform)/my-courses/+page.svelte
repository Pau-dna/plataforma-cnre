<script lang="ts">
	import CourseCard from '$lib/components/course/CourseCard.svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	let courses = $state((data?.enrollments || []).map((enrollment) => enrollment.course));
</script>

<div class="flex flex-col gap-6">
	<h1 class="text-h1">Mis Cursos</h1>

	{#if courses.length === 0}
		<div class="flex flex-col items-center justify-center py-12 text-center">
			<p class="text-muted-foreground text-lg">No estás inscrito en ningún curso</p>
		</div>
	{:else}
		<div class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
			{#each courses as course}
				<CourseCard {course} />
			{/each}
		</div>
	{/if}
</div>

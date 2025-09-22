<script lang="ts">
	import type { PageProps } from './$types';
	import Button from '$lib/components/ui/button/button.svelte';
	import { Plus } from '@lucide/svelte';
	import type { Course } from '$lib/types/models/course';
	import CourseCard from '$lib/components/course/CourseCard.svelte';

	let { data }: PageProps = $props();

	let courses = $state(data.courses);

	function handleUpdate(updated: Course) {
		const index = courses.findIndex((c) => c.id === updated.id);
		courses[index] = updated;
	}
</script>

<div class="flex flex-col gap-6">
	<div class="flex items-center justify-between">
		<h1 class="text-h1">Administrar Cursos</h1>
		<!--
		<Button href="/admin/courses/create" class="bg-pink-500 hover:bg-pink-900">
			<Plus class="h-4 w-4 leading-none" />
			<span class="leading-none">Crear Curso</span>
		</Button>
		-->
	</div>
	{#if courses.length === 0}
		<p>No estás inscrito en ningún curso.</p>
	{:else}
		<div class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
			{#each courses as course}
				<CourseCard onupdate={handleUpdate} {course} />
			{/each}
		</div>
	{/if}
</div>

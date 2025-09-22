<script lang="ts">
	import type { PageProps } from './$types';
	import Button from '$lib/components/ui/button/button.svelte';
	import { Plus } from '@lucide/svelte';
	import type { Course } from '$lib/types/models/course';
	import CourseCard from '$lib/components/course/CourseCard.svelte';

	let { data }: PageProps = $props();

	let courses: Course[] = [
		{
			id: 1,
			title: 'Curso de Ejemplo',
			description: 'Descripción del curso de ejemplo.',
			created_at: '2023-10-01T12:00:00Z',
			updated_at: '2023-10-01T12:00:00Z',
			modules: [],
			student_count: 25,
			module_count: 5
		},
		{
			id: 2,
			title: 'Otro Curso',
			description: 'Descripción de otro curso.',
			created_at: '2023-10-02T12:00:00Z',
			updated_at: '2023-10-02T12:00:00Z',
			modules: [],
			student_count: 10,
			module_count: 3
		},
		{
			id: 3,
			title: 'Tercer Curso',
			description: 'Descripción del tercer curso.',
			created_at: '2023-10-03T12:00:00Z',
			updated_at: '2023-10-03T12:00:00Z',
			modules: [],
			student_count: 5,
			module_count: 1
		}
	];
</script>

<div class="flex flex-col gap-6">
	<div class="flex items-center justify-between">
		<h1 class="text-h1">Administrar Cursos</h1>
		<Button href="/admin/courses/create" class="bg-pink-500 hover:bg-pink-900">
			<Plus class="h-4 w-4 leading-none" />
			<span class="leading-none">Crear Curso</span>
		</Button>
	</div>
	{#if courses.length === 0}
		<p>No estás inscrito en ningún curso.</p>
	{:else}
		<div class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
			{#each courses as course}
				<CourseCard {course} />
			{/each}
		</div>
	{/if}
</div>

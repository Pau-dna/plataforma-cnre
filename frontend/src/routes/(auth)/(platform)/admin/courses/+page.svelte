<script lang="ts">
	import type { PageProps } from './$types';
	import Button from '$lib/components/ui/button/button.svelte';
	import { Plus, GripVertical } from '@lucide/svelte';
	import type { Course } from '$lib/types/models/course';
	import CourseCard from '$lib/components/course/CourseCard.svelte';
	import { dndzone } from 'svelte-dnd-action';
	import { CourseController } from '$lib/controllers';
	import { toast } from 'svelte-sonner';

	let { data }: PageProps = $props();

	let courses = $state(data.courses.map((course, index) => ({ ...course, order: course.order ?? index })));
	let dragDisabled = $state(true);

	const courseController = new CourseController();

	function handleUpdate(updated: Course) {
		const index = courses.findIndex((c) => c.id === updated.id);
		courses[index] = updated;
	}

	function handleDndConsider(e: CustomEvent) {
		courses = e.detail.items;
	}

	async function handleDndFinalize(e: CustomEvent) {
		courses = e.detail.items;
		
		// Update order for all courses
		const courseOrders = courses.map((course, index) => ({
			id: course.id,
			order: index
		}));

		try {
			await courseController.reorderCourses(courseOrders);
			toast.success('Orden de cursos actualizado correctamente');
		} catch (error) {
			toast.error('Error al actualizar el orden de los cursos');
			console.error(error);
		}
	}

	function startDrag() {
		dragDisabled = false;
	}

	function endDrag() {
		dragDisabled = true;
	}
</script>

<div class="flex flex-col gap-6">
	<div class="flex items-center justify-between">
		<h1 class="text-h1">Administrar Cursos</h1>
		<div class="flex items-center gap-2">
			<Button variant="outline" size="sm" onclick={startDrag}>
				<GripVertical class="h-4 w-4 leading-none" />
				<span class="leading-none">Reordenar</span>
			</Button>
			<!--
			<Button href="/admin/courses/create" class="bg-pink-500 hover:bg-pink-900">
				<Plus class="h-4 w-4 leading-none" />
				<span class="leading-none">Crear Curso</span>
			</Button>
			-->
		</div>
	</div>
	{#if courses.length === 0}
		<p>No estás inscrito en ningún curso.</p>
	{:else}
		<div 
			class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3"
			use:dndzone={{
				items: courses,
				dragDisabled,
				dropTargetStyle: {}
			}}
			onconsider={handleDndConsider}
			onfinalize={handleDndFinalize}
		>
			{#each courses as course (course.id)}
				<div class="relative">
					{#if !dragDisabled}
						<div class="absolute -top-2 -left-2 z-10 bg-white border rounded-full p-1 shadow-md cursor-grab">
							<GripVertical class="h-4 w-4 text-gray-500" />
						</div>
					{/if}
					<CourseCard onupdate={handleUpdate} {course} />
				</div>
			{/each}
		</div>
		{#if !dragDisabled}
			<div class="flex justify-center">
				<Button variant="outline" onclick={endDrag}>
					Finalizar reordenamiento
				</Button>
			</div>
		{/if}
	{/if}
</div>

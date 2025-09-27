<script lang="ts">
	import { CheckCircle, Circle } from '@lucide/svelte';
	import ProgressBar from '$lib/components/course/ProgressBar.svelte';
	import CourseCompletion from '$lib/components/course/CourseCompletion.svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
</script>

<div class="flex flex-col gap-y-4">
	<h1 class="text-h1">{data.course.title}</h1>
	<p class="text-muted-foreground max-w-prose text-pretty">{data.course.description}</p>

	<div class="grid grid-cols-12 gap-8">
		<div class="col-span-8">
			<div class="aspect-[16/6] w-full rounded-lg bg-indigo-50"></div>
		</div>

		<div class="col-span-4 flex flex-col gap-y-4">
			<h2 class="text-h2">Tu progreso</h2>
			
			<!-- Overall Course Progress -->
			<div class="bg-white border rounded-lg p-4">
				<ProgressBar 
					progress={data.courseProgressPercentage} 
					title="Progreso del curso"
					subtitle="{data.courseProgressPercentage}% completado"
					size="lg"
				/>
			</div>

			<!-- Module Progress -->
			<div class="space-y-3">
				{#each data.modules as modulo}
					<div class="bg-muted/5 flex flex-col gap-3 rounded-lg border p-4">
						<div class="flex items-center justify-between">
							<div class="flex items-center gap-x-2">
								{#if modulo.progressPercentage >= 100}
									<CheckCircle class="size-4 text-green-600" />
								{:else}
									<Circle class="size-4 text-gray-400" />
								{/if}
								<p class="text-sm font-medium">{modulo.title}</p>
							</div>
							<span class="text-xs text-muted-foreground">{modulo.progressPercentage}%</span>
						</div>
						
						<!-- Module Progress Bar -->
						<div class="w-full bg-gray-200 rounded-full h-2">
							<div
								class="bg-gradient-to-r from-blue-500 to-purple-600 h-2 rounded-full transition-all duration-300"
								style="width: {modulo.progressPercentage}%"
							></div>
						</div>

						<!-- Content List -->
						{#if modulo.contents && modulo.contents.length > 0}
							<div class="mt-2 space-y-1">
								{#each modulo.contents as content}
									<div class="flex items-center gap-2 text-xs">
										{#if content.isCompleted}
											<CheckCircle class="size-3 text-green-500" />
											<span class="text-green-700">{content.title}</span>
										{:else}
											<Circle class="size-3 text-gray-300" />
											<span class="text-gray-500">{content.title}</span>
										{/if}
									</div>
								{/each}
							</div>
						{/if}
					</div>
				{/each}
			</div>

			<!-- Course Completion Status -->
			{#if data.courseProgressPercentage >= 100 && data.enrollment.completed_at}
				<CourseCompletion 
					courseName={data.course.title}
					completedAt={data.enrollment.completed_at}
					progressPercentage={data.courseProgressPercentage}
				/>
			{:else if data.courseProgressPercentage >= 100}
				<div class="bg-green-50 border border-green-200 rounded-lg p-4">
					<div class="flex items-center gap-2">
						<CheckCircle class="h-5 w-5 text-green-600" />
						<div>
							<p class="text-sm font-medium text-green-800">¡Curso completado!</p>
							<p class="text-xs text-green-600">Has terminado todos los módulos de este curso</p>
						</div>
					</div>
				</div>
			{/if}
		</div>
	</div>
</div>

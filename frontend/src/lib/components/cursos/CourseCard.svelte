<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import type { Course } from '$lib/types/models/course';
	import { Clock, Users } from '@lucide/svelte';

	type Props = {
		course: Course;
	};

	const { course }: Props = $props();
</script>

<a href="/courses/{course.id}">
	<Card.Root class="pt-0 hover:shadow-lg">
		<!-- Header con gradiente -->
		<div class="h-44 w-full rounded-t-xl bg-gradient-to-r from-indigo-400 to-pink-400"></div>

		<Card.Header>
			<Card.Title class="text-lg font-bold">{course.title}</Card.Title>
			<Card.Description>
				<div class="flex flex-col gap-4">
					<span class="line-clamp-2">{course.description}</span>
					<div>
						<div class="flex items-center gap-4">
							<div class="flex items-center justify-center gap-1">
								<Users class="h-4 w-4 leading-none" />
								<span class="leading-none">
									{course.student_count}
									{course.student_count === 1 ? 'estudiante' : 'estudiantes'}
								</span>
							</div>
							<div class="flex items-center justify-center gap-1">
								<Clock class="h-4 w-4 leading-none" />
								<span class="leading-none">
									{course.module_count}
									{course.module_count === 1 ? 'módulo' : 'módulos'}
								</span>
							</div>
						</div>
					</div>
				</div>
			</Card.Description>
		</Card.Header>
	</Card.Root>
</a>

<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import type { Course } from '$lib/types/models/course';

	type Props = {
		course: Course;
	};

	const { course }: Props = $props();
</script>

<Card.Root>
	<Card.Header>
		<Card.Title>{course.title}</Card.Title>
		<Card.Description>
			<div class="flex flex-col gap-4">
				<span class="line-clamp-2">{course.description}</span>
				<div>
					<span class="text-sm text-muted-foreground">
						{course.student_count} {course.student_count === 1 ? 'estudiante' : 'estudiantes'} · {course.module_count} {course.module_count === 1 ? 'módulo' : 'módulos'}
					</span>
				</div>
			</div>
		</Card.Description>
	</Card.Header>
</Card.Root>

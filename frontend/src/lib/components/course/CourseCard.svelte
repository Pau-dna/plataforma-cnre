<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import type { Course } from '$lib/types/models/course';
	import { Clock, Ellipsis, Users } from '@lucide/svelte';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';

	type Props = {
		course: Course;
	};

	const { course }: Props = $props();
</script>

<a href="/courses/{course.id}">
	<Card.Root class="pt-0 hover:shadow-lg">
		<!-- Header con gradiente -->
		<div class="h-44 w-full rounded-t-xl bg-gradient-to-r from-sky-500 to-pink-500"></div>

		<Card.Header>
			<div class="flex items-center justify-between">
				<Card.Title class="text-lg font-bold">{course.title}</Card.Title>
				<DropdownMenu.Root>
					<DropdownMenu.Trigger>
						<Button size="sm" variant="ghost">
							<Ellipsis class="h-4 w-4 leading-none" />
						</Button>
					</DropdownMenu.Trigger>
					<DropdownMenu.Content>
						<DropdownMenu.Group>
							<DropdownMenu.Item>Ver</DropdownMenu.Item>
							<DropdownMenu.Item>Editar</DropdownMenu.Item>
							<DropdownMenu.Separator />
							<DropdownMenu.Item class="text-destructive">Eliminar</DropdownMenu.Item>
						</DropdownMenu.Group>
					</DropdownMenu.Content>
				</DropdownMenu.Root>
			</div>
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

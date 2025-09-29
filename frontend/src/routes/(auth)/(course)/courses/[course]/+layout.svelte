<script lang="ts">
	import type { LayoutProps } from './$types';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import CourseSidebar from '$lib/components/course/CourseSidebar.svelte';
	import { Separator } from '$lib/components/ui/separator/index.js';
	import * as Breadcrumb from '$lib/components/ui/breadcrumb/index.js';
	import { page } from '$app/state';
	import { toast } from 'svelte-sonner';
	import { goto } from '$app/navigation';
	import { enrollmentController } from '$lib';

	let { data, children }: LayoutProps = $props();

	let isCompleted = $derived(data.enrollment.progress === 100);
	let prevIsCompleted = isCompleted;

	$effect(() => {
		if (isCompleted && !prevIsCompleted) {
			toast.success('¡Felicidades! Has completado el curso.', {
				duration: 5000,
				action: {
					label: 'Ver certificado',
					onClick: () => {
						goto(`/courses/${data.course.id}/completed`);
					}
				}
			});
		}

		prevIsCompleted = isCompleted;
	});	

</script>

<div class="min-h-auto relative h-full overflow-hidden rounded-md border">
	<Sidebar.Provider>
		<CourseSidebar enrollment={data.enrollment} course={data.course} modules={data.modules} />
		<main class="flex w-full flex-col gap-y-12 p-8">
			<div class="flex items-center gap-x-2">
				<Sidebar.Trigger />
				<Separator orientation="vertical" class="h-1 max-h-max" />

				<Breadcrumb.Root>
					<Breadcrumb.List>
						<Breadcrumb.Item>
							<Breadcrumb.Link href="/courses/{data.course.id}">{data.course.title}</Breadcrumb.Link
							>
						</Breadcrumb.Item>
						{#if page.params.module}
							<Breadcrumb.Separator />
							<Breadcrumb.Item>
								<Breadcrumb.Link href="/courses/{data.course.id}/{page.params.module}"
									>{data.modules.find((m) => m.id === parseInt(page.params.module || '-1'))
										?.title}</Breadcrumb.Link
								>
							</Breadcrumb.Item>
						{/if}
						{#if page.params.content}
							<Breadcrumb.Separator />
							<Breadcrumb.Item>
								<Breadcrumb.Link
									href="/courses/{data.course.id}/{page.params.module}/{page.params.content}"
								>
									{data.modules
										.find((m) => m.id === parseInt(page.params.module || '-1'))
										?.contents.find((c) => c.id === parseInt(page.params.content || '-1'))?.title ||
										'Content'}
								</Breadcrumb.Link>
							</Breadcrumb.Item>
						{/if}
					</Breadcrumb.List>
				</Breadcrumb.Root>
			</div>

			{#if data.enrollment.progress === 100}
				<div class="rounded-md bg-green-50 p-4 text-sm text-green-800">
					¡Felicidades! Has completado el curso. Puedes ver tu certificado{' '}
					<a
						href="/certificado/{data.enrollment.id}"
						class="font-medium underline underline-offset-4"
						> aquí</a
					>.
				</div>
			{/if}

			<div>
				{@render children?.()}
			</div>
		</main>
	</Sidebar.Provider>
</div>

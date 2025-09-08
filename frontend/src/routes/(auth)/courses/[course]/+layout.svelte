<script lang="ts">
	import type { LayoutProps } from './$types';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import CourseSidebar from '$lib/components/course/CourseSidebar.svelte';
	import { Separator } from '$lib/components/ui/separator/index.js';
	import * as Breadcrumb from '$lib/components/ui/breadcrumb/index.js';
	import { page } from '$app/state';

	let { data, children }: LayoutProps = $props();
</script>

<div class="min-h-auto relative h-full overflow-hidden rounded-md border">
	<Sidebar.Provider>
		<CourseSidebar course={data.course} modules={data.modules} />
		<main class="flex w-full flex-col gap-y-12 p-8">
			<div class="flex items-center gap-x-2">
				<Sidebar.Trigger />
				<Separator orientation="vertical" />

				<Breadcrumb.Root>
					<Breadcrumb.List>
						<Breadcrumb.Item>
							<Breadcrumb.Link href="/courses/{data.course.id}">{data.course.title}</Breadcrumb.Link
							>
						</Breadcrumb.Item>
						{#if page.params.module}
							<Breadcrumb.Separator />
							<Breadcrumb.Item>
								<Breadcrumb.Link href="/docs/components"
									>{data.modules.find((m) => m.id === parseInt(page.params.module || "-1"))
										?.title}</Breadcrumb.Link
								>
							</Breadcrumb.Item>
						{/if}
					</Breadcrumb.List>
				</Breadcrumb.Root>
			</div>

			<div>
				{@render children?.()}
			</div>
		</main>
	</Sidebar.Provider>
</div>

<script lang="ts">
	import type { LayoutProps } from './$types';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import CourseSidebar from '$lib/components/course/CourseSidebar.svelte';
	import Separator from '$lib/components/ui/separator/separator.svelte';

	let { data, children }: LayoutProps = $props();
</script>

<div class="min-h-auto relative h-full overflow-hidden rounded-md border">
	<Sidebar.Provider>
		<CourseSidebar />
		<main class="flex w-full flex-col gap-y-12 p-8">
			<div class="flex gap-x-2 items-center">
				<Sidebar.Trigger />
				<Separator orientation="vertical" />
			</div>

			<div>
				{@render children?.()}
			</div>
		</main>
	</Sidebar.Provider>
</div>

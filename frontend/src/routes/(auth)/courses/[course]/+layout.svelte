<script lang="ts">
	import type { LayoutProps } from './$types';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import CourseSidebar from '$lib/components/course/CourseSidebar.svelte';

	let { data, children }: LayoutProps = $props();
</script>

<div class="relative border rounded-md h-full min-h-auto overflow-hidden">
	<Sidebar.Provider>
		<CourseSidebar />
		<main class="p-4">
			<Sidebar.Trigger />
			<div>
				{@render children?.()}
			</div>
		</main>
	</Sidebar.Provider>
</div>

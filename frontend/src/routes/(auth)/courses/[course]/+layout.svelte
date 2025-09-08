<script lang="ts">
	import type { LayoutProps } from './$types';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import CourseSidebar from '$lib/components/course/CourseSidebar.svelte';

	let { data, children }: LayoutProps = $props();
</script>

<div class="relative">
	<Sidebar.Provider>
		<CourseSidebar />
		<main>
			<Sidebar.Trigger />
			<div>
				{@render children?.()}
			</div>
		</main>
	</Sidebar.Provider>
</div>

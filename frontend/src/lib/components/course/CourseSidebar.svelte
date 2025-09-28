<script lang="ts">
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import type { ComponentProps } from 'svelte';
	import { Book, Home } from '@lucide/svelte';
	import NavGroup from './NavGroup.svelte';
	import type { Course, Module } from '$lib/types/models/course';

	type Props = {
		course: Course;
		modules: Module[];
	};

	let {
		ref = $bindable(null),
		collapsible = 'icon',
		course,
		modules = [],
		...restProps
	}: ComponentProps<typeof Sidebar.Root> & Props = $props();

	const data = $derived({
		main: [
			{
				title: 'Resumen',
				url: `/courses/${course.id.toString()}`,
				icon: Home
			},
			{
				title: 'Modulos',
				url: '#',
				icon: Book,
				active: true,
				items: modules.map((m) => ({
					title: m.title,
					url: `/courses/${course.id.toString()}/${m.id}`
				}))
			}
		],
		info: [
			{
				title: 'Evaluaciones',
				url: '#',
				icon: Book
			}
		]
	});
</script>

<Sidebar.Root {collapsible} {...restProps}>
	<!--
	<Sidebar.Header>
		<TeamSwitcher />
	</Sidebar.Header>
	-->
	<Sidebar.Content>
		<NavGroup title="Curso" items={data.main} />
	</Sidebar.Content>
	<Sidebar.Rail />
</Sidebar.Root>

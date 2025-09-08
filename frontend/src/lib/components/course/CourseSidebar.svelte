<script lang="ts" module>
	import AudioWaveformIcon from '@lucide/svelte/icons/audio-waveform';
	import BookOpenIcon from '@lucide/svelte/icons/book-open';
	import BotIcon from '@lucide/svelte/icons/bot';
	import ChartPieIcon from '@lucide/svelte/icons/chart-pie';
	import CommandIcon from '@lucide/svelte/icons/command';
	import FrameIcon from '@lucide/svelte/icons/frame';
	import GalleryVerticalEndIcon from '@lucide/svelte/icons/gallery-vertical-end';
	import MapIcon from '@lucide/svelte/icons/map';
	import Settings2Icon from '@lucide/svelte/icons/settings-2';
	import SquareTerminalIcon from '@lucide/svelte/icons/square-terminal';

	const x = $state(12)
	// This is sample data.
	const data = {
		main: [
			{
				title: 'Resumen',
				url: '#',
				icon: Home
			},
			{
				title: 'Modulos',
				url: '#',
				icon: Book,
				active: true,
				items: [
					{
						title: 'History',
						url: '#'
					},
					{
						title: 'Starred',
						url: '#'
					},
					{
						title: 'Settings',
						url: '#'
					}
				]
			}
		],
		info: [
			{
				title: 'Evaluaciones',
				url: '#',
				icon: FrameIcon
			}
		]
	};
</script>

<script lang="ts">
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import type { ComponentProps } from 'svelte';
	import { Book, Home } from '@lucide/svelte';
	import NavGroup from './NavGroup.svelte';

	let {
		ref = $bindable(null),
		collapsible = 'icon',
		...restProps
	}: ComponentProps<typeof Sidebar.Root> = $props();
</script>

<Sidebar.Root {collapsible} {...restProps}>
	<!--
	<Sidebar.Header>
		<TeamSwitcher />
	</Sidebar.Header>
	-->
	<Sidebar.Content>
		<NavGroup title="Curso" items={data.main} />
		<NavGroup title="Informacion" items={data.info} />
	</Sidebar.Content>
	<Sidebar.Rail />
</Sidebar.Root>

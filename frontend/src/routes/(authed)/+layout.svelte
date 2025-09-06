<script lang="ts">
	import type { Snippet } from 'svelte';
	import type { LayoutData } from './$types';
	import { CirclePlay, House } from '@lucide/svelte';
	import NavigationBar from '$lib/components/layout/NavigationBar.svelte';
	import type { Icon as IconType } from '@lucide/svelte';

	let { children }: { data: LayoutData; children: Snippet } = $props();

	type NavLink = {
		href: string;
		name: string;
	};

	const links: NavLink[] = [
		{
			name: 'Inicio',
			href: '/home',
		},
		{
			href: '/my-courses',
			name: 'Mis Cursos',
		}
	];
</script>

<div class="hidden md:block">
	<NavigationBar {links} />
</div>

<div class="px-4 pt-12 pb-24 md:px-12 md:pb-12 xl:px-24">
	{@render children()}
</div>

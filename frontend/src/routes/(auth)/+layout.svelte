<script lang="ts">
	import type { Snippet } from 'svelte';
	import type { LayoutData } from './$types';
	import { CirclePlay, House } from '@lucide/svelte';
	import NavigationBar from '$lib/components/layout/NavigationBar.svelte';
	import type { Icon as IconType } from '@lucide/svelte';
	import { authStore } from '$lib/stores/auth.svelte';
	import { UserRole } from '$lib/types/models/course';

	let { children }: { data: LayoutData; children: Snippet } = $props();

	type NavLink = {
		href: string;
		name: string;
	};

	const links: NavLink[] = $derived.by(() => {
		const internalLinks = [
			{
				name: 'Inicio',
				href: '/'
			},
			{
				href: '/my-courses',
				name: 'Mis Cursos'
			}
		];

		if (authStore.user?.role === UserRole.ADMIN || authStore.user?.role === UserRole.INSTRUCTOR) {
			internalLinks.push({
				href: '/admin/courses',
				name: 'Admin'
			});
		}

		return internalLinks;
	});
</script>

<div class="hidden md:block">
	<NavigationBar {links} />
</div>

{@render children()}

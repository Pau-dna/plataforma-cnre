<script lang="ts">
	import { page } from '$app/state';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Avatar from '$lib/components/ui/avatar/index.js';
	import * as Popover from '$lib/components/ui/popover/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Bell, type Icon as IconType } from '@lucide/svelte';
	import DropdownAvatar from './DropdownAvatar.svelte';

	type NavLink = {
		href: string;
		name: string;
	};

	type Props = {
		links: NavLink[];
	};

	const { links }: Props = $props();

	let activeRoute = $derived(page?.url?.pathname);
</script>

<nav class="flex items-center justify-between border-b md:px-12 xl:px-24">
	<div class="flex items-center gap-x-20">
		<div class="flex shrink-0 text-lg font-semibold">
			<a href="/home" class="flex shrink-0 gap-4 py-0">
				<img src="/images/logo.png" alt="CNRE Logo" class="h-14 w-auto" />
			</a>
		</div>

		<div class="flex items-center gap-x-6">
			{#each links as link}
				<a href={link?.href || ''}>
					<div
						class={`relative flex items-center justify-center gap-2 border-b-2 px-2 py-6 text-base font-medium transition-colors duration-200  ${
							activeRoute === link.href
								? 'border-b-2 border-sky-500'
								: 'border-transparent text-primary/80 hover:border-gray-200 hover:text-primary'
						}`}
					>
						<span>{link.name}</span>
					</div>
				</a>
			{/each}
		</div>
	</div>

	<div class="py-4">
		<DropdownAvatar />
	</div>
</nav>

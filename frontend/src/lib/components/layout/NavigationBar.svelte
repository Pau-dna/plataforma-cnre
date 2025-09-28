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

<nav class="flex items-center justify-between border-b px-4 md:px-12 xl:px-24">
	<div class="flex items-center gap-x-8 md:gap-x-20">
		<div class="flex shrink-0 text-lg font-semibold">
			<a href="/home" class="flex shrink-0 gap-4 py-0">
				<img src="/images/logo.png" alt="Logo CNRE" class="h-12 w-auto md:h-14" />
			</a>
		</div>

		<div class="hidden items-center gap-x-4 sm:flex md:gap-x-6">
			{#each links as link}
				<a href={link?.href || ''}>
					<div
						class={`relative flex items-center justify-center gap-2 border-b-2 px-2 py-6 text-sm font-medium transition-colors duration-200 md:text-base  ${
							activeRoute === link.href
								? 'border-b-2 border-sky-500'
								: 'text-primary/80 hover:text-primary border-transparent hover:border-gray-200'
						}`}
					>
						<span>{link.name}</span>
					</div>
				</a>
			{/each}
		</div>
	</div>

	<div class="flex max-h-max">
		<DropdownAvatar />
	</div>
</nav>

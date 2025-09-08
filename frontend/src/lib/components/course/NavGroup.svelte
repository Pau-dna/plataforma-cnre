<script lang="ts">
	import * as Collapsible from '$lib/components/ui/collapsible/index.js';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import ChevronRightIcon from '@lucide/svelte/icons/chevron-right';

	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import FolderIcon from '@lucide/svelte/icons/folder';
	import ForwardIcon from '@lucide/svelte/icons/forward';
	import Trash2Icon from '@lucide/svelte/icons/trash-2';
	import type { Icon as IconType } from '@lucide/svelte';

	type NavItem = {
		title: string;
		url: string;
		// this should be `Component` after @lucide/svelte updates types
		// eslint-disable-next-line @typescript-eslint/no-explicit-any
		icon: typeof IconType;
		active?: boolean;
		items?: {
			title: string;
			url: string;
		}[];
	};

	type Props = {
		title: string;
		items: NavItem[];
	};

	const { title, items }: Props = $props();
</script>

<Sidebar.Group>
	<Sidebar.GroupLabel>{title}</Sidebar.GroupLabel>
	<Sidebar.Menu>
		{#each items as item (item.title)}
			{#if (item.items?.length || 0) > 0}
				{@render nestedItem(item)}
			{:else}
				{@render simpleItem(item)}
			{/if}
		{/each}
	</Sidebar.Menu>
</Sidebar.Group>

{#snippet simpleItem(item: NavItem)}
	<Sidebar.MenuItem>
		<Sidebar.MenuButton>
			{#snippet child({ props })}
				<a href={item.url} {...props}>
					<item.icon />
					<span>{item.title}</span>
				</a>
			{/snippet}
		</Sidebar.MenuButton>
	</Sidebar.MenuItem>
{/snippet}

{#snippet nestedItem(item: NavItem)}
	<Collapsible.Root open={item.active} class="group/collapsible">
		{#snippet child({ props })}
			<Sidebar.MenuItem {...props}>
				<Collapsible.Trigger>
					{#snippet child({ props })}
						<Sidebar.MenuButton {...props} tooltipContent={item.title}>
							{#if item.icon}
								<item.icon />
							{/if}
							<span>{item.title}</span>
							<ChevronRightIcon
								class="ml-auto transition-transform duration-200 group-data-[state=open]/collapsible:rotate-90"
							/>
						</Sidebar.MenuButton>
					{/snippet}
				</Collapsible.Trigger>
				<Collapsible.Content>
					<Sidebar.MenuSub>
						{#each item.items ?? [] as subItem (subItem.title)}
							<Sidebar.MenuSubItem>
								<Sidebar.MenuSubButton>
									{#snippet child({ props })}
										<a href={subItem.url} {...props}>
											<span>{subItem.title}</span>
										</a>
									{/snippet}
								</Sidebar.MenuSubButton>
							</Sidebar.MenuSubItem>
						{/each}
					</Sidebar.MenuSub>
				</Collapsible.Content>
			</Sidebar.MenuItem>
		{/snippet}
	</Collapsible.Root>
{/snippet}

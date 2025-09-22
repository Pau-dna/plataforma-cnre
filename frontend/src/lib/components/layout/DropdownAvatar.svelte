<script lang="ts">
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import * as Avatar from '$lib/components/ui/avatar/index.js';
	import { authStore } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { UserIcon, MailIcon, LogOutIcon } from '@lucide/svelte';

	function getInitials(fullname: string) {
		const names = fullname.split(' ');
		const initials = names.map((name) => name.charAt(0).toUpperCase());
		return initials.slice(0, 3).join('');
	}

	function shortName(fullname: string) {
		const names = fullname.split(' ');
		if (names.length === 0) return '';
		if (names.length === 1) return names[0];
		return `${names[0]} ${names[names.length - 2]}`;
	}
</script>

<DropdownMenu.Root>
	<DropdownMenu.Trigger>
		<Avatar.Root class="size-10">
			<Avatar.Fallback>
				{getInitials(authStore.user?.fullname || '')}
			</Avatar.Fallback>
			<Avatar.Image
				src={authStore.user?.avatar_url || '/images/default-avatar.png'}
				alt="User Avatar"
			/>
		</Avatar.Root>
	</DropdownMenu.Trigger>
	<DropdownMenu.Content class="w-64">
		<DropdownMenu.Label class="p-0 font-normal">
			<div class="flex items-center gap-3 px-2 py-3">
				<Avatar.Root class="size-12">
					<Avatar.Fallback>
						{getInitials(authStore.user?.fullname || '')}
					</Avatar.Fallback>
					<Avatar.Image
						src={authStore.user?.avatar_url || '/images/default-avatar.png'}
						alt="User Avatar"
					/>
				</Avatar.Root>
				<div class="grid flex-1 text-left text-sm leading-tight">
					<span class="truncate font-semibold">{shortName(authStore.user?.fullname || 'Usuario')}</span>
					<div class="mt-1 flex items-center gap-1">
						<MailIcon class="text-muted-foreground size-3" />
						<span class="text-muted-foreground truncate text-xs">{authStore.user?.email || ''}</span
						>
					</div>
					{#if authStore.user?.role}
						<div class="mt-1 flex items-center gap-1">
							<UserIcon class="text-muted-foreground size-3" />
							<Badge variant="secondary" class="h-4 px-1 text-xs">
								{authStore.user.role}
							</Badge>
						</div>
					{/if}
				</div>
			</div>
		</DropdownMenu.Label>
		<DropdownMenu.Separator />
		<DropdownMenu.Group>
			<a href="/logout" data-sveltekit-reload>
				<DropdownMenu.Item class="cursor-pointer">
					<LogOutIcon class="size-4" />
					Cerrar Sesión
				</DropdownMenu.Item>
			</a>
		</DropdownMenu.Group>
	</DropdownMenu.Content>
</DropdownMenu.Root>

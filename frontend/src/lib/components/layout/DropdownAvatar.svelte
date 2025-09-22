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
		return initials.slice(0,3).join('');
	}

	async function handleLogout() {
		await goto('/logout');
	}

</script>

<DropdownMenu.Root>
	<DropdownMenu.Trigger>
		<Avatar.Root class="size-10">
			<Avatar.Fallback>
				{getInitials(authStore.user?.fullname || "")}
			</Avatar.Fallback>
			<Avatar.Image src={authStore.user?.avatar_url || '/images/default-avatar.png'} alt="User Avatar" />
		</Avatar.Root>
	</DropdownMenu.Trigger>
	<DropdownMenu.Content class="w-64">
		<DropdownMenu.Label class="p-0 font-normal">
			<div class="flex items-center gap-3 px-2 py-3">
				<Avatar.Root class="size-12">
					<Avatar.Fallback>
						{getInitials(authStore.user?.fullname || "")}
					</Avatar.Fallback>
					<Avatar.Image src={authStore.user?.avatar_url || '/images/default-avatar.png'} alt="User Avatar" />
				</Avatar.Root>
				<div class="grid flex-1 text-left text-sm leading-tight">
					<span class="truncate font-semibold">{authStore.user?.fullname || "Usuario"}</span>
					<div class="flex items-center gap-1 mt-1">
						<MailIcon class="size-3 text-muted-foreground" />
						<span class="truncate text-xs text-muted-foreground">{authStore.user?.email || ""}</span>
					</div>
					{#if authStore.user?.role}
						<div class="flex items-center gap-1 mt-1">
							<UserIcon class="size-3 text-muted-foreground" />
							<Badge variant="secondary" class="text-xs h-4 px-1">
								{authStore.user.role}
							</Badge>
						</div>
					{/if}
				</div>
			</div>
		</DropdownMenu.Label>
		<DropdownMenu.Separator />
		<DropdownMenu.Group>
			<DropdownMenu.Item class="cursor-pointer" onclick={handleLogout}>
				<LogOutIcon class="size-4" />
				Cerrar Sesión
			</DropdownMenu.Item>
		</DropdownMenu.Group>
	</DropdownMenu.Content>
</DropdownMenu.Root>

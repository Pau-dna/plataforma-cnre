<script lang="ts">
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import * as Avatar from '$lib/components/ui/avatar/index.js';
	import { authStore } from '$lib/stores/auth.svelte';


	function getInitials(fullname: string) {
		const names = fullname.split(' ');
		const initials = names.map((name) => name.charAt(0).toUpperCase());
		return initials.slice(0,3).join('');
	}

</script>

{JSON.stringify(authStore.user)}
<DropdownMenu.Root>
	<DropdownMenu.Trigger>
		<Avatar.Root class="size-10">
			<Avatar.Fallback>
				{getInitials(authStore.user?.fullname || "")}
			</Avatar.Fallback>
			<Avatar.Image src={authStore.user?.avatar_url || '/images/default-avatar.png'} alt="User Avatar" />
		</Avatar.Root>
	</DropdownMenu.Trigger>
	<DropdownMenu.Content>
		<DropdownMenu.Group>
			<DropdownMenu.Label>Mi cuenta</DropdownMenu.Label>
			<DropdownMenu.Separator />
			<DropdownMenu.Item>Cerrar Sesión</DropdownMenu.Item>
		</DropdownMenu.Group>
	</DropdownMenu.Content>
</DropdownMenu.Root>

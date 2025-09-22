<script lang="ts">
	import { goto, invalidateAll } from '$app/navigation';
	import { onMount } from 'svelte';

	onMount(async () => {
		await invalidateAll();
		await goto('/login');
	});
</script>

<h1>Logging out...</h1>

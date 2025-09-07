<script lang="ts">
	import LoginCard from '$lib/components/login/LoginCard.svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
</script>

<div class="flex min-h-screen h-screen items-center justify-center">
	<LoginCard />
</div>

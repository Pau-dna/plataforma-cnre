<script lang="ts">
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	const modulo = $derived(data.modules.find(m => m.id === data.moduleID))
</script>

<div class="flex flex-col gap-y-4">
	<h1 class="text-3xl font-bold">{modulo.title}</h1>
	<p class="text-muted-foreground">{modulo.description}</p>
</div>
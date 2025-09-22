<script lang="ts">
	import { ArrowLeft } from '@lucide/svelte';
	import { Button } from '$lib/components/ui/button/index.js';
	import { browser } from '$app/environment';

	type Props = {
		href?: string;
	};

	const { href }: Props = $props();

	const hasBack = $derived(browser && window?.history && window.history.length > 2);
	const backJavascript = 'javascript:history.back()';
	const backLink = $derived(href || (hasBack ? backJavascript : backJavascript));
</script>

<div class="block h-14 md:hidden">
	<Button
		variant="ghost"
		size="icon"
		class="dark:bg-primary-foreground fixed top-10 z-50 flex size-14 items-center justify-center rounded-full bg-neutral-200 transition-all duration-200 hover:bg-neutral-100 active:scale-95  md:hidden dark:border"
		href={backLink}
	>
		<ArrowLeft class="h-5 w-5" />
		<span class="sr-only">Volver</span>
	</Button>
</div>

<Button
	variant="link"
	class="hidden max-w-max items-center justify-center px-0 transition-all duration-200 md:flex"
	href={backLink}
>
	<ArrowLeft class="h-5 w-5" />
	<span>Volver</span>
</Button>

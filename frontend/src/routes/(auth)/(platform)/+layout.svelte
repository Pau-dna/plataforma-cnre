<script lang="ts">
    import type { LayoutProps } from './$types';

    let { data, children }: LayoutProps = $props();
</script>

<div class="px-4 pb-24 pt-12 md:px-12 md:pb-12 xl:px-24 flex flex-col gap-y-4">
	{@render children()}
</div>
<script lang="ts">
    import type { LayoutProps } from './$types';

    let { data, children }: LayoutProps = $props();
</script>

<div class="pb-24 md:pb-12">
	{@render children()}
</div>
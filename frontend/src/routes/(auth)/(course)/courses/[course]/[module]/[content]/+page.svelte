<script lang="ts">
    import type { PageProps } from './$types';

    let { data }: PageProps = $props();

    let content = $state(data.content);
</script>

{JSON.stringify(content)}
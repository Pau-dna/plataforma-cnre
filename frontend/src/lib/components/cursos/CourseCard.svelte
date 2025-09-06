<script lang="ts">
	import * as Card from "$lib/components/ui/card/index.js";

	type Props = {
		title: string;
		description: string;
		userCount: number;
		duration: string;
	};

	const { title, description, userCount, duration }: Props = $props();
</script>


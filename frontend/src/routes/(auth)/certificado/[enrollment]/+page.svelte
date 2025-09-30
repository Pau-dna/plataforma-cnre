<script lang="ts">
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
</script>

Se certifica que {data?.enrollment?.user?.fullname || "N/A"} ha completado el curso {data?.enrollment?.course?.title || "N/A"}. (en proceso mejorar el diseño xd)
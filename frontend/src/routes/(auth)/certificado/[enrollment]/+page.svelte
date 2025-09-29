<script lang="ts">
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
</script>

Se certifica que {data.enrollment.user.name} ha completado el curso {data.enrollment.course.title}. (en proceso mejorar el diseño xd)
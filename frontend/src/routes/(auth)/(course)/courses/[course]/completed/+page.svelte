<script lang="ts">
	import CourseCompletion from '$lib/components/course/CourseCompletion.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Share2, ArrowLeft } from '@lucide/svelte';
	import type { PageData } from './$types';
	import { goto } from '$app/navigation';

	let { data }: { data: PageData } = $props();

	const handleDownloadCertificate = () => {
		goto(`/certificado/${data.enrollment.id}`);
	};

	const handleShareCertificate = () => {
		const certificateUrl = `/certificado/${data.enrollment.id}`;
		if (navigator.share) {
			navigator.share({
				title: `Certificado del curso ${data.course.title}`,
				text: `¡He completado exitosamente el curso "${data.course.title}"!`,
				url: `${window.location.origin}${certificateUrl}`
			});
		} else {
			// Fallback: copy to clipboard
			navigator.clipboard.writeText(`${window.location.origin}${certificateUrl}`);
			// Could add a toast notification here
		}
	};

	const handleBackToCourse = () => {
		goto(`/courses/${data.course.id}`);
	};
</script>

<svelte:head>
	<title>¡Felicidades! - {data.course.title} | CNRE</title>
	<meta name="description" content="Has completado exitosamente el curso {data.course.title}" />
</svelte:head>

<div class="container mx-auto max-w-4xl px-4 py-8">
	<!-- Back navigation -->
	<div class="mb-6">
		<Button variant="ghost" onclick={handleBackToCourse} class="gap-2">
			<ArrowLeft class="h-4 w-4" />
			Volver al curso
		</Button>
	</div>

	<!-- Main completion card -->
	<div class="mb-8">
		<CourseCompletion
			courseName={data.course.title}
			completedAt={data.enrollment.completed_at || ''}
			progressPercentage={data.enrollment.progress}
			onDownloadCertificate={handleDownloadCertificate}
		/>
	</div>

	<!-- Additional actions -->
	<div class="bg-card rounded-lg border p-6">
		<h2 class="mb-4 text-xl font-semibold">¡Comparte tu logro!</h2>
		<p class="text-muted-foreground mb-4">
			Comparte tu certificado con otros y demuestra tu nueva competencia en {data.course.title}.
		</p>

		<div class="flex flex-wrap gap-3">
			<Button onclick={handleDownloadCertificate} class="gap-2">Ver certificado</Button>

			<Button variant="outline" onclick={handleShareCertificate} class="gap-2">
				<Share2 class="h-4 w-4" />
				Compartir certificado
			</Button>
		</div>
	</div>

	<!-- Course summary (optional enhancement) -->
	<div class="bg-card mt-8 rounded-lg border p-6">
		<h3 class="mb-3 text-lg font-semibold">Resumen del curso</h3>
		<div class="grid gap-4 md:grid-cols-2">
			<div>
				<p class="text-muted-foreground text-sm">Nombre del curso</p>
				<p class="font-medium">{data.course.title}</p>
			</div>
			<div>
				<p class="text-muted-foreground text-sm">Progreso completado</p>
				<p class="font-medium">{data.enrollment.progress}%</p>
			</div>
			{#if data.course.description}
				<div class="md:col-span-2">
					<p class="text-muted-foreground text-sm">Descripción</p>
					<p class="text-sm">{data.course.description}</p>
				</div>
			{/if}
		</div>
	</div>
</div>

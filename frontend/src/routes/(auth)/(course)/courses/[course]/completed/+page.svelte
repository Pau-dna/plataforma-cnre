<script lang="ts">
	import type { PageProps } from './$types';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { Award, CheckCircle, Share2, ExternalLink, Trophy } from '@lucide/svelte';
	import { goto } from '$app/navigation';

	let { data }: PageProps = $props();

	// Format the completion date in Spanish
	const completionDate = data.enrollment.completed_at
		? new Date(data.enrollment.completed_at).toLocaleDateString('es-ES', {
				year: 'numeric',
				month: 'long',
				day: 'numeric'
			})
		: 'Fecha no disponible';

	// Certificate URL
	const certificateUrl = `/certificado/${data.enrollment.id}`;

	// Function to handle certificate viewing
	function viewCertificate() {
		goto(certificateUrl);
	}

	// Function to handle sharing (future implementation could include native sharing)
	function shareCertificate() {
		if (navigator.share) {
			navigator.share({
				title: `Certificado del curso: ${data.enrollment.course?.title}`,
				text: `¡He completado exitosamente el curso "${data.enrollment.course?.title}"!`,
				url: window.location.origin + certificateUrl
			});
		} else {
			// Fallback: copy to clipboard
			navigator.clipboard.writeText(window.location.origin + certificateUrl);
			// Note: In a real implementation, you'd show a toast notification here
		}
	}
</script>

<div class="container mx-auto max-w-4xl p-6">
	<!-- Main congratulations card -->
	<Card.Root class="mb-6 border-2 border-green-200 bg-gradient-to-br from-green-50 to-emerald-50">
		<Card.Content class="p-8">
			<div class="text-center" role="main" aria-labelledby="completion-title">
				<!-- Main trophy icon -->
				<div
					class="mx-auto mb-6 flex h-20 w-20 items-center justify-center rounded-full bg-gradient-to-r from-yellow-400 to-orange-500 shadow-lg"
					aria-hidden="true"
				>
					<Trophy class="h-10 w-10 text-white" />
				</div>

				<!-- Congratulations title -->
				<h1 id="completion-title" class="mb-4 text-3xl font-bold text-green-800 md:text-4xl">
					¡Felicidades!
				</h1>

				<!-- Success message -->
				<div class="mb-6 flex items-center justify-center gap-2">
					<CheckCircle class="h-6 w-6 text-green-600" aria-hidden="true" />
					<h2 class="text-xl font-semibold text-green-700">Has completado exitosamente el curso</h2>
				</div>

				<!-- Course title -->
				<p class="mb-4 text-lg text-gray-700">
					<strong class="text-green-800">{data.enrollment.course?.title}</strong>
				</p>

				<!-- Completion date -->
				<p class="mb-8 text-gray-600">
					Completado el {completionDate}
				</p>

				<!-- Progress bar showing 100% -->
				<div class="mx-auto mb-8 max-w-md">
					<div class="mb-2 flex items-center justify-between">
						<span class="text-sm font-medium text-green-700">Progreso del curso</span>
						<span class="text-sm font-bold text-green-800">100%</span>
					</div>
					<div
						class="h-3 w-full overflow-hidden rounded-full bg-green-200"
						role="progressbar"
						aria-valuenow="100"
						aria-valuemin="0"
						aria-valuemax="100"
						aria-label="Progreso del curso completado al 100%"
					>
						<div
							class="h-full rounded-full bg-gradient-to-r from-green-500 to-emerald-600 transition-all duration-1000"
							style="width: 100%"
						></div>
					</div>
				</div>
			</div>
		</Card.Content>
	</Card.Root>

	<!-- Certificate section -->
	<Card.Root class="border-amber-200 bg-gradient-to-r from-amber-50 to-yellow-50">
		<Card.Header class="text-center">
			<div
				class="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-amber-100"
				aria-hidden="true"
			>
				<Award class="h-8 w-8 text-amber-600" />
			</div>
			<Card.Title class="text-2xl text-amber-800">Tu Certificado</Card.Title>
			<Card.Description class="text-amber-700">
				¡Has ganado tu certificado de finalización! Puedes verlo y compartirlo con otros.
			</Card.Description>
		</Card.Header>

		<Card.Content class="text-center">
			<p class="mb-6 text-gray-600">
				Tu certificado está listo y disponible para descargar o compartir. Demuestra tu nuevo
				conocimiento y habilidades.
			</p>

			<!-- Action buttons -->
			<div class="flex flex-col items-center gap-4 sm:flex-row sm:justify-center">
				<Button
					onclick={viewCertificate}
					class="flex items-center gap-2 bg-amber-600 px-6 py-3 text-white hover:bg-amber-700"
					aria-label="Ver certificado de finalización del curso"
				>
					<Award class="h-5 w-5" />
					Ver Certificado
					<ExternalLink class="h-4 w-4" />
				</Button>

				<Button
					variant="outline"
					onclick={shareCertificate}
					class="flex items-center gap-2 border-amber-300 px-6 py-3 text-amber-700 hover:bg-amber-100"
					aria-label="Compartir certificado con otros"
				>
					<Share2 class="h-5 w-5" />
					Compartir Certificado
				</Button>
			</div>
		</Card.Content>
	</Card.Root>

	<!-- Additional encouragement -->
	<div class="mt-8 text-center">
		<p class="text-gray-600">
			¡Excelente trabajo! Continúa aprendiendo y explorando nuevos cursos para seguir creciendo.
		</p>
	</div>
</div>

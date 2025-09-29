<script lang="ts">
	import type { PageProps } from './$types';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { Award, Share2, ExternalLink, CheckCircle } from '@lucide/svelte';
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

<div class="container mx-auto max-w-2xl p-6">
	<div class="text-center">
		<!-- Simple completion checkmark -->
		<div
			class="border-border bg-background mx-auto mb-6 flex h-16 w-16 items-center justify-center rounded-full border"
		>
			<CheckCircle class="text-foreground h-8 w-8" />
		</div>

		<!-- Simple title -->
		<h1 class="text-foreground mb-2 text-2xl font-semibold">Curso completado</h1>

		<!-- Course title -->
		<p class="text-muted-foreground mb-4 text-lg">
			{data.enrollment.course?.title}
		</p>

		<!-- Completion date -->
		<p class="text-muted-foreground mb-8 text-sm">
			Finalizado el {completionDate}
		</p>

		<!-- Certificate section -->
		<Card.Root class="border-border bg-card mb-6">
			<Card.Content class="p-6">
				<div class="mb-4 flex items-center justify-center gap-2">
					<Award class="text-muted-foreground h-5 w-5" />
					<span class="text-card-foreground text-sm font-medium">Tu certificado está listo</span>
				</div>

				<div class="flex flex-col items-center gap-3 sm:flex-row sm:justify-center">
					<Button
						onclick={viewCertificate}
						variant="default"
						class="flex items-center gap-2"
						aria-label="Ver certificado de finalización del curso"
					>
						<Award class="h-4 w-4" />
						Ver Certificado
						<ExternalLink class="h-3 w-3" />
					</Button>

					<Button
						variant="outline"
						onclick={shareCertificate}
						class="flex items-center gap-2"
						aria-label="Compartir certificado con otros"
					>
						<Share2 class="h-4 w-4" />
						Compartir
					</Button>
				</div>
			</Card.Content>
		</Card.Root>

		<!-- Simple encouragement message -->
		<p class="text-muted-foreground text-sm">
			¡Excelente trabajo! Continúa explorando nuevos cursos.
		</p>
	</div>
</div>

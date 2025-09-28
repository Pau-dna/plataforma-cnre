<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { CheckCircle, Award, Download } from '@lucide/svelte';

	interface Props {
		courseName: string;
		completedAt: string;
		progressPercentage: number;
		onDownloadCertificate?: () => void;
	}

	let { courseName, completedAt, progressPercentage, onDownloadCertificate }: Props = $props();

	const completionDate = new Date(completedAt).toLocaleDateString('es-ES', {
		year: 'numeric',
		month: 'long',
		day: 'numeric'
	});
</script>

<div class="rounded-xl border border-green-200 bg-gradient-to-r from-green-50 to-emerald-50 p-6">
	<div class="flex items-start gap-4">
		<div class="flex-shrink-0">
			<div class="rounded-full bg-green-100 p-3">
				<Award class="h-8 w-8 text-green-600" />
			</div>
		</div>

		<div class="flex-1">
			<div class="mb-2 flex items-center gap-2">
				<CheckCircle class="h-5 w-5 text-green-600" />
				<h3 class="text-lg font-semibold text-green-800">¡Felicidades! Has completado el curso</h3>
			</div>

			<p class="mb-3 text-green-700">
				Has terminado <strong>{courseName}</strong> el {completionDate}
			</p>

			<div class="flex items-center gap-4">
				<div class="flex items-center gap-2">
					<div class="h-2 w-12 overflow-hidden rounded-full bg-green-200">
						<div
							class="h-full rounded-full bg-green-600 transition-all duration-500"
							style="width: {progressPercentage}%"
						></div>
					</div>
					<span class="text-sm font-medium text-green-700">{progressPercentage}% completado</span>
				</div>

				{#if onDownloadCertificate}
					<Button
						variant="outline"
						size="sm"
						class="border-green-300 text-green-700 hover:bg-green-100"
						onclick={onDownloadCertificate}
					>
						<Download class="mr-2 h-4 w-4" />
						Descargar certificado
					</Button>
				{/if}
			</div>
		</div>
	</div>
</div>

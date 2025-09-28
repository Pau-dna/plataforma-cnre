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

<div class="bg-gradient-to-r from-green-50 to-emerald-50 border border-green-200 rounded-xl p-6">
	<div class="flex items-start gap-4">
		<div class="flex-shrink-0">
			<div class="bg-green-100 rounded-full p-3">
				<Award class="h-8 w-8 text-green-600" />
			</div>
		</div>
		
		<div class="flex-1">
			<div class="flex items-center gap-2 mb-2">
				<CheckCircle class="h-5 w-5 text-green-600" />
				<h3 class="text-lg font-semibold text-green-800">¡Felicidades! Has completado el curso</h3>
			</div>
			
			<p class="text-green-700 mb-3">
				Has terminado <strong>{courseName}</strong> el {completionDate}
			</p>
			
			<div class="flex items-center gap-4">
				<div class="flex items-center gap-2">
					<div class="w-12 h-2 bg-green-200 rounded-full overflow-hidden">
						<div 
							class="h-full bg-green-600 rounded-full transition-all duration-500"
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
						<Download class="h-4 w-4 mr-2" />
						Descargar certificado
					</Button>
				{/if}
			</div>
		</div>
	</div>
</div>
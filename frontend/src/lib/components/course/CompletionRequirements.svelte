<script lang="ts">
	import { CheckCircle, Circle } from '@lucide/svelte';

	interface Props {
		modules: Array<{
			title: string;
			contents: Array<{ title: string; isCompleted: boolean }>;
			evaluations: Array<{ title: string; hasPassed: boolean }>;
			progressPercentage: number;
		}>;
	}

	let { modules }: Props = $props();

	// Calculate requirements
	const requirements = $derived(() => {
		let totalContent = 0;
		let completedContent = 0;
		let totalEvaluations = 0;
		let passedEvaluations = 0;

		modules.forEach(module => {
			totalContent += module.contents?.length || 0;
			completedContent += module.contents?.filter(c => c.isCompleted).length || 0;
			totalEvaluations += module.evaluations?.length || 0;
			passedEvaluations += module.evaluations?.filter(e => e.hasPassed).length || 0;
		});

		return {
			totalContent,
			completedContent,
			totalEvaluations,
			passedEvaluations,
			allContentCompleted: completedContent === totalContent,
			allEvaluationsPassed: passedEvaluations === totalEvaluations
		};
	});
</script>

<div class="bg-gray-50 border border-gray-200 rounded-lg p-4">
	<h3 class="text-sm font-medium text-gray-800 mb-3">Requisitos para completar el curso</h3>
	
	<div class="space-y-2">
		<!-- Content Requirement -->
		<div class="flex items-center gap-2">
			{#if requirements.allContentCompleted}
				<CheckCircle class="h-4 w-4 text-green-600" />
				<span class="text-sm text-green-700">
					Todo el contenido visto ({requirements.completedContent}/{requirements.totalContent})
				</span>
			{:else}
				<Circle class="h-4 w-4 text-gray-400" />
				<span class="text-sm text-gray-600">
					Ver todo el contenido ({requirements.completedContent}/{requirements.totalContent})
				</span>
			{/if}
		</div>

		<!-- Evaluations Requirement -->
		{#if requirements.totalEvaluations > 0}
			<div class="flex items-center gap-2">
				{#if requirements.allEvaluationsPassed}
					<CheckCircle class="h-4 w-4 text-green-600" />
					<span class="text-sm text-green-700">
						Todas las evaluaciones aprobadas ({requirements.passedEvaluations}/{requirements.totalEvaluations})
					</span>
				{:else}
					<Circle class="h-4 w-4 text-gray-400" />
					<span class="text-sm text-gray-600">
						Aprobar todas las evaluaciones ({requirements.passedEvaluations}/{requirements.totalEvaluations})
					</span>
				{/if}
			</div>
		{/if}
	</div>

	{#if requirements.allContentCompleted && (requirements.totalEvaluations === 0 || requirements.allEvaluationsPassed)}
		<div class="mt-3 pt-3 border-t border-gray-200">
			<div class="flex items-center gap-2">
				<CheckCircle class="h-4 w-4 text-green-600" />
				<span class="text-sm font-medium text-green-700">¡Listo! Has cumplido todos los requisitos</span>
			</div>
		</div>
	{/if}
</div>
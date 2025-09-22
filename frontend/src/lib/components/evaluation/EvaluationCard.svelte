<script lang="ts">
	import type { Evaluation } from '$lib/types';
	import { EvaluationController } from '$lib/controllers';
	import { toast } from 'svelte-sonner';
	import { Edit, Trash, ChevronUp, ChevronDown, FileQuestion } from 'lucide/svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import Card from '$lib/components/ui/card/card.svelte';
	import CardHeader from '$lib/components/ui/card/card-header.svelte';
	import CardTitle from '$lib/components/ui/card/card-title.svelte';
	import CardContent from '$lib/components/ui/card/card-content.svelte';

	interface Props {
		evaluation: Evaluation;
		courseId: number;
		onupdate?: (evaluation: Evaluation) => void;
		ondelete?: (evaluation: Evaluation) => void;
		onmoveup?: (evaluation: Evaluation) => void;
		onmovedown?: (evaluation: Evaluation) => void;
		canMoveUp?: boolean;
		canMoveDown?: boolean;
	}

	let {
		evaluation,
		courseId,
		onupdate,
		ondelete,
		onmoveup,
		onmovedown,
		canMoveUp = false,
		canMoveDown = false
	}: Props = $props();

	const evaluationController = new EvaluationController();

	async function handleDeleteEvaluation() {
		if (!confirm('¿Estás seguro de que quieres eliminar esta evaluación?')) {
			return;
		}

		try {
			await evaluationController.deleteEvaluation(evaluation.id);
			toast.success('Evaluación eliminada exitosamente.');
			ondelete?.(evaluation);
		} catch (error) {
			console.error('Error deleting evaluation:', error);
			toast.error('Error al eliminar la evaluación.', {
				description: error instanceof Error ? error.message : String(error)
			});
		}
	}

	function handleMoveUp() {
		onmoveup?.(evaluation);
	}

	function handleMoveDown() {
		onmovedown?.(evaluation);
	}
</script>

<Card class="w-full">
	<CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
		<div class="flex items-center gap-2">
			<FileQuestion class="h-5 w-5 text-purple-600" />
			<CardTitle class="text-lg font-semibold">{evaluation.title}</CardTitle>
		</div>
		<div class="flex items-center gap-2">
			{#if canMoveUp}
				<Button variant="outline" size="sm" onclick={handleMoveUp}>
					<ChevronUp class="h-4 w-4" />
				</Button>
			{/if}
			{#if canMoveDown}
				<Button variant="outline" size="sm" onclick={handleMoveDown}>
					<ChevronDown class="h-4 w-4" />
				</Button>
			{/if}
			<Button
				variant="outline"
				size="sm"
				href="/admin/courses/{courseId}/{evaluation.module_id}/evaluations/{evaluation.id}"
			>
				<Edit class="h-4 w-4" />
				<span>Gestionar</span>
			</Button>
			<Button variant="destructive" size="sm" onclick={handleDeleteEvaluation}>
				<Trash class="h-4 w-4" />
			</Button>
		</div>
	</CardHeader>
	<CardContent>
		{#if evaluation.description}
			<p class="text-muted-foreground text-sm mb-2">{evaluation.description}</p>
		{/if}
		<div class="flex flex-wrap gap-2 text-sm text-muted-foreground">
			<span>Preguntas: {evaluation.question_count}</span>
			<span>•</span>
			<span>Opciones: {evaluation.answer_options_count}</span>
			<span>•</span>
			<span>Puntaje mínimo: {evaluation.passing_score}%</span>
			{#if evaluation.time_limit}
				<span>•</span>
				<span>Tiempo límite: {evaluation.time_limit} minutos</span>
			{/if}
			{#if evaluation.max_attempts}
				<span>•</span>
				<span>Intentos máximos: {evaluation.max_attempts}</span>
			{/if}
		</div>
	</CardContent>
</Card>
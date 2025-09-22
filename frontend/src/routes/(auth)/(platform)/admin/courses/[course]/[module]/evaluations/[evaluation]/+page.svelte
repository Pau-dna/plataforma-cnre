<script lang="ts">
	import type { PageProps } from './$types';
	import type { Question } from '$lib/types';
	import { Plus, Settings } from '@lucide/svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import Back from '$lib/components/kit/Back.svelte';
	import QuestionCard from '$lib/components/evaluation/QuestionCard.svelte';

	let { data }: PageProps = $props();

	const evaluation = $state(data.evaluation);
	let questions = $state(data.questions);
	const courseId = data.courseId;
	const moduleId = data.moduleId;

	function handleQuestionUpdate(updated: Question) {
		const index = questions.findIndex((q) => q.id === updated.id);
		if (index !== -1) {
			questions[index] = updated;
		}
	}

	function handleQuestionDelete(deleted: Question) {
		questions = questions.filter((q) => q.id !== deleted.id);
	}

	function handleQuestionAdd(newQuestion: Question) {
		questions = [...questions, newQuestion];
	}
</script>

<Back href="/admin/courses/{courseId}/{moduleId}" />

<div class="flex flex-col gap-6">
	<div class="flex items-start justify-between">
		<div class="flex flex-col gap-2">
			<h1 class="text-h1">Gestión de Evaluación</h1>
			<h2 class="text-2xl font-semibold">{evaluation.title}</h2>
			{#if evaluation.description}
				<p class="text-muted-foreground">{evaluation.description}</p>
			{/if}
			<div class="flex flex-wrap gap-2 text-sm text-muted-foreground">
				<span>Preguntas: {evaluation.question_count}</span>
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
		</div>
		<div class="flex gap-2">
			<Button
				href="/admin/courses/{courseId}/{moduleId}/evaluations/{evaluation.id}/new-question"
				class="bg-green-500 hover:bg-green-600"
			>
				<Plus class="h-4 w-4" />
				<span>Nueva Pregunta</span>
			</Button>
			<Button
				variant="outline"
				href="/admin/courses/{courseId}/{moduleId}/evaluations/{evaluation.id}/settings"
			>
				<Settings class="h-4 w-4" />
				<span>Configuración</span>
			</Button>
		</div>
	</div>

	<div class="flex flex-col gap-4">
		<h3 class="text-xl font-semibold">Preguntas de la Evaluación</h3>

		{#if questions.length > 0}
			<div class="grid grid-cols-1 gap-4">
				{#each questions as question, index (question.id)}
					<QuestionCard
						{question}
						questionNumber={index + 1}
						evaluationId={evaluation.id}
						{courseId}
						{moduleId}
						onupdate={handleQuestionUpdate}
						ondelete={handleQuestionDelete}
					/>
				{/each}
			</div>
		{:else}
			<div class="flex flex-col items-center justify-center py-12 text-center border-2 border-dashed border-gray-300 rounded-lg">
				<p class="text-muted-foreground mb-4 text-lg">
					Esta evaluación no tiene preguntas aún
				</p>
				<p class="text-muted-foreground mb-6 text-sm">
					Comienza agregando preguntas para que los estudiantes puedan realizar la evaluación
				</p>
				<Button
					href="/admin/courses/{courseId}/{moduleId}/evaluations/{evaluation.id}/new-question"
					class="bg-green-500 hover:bg-green-600"
				>
					<Plus class="h-4 w-4" />
					<span>Crear Primera Pregunta</span>
				</Button>
			</div>
		{/if}
	</div>
</div>
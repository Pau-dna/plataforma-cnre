<script lang="ts">
	import type { Question } from '$lib/types';
	import { QuestionController } from '$lib/controllers';
	import { toast } from 'svelte-sonner';
	import { Edit, Trash, MessageSquare } from '@lucide/svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import Card from '$lib/components/ui/card/card.svelte';
	import CardHeader from '$lib/components/ui/card/card-header.svelte';
	import CardTitle from '$lib/components/ui/card/card-title.svelte';
	import CardContent from '$lib/components/ui/card/card-content.svelte';
	import Badge from '$lib/components/ui/badge/badge.svelte';

	interface Props {
		question: Question;
		questionNumber: number;
		evaluationId: number;
		courseId: number;
		moduleId: number;
		onupdate?: (question: Question) => void;
		ondelete?: (question: Question) => void;
	}

	let { question, questionNumber, evaluationId, courseId, moduleId, onupdate, ondelete }: Props =
		$props();

	const questionController = new QuestionController();

	async function handleDeleteQuestion() {
		if (!confirm('¿Estás seguro de que quieres eliminar esta pregunta?')) {
			return;
		}

		try {
			await questionController.deleteQuestion(question.id);
			toast.success('Pregunta eliminada exitosamente.');
			ondelete?.(question);
		} catch (error) {
			console.error('Error deleting question:', error);
			toast.error('Error al eliminar la pregunta.', {
				description: error instanceof Error ? error.message : String(error)
			});
		}
	}

	function getQuestionTypeLabel(type: string): string {
		switch (type) {
			case 'single_choice':
				return 'Opción única';
			case 'multiple_choice':
				return 'Opción múltiple';
			default:
				return type;
		}
	}

	function getQuestionTypeBadgeVariant(
		type: string
	): 'default' | 'secondary' | 'destructive' | 'outline' {
		switch (type) {
			case 'single_choice':
				return 'default';
			case 'multiple_choice':
				return 'secondary';
			default:
				return 'outline';
		}
	}
</script>

<Card class="w-full">
	<CardHeader class="flex flex-row items-start justify-between space-y-0 pb-2">
		<div class="flex items-start gap-3">
			<div
				class="bg-primary text-primary-foreground flex h-8 w-8 items-center justify-center rounded-full text-sm font-semibold"
			>
				{questionNumber}
			</div>
			<div class="flex flex-col gap-2">
				<div class="flex items-center gap-2">
					<MessageSquare class="text-muted-foreground h-4 w-4" />
					<Badge variant={getQuestionTypeBadgeVariant(question.type)}>
						{getQuestionTypeLabel(question.type)}
					</Badge>
					<Badge variant="outline">{question.points} pts</Badge>
				</div>
				<CardTitle class="text-base font-medium leading-relaxed">
					{question.text}
				</CardTitle>
			</div>
		</div>
		<div class="flex items-center gap-2">
			<Button
				variant="outline"
				size="sm"
				href="/admin/courses/{courseId}/{moduleId}/evaluations/{evaluationId}/questions/{question.id}"
			>
				<Edit class="h-4 w-4" />
				<span>Editar</span>
			</Button>
			<Button variant="destructive" size="sm" onclick={handleDeleteQuestion}>
				<Trash class="h-4 w-4" />
			</Button>
		</div>
	</CardHeader>
	<CardContent>
		{#if question.explanation}
			<div class="text-muted-foreground text-sm">
				<strong>Explicación:</strong>
				{question.explanation}
			</div>
		{/if}
		{#if question.answers && question.answers.length > 0}
			<div class="mt-3">
				<p class="text-muted-foreground mb-2 text-sm font-medium">
					Respuestas ({question.answers.length}):
				</p>
				<div class="space-y-1">
					{#each question.answers.slice(0, 3) as answer}
						<div class="flex items-center gap-2 text-sm">
							<div
								class="h-2 w-2 rounded-full {answer.is_correct ? 'bg-green-500' : 'bg-gray-300'}"
							></div>
							<span
								class={answer.is_correct ? 'font-medium text-green-700' : 'text-muted-foreground'}
							>
								{answer.text}
							</span>
						</div>
					{/each}
					{#if question.answers.length > 3}
						<p class="text-muted-foreground text-xs">
							... y {question.answers.length - 3} respuesta(s) más
						</p>
					{/if}
				</div>
			</div>
		{:else}
			<div class="mt-3 rounded-md border border-yellow-200 bg-yellow-50 p-3">
				<p class="text-sm text-yellow-800">Esta pregunta no tiene respuestas configuradas aún.</p>
			</div>
		{/if}
	</CardContent>
</Card>

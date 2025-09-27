<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import Back from '$lib/components/kit/Back.svelte';
	import type { PageProps } from './$types';
	import { QuestionType } from '$lib/types';
	import { Trophy, X, Check, ArrowLeft, Eye, Clock, Target, FileText } from '@lucide/svelte';

	let { data }: PageProps = $props();

	const passed = $derived(data.attempt.passed);
	const percentage = $derived(Math.round((data.attempt.score / data.attempt.total_points) * 100));

	function formatDuration(minutes: number): string {
		if (!minutes) return 'N/A';
		if (minutes < 60) return `${minutes} min`;
		const hours = Math.floor(minutes / 60);
		const mins = minutes % 60;
		return `${hours}h ${mins > 0 ? mins + 'm' : ''}`;
	}

	function getQuestionResult(questionId: number) {
		const answer = data.attempt.answers.find((a) => a.attempt_question_id === questionId);
		return answer || null;
	}

	function getSelectedOptions(questionId: number) {
		const result = getQuestionResult(questionId);
		if (!result) return [];
		const question = data.attempt.questions.find((q) => q.id === questionId);
		if (!question) return [];

		return result.selected_option_ids
			.map((id) => question.answer_options.find((opt) => opt.id === id))
			.filter(Boolean);
	}
</script>

<div class="mx-auto max-w-4xl p-6">
	<!-- Header -->
	<div class="mb-6">
		<Back href="/courses/{data.courseId}/{data.moduleId}" class="mb-4" />
		<div class="text-center">
			<h1 class="mb-2 text-3xl font-bold">{data.attempt.evaluation?.title}</h1>
			<p class="text-muted-foreground">Resultados del examen</p>
		</div>
	</div>

	<!-- Results Summary -->
	<div class="mb-8 grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-4">
		<Card.Root>
			<Card.Header class="pb-2">
				<div class="flex items-center gap-2">
					<Trophy class="h-5 w-5 text-yellow-600" />
					<Card.Title class="text-sm font-medium">Puntuación</Card.Title>
				</div>
			</Card.Header>
			<Card.Content>
				<div class="text-2xl font-bold {passed ? 'text-green-600' : 'text-red-600'}">
					{data.attempt.score}/{data.attempt.total_points}
				</div>
				<p class="text-muted-foreground text-xs">
					{percentage}%
				</p>
			</Card.Content>
		</Card.Root>

		<Card.Root>
			<Card.Header class="pb-2">
				<div class="flex items-center gap-2">
					<Target class="h-5 w-5 text-blue-600" />
					<Card.Title class="text-sm font-medium">Estado</Card.Title>
				</div>
			</Card.Header>
			<Card.Content>
				<Badge class={passed ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}>
					{passed ? 'Aprobado' : 'No Aprobado'}
				</Badge>
				<p class="text-muted-foreground mt-1 text-xs">
					Mínimo: {data.attempt.evaluation?.passing_score}%
				</p>
			</Card.Content>
		</Card.Root>

		<Card.Root>
			<Card.Header class="pb-2">
				<div class="flex items-center gap-2">
					<FileText class="h-5 w-5 text-purple-600" />
					<Card.Title class="text-sm font-medium">Preguntas</Card.Title>
				</div>
			</Card.Header>
			<Card.Content>
				<div class="text-2xl font-bold">
					{data.attempt.questions.length}
				</div>
				<p class="text-muted-foreground text-xs">
					{data.attempt.answers.filter((a) => a.is_correct).length} correctas
				</p>
			</Card.Content>
		</Card.Root>

		<Card.Root>
			<Card.Header class="pb-2">
				<div class="flex items-center gap-2">
					<Clock class="h-5 w-5 text-orange-600" />
					<Card.Title class="text-sm font-medium">Tiempo</Card.Title>
				</div>
			</Card.Header>
			<Card.Content>
				<div class="text-2xl font-bold">
					{formatDuration(data.attempt.time_spent)}
				</div>
				<p class="text-muted-foreground text-xs">
					{data.attempt.evaluation?.time_limit
						? `Límite: ${formatDuration(data.attempt.evaluation.time_limit)}`
						: 'Sin límite'}
				</p>
			</Card.Content>
		</Card.Root>
	</div>

	<!-- Review Questions -->
	<div class="space-y-6">
		<div class="flex items-center justify-between">
			<h2 class="text-xl font-semibold">Revisión detallada</h2>
			<Button
				variant="outline"
				href="/courses/{data.courseId}/{data.moduleId}/evaluation/{data.evaluationId}/attempts"
			>
				<Eye class="mr-2 h-4 w-4" />
				Ver todos los intentos
			</Button>
		</div>

		{#each data.attempt.questions as question, index}
			{@const result = getQuestionResult(question.id)}
			{@const selectedOptions = getSelectedOptions(question.id)}
			{@const isCorrect = result?.is_correct || false}

			<Card.Root class="overflow-hidden">
				<Card.Header>
					<div class="flex items-start justify-between">
						<div class="flex-1">
							<div class="mb-2 flex items-center gap-2">
								<span class="text-muted-foreground text-sm font-medium">Pregunta {index + 1}</span>
								<Badge variant={isCorrect ? 'default' : 'destructive'} class="text-xs">
									{#if isCorrect}
										<Check class="mr-1 h-3 w-3" />
										Correcta
									{:else}
										<X class="mr-1 h-3 w-3" />
										Incorrecta
									{/if}
								</Badge>
								<Badge variant="secondary" class="text-xs">
									{result?.points || 0}/{question.points} puntos
								</Badge>
							</div>
							<Card.Title class="text-base">{question.text}</Card.Title>
							{#if question.explanation}
								<Card.Description class="mt-2">{question.explanation}</Card.Description>
							{/if}
						</div>
						<Badge
							variant={question.type === QuestionType.SINGLE ? 'outline' : 'secondary'}
							class="text-xs"
						>
							{question.type === QuestionType.SINGLE ? 'Una opción' : 'Múltiple opción'}
						</Badge>
					</div>
				</Card.Header>

				<Card.Content>
					<div class="space-y-3">
						{#each question.answer_options as option}
							{@const isSelected = selectedOptions.some((opt) => opt?.id === option.id)}
							{@const isCorrectOption = option.is_correct}

							<div
								class="flex items-center space-x-3 rounded-lg border-2 p-3 transition-colors
								{isSelected && isCorrectOption
									? 'border-green-500 bg-green-50'
									: isSelected && !isCorrectOption
										? 'border-red-500 bg-red-50'
										: !isSelected && isCorrectOption
											? 'bg-green-25 border-green-300'
											: 'border-gray-200'}
							"
							>
								<div class="flex items-center">
									{#if isSelected}
										{#if isCorrectOption}
											<Check class="h-4 w-4 text-green-600" />
										{:else}
											<X class="h-4 w-4 text-red-600" />
										{/if}
									{:else if isCorrectOption}
										<div class="h-4 w-4 rounded-full border-2 border-green-500 bg-green-100"></div>
									{:else}
										<div class="h-4 w-4 rounded-full border-2 border-gray-300"></div>
									{/if}
								</div>

								<span
									class="flex-1 {isSelected && isCorrectOption
										? 'font-medium text-green-800'
										: isSelected && !isCorrectOption
											? 'text-red-800'
											: !isSelected && isCorrectOption
												? 'font-medium text-green-700'
												: 'text-gray-700'}"
								>
									{option.text}
								</span>

								{#if !isSelected && isCorrectOption}
									<Badge variant="outline" class="border-green-300 text-xs text-green-700">
										Respuesta correcta
									</Badge>
								{:else if isSelected && !isCorrectOption}
									<Badge variant="outline" class="border-red-300 text-xs text-red-700">
										Tu respuesta
									</Badge>
								{:else if isSelected && isCorrectOption}
									<Badge class="bg-green-600 text-xs text-white">Correcta ✓</Badge>
								{/if}
							</div>
						{/each}

						{#if result && result.selected_option_ids.length === 0}
							<div class="rounded-lg border border-yellow-200 bg-yellow-50 p-3">
								<p class="text-sm text-yellow-800">No respondiste esta pregunta</p>
							</div>
						{/if}
					</div>
				</Card.Content>
			</Card.Root>
		{/each}
	</div>

	<!-- Action buttons -->
	<div class="mt-8 flex justify-center gap-4">
		<Button variant="outline" href="/courses/{data.courseId}/{data.moduleId}">
			<ArrowLeft class="mr-2 h-4 w-4" />
			Volver al módulo
		</Button>

		<Button href="/courses/{data.courseId}/{data.moduleId}/evaluation/{data.evaluationId}/attempts">
			<Eye class="mr-2 h-4 w-4" />
			Ver todos los intentos
		</Button>
	</div>
</div>

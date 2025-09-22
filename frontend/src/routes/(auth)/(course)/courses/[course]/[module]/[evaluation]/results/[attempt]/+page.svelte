<script lang="ts">
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$components/ui/card';
	import { Badge } from '$components/ui/badge';
	import { Button } from '$components/ui/button';
	// Remove Progress import since it doesn't exist
	import { 
		CheckCircle, 
		XCircle, 
		Clock, 
		Target, 
		Users, 
		Trophy,
		AlertCircle,
		ArrowRight,
		RotateCcw,
		Home
	} from '@lucide/svelte';
	import Back from '$lib/components/kit/Back.svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	// Computed properties
	const evaluation = $derived(data.attempt.evaluation);
	const totalQuestions = $derived(data.attempt.questions.length);
	const totalPoints = $derived(data.attempt.total_points);
	const score = $derived(data.attempt.score);
	const passed = $derived(data.attempt.passed);
	const passingScore = $derived(evaluation?.passing_score || 0);

	// Calculate detailed results
	const correctAnswers = $derived(data.attempt.answers.filter(answer => answer.is_correct).length);
	const wrongAnswers = $derived(data.attempt.answers.filter(answer => !answer.is_correct).length);
	const unansweredQuestions = $derived(totalQuestions - data.attempt.answers.length);

	// Time calculations
	const timeSpent = $derived(() => {
		if (!data.attempt.started_at || !data.attempt.submitted_at) return null;
		const start = new Date(data.attempt.started_at).getTime();
		const end = new Date(data.attempt.submitted_at).getTime();
		return Math.round((end - start) / 60000); // minutes
	});

	function formatDuration(minutes: number): string {
		if (minutes < 60) {
			return `${minutes} minutos`;
		}
		const hours = Math.floor(minutes / 60);
		const remainingMinutes = minutes % 60;
		return `${hours}h ${remainingMinutes}m`;
	}

	function getScoreColor(): string {
		if (passed) return 'text-green-600';
		if (score >= passingScore * 0.8) return 'text-orange-600';
		return 'text-red-600';
	}

	function getPerformanceLevel(): { label: string; color: string } {
		if (score >= 90) return { label: 'Excelente', color: 'bg-green-100 text-green-800' };
		if (score >= 80) return { label: 'Muy Bueno', color: 'bg-blue-100 text-blue-800' };
		if (score >= passingScore) return { label: 'Bueno', color: 'bg-yellow-100 text-yellow-800' };
		return { label: 'Necesita Mejorar', color: 'bg-red-100 text-red-800' };
	}

	// Get answer for specific question
	function getAnswerForQuestion(questionId: number) {
		return data.attempt.answers.find(answer => answer.attempt_question_id === questionId);
	}

	// Get selected options text for a question
	function getSelectedOptionsText(questionId: number): string {
		const answer = getAnswerForQuestion(questionId);
		if (!answer) return 'Sin respuesta';
		
		const question = data.attempt.questions.find(q => q.id === questionId);
		if (!question) return 'Sin respuesta';

		const selectedOptions = question.answer_options.filter(option => 
			answer.selected_option_ids.includes(option.id)
		);

		return selectedOptions.map(option => option.text).join(', ') || 'Sin respuesta';
	}
</script>

<Back href="/courses/{data.courseId}/{data.moduleId}/{data.evaluationId}/start" />

<div class="flex flex-col gap-8">
	<!-- Header with result status -->
	<div class="text-center">
		{#if passed}
			<div class="mx-auto w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mb-4">
				<Trophy class="h-8 w-8 text-green-600" />
			</div>
			<h1 class="text-h1 text-green-600">¡Felicitaciones!</h1>
			<p class="text-subtitle">Has aprobado la evaluación</p>
		{:else}
			<div class="mx-auto w-16 h-16 bg-red-100 rounded-full flex items-center justify-center mb-4">
				<XCircle class="h-8 w-8 text-red-600" />
			</div>
			<h1 class="text-h1 text-red-600">Evaluación No Aprobada</h1>
			<p class="text-subtitle">Necesitas {passingScore}% para aprobar</p>
		{/if}
	</div>

	<!-- Score Summary -->
	<div class="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
		<Card>
			<CardContent class="p-6">
				<div class="flex items-center justify-between">
					<div>
						<p class="text-sm text-muted-foreground">Puntuación</p>
						<p class="text-3xl font-bold {getScoreColor()}">{score}%</p>
					</div>
					<Target class="h-8 w-8 {getScoreColor()}" />
				</div>
				<div class="w-full bg-gray-200 rounded-full h-2 mt-3">
					<div class="bg-primary h-2 rounded-full transition-all duration-300" style="width: {score}%"></div>
				</div>
			</CardContent>
		</Card>

		<Card>
			<CardContent class="p-6">
				<div class="flex items-center justify-between">
					<div>
						<p class="text-sm text-muted-foreground">Respuestas Correctas</p>
						<p class="text-3xl font-bold text-green-600">{correctAnswers}</p>
						<p class="text-xs text-muted-foreground">de {totalQuestions}</p>
					</div>
					<CheckCircle class="h-8 w-8 text-green-600" />
				</div>
			</CardContent>
		</Card>

		<Card>
			<CardContent class="p-6">
				<div class="flex items-center justify-between">
					<div>
						<p class="text-sm text-muted-foreground">Puntos Obtenidos</p>
						<p class="text-3xl font-bold text-blue-600">{Math.round(score * totalPoints / 100)}</p>
						<p class="text-xs text-muted-foreground">de {totalPoints}</p>
					</div>
					<Users class="h-8 w-8 text-blue-600" />
				</div>
			</CardContent>
		</Card>

		{#if timeSpent()}
			<Card>
				<CardContent class="p-6">
					<div class="flex items-center justify-between">
						<div>
							<p class="text-sm text-muted-foreground">Tiempo Utilizado</p>
							<p class="text-3xl font-bold text-purple-600">{timeSpent()}</p>
							<p class="text-xs text-muted-foreground">minutos</p>
						</div>
						<Clock class="h-8 w-8 text-purple-600" />
					</div>
				</CardContent>
			</Card>
		{/if}
	</div>

	<!-- Performance Analysis -->
	<Card>
		<CardHeader>
			<CardTitle>Análisis de Rendimiento</CardTitle>
			<CardDescription>Resumen detallado de tu evaluación</CardDescription>
		</CardHeader>
		<CardContent>
			<div class="grid gap-4 md:grid-cols-2">
				<div>
					<div class="flex items-center justify-between mb-2">
						<span class="text-sm text-muted-foreground">Nivel de Rendimiento</span>
						<Badge class={getPerformanceLevel().color}>
							{getPerformanceLevel().label}
						</Badge>
					</div>

					<div class="space-y-2 text-sm">
						<div class="flex justify-between">
							<span>Respuestas Correctas:</span>
							<span class="font-medium text-green-600">{correctAnswers}</span>
						</div>
						<div class="flex justify-between">
							<span>Respuestas Incorrectas:</span>
							<span class="font-medium text-red-600">{wrongAnswers}</span>
						</div>
						{#if unansweredQuestions > 0}
							<div class="flex justify-between">
								<span>Preguntas Sin Responder:</span>
								<span class="font-medium text-orange-600">{unansweredQuestions}</span>
							</div>
						{/if}
						{#if evaluation?.time_limit && timeSpent()}
							<div class="flex justify-between">
								<span>Tiempo Límite:</span>
								<span class="font-medium">{formatDuration(evaluation.time_limit)}</span>
							</div>
							<div class="flex justify-between">
								<span>Tiempo Utilizado:</span>
								<span class="font-medium">{formatDuration(timeSpent())}</span>
							</div>
						{/if}
					</div>
				</div>

				<div>
					<h4 class="text-sm font-medium mb-3">Detalles de la Evaluación</h4>
					<div class="space-y-2 text-sm">
						<div class="flex justify-between">
							<span>Fecha de Realización:</span>
							<span class="font-medium">
								{new Date(data.attempt.submitted_at || '').toLocaleDateString('es-ES')}
							</span>
						</div>
						<div class="flex justify-between">
							<span>Hora de Finalización:</span>
							<span class="font-medium">
								{new Date(data.attempt.submitted_at || '').toLocaleTimeString('es-ES')}
							</span>
						</div>
						<div class="flex justify-between">
							<span>Estado:</span>
							<Badge variant={passed ? "default" : "destructive"}>
								{passed ? "Aprobado" : "No Aprobado"}
							</Badge>
						</div>
					</div>
				</div>
			</div>
		</CardContent>
	</Card>

	<!-- Question-by-Question Review -->
	<Card>
		<CardHeader>
			<CardTitle>Revisión por Pregunta</CardTitle>
			<CardDescription>
				Revisa tus respuestas y las explicaciones de cada pregunta
			</CardDescription>
		</CardHeader>
		<CardContent>
			<div class="space-y-6">
				{#each data.attempt.questions as question, index}
					{@const answer = getAnswerForQuestion(question.id)}
					{@const isCorrect = answer?.is_correct || false}
					{@const selectedOptions = question.answer_options.filter(option => 
						answer?.selected_option_ids.includes(option.id)
					)}
					{@const correctOptions = question.answer_options.filter(option => option.is_correct)}
					
					<div class="border rounded-lg p-4">
						<div class="flex items-start justify-between mb-3">
							<div class="flex-1">
								<div class="flex items-center gap-2 mb-2">
									<span class="text-sm font-medium">Pregunta {index + 1}</span>
									<Badge variant="outline">{question.points} puntos</Badge>
									{#if isCorrect}
										<CheckCircle class="h-4 w-4 text-green-600" />
									{:else}
										<XCircle class="h-4 w-4 text-red-600" />
									{/if}
								</div>
								<p class="text-base">{question.text}</p>
							</div>
						</div>

						<div class="grid gap-4 md:grid-cols-2">
							<!-- Your Answer -->
							<div>
								<h5 class="text-sm font-medium mb-2 text-muted-foreground">Tu Respuesta:</h5>
								{#if selectedOptions.length > 0}
									<ul class="space-y-1">
										{#each selectedOptions as option}
											<li class="flex items-center gap-2">
												{#if isCorrect}
													<CheckCircle class="h-3 w-3 text-green-600 flex-shrink-0" />
												{:else}
													<XCircle class="h-3 w-3 text-red-600 flex-shrink-0" />
												{/if}
												<span class="text-sm">{option.text}</span>
											</li>
										{/each}
									</ul>
								{:else}
									<p class="text-sm text-muted-foreground italic">Sin respuesta</p>
								{/if}
							</div>

							<!-- Correct Answer -->
							<div>
								<h5 class="text-sm font-medium mb-2 text-muted-foreground">Respuesta Correcta:</h5>
								<ul class="space-y-1">
									{#each correctOptions as option}
										<li class="flex items-center gap-2">
											<CheckCircle class="h-3 w-3 text-green-600 flex-shrink-0" />
											<span class="text-sm">{option.text}</span>
										</li>
									{/each}
								</ul>
							</div>
						</div>

						{#if question.explanation}
							<div class="mt-4 p-3 bg-blue-50 border border-blue-200 rounded-md">
								<h5 class="text-sm font-medium text-blue-800 mb-1">Explicación:</h5>
								<p class="text-sm text-blue-700">{question.explanation}</p>
							</div>
						{/if}
					</div>
				{/each}
			</div>
		</CardContent>
	</Card>

	<!-- Actions -->
	<div class="flex flex-col gap-4 sm:flex-row sm:justify-center">
		<Button 
			href="/courses/{data.courseId}/{data.moduleId}" 
			variant="outline"
		>
			<Home class="mr-2 h-4 w-4" />
			Volver al Módulo
		</Button>

		{#if !passed && evaluation?.max_attempts && data.attempt.user?.attemptCount < evaluation.max_attempts}
			<Button 
				href="/courses/{data.courseId}/{data.moduleId}/{data.evaluationId}/start"
				class="bg-blue-600 hover:bg-blue-700"
			>
				<RotateCcw class="mr-2 h-4 w-4" />
				Intentar Nuevamente
			</Button>
		{/if}

		<Button 
			href="/courses/{data.courseId}" 
			variant="outline"
		>
			<ArrowRight class="mr-2 h-4 w-4" />
			Continuar Curso
		</Button>
	</div>
</div>
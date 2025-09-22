<script lang="ts">
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$components/ui/card';
	import { Badge } from '$components/ui/badge';
	import { Button } from '$components/ui/button';
	import { 
		Clock, 
		Target, 
		Users, 
		AlertTriangle, 
		CheckCircle, 
		PlayCircle,
		RotateCcw,
		Info
	} from '@lucide/svelte';
	import Back from '$lib/components/kit/Back.svelte';
	import { EvaluationAttemptController } from '$lib';
	import { toast } from 'svelte-sonner';
	import { goto } from '$app/navigation';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	const evaluationAttemptController = new EvaluationAttemptController();

	let isStarting = $state(false);

	const completedAttempts = $derived(data.userAttempts.filter(attempt => attempt.submitted_at));
	const hasPassedAttempt = $derived(completedAttempts.some(attempt => attempt.passed));
	const remainingAttempts = $derived(
		data.evaluation.max_attempts ? data.evaluation.max_attempts - completedAttempts.length : null
	);
	const bestScore = $derived(() => {
		if (completedAttempts.length === 0) return null;
		return completedAttempts.reduce((best, current) => 
			current.score > best.score ? current : best
		);
	});

	async function startEvaluation() {
		if (!data.canAttempt.can_attempt) {
			toast.error(data.canAttempt.reason);
			return;
		}

		isStarting = true;
		
		try {
			const attempt = await evaluationAttemptController.startAttempt({
				evaluation_id: data.evaluation.id,
				user_id: data.userId
			});

			// Navigate to the exam interface
			goto(`/courses/${data.courseId}/${data.moduleId}/${data.evaluation.id}/attempt/${attempt.id}`);
		} catch (error) {
			console.error('Error starting evaluation:', error);
			toast.error('Error al iniciar la evaluación. Intenta nuevamente.');
		} finally {
			isStarting = false;
		}
	}

	function formatDuration(minutes: number): string {
		if (minutes < 60) {
			return `${minutes} minutos`;
		}
		const hours = Math.floor(minutes / 60);
		const remainingMinutes = minutes % 60;
		return `${hours}h ${remainingMinutes}m`;
	}
</script>

<Back href="/courses/{data.courseId}/{data.moduleId}" />

<div class="flex flex-col gap-6">
	<div>
		<h1 class="text-h1">{data.evaluation.title}</h1>
		{#if data.evaluation.description}
			<p class="text-subtitle">{data.evaluation.description}</p>
		{/if}
	</div>

	<div class="grid gap-6 md:grid-cols-3">
		<!-- Evaluation Info -->
		<div class="md:col-span-2">
			<Card>
				<CardHeader>
					<CardTitle class="flex items-center gap-2">
						<Info class="h-5 w-5" />
						Información de la Evaluación
					</CardTitle>
				</CardHeader>
				<CardContent class="space-y-4">
					<div class="grid grid-cols-2 gap-4">
						<div class="flex items-center gap-2">
							<Users class="h-4 w-4 text-muted-foreground" />
							<span class="text-sm">
								<strong>{data.evaluation.question_count}</strong> preguntas
							</span>
						</div>
						
						<div class="flex items-center gap-2">
							<Target class="h-4 w-4 text-muted-foreground" />
							<span class="text-sm">
								<strong>{data.evaluation.passing_score}%</strong> para aprobar
							</span>
						</div>

						{#if data.evaluation.time_limit}
							<div class="flex items-center gap-2">
								<Clock class="h-4 w-4 text-muted-foreground" />
								<span class="text-sm">
									<strong>{formatDuration(data.evaluation.time_limit)}</strong> límite de tiempo
								</span>
							</div>
						{/if}

						{#if data.evaluation.max_attempts}
							<div class="flex items-center gap-2">
								<RotateCcw class="h-4 w-4 text-muted-foreground" />
								<span class="text-sm">
									<strong>{data.evaluation.max_attempts}</strong> intentos máximo
								</span>
							</div>
						{/if}
					</div>

					{#if data.evaluation.time_limit}
						<div class="rounded-md border border-yellow-200 bg-yellow-50 p-3">
							<div class="flex items-center gap-2 text-yellow-800">
								<AlertTriangle class="h-4 w-4" />
								<span class="text-sm font-medium">Tiempo Límite</span>
							</div>
							<p class="text-yellow-700 text-sm mt-1">
								Tienes {formatDuration(data.evaluation.time_limit)} para completar esta evaluación. 
								El tiempo empezará a correr cuando inicies la evaluación.
							</p>
						</div>
					{/if}

					{#if !data.evaluation.time_limit}
						<div class="rounded-md border border-green-200 bg-green-50 p-3">
							<div class="flex items-center gap-2 text-green-800">
								<CheckCircle class="h-4 w-4" />
								<span class="text-sm font-medium">Sin Límite de Tiempo</span>
							</div>
							<p class="text-green-700 text-sm mt-1">
								Puedes tomar el tiempo que necesites para completar esta evaluación.
							</p>
						</div>
					{/if}
				</CardContent>
			</Card>
		</div>

		<!-- Status & Actions -->
		<div class="space-y-6">
			<!-- Current Status -->
			<Card>
				<CardHeader>
					<CardTitle class="text-base">Estado Actual</CardTitle>
				</CardHeader>
				<CardContent class="space-y-3">
					{#if hasPassedAttempt}
						<div class="flex items-center gap-2 text-green-700">
							<CheckCircle class="h-4 w-4" />
							<span class="text-sm font-medium">Evaluación Aprobada</span>
						</div>
					{:else if completedAttempts.length > 0}
						<div class="flex items-center gap-2 text-orange-700">
							<AlertTriangle class="h-4 w-4" />
							<span class="text-sm font-medium">Pendiente de Aprobar</span>
						</div>
					{:else}
						<div class="flex items-center gap-2 text-blue-700">
							<PlayCircle class="h-4 w-4" />
							<span class="text-sm font-medium">Listo para Iniciar</span>
						</div>
					{/if}

					{#if bestScore()}
						<div class="text-sm">
							<span class="text-muted-foreground">Mejor puntuación:</span>
							<span class="font-medium">{bestScore()?.score}%</span>
						</div>
					{/if}

					{#if remainingAttempts !== null}
						<div class="text-sm">
							<span class="text-muted-foreground">Intentos restantes:</span>
							<span class="font-medium">{remainingAttempts}</span>
						</div>
					{/if}

					<div class="text-sm">
						<span class="text-muted-foreground">Intentos realizados:</span>
						<span class="font-medium">{completedAttempts.length}</span>
					</div>
				</CardContent>
			</Card>

			<!-- Action Button -->
			<Card>
				<CardContent class="pt-6">
					{#if !data.canAttempt.can_attempt}
						<div class="text-center space-y-2">
							<p class="text-sm text-muted-foreground">{data.canAttempt.reason}</p>
							{#if hasPassedAttempt}
								<Badge variant="default" class="bg-green-100 text-green-800">
									Completado
								</Badge>
							{:else}
								<Badge variant="destructive">
									No Disponible
								</Badge>
							{/if}
						</div>
					{:else}
						<Button 
							onclick={startEvaluation} 
							disabled={isStarting}
							class="w-full"
							size="lg"
						>
							{#if isStarting}
								Iniciando...
							{:else}
								<PlayCircle class="mr-2 h-4 w-4" />
								Iniciar Evaluación
							{/if}
						</Button>
					{/if}
				</CardContent>
			</Card>
		</div>
	</div>

	<!-- Previous Attempts -->
	{#if completedAttempts.length > 0}
		<Card>
			<CardHeader>
				<CardTitle>Intentos Anteriores</CardTitle>
				<CardDescription>Historial de tus intentos en esta evaluación</CardDescription>
			</CardHeader>
			<CardContent>
				<div class="space-y-3">
					{#each completedAttempts.sort((a, b) => new Date(b.submitted_at || 0).getTime() - new Date(a.submitted_at || 0).getTime()) as attempt}
						<div class="flex items-center justify-between rounded-lg border p-3">
							<div class="flex items-center gap-3">
								<div class="flex-shrink-0">
									{#if attempt.passed}
										<CheckCircle class="h-5 w-5 text-green-600" />
									{:else}
										<AlertTriangle class="h-5 w-5 text-orange-600" />
									{/if}
								</div>
								<div>
									<p class="text-sm font-medium">
										{attempt.score}% ({attempt.score * attempt.total_points / 100} / {attempt.total_points} puntos)
									</p>
									<p class="text-xs text-muted-foreground">
										{new Date(attempt.submitted_at || '').toLocaleString('es-ES')}
									</p>
								</div>
							</div>
							<div class="flex items-center gap-2">
								<Badge variant={attempt.passed ? "default" : "secondary"} 
								       class={attempt.passed ? "bg-green-100 text-green-800" : ""}>
									{attempt.passed ? "Aprobado" : "No Aprobado"}
								</Badge>
							</div>
						</div>
					{/each}
				</div>
			</CardContent>
		</Card>
	{/if}
</div>
<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import type { Evaluation } from '$lib/types';
	import { Clock, FileText, Target, Trophy, User } from '@lucide/svelte';
	import { authStore } from '$lib/stores/auth.svelte';
	import { EvaluationAttemptController } from '$lib/controllers/evaluationAttempt';
	import { onMount } from 'svelte';
	import { toast } from 'svelte-sonner';

	type Props = {
		evaluation: Evaluation;
		courseId: number;
	};

	const { evaluation, courseId }: Props = $props();
	const evaluationAttemptController = new EvaluationAttemptController();

	let canAttempt = $state(true);
	let attemptStatus = $state<{ can_attempt: boolean; reason: string } | null>(null);
	let userAttempts = $state<any[]>([]);
	let loading = $state(false);

	onMount(async () => {
		if (!authStore.user) return;

		try {
			const [statusResponse, attemptsResponse] = await Promise.all([
				evaluationAttemptController.canUserAttempt(authStore.user.id, evaluation.id),
				evaluationAttemptController.getUserAttempts(authStore.user.id, evaluation.id)
			]);

			attemptStatus = statusResponse;
			userAttempts = attemptsResponse;
			canAttempt = statusResponse.can_attempt;
		} catch (error) {
			console.error('Error loading evaluation status:', error);
		}
	});

	async function startAttempt() {
		if (!authStore.user || !canAttempt) return;

		loading = true;
		try {
			const attempt = await evaluationAttemptController.startAttempt({
				user_id: authStore.user.id,
				evaluation_id: evaluation.id
			});

			// Navigate to the exam taking page
			window.location.href = `/courses/${courseId}/module/${evaluation.module_id}/evaluation/${evaluation.id}/attempt/${attempt.id}`;
		} catch (error) {
			console.error('Error starting attempt:', error);
			toast.error('No se pudo iniciar el examen. Por favor intenta de nuevo.');
		} finally {
			loading = false;
		}
	}

	function viewAttempts() {
		window.location.href = `/courses/${courseId}/module/${evaluation.module_id}/evaluation/${evaluation.id}/attempts`;
	}

	function formatDuration(minutes: number): string {
		if (minutes < 60) return `${minutes} min`;
		const hours = Math.floor(minutes / 60);
		const mins = minutes % 60;
		return `${hours}h ${mins > 0 ? mins + 'm' : ''}`;
	}

	const bestAttempt = $derived(
		userAttempts.length > 0
			? userAttempts.reduce((best, current) => (current.score > best.score ? current : best))
			: null
	);

	const passed = $derived(bestAttempt && bestAttempt.passed);
</script>

<Card.Root class="hover:shadow-md transition-shadow">
	<Card.Header>
		<div class="flex items-start justify-between">
			<div class="flex-1">
				<Card.Title class="flex items-center gap-2 text-lg">
					<FileText class="h-5 w-5" />
					{evaluation.title}
				</Card.Title>
				{#if evaluation.description}
					<Card.Description class="mt-2">{evaluation.description}</Card.Description>
				{/if}
			</div>
			{#if passed}
				<Badge variant="secondary" class="bg-green-100 text-green-800">
					<Trophy class="h-3 w-3 mr-1" />
					Aprobado
				</Badge>
			{:else if userAttempts.length > 0}
				<Badge variant="secondary" class="bg-orange-100 text-orange-800">
					Intentos: {userAttempts.length}
				</Badge>
			{/if}
		</div>
	</Card.Header>

	<Card.Content>
		<div class="grid grid-cols-2 gap-4 text-sm text-muted-foreground">
			<div class="flex items-center gap-2">
				<FileText class="h-4 w-4" />
				{evaluation.question_count} preguntas
			</div>
			<div class="flex items-center gap-2">
				<Target class="h-4 w-4" />
				{evaluation.passing_score}% para aprobar
			</div>
			{#if evaluation.time_limit}
				<div class="flex items-center gap-2">
					<Clock class="h-4 w-4" />
					{formatDuration(evaluation.time_limit)}
				</div>
			{/if}
			{#if evaluation.max_attempts}
				<div class="flex items-center gap-2">
					<User class="h-4 w-4" />
					Máx. {evaluation.max_attempts} intentos
				</div>
			{/if}
		</div>

		{#if bestAttempt}
			<div class="mt-4 p-3 bg-muted/50 rounded-lg">
				<div class="text-sm font-medium">Mejor intento:</div>
				<div class="flex items-center justify-between mt-1">
					<span class="text-2xl font-bold {passed ? 'text-green-600' : 'text-orange-600'}">
						{bestAttempt.score}/{bestAttempt.total_points}
					</span>
					<span class="text-sm text-muted-foreground">
						{Math.round((bestAttempt.score / bestAttempt.total_points) * 100)}%
					</span>
				</div>
			</div>
		{/if}

		{#if !canAttempt && attemptStatus}
			<div class="mt-4 p-3 bg-yellow-50 border border-yellow-200 rounded-lg">
				<p class="text-sm text-yellow-800">{attemptStatus.reason}</p>
			</div>
		{/if}
	</Card.Content>

	<Card.Footer class="flex gap-2">
		{#if canAttempt}
			<Button onclick={startAttempt} disabled={loading} class="flex-1">
				{loading ? 'Iniciando...' : 'Tomar Examen'}
			</Button>
		{/if}
		
		{#if userAttempts.length > 0}
			<Button variant="outline" onclick={viewAttempts}>
				Ver Intentos
			</Button>
		{/if}
	</Card.Footer>
</Card.Root>
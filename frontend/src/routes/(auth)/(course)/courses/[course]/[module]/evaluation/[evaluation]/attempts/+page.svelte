<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import Back from '$lib/components/kit/Back.svelte';
	import type { PageProps } from './$types';
	import { Trophy, X, Clock, Calendar, Eye, Plus, Target, FileText } from '@lucide/svelte';
	import { EvaluationAttemptController } from '$lib/controllers/evaluationAttempt';
	import { authStore } from '$lib/stores/auth.svelte';
	import { toast } from 'svelte-sonner';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import type { EvaluationAttempt } from '$lib';
	import { page } from '$app/state';

	let { data }: PageProps = $props();

	const evaluationUrl = $derived(`/courses/${page.params.course}/${page.params.module}/evaluation/${page.params.evaluation}`)

	const evaluationAttemptController = new EvaluationAttemptController();
	let loading = $state(false);
	let attempts = $state<EvaluationAttempt[]>(data.attempts || []);
	let userId = $state<number>(data.userId);

	// Sort attempts by date (newest first)
	const sortedAttempts = $derived(
		attempts.toSorted((a, b) => new Date(b.started_at).getTime() - new Date(a.started_at).getTime())
	);

	const bestAttempt = $derived(
		attempts.length > 0
			? attempts.reduce((best, current) => (current.score > best.score ? current : best))
			: null
	);

	const passed = $derived(bestAttempt && bestAttempt.passed);
	const canAttempt = $derived(
		!data.evaluation.max_attempts || attempts.length < data.evaluation.max_attempts
	);

	async function startNewAttempt() {
		if (loading || !canAttempt || !userId) return;

		loading = true;
		try {
			const attempt = await evaluationAttemptController.startAttempt({
				user_id: userId,
				evaluation_id: data.evaluationId
			});

			toast.success('Nuevo intento iniciado');
			await goto(
				`${evaluationUrl}/attempt/${attempt.id}`
			);
		} catch (error) {
			console.error('Error starting attempt:', error);
			toast.error('No se pudo iniciar el examen. Por favor intenta de nuevo.');
		} finally {
			loading = false;
		}
	}

	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleString('es-ES', {
			year: 'numeric',
			month: 'long',
			day: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}

	function formatDuration(minutes: number): string {
		if (!minutes) return 'N/A';
		if (minutes < 60) return `${minutes} min`;
		const hours = Math.floor(minutes / 60);
		const mins = minutes % 60;
		return `${hours}h ${mins > 0 ? mins + 'm' : ''}`;
	}

	function getStatusColor(attempt: any): string {
		if (!attempt.submitted_at) return 'bg-yellow-100 text-yellow-800';
		return attempt.passed ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800';
	}

	function getStatusText(attempt: any): string {
		if (!attempt.submitted_at) return 'En progreso';
		return attempt.passed ? 'Aprobado' : 'No aprobado';
	}

</script>

<div class="mx-auto max-w-4xl p-6">
	<!-- Header -->
	<div class="mb-6 space-y-4">
		<Back href="/courses/{data.courseId}/{data.moduleId}" />
		<div class="flex items-center justify-between">
			<div>
				<h1 class="mb-2 text-2xl font-bold">{data.evaluation.title}</h1>
				<p class="text-muted-foreground">Historial de intentos</p>
			</div>

			{#if canAttempt && !loading && userId}
				<Button onclick={startNewAttempt} class="bg-blue-600 hover:bg-blue-700">
					<Plus class="mr-2 h-4 w-4" />
					Nuevo Intento
				</Button>
			{/if}
		</div>
	</div>

	<!-- Summary Stats -->
	<div class="mb-8 grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-4">
		<Card.Root>
			<Card.Header class="pb-2">
				<div class="flex items-center gap-2">
					<FileText class="h-5 w-5 text-blue-600" />
					<Card.Title class="text-sm font-medium">Total de Intentos</Card.Title>
				</div>
			</Card.Header>
			<Card.Content>
				<div class="text-2xl font-bold">{attempts.length}</div>
				<p class="text-muted-foreground text-xs">
					{data.evaluation.max_attempts
						? `de ${data.evaluation.max_attempts} máximo`
						: 'Sin límite'}
				</p>
			</Card.Content>
		</Card.Root>

		{#if bestAttempt}
			<Card.Root>
				<Card.Header class="pb-2">
					<div class="flex items-center gap-2">
						<Trophy class="h-5 w-5 text-yellow-600" />
						<Card.Title class="text-sm font-medium">Mejor Puntuación</Card.Title>
					</div>
				</Card.Header>
				<Card.Content>
					<div class="text-2xl font-bold {passed ? 'text-green-600' : 'text-red-600'}">
						{bestAttempt.score}/{bestAttempt.total_points}
					</div>
					<p class="text-muted-foreground text-xs">
						{Math.round((bestAttempt.score / bestAttempt.total_points) * 100)}%
					</p>
				</Card.Content>
			</Card.Root>
		{/if}

		<Card.Root>
			<Card.Header class="pb-2">
				<div class="flex items-center gap-2">
					<Target class="h-5 w-5 text-green-600" />
					<Card.Title class="text-sm font-medium">Estado</Card.Title>
				</div>
			</Card.Header>
			<Card.Content>
				<Badge class={passed ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}>
					{passed ? 'Aprobado' : 'No Aprobado'}
				</Badge>
				<p class="text-muted-foreground mt-1 text-xs">
					Mínimo: {data.evaluation.passing_score}%
				</p>
			</Card.Content>
		</Card.Root>

		<Card.Root>
			<Card.Header class="pb-2">
				<div class="flex items-center gap-2">
					<Clock class="h-5 w-5 text-purple-600" />
					<Card.Title class="text-sm font-medium">Límite de Tiempo</Card.Title>
				</div>
			</Card.Header>
			<Card.Content>
				<div class="text-2xl font-bold">
					{data.evaluation.time_limit ? formatDuration(data.evaluation.time_limit) : 'Sin límite'}
				</div>
			</Card.Content>
		</Card.Root>
	</div>

	<!-- Attempts List -->
	<div class="space-y-4">
		<div class="flex items-center justify-between">
			<h2 class="text-xl font-semibold">Intentos realizados</h2>
			{#if !canAttempt}
				<Badge variant="secondary" class="bg-orange-100 text-orange-800">
					Máximo de intentos alcanzado
				</Badge>
			{/if}
		</div>

		{#if attempts.length === 0}
			<Card.Root>
				<Card.Content class="py-12 text-center">
					<FileText class="text-muted-foreground mx-auto mb-4 h-12 w-12" />
					<h3 class="mb-2 text-lg font-semibold">No hay intentos aún</h3>
					<p class="text-muted-foreground mb-4">
						¿Listo para comenzar tu primer intento de este examen?
					</p>
					{#if canAttempt && userId}
						<Button onclick={startNewAttempt} disabled={loading}>
							<Plus class="mr-2 h-4 w-4" />
							{loading ? 'Iniciando...' : 'Comenzar Primer Intento'}
						</Button>
					{/if}
				</Card.Content>
			</Card.Root>
		{:else}
			{#each sortedAttempts as attempt, index (attempt.id)}
				<Card.Root class="transition-shadow hover:shadow-md">
					<Card.Header>
						<div class="flex items-center justify-between">
							<div class="flex items-center gap-3">
								<div class="flex items-center gap-2">
									<span class="font-semibold">Intento #{attempts.length - index}</span>
									<Badge class={getStatusColor(attempt)}>
										{getStatusText(attempt)}
									</Badge>
									{#if attempt.id === bestAttempt?.id && attempts.length > 1}
										<Badge class="bg-yellow-100 text-yellow-800">
											<Trophy class="mr-1 h-3 w-3" />
											Mejor
										</Badge>
									{/if}
								</div>
							</div>

							{#if attempt.submitted_at}
								<div class="text-right">
									<div
										class="text-lg font-bold {attempt.passed ? 'text-green-600' : 'text-red-600'}"
									>
										{attempt.score}/{attempt.total_points}
									</div>
									<div class="text-muted-foreground text-sm">
										{Math.round((attempt.score / attempt.total_points) * 100)}%
									</div>
								</div>
							{:else}
								<Badge variant="secondary" class="bg-yellow-100 text-yellow-800">En progreso</Badge>
							{/if}
						</div>
					</Card.Header>

					<Card.Content>
						<div class="grid grid-cols-1 gap-4 text-sm md:grid-cols-3">
							<div class="flex items-center gap-2">
								<Calendar class="text-muted-foreground h-4 w-4" />
								<div>
									<div class="font-medium">Iniciado</div>
									<div class="text-muted-foreground">
										{formatDate(attempt.started_at)}
									</div>
								</div>
							</div>

							{#if attempt.submitted_at}
								<div class="flex items-center gap-2">
									<Clock class="text-muted-foreground h-4 w-4" />
									<div>
										<div class="font-medium">Completado</div>
										<div class="text-muted-foreground">
											{formatDate(attempt.submitted_at)}
										</div>
									</div>
								</div>
							{:else}
								<div class="flex items-center gap-2">
									<X class="h-4 w-4 text-yellow-600" />
									<div>
										<div class="font-medium text-yellow-600">Sin completar</div>
										<div class="text-muted-foreground">Intento en progreso</div>
									</div>
								</div>
							{/if}
						</div>
					</Card.Content>

					<Card.Footer class="flex gap-2">
						{#if attempt.submitted_at}
							<Button
								variant="outline"
								href="{evaluationUrl}/attempt/{attempt.id}/results"
								class="flex-1"
								data-sveltekit-reload
							>
								<Eye class="mr-2 h-4 w-4" />
								Ver Resultados
							</Button>
						{:else}
							<Button
								href="{evaluationUrl}/attempt/{attempt.id}"
								class="flex-1"
								data-sveltekit-reload
							>
								Continuar Intento
							</Button>
						{/if}
					</Card.Footer>
				</Card.Root>
			{/each}
		{/if}
	</div>
</div>

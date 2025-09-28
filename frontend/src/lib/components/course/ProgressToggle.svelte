<script lang="ts">
	import { UserProgressController } from '$lib/controllers/userProgress';
	import { Button } from '$lib/components/ui/button';
	import { CheckCircle, Circle } from '@lucide/svelte';
	import { toast } from 'svelte-sonner';

	interface Props {
		userId: number;
		courseId: number;
		moduleId: number;
		contentId: number;
		isCompleted: boolean;
		accessToken: string;
		onProgressChange?: (contentId: number, completed: boolean) => void;
	}

	let {
		userId,
		courseId,
		moduleId,
		contentId,
		isCompleted = false,
		accessToken,
		onProgressChange
	}: Props = $props();

	let loading = $state(false);
	let completed = $state(isCompleted);

	async function toggleCompletion() {
		if (loading) return;

		loading = true;
		const progressController = new UserProgressController(accessToken);

		try {
			if (completed) {
				await progressController.markContentIncomplete(userId, courseId, moduleId, contentId);
				completed = false;
				toast.success('Contenido marcado como incompleto');
			} else {
				await progressController.markContentComplete(userId, courseId, moduleId, contentId);
				completed = true;
				toast.success('Contenido marcado como completado');
			}

			// Call the callback if provided
			if (onProgressChange) {
				onProgressChange(contentId, completed);
			}
		} catch (error) {
			console.error('Error al actualizar el progreso:', error);
			toast.error('Error al actualizar el progreso');
		} finally {
			loading = false;
		}
	}
</script>

<Button
	variant="ghost"
	size="sm"
	class="flex h-auto items-center gap-2 p-2"
	onclick={toggleCompletion}
	disabled={loading}
>
	{#if completed}
		<CheckCircle class="h-5 w-5 text-green-600" />
		<span class="text-sm font-medium text-green-700">Completado</span>
	{:else}
		<Circle class="h-5 w-5 text-gray-400" />
		<span class="text-sm font-medium text-gray-600">Marcar como completado</span>
	{/if}
</Button>

<script lang="ts">
	import VideoPlayer from '$lib/components/course/VideoPlayer.svelte';
	import TextWithLinks from '$lib/components/ui/TextWithLinks.svelte';
	import ProgressToggle from '$lib/components/course/ProgressToggle.svelte';
	import { Button } from '$lib/components/ui/button';
	import { ChevronRight } from '@lucide/svelte';
	import { goto, invalidateAll } from '$app/navigation';
	import type { PageProps } from './$types';
	import type { Module, ModuleContent } from '$lib/types/models/course';
	import { UserProgressController } from '$lib';

	let { data }: PageProps = $props();

	let content = $derived(data.content);
	const modules = data.modules;
	const courseId = data.courseId;
	const moduleId = data.moduleId;

	// Reactive state for completion status
	let isCompleted = $state(data.isCompleted);

	// Handle progress change
	function handleProgressChange(contentId: number, completed: boolean) {
		isCompleted = completed;
		// Invalidate to refresh course progress data
		invalidateAll();
	}

	// Navigation logic to find next content
	function getNextContent(): { url: string; label: string } | null {
		// Validate that we have all necessary data
		if (!modules || !Array.isArray(modules) || modules.length === 0) return null;

		// Find current module
		const currentModule = modules.find((m: Module) => m.id === moduleId);
		if (!currentModule?.contents) return null;

		// Sort contents by order
		const sortedContents = [...currentModule.contents].sort((a, b) => a.order - b.order);
		const currentContentIndex = sortedContents.findIndex((c) => c.id === content.id);

		if (currentContentIndex === -1) return null;

		// Check if there's a next content in the current module
		if (currentContentIndex < sortedContents.length - 1) {
			const nextContent = sortedContents[currentContentIndex + 1];
			return {
				url: `/courses/${courseId}/${moduleId}/${nextContent.id}`,
				label: `Siguiente: ${nextContent.title}`
			};
		}

		// If no more content in current module, find next module
		const sortedModules = [...modules].sort((a, b) => a.order - b.order);
		const currentModuleIndex = sortedModules.findIndex((m) => m.id === moduleId);

		if (currentModuleIndex < sortedModules.length - 1) {
			const nextModule = sortedModules[currentModuleIndex + 1];
			if (nextModule.contents && nextModule.contents.length > 0) {
				const firstContent = [...nextModule.contents].sort((a, b) => a.order - b.order)[0];
				return {
					url: `/courses/${courseId}/${nextModule.id}/${firstContent.id}`,
					label: `Siguiente módulo: ${nextModule.title}`
				};
			} else {
				// If next module has no content, go to module page
				return {
					url: `/courses/${courseId}/${nextModule.id}`,
					label: `Siguiente módulo: ${nextModule.title}`
				};
			}
		}

		return null; // No more content
	}

	const nextContentInfo = $derived(getNextContent());

	async function handleCompleteAndNext() {
		// Mark as completed
		const progressController = new UserProgressController();
		progressController
			.markContentComplete(data.userId, courseId, moduleId, content.id)
			.then(() => {
				isCompleted = true;
			})
			.catch((error) => {
				console.error('Error marking content as complete:', error);
			});

		if (nextContentInfo) {
			await goto(nextContentInfo.url);
		}
	}
</script>

<div class="flex flex-col gap-y-8">
	<div class="flex flex-col gap-y-4">
		<div class="flex items-start justify-between">
			<div class="flex flex-col gap-y-2">
				<h1 class="text-3xl font-bold">{content.title}</h1>
				<p class="text-muted-foreground">{content.description}</p>
			</div>
		</div>

		<!-- Completion Status Badge -->
		{#if isCompleted}
			<div class="inline-flex w-fit items-center gap-2 rounded-full bg-green-100 px-3 py-1 text-sm">
				<div class="h-2 w-2 rounded-full bg-green-600"></div>
				<span class="font-medium text-green-800">Completado</span>
			</div>
		{/if}
	</div>

	{#if content?.media_url}
		<div class="">
			<VideoPlayer url={content.media_url} />
		</div>
	{/if}

	{#if content.body}
		<div class="prose prose-sm max-w-none">
			<TextWithLinks text={content.body} class="whitespace-pre-wrap text-sm leading-relaxed" />
		</div>
	{/if}

	<!-- Action Buttons -->
	<div class="border-border mt-8 flex justify-between border-t pt-8">
		<div class="flex gap-3">
			{#if !isCompleted}
				<ProgressToggle
					userId={data.userId}
					courseId={data.courseId}
					moduleId={data.moduleId}
					contentId={content.id}
					{isCompleted}
					accessToken={data.accessToken}
					onProgressChange={handleProgressChange}
				/>
			{/if}
		</div>

		{#if nextContentInfo}
			<Button
				data-sveltekit-reload
				href={nextContentInfo.url}
				class="bg-primary hover:bg-primary/90 flex items-center gap-2"
				onclick={isCompleted ? undefined : handleCompleteAndNext}
			>
				<span>{nextContentInfo.label}</span>
				<ChevronRight class="h-4 w-4" />
			</Button>
		{:else if isCompleted}
			<Button
				href="/courses/{courseId}"
				class="flex items-center gap-2 bg-green-600 hover:bg-green-700"
			>
				<span>Volver al curso</span>
			</Button>
		{/if}
	</div>
</div>

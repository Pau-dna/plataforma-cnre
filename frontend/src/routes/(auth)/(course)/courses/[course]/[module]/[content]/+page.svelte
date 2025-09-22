<script lang="ts">
	import VideoPlayer from '$lib/components/course/VideoPlayer.svelte';
	import { Button } from '$lib/components/ui/button';
	import { ChevronRight } from '@lucide/svelte';
	import type { PageProps } from './$types';
	import type { Module, ModuleContent } from '$lib/types/models/course';

	let { data }: PageProps = $props();

	let content = $state(data.content);
	const modules = data.modules;
	const courseId = data.courseId;
	const moduleId = data.moduleId;

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

	const nextContentInfo = getNextContent();
</script>

<div class="flex flex-col gap-y-8">
	<div class="flex flex-col gap-y-2">
		<h1 class="text-3xl font-bold">{content.title}</h1>
		<p class="text-muted-foreground">{content.description}</p>
	</div>

	{#if content?.media_url}
		<div class="aspect-video w-full">
			<VideoPlayer url={content.media_url} />
		</div>
	{/if}

	{#if content.body}
		<div class="prose prose-sm max-w-none">
			<div class="whitespace-pre-wrap text-sm leading-relaxed">
				{content.body}
			</div>
		</div>
	{/if}

	<!-- Next Content Button -->
	{#if nextContentInfo}
		<div class="flex justify-end mt-8 pt-8 border-t border-border">
			<Button
				href={nextContentInfo.url}
				class="flex items-center gap-2 bg-primary hover:bg-primary/90"
			>
				<span>{nextContentInfo.label}</span>
				<ChevronRight class="h-4 w-4" />
			</Button>
		</div>
	{/if}
</div>
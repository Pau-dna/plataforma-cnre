<script lang="ts">
	import type { ModuleContent } from '$lib/types/models/course';
	import { ContentType } from '$lib/types/models/course';
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$components/ui/card';
	import { Badge } from '$components/ui/badge';
	import { Button } from '$components/ui/button';
	import { PlayCircle, FileText, CheckCircle, Clock, Users, Target, Video } from '@lucide/svelte';
	import { page } from '$app/state';

	type Props = {
		content: ModuleContent;
		active?: boolean;
	};

	const getContentIcon = (type: ContentType) => {
		switch (type) {
			case ContentType.CONTENT:
				return FileText;
			case ContentType.EVALUATION:
				return CheckCircle;
			default:
				return FileText;
		}
	};

	const getContentTypeLabel = (type: ContentType) => {
		switch (type) {
			case ContentType.CONTENT:
				return 'Contenido';
			case ContentType.EVALUATION:
				return 'Evaluación';
			default:
				return 'Contenido';
		}
	};

	const getMediaIcon = (mediaUrl: string) => {
		// Simple check based on file extension
		const extension = mediaUrl.split('.').pop()?.toLowerCase();
		if (extension) {
			if (['mp4', 'mov', 'avi'].includes(extension)) {
				return PlayCircle;
			}
			// Add more media types if needed
		}
		return FileText; // Default icon
	};

	const { content, active }: Props = $props();
</script>

<Card class="group cursor-pointer transition-shadow hover:shadow-md">
	<CardHeader class="pb-3">
		<div class="flex items-start justify-between gap-4">
			<div class="flex min-w-0 flex-1 items-center gap-3">
				<div class="bg-primary/10 text-primary flex-shrink-0 rounded-lg p-2">
					{#if content.type === ContentType.CONTENT}
						<Video class="size-4" />
					{:else}
						<CheckCircle class="size-4" />
					{/if}
				</div>
				<div class="min-w-0 flex-1">
					<div class="mb-1 flex items-center gap-2">
						<span class="text-muted-foreground text-sm font-medium">#{content.order}</span>
						<Badge variant="outline" class="text-xs">
							{getContentTypeLabel(content.type)}
						</Badge>
					</div>
					<CardTitle class="group-hover:text-primary text-pretty text-lg transition-colors">
						{content.title}
					</CardTitle>
				</div>
			</div>

			{#if content.type === ContentType.CONTENT && content.media_url}
				<div class="text-muted-foreground flex flex-shrink-0 items-center gap-1">
					<Video class="size-4" />
					<span class="text-xs">Video</span>
				</div>
			{/if}
		</div>

		<CardDescription class="text-pretty">{content.description}</CardDescription>
	</CardHeader>

	<CardContent class="pt-0">
		{#if content.type === ContentType.EVALUATION}
			<div class="text-muted-foreground mb-4 flex flex-wrap gap-4 text-sm">
				{#if content.question_count}
					<div class="flex items-center gap-1">
						<Users class="h-4 w-4" />
						<span>{content.question_count} preguntas</span>
					</div>
				{/if}
				{#if content.passing_score}
					<div class="flex items-center gap-1">
						<Target class="h-4 w-4" />
						<span>{content.passing_score} pts mínimo</span>
					</div>
				{/if}
				{#if content.time_limit}
					<div class="flex items-center gap-1">
						<Clock class="h-4 w-4" />
						<span>{content.time_limit} min</span>
					</div>
				{/if}
				{#if content.max_attempts}
					<div class="flex items-center gap-1">
						<CheckCircle class="h-4 w-4" />
						<span>{content.max_attempts} intentos</span>
					</div>
				{/if}
			</div>
		{/if}

		<!--
		{#if content.body && content.type === ContentType.CONTENT}
			<p class="text-muted-foreground mb-4 line-clamp-2 text-sm">{content.body}</p>
		{/if}
        -->

		<div class="flex items-center justify-between">
			<div class="text-muted-foreground text-xs">
				Actualizado: {content.updated_at
					? new Date(content.updated_at).toLocaleDateString('es-ES')
					: 'N/A'}
			</div>

			<Button
				variant="ghost"
				size="sm"
				class="group-hover:bg-primary group-hover:text-primary-foreground transition-colors"
				href={content.type === ContentType.EVALUATION ? `/courses/${page.params.course}/${page.params.module}/${content.id}/start` : `/courses/${page.params.course}/${page.params.module}/${content.id}`}
			>
				{content.type === ContentType.EVALUATION ? 'Iniciar Evaluación' : 'Ver Contenido'}
			</Button>
		</div>
	</CardContent>
</Card>

import { ContentController } from '$lib';
import { UserProgressController } from '$lib/controllers/userProgress';
import type { PageServerLoad } from './$types';

export const load = (async ({ locals, params, parent }) => {
	const contentController = new ContentController(locals.accessToken || '');
	const progressController = new UserProgressController(locals.accessToken || '');
	
	const contentId = parseInt(params.content);
	const userId = locals.user.id;

	const [content, parentData] = await Promise.all([
		contentController.getContent(contentId),
		parent()
	]);

	// Check if this content is completed
	let isCompleted = false;
	try {
		isCompleted = await progressController.isContentCompleted(userId, contentId);
	} catch (error) {
		console.log('Content not yet tracked:', error);
		isCompleted = false;
	}

	return {
		content,
		modules: parentData.modules,
		courseId: parseInt(params.course),
		moduleId: parseInt(params.module),
		isCompleted,
		userId,
		accessToken: locals.accessToken || ''
	};
}) satisfies PageServerLoad;

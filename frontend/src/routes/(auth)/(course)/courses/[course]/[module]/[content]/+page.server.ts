import { ContentController } from '$lib';
import { UserProgressController } from '$lib/controllers/userProgress';
import type { PageServerLoad } from './$types';

export const load = (async ({ locals, params, parent }) => {
	const contentController = new ContentController(locals.accessToken || '');
	const progressController = new UserProgressController(locals.accessToken || '');

	const contentID = parseInt(params.content);
	const userID = locals?.user?.id as number;

	const [content, parentData, completed] = await Promise.all([
		contentController.getContent(contentID),
		parent(),
		progressController.isContentCompleted(userID, contentID)
	]);

	return {
		content,
		modules: parentData.modules,
		courseId: parseInt(params.course),
		moduleId: parseInt(params.module),
		isCompleted: completed,
		userId: userID,
		accessToken: locals.accessToken || ''
	};
}) satisfies PageServerLoad;

import { ContentController } from '$lib';
import type { PageServerLoad } from './$types';

export const load = (async ({ locals, params, parent }) => {
	const contentController = new ContentController(locals.accessToken || '');
	const content = await contentController.getContent(parseInt(params.content));

	// Get parent layout data to access modules and navigation info
	const parentData = await parent();

	return {
		content,
		modules: parentData.modules,
		courseId: parseInt(params.course),
		moduleId: parseInt(params.module)
	};
}) satisfies PageServerLoad;

import { ContentController } from '$lib';
import type { PageServerLoad } from './$types';

export const load = (async ({ locals, params }) => {
	const contentController = new ContentController(locals.accessToken || '');
	const content = await contentController.getContent(parseInt(params.content));

	return {
		content
	};
}) satisfies PageServerLoad;

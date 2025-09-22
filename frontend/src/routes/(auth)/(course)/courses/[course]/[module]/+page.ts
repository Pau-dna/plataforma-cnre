import { ContentController } from '$lib/controllers/content';
import type { PageLoad } from './$types';

export const load = (async ({ params }) => {
	const contentController = new ContentController();
	const contents = await contentController.getContentsByModule(parseInt(params.module));

	return {
		moduleID: parseInt(params.module),
		contents
	};
}) satisfies PageLoad;

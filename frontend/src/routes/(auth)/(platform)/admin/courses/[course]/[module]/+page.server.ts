import { ModuleController } from '$lib/controllers';
import { ContentController } from '$lib/controllers';
import { error } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ params }) => {
	const courseId = parseInt(params.course);
	const moduleId = parseInt(params.module);

	if (isNaN(courseId) || isNaN(moduleId)) {
		throw error(404, 'Invalid course or module ID');
	}

	try {
		const moduleController = new ModuleController();
		const contentController = new ContentController();

		const module = await moduleController.getModule(moduleId);
		const contents = await contentController.getContentsByModule(moduleId);

		return {
			module,
			contents,
			courseId
		};
	} catch (err) {
		console.error('Error loading module contents:', err);
		throw error(404, 'Module not found');
	}
};
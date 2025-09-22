import { ModuleController } from '$lib/controllers';
import { ContentController } from '$lib/controllers';
import { EvaluationController } from '$lib/controllers';
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
		const evaluationController = new EvaluationController();

		const module = await moduleController.getModule(moduleId);
		const contents = await contentController.getContentsByModule(moduleId);
		const evaluations = await evaluationController.getEvaluationsByModule(moduleId);

		return {
			module,
			contents,
			evaluations,
			courseId
		};
	} catch (err) {
		console.error('Error loading module contents:', err);
		throw error(404, 'Module not found');
	}
};

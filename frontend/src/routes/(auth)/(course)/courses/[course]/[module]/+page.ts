import { ContentController } from '$lib/controllers/content';
import { EvaluationController } from '$lib/controllers/evaluation';
import type { PageLoad } from './$types';

export const load = (async ({ params }) => {
	const contentController = new ContentController();
	const evaluationController = new EvaluationController();
	
	const [contents, evaluations] = await Promise.all([
		contentController.getContentsByModule(parseInt(params.module)),
		evaluationController.getEvaluationsByModule(parseInt(params.module))
	]);

	// Combine content and evaluations as ModuleContent array, sorted by order
	const allContent = [...contents, ...evaluations].sort((a, b) => a.order - b.order);

	return {
		moduleID: parseInt(params.module),
		contents: allContent
	};
}) satisfies PageLoad;

import { ContentController } from '$lib/controllers/content';
import { EvaluationController } from '$lib/controllers/evaluation';
import type { PageLoad } from './$types';

export const load = (async ({ params, parent }) => {
	const contentController = new ContentController();
	const evaluationController = new EvaluationController();

	// Get parent data (from layout) which includes course info
	const parentData = await parent();

	const [contents, evaluations] = await Promise.all([
		contentController.getContentsByModule(parseInt(params.module)),
		evaluationController.getEvaluationsByModule(parseInt(params.module))
	]);

	return {
		...parentData, // Include course, modules, etc. from parent layout
		moduleID: parseInt(params.module),
		contents,
		evaluations
	};
}) satisfies PageLoad;

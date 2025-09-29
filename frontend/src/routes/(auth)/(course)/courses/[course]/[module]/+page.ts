import { UserProgressController } from '$lib';
import { ContentController } from '$lib/controllers/content';
import { EvaluationController } from '$lib/controllers/evaluation';
import type { PageLoad } from './$types';

export const load = (async ({ params, parent }) => {
	const contentController = new ContentController();
	const evaluationController = new EvaluationController();
	const progressController = new UserProgressController();

	// Get parent data (from layout) which includes course info
	const parentData = await parent();

	const [contents, evaluations, progress] = await Promise.all([
		contentController.getContentsByModule(parseInt(params.module)),
		evaluationController.getEvaluationsByModule(parseInt(params.module)),
		progressController.getModuleContentProgress(
			parentData.user?.id as number,
			parseInt(params.module)
		)
	]);

	return {
		...parentData, // Include course, modules, etc. from parent layout
		moduleID: parseInt(params.module),
		contents,
		evaluations,
		progress
	};
}) satisfies PageLoad;

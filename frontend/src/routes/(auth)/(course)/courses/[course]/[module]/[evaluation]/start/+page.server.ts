import { EvaluationController, EvaluationAttemptController } from '$lib';
import type { PageServerLoad } from './$types';

export const load = (async ({ locals, params, parent }) => {
	const evaluationController = new EvaluationController(locals.accessToken || '');
	const evaluationAttemptController = new EvaluationAttemptController(locals.accessToken || '');
	
	const evaluationId = parseInt(params.evaluation);
	const userId = locals.user.id;

	// Get parent layout data to access course and modules info
	const parentData = await parent();

	const [evaluation, canAttempt, userAttempts] = await Promise.all([
		evaluationController.getEvaluation(evaluationId),
		evaluationAttemptController.canUserAttempt(userId, evaluationId),
		evaluationAttemptController.getUserAttempts(userId, evaluationId)
	]);

	return {
		evaluation,
		canAttempt,
		userAttempts,
		courseId: parseInt(params.course),
		moduleId: parseInt(params.module),
		userId: userId,
		course: parentData.course,
		modules: parentData.modules
	};
}) satisfies PageServerLoad;
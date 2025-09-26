import { EvaluationAttemptController } from '$lib/controllers/evaluationAttempt';
import type { PageLoad } from './$types';

export const load = (async ({ params }) => {
	const evaluationAttemptController = new EvaluationAttemptController();

	const attempt = await evaluationAttemptController.getAttempt(parseInt(params.attempt_id));

	return {
		attempt,
		courseId: parseInt(params.course),
		moduleId: parseInt(params.module_id),
		evaluationId: parseInt(params.evaluation_id),
		attemptId: parseInt(params.attempt_id)
	};
}) satisfies PageLoad;

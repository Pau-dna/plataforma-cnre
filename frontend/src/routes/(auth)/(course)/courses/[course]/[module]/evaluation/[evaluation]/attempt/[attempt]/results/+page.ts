import { EvaluationAttemptController } from '$lib/controllers/evaluationAttempt';
import type { PageLoad } from './$types';

export const load = (async ({ params }) => {
	const evaluationAttemptController = new EvaluationAttemptController();

	const attempt = await evaluationAttemptController.getAttempt(parseInt(params.attempt));

	return {
		attempt,
		courseId: parseInt(params.course),
		moduleId: parseInt(params.module),
		evaluationId: parseInt(params.evaluation),
		attemptId: parseInt(params.attempt)
	};
}) satisfies PageLoad;

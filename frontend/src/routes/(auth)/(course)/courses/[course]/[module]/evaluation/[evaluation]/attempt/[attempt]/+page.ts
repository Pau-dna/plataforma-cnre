import { EvaluationAttemptController } from '$lib/controllers/evaluationAttempt';
import type { PageLoad } from './$types';

export const load = (async ({ params }) => {
	const evaluationAttemptController = new EvaluationAttemptController();
	const attempt = await evaluationAttemptController.getAttempt(parseInt(params.attempt));

	return {
		attempt
	};
}) satisfies PageLoad;

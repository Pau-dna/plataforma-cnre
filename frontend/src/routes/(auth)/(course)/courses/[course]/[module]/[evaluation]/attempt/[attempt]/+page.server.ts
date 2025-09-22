import { EvaluationAttemptController } from '$lib';
import type { PageServerLoad } from './$types';

export const load = (async ({ locals, params, parent }) => {
	const evaluationAttemptController = new EvaluationAttemptController(locals.accessToken || '');
	
	const attemptId = parseInt(params.attempt);

	// Get parent layout data
	const parentData = await parent();

	const attempt = await evaluationAttemptController.getAttempt(attemptId);

	// Verify this attempt belongs to the current user
	if (attempt.user_id !== locals.user.id) {
		throw new Error('Unauthorized: This attempt does not belong to you');
	}

	// Check if attempt is still active (not submitted)
	if (attempt.submitted_at) {
		// Redirect to results if already submitted
		throw new Error('Attempt already submitted');
	}

	return {
		attempt,
		courseId: parseInt(params.course),
		moduleId: parseInt(params.module),
		evaluationId: parseInt(params.evaluation),
		course: parentData.course,
		modules: parentData.modules
	};
}) satisfies PageServerLoad;
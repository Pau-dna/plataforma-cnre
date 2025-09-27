import { EvaluationAttemptController } from '$lib/controllers/evaluationAttempt';
import { EvaluationController } from '$lib/controllers/evaluation';
import type { PageLoad } from './$types';
import { get } from 'svelte/store';
import { browser } from '$app/environment';

export const load = (async ({ params }) => {
	// Skip auth check during SSR
	if (!browser) {
		return {
			evaluation: null,
			attempts: [],
			courseId: parseInt(params.course),
			moduleId: parseInt(params.module),
			evaluationId: parseInt(params.evaluation),
			userId: null
		};
	}

	const evaluationAttemptController = new EvaluationAttemptController();
	const evaluationController = new EvaluationController();

	// Get evaluation first, then attempts if user is authenticated
	const evaluation = await evaluationController.getEvaluation(parseInt(params.evaluation));

	// For now, we'll load attempts in the component where we have access to authStore
	return {
		evaluation,
		attempts: [], // Will be loaded in component
		courseId: parseInt(params.course),
		moduleId: parseInt(params.module),
		evaluationId: parseInt(params.evaluation),
		userId: null // Will be set in component
	};
}) satisfies PageLoad;

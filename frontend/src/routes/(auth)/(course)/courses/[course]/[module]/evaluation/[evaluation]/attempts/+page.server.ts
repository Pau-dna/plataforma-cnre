import { EvaluationAttemptController, EvaluationController } from '$lib';
import type { PageServerLoad } from './$types';

export const load = (async ({ params, locals }) => {
    const evaluationAttemptController = new EvaluationAttemptController();
    const evaluationController = new EvaluationController();

    // Get evaluation first, then attempts if user is authenticated
    const evaluation = await evaluationController.getEvaluation(parseInt(params.evaluation));
    const attempts = await evaluationAttemptController.getUserAttempts(
        locals.user.id,
        params.evaluation,
    );

    // For now, we'll load attempts in the component where we have access to authStore
    return {
        evaluation,
        attempts: attempts,
        courseId: parseInt(params.course),
        moduleId: parseInt(params.module),
        evaluationId: parseInt(params.evaluation),
        userId: locals.user.id,
    };
}) satisfies PageServerLoad;
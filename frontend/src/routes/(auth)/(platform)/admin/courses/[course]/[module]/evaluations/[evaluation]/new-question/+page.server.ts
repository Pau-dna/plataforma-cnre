import { EvaluationController } from '$lib/controllers';
import { error } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ params }) => {
	const courseId = parseInt(params.course);
	const moduleId = parseInt(params.module);
	const evaluationId = parseInt(params.evaluation);

	if (isNaN(courseId) || isNaN(moduleId) || isNaN(evaluationId)) {
		throw error(404, 'Invalid course, module, or evaluation ID');
	}

	try {
		const evaluationController = new EvaluationController();
		const evaluation = await evaluationController.getEvaluation(evaluationId);

		return {
			evaluation,
			courseId,
			moduleId
		};
	} catch (err) {
		console.error('Error al cargar evaluación:', err);
		throw error(404, 'Evaluation not found');
	}
};

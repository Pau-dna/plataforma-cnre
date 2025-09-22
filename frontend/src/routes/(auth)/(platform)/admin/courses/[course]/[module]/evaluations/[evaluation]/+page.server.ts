import { EvaluationController, QuestionController } from '$lib/controllers';
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
		const questionController = new QuestionController();

		const evaluation = await evaluationController.getEvaluation(evaluationId);
		const questions = await questionController.getQuestionsByEvaluation(evaluationId);

		return {
			evaluation,
			questions,
			courseId,
			moduleId
		};
	} catch (err) {
		console.error('Error loading evaluation:', err);
		throw error(404, 'Evaluation not found');
	}
};
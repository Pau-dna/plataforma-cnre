import { QuestionController, AnswerController } from '$lib/controllers';
import { error } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ params }) => {
	const courseId = parseInt(params.course);
	const moduleId = parseInt(params.module);
	const evaluationId = parseInt(params.evaluation);
	const questionId = parseInt(params.question);

	if (isNaN(courseId) || isNaN(moduleId) || isNaN(evaluationId) || isNaN(questionId)) {
		throw error(404, 'Invalid course, module, evaluation, or question ID');
	}

	try {
		const questionController = new QuestionController();
		const answerController = new AnswerController();

		const question = await questionController.getQuestionWithAnswers(questionId);

		return {
			question,
			courseId,
			moduleId,
			evaluationId
		};
	} catch (err) {
		console.error('Error loading question:', err);
		throw error(404, 'Question not found');
	}
};

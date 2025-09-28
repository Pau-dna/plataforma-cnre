import type { PageServerLoad } from './$types';
import { EvaluationController, ModuleController } from '$lib/controllers';

export const load: PageServerLoad = async ({ params, locals }) => {
	const courseId = Number(params.course);
	const moduleId = Number(params.module);
	const evaluationId = Number(params.evaluation);

	const evaluationController = new EvaluationController();
	const moduleController = new ModuleController();

	try {
		const [evaluation, module] = await Promise.all([
			evaluationController.getEvaluation(evaluationId),
			moduleController.getModule(moduleId)
		]);

		return {
			evaluation,
			module,
			courseId,
			moduleId,
			evaluationId
		};
	} catch (error) {
		console.error('Error loading evaluation for settings:', error);
		throw error;
	}
};
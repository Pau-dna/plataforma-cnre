import { ModuleController } from '$lib/controllers';
import { error } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ params }) => {
	const courseId = parseInt(params.course);
	const moduleId = parseInt(params.module);

	if (isNaN(courseId) || isNaN(moduleId)) {
		throw error(404, 'Invalid course or module ID');
	}

	try {
		const moduleController = new ModuleController();
		const module = await moduleController.getModule(moduleId);

		return {
			module,
			courseId,
			moduleId
		};
	} catch (err) {
		console.error('Error al cargar el módulo:', err);
		throw error(404, 'Module not found');
	}
};

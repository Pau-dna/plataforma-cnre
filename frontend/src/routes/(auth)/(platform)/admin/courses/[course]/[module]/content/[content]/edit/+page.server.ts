import type { PageServerLoad } from './$types';
import { ContentController, ModuleController } from '$lib/controllers';

export const load: PageServerLoad = async ({ params, locals }) => {
	const courseId = Number(params.course);
	const moduleId = Number(params.module);
	const contentId = Number(params.content);

	const contentController = new ContentController();
	const moduleController = new ModuleController();

	try {
		const [content, module] = await Promise.all([
			contentController.getContent(contentId),
			moduleController.getModule(moduleId)
		]);

		return {
			content,
			module,
			courseId,
			moduleId,
			contentId
		};
	} catch (error) {
		console.error('Error al cargar el contenido para editar:', error);
		throw error;
	}
};

import { CourseController } from '$lib/controllers/course';
import { ModuleController } from '$lib/controllers/module';
import type { LayoutServerLoad } from './$types';

export const load = (async () => {
	const courseController = new CourseController();
	const moduleController = new ModuleController();

	const course = (await courseController.getCourses())[0];
	const modules = await moduleController.getModules(course.id);

	return {
		course,
		modules
	};
}) satisfies LayoutServerLoad;

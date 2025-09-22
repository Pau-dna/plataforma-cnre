import type { Course, Module } from '$lib';
import { CourseController } from '$lib/controllers/course';
import { ModuleController } from '$lib/controllers/module';
import type { LayoutServerLoad } from './$types';

export const load = (async ({ params }) => {
	const courseController = new CourseController();
	const moduleController = new ModuleController();

	let course: Course = null as unknown as Course;
	let modules: Module[] = [];

	try {
		course = await courseController.getCourse(params.course);
		console.log(course);

		modules = await moduleController.getModulesByCourse(course.id);
		console.log(modules);
	} catch (error) {
		console.log(error);
	}

	return {
		course,
		modules
	};
}) satisfies LayoutServerLoad;

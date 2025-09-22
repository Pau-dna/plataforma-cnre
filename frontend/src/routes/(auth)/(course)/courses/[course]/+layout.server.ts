import { EnrollmentController, type Course, type Module } from '$lib';
import { CourseController } from '$lib/controllers/course';
import { ModuleController } from '$lib/controllers/module';
import type { LayoutServerLoad } from './$types';

export const load = (async ({ params, locals }) => {
	const courseController = new CourseController(locals?.accessToken || "");
	const moduleController = new ModuleController(locals?.accessToken || "");
	const enrollmentController = new EnrollmentController(locals?.accessToken || "");

	const [course, modules, enrollment] = await Promise.all([
		courseController.getCourse(parseInt(params.course)),
		moduleController.getModulesByCourse(parseInt(params.course)),
		enrollmentController.getUserCourseEnrollment(locals.user.id, parseInt(params.course))
	])
	
	return {
		course,
		modules,
		enrollment
	};
}) satisfies LayoutServerLoad;

import { EnrollmentController, type Course, type Module } from '$lib';
import { CourseController } from '$lib/controllers/course';
import { ModuleController } from '$lib/controllers/module';
import type { LayoutServerLoad } from './$types';

export const load = (async ({ params, locals }) => {
	const courseController = new CourseController(locals?.accessToken || '');
	const moduleController = new ModuleController(locals?.accessToken || '');
	const enrollmentController = new EnrollmentController(locals?.accessToken || '');

	const [enrollment, modules] = await Promise.all([
		enrollmentController.getUserCourseEnrollment(locals.user.id, parseInt(params.course)),
		moduleController.getModulesByCourse(parseInt(params.course))
	]);

	return {
		course: enrollment.course as Course,
		modules,
		enrollment
	};
}) satisfies LayoutServerLoad;

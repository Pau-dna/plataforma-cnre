import { EnrollmentController, type Course, type Module } from '$lib';
import { CourseController } from '$lib/controllers/course';
import { ModuleController } from '$lib/controllers/module';
import { UserProgressController } from '$lib/controllers/userProgress';
import type { LayoutServerLoad } from './$types';

export const load = (async ({ params, locals }) => {
	
	const moduleController = new ModuleController(locals?.accessToken || '');
	const enrollmentController = new EnrollmentController(locals?.accessToken || '');
	const progressController = new UserProgressController(locals?.accessToken || '');

	const courseId = parseInt(params.course);
	const userID = locals?.user?.id as number;

	const [enrollment, modules, courseProgress, courseProgressPercentage] = await Promise.all([
		enrollmentController.getUserCourseEnrollment(userID, courseId),
		moduleController.getModulesByCourse(courseId),
		progressController.getUserCourseProgress(userID, courseId),
		progressController.calculateCourseProgress(userID, courseId)
	]);

	return {
		course: enrollment.course as Course,
		modules,
		enrollment,
		courseProgress,
		courseProgressPercentage: Math.round(courseProgressPercentage)
	};
}) satisfies LayoutServerLoad;

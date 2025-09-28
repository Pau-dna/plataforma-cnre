import { EnrollmentController, type Course, type Module } from '$lib';
import { CourseController } from '$lib/controllers/course';
import { ModuleController } from '$lib/controllers/module';
import { UserProgressController } from '$lib/controllers/userProgress';
import type { LayoutServerLoad } from './$types';

export const load = (async ({ params, locals }) => {
	const courseController = new CourseController(locals?.accessToken || '');
	const moduleController = new ModuleController(locals?.accessToken || '');
	const enrollmentController = new EnrollmentController(locals?.accessToken || '');
	const progressController = new UserProgressController(locals?.accessToken || '');

	const courseId = parseInt(params.course);
	const userId = locals.user.id;

	const [enrollment, modules, courseProgress, courseProgressPercentage] = await Promise.all([
		enrollmentController.getUserCourseEnrollment(userId, courseId),
		moduleController.getModulesByCourse(courseId),
		progressController.getUserCourseProgress(userId, courseId),
		progressController.calculateCourseProgress(userId, courseId)
	]);

	// Calculate progress for each module
	const modulesWithProgress = await Promise.all(
		modules.map(async (module) => {
			try {
				const moduleProgressPercentage = await progressController.calculateModuleProgress(
					userId,
					module.id
				);
				const moduleProgress = await progressController.getUserModuleProgress(userId, module.id);

				// Get content completion status
				const contentsWithProgress = await Promise.all(
					(module.contents || []).map(async (content) => {
						try {
							const isCompleted = await progressController.isContentCompleted(userId, content.id);
							return { ...content, isCompleted };
						} catch {
							return { ...content, isCompleted: false };
						}
					})
				);

				// Get evaluation pass status
				const evaluationsWithStatus = await Promise.all(
					(module.evaluations || []).map(async (evaluation) => {
						try {
							const hasPassed = await progressController.hasUserPassedEvaluation(
								userId,
								evaluation.id
							);
							return { ...evaluation, hasPassed };
						} catch {
							return { ...evaluation, hasPassed: false };
						}
					})
				);

				return {
					...module,
					progressPercentage: Math.round(moduleProgressPercentage),
					contents: contentsWithProgress,
					evaluations: evaluationsWithStatus,
					progressRecords: moduleProgress
				};
			} catch {
				return {
					...module,
					progressPercentage: 0,
					contents: (module.contents || []).map((content) => ({ ...content, isCompleted: false })),
					evaluations: (module.evaluations || []).map((evaluation) => ({
						...evaluation,
						hasPassed: false
					})),
					progressRecords: []
				};
			}
		})
	);

	return {
		course: enrollment.course as Course,
		modules: modulesWithProgress,
		enrollment,
		courseProgress,
		courseProgressPercentage: Math.round(courseProgressPercentage)
	};
}) satisfies LayoutServerLoad;

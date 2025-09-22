import { EnrollmentController } from '$lib';
import { CourseController } from '$lib/controllers/course';
import type { PageServerLoad } from './$types';

export const load = (async ({ locals }) => {
	const enrollmentController = new EnrollmentController(locals?.accessToken || '');
	const enrollments = await enrollmentController.getUserEnrollments(locals.user.id);

	return { enrollments };
}) satisfies PageServerLoad;

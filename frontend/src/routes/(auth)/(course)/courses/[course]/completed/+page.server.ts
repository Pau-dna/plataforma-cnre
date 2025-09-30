import { EnrollmentController, type Enrollment } from '$lib';
import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load = (async ({ locals, params }) => {
	const enrollmentController = new EnrollmentController();

	let enrollment: Enrollment | null = null;
	let redirectPath = '';
	const fallbackRedirect = '/courses/' + params.course;
	try {
		enrollment = await enrollmentController.getUserCourseEnrollment(
			locals.user?.id as number,
			parseInt(params.course)
		);
		if (!enrollment || enrollment.progress < 100) {
			redirectPath = fallbackRedirect;
		}
	} catch (error) {
		redirectPath = fallbackRedirect;
	}

	if (redirectPath) {
		redirect(307, redirectPath);
	}

	return { enrollment: enrollment as Enrollment };
}) satisfies PageServerLoad;

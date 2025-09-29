import { EnrollmentController } from '$lib';
import type { PageServerLoad } from './$types';
import { error } from '@sveltejs/kit';

export const load = (async ({ params, parent }) => {
	// Get parent data which includes course info and enrollment
	const parentData = await parent();

	// Check if the enrollment is actually completed
	if (!parentData.enrollment.completed_at) {
		throw error(404, 'Curso no completado');
	}

	// Return the enrollment data needed for the completion page
	return {
		...parentData, // Includes course, enrollment, and modules from parent
		enrollment: parentData.enrollment
	};
}) satisfies PageServerLoad;

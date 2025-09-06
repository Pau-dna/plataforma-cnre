import { CourseController } from '$lib/controllers/course';
import type { PageServerLoad } from './$types';

export const load = (async () => {

    const courseController = new CourseController();
    const courses = await courseController.getCourses();

    return { courses };
}) satisfies PageServerLoad;
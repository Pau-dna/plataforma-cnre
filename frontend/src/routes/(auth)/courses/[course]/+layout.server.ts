import { CourseController } from '$lib/controllers/course';
import type { LayoutServerLoad } from './$types';

export const load = (async () => {

    const courseController = new CourseController()

    const course = (await courseController.getCourses())[0]

    return {
course
    };
}) satisfies LayoutServerLoad;
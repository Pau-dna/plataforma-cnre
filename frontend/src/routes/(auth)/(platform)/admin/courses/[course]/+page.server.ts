import { CourseController, ModuleController } from '$lib';
import type { PageServerLoad } from './$types';

export const load = (async ({ params }) => {

    const courseController = new CourseController();
    const moduleController = new ModuleController();

    const course = await courseController.getCourse(params.course);
    const modules = await moduleController.getModulesByCourse(course.id);
    
    return {
        course,
        modules
    };
}) satisfies PageServerLoad;
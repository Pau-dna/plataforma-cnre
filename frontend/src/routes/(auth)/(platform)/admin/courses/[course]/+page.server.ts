import { CourseController, ModuleController } from '$lib';
import type { PageServerLoad } from './$types';

export const load = (async ({ params }) => {

    const courseController = new CourseController();
    const moduleController = new ModuleController();

    const course = await courseController.getCourse(Number(params.course));
    const modules = await moduleController.getModulesByCourse(course.id);
    
    // Ensure modules are sorted by order
    modules.sort((a, b) => a.order - b.order);
    
    return {
        course,
        modules
    };
}) satisfies PageServerLoad;
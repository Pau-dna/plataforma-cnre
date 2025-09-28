import { UserProgressController } from '$lib';
import type { PageServerLoad } from './$types';

export const load = (async ({ params, locals }) => {

    const progressController = new UserProgressController();
    const progress = await progressController.getComprehensiveCourseProgress(locals.user?.id as number, parseInt(params.course));

    return { progress };
}) satisfies PageServerLoad;
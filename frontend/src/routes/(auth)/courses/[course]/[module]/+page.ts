import type { PageLoad } from './$types';

export const load = (async ({ params }) => {
    return {
        moduleID: parseInt(params.module)
    };
}) satisfies PageLoad;
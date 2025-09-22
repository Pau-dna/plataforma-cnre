import type { PageServerLoad } from './$types';
import { redirect } from '@sveltejs/kit';
import { AuthCookies } from '$lib/server/auth-cookies';

export const load = (async ({ cookies }) => {
	AuthCookies.deleteAuthCookies(cookies);
	redirect(303, '/login');
}) satisfies PageServerLoad;

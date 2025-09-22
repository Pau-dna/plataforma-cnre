import { authCookiesManager } from '$lib/server/cookies/manager';
import type { User } from '$lib/types/models/user';
import type { LayoutServerLoad } from './$types';

export const load = (async ({ locals, cookies }) => {
	const user: User | undefined = locals?.user;
	const { accessToken } = authCookiesManager.getTokens(cookies);

	return { user, accessToken };
}) satisfies LayoutServerLoad;

import { redirect, type Handle } from '@sveltejs/kit';
import type { User } from '$lib/types/models/user';
import { authCookiesManager } from '$lib/server/cookies/manager';
import { AuthController } from '$lib/controllers/auth';

export const handle: Handle = async ({ event, resolve }) => {
	const isLogin =
		event.url.pathname === '/login' ||
		event.url.pathname === '/authorize' ||
		event.url.pathname === '/logout';

	if (isLogin) {
		if (event.url.pathname === '/login' && event.request.method === 'GET') {
			authCookiesManager.logout(event.cookies);
		}

		return await resolve(event);
	}

	const hasAuthCookies = authCookiesManager.isAuthenticated(event.cookies);
	if (!hasAuthCookies) {
		// Only add redirect parameter if it's not the my-courses page, or if my-courses page has search parameters
		if ((event.url.pathname !== '/' && event.url.pathname !== '/my-courses') || event.url.search) {
			const redirectTo = btoa(event.url.pathname + event.url.search);
			redirect(303, `/login?redirect=${redirectTo}`);
		} else {
			redirect(303, '/login');
		}
	}

	if (event.url.pathname === '/') {
		redirect(303, '/my-courses');
	}

	const authTokens = authCookiesManager.getTokens(event.cookies);
	event.locals.accessToken = authTokens?.accessToken;
	event.locals.refreshToken = authTokens?.refreshToken;

	try {
		const authController = new AuthController(authTokens.accessToken);
		const user = (await authController.getUserInfo()) as User;
		event.locals.user = user;
		const response = await resolve(event);
		return response;
	} catch (error) {
		console.error('Error fetching user data:', error);
		authCookiesManager.logout(event.cookies);

		// Only add redirect parameter if it's not the home page, or if home page has search parameters
		if ((event.url.pathname !== '/' && event.url.pathname !== '/home') || event.url.search) {
			const redirectTo = btoa(event.url.pathname + event.url.search);
			redirect(303, `/login?redirect=${redirectTo}`);
		} else {
			redirect(303, '/login');
		}
	}
};

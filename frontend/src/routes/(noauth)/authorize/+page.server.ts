import { AuthController } from '$lib/controllers/auth';
import { authCookiesManager } from '$lib/server/cookies/manager';
import type { PageServerLoad } from './$types';
import { redirect } from '@sveltejs/kit';

type GoogleOAuthResponse = {
	code: string;
	scope: string;
	authuser: string;
	hd: string;
	prompt: string;
};

export const load = (async ({ url, cookies }) => {

	const isAuthenticated = authCookiesManager.isAuthenticated(cookies);
	const redirectParam = url.searchParams.get('redirect');
	let redirectTo: null | string = null;

	if (redirectParam) {
		try {
			redirectTo = atob(redirectParam);
		} catch (error) {
			// If base64 decoding fails, ignore the redirect parameter
			console.warn('Failed to decode redirect parameter.');
			redirectTo = null;
		}
	}

	if (isAuthenticated) {
		redirect(303, redirectTo || '/my-courses');
	}

	const credentialsText = url.searchParams.toString();
	const decoded = new URLSearchParams(credentialsText);
	const credentials = Object.fromEntries(decoded) as GoogleOAuthResponse;

	const authController = new AuthController();

	try {
		const response = await authController.loginWithGoogle(credentials.code);
		authCookiesManager.login(
			cookies,
			response.tokens.access_token,
			response.tokens.refresh_token,
		);

		console.log(response);

		// Redirect to the original destination or default to home
		const destination = redirectTo && redirectTo !== '' ? redirectTo : '/my-courses';
		redirect(303, destination);
	} catch (error) {
		console.log(error);
		redirect(303, '/logout');
	}
}) satisfies PageServerLoad;

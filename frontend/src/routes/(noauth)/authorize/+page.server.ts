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
	const credentialsText = url.searchParams.toString();
	const decoded = new URLSearchParams(credentialsText);
	const credentials = Object.fromEntries(decoded) as GoogleOAuthResponse;

	const authController = new AuthController();

	try {
		const response = await authController.login(credentials.code);
		authCookiesManager.setAuthTokens(cookies, response.access_token, response.refresh_token);
	} catch (error) {
		redirect(303, '/logout');
	}

	redirect(303, '/');
}) satisfies PageServerLoad;

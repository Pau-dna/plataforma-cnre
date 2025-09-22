import { AuthCookiesManager } from 'svelte-auth-tools/dist/server/cookies/cookies.js';

export const authCookiesManager = new AuthCookiesManager({
	cookies: {
		accessTokenCookieName: 'access_token',
		refreshTokenCookieName: 'refresh_token',
		maxAgeSeconds: 3600,
		sameSite: 'lax'
	}
});

import type { User } from '$lib/types/models/user';
import type { AuthTokens } from '$lib/types/tokens';

/**
 * A generic authentication store for managing user sessiwon and tokens.
 *
 * @template T - The type representing user data.
 */
export class AuthStore<T> {
	private userData: T | null = $state(null);
	private tokens: AuthTokens | null = $state(null);

	readonly authenticated = $derived(this.isLoggedIn());
	readonly user = $derived(this.userData);

	constructor() {}

	login(accessToken: string, refreshToken: string, user: T | null) {
		if (accessToken == '') throw new Error('El token de acceso no puede estar vacío');
		if (refreshToken == '') throw new Error('El token de actualización no puede estar vacío');
		if (user === null || user === undefined)
			throw new Error('El usuario no puede ser nulo o indefinido');

		this.userData = user;
		this.tokens = { accessToken, refreshToken };
	}

	initiateUserSession(accessToken: string, refreshToken: string, user: T | null) {
		try {
			this.login(accessToken, refreshToken, user);
		} catch (error) {
			console.error('[sv-auth] Error de inicio de sesión en auth store:', error);
			this.logout();
		}
	}

	logout() {
		this.userData = null;
		this.tokens = null;
	}

	isLoggedIn(): boolean {
		const userExists = this.userData !== null && this.userData !== undefined;
		const tokensExist = this.tokens !== null && this.tokens !== undefined;
		const validTokens = this.tokens?.accessToken !== '' && this.tokens?.refreshToken !== '';
		return userExists && tokensExist && validTokens;
	}

	getUser(): T | null {
		return this.userData;
	}

	getAccessToken(): string | null {
		return this.tokens?.accessToken || null;
	}

	getRefreshToken(): string | null {
		return this.tokens?.refreshToken || null;
	}
}

export const authStore = new AuthStore<User>();

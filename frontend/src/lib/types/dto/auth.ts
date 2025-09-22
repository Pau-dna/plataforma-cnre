import type { User } from '../models/user';

export type SignInResponse = {
	user: User;
	tokens: {
		access_token: string;
		refresh_token: string;
		expires_at: number;
	};
};

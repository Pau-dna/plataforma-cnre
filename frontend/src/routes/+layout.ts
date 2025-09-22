import { browser } from '$app/environment';
import type { LayoutLoad } from './$types';
import type { User } from '$lib/types/models/user';
import { authStore } from '$lib/stores/auth.svelte';

export const load = (async ({ data }) => {
	const user: User | undefined = data?.user;
	const accessToken: string | any = data?.accessToken;

	if (browser) {
		if (user === undefined || accessToken === undefined) {
			authStore.logout();
		} else {
			authStore.login(
				accessToken, 
				accessToken, 
				user
			);
		}
	}

	return {
		...data
	};
}) satisfies LayoutLoad;

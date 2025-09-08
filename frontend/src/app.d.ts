// See https://svelte.dev/docs/kit/types#app.d.ts

import type { User } from "$lib/types/models/user";
import { AuthLocals } from "svelte-auth-tools/dist/types/locals.d.ts";


// for information about these interfaces
declare global {
	namespace App {
		// interface Error {}
		interface Locals extends AuthLocals<User> {}
		// interface PageData {}
		// interface PageState {}
		// interface Platform {}
	}
}

export { };

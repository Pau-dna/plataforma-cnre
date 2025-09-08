import type { User } from "$lib/types/models/user";
import { AuthStore } from "svelte-auth-tools/dist/store/auth.svelte";

export const authStore = new AuthStore<User>();
<script lang="ts">
	import type { PageProps } from './$types';
	import * as Card from '$lib/components/ui/card/index.js';
	import { BookOpen, CircleCheckBig, TrendingUp, Users } from '@lucide/svelte';
	import { authStore } from '$lib/stores/auth.svelte';
	import { UserRole } from '$lib/types/models/course';
	import AdminDashboard from '$lib/components/dashboard/AdminDashboard.svelte';
	import StudentDashboard from '$lib/components/dashboard/StudentDashboard.svelte';

	let { data }: PageProps = $props();
</script>

{#if authStore.user?.role === UserRole.ADMIN || authStore.user?.role === UserRole.INSTRUCTOR}
	<AdminDashboard />

{:else}
	<StudentDashboard />
{/if}

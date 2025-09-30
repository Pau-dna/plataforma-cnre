<script lang="ts">
	import type { Enrollment } from '$lib';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	let enrollment = $state<Enrollment>(data?.enrollment);

	function formatDate(date: string | Date | undefined): string {
		if (!date) return 'N/A';
		const d = typeof date === 'string' ? new Date(date) : date;
		return d.toLocaleDateString('es-ES', {
			year: 'numeric',
			month: 'short',
			day: 'numeric'
		});
	}
</script>

<div class="flex min-h-screen w-full items-center justify-center bg-sky-100 p-4 sm:p-6 lg:p-8">
	<div
		class="w-full max-w-[1200px] rounded-lg bg-gradient-to-b from-sky-500 to-pink-500 p-1.5 sm:p-2"
	>
		<div class="flex flex-col gap-12 rounded-lg bg-white p-6 sm:gap-16 sm:p-12 lg:gap-28 lg:p-24">
			<div class="flex flex-col items-start gap-8 sm:gap-12 lg:gap-16">
				<img
					src="/images/logo.png"
					alt="logo-cnre"
					class="mb-4 h-16 sm:mb-6 sm:h-24 lg:mb-8 lg:h-36"
				/>
				<div class="flex flex-col gap-4 sm:gap-8 lg:gap-12">
					<span class="text-muted-foreground text-base font-semibold sm:text-xl lg:text-2xl">
						Por la presente se certifica que
					</span>
					<span class="text-2xl font-bold sm:text-4xl lg:text-6xl">
						{enrollment?.user?.fullname}
					</span>
					<span class="text-muted-foreground text-base font-semibold sm:text-xl lg:text-2xl">
						ha completado exitosamente el curso
					</span>
					<span class="text-xl font-semibold sm:text-2xl lg:text-4xl">
						{enrollment?.course?.title}
					</span>
				</div>
			</div>
			<div class="flex flex-col items-start gap-1 sm:gap-2">
				<span class="text-xl font-semibold sm:text-2xl lg:text-3xl">
					{formatDate(enrollment?.completed_at)}
				</span>
				<span class="text-muted-foreground text-sm font-semibold sm:text-base lg:text-xl">
					Fecha de emisión
				</span>
			</div>
		</div>
	</div>
</div>

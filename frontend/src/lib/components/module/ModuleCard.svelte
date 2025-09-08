<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import { ArrowDown, ArrowUp, BookOpen, CirclePlay, Ellipsis, GripVertical } from '@lucide/svelte';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';

	type Props = {
		id: number;
		title: string;
		description: string;
		actDate?: string;
	};

	const { id, title, description, actDate }: Props = $props();
</script>

<Card.Root>
	<Card.Header>
		<div class="flex items-center justify-between">
			<div class="flex items-center gap-2">
				<Button size="sm" variant="ghost" class="bg-muted">
					<GripVertical class="text-muted-foreground h-4 w-4" />
				</Button>

				<span class="text-muted-foreground font-semibold">Modulo {id}</span>
			</div>
			<DropdownMenu.Root>
				<DropdownMenu.Trigger>
					<Button size="sm" variant="ghost">
						<Ellipsis class="h-4 w-4 leading-none" />
					</Button>
				</DropdownMenu.Trigger>
				<DropdownMenu.Content>
					<DropdownMenu.Group>
						<DropdownMenu.Item>Ver Detalles</DropdownMenu.Item>
						<DropdownMenu.Item>Editar</DropdownMenu.Item>
						<DropdownMenu.Separator />
						<DropdownMenu.Item class="text-destructive">Eliminar</DropdownMenu.Item>
					</DropdownMenu.Group>
				</DropdownMenu.Content>
			</DropdownMenu.Root>
		</div>
		<div class="flex items-center gap-4">
			<div class="flex flex-col gap-0">
				<Button size="sm" variant="ghost">
					<ArrowUp class="h-4 w-4" />
				</Button>
				<Button size="sm" variant="ghost">
					<ArrowDown class="h-4 w-4" />
				</Button>
			</div>
			<div class="flex flex-col gap-1">
				<Card.Title class="text-lg">{title}</Card.Title>
				<Card.Description>{description}</Card.Description>
			</div>
		</div>
	</Card.Header>
	<Card.Content class="flex justify-between w-full">
		<div class="flex items-center gap-4">
			<div class="flex items-center gap-2">
				<CirclePlay class="text-muted-foreground h-4 w-4 leading-none" />
				<span class="text-muted-foreground text-sm leading-none">Video Incluido</span>
			</div>
			<div class="flex items-center gap-2 text-cyan-900">
				<BookOpen class="h-4 w-4 leading-none" />
				<span class="text-sm leading-none">Ver Contenido</span>
			</div>
		</div>
        <span class="text-muted-foreground text-sm leading-none">Actualizado el {actDate}</span>
	</Card.Content>
</Card.Root>

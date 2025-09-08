<script lang="ts">
	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';
	import type { Snippet } from 'svelte';

	type Props = {
		openDelete?: boolean;
		children?: Snippet;
	};

	let { children, openDelete = $bindable() }: Props = $props();
</script>

<AlertDialog.Root bind:open={openDelete}>
	<AlertDialog.Content>
		<AlertDialog.Header>
			<AlertDialog.Title>¿Estás seguro de que deseas eliminar este curso?</AlertDialog.Title>
			<AlertDialog.Description>
				Esta acción no se puede deshacer. Esto eliminará permanentemente el curso y eliminará todos
				los datos asociados.
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Cancel>Cancelar</AlertDialog.Cancel>
			<AlertDialog.Action class="bg-destructive hover:bg-red-800">Continuar</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>

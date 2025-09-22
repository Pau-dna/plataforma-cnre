<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { ArrowRight } from '@lucide/svelte';
	import { googleClientID } from '$lib/constants';

	async function handleLogin() {
		// Google's OAuth 2.0 endpoint for requesting an access token
		const oauth2Endpoint = 'https://accounts.google.com/o/oauth2/v2/auth';

		// Create <form> element to submit parameters to OAuth 2.0 endpoint.
		const form = document.createElement('form');
		form.setAttribute('method', 'GET'); // Send as a GET request.
		form.setAttribute('action', oauth2Endpoint);

		// Parameters to pass to OAuth 2.0 endpoint.
		const params: Record<string, string> = {
			client_id: googleClientID,
			redirect_uri: `${window.location.origin}/authorize`,
			response_type: 'code',
			prompt: 'select_account',
			scope: 'openid profile email',
			include_granted_scopes: 'true'
		};

		// Add form parameters as hidden input values.
		for (const p in params) {
			const input = document.createElement('input');
			input.setAttribute('type', 'hidden');
			input.setAttribute('name', p);
			input.setAttribute('value', params[p]);
			form.appendChild(input);
		}

		// Add form to page and submit it to open the OAuth 2.0 endpoint.
		document.body.appendChild(form);
		form.submit();
	}

</script>

<Card.Root class="w-md border-sky-500/25 shadow-lg shadow-sky-500/25">
	<img src="/images/logo.png" alt="CNRE Logo" class="mx-auto h-40 w-auto" />
	<Card.Content class="flex flex-col items-center gap-4 p-6">
		<Card.Title class="text-h2 text-center"
			>Bienvenido(a) a la Plataforma de Cursos del CNRE</Card.Title
		>
		<Card.Description class="text-center"
			>Ingrese con su correo institucional para acceder a todos los recursos disponibles.</Card.Description
		>
	</Card.Content>
	<Card.Footer>
		<Button class="w-full bg-pink-500 hover:bg-pink-900" onclick={handleLogin}>
			<span>Continuar con Google</span>
			<ArrowRight />
		</Button>
	</Card.Footer>
</Card.Root>

import { ApiError, getErrorMessage } from '$lib/utils/error';
import { toast } from 'svelte-sonner';

/**
 * Standard error handler for API operations
 * Can be used across the application for consistent error handling
 *
 * @param error - The error caught from API calls
 * @param operation - Description of the operation for logging
 * @param showToast - Whether to show a toast notification
 * @returns ApiError instance for further handling
 */
export function handleApiError(
	error: unknown,
	operation: string = 'operation',
	showToast: boolean = true
): ApiError {
	const apiError = ApiError.from(error);

	// Log error for debugging
	console.error(`Error during ${operation}:`, {
		code: apiError.code,
		message: apiError.message,
		payload: apiError.payload
	});

	if (showToast) {
		// Show appropriate toast based on error type
		if (apiError.isValidationError()) {
			toast.error('Datos inválidos', {
				description: apiError.message
			});
		} else if (apiError.isAuthError()) {
			toast.error('Sin autorización', {
				description: 'Debes iniciar sesión para continuar.'
			});
		} else if (apiError.isNotFoundError()) {
			toast.error('No encontrado', {
				description: 'El recurso solicitado no existe.'
			});
		} else if (apiError.isConflictError()) {
			toast.error('Conflicto', {
				description: apiError.message
			});
		} else if (apiError.isNetworkError()) {
			toast.error('Error de conexión', {
				description: 'Verifica tu conexión a internet.'
			});
		} else {
			toast.error('Error', {
				description: apiError.getUserMessage()
			});
		}
	}

	return apiError;
}

/**
 * Wrapper for API calls with consistent error handling
 *
 * @param apiCall - The API call function
 * @param operation - Description of the operation
 * @param onSuccess - Success callback
 * @param onError - Optional error callback
 */
export async function executeApiCall<T>(
	apiCall: () => Promise<T>,
	operation: string,
	onSuccess?: (result: T) => void,
	onError?: (error: ApiError) => void
): Promise<T | null> {
	try {
		const result = await apiCall();
		if (onSuccess) {
			onSuccess(result);
		}
		return result;
	} catch (error) {
		const apiError = handleApiError(error, operation);
		if (onError) {
			onError(apiError);
		}
		return null;
	}
}

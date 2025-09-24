import type { ApiErrorResponse } from '$lib/services/api';

/**
 * API Error class for consistent error handling
 */
export class ApiError extends Error implements ApiErrorResponse {
	code: string;
	payload?: Record<string, any>;

	constructor(code: string, message: string, payload?: Record<string, any>) {
		super(message);
		this.name = 'ApiError';
		this.code = code;
		this.payload = payload;
	}

	/**
	 * Create ApiError from unknown error type
	 */
	static from(error: unknown): ApiError {
		if (error instanceof ApiError) {
			return error;
		}

		if (isApiErrorResponse(error)) {
			return new ApiError(error.code, error.message, error.payload);
		}

		if (error instanceof Error) {
			return new ApiError('CLIENT_ERROR', error.message, { originalError: error.name });
		}

		return new ApiError('UNKNOWN_ERROR', String(error));
	}

	/**
	 * Check if the error is a network error
	 */
	isNetworkError(): boolean {
		return this.code === 'NETWORK_ERROR';
	}

	/**
	 * Check if the error is a validation error
	 */
	isValidationError(): boolean {
		return this.code === 'BAD_REQUEST' || this.code === 'BIND_JSON';
	}

	/**
	 * Check if the error is an authorization error
	 */
	isAuthError(): boolean {
		return this.code === 'UNAUTHORIZED';
	}

	/**
	 * Check if the error is a not found error
	 */
	isNotFoundError(): boolean {
		return this.code === 'NOT_FOUND';
	}

	/**
	 * Check if the error is a conflict error
	 */
	isConflictError(): boolean {
		return this.code === 'CONFLICT';
	}

	/**
	 * Get user-friendly error message
	 */
	getUserMessage(): string {
		switch (this.code) {
			case 'NETWORK_ERROR':
				return 'Error de conexión. Por favor, verifica tu conexión a internet.';
			case 'UNAUTHORIZED':
				return 'No tienes permisos para realizar esta acción.';
			case 'NOT_FOUND':
				return 'El recurso solicitado no fue encontrado.';
			case 'CONFLICT':
				return 'El recurso ya existe o hay un conflicto.';
			case 'BAD_REQUEST':
			case 'BIND_JSON':
				return 'Los datos enviados no son válidos.';
			case 'INTERNAL_SERVER_ERROR':
				return 'Error interno del servidor. Intenta de nuevo más tarde.';
			default:
				return this.message || 'Ha ocurrido un error inesperado.';
		}
	}
}

/**
 * Type guard to check if an error is an ApiErrorResponse
 */
export function isApiErrorResponse(error: unknown): error is ApiErrorResponse {
	return (
		typeof error === 'object' &&
		error !== null &&
		typeof (error as any).code === 'string' &&
		typeof (error as any).message === 'string'
	);
}

/**
 * Cast an unknown error to ApiError
 */
export function castToApiError(error: unknown): ApiError {
	return ApiError.from(error);
}

/**
 * Handle error and return user-friendly message
 */
export function getErrorMessage(error: unknown): string {
	return ApiError.from(error).getUserMessage();
}
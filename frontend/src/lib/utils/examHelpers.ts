import type { EvaluationAttempt } from '$lib/types';

// Helper function to validate attempt data
export function validateAttempt(attempt: EvaluationAttempt): string[] {
	const errors: string[] = [];

	if (!attempt.questions || attempt.questions.length === 0) {
		errors.push('El intento no tiene preguntas disponibles');
	}

	if (attempt.submitted_at && !attempt.answers) {
		errors.push('El intento enviado no tiene respuestas');
	}

	return errors;
}

// Helper function to check if attempt is still active (not submitted or expired)
export function isAttemptActive(attempt: EvaluationAttempt, timeLimit?: number): boolean {
	if (attempt.submitted_at) return false;

	if (timeLimit) {
		const startTime = new Date(attempt.started_at).getTime();
		const now = Date.now();
		const elapsed = now - startTime;
		const timeLimitMs = timeLimit * 60 * 1000;

		return elapsed < timeLimitMs;
	}

	return true;
}

// Helper function to calculate remaining time
export function calculateRemainingTime(attempt: EvaluationAttempt, timeLimit?: number): number {
	if (!timeLimit) return -1; // No time limit

	const startTime = new Date(attempt.started_at).getTime();
	const now = Date.now();
	const elapsed = now - startTime;
	const timeLimitMs = timeLimit * 60 * 1000;

	return Math.max(0, Math.floor((timeLimitMs - elapsed) / 1000));
}

// Helper function to format attempt status
export function getAttemptStatusInfo(attempt: EvaluationAttempt) {
	if (!attempt.submitted_at) {
		return {
			status: 'En progreso',
			color: 'bg-yellow-100 text-yellow-800',
			icon: 'clock'
		};
	}

	if (attempt.passed) {
		return {
			status: 'Aprobado',
			color: 'bg-green-100 text-green-800',
			icon: 'check'
		};
	}

	return {
		status: 'No aprobado',
		color: 'bg-red-100 text-red-800',
		icon: 'x'
	};
}

// Helper function to validate answers before submission
export function validateAnswers(
	questions: any[],
	answers: Record<string, number[]>
): {
	isValid: boolean;
	errors: string[];
	warnings: string[];
} {
	const errors: string[] = [];
	const warnings: string[] = [];

	// Check for unanswered questions
	const unansweredCount = questions.length - Object.values(answers).filter((a) => a.length > 0).length;
	if (unansweredCount > 0) {
		warnings.push(
			`Tienes ${unansweredCount} pregunta${unansweredCount > 1 ? 's' : ''} sin responder`
		);
	}

	// Validate individual answers
	for (const [questionId, selectedOptions] of Object.entries(answers)) {
		const question = questions.find((q) => q.id === Number(questionId));
		if (!question) {
			errors.push(`Pregunta con ID ${questionId} no encontrada`);
			continue;
		}

		if (selectedOptions.length === 0) {
			warnings.push(`La pregunta "${question.text}" no tiene respuesta seleccionada`);
		}

		// Validate option IDs exist
		const validOptionIds = question.answer_options.map((opt: any) => opt.id);
		const invalidOptions = selectedOptions.filter((id) => !validOptionIds.includes(id));
		if (invalidOptions.length > 0) {
			errors.push(`Opciones inválidas seleccionadas en la pregunta: ${question.text}`);
		}
	}

	return {
		isValid: errors.length === 0,
		errors,
		warnings
	};
}

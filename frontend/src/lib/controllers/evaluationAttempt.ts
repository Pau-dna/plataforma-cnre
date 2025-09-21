import { BaseController } from './base';
import type {
	EvaluationAttempt,
	StartEvaluationAttemptDTO,
	SubmitEvaluationAttemptDTO
} from '$lib/types';

export class EvaluationAttemptController extends BaseController {
	/**
	 * Start a new evaluation attempt
	 */
	async startAttempt(attemptData: StartEvaluationAttemptDTO): Promise<EvaluationAttempt> {
		return this.post<EvaluationAttempt>('/api/v1/evaluation-attempts/start', attemptData);
	}

	/**
	 * Submit evaluation attempt with answers
	 */
	async submitAttempt(attemptId: number, submissionData: SubmitEvaluationAttemptDTO): Promise<EvaluationAttempt> {
		return this.post<EvaluationAttempt>(`/api/v1/evaluation-attempts/${attemptId}/submit`, submissionData);
	}

	/**
	 * Get a specific evaluation attempt by ID
	 */
	async getAttempt(id: number): Promise<EvaluationAttempt> {
		return this.get<EvaluationAttempt>(`/api/v1/evaluation-attempts/${id}`);
	}

	/**
	 * Get all attempts for a user and evaluation
	 */
	async getUserAttempts(userId: number, evaluationId: number): Promise<EvaluationAttempt[]> {
		return this.get<EvaluationAttempt[]>(`/api/v1/users/${userId}/evaluations/${evaluationId}/attempts`);
	}

	/**
	 * Check if user can attempt an evaluation
	 */
	async canUserAttempt(userId: number, evaluationId: number): Promise<{ can_attempt: boolean; reason: string }> {
		return this.get<{ can_attempt: boolean; reason: string }>(`/api/v1/users/${userId}/evaluations/${evaluationId}/can-attempt`);
	}

	/**
	 * Score an evaluation attempt (admin only)
	 */
	async scoreAttempt(attemptId: number): Promise<EvaluationAttempt> {
		return this.post<EvaluationAttempt>(`/api/v1/evaluation-attempts/${attemptId}/score`);
	}

	/**
	 * Helper method to get the latest attempt for a user and evaluation
	 */
	async getLatestAttempt(userId: number, evaluationId: number): Promise<EvaluationAttempt | null> {
		const attempts = await this.getUserAttempts(userId, evaluationId);
		if (attempts.length === 0) return null;
		
		// Sort by started_at descending and return the most recent
		return attempts.sort((a, b) => 
			new Date(b.started_at).getTime() - new Date(a.started_at).getTime()
		)[0];
	}

	/**
	 * Helper method to check if user has passed an evaluation
	 */
	async hasUserPassed(userId: number, evaluationId: number): Promise<boolean> {
		const attempts = await this.getUserAttempts(userId, evaluationId);
		return attempts.some(attempt => attempt.passed && attempt.submitted_at);
	}

	/**
	 * Helper method to get user's best score for an evaluation
	 */
	async getBestScore(userId: number, evaluationId: number): Promise<{ score: number; total_points: number } | null> {
		const attempts = await this.getUserAttempts(userId, evaluationId);
		const completedAttempts = attempts.filter(attempt => attempt.submitted_at);
		
		if (completedAttempts.length === 0) return null;
		
		const bestAttempt = completedAttempts.reduce((best, current) => 
			current.score > best.score ? current : best
		);
		
		return {
			score: bestAttempt.score,
			total_points: bestAttempt.total_points
		};
	}
}
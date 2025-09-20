import { apiClient } from '$lib/client';
import type {
	Evaluation,
	CreateEvaluationDTO,
	UpdateEvaluationDTO
} from '$lib/types';

export class EvaluationController {
	/**
	 * Get a specific evaluation by ID
	 */
	async getEvaluation(id: number): Promise<Evaluation> {
		return apiClient.get<Evaluation>(`/api/v1/evaluations/${id}`);
	}

	/**
	 * Get all evaluations for a specific module
	 */
	async getEvaluationsByModule(moduleId: number): Promise<Evaluation[]> {
		return apiClient.get<Evaluation[]>(`/api/v1/modules/${moduleId}/evaluations`);
	}

	/**
	 * Create new evaluation
	 */
	async createEvaluation(evaluationData: CreateEvaluationDTO): Promise<Evaluation> {
		return apiClient.post<Evaluation>('/api/v1/evaluations', evaluationData);
	}

	/**
	 * Update existing evaluation
	 */
	async updateEvaluation(id: number, evaluationData: UpdateEvaluationDTO): Promise<Evaluation> {
		return apiClient.put<Evaluation>(`/api/v1/evaluations/${id}`, evaluationData);
	}

	/**
	 * Delete evaluation
	 */
	async deleteEvaluation(id: number): Promise<void> {
		return apiClient.delete(`/api/v1/evaluations/${id}`);
	}
}
import { BaseController } from './base';
import type { Evaluation, CreateEvaluationDTO, UpdateEvaluationDTO } from '$lib/types';

export class EvaluationController extends BaseController {
	/**
	 * Get a specific evaluation by ID
	 */
	async getEvaluation(id: number): Promise<Evaluation> {
		return this.get<Evaluation>(`/api/v1/evaluations/${id}`);
	}

	/**
	 * Get all evaluations for a specific module
	 */
	async getEvaluationsByModule(moduleId: number): Promise<Evaluation[]> {
		return this.get<Evaluation[]>(`/api/v1/modules/${moduleId}/evaluations`);
	}

	/**
	 * Create new evaluation
	 */
	async createEvaluation(evaluationData: CreateEvaluationDTO): Promise<Evaluation> {
		return this.post<Evaluation>('/api/v1/evaluations', evaluationData);
	}

	/**
	 * Update existing evaluation
	 */
	async updateEvaluation(id: number, evaluationData: UpdateEvaluationDTO): Promise<Evaluation> {
		return this.put<Evaluation>(`/api/v1/evaluations/${id}`, evaluationData);
	}

	/**
	 * Partially update existing evaluation (PATCH)
	 */
	async updateEvaluationPatch(
		id: number,
		evaluationData: Partial<UpdateEvaluationDTO>
	): Promise<Evaluation> {
		return this.patch<Evaluation>(`/api/v1/evaluations/${id}`, evaluationData);
	}

	/**
	 * Delete evaluation
	 */
	async deleteEvaluation(id: number): Promise<void> {
		return this.delete(`/api/v1/evaluations/${id}`);
	}
}

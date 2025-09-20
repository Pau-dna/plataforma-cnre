import { apiClient } from '$lib/client';
import type {
	Module,
	CreateModuleDTO,
	UpdateModuleDTO,
	Content,
	Evaluation,
	ReorderItemDTO
} from '$lib/types';

export class ModuleController {
	/**
	 * Get a specific module by ID
	 */
	async getModule(id: number): Promise<Module> {
		return apiClient.get<Module>(`/api/v1/modules/${id}`);
	}

	/**
	 * Get module with all content loaded
	 */
	async getModuleWithContent(id: number): Promise<Module> {
		return apiClient.get<Module>(`/api/v1/modules/${id}/content`);
	}

	/**
	 * Get all modules for a specific course
	 */
	async getModulesByCourse(courseId: number): Promise<Module[]> {
		return apiClient.get<Module[]>(`/api/v1/courses/${courseId}/modules`);
	}

	/**
	 * Create a new module
	 */
	async createModule(moduleData: CreateModuleDTO): Promise<Module> {
		return apiClient.post<Module>('/api/v1/modules', moduleData);
	}

	/**
	 * Update an existing module
	 */
	async updateModule(id: number, moduleData: UpdateModuleDTO): Promise<Module> {
		return apiClient.put<Module>(`/api/v1/modules/${id}`, moduleData);
	}

	/**
	 * Delete a module
	 */
	async deleteModule(id: number): Promise<void> {
		return apiClient.delete(`/api/v1/modules/${id}`);
	}

	/**
	 * Reorder modules within a course
	 */
	async reorderModules(courseId: number, moduleOrders: ReorderItemDTO[]): Promise<void> {
		return apiClient.post(`/api/v1/courses/${courseId}/modules/reorder`, moduleOrders);
	}

	/**
	 * Get all content for a specific module
	 */
	async getModuleContents(moduleId: number): Promise<Content[]> {
		return apiClient.get<Content[]>(`/api/v1/modules/${moduleId}/content`);
	}

	/**
	 * Get all evaluations for a specific module
	 */
	async getModuleEvaluations(moduleId: number): Promise<Evaluation[]> {
		return apiClient.get<Evaluation[]>(`/api/v1/modules/${moduleId}/evaluations`);
	}

	/**
	 * Legacy method for backwards compatibility
	 * @deprecated Use getModulesByCourse instead
	 */
	async getModules(courseID: number): Promise<Module[]> {
		return this.getModulesByCourse(courseID);
	}
}

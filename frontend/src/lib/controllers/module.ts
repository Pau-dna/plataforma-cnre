import { BaseController } from './base';
import type {
	Module,
	CreateModuleDTO,
	UpdateModuleDTO,
	Content,
	Evaluation,
	ReorderItemDTO
} from '$lib/types';

export class ModuleController extends BaseController {
	/**
	 * Get a specific module by ID
	 */
	async getModule(id: number): Promise<Module> {
		return this.get<Module>(`/api/v1/modules/${id}`);
	}

	/**
	 * Get module with all content loaded
	 */
	async getModuleWithContent(id: number): Promise<Module> {
		return this.get<Module>(`/api/v1/modules/${id}/content`);
	}

	/**
	 * Get all modules for a specific course
	 */
	async getModulesByCourse(courseId: number): Promise<Module[]> {
		return this.get<Module[]>(`/api/v1/courses/${courseId}/modules`);
	}

	/**
	 * Create a new module
	 */
	async createModule(moduleData: CreateModuleDTO): Promise<Module> {
		return this.post<Module>('/api/v1/modules', moduleData);
	}

	/**
	 * Update an existing module
	 */
	async updateModule(id: number, moduleData: UpdateModuleDTO): Promise<Module> {
		return this.put<Module>(`/api/v1/modules/${id}`, moduleData);
	}

	/**
	 * Partially update an existing module (PATCH)
	 */
	async updateModulePatch(id: number, moduleData: Partial<UpdateModuleDTO>): Promise<Module> {
		return this.patch<Module>(`/api/v1/modules/${id}`, moduleData);
	}

	/**
	 * Delete a module
	 */
	async deleteModule(id: number): Promise<void> {
		return this.delete(`/api/v1/modules/${id}`);
	}

	/**
	 * Reorder modules within a course
	 */
	async reorderModules(courseId: number, moduleOrders: ReorderItemDTO[]): Promise<void> {
		return this.post(`/api/v1/courses/${courseId}/modules/reorder`, moduleOrders);
	}

	/**
	 * Get all content for a specific module
	 */
	async getModuleContents(moduleId: number): Promise<Content[]> {
		return this.get<Content[]>(`/api/v1/modules/${moduleId}/content`);
	}

	/**
	 * Get all evaluations for a specific module
	 */
	async getModuleEvaluations(moduleId: number): Promise<Evaluation[]> {
		return this.get<Evaluation[]>(`/api/v1/modules/${moduleId}/evaluations`);
	}

	/**
	 * Legacy method for backwards compatibility
	 * @deprecated Use getModulesByCourse instead
	 */
	async getModules(courseID: number): Promise<Module[]> {
		return this.getModulesByCourse(courseID);
	}
}

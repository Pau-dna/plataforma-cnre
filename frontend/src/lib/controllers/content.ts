import { apiClient } from '$lib/client';
import type {
	Content,
	CreateContentDTO,
	UpdateContentDTO,
	ReorderItemDTO,
	ModuleContent
} from '$lib/types';

export class ContentController {
	/**
	 * Get a specific content by ID
	 */
	async getContent(id: number): Promise<Content> {
		return apiClient.get<Content>(`/api/v1/content/${id}`);
	}

	/**
	 * Get all content for a specific module
	 */
	async getContentsByModule(moduleId: number): Promise<Content[]> {
		return apiClient.get<Content[]>(`/api/v1/modules/${moduleId}/content`);
	}

	/**
	 * Create new content
	 */
	async createContent(contentData: CreateContentDTO): Promise<Content> {
		return apiClient.post<Content>('/api/v1/content', contentData);
	}

	/**
	 * Update existing content
	 */
	async updateContent(id: number, contentData: UpdateContentDTO): Promise<Content> {
		return apiClient.put<Content>(`/api/v1/content/${id}`, contentData);
	}

	/**
	 * Delete content
	 */
	async deleteContent(id: number): Promise<void> {
		return apiClient.delete(`/api/v1/content/${id}`);
	}

	/**
	 * Reorder content within a module
	 */
	async reorderContent(moduleId: number, contentOrders: ReorderItemDTO[]): Promise<void> {
		return apiClient.post(`/api/v1/modules/${moduleId}/content/reorder`, contentOrders);
	}

	/**
	 * Legacy method for backwards compatibility
	 * @deprecated Use getContentsByModule instead
	 */
	async getContents(moduleID: number): Promise<ModuleContent[]> {
		// This method previously returned mock data including evaluations
		// For now, we'll just return contents and handle evaluations separately
		return this.getContentsByModule(moduleID);
	}
}

import { BaseController } from './base';
import type {
	Content,
	CreateContentDTO,
	UpdateContentDTO,
	ReorderItemDTO,
	ModuleContent
} from '$lib/types';

export class ContentController extends BaseController {
	/**
	 * Get a specific content by ID
	 */
	async getContent(id: number): Promise<Content> {
		return this.get<Content>(`/api/v1/content/${id}`);
	}

	/**
	 * Get all content for a specific module
	 */
	async getContentsByModule(moduleId: number): Promise<Content[]> {
		return this.get<Content[]>(`/api/v1/modules/${moduleId}/content`);
	}

	/**
	 * Create new content
	 */
	async createContent(contentData: CreateContentDTO): Promise<Content> {
		return this.post<Content>('/api/v1/content', contentData);
	}

	/**
	 * Update existing content
	 */
	async updateContent(id: number, contentData: UpdateContentDTO): Promise<Content> {
		return this.put<Content>(`/api/v1/content/${id}`, contentData);
	}

	/**
	 * Partially update existing content (PATCH)
	 */
	async updateContentPatch(id: number, contentData: Partial<UpdateContentDTO>): Promise<Content> {
		return this.patch<Content>(`/api/v1/content/${id}`, contentData);
	}

	/**
	 * Delete content
	 */
	async deleteContent(id: number): Promise<void> {
		return this.delete(`/api/v1/content/${id}`);
	}

	/**
	 * Reorder content within a module
	 */
	async reorderContent(moduleId: number, contentOrders: ReorderItemDTO[]): Promise<void> {
		for (const content of contentOrders) {
			this.updateContentPatch(content.id, {
				order: content.order
			});
		}
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

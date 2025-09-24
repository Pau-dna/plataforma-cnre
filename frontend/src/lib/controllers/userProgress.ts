import { BaseController } from './base';
import type { UserProgress } from '$lib/types';

export class UserProgressController extends BaseController {
	/**
	 * Get user progress for a specific course
	 * WARNING: This endpoint doesn't exist in backend
	 */
	async getUserCourseProgress(userId: number, courseId: number): Promise<UserProgress[]> {
		throw new Error('User progress endpoints not implemented in backend');
	}

	/**
	 * Get user progress for a specific module
	 * WARNING: This endpoint doesn't exist in backend
	 */
	async getUserModuleProgress(userId: number, moduleId: number): Promise<UserProgress[]> {
		throw new Error('User progress endpoints not implemented in backend');
	}

	/**
	 * Mark content as completed
	 * WARNING: This endpoint doesn't exist in backend
	 */
	async markContentComplete(
		userId: number,
		courseId: number,
		moduleId: number,
		contentId: number
	): Promise<UserProgress> {
		throw new Error('User progress endpoints not implemented in backend');
	}

	/**
	 * Mark content as incomplete
	 * WARNING: This endpoint doesn't exist in backend
	 */
	async markContentIncomplete(
		userId: number,
		courseId: number,
		moduleId: number,
		contentId: number
	): Promise<void> {
		throw new Error('User progress endpoints not implemented in backend');
	}

	/**
	 * Calculate course progress percentage
	 * WARNING: This endpoint doesn't exist in backend
	 */
	async calculateCourseProgress(userId: number, courseId: number): Promise<number> {
		throw new Error('User progress endpoints not implemented in backend');
	}

	/**
	 * Calculate module progress percentage
	 * WARNING: This endpoint doesn't exist in backend
	 */
	async calculateModuleProgress(userId: number, moduleId: number): Promise<number> {
		throw new Error('User progress endpoints not implemented in backend');
	}

	/**
	 * Helper method to check if content is completed by user
	 * WARNING: This endpoint doesn't exist in backend
	 */
	async isContentCompleted(userId: number, contentId: number): Promise<boolean> {
		throw new Error('User progress endpoints not implemented in backend');
	}

	/**
	 * Helper method to get completion percentage for a list of content items
	 */
	async getContentCompletionRate(userId: number, contentIds: number[]): Promise<number> {
		if (contentIds.length === 0) return 0;

		const completionChecks = await Promise.all(
			contentIds.map(async (contentId) => {
				try {
					return await this.isContentCompleted(userId, contentId);
				} catch {
					return false;
				}
			})
		);

		const completedCount = completionChecks.filter(Boolean).length;
		return (completedCount / contentIds.length) * 100;
	}

	/**
	 * Update course enrollment progress (WORKING ENDPOINT)
	 * This is the only progress endpoint that exists in backend
	 */
	async updateCourseProgress(userId: number, courseId: number, progress: number): Promise<void> {
		return this.put(`/api/v1/users/${userId}/courses/${courseId}/progress`, { progress });
	}
}

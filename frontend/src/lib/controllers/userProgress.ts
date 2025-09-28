import { BaseController } from './base';
import type { UserProgress } from '$lib/types';

export class UserProgressController extends BaseController {
	/**
	 * Get user progress for a specific course
	 */
	async getUserCourseProgress(userId: number, courseId: number): Promise<UserProgress[]> {
		return this.get<UserProgress[]>(`/api/v1/users/${userId}/courses/${courseId}/progress`);
	}

	/**
	 * Get user progress for a specific module
	 */
	async getUserModuleProgress(userId: number, moduleId: number): Promise<UserProgress[]> {
		return this.get<UserProgress[]>(`/api/v1/users/${userId}/modules/${moduleId}/progress`);
	}

	/**
	 * Mark content as completed
	 */
	async markContentComplete(
		userId: number,
		courseId: number,
		moduleId: number,
		contentId: number
	): Promise<UserProgress> {
		return this.post<UserProgress>('/api/v1/user-progress/complete', {
			user_id: userId,
			course_id: courseId,
			module_id: moduleId,
			content_id: contentId
		});
	}

	/**
	 * Mark content as incomplete
	 */
	async markContentIncomplete(
		userId: number,
		courseId: number,
		moduleId: number,
		contentId: number
	): Promise<void> {
		return this.post('/api/v1/user-progress/incomplete', {
			user_id: userId,
			course_id: courseId,
			module_id: moduleId,
			content_id: contentId
		});
	}

	/**
	 * Calculate course progress percentage
	 */
	async calculateCourseProgress(userId: number, courseId: number): Promise<number> {
		return this.get<number>(`/api/v1/users/${userId}/courses/${courseId}/progress-percentage`);
	}

	/**
	 * Calculate module progress percentage
	 */
	async calculateModuleProgress(userId: number, moduleId: number): Promise<number> {
		return this.get<number>(`/api/v1/users/${userId}/modules/${moduleId}/progress-percentage`);
	}

	/**
	 * Helper method to check if content is completed by user
	 */
	async isContentCompleted(userId: number, contentId: number): Promise<boolean> {
		try {
			const progress = await this.get<UserProgress>(
				`/api/v1/users/${userId}/content/${contentId}/progress`
			);
			return !!progress.completed_at;
		} catch (error) {
			return false;
		}
	}

	/**
	 * Helper method to get completion percentage for a list of content items
	 */
	async getContentCompletionRate(userId: number, contentIds: number[]): Promise<number> {
		if (contentIds.length === 0) return 0;

		const completionChecks = await Promise.all(
			contentIds.map((contentId) => this.isContentCompleted(userId, contentId))
		);

		const completedCount = completionChecks.filter(Boolean).length;
		return (completedCount / contentIds.length) * 100;
	}

	/**
	 * Helper method to check if user has passed an evaluation
	 */
	async hasUserPassedEvaluation(userId: number, evaluationId: number): Promise<boolean> {
		try {
			const result = await this.get<{ has_passed: boolean }>(
				`/api/v1/users/${userId}/evaluations/${evaluationId}/passed`
			);
			return result.has_passed;
		} catch (error) {
			return false;
		}
	}
}

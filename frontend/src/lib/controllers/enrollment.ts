import { BaseController } from './base';
import type {
	Enrollment,
	CreateEnrollmentDTO,
	UpdateEnrollmentProgressDTO
} from '$lib/types';

export class EnrollmentController extends BaseController {
	/**
	 * Get a specific enrollment by ID
	 */
	async getEnrollment(id: number): Promise<Enrollment> {
		return this.get<Enrollment>(`/api/v1/enrollments/${id}`);
	}

	/**
	 * Create new enrollment (enroll user in course)
	 */
	async createEnrollment(enrollmentData: CreateEnrollmentDTO): Promise<Enrollment> {
		return this.post<Enrollment>('/api/v1/enrollments', enrollmentData);
	}

	/**
	 * Delete enrollment (unenroll user from course)
	 */
	async deleteEnrollment(id: number): Promise<void> {
		return this.delete(`/api/v1/enrollments/${id}`);
	}

	/**
	 * Get all enrollments for a specific user
	 */
	async getUserEnrollments(userId: number): Promise<Enrollment[]> {
		return this.get<Enrollment[]>(`/api/v1/users/${userId}/enrollments`);
	}

	/**
	 * Get all enrollments for a specific course
	 */
	async getCourseEnrollments(courseId: number): Promise<Enrollment[]> {
		return this.get<Enrollment[]>(`/api/v1/courses/${courseId}/enrollments`);
	}

	/**
	 * Complete enrollment (mark course as completed)
	 */
	async completeEnrollment(userId: number, courseId: number): Promise<void> {
		return this.post(`/api/v1/users/${userId}/courses/${courseId}/complete`);
	}

	/**
	 * Update enrollment progress
	 */
	async updateProgress(userId: number, courseId: number, progressData: UpdateEnrollmentProgressDTO): Promise<void> {
		return this.put(`/api/v1/users/${userId}/courses/${courseId}/progress`, progressData);
	}

	/**
	 * Helper method to enroll current user in a course
	 */
	async enrollInCourse(courseId: number, userId: number): Promise<Enrollment> {
		return this.createEnrollment({ user_id: userId, course_id: courseId });
	}

	/**
	 * Helper method to check if user is enrolled in a course
	 */
	async isUserEnrolled(userId: number, courseId: number): Promise<boolean> {
		try {
			const enrollments = await this.getUserEnrollments(userId);
			return enrollments.some(enrollment => enrollment.course_id === courseId);
		} catch (error) {
			return false;
		}
	}
}
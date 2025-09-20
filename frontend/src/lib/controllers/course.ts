import { apiClient } from '$lib/client';
import type {
	Course,
	CreateCourseDTO,
	UpdateCourseDTO,
	Module,
	Enrollment
} from '$lib/types';

export class CourseController {
	/**
	 * Get all courses
	 */
	async getCourses(): Promise<Course[]> {
		return apiClient.get<Course[]>('/api/v1/courses');
	}

	/**
	 * Get a specific course by ID
	 */
	async getCourse(id: number): Promise<Course> {
		return apiClient.get<Course>(`/api/v1/courses/${id}`);
	}

	/**
	 * Get course with all modules loaded
	 */
	async getCourseWithModules(id: number): Promise<Course> {
		return apiClient.get<Course>(`/api/v1/courses/${id}/modules`);
	}

	/**
	 * Create a new course
	 */
	async createCourse(courseData: CreateCourseDTO): Promise<Course> {
		return apiClient.post<Course>('/api/v1/courses', courseData);
	}

	/**
	 * Update an existing course
	 */
	async updateCourse(id: number, courseData: UpdateCourseDTO): Promise<Course> {
		return apiClient.put<Course>(`/api/v1/courses/${id}`, courseData);
	}

	/**
	 * Delete a course
	 */
	async deleteCourse(id: number): Promise<void> {
		return apiClient.delete(`/api/v1/courses/${id}`);
	}

	/**
	 * Get all modules for a specific course
	 */
	async getCourseModules(courseId: number): Promise<Module[]> {
		return apiClient.get<Module[]>(`/api/v1/courses/${courseId}/modules`);
	}

	/**
	 * Get all enrollments for a specific course
	 */
	async getCourseEnrollments(courseId: number): Promise<Enrollment[]> {
		return apiClient.get<Enrollment[]>(`/api/v1/courses/${courseId}/enrollments`);
	}
}

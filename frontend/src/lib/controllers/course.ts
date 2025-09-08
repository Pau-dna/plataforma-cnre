import type { Course } from '$lib/types/models/course';

export class CourseController {
	async getCourses(): Promise<Course[]> {
		return [
			{
				id: 1,
				title: 'Course 1',
				description: 'Description for Course 1',
				student_count: 120,
				module_count: 10
			}
		];
	}
}

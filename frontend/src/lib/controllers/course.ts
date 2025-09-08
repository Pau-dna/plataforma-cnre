import type { Course } from '$lib/types/models/course';

export class CourseController {
	async getCourses(): Promise<Course[]> {
		return [
			{
				id: 1,
				title: 'Inducción al CNRE 2025',
				description: 'Curso introductorio al CNRE para estudiantes de la unal.',
				student_count: 120,
				module_count: 10,
				modules: [],
			}
		];
	}
}

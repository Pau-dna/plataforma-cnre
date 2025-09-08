import type { Course } from '$lib/types/models/course';

export class CourseController {
	async getCourses(): Promise<Course[]> {
		return [
			{
				id: 1,
				title: 'Capacitación CNRE Medellín 2025',
				description: 'Curso de formación para el personal del CNRE en Medellín, enfocado en competencias técnicas y administrativas para el año 2025.',
				student_count: 80,
				module_count: 6
			}
		];
	}
}

import type { Module } from '$lib/types/models/course';

export class ModuleController {
	async getModules(courseID: number): Promise<Module[]> {
		return [
			{
				id: 1,
				order: 1,
				course_id: courseID,
				contents: [],
				title: 'Module 1',
				description: 'Description for Module 1'
			},
			{
				id: 2,
				order: 2,
				course_id: courseID,
				contents: [],
				title: 'Module 2',
				description: 'Description for Module 2'
			}
		];
	}
}

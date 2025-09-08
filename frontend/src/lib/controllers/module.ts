import type { Module } from '$lib/types/models/course';

export class ModuleController {
	async getModules(courseID: number): Promise<Module[]> {
		return [
			{
				id: 1,
				order: 1,
				course_id: courseID,
				contents: [],
				title: 'Introducción al CNRE',
				description: 'Bienvenida y presentación general del CNRE Medellín 2025.'
			},
			{
				id: 2,
				order: 2,
				course_id: courseID,
				contents: [],
				title: 'Competencias Técnicas',
				description: 'Formación en habilidades técnicas requeridas para el personal del CNRE.'
			},
			{
				id: 3,
				order: 3,
				course_id: courseID,
				contents: [],
				title: 'Gestión Administrativa',
				description: 'Módulo sobre procesos administrativos y gestión interna.'
			},
			{
				id: 4,
				order: 4,
				course_id: courseID,
				contents: [],
				title: 'Evaluación y Seguimiento',
				description: 'Herramientas y metodologías para la evaluación y seguimiento del desempeño.'
			}
		];
	}
}

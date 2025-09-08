import type { ModuleContent } from "$lib/types/models/course";
import { ContentType, QuestionType } from "$lib/types/models/course";

export class ContentController {
	async getContents(moduleID: number): Promise<ModuleContent[]> {
		return [
			{
				id: 1,
				order: 1,
				title: "Bienvenida",
				description: "Introducción general al módulo y objetivos.",
				type: ContentType.CONTENT,
				body: "Este contenido da la bienvenida a los participantes y explica el propósito del módulo.",
				mediaUrl: "https://cdn.example.com/videos/bienvenida.mp4",
				created_at: "2025-01-10T08:00:00Z",
				updated_at: "2025-01-15T12:00:00Z",
			},
			{
				id: 2,
				order: 2,
				title: "Lectura de Referencia",
				description: "Documento de apoyo para reforzar conceptos clave.",
				type: ContentType.CONTENT,
				body: "Consulta este documento PDF con información detallada.",
				mediaUrl: "https://cdn.example.com/docs/referencia.pdf",
				created_at: "2025-01-11T08:00:00Z",
				updated_at: "2025-01-16T12:00:00Z",
			},
			{
				id: 3,
				order: 3,
				title: "Cuestionario de Diagnóstico",
				description: "Evalúa tus conocimientos iniciales antes de continuar.",
				type: ContentType.EVALUATION,
				questions: [
					{
						id: 101,
						text: "¿Cuál es el objetivo principal del CNRE?",
						type: QuestionType.SINGLE,
						points: 5,
						answers: [
							{ id: 201, text: "Promover actividades culturales", isCorrect: false, order: 1 },
							{ id: 202, text: "Formar y capacitar al personal", isCorrect: true, order: 2 },
							{ id: 203, text: "Administrar los recursos financieros", isCorrect: false, order: 3 },
						],
						explanation: "El CNRE busca principalmente la formación y capacitación.",
						created_at: "2025-01-12T10:00:00Z",
					},
					{
						id: 102,
						text: "Selecciona los elementos que forman parte de la inducción.",
						type: QuestionType.MULTIPLE,
						points: 10,
						answers: [
							{ id: 204, text: "Bienvenida institucional", isCorrect: true, order: 1 },
							{ id: 205, text: "Prueba de laboratorio", isCorrect: false, order: 2 },
							{ id: 206, text: "Presentación de normas", isCorrect: true, order: 3 },
						],
						explanation: "La inducción incluye bienvenida y normas generales.",
						created_at: "2025-01-12T11:00:00Z",
					},
				],
				question_count: 2,
				passing_score: 10,
				max_attempts: 3,
				time_limit: 20,
				created_at: "2025-01-12T08:00:00Z",
				updated_at: "2025-01-17T12:00:00Z",
			},
		];
	}
}

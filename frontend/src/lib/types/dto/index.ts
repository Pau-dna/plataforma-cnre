// DTOs for API requests and responses
import { ContentType, QuestionType } from '../models/course';

export interface CreateCourseDTO {
	title: string;
	description: string;
	short_description?: string;
	image_url?: string;
}

export interface UpdateCourseDTO {
	title?: string;
	description?: string;
	short_description?: string;
	image_url?: string;
}

export interface CreateModuleDTO {
	title: string;
	description?: string;
	course_id: number;
}

export interface UpdateModuleDTO {
	title?: string;
	description?: string;
	order?: number;
}

export interface CreateContentDTO {
	title: string;
	description?: string;
	type: ContentType.CONTENT;
	body: string;
	media_url?: string;
	module_id: number;
}

export interface UpdateContentDTO {
	order?: number;
	title?: string;
	description?: string;
	body?: string;
	media_url?: string;
}

export interface CreateEvaluationDTO {
	order: number;
	title: string;
	description?: string;
	type: ContentType.EVALUATION;
	question_count: number;
	passing_score: number;
	max_attempts?: number;
	time_limit?: number;
	module_id: number;
}

export interface UpdateEvaluationDTO {
	order?: number;
	title?: string;
	description?: string;
	question_count?: number;
	passing_score?: number;
	max_attempts?: number;
	time_limit?: number;
}

export interface CreateQuestionDTO {
	text: string;
	type: QuestionType;
	explanation?: string;
	points: number;
	evaluation_id: number;
}

export interface UpdateQuestionDTO {
	text?: string;
	type?: QuestionType;
	explanation?: string;
	points?: number;
}

export interface CreateAnswerDTO {
	text: string;
	is_correct: boolean;
	order: number;
	question_id: number;
}

export interface UpdateAnswerDTO {
	text?: string;
	is_correct?: boolean;
	order?: number;
}

export interface CreateEnrollmentDTO {
	user_id: number;
	course_id: number;
}

export interface UpdateEnrollmentProgressDTO {
	progress: number;
}

export interface StartEvaluationAttemptDTO {
	user_id: number;
	evaluation_id: number;
}

export interface SubmitEvaluationAttemptDTO {
	answers: {
		question_id: number;
		selected_answer_ids: number[];
		is_correct: boolean;
		points: number;
	}[];
}

export interface ReorderItemDTO {
	id: number;
	order: number;
}

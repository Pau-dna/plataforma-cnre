// Base types and enums
export enum QuestionType {
	SINGLE = 'single_choice',
	MULTIPLE = 'multiple_choice'
}

export enum ContentType {
	CONTENT = 'content',
	EVALUATION = 'evaluation'
}

export enum ModuleStatus {
	LOCKED = 'locked',
	AVAILABLE = 'available',
	COMPLETED = 'completed'
}

export enum UserRole {
	STUDENT = 'student',
	INSTRUCTOR = 'instructor',
	ADMIN = 'admin'
}

// Base interface for entities with common properties
export interface BaseEntity {
	id: number;
	created_at?: string;
	updated_at?: string;
}

// Course interface - matches backend Course model
export interface Course extends BaseEntity {
	title: string;
	description: string;
	short_description?: string;
	image_url?: string;
	student_count?: number;
	module_count?: number;
	modules?: Module[];
	enrollments?: Enrollment[];
}

// Module interface - matches backend Module model
export interface Module extends BaseEntity {
	title: string;
	description?: string;
	order: number;
	course_id: number;
	course?: Course;
	contents?: Content[];
	evaluations?: Evaluation[];
}

// Content interface - matches backend Content model
export interface Content extends BaseEntity {
	order: number;
	title: string;
	description?: string;
	type: ContentType.CONTENT;
	body: string;
	media_url?: string;
	module_id: number;
	module?: Module;
	user_progress?: UserProgress[];
}

// Evaluation interface - matches backend Evaluation model
export interface Evaluation extends BaseEntity {
	order: number;
	title: string;
	description?: string;
	type: ContentType.EVALUATION;
	question_count: number;
	passing_score: number;
	max_attempts?: number;
	time_limit?: number;
	module_id: number;
	module?: Module;
	questions?: Question[];
	evaluation_attempts?: EvaluationAttempt[];
	user_progress?: UserProgress[];
}

// Question interface - matches backend Question model
export interface Question extends BaseEntity {
	text: string;
	type: QuestionType;
	explanation?: string;
	points: number;
	evaluation_id: number;
	evaluation?: Evaluation;
	answers?: Answer[];
}

// Answer interface - matches backend Answer model
export interface Answer extends BaseEntity {
	text: string;
	is_correct: boolean;
	order: number;
	question_id: number;
	question?: Question;
}

// Enrollment interface - matches backend Enrollment model
export interface Enrollment extends BaseEntity {
	user_id: number;
	course_id: number;
	enrolled_at: string;
	completed_at?: string;
	progress: number;
	user?: User;
	course?: Course;
}

// EvaluationAttempt interface - matches backend EvaluationAttempt model
export interface EvaluationAttempt extends BaseEntity {
	user_id: number;
	evaluation_id: number;
	answers: {
		question_id: number;
		selected_answer_ids: number[];
		is_correct: boolean;
		points: number;
	}[];
	score: number;
	total_points: number;
	passed: boolean;
	started_at: string;
	submitted_at?: string;
	time_spent?: number;
	user?: User;
	evaluation?: Evaluation;
}

// UserProgress interface - matches backend UserProgress model
export interface UserProgress extends BaseEntity {
	user_id: number;
	course_id: number;
	module_id: number;
	content_id: number;
	completed_at?: string;
	score?: number;
	attempts: number;
	user?: User;
	course?: Course;
	module?: Module;
	content?: Content;
}

// User interface - matches backend User model
export interface User extends BaseEntity {
	email: string;
	first_name: string;
	last_name: string;
	avatar_url?: string;
	role: UserRole;
	enrollments?: Enrollment[];
}

// Union types
export type ModuleContent = Content | Evaluation;
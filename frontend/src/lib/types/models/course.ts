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

// Base interface for entities with common properties
interface BaseEntity {
	id: number;
	created_at: string;
	updated_at: string;
}

// Answer interface for question options
export interface Answer extends BaseEntity {
	text: string;
	isCorrect: boolean;
	order: number;
}

// Question interface
export interface Question extends BaseEntity {
	text: string;
	type: QuestionType;
	answers: Answer[];
	explanation?: string;
	points: number;
}

// Base content interface
interface BaseContent extends BaseEntity {
	title: string;
	description?: string;
	type: ContentType;
	isRequired: boolean;
}

// Content (lessons, videos, readings, etc.)
export interface Content extends BaseContent {
	type: ContentType.CONTENT;
	body: string;
	mediaUrl?: string;
	resources?: string[];
}

// Evaluation (quizzes, exams, etc.)
export interface Evaluation extends BaseContent {
	type: ContentType.EVALUATION;
	questions: Question[];
	passingScore: number;
	maxAttempts?: number;
	timeLimit?: number; // in minutes
	shuffleQuestions: boolean;
	shuffleAnswers: boolean;
}

// Union type for module content
export type ModuleContent = Content | Evaluation;


// Course interface
export interface Course extends BaseEntity {
	title: string;
	description: string;
	shortDescription?: string;
	imageUrl?: string;
	modules: Module[];
	student_count?: number;
	module_count?: number;
}

// Module interface
export interface Module extends BaseEntity {
	title: string;
	description?: string;
	order: number;
	contents: ModuleContent[];
	course_id: number;
}

// Course enrollment
export interface Enrollment extends BaseEntity {
	user_id: number;
	course_id: number;
	enrolled_at: string;
	completed_at?: string;
	progress: number;
}

// Evaluation attempt
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
	started_at: Date;
	submitted_at?: Date;
	time_spent?: number; // in minutes
}

// User progress tracking
export interface UserProgress extends BaseEntity {
	user_id: number;
	course_id: number;
	module_id: number;
	content_id: number;
	completed_at?: string;
	score?: number;
	attempts: number;
}

// User interface
export interface User extends BaseEntity {
	email: string;
	firstName: string;
	lastName: string;
	avatarUrl?: string;
	role: 'student' | 'instructor' | 'admin';
}

// Instructor interface (extends User)
export interface Instructor extends User {
	role: 'instructor';
	bio?: string;
	specializations: string[];
	socialLinks?: {
		linkedin?: string;
		twitter?: string;
		website?: string;
	};
}

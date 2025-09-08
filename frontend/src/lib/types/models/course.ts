// Base types and enums
export enum QuestionType {
    SINGLE_CHOICE = 'single_choice',
    MULTIPLE_CHOICE = 'multiple_choice'
}

export enum ContentType {
    CONTENT = 'content',
    EVALUATION = 'evaluation'
}

export enum CourseStatus {
    DRAFT = 'draft',
    PUBLISHED = 'published',
    ARCHIVED = 'archived'
}

export enum ModuleStatus {
    LOCKED = 'locked',
    AVAILABLE = 'available',
    COMPLETED = 'completed'
}

// Base interface for entities with common properties
interface BaseEntity {
    id: string;
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
    order: number;
}

// Base content interface
interface BaseContent extends BaseEntity {
    title: string;
    description?: string;
    type: ContentType;
    order: number;
    isRequired: boolean;
    estimatedDuration?: number; // in minutes
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

// Module interface
export interface Module extends BaseEntity {
    title: string;
    description?: string;
    order: number;
    status: ModuleStatus;
    contents: ModuleContent[];
    prerequisites?: string[]; // Module IDs
    estimatedDuration?: number; // in minutes
    courseId: string;
}

// Course interface
export interface Course extends BaseEntity {
    title: string;
    description: string;
    shortDescription?: string;
    imageUrl?: string;
    modules: Module[];
    instructorId: string;
    student_count?: number;
    module_count?: number;
}

// User progress tracking
export interface UserProgress extends BaseEntity {
    userId: string;
    courseId: string;
    moduleId: string;
    contentId: string;
    completedAt?: Date;
    score?: number;
    attempts: number;
    timeSpent?: number; // in minutes
}

// Course enrollment
export interface Enrollment extends BaseEntity {
    userId: string;
    courseId: string;
    enrolledAt: Date;
    completedAt?: Date;
    progress: number; // percentage 0-100
    lastAccessedAt?: Date;
}

// Evaluation attempt
export interface EvaluationAttempt extends BaseEntity {
    userId: string;
    evaluationId: string;
    answers: {
        questionId: string;
        selectedAnswerIds: string[];
        isCorrect: boolean;
        points: number;
    }[];
    score: number;
    totalPoints: number;
    passed: boolean;
    startedAt: Date;
    submittedAt?: Date;
    timeSpent?: number; // in minutes
}

// User interface
export interface User extends BaseEntity {
    email: string;
    firstName: string;
    lastName: string;
    avatarUrl?: string;
    role: 'student' | 'instructor' | 'admin';
    isActive: boolean;
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

// API Response types
export interface PaginatedResponse<T> {
    data: T[];
    total: number;
    page: number;
    limit: number;
    totalPages: number;
}

export interface ApiResponse<T> {
    success: boolean;
    data?: T;
    error?: string;
    message?: string;
}

// DTOs (Data Transfer Objects) for API operations
export interface CreateCourseDto {
    title: string;
    description: string;
    shortDescription?: string;
    imageUrl?: string;
    instructorId: string;
    price?: number;
    currency?: string;
    tags: string[];
    difficulty: 'beginner' | 'intermediate' | 'advanced';
    language: string;
}

export interface UpdateCourseDto extends Partial<CreateCourseDto> {
    status?: CourseStatus;
}

export interface CreateModuleDto {
    title: string;
    description?: string;
    order: number;
    courseId: string;
    prerequisites?: string[];
}

export interface CreateContentDto {
    title: string;
    description?: string;
    type: ContentType.CONTENT;
    order: number;
    body: string;
    mediaUrl?: string;
    resources?: string[];
    isRequired: boolean;
    estimatedDuration?: number;
}

export interface CreateEvaluationDto {
    title: string;
    description?: string;
    type: ContentType.EVALUATION;
    order: number;
    passingScore: number;
    maxAttempts?: number;
    timeLimit?: number;
    shuffleQuestions: boolean;
    shuffleAnswers: boolean;
    isRequired: boolean;
    estimatedDuration?: number;
}

export interface CreateQuestionDto {
    text: string;
    type: QuestionType;
    explanation?: string;
    points: number;
    order: number;
    answers: Omit<Answer, 'id' | 'createdAt' | 'updatedAt'>[];
}
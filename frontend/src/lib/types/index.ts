// Export types with explicit re-exports to avoid conflicts
export type {
	BaseEntity,
	Course,
	Module,
	Content,
	Evaluation,
	Question,
	Answer,
	Enrollment,
	EvaluationAttempt,
	UserProgress,
	ModuleContent,
	ModuleStatus,
	UserRole
} from './models/course';

// Export enums as values
export { ContentType, QuestionType } from './models/course';

export type { User } from './models/user';
export * from './dto';

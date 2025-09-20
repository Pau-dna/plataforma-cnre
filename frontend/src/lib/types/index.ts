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
	QuestionType,
	ContentType,
	ModuleStatus,
	UserRole
} from './models/course';

export type { User } from './models/user';
export * from './dto';

// Export all controllers
export * from './course';
export * from './module';
export * from './content';
export * from './evaluation';
export * from './question';
export * from './enrollment';
export * from './evaluationAttempt';
export * from './userProgress';
export * from './user';

// Create controller instances for easy import
import { CourseController } from './course';
import { ModuleController } from './module';
import { ContentController } from './content';
import { EvaluationController } from './evaluation';
import { QuestionController, AnswerController } from './question';
import { EnrollmentController } from './enrollment';
import { EvaluationAttemptController } from './evaluationAttempt';
import { UserProgressController } from './userProgress';
import { UserController } from './user';

// Controller instances
export const courseController = new CourseController();
export const moduleController = new ModuleController();
export const contentController = new ContentController();
export const evaluationController = new EvaluationController();
export const questionController = new QuestionController();
export const answerController = new AnswerController();
export const enrollmentController = new EnrollmentController();
export const evaluationAttemptController = new EvaluationAttemptController();
export const userProgressController = new UserProgressController();
export const userController = new UserController();
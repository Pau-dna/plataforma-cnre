import { BaseController } from './base';
import type {
	Question,
	CreateQuestionDTO,
	UpdateQuestionDTO,
	Answer,
	CreateAnswerDTO,
	UpdateAnswerDTO
} from '$lib/types';

export class QuestionController extends BaseController {
	/**
	 * Get a specific question by ID
	 */
	async getQuestion(id: number): Promise<Question> {
		return this.get<Question>(`/api/v1/questions/${id}`);
	}

	/**
	 * Create new question
	 */
	async createQuestion(questionData: CreateQuestionDTO): Promise<Question> {
		return this.post<Question>('/api/v1/questions', questionData);
	}

	/**
	 * Update existing question
	 */
	async updateQuestion(id: number, questionData: UpdateQuestionDTO): Promise<Question> {
		return this.put<Question>(`/api/v1/questions/${id}`, questionData);
	}

	/**
	 * Partially update a question (PATCH)
	 */
	async updateQuestionPatch(
		id: number,
		questionData: Partial<UpdateQuestionDTO>
	): Promise<Question> {
		return this.patch<Question>(`/api/v1/questions/${id}`, questionData);
	}

	/**
	 * Delete question
	 */
	async deleteQuestion(id: number): Promise<void> {
		return this.delete(`/api/v1/questions/${id}`);
	}

	/**
	 * Get questions by evaluation
	 */
	async getQuestionsByEvaluation(evaluationId: number): Promise<Question[]> {
		return this.get<Question[]>(`/api/v1/evaluations/${evaluationId}/questions`);
	}

	/**
	 * Get question with answers
	 * Note: Backend endpoint returns Answer[] not Question, so we combine the data
	 */
	async getQuestionWithAnswers(questionId: number): Promise<Question> {
		// Get question and answers separately since backend doesn't return combined data
		const question = await this.getQuestion(questionId);
		const answers = await this.get<Answer[]>(`/api/v1/questions/${questionId}/answers`);
		
		// Combine them into the expected structure
		return {
			...question,
			answers: answers || []
		};
	}
}

export class AnswerController extends BaseController {
	/**
	 * Get a specific answer by ID
	 */
	async getAnswer(id: number): Promise<Answer> {
		return this.get<Answer>(`/api/v1/answers/${id}`);
	}

	/**
	 * Create new answer
	 */
	async createAnswer(answerData: CreateAnswerDTO): Promise<Answer> {
		return this.post<Answer>('/api/v1/answers', answerData);
	}

	/**
	 * Update existing answer
	 */
	async updateAnswer(id: number, answerData: UpdateAnswerDTO): Promise<Answer> {
		return this.put<Answer>(`/api/v1/answers/${id}`, answerData);
	}

	/**
	 * Partially update an answer (PATCH)
	 */
	async updateAnswerPatch(id: number, answerData: Partial<UpdateAnswerDTO>): Promise<Answer> {
		return this.patch<Answer>(`/api/v1/answers/${id}`, answerData);
	}

	/**
	 * Delete answer
	 */
	async deleteAnswer(id: number): Promise<void> {
		return this.delete(`/api/v1/answers/${id}`);
	}

	/**
	 * Get all answers for a specific question
	 */
	async getAnswersByQuestion(questionId: number): Promise<Answer[]> {
		return this.get<Answer[]>(`/api/v1/questions/${questionId}/answers`);
	}
}

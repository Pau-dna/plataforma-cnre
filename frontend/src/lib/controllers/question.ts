import { apiClient } from '$lib/client';
import type {
	Question,
	CreateQuestionDTO,
	UpdateQuestionDTO,
	Answer,
	CreateAnswerDTO,
	UpdateAnswerDTO
} from '$lib/types';

export class QuestionController {
	/**
	 * Get a specific question by ID
	 */
	async getQuestion(id: number): Promise<Question> {
		return apiClient.get<Question>(`/api/v1/questions/${id}`);
	}

	/**
	 * Create new question
	 */
	async createQuestion(questionData: CreateQuestionDTO): Promise<Question> {
		return apiClient.post<Question>('/api/v1/questions', questionData);
	}

	/**
	 * Update existing question
	 */
	async updateQuestion(id: number, questionData: UpdateQuestionDTO): Promise<Question> {
		return apiClient.put<Question>(`/api/v1/questions/${id}`, questionData);
	}

	/**
	 * Delete question
	 */
	async deleteQuestion(id: number): Promise<void> {
		return apiClient.delete(`/api/v1/questions/${id}`);
	}
}

export class AnswerController {
	/**
	 * Get a specific answer by ID
	 */
	async getAnswer(id: number): Promise<Answer> {
		return apiClient.get<Answer>(`/api/v1/answers/${id}`);
	}

	/**
	 * Create new answer
	 */
	async createAnswer(answerData: CreateAnswerDTO): Promise<Answer> {
		return apiClient.post<Answer>('/api/v1/answers', answerData);
	}

	/**
	 * Update existing answer
	 */
	async updateAnswer(id: number, answerData: UpdateAnswerDTO): Promise<Answer> {
		return apiClient.put<Answer>(`/api/v1/answers/${id}`, answerData);
	}

	/**
	 * Delete answer
	 */
	async deleteAnswer(id: number): Promise<void> {
		return apiClient.delete(`/api/v1/answers/${id}`);
	}

	/**
	 * Get all answers for a specific question
	 */
	async getAnswersByQuestion(questionId: number): Promise<Answer[]> {
		return apiClient.get<Answer[]>(`/api/v1/questions/${questionId}/answers`);
	}
}
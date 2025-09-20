import { apiClient } from '$lib/client';
import type { User } from '$lib/types';

export class UserController {
	/**
	 * Get current user profile
	 */
	async getCurrentUser(): Promise<User> {
		return apiClient.get<User>('/auth/me');
	}

	/**
	 * Get user by ID
	 */
	async getUser(id: number): Promise<User> {
		return apiClient.get<User>(`/api/v1/users/${id}`);
	}

	/**
	 * Update user profile
	 */
	async updateUser(id: number, userData: Partial<User>): Promise<User> {
		return apiClient.put<User>(`/api/v1/users/${id}`, userData);
	}

	/**
	 * Delete user
	 */
	async deleteUser(id: number): Promise<void> {
		return apiClient.delete(`/api/v1/users/${id}`);
	}

	/**
	 * Get all users (admin only)
	 */
	async getAllUsers(): Promise<User[]> {
		return apiClient.get<User[]>('/api/v1/users');
	}

	/**
	 * Search users by email or name
	 */
	async searchUsers(query: string): Promise<User[]> {
		return apiClient.get<User[]>(`/api/v1/users/search?q=${encodeURIComponent(query)}`);
	}
}
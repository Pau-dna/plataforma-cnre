import { BaseController } from './base';
import type { User } from '$lib/types';

export class UserController extends BaseController {
	/**
	 * Get current user profile
	 */
	async getCurrentUser(): Promise<User> {
		return this.get<User>('/auth/me');
	}

	/**
	 * Get user by ID
	 */
	async getUser(id: number): Promise<User> {
		return this.get<User>(`/api/v1/users/${id}`);
	}

	/**
	 * Update user profile
	 */
	async updateUser(id: number, userData: Partial<User>): Promise<User> {
		return this.put<User>(`/api/v1/users/${id}`, userData);
	}

	/**
	 * Partially update a user (PATCH)
	 */
	async updateUserPatch(id: number, userData: Partial<User>): Promise<User> {
		return this.patch<User>(`/api/v1/users/${id}`, userData);
	}

	/**
	 * Delete user
	 */
	async deleteUser(id: number): Promise<void> {
		return this.delete(`/api/v1/users/${id}`);
	}

	/**
	 * Get all users (admin only)
	 */
	async getAllUsers(): Promise<User[]> {
		return this.get<User[]>('/api/v1/users');
	}

	/**
	 * Search users by email or name
	 */
	async searchUsers(query: string): Promise<User[]> {
		return this.get<User[]>(`/api/v1/users/search?q=${encodeURIComponent(query)}`);
	}
}

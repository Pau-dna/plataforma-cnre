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
	 * WARNING: This endpoint doesn't exist in backend - user handler not implemented
	 */
	async getUser(id: number): Promise<User> {
		throw new Error('User endpoints not implemented in backend');
	}

	/**
	 * Update user profile
	 * WARNING: This endpoint doesn't exist in backend - user handler not implemented
	 */
	async updateUser(id: number, userData: Partial<User>): Promise<User> {
		throw new Error('User endpoints not implemented in backend');
	}

	/**
	 * Partially update a user (PATCH)
	 * WARNING: This endpoint doesn't exist in backend - user handler not implemented
	 */
	async updateUserPatch(id: number, userData: Partial<User>): Promise<User> {
		throw new Error('User endpoints not implemented in backend');
	}

	/**
	 * Delete user
	 * WARNING: This endpoint doesn't exist in backend - user handler not implemented
	 */
	async deleteUser(id: number): Promise<void> {
		throw new Error('User endpoints not implemented in backend');
	}

	/**
	 * Get all users (admin only)
	 * WARNING: This endpoint doesn't exist in backend - user handler not implemented
	 */
	async getAllUsers(): Promise<User[]> {
		throw new Error('User endpoints not implemented in backend');
	}

	/**
	 * Search users by email or name
	 * WARNING: This endpoint doesn't exist in backend - user handler not implemented
	 */
	async searchUsers(query: string): Promise<User[]> {
		throw new Error('User endpoints not implemented in backend');
	}
}

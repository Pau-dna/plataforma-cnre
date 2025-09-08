import type { BaseEntity } from "./course";

// User interface
export interface User extends BaseEntity {
    id: number;
	email: string;
	firstName: string;
	lastName: string;
	avatarUrl?: string;
	role: 'student' | 'instructor' | 'admin';
}
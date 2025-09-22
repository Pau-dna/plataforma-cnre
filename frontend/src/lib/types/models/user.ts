import type { BaseEntity, UserRole, Enrollment } from './course';

// User interface - matches backend User model
export interface User extends BaseEntity {
	email: string;
	first_name: string;
	last_name: string;
	avatar_url?: string;
	role: UserRole;
	enrollments?: Enrollment[];
}

import type { BaseEntity, UserRole, Enrollment } from './course';

// User interface - matches backend User model
export interface User extends BaseEntity {
	email: string;
	fullname: string;
	avatar_url: string;
	role: UserRole;
	enrollments?: Enrollment[];
}

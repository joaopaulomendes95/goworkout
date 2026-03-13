export interface User {
	id: number;
	username: string;
	email: string;
	bio?: string;
	created_at: string;
	updated_at: string;
}

export interface RegisterRequest {
	username: string;
	email: string;
	password: string;
	bio?: string;
}

export interface LoginRequest {
	username: string;
	password: string;
}

export interface AuthResponse {
	user: User;
	token: string;
}

export interface ApiError {
	error: string;
}

export interface FormFailure {
	username?: string;
	email?: string;
	bio?: string;
	error: string;
}

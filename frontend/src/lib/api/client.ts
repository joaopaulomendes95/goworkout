import type { RegisterRequest, LoginRequest, AuthResponse, ApiError, User } from './types';

const API_URL = 'http://localhost:8080/api';

class ApiClient {
	private token: string | null = null;

	setToken(token: string | null) {
		this.token = token;
	}

	private async request<T>(
		endpoint: string,
		options: RequestInit = {}
	): Promise<T> {
		const headers: HeadersInit = {
			'Content-Type': 'application/json',
			...options.headers,
		};

		if (this.token) {
			(headers as Record<string, string>)['Authorization'] = `Bearer ${this.token}`;
		}

		const response = await fetch(`${API_URL}${endpoint}`, {
			...options,
			headers,
		});

		if (!response.ok) {
			const error: ApiError = await response.json().catch(() => ({
				error: 'An unexpected error occurred',
			}));
			throw new Error(error.error);
		}

		return response.json();
	}

	async register(data: RegisterRequest): Promise<AuthResponse> {
		return this.request<AuthResponse>('/users', {
			method: 'POST',
			body: JSON.stringify(data),
		});
	}

	async login(data: LoginRequest): Promise<AuthResponse> {
		return this.request<AuthResponse>('/tokens/authentication', {
			method: 'POST',
			body: JSON.stringify(data),
		});
	}

	async getMe(): Promise<{ user: User }> {
		return this.request<{ user: User }>('/me');
	}

	async getHello(): Promise<string> {
		const response = await fetch(`${API_URL}/hello`, {
			headers: this.token
				? { Authorization: `Bearer ${this.token}` }
				: {},
		});
		return response.text();
	}
}

export const api = new ApiClient();

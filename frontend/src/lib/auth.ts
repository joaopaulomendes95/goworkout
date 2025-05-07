export interface AuthResponse {
  sucess: boolean,
  message: string,
  token?: string,
  username?: string,
}

export interface LoginCredentials {
  username: string;
  password: string;
}

export const AUTH_TOKEN_KEY = 'auth_token';

export function saveToken(token:string): void {
  localStorage.setItem(AUTH_TOKEN_KEY, token);
}

export function getToken(): string | null {
  return localStorage.getItem(AUTH_TOKEN_KEY);
}

export function removeToken(): void {
  localStorage.removeItem(AUTH_TOKEN_KEY);
}

export function isLoggedIn(): boolean {
  return !!getToken();
}

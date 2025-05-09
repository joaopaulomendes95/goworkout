// NOTE: With HttpOnly cookies for authentication (which is generally recommended for security),
// client-side JavaScript cannot directly access or manage the authentication token.
// The functions below (saveToken, getToken, removeToken, isLoggedIn) using localStorage
// would be for a scenario where the token is *not* HttpOnly and is managed client-side.
// This is generally less secure for auth tokens.

// For an HttpOnly cookie setup, the authentication state (`isLoggedIn`) should be derived
// from `$page.data.authenticated` (set by server-side hooks/layouts), and login/logout
// should be handled by SvelteKit actions that interact with the backend and manage cookies.

export interface AuthResponse {
  sucess: boolean,
  message: string,
  token?: string, // This would be relevant if token is sent in response body for client storage
  username?: string,
}

export interface LoginCredentials {
  username: string;
  password: string;
}

export const AUTH_TOKEN_KEY_LOCAL_STORAGE = 'auth_token_local'; // Use a different key if mixing strategies

// Example: If you had a non-HttpOnly token you wanted to store
export function saveNonHttpOnlyToken(token: string): void {
  if (typeof localStorage !== 'undefined') {
    localStorage.setItem(AUTH_TOKEN_KEY_LOCAL_STORAGE, token);
  }
}

export function getNonHttpOnlyToken(): string | null {
  if (typeof localStorage !== 'undefined') {
    return localStorage.getItem(AUTH_TOKEN_KEY_LOCAL_STORAGE);
  }
  return null;
}

export function removeNonHttpOnlyToken(): void {
  if (typeof localStorage !== 'undefined') {
    localStorage.removeItem(AUTH_TOKEN_KEY_LOCAL_STORAGE);
  }
}

// This would reflect the presence of a non-HttpOnly token.
// For HttpOnly, use the derived store: `import { isAuthenticated } from '$lib/stores/auth_store';`
export function isClientSideTokenPresent(): boolean {
  return !!getNonHttpOnlyToken();
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

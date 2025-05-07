import { goto } from '$app/navigation';
import { getToken } from '$lib/auth';

// Base API URL
const API_BASE_URL = 'http://app:8080/';

// fetch wrapper that adds authentication
export async function apiRequest<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const token = getToken();

  const headers = {
    'Content-type': 'application/json',
    ...(token ? { 'Authorization': `Bearer ${token}` } : {}),
    ...(options.headers || {} )
  };

  const response = await fetch (`${API_BASE_URL}${endpoint}`, {
    ...options,
    headers
  });

  // Handle 401 Unauthorized
  if (response.status === 401) {
    removeToken();
    goto('/login');
    throw new Error('Authentication required');
  }

  const data = await response.json();

  if (!response.ok) {
    throw new Error(data.message || 'API request failed');
  }

  return data as T;
}  

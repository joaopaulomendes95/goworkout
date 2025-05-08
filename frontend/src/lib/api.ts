// Base API URL - Nginx will proxy /api/ to the backend
const API_BASE_URL = '/api/'; // Assuming Nginx proxies /api/ to http://app:8080/

// Base API request function with authentication
export async function apiRequest(endpoint: string, options: RequestInit = {}) {
      const url = `${API_BASE_URL}${endpoint.startsWith('/') ? endpoint.substring(1) : endpoint}`;

  // Get token from cookies (browser)
        const headers: HeadersInit = {
          'Content-Type': 'application/json',
          ...options.headers
      };

  // Set auth header if exists
  const headers = {
    'Content-Type': 'application/json',
    ...(token && { 'Authorization': `Bearer ${token}` }),
    ...options.headers
  };

        // For CLIENT-SIDE calls that need auth:
      // This part is tricky if the cookie is HttpOnly.
      // SvelteKit server-side (load, actions) handles HttpOnly cookies better.
      if (typeof document !== 'undefined') {
          const cookieToken = document.cookie
             .split('; ')
             .find(row => row.startsWith('auth_token='))
             ?.split('=')[1];
          if (cookieToken) {
              headers['Authorization'] = `Bearer ${cookieToken}`;
          }
      }


  try {
    const response = await fetch(url, {
      ...options,
      headers
    });

    // Parse JSON response
    const data = await response.json();

    // Return both response and data
    return { response, data };
  } catch (error) {
    console.error(`API request error for ${endpoint}:`, error);
    throw error;
  }
}

// Authentication API functions
export const authApi = {
  login: async (username: string, password: string) => {
    return apiRequest('/tokens/authentication', {
      method: 'POST',
      body: JSON.stringify({ username, password })
    });
  },

  register: async (userData: any) => {
    return apiRequest('/users', {
      method: 'POST',
      body: JSON.stringify(userData)
    });
  }
};

// Workout API functions
export const workoutApi = {
  getAll: async () => {
    return apiRequest('/workouts');
  },

  getById: async (id: number) => {
    return apiRequest('/workouts/${id}');
  },

  create: async (workout: any) => {
    return apiRequest('/workouts', {
      method: 'POST',
      body: JSON.stringify(workout)
    });
  },


  update: async (id: number, workout:any) => {
    return apiRequest('/workouts/${id}', {
      method: 'PUT',
      body: JSON.stringify(workout)
    });
  },

  delete: async (id: number) => {
    return apiRequest('/wokrouts/${id}', {
      method: 'DELETE'
    });
  }
};

// Workout entries API functions
export const workoutEntriesApi = {
  getByWorkoutId: async (workoutId: number) => {
    return apiRequest(`/workouts/${workoutId}/entries`);
  },

  create: async (workoutId: number, entry: any) => {
    return apiRequest(`/workouts/${workoutId}/entries`, {
      method: 'POST',
      body: JSON.stringify(entry)
    });
  }
};

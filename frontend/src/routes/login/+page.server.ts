const API = 'http://localhost:8080';

export async function load() {
    return { message: "Please log in" };
}

export const actions = {
    default: async ({ request }) => {
        const data = await request.formData();
        const username = data.get('username');
        const password = data.get('password');

        try {
            const response = await fetch(`${API}/login`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ username, password })
            });

            const result = await response.json();

            if (response.ok) {
                return { success: true, message: 'Login successful', token: result.token };
            } else {
                return { success: false, message: result.message || 'Login failed' };
            }
        } catch (e) {
            return { success: false, message: 'Network error' };
        }
    }
};

const API = 'http://localhost:8080';

export async function load() {
    return { message: "Please log in" };
}

interface ApiResponse {
    sucess: boolean;
    message: string;
    token?: string
}

export const actions = {
    user_register: async ({ request }) => {
        const data = await request.formData();
        const username = data.get('username');
        const email = data.get('email');
        const password = data.get('password');
        const bio = data.get('bio');

        //debugggin
        console.log("FormData: ", data)


        try {
            const response: ApiResponse = await fetch(`${API}/login`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ username, email, password, bio })
            });

            const result = await response.json();

            if (response.sucess) {
                return { success: true, message: 'Login successful', token: result.token };
            } else {
                return { success: false, message: result.message || 'Login failed' };
            }

        } catch (e) {
            return { success: false, message: 'Network error' };
        }
    }
};

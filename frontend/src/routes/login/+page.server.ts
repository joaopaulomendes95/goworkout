const API = 'http://app:8080';

// TODO: redirect if sucess
// import { redirect } from '@sveltejs/kit';

export async function load() {
}

export const actions = {
    user_login: async ({ request }) => {
        const data = await request.formData();
        const username = data.get('username');
        const password = data.get('password');

        //debugggin
        console.log("FormData: ", data)


        try {
            const response = await fetch(`${API}/tokens/authentication`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ username, password })
            });

            console.log("Response: ", response);
            const result = await response.json();
            console.log("Result: ", result);


            if (response.ok) {
                return {
                    success: true,
                    message: 'Login successful',
                    token: result.token,
                };
            } else {
                return {
                    success: false,
                    message: result.message || 'login failed',
                    username // Return data to repopulate form
                };
            }

        } catch (e) {
            return {
                success: false,
                message: 'Network error',
                username // Return data to repopulate form
            };
        }
    }
};

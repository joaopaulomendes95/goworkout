const API = 'http://app:8080';

export async function load() {
    return { message: "Please log in" };
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
            const response = await fetch(`${API}/users`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ username, email, password, bio })
            });

            console.log("Response: ", response);
            const result = await response.json();

            if (response.ok) {
                return {
                    success: true,
                    message: 'Registation successful',
                    token: result.token,
                };
            } else {
                return {
                    success: false,
                    message: result.message || 'Registation failed',
                    username, email, bio // Return data to repopulate form
                };
            }

        } catch (e) {
            return {
                success: false,
                message: 'Network error',
                username, email, bio // Return data to repopulate form
            };
        }
    }
};

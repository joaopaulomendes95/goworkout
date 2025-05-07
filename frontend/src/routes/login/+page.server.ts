const API = 'http://app:8080';

// TODO: redirect if sucess and store the auth_token
import { redirect } from '@sveltejs/kit';
import type { Actions } from './$types';


export const actions: Actions = {
    login: async ({ request, cookies, fetch, locals }) => {
        // Get form data
        const data = await request.formData();
        const username = data.get('username')?.toString() || '';
        const password = data.get('password')?.toString() || '';

        // Make API request
        const response = await fetch(`${API}/tokens/authentication`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password })
        });

        console.log("Response: ", response);
        const result = await response.json();
        console.log("Result: ", result);

        if (response.ok && result.token) {
            locals.token = result.token;
            console.log('token: ', locals.token)
            locals.authenticated = true;
            cookies.set('auth_token', result.token, {
                path: '/',
                httpOnly: true,
                secure: true,
                sameSite: 'strict',
                maxAge: 60 * 60 * 24 // 24 hours
            });

            throw redirect(303, '/');
        }

        // Authentication failed
        return {
            success: false,
            message: result.message || 'authentication failed'
        };
    }
};

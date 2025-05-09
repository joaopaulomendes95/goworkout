const API = 'http://app:8080'; // Changed from 'http://localhost:8080'
import { error } from '@sveltejs/kit';

export async function load({ }) {
    async function getHealth() {
        try {
            const health = await fetch(`${API}/health`);
            return health.json();
        } catch (e) {
            console.error('Error fetching health', e);
            error(400, 'Failed to fetch health');
        }
    }

    const health = await getHealth();

    console.log('health', health);

    // const { user } = locals;
    // if (!user) {
    // 	return {
    // 		status: 302,
    // 		redirect: '/login',
    // 		message: 'You must be logged in to view this page'
    // 	};
    // }
    return {
        // user,
        health
    };
}

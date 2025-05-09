import { error } from '@sveltejs/kit';

const GO_API_URL = process.env.PRIVATE_GO_API_URL || 'http://app:8080';

export async function load({ fetch: svelteKitFetch }) {
	async function getHealth() {
		try {
			const healthResponse = await svelteKitFetch(`${GO_API_URL}/health`); // Use SvelteKit's fetch
			if (!healthResponse.ok) {
				console.error('Health check failed:', healthResponse.status, await healthResponse.text());
				// Return an error structure that the page can display
				return { status: 'down', error: `Backend health check failed with status ${healthResponse.status}` };
			}
			return await healthResponse.json();
		} catch (e: any) {
			console.error('Error fetching health:', e.message || e);
			return { status: 'down', error: 'Network error or backend unavailable for health check.' };
		}
	}

	const healthData = await getHealth();
	return {
		health: healthData // This will include status and potentially an error message
	};
}

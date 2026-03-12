const API = 'http://app:8080';
import { error } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async () => {
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
};

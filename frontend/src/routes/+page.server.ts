import { error } from '@sveltejs/kit';
import type { PageServerLoad, PageServerLoadEvent } from './$types';

export const load: PageServerLoad = async (event: PageServerLoadEvent) => { // Add event type
    async function getHealth() {
        try {
            // Use event.fetch for server-side requests in load functions
            const healthResponse = await event.fetch(`/api/health`); // Relative path
            if (!healthResponse.ok) {
                throw new Error(`Failed to fetch health: ${healthResponse.status}`);
            }
            return healthResponse.json();
        } catch (e: any) {
            console.error('Error fetching health in /+page.server.ts:', e.message);
            // SvelteKit's error helper creates a proper error page
            error(500, `Failed to fetch health: ${e.message}`);
        }
    }

    const health = await getHealth();

    return {
        health
    }
};

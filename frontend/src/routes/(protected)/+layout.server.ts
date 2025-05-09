import { redirect } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ url, parent }: { url: URL; parent: () => Promise<{ authenticated: boolean; user: any }> }) => {
  const parentData = await parent(); // Gets { authenticated, user } from root layout

  if (!parentData.authenticated) {
    // Redirect to login page and store intended destination
    throw redirect(303, `/login?redirectTo=${encodeURIComponent(url.pathname)}`);
  }
  
  // User is authenticated, proceed. Data is already available via parent().
  return {
    // authenticated: parentData.authenticated, // Already in parent data
    // user: parentData.user // Already in parent data
  };
};

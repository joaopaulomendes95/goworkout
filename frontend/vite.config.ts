import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig, loadEnv } from 'vite';

export default defineConfig(({ mode }) => {
	// Load .env files for the current mode
	// useful to dynamically set the proxy target
	// but docker 'http://app:8080' is static
	const env = loadEnv(mode, process.cwd(), '');

	return {
	plugins: [tailwindcss(), sveltekit()],
	server: {
		// Proxy /api makes requests to Go backend
		proxy: {
			'/api': {
				target: 'http://app:8080',
				changeOrigin: true,
				rewrite: (path) => path.replace(/^\/api/, ''),
			}
		},
		host: '0.0.0.0',
		port: 5173,
	}
	};
});

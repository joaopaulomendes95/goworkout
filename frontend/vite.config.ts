import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig, loadEnv } from 'vite';

export default defineConfig({
	plugins: [tailwindcss(), sveltekit()],
	server: {
		proxy: {
			'/api': {
				target: 'http://app:8080',
				changeOrigin: true,
				rewrite: (path) => path.replace(/^\/api/, ''),
			}
		},
	}
});

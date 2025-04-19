import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [tailwindcss(), sveltekit()],
	optimizeDeps: {
		include: ['@tanstack/svelte-query']
	},
	ssr: {
		noExternal: ['@tanstack/svelte-query']
	}
});

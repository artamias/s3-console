import { svelte } from '@sveltejs/vite-plugin-svelte';
import tailwindcss from '@tailwindcss/vite';
import { defineConfig, loadEnv } from 'vite';

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '');

  return {
    plugins: [tailwindcss(), svelte()],
    server: {
      host: env.VITE_DEV_HOST || 'localhost',
      proxy: {
        '/api': {
          target: env.VITE_BACKEND_ORIGIN || 'http://localhost:8080',
          changeOrigin: true,
        },
        '/health': {
          target: env.VITE_BACKEND_ORIGIN || 'http://localhost:8080',
          changeOrigin: true,
        },
      },
    },
  };
});

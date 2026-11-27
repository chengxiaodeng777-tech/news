import { defineConfig } from '@vben/vite-config';

export default defineConfig(async () => {
  return {
    application: {},
    vite: {
      resolve: {
        // Workspace packages must share Vue's component and injection context.
        dedupe: [
          'vue',
          '@vue/runtime-core',
          '@vue/runtime-dom',
          '@vue/reactivity',
        ],
      },
      server: {
        // Fail clearly instead of falling back to browser-blocked port 6000.
        strictPort: true,
        proxy: {
          '/api': {
            changeOrigin: true,
            target: 'http://localhost:8082',
            ws: true,
          },
        },
      },
    },
  };
});

import { fileURLToPath } from 'node:url';
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import legacy from '@vitejs/plugin-legacy'
import { nodePolyfills } from 'vite-plugin-node-polyfills'

import Components from 'unplugin-vue-components/vite'
import { NaiveUiResolver } from 'unplugin-vue-components/resolvers'

// https://vitejs.dev/config/
export default defineConfig({
  base: './',
  resolve: {
    alias: {
      '~': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  plugins: [
    vue(),
    Components({
      dirs: ['src/components'], // allow auto load Markdown components under `./src/components/`
      extensions: ['vue', 'md', 'jsx', 'tsx'], // allow auto import and register components used in Markdown
      include: [
        /\.vue$/,
        /\.vue\?vue/,
        /\.vue\.[tj]sx?\?vue/, // for vue-loader with experimentalInlineMatchResource enabled
        /\.vue\?v=/,
        /\.[tj]sx/,
      ],
      resolvers: [
        NaiveUiResolver(),
      ],
      dts: 'src/components.d.ts',
    }),
    legacy({
      targets: ['defaults', 'not IE 11']
    }),
    nodePolyfills(),
  ],
  build: {
    chunkSizeWarningLimit: 1024,
    rolldownOptions: {
      output: {
        advancedChunks: {
          groups: [
            { name: 'editer', test: /(codemirror|highlight)/ },
            { name: 'import-export', test: /(docx|fflate|hyparquet)/ },
            { name: 'naive-ui', test: /naive-ui/ },
          ],
        },
      },
    },
    sourcemap: false,
  },
  server: {
    proxy: {

      '/api': {
        changeOrigin: true,
        target: 'https://worker.firehomework.top/dice/api',
        // target: 'http://8.130.140.128:8082',
        // target: 'http://localhost:8088',

        rewrite: (path) => path.replace(/^\/api/, ''),

      },
    }
  },
})

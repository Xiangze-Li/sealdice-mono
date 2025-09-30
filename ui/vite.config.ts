import { fileURLToPath } from 'node:url';
import { defineConfig, loadEnv } from 'vite';
import vue from '@vitejs/plugin-vue';
import vueJsx from '@vitejs/plugin-vue-jsx';
import vueDevTools from 'vite-plugin-vue-devtools';
import legacy from '@vitejs/plugin-legacy';
import tailwindcss from '@tailwindcss/vite';
import AutoImport from 'unplugin-auto-import/vite';
import Components from 'unplugin-vue-components/vite';
import { NaiveUiResolver } from 'unplugin-vue-components/resolvers';
import Icons from 'unplugin-icons/vite';
import IconsResolver from 'unplugin-icons/resolver';

// https://vitejs.dev/config/
export default defineConfig(({ mode }) => ({
  base: './',
  resolve: {
    alias: [
      {
        find: /\/#\//,
        replacement: fileURLToPath(new URL('./types', import.meta.url)),
      },
      {
        find: '~',
        replacement: fileURLToPath(new URL('./src', import.meta.url)),
      },
      {
        find: '@',
        replacement: fileURLToPath(new URL('./src', import.meta.url)),
      },
    ],
  },
  server: {
    proxy: {
      '/sd-api': {
        target: loadEnv(mode, './').VITE_APP_APIURL,
        changeOrigin: true,
        rewrite: path => path,
      },
    },
  },
  plugins: [
    vue(),
    vueJsx(),
    vueDevTools(),
    tailwindcss(),
    AutoImport({
      include: [
        /\.[tj]sx?$/, // .ts, .tsx, .js, .jsx
        /\.vue$/,
        /\.vue\?vue/, // .vue
      ],
      imports: ['vue', 'pinia', 'vue-router', '@vueuse/core'],
      dirs: ['src/api/v2/**'],
      dts: true,
      vueTemplate: true,
      resolvers: [NaiveUiResolver(), IconsResolver()],
    }),
    Components({
      dirs: ['src/components', 'src/pages', 'src/views'],
      extensions: ['vue', 'jsx', 'tsx'],
      include: [
        /\.vue$/,
        /\.vue\?vue/,
        /\.vue\.[tj]sx?\?vue/, // for vue-loader with experimentalInlineMatchResource enabled
        /\.vue\?v=/,
        /\.[tj]sx/,
      ],
      resolvers: [NaiveUiResolver(), IconsResolver()],
    }),
    Icons({
      compiler: 'vue3',
      autoInstall: true,
    }),
    legacy({
      targets: ['defaults'],
    }),
  ],
  build: {
    sourcemap: false,
    chunkSizeWarningLimit: 1024,
    rolldownOptions: {
      output: {
        advancedChunks: {
          groups: [
            { name: 'codemirror', test: /codemirror/ },
            { name: 'naive-ui', test: /naive-ui/ },
          ],
        },
      },
    },
  },
}));

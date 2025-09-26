import '~/styles/index.css';
import { createApp } from 'vue';
import { useStore } from '~/store';
import App from './App.vue';
import router from '~/router';
import VueDiff from 'vue-diff';
import { VueQueryPlugin } from '@tanstack/vue-query';
import 'vue-diff/dist/index.css';
import yaml from 'highlight.js/lib/languages/yaml';
import toml from 'highlight.js/lib/languages/ini';

const loading = useStorage('router-view-loading', true);

router.beforeEach((to, from, next) => {
  loading.value = true;
  next();
});
// eslint-disable-next-line @typescript-eslint/no-unused-vars
router.afterEach((to, from) => {
  setTimeout(() => {
    loading.value = false;
  }, 300);
});

const app = createApp(App);
app.use(createPinia());

VueDiff.hljs.registerLanguage('yaml', yaml);
VueDiff.hljs.registerLanguage('toml', toml);
app
  .use(VueDiff, {
    componentName: 'VueDiff',
  })
  .use(router)
  .use(VueQueryPlugin);

const store = useStore();

store.trySignIn().then(() => {
  app.mount('#app');
});

try {
  (window as any).store = store;
  (globalThis as any).store = store;
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
} catch (e) {}

// app.use(ElementPlus);

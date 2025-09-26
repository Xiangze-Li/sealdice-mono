<!-- eslint-disable vue/multi-word-component-names -->
<template>
  <n-menu :options="options" :value="route.path" accordion />
</template>

<script lang="tsx" setup>
import { useStore } from '~/store';
import type { ComputedRef, ModelRef } from 'vue';
import type { MenuOption } from 'naive-ui';

const advancedConfigCounter: ModelRef<number> = defineModel('advancedConfigCounter', {
  default: 0,
});

const store = useStore();

const route = useRoute();

const customTextOptions: ComputedRef<MenuOption[]> = computed(() =>
  Object.keys(store.curDice.customTexts ?? {}).map(k => ({
    key: `/custom-text/${k}`,
    label: () => <router-link to={`/custom-text/${k}`}>{k}</router-link>,
  })),
);

const miscOptions: ComputedRef<MenuOption[]> = computed(() => {
  const list = [
    {
      key: '/misc/base-setting',
      label: () => <router-link to="/misc/base-setting">基本设置</router-link>,
    },
    {
      key: '/misc/group',
      label: () => <router-link to="/misc/group">群组管理</router-link>,
    },
    {
      key: '/misc/ban',
      label: () => <router-link to="/misc/ban">黑白名单</router-link>,
    },
    {
      key: '/misc/dice-public',
      label: () => <router-link to="/misc/dice-public">公骰设置</router-link>,
    },
    {
      key: '/misc/backup',
      label: () => <router-link to="/misc/backup">备份</router-link>,
    },
  ];
  if (advancedConfigCounter.value >= 8) {
    list.push({
      key: '/misc/advanced-setting',
      label: () => <router-link to="/misc/advanced-setting">高级设置</router-link>,
    });
  }
  return list;
});

const options: ComputedRef<MenuOption[]> = computed(() => [
  {
    key: '/home',
    label: () => <router-link to="/home">主页</router-link>,
    icon: () => (
      <n-icon>
        <i-carbon-home />
      </n-icon>
    ),
  },
  {
    key: '/connect',
    label: () => <router-link to="/connect">账号设置</router-link>,
    icon: () => (
      <n-icon>
        <i-carbon-link />
      </n-icon>
    ),
  },
  {
    key: '/custom-text',
    label: '自定义文案',
    icon: () => (
      <n-icon>
        <i-carbon-edit />
      </n-icon>
    ),
    children: customTextOptions.value,
  },
  {
    key: '/mod',
    label: '扩展功能',
    icon: () => (
      <n-icon>
        <i-carbon-settings />
      </n-icon>
    ),
    children: [
      {
        key: '/mod/reply',
        label: () => <router-link to="/mod/reply">自定义回复</router-link>,
      },
      {
        key: '/mod/deck',
        label: () => <router-link to="/mod/deck">牌堆管理</router-link>,
      },
      {
        key: '/mod/story',
        label: () => <router-link to="/mod/story">跑团日志</router-link>,
      },
      {
        key: '/mod/js',
        label: () => <router-link to="/mod/js">JS 扩展</router-link>,
      },
      {
        key: '/mod/helpdoc',
        label: () => <router-link to="/mod/helpdoc">帮助文档</router-link>,
      },
      {
        key: '/mod/censor',
        label: () => <router-link to="/mod/censor">拦截管理</router-link>,
      },
    ],
  },
  {
    key: '/misc',
    label: '综合设置',
    icon: () => (
      <n-icon>
        <i-carbon-settings-adjust />
      </n-icon>
    ),
    children: miscOptions.value,
  },
  {
    key: '/tool',
    label: '辅助工具',
    icon: () => (
      <n-icon>
        <i-carbon-tool-kit />
      </n-icon>
    ),
    children: [
      {
        key: '/tool/test',
        label: () => <router-link to="/tool/test">测试</router-link>,
      },
      {
        key: '/tool/resource',
        label: () => <router-link to="/tool/resource">资源管理</router-link>,
      },
    ],
  },
  {
    key: '/about',
    label: () => <router-link to="/about">关于</router-link>,
    icon: () => (
      <n-icon>
        <i-carbon-star />
      </n-icon>
    ),
  },
]);
</script>

<style lang="css" scoped></style>

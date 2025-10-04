<script setup lang="ts">
import { urlBase } from '~/backend';
import imgHaibao from '~/assets/haibao1.png';

const props = withDefaults(
  defineProps<{
    username: string;
    user?: string;
    src?: string;
    onlyName?: boolean;
    href?: string;
    info?: string;
  }>(),
  {
    onlyName: false,
  },
);

const names = new Map<string, string>([
  ['木落', 'fy0'],
  ['逐风', 'MintCider'],
  ['暮星', 'MX-fox'],
  ['JohNSoN', 'Xiangze-Li'],
  ['Bugtower100', 'bugtower100'],
  ['只是另一个 ID', 'JustAnotherID'],
  ['檀轶步棋', 'oissevalt'],
  ['脑', 'f44r'],
]);

const getUser = () => {
  return props.user ?? names.get(props.username) ?? props.username;
};

const getAvatarSrc = () => {
  if (props.src) {
    return props.src;
  }
  return `${urlBase}/sd-api/utils/ga/${getUser()}`;
};

const getHref = () => {
  if (props.onlyName) {
    return 'javascript:void(0)';
  }
  if (props.href) {
    return props.href;
  }
  return `https://github.com/${getUser()}`;
};
</script>

<template>
  <n-flex size="small" vertical>
    <n-button v-if="!props.onlyName" text tag="a" :href="getHref()" target="_blank">
      <n-flex size="small" align="center">
        <n-avatar round size="large" :src="getAvatarSrc()" :fallback-src="imgHaibao" />
        <n-flex size="small" vertical align="start">
          <n-text strong>{{ props.username ?? getUser() }}</n-text>
          <n-text class="text-xs" italic v-if="props.info">{{ props.info }}</n-text>
        </n-flex>
      </n-flex>
    </n-button>
    <n-flex v-else size="small" align="center">
      <n-avatar round size="large" :src="imgHaibao" />
      <n-flex size="small" vertical align="start">
        <n-text class="leading-none" strong>{{ props.username ?? getUser() }}</n-text>
        <n-text class="text-xs" italic v-if="props.info">{{ props.info }}</n-text>
      </n-flex>
    </n-flex>
  </n-flex>
</template>

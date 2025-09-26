<script setup lang="tsx">
import dayjs from 'dayjs';
import { getUtilsCheckNetWorkHealth } from '~/api/v1/utils';

const now = ref<dayjs.Dayjs>(dayjs());
let timerId: number;
let checkTimerId: number;

const autoRefresh = ref(true);
const networkHealth = ref({
  total: 0,
  ok: [],
  timestamp: 0,
} as {
  total: number;
  ok: string[];
  timestamp: number;
});

const getWebsiteHealthComponent = (key: string, name: string): VNode => {
  const ok = networkHealth.value.ok?.includes(key);
  if (ok) {
    return (
      <n-tag round bordered={false} type="success" size="small">
        {{
          default: () => name,
          icon: () => (
            <n-icon>
              <i-carbon-checkmark-filled />
            </n-icon>
          ),
        }}
      </n-tag>
    );
  } else {
    return (
      <n-tag round bordered={false} type="error" size="small">
        {{
          default: () => name,
          icon: () => (
            <n-icon>
              <i-carbon-close-filled />
            </n-icon>
          ),
        }}
      </n-tag>
    );
  }
};

const refreshNetworkHealth = async () => {
  networkHealth.value.timestamp = 0;
  const ret = await getUtilsCheckNetWorkHealth();
  if (ret.result) {
    networkHealth.value = ret;
  }
};

onBeforeMount(async () => {
  if (autoRefresh.value) {
    await refreshNetworkHealth();
  }

  timerId = setInterval(() => {
    now.value = dayjs();
  }, 5000) as any;
  checkTimerId = setInterval(
    async () => {
      await refreshNetworkHealth();
    },
    5 * 60 * 1000,
  ) as any; // 5 min 一次
});

onBeforeUnmount(() => {
  clearInterval(timerId);
  clearInterval(checkTimerId);
});
</script>

<template>
  <div class="flex flex-wrap items-center gap-1" @click="refreshNetworkHealth">
    <n-tooltip trigger="hover">
      <template #trigger>
        <n-text>网络质量：</n-text>
      </template>
      点击重新进行检测
    </n-tooltip>

    <n-text v-if="networkHealth.timestamp === 0">检测中…… 🤔</n-text>
    <n-text
      type="success"
      v-else-if="networkHealth.total !== 0 && networkHealth.total === networkHealth.ok?.length">
      优 😄
    </n-text>
    <n-text
      type="default"
      v-else-if="networkHealth.ok?.includes('sign') && networkHealth.ok?.includes('seal')">
      一般 😐️
    </n-text>
    <n-text
      type="error"
      v-else-if="networkHealth.total !== 0 && (networkHealth.ok ?? []).length === 0">
      网络中断 😱
    </n-text>
    <template v-else>
      <n-text type="warning" class="mr-4">差 ☹️</n-text>
      <n-text type="warning" class="text-xs">
        这意味着你可能无法正常使用内置客户端/Lagrange 连接 QQ
        平台，有时会出现消息无法正常发送的现象。
      </n-text>
    </template>

    <n-tooltip v-if="networkHealth.timestamp !== 0">
      <template #trigger>
        <span>
          <n-text class="ml-auto text-xs" type="info">
            检测于 {{ dayjs.unix(networkHealth.timestamp).from(now) }}
          </n-text>
        </span>
      </template>
      {{ dayjs.unix(networkHealth.timestamp).format('YYYY-MM-DD HH:mm:ss') }}
    </n-tooltip>
  </div>

  <n-flex v-if="networkHealth.timestamp !== 0" size="small" align="center" class="mx-2">
    <component :is="getWebsiteHealthComponent('seal', '官网')"></component>
    <component :is="getWebsiteHealthComponent('sign', 'Lagrange Sign')"></component>
    <component :is="getWebsiteHealthComponent('google', 'Google')"></component>
    <component :is="getWebsiteHealthComponent('github', 'GitHub')"></component>
  </n-flex>
</template>

<style scoped lang="css"></style>

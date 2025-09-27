<script setup lang="tsx">
import dayjs from 'dayjs';
import { getUtilsCheckNetWorkHealth } from '~/api/v1/utils';

const now = ref<dayjs.Dayjs>(dayjs());
let timerId: number;
let checkTimerId: number;

const autoRefresh = ref(true);
const networkHealth = ref({
  total: 0,
  timestamp: 0,
  targets: [],
} as {
  total: number;
  timestamp: number;
  targets: {
    target: string;
    ok: boolean;
    duration: number;
  }[];
});
const networkHealths = ref<Map<string, number>>(new Map());

const getWebsiteHealthComponent = (key: string, name: string): VNode => {
  const duration = networkHealths.value.get(key);
  if (duration) {
    return (
      <n-tag round bordered={false} type="success" size="small">
        {{
          default: () => <span>{`${name} - ${(duration / 1000000).toFixed(2)}ms`}</span>,
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
    for (let target of ret.targets) {
      if (target.ok) {
        networkHealths.value.set(target.target, target.duration);
      }
    }
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
      v-else-if="networkHealth.total !== 0 && networkHealth.total === networkHealths?.size">
      优 😄
    </n-text>
    <n-text type="default" v-else-if="networkHealths?.get('sign') && networkHealths?.get('seal')">
      一般 😐️
    </n-text>
    <n-text type="error" v-else-if="networkHealth.total !== 0 && networkHealths?.size === 0">
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

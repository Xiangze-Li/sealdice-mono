<template>
  <n-flex align="center" justify="space-between" wrap>
    <n-switch v-model:value="censorEnable" @update:value="enableChange">
      <template #checked>启用</template>
      <template #unchecked>关闭</template>
    </n-switch>
    <n-button v-show="censorEnable" type="primary" @click="restartCensor">
      <template #icon>
        <i-carbon-renew />
      </template>
      重载拦截
    </n-button>
  </n-flex>

  <n-affix v-if="censorStore.needReload" :top="60">
    <tip-box type="error">
      <n-text type="error" class="text-base" tag="strong">存在修改，需要重载后生效！</n-text>
    </tip-box>
  </n-affix>

  <template v-if="censorEnable">
    <n-tabs v-model:value="tab" justify-content="space-evenly">
      <n-tab-pane tab="拦截设置" name="setting">
        <censor-config />
      </n-tab-pane>

      <n-tab-pane tab="敏感词管理" name="word">
        <censor-word-tip />
        <censor-files />
        <censor-words />
      </n-tab-pane>

      <n-tab-pane tab="拦截日志" name="log">
        <censor-log />
      </n-tab-pane>
    </n-tabs>
  </template>
  <template v-else>
    <n-text type="error" class="mt-4 block text-2xl">请先启用拦截！</n-text>
  </template>
</template>

<script lang="ts" setup>
import { getCensorStatus } from '~/api/v1/censor';
onBeforeMount(() => {
  refreshCensorStatus();
});

import { useCensorStore } from '~/components/censor/censor.ts';

const censorEnable = ref<boolean>(false);

const censorStore = useCensorStore();

const refreshCensorStatus = async () => {
  const status:
    | { result: false }
    | {
        result: true;
        enable: boolean;
        isLoading: boolean;
      } = await getCensorStatus();
  if (status.result) {
    censorEnable.value = status.enable;
  }
};

const restartCensor = async () => {
  const restart = await censorStore.restartCensor();
  if (restart.result) {
    censorEnable.value = restart.enable;
    censorStore.reload();
  }
};

const stopCensor = async () => {
  const stop = await censorStore.stopCensor();
  if (stop.result) {
    censorEnable.value = false;
  }
};

const enableChange = async (value: boolean | number | string) => {
  if (value === true) {
    await restartCensor();
  } else {
    await stopCensor();
  }
};

const tab = ref('setting');
</script>

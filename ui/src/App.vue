<script setup lang="ts">
import { darkTheme, dateZhCN, lightTheme, zhCN } from 'naive-ui';
import { useStore } from '~/store';
import { VueQueryDevtools } from '@tanstack/vue-query-devtools';
import MainBox from '~/pages/MainBox.vue';

const isDark = useDark({ disableTransition: false });
const store = useStore();

const testMode = ref<boolean>(false);

onMounted(async () => {
  const info = await store.getPreInfo();
  testMode.value = info.testMode;
});
</script>

<template>
  <n-config-provider
    :theme="isDark ? darkTheme : lightTheme"
    :locale="zhCN"
    :date-locale="dateZhCN"
    class="mx-auto size-full">
    <n-dialog-provider>
      <n-message-provider>
        <n-modal-provider>
          <n-notification-provider>
            <n-watermark
              :content="testMode ? '仅用于展示 修改无效' : ''"
              :cross="true"
              :font-size="22"
              :height="300"
              :line-height="22"
              :rotate="-15"
              :width="240"
              :x-offset="12"
              :y-offset="60"
              full-screen>
              <MainBox />
            </n-watermark>
          </n-notification-provider>
        </n-modal-provider>
      </n-message-provider>
    </n-dialog-provider>
  </n-config-provider>
  <VueQueryDevtools />
</template>

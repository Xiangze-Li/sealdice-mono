<template>
  <n-affix v-if="modified" :top="120" :trigger-top="60" class="w-full">
    <tip-box type="error">
      <n-text type="error" tag="strong" class="text-lg">内容已修改，不要忘记保存！</n-text>
    </tip-box>
  </n-affix>

  <h2 class="h-2">高级设置</h2>
  <tip-box type="warning" class="my-4">
    <n-text>
      此处是面向开发者或进阶用户的隐藏设置页，下列的设置项可能会对海豹核心的功能造成重大影响。<br />
      一些尚在测试的不稳定设置项，以及
      <strong>普通骰主无需关注</strong> 的设置项会被放在此处。<br />
      此处的设置项不保证稳定提供，在未来版本随时可能会被移除。<br />
      <br />
      <strong>除非你知道自己在做什么，否则不要修改此处的任何设置项！</strong><br />
      <br />
      <em>
        如果你误操作修改了此处设置，希望恢复默认，请手动删除
        <n-text code>data/default/advanced.yaml</n-text>
        文件。
      </em>
    </n-text>
  </tip-box>

  <n-form label-placement="left" label-width="auto">
    <n-form-item label="显示高级设置页">
      <template #label>
        <span>显示高级设置页</span>
        <n-tooltip>
          <template #trigger>
            <n-icon>
              <i-carbon-help-filled />
            </n-icon>
          </template>
          设置是否显示高级设置页，只影响展示
        </n-tooltip>
      </template>
      <n-switch v-model:value="config.show" />
    </n-form-item>
    <n-form-item label="启用高级设置">
      <template #label>
        <span>启用高级设置</span>
        <n-tooltip>
          <template #trigger>
            <n-icon>
              <i-carbon-help-filled />
            </n-icon>
          </template>
          设置是否启用高级设置，关闭时下列设置无效
        </n-tooltip>
      </template>
      <n-switch v-model:value="config.enable" />
    </n-form-item>

    <h3>自定义回复</h3>
    <n-form-item label="开启回复调试日志">
      <template #label>
        <span>回复调试日志</span>
        <n-tooltip>
          <template #trigger>
            <n-icon>
              <i-carbon-help-filled />
            </n-icon>
          </template>
          开启自定义回复调试日志，打印字符细节
        </n-tooltip>
      </template>
      <n-checkbox v-model:checked="replyDebugMode">开启</n-checkbox>
    </n-form-item>

    <h3>跑团日志</h3>
    <n-form-item label="自定义后端 URL">
      <template #label>
        <span>自定义后端 URL</span>
        <n-tooltip>
          <template #trigger>
            <n-icon>
              <i-carbon-help-filled />
            </n-icon>
          </template>
          设置第三方跑团日志后端 URL
        </n-tooltip>
      </template>
      <n-input v-model:value="config.storyLogBackendUrl" placeholder="" style="width: 30rem" />
    </n-form-item>
    <n-form-item label="API 版本">
      <template #label>
        <span>API 版本</span>
        <n-tooltip>
          <template #trigger>
            <n-icon>
              <i-carbon-help-filled />
            </n-icon>
          </template>
          指定后端的 API 版本
        </n-tooltip>
      </template>
      <n-input v-model:value="config.storyLogApiVersion" placeholder="" style="width: 10rem" />
    </n-form-item>
    <n-form-item label="Token">
      <template #label>
        <span>Token</span>
        <n-tooltip>
          <template #trigger>
            <n-icon>
              <i-carbon-help-filled />
            </n-icon>
          </template>
          指定传递给后端的 token
        </n-tooltip>
      </template>
      <n-input v-model:value="config.storyLogBackendToken" placeholder="" style="width: 30rem" />
    </n-form-item>

    <n-form-item v-if="modified" label="" label-width="1rem" class="mt-4">
      <n-flex>
        <n-button type="error" @click="submitGiveup">放弃改动</n-button>
        <n-button type="success" @click="submit">保存设置</n-button>
      </n-flex>
    </n-form-item>
  </n-form>
</template>

<script setup lang="ts">
import { useStore } from '~/store';
import type { AdvancedConfig } from '#';
import { getCustomReplyDebug, postCustomReplyDebug } from '~/api/v1/configs';

const emit = defineEmits(['update:advanced-settings-show']);

const store = useStore();

const config = ref<AdvancedConfig>({
  show: false,
  enable: false,
  storyLogBackendUrl: '',
  storyLogApiVersion: '',
  storyLogBackendToken: '',
});
const replyDebugMode = ref(false);

onBeforeMount(async () => {
  config.value = await store.diceAdvancedConfigGet();
  replyDebugMode.value = (await getCustomReplyDebug()).value;
  nextTick(() => {
    modified.value = false;
  });
});

const modified = ref(false);
watch(
  () => config,
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  (newValue, oldValue) => {
    //直接监听
    modified.value = true;
  },
  {
    deep: true,
  },
);
watch(
  () => replyDebugMode.value,
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  (newValue, oldValue) => {
    //直接监听
    modified.value = true;
  },
);

const submit = async () => {
  await store.diceAdvancedConfigSet(config.value);
  await postCustomReplyDebug(replyDebugMode.value);
  config.value = await store.diceAdvancedConfigGet();
  modified.value = false;
  emit('update:advanced-settings-show', config.value.show);
  nextTick(async () => {
    modified.value = false;
  });
};

const submitGiveup = async () => {
  config.value = await store.diceAdvancedConfigGet();
  replyDebugMode.value = (await getCustomReplyDebug()).value;
  modified.value = false;
  nextTick(() => {
    modified.value = false;
  });
};
</script>

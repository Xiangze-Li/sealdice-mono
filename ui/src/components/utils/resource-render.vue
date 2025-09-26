<script setup lang="ts">
import { useStore, type Resource } from '~/store';
import { getResourceData } from '~/api/v1/resource';
import { urlBase } from '~/backend';
const store = useStore();

const url = ref<string>('');

const getImageUrl = async (path: string) => {
  const blob = await getResourceData(path, props.mini);
  return window.URL.createObjectURL(blob);
};

onMounted(async () => {
  url.value = await getImageUrl(props.data.path);
});

const props = withDefaults(
  defineProps<{
    data: Resource;
    mini?: boolean;
  }>(),
  {
    mini: false,
  },
);
</script>

<template>
  <template v-if="data.type === 'image'">
    <n-image
      :key="url"
      :alt="data.name"
      :src="url"
      :preview-src="`${urlBase}/sd-api/resource/download?path=${encodeURIComponent(data.path)}&token=${encodeURIComponent(store.token)}`"
      object-fit="contain"
      lazy>
      <template #placeholder>
        <div>
          <n-icon>
            <i-carbon-image />
          </n-icon>
        </div>
      </template>
      <template #error>
        <div>
          <n-icon>
            <i-carbon-image />
          </n-icon>
        </div>
      </template>
    </n-image>
  </template>
  <template v-else-if="data.type === 'audio'">
    <n-icon>
      <i-carbon-microphone />
    </n-icon>
  </template>
  <template v-else-if="data.type === 'video'">
    <n-icon>
      <i-carbon-video />
    </n-icon>
  </template>
  <template v-else>
    {{ '未知格式' }}
  </template>
</template>

<style scoped lang="css"></style>

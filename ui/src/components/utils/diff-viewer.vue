<template>
  <n-flex
    style="margin-left: 1rem; display: flex; justify-content: space-between; align-items: center">
    <n-text v-if="changed" type="info" class="flex items-center gap-x-1">
      <n-icon><i-carbon-help-filled /></n-icon>
      变更如下：
    </n-text>
    <n-text v-else type="info" size="large">
      <n-icon><i-carbon-information-filled /></n-icon>
      无变更
    </n-text>
    <n-flex v-if="changed" vertical align="center" wrap>
      <n-switch v-model:value="split">
        <template #checked>双列</template>
        <template #unchecked>单列</template>
      </n-switch>
      <n-checkbox v-model:checked="folding">折叠无变更</n-checkbox>
    </n-flex>
  </n-flex>
  <div v-show="split" style="display: flex; justify-content: space-around; align-items: center">
    <h3 style="padding-left: 2rem">原内容</h3>
    <n-icon>
      <i-carbon-chevron-right />
    </n-icon>
    <h3 style="padding-right: 2rem">新内容</h3>
  </div>
  <vue-diff
    v-if="changed"
    :mode="mode"
    theme="light"
    :language="props.lang"
    :folding="folding"
    :prev="props.old"
    :current="props.new" />
</template>

<script lang="ts" setup>
interface Props {
  old: string;
  new: string;
  lang?: string;
}

const props = withDefaults(defineProps<Props>(), {
  lang: 'text',
  old: '',
  new: '',
});

const changed = computed(() => !(props.old === props.new));
const mode = computed(() => (split.value ? 'split' : 'unified'));
const split = ref<boolean>(false);
const folding = ref<boolean>(false);
</script>

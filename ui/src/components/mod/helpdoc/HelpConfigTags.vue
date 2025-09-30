<template>
  <n-flex align="center" ref="tagsRef" wrap>
    <n-tag type="success" size="small" :bordered="false">{{ group.value }}</n-tag>
    <n-tag
      size="small"
      v-for="a in aliases.get(group.value)"
      :key="a"
      closable
      :bordered="false"
      @close="handleClose(group.value, a)">
      {{ a }}
    </n-tag>

    <n-input
      v-if="inputVisible"
      ref="inputRef"
      v-model:value="inputValue"
      size="tiny"
      autosize
      style="min-width: 6rem"
      placeholder=""
      @keyup.enter="handleInputConfirm(group.value)"
      @blur="handleInputConfirm(group.value)" />
    <n-button v-if="inputVisible" size="tiny" tertiary @click="handleInputConfirm(group.value)">
      确定
    </n-button>
    <n-button v-else size="tiny" tertiary @click="showInput">
      <template #icon>
        <n-icon><i-carbon-add-large /></n-icon>
      </template>
      新别名
    </n-button>
  </n-flex>
</template>

<script setup lang="ts">
import type { ElInput, ElSpace } from 'element-plus';

const { group, aliases } = defineProps<{
  group: {
    value: string;
    label: string;
  };
  aliases: Map<string, string[]>;
}>();
const emit = defineEmits(['addAlias', 'removeAlias']);

const tagsRef = ref<InstanceType<typeof ElSpace>>();

const inputVisible = ref<boolean>(false);
const inputRef = ref<InstanceType<typeof ElInput>>();
const inputValue = ref('');
const showInput = () => {
  inputVisible.value = true;
  nextTick(() => {
    inputRef.value!.input!.focus();
  });
};

const handleInputConfirm = (groupKey: string) => {
  if (inputValue.value) {
    emit('addAlias', groupKey, inputValue.value);
  }
  inputVisible.value = false;
  inputValue.value = '';
};

const handleClose = (groupKey: string, tag: string) => {
  emit('removeAlias', groupKey, tag);
};
</script>

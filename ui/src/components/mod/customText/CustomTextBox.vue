<template>
  <template v-if="group !== '__others__'">
    <n-collapse-item :title="group" :name="group">
      <slot name="values" />
    </n-collapse-item>
  </template>
  <template v-else>
    <n-collapse-item title="无分组" :name="group">
      <slot name="values" />
    </n-collapse-item>
  </template>
</template>

<script setup lang="ts">
defineProps<{
  group: string | '__others__';
}>();
</script>

<style scoped lang="css">
.text-group-title {
  border-top: 1px solid var(--el-border-color);
}
</style>

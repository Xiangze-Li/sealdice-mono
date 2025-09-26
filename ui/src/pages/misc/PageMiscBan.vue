<script setup lang="ts">
const mode = ref<string>('list');
</script>

<template>
  <n-tabs v-model:value="mode" justify-content="space-evenly">
    <n-tab-pane tab="黑白名单" name="list">
      <ban-list />
    </n-tab-pane>
    <n-tab-pane tab="拉黑设置" name="config">
      <ban-config />
    </n-tab-pane>
  </n-tabs>
</template>

<style scoped lang="css"></style>

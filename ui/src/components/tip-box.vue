<script setup lang="ts">
import { useThemeVars } from 'naive-ui';
import tinycolor from 'tinycolor2';

interface Props {
  type?: 'default' | 'success' | 'info' | 'warning' | 'error';
}
const { type = 'default' } = defineProps<Props>();

const themeVars = useThemeVars();
const borderColor = computed(() => {
  switch (type) {
    case 'success':
      return themeVars.value.successColor;
    case 'info':
      return themeVars.value.infoColor;
    case 'warning':
      return themeVars.value.warningColor;
    case 'error':
      return themeVars.value.errorColor;
    default:
      return themeVars.value.primaryColor;
  }
});
const bgColor = computed(() => {
  return tinycolor(borderColor.value).setAlpha(0.06).toRgbString();
});
</script>

<template>
  <div class="tip-box rounded-sm border-l-[3px] p-4">
    <slot :type="type" />
  </div>
</template>

<style scoped lang="css">
.tip-box {
  border-color: v-bind('borderColor');
  background-color: v-bind('bgColor');
}
</style>

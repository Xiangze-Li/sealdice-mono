<script setup lang="ts">
import { breakpointsTailwind } from '@vueuse/core';

interface ReplyCondition {
  condType: string;
  value: string | number | undefined;
  matchType: string;
  matchOp?: string;
}

const listModel = defineModel<ReplyCondition[]>();

const breakpoints = useBreakpoints(breakpointsTailwind);
const notMobile = breakpoints.greater('sm');

const deleteByIndex = (index: number) => {
  listModel.value?.splice(index, 1);
};
</script>

<template>
  <div
    v-for="(cond, index) in listModel"
    :key="index"
    class="mb-3 border-l-2 border-orange-500 pl-2">
    <div class="flex justify-between border-b pb-2">
      <n-flex align="center">
        <n-text>模式</n-text>
        <n-radio-group v-model:value="cond.condType" size="small">
          <n-radio-button value="textMatch" label="文本匹配" />
          <n-radio-button value="textLenLimit" label="文本长度" />
          <n-radio-button value="exprTrue" label="表达式为真" />
        </n-radio-group>
      </n-flex>
      <n-button type="error" size="tiny" secondary @click="deleteByIndex(index)">
        <template #icon>
          <n-icon><i-carbon-row-delete /></n-icon>
        </template>
        <template v-if="notMobile" #default>删除条件</template>
      </n-button>
    </div>

    <div v-if="cond.condType === 'textMatch'" style="display: flex" class="mobile-changeline">
      <div style="width: 7rem; margin-right: 0.5rem">
        <n-text>
          方式
          <n-tooltip>
            <template #trigger>
              <n-icon><i-carbon-help-filled /> </n-icon>
            </template>
            匹配方式一览:<br />
            精确匹配: 完全相同时触发。<br />
            任意相符: 如aa|bb，则aa或bb都能触发。<br />
            包含文本: 包含此文本触发。<br />
            不含文本: 不包含此文本触发。<br />
            模糊匹配: 文本相似时触发<br />
            正则匹配: 正则表达式匹配，语法请自行查阅<br />
            前缀匹配: 文本以内容为开头<br />
            后缀匹配: 文本以此内容为结尾
          </n-tooltip>
        </n-text>
        <n-select
          v-model:value="cond.matchType"
          size="small"
          placeholder="选择方式"
          :options="[
            { label: '精确匹配', value: 'matchExact' },
            { label: '任意相符', value: 'matchMulti' },
            { label: '包含文本', value: 'matchContains' },
            { label: '不含文本', value: 'matchNotContains' },
            { label: '模糊匹配', value: 'matchFuzzy' },
            { label: '正则匹配', value: 'matchRegex' },
            { label: '前缀匹配', value: 'matchPrefix' },
            { label: '后缀匹配', value: 'matchSuffix' },
          ]" />
      </div>

      <div style="flex: 1">
        <n-text>内容</n-text>
        <n-input v-model:value="cond.value as string" />
      </div>
    </div>

    <div v-else-if="cond.condType === 'exprTrue'" style="display: flex" class="mobile-changeline">
      <div style="flex: 1">
        <n-text>
          表达式
          <n-tooltip>
            <template #trigger>
              <n-icon><i-carbon-help-filled /> </n-icon>
            </template>
            举例：<br />
            $t1 == '张三' // 正则匹配的第一个组内容是张三<br />
            $m个人计数器 >= 10<br />
            友情提醒，匹配失败时无提示，请先自行在“指令测试”测好
          </n-tooltip>
        </n-text>
        <n-input
          v-model:value="cond.value as string"
          type="textarea"
          :autosize="{ minRows: 1, maxRows: 10 }" />
      </div>
    </div>

    <n-flex v-else-if="cond.condType === 'textLenLimit'" class="mt-2 flex">
      <n-radio-group v-model:value="cond.matchOp" size="small">
        <n-radio-button value="ge" label="大于等于" />
        <n-radio-button value="le" label="小于等于" />
      </n-radio-group>
      <n-input-number v-model:value="cond.value as number" :min="0" />
    </n-flex>
  </div>
</template>

<style scoped lang="css">
@media screen and (max-width: 700px) {
  .mobile-changeline {
    flex-direction: column;
  }
}
</style>

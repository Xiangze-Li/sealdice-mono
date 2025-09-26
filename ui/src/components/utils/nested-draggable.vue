<!-- eslint-disable vue/multi-word-component-names -->
<template>
  <draggable
    class="dragArea"
    tag="div"
    :list="tasks"
    handle=".handle"
    :group="{ name: 'g1' }"
    item-key="name">
    <template #item="{ element: el, index }">
      <li class="reply-item mb-2 list-none">
        <foldable-card type="div" :default-fold="true" compact>
          <template #title>
            <n-checkbox v-model:checked="el.enable">开启</n-checkbox>
          </template>
          <template #title-extra>
            <n-flex size="large" align="center">
              <n-icon class="handle">
                <i-carbon-move />
              </n-icon>
              <n-button type="error" size="small" secondary @click="deleteItem(index)">
                <template #icon>
                  <n-icon><i-carbon-row-delete /></n-icon>
                </template>
                删除
              </n-button>
            </n-flex>
          </template>

          <template #unfolded-extra>
            <div class="border-l-4 border-orange-500 pl-4">
              <div v-for="(cond, index2) in el.conditions || []" :key="index2">
                <n-text
                  v-if="cond.condType === 'textMatch'"
                  style="display: flex"
                  class="mobile-changeline text-base">
                  文本匹配：{{ cond.value }}
                </n-text>
                <n-text
                  v-else-if="cond.condType === 'exprTrue'"
                  style="display: flex"
                  class="mobile-changeline text-base">
                  <div style="flex: 1">表达式：{{ cond.value }}</div>
                </n-text>
                <n-text
                  v-else-if="cond.condType === 'textLenLimit'"
                  style="display: flex"
                  class="mobile-changeline">
                  <div style="flex: 1">
                    长度：
                    {{ cond.matchOp === 'ge' ? '大于等于' : '' }}
                    {{ cond.matchOp === 'le' ? '小于等于' : '' }}
                    {{ cond.value }}
                  </div>
                </n-text>
              </div>
            </div>
          </template>

          <n-text class="mb-2 block text-base">条件（需同时满足，即 and）</n-text>
          <div class="border-l-4 border-orange-500 pl-4">
            <custom-reply-conditions v-model="el.conditions" />
            <n-button type="success" size="small" @click="addCond(el.conditions)">
              <template #icon>
                <n-icon><i-carbon-add-large /></n-icon>
              </template>
              增加
            </n-button>
          </div>

          <n-text class="my-2 block text-base">结果（顺序执行）</n-text>
          <div class="border-l-4 border-blue-500 pl-4">
            <div
              v-for="(i, index) in el.results || []"
              :key="index"
              class="mb-3 border-l-2 border-blue-500 pl-2">
              <div style="display: flex; justify-content: space-between">
                <n-flex align="center">
                  <n-text>模式</n-text>
                  <n-radio-group v-model:value="i.resultType" size="small">
                    <n-radio-button value="replyToSender" label="回复" />
                    <n-radio-button value="replyPrivate" label="私聊回复" />
                    <n-radio-button value="replyGroup" label="群内回复" />
                  </n-radio-group>
                </n-flex>

                <n-button
                  type="error"
                  size="tiny"
                  secondary
                  @click="deleteAnyItem(el.results, index)">
                  <template #icon>
                    <n-icon><i-carbon-row-delete /></n-icon>
                  </template>
                  <template v-if="notMobile" #default>删除结果</template>
                </n-button>
              </div>

              <div v-if="['replyToSender', 'replyPrivate', 'replyGroup'].includes(i.resultType)">
                <div class="mobile-changeline my-2 flex justify-between">
                  <div style="display: flex; align-items: center">
                    <n-text>回复文本（随机选择）</n-text>
                  </div>
                  <n-flex align="center">
                    <n-text>
                      延迟
                      <n-tooltip>
                        <template #trigger>
                          <n-icon><i-carbon-help-filled /></n-icon>
                        </template>
                        文本将在此延迟后发送，单位秒，可小数。<br />
                        注意随机延迟仍会被加入，如果你希望保证发言顺序，记得考虑这点。
                      </n-tooltip>
                    </n-text>
                    <n-input-number v-model:value="i.delay" size="small" class="w-20" />
                  </n-flex>
                </div>

                <div v-for="(k2, index) in i.message" :key="index" class="my-2 w-full">
                  <!-- 这里面是单条修改项 -->
                  <div style="display: flex">
                    <div
                      style="
                        display: flex;
                        align-items: center;
                        width: 1.3rem;
                        margin-left: 0.2rem;
                      ">
                      <n-tooltip placement="bottom-start">
                        <template #trigger>
                          <n-icon>
                            <i-carbon-add-filled v-if="index == 0" @click="addItem(i.message)" />
                            <i-carbon-close-outline v-else @click="removeItem(i.message, index)" />
                          </n-icon>
                        </template>
                        {{
                          index === 0
                            ? '点击添加一个回复语，海豹将会随机抽取一个回复'
                            : '点击删除你不想要的回复语'
                        }}
                      </n-tooltip>
                    </div>
                    <div style="flex: 1">
                      <n-input v-model:value="k2[0]" type="textarea" class="reply-text" autosize />
                    </div>
                  </div>
                </div>
              </div>
            </div>
            <n-button type="success" size="small" @click="addResult(el.results)">
              <template #icon>
                <n-icon><i-carbon-add-large /> </n-icon>
              </template>
              增加
            </n-button>
          </div>
        </foldable-card>
      </li>
    </template>
  </draggable>
</template>

<script lang="ts" setup>
import draggable from 'vuedraggable';
import { breakpointsTailwind } from '@vueuse/core';
import { useThemeVars } from 'naive-ui';

const themeVars = useThemeVars();

const props = defineProps<{ tasks: Array<any> }>();

const breakpoints = useBreakpoints(breakpointsTailwind);
const notMobile = breakpoints.greater('sm');

const deleteItem = (index: number) => {
  ElMessageBox.confirm('确认删除此项？', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  }).then(async () => {
    // eslint-disable-next-line vue/no-mutating-props
    props.tasks.splice(index, 1);
  });
};

const deleteAnyItem = (lst: any[], index: number) => {
  lst.splice(index, 1);
};

const addCond = (i: any) => {
  i.push({
    condType: 'textMatch',
    matchType: 'matchExact',
    value: '要匹配的文本',
  });
};

const addResult = (i: any) => {
  i.push({ resultType: 'replyToSender', delay: 0, message: [['说点什么', 1]] });
};

const addItem = (k: any) => {
  k.push(['怎么辉石呢', 1]);
};

const removeItem = (v: any[], index: number | any) => {
  v.splice(index, 1);
};
</script>

<style scoped lang="css">
.dragArea {
  min-height: 50px;
  /* outline: 1px dashed; */
  padding-top: 1rem;
  padding-bottom: 1rem;

  .reply-item:not(:last-child) {
    border-bottom: 1px solid v-bind(themeVars.borderColor);
    padding-bottom: 1rem;
  }
}

@media screen and (max-width: 700px) {
  .mobile-changeline {
    flex-direction: column;
  }
}
</style>

<script setup lang="ts">
import type { JsPluginConfigItem } from '#';
import { postUtilsCheckCronExpr } from '~/api/v1/utils';
import { isEqual, startsWith } from 'es-toolkit/compat';

const model = defineModel<JsPluginConfigItem>({ required: true });

defineEmits(['update:value', 'update:template-add', 'update:template-remove', 'reset', 'delete']);

const checkCronExpr = (expr: string) => {
  try {
    postUtilsCheckCronExpr(expr);
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
  } catch (_err) {
    return false;
  }
  return true;
};

const checkDailyExpr = (expr: string) => {
  const pattern = /^([0-1]?[0-9]|2[0-3]):([0-5][0-9])$/;
  return pattern.test(expr);
};
</script>

<template>
  <n-flex vertical>
    <n-flex align="center" justify="space-between">
      <n-text v-if="model.type === 'string'">字符串配置项「{{ model.key }}」</n-text>
      <n-text v-if="model.type === 'int'">整数配置项「{{ model.key }}」</n-text>
      <n-text v-if="model.type === 'float'">浮点数配置项「{{ model.key }}」</n-text>
      <n-text v-if="model.type === 'bool'">布尔配置项「{{ model.key }}」</n-text>
      <n-text v-if="model.type === 'option'">选项配置项「{{ model.key }}」</n-text>
      <n-text v-if="model.type === 'task:cron'">Cron 型定时任务「{{ model.key }}」</n-text>
      <n-text v-if="model.type === 'task:daily'">每日定时任务「{{ model.key }}」</n-text>
      <n-text v-if="model.type === 'template'">模板配置项「{{ model.key }}」</n-text>

      <n-text type="info" class="text-xs">{{ model.description }}</n-text>
    </n-flex>

    <div class="my-1">
      <n-input
        v-if="model.type === 'string'"
        v-model:value="model.value"
        type="textarea"
        @update:value="$emit('update:value')" />
      <n-input-number
        v-else-if="model.type === 'int'"
        v-model:value="model.value"
        :precision="0"
        @update:value="$emit('update:value')" />
      <n-input-number
        v-else-if="model.type === 'float'"
        v-model:value="model.value"
        @update:value="$emit('update:value')" />
      <n-switch
        v-else-if="model.type === 'bool'"
        v-model:value="model.value"
        @update:value="$emit('update:value')" />
      <n-radio-group
        v-else-if="model.type === 'option'"
        v-model:value="model.value"
        @update:value="$emit('update:value')">
        <n-radio v-for="s in model.option" :key="s" :value="s" />
      </n-radio-group>
      <n-input
        v-else-if="model.type === 'task:cron'"
        v-model:value="model.value"
        type="textarea"
        :allow-input="checkCronExpr"
        @update:value="$emit('update:value')" />
      <n-input
        v-else-if="model.type === 'task:daily'"
        v-model:value="model.value"
        type="textarea"
        :allow-input="checkDailyExpr"
        @update:value="$emit('update:value')" />
      <div v-else-if="model.type === 'template'">
        <div v-for="(_, index) in model.value" :key="index">
          <!-- 这里面是单条修改项 -->
          <n-row>
            <n-col>
              <n-input
                v-model:value="model.value[index]"
                type="textarea"
                :autosize="true"
                @update:value="$emit('update:value')" />
            </n-col>
            <n-col :span="5">
              <div>
                <n-tooltip placement="bottom-start">
                  <template #trigger>
                    <n-button text v-if="index == 0" @click="$emit('update:template-add')">
                      <template #icon>
                        <n-icon><i-carbon-add-filled /></n-icon>
                      </template>
                    </n-button>
                    <n-button text v-else @click="$emit('update:template-remove', index)">
                      <template #icon>
                        <n-icon><i-carbon-close-outline /></n-icon>
                      </template>
                    </n-button>
                  </template>
                  {{ index === 0 ? '点击添加一项' : '点击删除你不想要的配置项' }}
                </n-tooltip>
              </div>
            </n-col>
          </n-row>
        </div>
      </div>
    </div>

    <n-flex align="center" class="ml-auto">
      <template v-if="!isEqual(model.value, model.defaultValue)">
        <n-tooltip placement="bottom-end">
          <template #trigger>
            <n-button text @click="$emit('reset')">
              <template #icon>
                <n-icon><i-carbon-paint-brush /></n-icon>
              </template>
            </n-button>
          </template>
          重置为初始值
        </n-tooltip>
      </template>
      <template v-if="model.deprecated">
        <n-tooltip placement="bottom-end">
          <template #trigger>
            <n-button text @click="$emit('delete')">
              <template #icon>
                <n-icon><i-carbon-row-delete /></n-icon>
              </template>
            </n-button>
          </template>

          <n-text v-if="startsWith(model.type, 'cron:')">
            移除 - 这个定时任务在当前版本中不被使用，<br />但未来版本中仍可能被使用，请确认无用后删除
          </n-text>
          <n-text v-else>
            移除 - 这个配置在新版的默认配置中不被使用，<br />但升级而来时仍可能被使用，请确认无用后删除
          </n-text>
        </n-tooltip>
      </template>
    </n-flex>
  </n-flex>
</template>

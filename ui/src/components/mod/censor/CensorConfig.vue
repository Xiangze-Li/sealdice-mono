<template>
  <header class="flex items-center">
    <n-button type="info" secondary @click="submit">
      <template #icon>
        <i-carbon-save />
      </template>
      保存设置
    </n-button>
    <n-text v-if="modified" type="error" tag="strong" class="ml-4 text-base">
      内容已修改，不要忘记保存！
    </n-text>
  </header>
  <n-form label-placement="left" label-width="auto">
    <h4>匹配选项</h4>
    <n-form-item>
      <template #label>
        <n-text>拦截范围</n-text>
        <n-tooltip raw-content content="">
          <template #trigger>
            <n-icon>
              <i-carbon-help-filled />
            </n-icon>
          </template>
          发出的消息： 拦截骰子发出的内容，进行检查。未通过检查，替换为
          <n-tag size="small" type="info" :bordered="false">拦截_完全拦截_发出的消息</n-tag>
          的内容。<br />
          收到的指令： 拦截骰子收到的命令文本进行检查，如收到「.rd
          进行一次骰点」时，会检查其中的「进行一次骰点」，未通过检查则发送
          <n-tag size="small" type="info" :bordered="false">拦截_完全拦截_收到的指令</n-tag>
          的内容<br />
          收到的所有消息： 会对所有收到的消息(所有群内聊天)进行检查，未通过检查默认不做响应，如
          <n-tag size="small" type="info" :bordered="false">拦截_完全拦截_收到的所有消息</n-tag>
          不为空时会发送拦截提示。
        </n-tooltip>
      </template>
      <n-radio-group v-model:value="config.mode" size="small">
        <n-radio :value="Mode.ReplyOutput">{{ '发出的消息' }}</n-radio>
        <n-radio :value="Mode.CommandInput">{{ '收到的指令' }}</n-radio>
        <n-radio :value="Mode.AllInput">{{ '收到的所有消息(慎用)' }}</n-radio>
      </n-radio-group>
    </n-form-item>
    <n-form-item label="大小写敏感">
      <n-checkbox v-model:checked="config.caseSensitive" label="开启" />
    </n-form-item>
    <n-form-item>
      <template #label>
        <n-text>匹配拼音</n-text>
        <n-tooltip>
          <template #trigger>
            <n-icon>
              <i-carbon-help-filled />
            </n-icon>
          </template>
          匹配敏感词拼音，勾选大小写敏感时该项无效。
        </n-tooltip>
      </template>
      <n-checkbox v-model:checked="config.matchPinyin" label="开启" />
    </n-form-item>
    <n-form-item>
      <template #label>
        <n-text>过滤字符正则</n-text>
        <n-tooltip>
          <template #trigger>
            <n-icon>
              <i-carbon-help-filled />
            </n-icon>
          </template>
          <!-- eslint-disable-next-line prettier/prettier-->
          判断敏感词时，忽略过滤字符。如敏感词为"114514"，指定过滤字符为空白，则"114&nbsp;&nbsp;&nbsp;514"也会命中敏感词。
        </n-tooltip>
      </template>
      <n-input v-model:value="config.filterRegex" placeholder="" style="width: 12rem" />
    </n-form-item>
  </n-form>
  <h4>响应设置</h4>
  <tip-box type="warning" class="my-4">
    <n-text type="warning">
      <span>提示：</span>
      <ul class="ml-4 list-disc">
        <li><p>超过阈值时，对应用户该等级的计数会被清空重新计算。</p></li>
        <li>
          <p>
            增加怒气值时，会计算群组和邀请人的连带责任。连带责任比例在
            <strong>综合设置 > 黑白名单 > 设置选项</strong> 中调整。
          </p>
        </li>
      </ul>
    </n-text>
  </tip-box>
  <n-form label-placement="left" label-width="auto">
    <n-form-item>
      <template #label>
        <sensitive-tag type="default" />
      </template>
      <n-flex align="center">
        <n-flex align="center" wrap>
          <n-text>用户触发超过</n-text>
          <n-input-number
            v-model:value="config.levelConfig.notice.threshold"
            class="w-28"
            size="small"
            :step="1"
            :min="0"
            :precision="0" />
          <n-text>次时：</n-text>
        </n-flex>
        <n-flex justify="center" vertical wrap>
          <n-checkbox-group v-model:value="config.levelConfig.notice.handlers">
            <n-checkbox v-for="handle in defaultHandles" :key="handle.key" :value="handle.key">
              {{ handle.name }}
            </n-checkbox>
          </n-checkbox-group>
          <n-flex align="center">
            <n-text>怒气值</n-text>
            <n-input-number
              v-model:value="config.levelConfig.notice.score"
              class="w-28"
              size="small"
              :step="1"
              :min="0"
              :precision="0" />
          </n-flex>
        </n-flex>
      </n-flex>
    </n-form-item>
    <n-form-item>
      <template #label>
        <sensitive-tag type="info" />
      </template>
      <n-flex align="center">
        <n-flex align="center" wrap>
          <n-text>用户触发超过</n-text>
          <n-input-number
            v-model:value="config.levelConfig.caution.threshold"
            class="w-28"
            size="small"
            :step="1"
            :min="0"
            :precision="0" />
          <n-text>次时：</n-text>
        </n-flex>
        <n-flex justify="center" vertical wrap>
          <n-checkbox-group v-model:value="config.levelConfig.caution.handlers">
            <n-checkbox v-for="handle in defaultHandles" :key="handle.key" :value="handle.key">
              {{ handle.name }}
            </n-checkbox>
          </n-checkbox-group>
          <n-flex align="center">
            <n-text>怒气值</n-text>
            <n-input-number
              v-model:value="config.levelConfig.caution.score"
              class="w-28"
              size="small"
              :step="1"
              :min="0"
              :precision="0" />
          </n-flex>
        </n-flex>
      </n-flex>
    </n-form-item>
    <n-form-item>
      <template #label>
        <sensitive-tag type="warning" />
      </template>
      <n-flex align="center">
        <n-flex align="center" wrap>
          <n-text>用户触发超过</n-text>
          <n-input-number
            v-model:value="config.levelConfig.warning.threshold"
            class="w-28"
            size="small"
            :step="1"
            :min="0"
            :precision="0" />
          <n-text>次时：</n-text>
        </n-flex>
        <n-flex justify="center" vertical wrap>
          <n-checkbox-group v-model:value="config.levelConfig.warning.handlers">
            <n-checkbox v-for="handle in defaultHandles" :key="handle.key" :value="handle.key">
              {{ handle.name }}
            </n-checkbox>
          </n-checkbox-group>
          <n-flex align="center">
            <n-text>怒气值</n-text>
            <n-input-number
              v-model:value="config.levelConfig.warning.score"
              class="w-28"
              size="small"
              :step="1"
              :min="0"
              :precision="0" />
          </n-flex>
        </n-flex>
      </n-flex>
    </n-form-item>
    <n-form-item>
      <template #label>
        <sensitive-tag type="error" />
      </template>
      <n-flex align="center">
        <n-flex align="center" wrap>
          <n-text>用户触发超过</n-text>
          <n-input-number
            v-model:value="config.levelConfig.danger.threshold"
            class="w-28"
            size="small"
            :step="1"
            :min="0"
            :precision="0" />
          <n-text>次时：</n-text>
        </n-flex>
        <n-flex justify="center" vertical wrap>
          <n-checkbox-group v-model:value="config.levelConfig.danger.handlers">
            <n-checkbox v-for="handle in defaultHandles" :key="handle.key" :value="handle.key">
              {{ handle.name }}
            </n-checkbox>
          </n-checkbox-group>
          <n-flex align="center">
            <n-text>怒气值</n-text>
            <n-input-number
              v-model:value="config.levelConfig.danger.score"
              class="w-28"
              size="small"
              :step="1"
              :min="0"
              :precision="0" />
          </n-flex>
        </n-flex>
      </n-flex>
    </n-form-item>
  </n-form>
</template>

<script lang="ts" setup>
import { isArray, isEqual, isObject, transform } from 'es-toolkit/compat';
import { useMessage } from 'naive-ui';
import { getCensorConfig, postCensorConfig } from '~/api/v1/censor';
import { useCensorStore } from '~/components/mod/censor/censor';

const message = useMessage();

onBeforeMount(async () => {
  await refreshCensorConfig();
  await nextTick(() => {
    modified.value = false;
  });
});

onBeforeUnmount(() => {
  clearInterval(timerId);
});

const censorStore = useCensorStore();

const enum Mode {
  ReplyOutput = 0,
  CommandInput = 1,
  AllInput = 2,
}

interface Config {
  mode: Mode;
  caseSensitive: boolean;
  matchPinyin: boolean;
  filterRegex: string;
  levelConfig: LevelConfigs;
}

const config = ref<Config>({
  mode: Mode.AllInput,
  caseSensitive: false,
  matchPinyin: false,
  filterRegex: '',
  levelConfig: {
    notice: { threshold: 0, handlers: [], score: 0 },
    caution: { threshold: 0, handlers: [], score: 0 },
    warning: { threshold: 0, handlers: [], score: 0 },
    danger: { threshold: 0, handlers: [], score: 0 },
  },
});

interface LevelConfigs {
  notice: LevelConfig;
  caution: LevelConfig;
  warning: LevelConfig;
  danger: LevelConfig;
}

interface LevelConfig {
  threshold: number;
  handlers: string[];
  score: number;
}

const defaultHandles: { key: string; name: string }[] = [
  { key: 'SendWarning', name: '发送警告' },
  { key: 'SendNotice', name: '通知 Master' },
  { key: 'BanUser', name: '拉黑用户' },
  { key: 'BanGroup', name: '拉黑群' },
  { key: 'BanInviter', name: '拉黑邀请人' },
  { key: 'AddScore', name: '增加怒气值' },
];

const modified = ref<boolean>(false);

watch(
  config,
  () => {
    modified.value = true;
  },
  { deep: true },
);

censorStore.$subscribe(async (_, state) => {
  if (state.settingsNeedRefresh === true) {
    await refreshCensorConfig();
    state.settingsNeedRefresh = false;
  }
});

const timerId: number = 0;
const refreshCensorConfig = async () => {
  const c = await getCensorConfig();
  if (c.result) {
    config.value = c;
  }
  modified.value = false;
  await nextTick(() => {
    modified.value = false;
  });
};

const confDiff = (object: any, base: any) => {
  const changes = function (object: any, base: any) {
    return transform(object, (result: any, value, key) => {
      if (isArray(value)) {
        result[key] = value;
      } else if (!isEqual(value, base[key])) {
        result[key] = isObject(value) && isObject(base[key]) ? changes(value, base[key]) : value;
      }
    });
  };
  return changes(object, base);
};

const submit = async () => {
  const conf = await getCensorConfig();
  const modify = confDiff(config.value, conf);

  const resp = await postCensorConfig(modify);
  if (resp.result) {
    message.success('保存设置成功');
  } else {
    message.error('保存设置失败，' + resp.err);
  }

  await refreshCensorConfig();
  censorStore.markReload();
  modified.value = false;
  await nextTick(() => {
    modified.value = false;
  });
};
</script>

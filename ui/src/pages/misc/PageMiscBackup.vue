<template>
  <div style="display: flex; justify-content: space-between; align-items: center">
    <h2>备份</h2>
    <n-flex>
      <n-button type="success" @click="doSave">
        <template #icon>
          <i-carbon-save />
        </template>
        保存设置
      </n-button>
      <n-button type="info" @click="showBackup = true">立即备份</n-button>
    </n-flex>
  </div>
  <div>
    <n-form label-placement="left" label-width="5.5rem" label-align="left">
      <h3>自动备份</h3>
      <n-checkbox v-model:checked="cfg.autoBackupEnable">开启</n-checkbox>
      <div v-if="cfg.autoBackupEnable" style="margin-top: 1rem">
        <n-form-item>
          <template #label>
            <span>备份间隔</span>
            <n-tooltip>
              <template #trigger>
                <n-icon><i-carbon-help-filled /></n-icon>
              </template>
              备份间隔表达式请参阅
              <a href="https://pkg.go.dev/github.com/robfig/cron" target="_blank">cron 文档</a>
            </n-tooltip>
          </template>
          <n-input v-model:value="cfg.autoBackupTime" style="width: 12rem" />
        </n-form-item>
        <n-form-item label="备份范围">
          <n-checkbox-group v-model:value="cfg.autoBackupSelectionList">
            <n-checkbox label="基础（含自定义回复）" value="base" checked disabled />
            <n-checkbox label="JS 插件" value="js" />
            <n-checkbox label="牌堆" value="deck" />
            <n-checkbox label="帮助文档" value="helpdoc" />
            <n-checkbox label="敏感词库" value="censor" />
            <n-checkbox label="人名信息" value="name" />
            <n-checkbox label="图片" value="image" />
          </n-checkbox-group>
        </n-form-item>
        <n-form-item label="备份文件名预览">
          <n-text type="info">
            bak_{{ now }}_auto_r{{ cfg.autoBackupSelection.toString(16) }}_&lt;随机值&gt;.zip
          </n-text>
        </n-form-item>
      </div>
      <h3>自动清理</h3>
      <n-form-item label="清理模式">
        <n-radio-group v-model:value="cfg.backupCleanStrategy" size="small">
          <n-radio-button :value="0">关闭</n-radio-button>
          <n-radio-button :value="1">保留一定数量</n-radio-button>
          <n-radio-button :value="2">保留一定时间内</n-radio-button>
        </n-radio-group>
      </n-form-item>
      <n-form-item v-if="cfg.backupCleanStrategy === 1" label="保留数量">
        <n-input-number v-model:value="cfg.backupCleanKeepCount" :min="1" :step="1" />
      </n-form-item>
      <n-form-item v-if="cfg.backupCleanStrategy === 2">
        <template #label>
          <span>保留时间</span>
          <n-tooltip>
            <template #trigger>
              <n-icon><i-carbon-help-filled /></n-icon>
            </template>
            请输入带时间单位的时间间隔。支持的时间单位只有 h m s（分别代表小时、分钟、秒）。<br />
            示例：<br />
            720h：代表保留 720 小时（即 30 天）内的备份<br />
            10.5h：代表保留 10.5 小时（即 10 小时 30 分）内的备份<br />
            10h30m：保留 10 小时 30 分内备份的另一种写法
          </n-tooltip>
        </template>
        <n-input v-model:value="cfg.backupCleanKeepDur" style="width: 12rem" placeholder="" />
      </n-form-item>
      <n-form-item v-if="cfg.backupCleanStrategy !== 0">
        <template #label>
          <span>触发方式</span>
          <n-tooltip raw-content content="">
            <template #trigger>
              <n-icon><i-carbon-help-filled /></n-icon>
            </template>
            自动备份后：在每次自动备份完成后，顺便进行备份清理。<br />
            定时：按照给定的 cron 表达式，单独触发清理。
          </n-tooltip>
        </template>
        <n-checkbox-group v-model:value="backupCleanTriggers">
          <n-checkbox :value="CleanTrigger.AfterAutoBackup">自动备份后</n-checkbox>
          <n-checkbox :value="CleanTrigger.Cron">定时</n-checkbox>
        </n-checkbox-group>
      </n-form-item>
      <n-form-item v-if="cfg.backupCleanStrategy !== 0">
        <template #label>
          <span>定时间隔</span>
          <n-tooltip>
            <template #trigger>
              <n-icon><i-carbon-help-filled /></n-icon>
            </template>
            定时间隔表达式请参阅
            <a href="https://pkg.go.dev/github.com/robfig/cron" target="_blank">cron 文档</a>
          </n-tooltip>
        </template>
        <n-input v-model:value="cfg.backupCleanCron" style="width: 12rem" placeholder="" />
      </n-form-item>
    </n-form>
    <h4>如何恢复备份？</h4>
    <n-text>
      将骰子彻底关闭，解压备份压缩包到骰子目录。若提示“是否覆盖？”选择“全部”即可 (覆盖 data 目录)。
    </n-text>
  </div>

  <div style="display: flex; justify-content: space-between; align-items: center">
    <h2>已备份文件</h2>
    <n-button type="error" @click="enterBatchDelete">
      <template #icon>
        <i-carbon-row-delete />
      </template>
      进入批量删除页面
    </n-button>
  </div>

  <n-list>
    <n-list-item v-for="i in data.items" :key="i.name">
      <n-flex align="center" justify="space-between" class="mx-4">
        <div class="flex flex-col">
          <n-text class="self-start text-base">{{ i.name }}</n-text>
          <n-text v-if="(i?.selection ?? 0) >= 0" class="self-start text-xs" type="info">
            此备份包含：{{ parseSelectionDesc(i.selection).join('、') }}
          </n-text>
          <n-text v-else class="self-start text-xs" type="warning">此备份内容无法识别</n-text>
        </div>
        <n-flex size="small" justify="end">
          <n-button
            secondary
            size="tiny"
            tag="a"
            style="text-decoration: none; width: 8rem"
            :href="`${urlBase}/sd-api/backup/download?name=${encodeURIComponent(i.name)}&token=${encodeURIComponent(store.token)}`">
            下载 - {{ filesize(i.fileSize) }}
          </n-button>
          <n-button type="error" size="tiny" secondary @click="bakDeleteConfirm(i.name)">
            <template #icon>
              <i-carbon-row-delete />
            </template>
          </n-button>
        </n-flex>
      </n-flex>
    </n-list-item>
  </n-list>

  <n-modal v-model:show="showBatchDelete" preset="card" title="批量删除备份" class="diff-dialog">
    <div>
      <n-alert :closable="false" class="mb-4">
        默认勾选最近的 5 个备份之前的历史备份，可自行调整。
      </n-alert>
      <n-flex size="large" align="center" class="mb-2">
        <n-checkbox
          v-model:checked="checkAllBaks"
          :indeterminate="isIndeterminate"
          @change="handleCheckAllChange">
          {{ checkAllBaks ? '取消全选' : '全选' }}
        </n-checkbox>
        <n-text type="info" class="text-xs">
          已勾选 {{ selectedBaks.length }} 个备份，共
          {{ filesize(selectedBaks.map(bak => bak.fileSize).reduce((a, b) => a + b, 0)) }}
        </n-text>
      </n-flex>
      <n-checkbox-group v-model:value="selectedBaks" @change="handleCheckedBakChange">
        <n-list :show-divider="false">
          <n-list-item v-for="i of data.items" :key="i.name">
            <n-checkbox :value="i">
              {{ i.name }}
            </n-checkbox>
          </n-list-item>
        </n-list>
      </n-checkbox-group>
    </div>
    <template #footer>
      <n-flex wrap>
        <n-button @click="showBatchDelete = false">取消</n-button>
        <n-button
          type="error"
          :disabled="!(selectedBaks && selectedBaks.length > 0)"
          @click="bakBatchDeleteConfirm">
          <template #icon>
            <i-carbon-row-delete />
          </template>
          删除所选
        </n-button>
      </n-flex>
    </template>
  </n-modal>

  <n-modal v-model:show="showBackup" title="立即备份" preset="card" class="diff-dialog">
    <n-flex vertical>
      <n-flex vertical>
        <span>备份范围：</span>
        <n-checkbox-group v-model:value="backupSelections">
          <n-checkbox label="基础（含自定义回复）" value="base" checked disabled />
          <n-checkbox label="JS 插件" value="js" />
          <n-checkbox label="牌堆" value="deck" />
          <n-checkbox label="帮助文档" value="helpdoc" />
          <n-checkbox label="敏感词库" value="censor" />
          <n-checkbox label="人名信息" value="name" />
          <n-checkbox label="图片" value="image" />
        </n-checkbox-group>
      </n-flex>
      <n-flex align="center">
        <span>备份文件名预览：</span>
        <n-text type="info">
          bak_{{ now }}_r{{ formatSelection(backupSelections).toString(16) }}_&lt;随机值&gt;.zip
        </n-text>
      </n-flex>
    </n-flex>
    <template #footer>
      <n-flex wrap>
        <n-button @click="showBackup = false">取消</n-button>
        <n-button type="info" @click="doBackup">立即备份</n-button>
      </n-flex>
    </template>
  </n-modal>
</template>

<script lang="ts" setup>
import type { CheckboxValueType } from 'element-plus';
import { useStore } from '~/store';
import { urlBase } from '~/backend';
import { filesize } from 'filesize';
import { sum } from 'es-toolkit/compat';
import { dayjs } from 'element-plus';
import {
  getBackupConfig,
  getBackupList,
  postBackupBatchDel,
  postBackupDel,
  postDoBackup,
  setBackupConfig,
} from '~/api/v1/backup';
import { useDialog, useMessage } from 'naive-ui';

const message = useMessage();
const dialog = useDialog();
const store = useStore();

const data = ref<{
  items: any[];
}>({
  items: [],
});

const cfg = ref<any>({});
const now = ref(dayjs().format('YYMMDD_HHmmss'));
const showBackup = ref<boolean>(false);
const backupSelections = ref<string[]>([
  'base',
  'js',
  'deck',
  'helpdoc',
  'censor',
  'name',
  'image',
]);

const parseSelection = (selection: number): string[] => {
  const list = ['base'];
  const jsMark = selection & 0b000001;
  if (jsMark) {
    list.push('js');
  }
  const deckMark = selection & 0b000010;
  if (deckMark) {
    list.push('deck');
  }
  const helpdocMark = selection & 0b000100;
  if (helpdocMark) {
    list.push('helpdoc');
  }
  const censorMark = selection & 0b001000;
  if (censorMark) {
    list.push('censor');
  }
  const nameMark = selection & 0b010000;
  if (nameMark) {
    list.push('name');
  }
  const resourceMark = selection & 0b100000;
  if (resourceMark) {
    list.push('image');
  }
  return list;
};

const parseSelectionDesc = (selection: number): string[] => {
  const list = ['基础'];
  const jsMark = selection & 0b000001;
  if (jsMark) {
    list.push('JS 插件');
  }
  const deckMark = selection & 0b000010;
  if (deckMark) {
    list.push('牌堆');
  }
  const helpdocMark = selection & 0b000100;
  if (helpdocMark) {
    list.push('帮助文档');
  }
  const censorMark = selection & 0b001000;
  if (censorMark) {
    list.push('敏感词库');
  }
  const nameMark = selection & 0b010000;
  if (nameMark) {
    list.push('人名信息');
  }
  const resourceMark = selection & 0b100000;
  if (resourceMark) {
    list.push('图片');
  }
  return list;
};

const formatSelection = (selections: string[]): number => {
  let mark = 0;
  if (selections.includes('js')) {
    mark |= 0b000001;
  }
  if (selections.includes('deck')) {
    mark |= 0b000010;
  }
  if (selections.includes('helpdoc')) {
    mark |= 0b000100;
  }
  if (selections.includes('censor')) {
    mark |= 0b001000;
  }
  if (selections.includes('name')) {
    mark |= 0b010000;
  }
  if (selections.includes('image')) {
    mark |= 0b100000;
  }
  return mark;
};

watch(
  () => cfg.value.autoBackupSelectionList,
  v => {
    cfg.value.autoBackupSelection = formatSelection(v);
  },
);

const refreshList = async () => {
  const lst = await getBackupList();
  data.value = lst;
};

const configGet = async () => {
  const data = await getBackupConfig();
  cfg.value = data;
  cfg.value.autoBackupSelectionList = parseSelection(data.autoBackupSelection);
  if (data.backupCleanTrigger) {
    const triggers: CleanTrigger[] = [];
    if (data.backupCleanTrigger & CleanTrigger.Cron) {
      triggers.push(CleanTrigger.Cron);
    }
    if (data.backupCleanTrigger & CleanTrigger.AfterAutoBackup) {
      triggers.push(CleanTrigger.AfterAutoBackup);
    }
    backupCleanTriggers.value = triggers;
  }
};

const bakDeleteConfirm = async (name: string) => {
  dialog.warning({
    title: '提示',
    content: `确认删除「${name}」？`,
    positiveText: '确定',
    negativeText: '取消',
    closable: false,
    onPositiveClick: async () => {
      const r = await postBackupDel(name);
      if (!r.success) {
        message.error('删除失败');
      } else {
        message.success('已删除');
      }
    },
  });
  await refreshList();
};

const showBatchDelete = ref<boolean>(false);
const selectedBaks = ref<any[]>([]); // 他不是 string[]，是备份项的一种格式
const checkAllBaks = ref(false);
const isIndeterminate = ref(true);

const enterBatchDelete = async () => {
  selectedBaks.value = data.value.items.filter((_, index) => index >= 5);
  showBatchDelete.value = true;
};

const handleCheckAllChange = (val: CheckboxValueType) => {
  selectedBaks.value = val ? data.value.items : [];
  isIndeterminate.value = false;
};

const handleCheckedBakChange = (value: CheckboxValueType[]) => {
  const checkedCount = value.length;
  checkAllBaks.value = checkedCount === data.value.items.length;
  isIndeterminate.value = checkedCount > 0 && checkedCount < data.value.items.length;
};

const bakBatchDeleteConfirm = async () => {
  dialog.warning({
    title: '提示',
    content: '确认删除所选备份？删除的内容无法找回！',
    positiveText: '确定',
    negativeText: '取消',
    closable: false,
    onPositiveClick: async () => {
      const res = await postBackupBatchDel(selectedBaks.value.map(bak => bak.name));
      if (res.result) {
        message.success('已删除所选备份');
      } else {
        message.error('有备份删除失败！失败文件：\n' + res.fails.join('\n'));
      }
      showBatchDelete.value = false;
      await refreshList();
    },
  });
};

const doBackup = async () => {
  const ret = await postDoBackup(formatSelection(backupSelections.value));
  showBackup.value = false;
  await refreshList();
  if (ret.testMode) {
    message.success('展示模式无法备份');
  } else {
    message.success('已进行备份');
  }
};

const doSave = async () => {
  await setBackupConfig(cfg.value);
  message.success('已保存');
};

const enum CleanTrigger {
  // 定时
  Cron = 1 << 0,
  // 自动备份后
  AfterAutoBackup = 1 << 1,
}

const backupCleanTriggers = ref<CleanTrigger[]>();

watch(backupCleanTriggers, newStrategies => {
  cfg.value.backupCleanTrigger = sum(newStrategies);
});

let timerId: number;

const refreshNow = async () => {
  now.value = dayjs().format('YYMMDD_HHmmss');
  timerId = setTimeout(refreshNow, 1000);
};

onBeforeMount(async () => {
  await configGet();
  await refreshList();
  await refreshNow();
});

onBeforeUnmount(() => {
  clearTimeout(timerId);
});
</script>

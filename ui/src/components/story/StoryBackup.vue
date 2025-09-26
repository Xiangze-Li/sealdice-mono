<template>
  <tip-box type="info">
    <n-text>
      每次向染色器上传跑团日志之前，都会在本地先保留一份备份，再进行上传。<br />
      确定不再需要时，你可以在此处删除这些备份文件。<br />
      <br />
      <strong>删除此处的备份文件不会使日志丢失。</strong>
    </n-text>
  </tip-box>

  <header class="my-2 flex items-center justify-between">
    <n-flex size="large" align="center">
      <n-checkbox
        v-model:checked="checkAllBackups"
        :indeterminate="isIndeterminate"
        :disabled="!(backups && backups.size > 0)"
        @update:checked="handleCheckAllChange">
        {{ checkAllBackups ? '取消全选' : '全选' }}
      </n-checkbox>
      <n-text type="info" class="text-xs">
        已勾选 {{ selectedBackups.length }} 个备份，共
        {{
          filesize(
            selectedBackups
              .map(n => backups.get(n))
              .map(bak => bak?.fileSize ?? 0)
              .reduce((a, b) => a + b, 0),
          )
        }}
      </n-text>
    </n-flex>
    <n-button
      type="error"
      :disabled="!(selectedBackups && selectedBackups.length > 0)"
      @click="backupBatchDeleteConfirm">
      <template #icon>
        <n-icon><i-carbon-row-delete /></n-icon>
      </template>
      删除所选
    </n-button>
  </header>

  <n-text v-if="backups.size === 0" class="text-base">无备份文件</n-text>
  <n-list v-else class="backup-list">
    <n-checkbox-group v-model:value="selectedBackups" @update:value="handleCheckedBackupChange">
      <n-list-item>
        <n-flex
          align="center"
          justify="space-between"
          v-for="backup in backups.values()"
          :key="backup.name"
          class="backup-line">
          <n-checkbox :value="backup.name" :label="backup.name" />
          <n-flex size="small" wrap style="margin-left: 1px; justify-content: flex-end">
            <n-button
              size="small"
              secondary
              tag="a"
              style="text-decoration: none; width: 8rem"
              :href="`${urlBase}/sd-api/story/backup/download?name=${encodeURIComponent(backup.name)}&token=${encodeURIComponent(store.token)}`">
              下载 - {{ filesize(backup.fileSize) }}
            </n-button>
            <n-button type="error" size="small" secondary @click="bakDeleteConfirm(backup.name)">
              <template #icon>
                <n-icon><i-carbon-row-delete /></n-icon>
              </template>
              删除
            </n-button>
          </n-flex>
        </n-flex>
      </n-list-item>
    </n-checkbox-group>
  </n-list>
</template>

<script setup lang="ts">
import { filesize } from 'filesize';
import { useStore } from '~/store';
import type { Backup } from '~/api/v1/story';
import { getStoryBackUpList, postStoryBatchDel } from '~/api/v1/story';
import { urlBase } from '~/backend';
import { useDialog, useMessage } from 'naive-ui';

const store = useStore();
const message = useMessage();
const dialog = useDialog();

const backups = ref<Map<string, Backup>>(new Map());

const refreshList = async () => {
  const resp = await getStoryBackUpList();
  if (resp?.result) {
    backups.value = new Map();
    for (const backup of resp.data) {
      backups.value.set(backup.name, backup);
    }
  }
  selectedBackups.value = [];
};

// const bakDownloadConfirm = async (name: string) => {
//   const res = await postStoryBatchDel([name]);
//   if (res?.result) {
//     ElMessage.success('已删除');
//   } else {
//     ElMessage.error('删除失败');
//   }
// };

const bakDeleteConfirm = async (name: string) => {
  dialog.warning({
    title: '提示',
    content: `确认删除备份文件「${name}」？`,
    positiveText: '确定',
    negativeText: '取消',
    closable: false,
    onPositiveClick: async () => {
      const res = await postStoryBatchDel([name]);
      if (res?.result) {
        message.success('已删除');
      } else {
        message.error('删除失败');
      }
      nextTick(async () => {
        await refreshList();
      });
    },
  });
};

const selectedBackups = ref<string[]>([]);
const checkAllBackups = ref(false);
const isIndeterminate = ref<boolean>();

const handleCheckAllChange = (val: boolean) => {
  selectedBackups.value = val ? Array.from(backups.value.values().map(backup => backup.name)) : [];
  isIndeterminate.value = false;
};

const handleCheckedBackupChange = (
  value: (string | number)[],
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  meta: { actionType: 'check' | 'uncheck'; value: string | number },
) => {
  const checkedCount = value.length;
  checkAllBackups.value = checkedCount === backups.value.size;
  isIndeterminate.value = checkedCount > 0 && checkedCount < backups.value.size;
};

const backupBatchDeleteConfirm = async () => {
  dialog.warning({
    title: '提示',
    content: `确认删除选择的所有跑团日志备份？`,
    positiveText: '确定',
    negativeText: '取消',
    closable: false,
    onPositiveClick: async () => {
      const res = await postStoryBatchDel(selectedBackups.value);
      if (res.result) {
        message.success('已删除所选备份');
      } else {
        message.error('有备份删除失败！失败文件：\n' + res.fails.join('\n'));
      }
      nextTick(async () => {
        await refreshList();
      });
    },
  });
};

onBeforeMount(async () => {
  await refreshList();
});
</script>

<style scoped lang="css">
.backup-list {
  display: flex;
  flex-direction: column;

  .backup-line {
    padding: 3px 1rem;
    display: flex;
    justify-content: space-between;
    flex-wrap: wrap;
  }

  .backup-line:not(:first-child) {
    border-top: 1px solid var(--el-border-color);
  }
}
</style>

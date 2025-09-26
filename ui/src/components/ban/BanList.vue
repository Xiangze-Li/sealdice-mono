<script setup lang="ts">
import type { UploadUserFile } from 'element-plus';
import { urlBase } from '~/backend';
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import {
  getBanConfigList,
  importBanConfig,
  postMapAddOne,
  postMapDelOne,
} from '~/api/v1/banconfig';
import { type UploadFileInfo, useDialog, useMessage } from 'naive-ui';

dayjs.extend(relativeTime);

const message = useMessage();
const dialog = useDialog();

const recordList = ref<any[]>([]);
const recordPage = ref({
  no: 1,
  total: 0,
  pageSize: 10,
});

const showBanned = ref(true);
const showWarn = ref(true);
const showTrusted = ref(true);
const showOthers = ref(true);
const dialogAddShow = ref(false);

const banRankText = new Map<number, string>();
banRankText.set(-30, '禁止');
banRankText.set(-10, '警告');
banRankText.set(+30, '信任');
banRankText.set(0, '常规');

const addData = ref<{ id: string; rank: number; name: string; reason: string }>({
  id: '',
  rank: -30,
  reason: '',
  name: '',
});

const doAdd = async () => {
  if (addData.value.id === '') return;
  await postMapAddOne(
    addData.value.id,
    addData.value.rank,
    addData.value.name,
    addData.value.reason,
  );
  await refreshList();
  message.success('已保存');
  dialogAddShow.value = false;
};

const searchBy = ref('');

const filteredRecordList = computed<any[]>(() => {
  if (recordList.value) {
    const items = [];
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    for (const [k, _v] of Object.entries(recordList.value)) {
      const v = _v as any;
      let ok = false;
      if (v.rank === -30 && showBanned.value) {
        ok = true;
      }
      if (v.rank === -10 && showWarn.value) {
        ok = true;
      }
      if (v.rank === 30 && showTrusted.value) {
        ok = true;
      }
      if (v.rank === 0 && showOthers.value) {
        ok = true;
      }

      if (ok && searchBy.value !== '') {
        let a = false;
        let b = false;
        if (v.ID.indexOf(searchBy.value) !== -1) {
          a = true;
        }
        if (v.name.indexOf(searchBy.value) !== -1) {
          b = true;
        }
        ok = a || b;
      }

      v.rankText = banRankText.get(v.rank);

      if (ok) items.push(v);
    }

    return items;
  }
  return [];
});

const groupItems = computed(() => {
  const start = (recordPage.value.no - 1) * recordPage.value.pageSize;
  const end = start + recordPage.value.pageSize;
  return filteredRecordList.value.slice(start, end);
});

const refreshList = async () => {
  const lst = await getBanConfigList();
  recordList.value = lst;
  recordPage.value.total = lst.length;
  recordPage.value.no = 1;
};

const deleteOne = async (i: any) => {
  dialog.warning({
    title: '删除',
    content: '是否删除此记录？',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      await postMapDelOne(i);
      await refreshList();
      message.success('已保存');
    },
  });
};

const beforeUpload = async (data: { file: UploadFileInfo }) => {
  const file = data.file.file as File;
  const c = await importBanConfig(file);
  if (c.result) {
    message.success('导入黑白名单完成');
    await nextTick(async () => {
      await refreshList();
    });
  } else {
    message.error('导入黑白名单失败！' + c.err);
  }
};

onBeforeMount(async () => {
  await refreshList();
});
</script>

<template>
  <header class="flex flex-wrap-reverse justify-between gap-y-4">
    <n-flex size="small" align="center">
      <n-text class="text-base">搜索：</n-text>
      <span>
        <n-input v-model:value="searchBy" class="w-64" placeholder="请输入帐号或名字的一部分" />
      </span>
    </n-flex>

    <n-flex align="center">
      <n-button type="success" secondary @click="dialogAddShow = true">
        <template #icon>
          <i-carbon-add-large />
        </template>
        添加
      </n-button>
      <span>
        <n-upload
          action=""
          multiple
          accept="application/json,.json"
          :show-file-list="false"
          @before-upload="beforeUpload">
          <n-button type="info" secondary>
            <template #icon>
              <n-icon><i-carbon-upload /></n-icon>
            </template>
            导入
          </n-button>
        </n-upload>
      </span>
      <n-button
        secondary
        type="info"
        tag="a"
        target="_blank"
        :href="`${urlBase}/sd-api/banconfig/export`"
        style="text-decoration: none">
        <template #icon>
          <n-icon><i-carbon-download /></n-icon>
        </template>
        导出
      </n-button>
    </n-flex>
  </header>

  <n-flex align="center" class="my-2">
    <n-text class="text-base">级别：</n-text>
    <n-checkbox v-model:checked="showBanned">拉黑</n-checkbox>
    <n-checkbox v-model:checked="showWarn">警告</n-checkbox>
    <n-checkbox v-model:checked="showTrusted">信任</n-checkbox>
    <n-checkbox v-model:checked="showOthers">其它</n-checkbox>
  </n-flex>

  <main class="mt-4">
    <n-list>
      <!-- eslint-disable-next-line vue/no-unused-vars-->
      <n-list-item v-for="(i, index) in groupItems" :key="i.ID" shadow="hover" class="w-full">
        <n-thing class="mx-4 my-2">
          <template #header>
            <n-flex size="small" align="center">
              <n-tag v-if="i.rankText === '禁止'" type="error" :bordered="false">
                {{ i.rankText }}
              </n-tag>
              <n-tag v-else-if="i.rankText === '警告'" type="warning" :bordered="false">
                {{ i.rankText }}
              </n-tag>
              <n-tag v-else-if="i.rankText === '信任'" type="success" :bordered="false">
                {{ i.rankText }}
              </n-tag>
              <n-tag v-else :bordered="false">{{ i.rankText }}</n-tag>
              <n-text tag="strong" class="text-base">{{ i.ID }}</n-text>
            </n-flex>
          </template>
          <template #header-extra>
            <n-flex justify="flex-end">
              <n-button type="error" size="small" secondary @click="deleteOne(i)">
                <template #icon>
                  <n-icon><i-carbon-row-delete /></n-icon>
                </template>
                删除
              </n-button>
            </n-flex>
          </template>
          <template #description>
            <n-flex size="small" align="center" wrap>
              <n-text class="text-sm">「{{ i.name }}」</n-text>
              <n-text tag="em" class="text-sm">怒气值：{{ i.score }}</n-text>
            </n-flex>
          </template>
          <n-flex vertical>
            <div v-for="(j, index) in i.reasons" :key="index">
              <n-flex size="small">
                <n-tooltip>
                  <template #trigger>
                    <n-tag size="small" type="info" :bordered="false">
                      {{ dayjs.unix(i.times[index]).fromNow() }}
                    </n-tag>
                  </template>
                  {{ dayjs.unix(i.times[index]).format('YYYY-MM-DD HH:mm:ssZ[Z]') }}
                </n-tooltip>
                <n-text>在&lt;{{ i.places[index] }}>，原因：「{{ j }}」</n-text>
              </n-flex>
            </div>
          </n-flex>
        </n-thing>
      </n-list-item>
    </n-list>
  </main>

  <footer class="mt-4 flex justify-center">
    <n-pagination
      v-model:page="recordPage.no"
      v-model:page-size="recordPage.pageSize"
      show-size-picker
      :page-sizes="[10, 20, 30, 50]"
      :page-slot="5"
      :item-count="filteredRecordList.length" />
  </footer>

  <n-modal v-model:show="dialogAddShow" preset="card" title="添加用户/群组" :closable="false">
    <n-flex>
      <n-form label-placement="left" label-width="auto">
        <n-form-item label="用户ID" required>
          <n-input
            v-model:value="addData.id"
            placeholder="必须为 QQ:12345 或 QQ-Group:12345 格式" />
        </n-form-item>
        <n-form-item label="名称">
          <n-input v-model:value="addData.name" placeholder="自动" />
        </n-form-item>
        <n-form-item label="原因">
          <n-input v-model:value="addData.reason" placeholder="骰主后台设置" />
        </n-form-item>
        <n-form-item label="身份">
          <n-radio-group v-model:value="addData.rank">
            <n-radio
              v-for="item in [
                { label: '禁用', value: -30 },
                { label: '信任', value: 30 },
              ]"
              :key="item.value"
              :label="item.label"
              :value="item.value" />
          </n-radio-group>
        </n-form-item>
      </n-form>
    </n-flex>
    <template #footer>
      <n-flex>
        <n-button @click="dialogAddShow = false">取消</n-button>
        <n-button type="success" @click="doAdd">添加</n-button>
      </n-flex>
    </template>
  </n-modal>
</template>

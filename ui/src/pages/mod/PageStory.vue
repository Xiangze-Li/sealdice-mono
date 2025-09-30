<script lang="ts" setup>
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import randomColor from 'randomcolor';
import {
  deleteStoryLog,
  getStoryInfo,
  getStoryItemPage,
  getStoryLogPage,
  postStoryLog,
} from '~/api/v1/story';
import { useDialog, useMessage } from 'naive-ui';

const message = useMessage();
const dialog = useDialog();

interface Log {
  id: number;
  name: string;
  groupId: string;
  createdAt: number;
  updatedAt: number;
  size: number;
  pitch?: boolean;
  current?: number;
}

interface Item {
  id: number;
  logId: number;
  nickname: string;
  IMUserId: string;
  time: number;
  message: string;
  isDice: boolean;
  isEdit?: boolean;
}

async function getInfo() {
  return getStoryInfo();
  //   return backend.get(url("info")) as any
}

// async function getLogs() {
//     return apiFetch(url("logs"), {
//         method: "get", headers: { token: token }
//     })
// }

const queryLogPage = ref({
  pageNum: 1,
  pageSize: 20,
  total: 0,
  name: '',
  groupId: '',
  createdTime: null as unknown as [number, number],
});

const getLogPage = getStoryLogPage;
// async function getLogPage(params: { pageNum: number, pageSize: number, name?: string, groupId?: string, createdTimeBegin?: number, createdTimeEnd?: number }) {
//     return await backend.get(url("logs/page"), {
//       headers: {token: token}, params: params
//     }) as any
// }

// async function getItems(v: Log) {
//     // ofetch get+params 至少在开发模式有莫名奇妙的 bug，会丢失 baseURL
//     // 下面的接口就先不更换了
//     return await backend.get(url('items'), { params: v, headers: { token } }) as unknown as Item[]
// }

const logItemPage = ref({
  pageNum: 1,
  pageSize: 100,
  size: 0,
  logName: '',
  groupId: '',
});

// async function getItemPage(params: { pageNum: number, pageSize: number, logName: string, groupId: string }) {
//     return backend.get(url("items/page"), {
//         headers: { token: token }, params: params
//     })
// }
const getItemPage = getStoryItemPage;

async function delLog(v: Log) {
  return deleteStoryLog(v);
  // return backend.delete(url('log'), { headers: { token }, data: v }) as unknown as boolean
}

const tab: Ref<'list' | 'backup'> = ref('list');
const mode: Ref<'logs' | 'items'> = ref('logs');
const sum_log = ref(0),
  sum_item = ref(0),
  cur_log = ref(0),
  cur_item = ref(0);
dayjs.extend(relativeTime);

const logs: Ref<Log[]> = ref([]);

async function searchLogs() {
  const params = {
    ...queryLogPage.value,
    createdTimeBegin: queryLogPage.value.createdTime?.[0]
      ? dayjs(queryLogPage.value.createdTime?.[0]).startOf('date').unix()
      : undefined,
    createdTimeEnd: queryLogPage.value.createdTime?.[1]
      ? dayjs(queryLogPage.value.createdTime?.[1]).endOf('date').unix()
      : undefined,
  };
  const result:
    | { result: false; err?: string }
    | {
        result: true;
        total: number;
        data: Log[];
        pageNum: number;
        pageSize: number;
      } = await getLogPage(params);
  if (result.result) {
    logs.value = result.data;
    queryLogPage.value.total = result.total;
  } else {
    message.error('无法获取跑团日志' + result.err);
  }
}

const refreshLogs = async () => {
  [sum_log.value, sum_item.value, cur_log.value, cur_item.value] = await getInfo();
  await searchLogs();
};

const handleLogPageChange = async (val: number) => {
  queryLogPage.value.pageNum = val;
  await searchLogs();
};

const handlePageSizeChange = async (val: number) => {
  queryLogPage.value.pageNum = 1;
  queryLogPage.value.pageSize = val;
  await searchLogs();
};

async function DelLog(v: Log, flag = true) {
  dialog.warning({
    title: '删除',
    content: '是否删除此跑团日志？',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      const info = await delLog(v);
      if (info) {
        message.success('删除成功');
        if (flag) await refreshLogs();
      } else {
        message.error('删除失败');
      }
    },
  });
}

async function DelLogs() {
  dialog.warning({
    title: '删除',
    content: '是否删除所选跑团日志？',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      const ls = [];
      for (const v of logs.value) {
        if (v.pitch == true) {
          ls.push(v);
        }
      }
      for (const v of ls) {
        const info = await delLog(v);
        if (info) {
          message.success('删除成功');
        } else {
          message.error('删除失败');
        }
      }
      nextTick(async () => {
        await refreshLogs();
      });
    },
  });
}

async function UploadLog(v: Log) {
  dialog.warning({
    title: '警告',
    content: '将此跑团日志上传至海豹服务器？',
    positiveText: '确定',
    negativeText: '取消',
    closable: false,
    onPositiveClick: async () => {
      const info = await postStoryLog(v);
      message.info(() => h('span', { innerHTML: info }), {
        closable: true,
        duration: 0,
      });
    },
  });
}
//

const item_data: Ref<Item[]> = ref([]);

const users = ref({}) as Ref<Record<string, Array<string>>>;

async function openItem(log: Log) {
  logItemPage.value.logName = log.name;
  logItemPage.value.groupId = log.groupId;
  logItemPage.value.size = log.size;
  logItemPage.value.pageNum = 1;
  item_data.value = (await getItemPage({
    pageNum: logItemPage.value.pageNum,
    pageSize: logItemPage.value.pageSize,
    logName: logItemPage.value.logName,
    groupId: logItemPage.value.groupId,
  })) as unknown as Item[];
  mode.value = 'items';
}

const handleItemPageChange = async (val: number) => {
  logItemPage.value.pageNum = val;
  item_data.value = (await getItemPage(logItemPage.value)) as unknown as Item[];
};

function closeItem() {
  item_data.value = [];
  mode.value = 'logs';
  users.value = {};
}

const randomColorWithIndex = (i: number): string => {
  const presets = [
    'var(--color-red-600)',
    'var(--color-orange-600)',
    'var(--color-yellow-600)',
    'var(--color-green-600)',
    'var(--color-cyan-600)',
    'var(--color-blue-600)',
    'var(--color-purple-600)',
    'var(--color-pink-600)',
    'var(--color-slate-600)',
  ];
  const randomColorSystems = [
    'red',
    'orange',
    'yellow',
    'green',
    'blue',
    'purple',
    'pink',
    'monochrome',
  ];
  if (i < presets.length) {
    return presets[i];
  } else {
    return randomColor({
      hue: randomColorSystems[i % randomColorSystems.length],
      luminosity: 'dark',
    });
  }
};

const items = computed(() => {
  const items: Item[] = [];
  item_data.value.forEach((v, i) => {
    if (!users.value[v.IMUserId]) {
      users.value[v.IMUserId] = [randomColorWithIndex(i), v.nickname];
    }
    items.push(v);
  });
  return items;
});

//

onBeforeMount(async () => {
  await refreshLogs();
});
</script>

<template>
  <n-tabs v-model:value="tab" pane-class="mb-8" justify-content="space-evenly">
    <n-tab-pane tab="跑团日志" name="list">
      <template v-if="mode == 'logs'">
        <header>
          <n-card title="跑团日志 / Story">
            <n-flex vertical align="flex-start">
              <n-text>记录过 {{ sum_log }} 份日志，共计 {{ sum_item }} 条消息</n-text>
              <n-text>现有 {{ cur_log }} 份日志，共计 {{ cur_item }} 条消息</n-text>
            </n-flex>
          </n-card>
        </header>
        <n-divider />
        <main>
          <n-form
            v-model:model="queryLogPage"
            size="small"
            class="flex flex-wrap"
            label-width="auto"
            label-placement="left"
            inline>
            <n-form-item label="日志名">
              <n-input v-model:value="queryLogPage.name" placeholder="" clearable />
            </n-form-item>
            <n-form-item label="群号">
              <n-input v-model:value="queryLogPage.groupId" placeholder="" clearable />
            </n-form-item>
            <n-form-item label="创建时间">
              <n-date-picker v-model:value="queryLogPage.createdTime" type="daterange" clearable />
            </n-form-item>
            <n-form-item>
              <n-button type="info" secondary @click="searchLogs">查询</n-button>
            </n-form-item>
          </n-form>
          <n-button-group cl style="margin-top: 5px; display: block">
            <n-button type="primary" size="small" @click="logs.forEach(v => (v.pitch = !v.pitch))">
              <template #icon>
                <n-icon><i-carbon-checkmark /></n-icon>
              </template>
              全选
            </n-button>
            <n-button
              v-show="logs?.filter(v => v.pitch)?.length ?? 0 > 0"
              type="error"
              size="small"
              @click="DelLogs()">
              <template #icon>
                <n-icon><i-carbon-row-delete /></n-icon>
              </template>
              删除所选
            </n-button>
          </n-button-group>
          <template v-for="i in logs" :key="i.id">
            <foldable-card style="margin-top: 10px">
              <template #title>
                <n-flex align="center">
                  <n-checkbox v-model:checked="i.pitch" />
                  <n-flex align="center" wrap>
                    <n-text class="text-base" tag="strong">{{ i.name }}</n-text>
                    <n-text>({{ i.groupId }})</n-text>
                  </n-flex>
                </n-flex>
              </template>
              <template #action>
                <n-flex>
                  <n-button size="small" secondary @click="openItem(i)">查看</n-button>
                  <n-button size="small" type="primary" secondary @click="UploadLog(i)">
                    <template #icon>
                      <n-icon><i-carbon-upload /></n-icon>
                    </template>
                    提取日志
                  </n-button>
                  <n-button size="small" type="error" secondary @click="DelLog(i)">
                    <template #icon>
                      <n-icon><i-carbon-row-delete /></n-icon>
                    </template>
                    删除
                  </n-button>
                </n-flex>
              </template>

              <n-flex vertical align="flex-start">
                <n-flex>
                  <n-text>包含 {{ i.size }} 条消息</n-text>
                </n-flex>
                <n-flex>
                  <n-text>创建于：{{ dayjs.unix(i.createdAt).format('YYYY-MM-DD') }}</n-text>
                  <n-tag type="info" size="small" :bordered="false">
                    {{ dayjs.unix(i.createdAt).fromNow() }}
                  </n-tag>
                </n-flex>
                <n-flex>
                  <n-text>更新于：{{ dayjs.unix(i.updatedAt).format('YYYY-MM-DD') }}</n-text>
                  <n-tag type="info" size="small" :bordered="false">
                    {{ dayjs.unix(i.updatedAt).fromNow() }}
                  </n-tag>
                </n-flex>
              </n-flex>
            </foldable-card>
          </template>
        </main>
        <div class="mt-4 flex justify-center">
          <n-pagination
            v-model:page="queryLogPage.pageNum"
            v-model:page-size="queryLogPage.pageSize"
            show-size-picker
            :page-sizes="[10, 20, 30, 50]"
            :page-slot="3"
            :item-count="queryLogPage.total"
            @update:page="handleLogPageChange"
            @update:page-size="handlePageSizeChange" />
        </div>
      </template>
      <template v-if="mode == 'items'">
        <n-card title="跑团日志 / Story">
          <template #header-extra>
            <n-button type="primary" @click="closeItem()">
              <template #icon>
                <n-icon><i-carbon-chevron-left /></n-icon>
              </template>
              返回列表
            </n-button>
          </template>
          <n-collapse>
            <n-collapse-item title="颜色设置">
              <template v-for="(_, id) in users" :key="id">
                <n-descriptions label-placement="top">
                  <n-descriptions-item :label="users[id][1]">
                    <n-color-picker
                      class="w-32"
                      v-model:value="users[id][0]"
                      :modes="['hex']"
                      :show-alpha="false"
                      :swatches="randomColor({ count: 16 })" />
                  </n-descriptions-item>
                </n-descriptions>
              </template>
            </n-collapse-item>
          </n-collapse>
        </n-card>
        <div class="my-4 px-4">
          <template v-for="(v, i1) in items" :key="i1">
            <p :style="{ color: users[v.IMUserId][0] }">
              <span>{{ v.nickname }}：</span>
              <template v-for="(p1, i2) in v.message.split('\n')" :key="i2">
                <span>{{ p1 }}</span
                ><br />
              </template>
            </p>
          </template>
        </div>
        <div class="mt-4 flex justify-center">
          <n-pagination
            v-model:page="logItemPage.pageNum"
            v-model:page-size="logItemPage.pageSize"
            show-size-picker
            :page-sizes="[50, 100, 200]"
            :page-slot="5"
            :item-count="logItemPage.size"
            @update:page="handleItemPageChange"
            @update:page-size="handleItemPageChange" />
        </div>
      </template>
    </n-tab-pane>
    <n-tab-pane tab="日志备份" name="backup">
      <story-backup />
    </n-tab-pane>
  </n-tabs>
</template>

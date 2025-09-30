<template>
  <div class="censor-log-container">
    <header class="censor-log-header">
      <n-button type="info" secondary @click="refreshCensorLog">
        <template #icon>
          <n-icon><i-carbon-renew /></n-icon>
        </template>
        刷新
      </n-button>
      <n-pagination size="small"
        v-model:page="logQuery.pageNum"
        v-model:page-size="logQuery.pageSize"
        :item-count="logQuery.total"
        :page-slot="5"
        :default-page-size="20"
        @update:page="handleCurrentPageChange"
        @update:page-size="handlePageSizeChange" />
    </header>
    <n-data-table :columns="columns" :data="logs" class="mt-4" />
  </div>
</template>

<script lang="tsx" setup>
import dayjs from 'dayjs';
import { useCensorStore } from '~/components/mod/censor/censor';
import SensitiveTag from '~/components/mod/censor/sensitive-tag';
import { getCensorLogs } from '~/api/v1/censor';
import { useMessage, type DataTableColumns } from 'naive-ui';

const message = useMessage();
const censorStore = useCensorStore();

interface CensorLog {
  id: number;
  msgType: string;
  userId: string;
  groupId: string;
  content: string;
  highestLevel: string;
  createAt: number;
}

const logs = ref<CensorLog[]>();
const logQuery = ref({
  pageNum: 1,
  pageSize: 10,
  total: 0,
});

const columns: DataTableColumns<CensorLog> = [
  {
    title: '命中级别',
    key: 'highestLevel',
    render: ({ highestLevel }) => {
      switch (highestLevel) {
        case '1':
          return <SensitiveTag type="default" />;
        case '2':
          return <SensitiveTag type="info" />;
        case '3':
          return <SensitiveTag type="warning" />;
        case '4':
          return <SensitiveTag type="error" />;
        default:
          return <SensitiveTag type="default" message="未知" />;
      }
    },
  },
  {
    title: '消息类型',
    key: 'msgType',
    render: ({ msgType }) => {
      if (msgType === 'private') {
        return <n-text>私聊</n-text>;
      } else if (msgType === 'group') {
        return <n-text>群</n-text>;
      } else {
        return <n-text>未知</n-text>;
      }
    },
  },
  {
    title: '用户',
    key: 'userId',
  },
  {
    title: '群',
    key: 'groupId',
  },
  {
    title: '内容',
    key: 'content',
  },
  {
    title: '消息时间',
    key: 'createAt',
    render: ({ createAt }) => {
      return <>{dayjs.unix(createAt).format('YYYY-MM-DD HH:mm:ss')}</>;
    },
  },
];

censorStore.$subscribe(async (_, state) => {
  if (state.logsNeedRefresh === true) {
    await refreshCensorLog();
    state.logsNeedRefresh = false;
  }
});

const handleCurrentPageChange = async (val: number) => {
  logQuery.value.pageNum = val;
  await refreshCensorLog();
};

const handlePageSizeChange = async (val: number) => {
  logQuery.value.pageNum = 1;
  logQuery.value.pageSize = val;
  await refreshCensorLog();
};

const refreshCensorLog = async () => {
  const c:
    | { result: false }
    | {
        result: true;
        data: CensorLog[];
        total: number;
      } = await getCensorLogs(logQuery.value.pageNum, logQuery.value.pageSize);
  if (c?.result) {
    logs.value = c.data;
    logQuery.value.total = c.total;
  } else {
    message.error('无法获取拦截日志');
  }
};

onMounted(async () => {
  await refreshCensorLog();
});
</script>

<style scoped lang="css">
.censor-log-container {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.censor-log-header {
  display: flex;
  width: 100%;
  align-items: center;
  justify-content: space-between;
  text-align: center;
  flex-wrap: wrap;
  gap: 1rem;
}

.pagination {
  background-color: #f3f5f7;
}
</style>

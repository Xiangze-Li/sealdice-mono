<template>
  <Teleport v-if="store.curDice.logs.length" to="#root">
    <n-button type="default" class="btn-scrolldown" circle @click="scrollDown" content="最新日志">
      <template #icon>
        <n-icon><i-carbon-chevron-down /></n-icon>
      </template>
    </n-button>
  </Teleport>

  <div style="display: flex; justify-content: flex-end; align-items: center">
    <div style="display: flex; flex-direction: column">
      <n-tooltip
        v-if="
          store.curDice.baseInfo.versionCode < store.curDice.baseInfo.versionNewCode &&
          store.curDice.baseInfo.containerMode
        ">
        <template #trigger>
          <n-button type="primary" disabled>升级新版</n-button>
        </template>
        容器模式下禁止直接更新，请手动拉取最新镜像
      </n-tooltip>
      <n-button
        v-else-if="store.curDice.baseInfo.versionCode < store.curDice.baseInfo.versionNewCode"
        type="primary"
        @click="upgradeDialogVisible = true">
        升级新版
      </n-button>
    </div>
  </div>

  <h4>状态</h4>
  <div class="flex flex-col justify-center gap-4">
    <div class="flex flex-wrap items-center gap-1">
      <n-text>内存占用：</n-text>
      <n-text class="mr-2">{{ filesize(store.curDice.baseInfo.memoryUsedSys || 0) }}</n-text>
      <n-text type="info" class="text-xs">
        理论内存占用，数值偏大。系统任务管理器中的「活动内存」才是实际使用的系统内存。
      </n-text>
    </div>
    <website-health-check />
  </div>

  <div class="flex items-center justify-between">
    <h4>日志</h4>
    <n-checkbox v-model:checked="autoRefresh">保持刷新</n-checkbox>
  </div>

  <n-divider>
    <n-text type="warning" class="text-xs hover:cursor-pointer" @click="scrollDown">
      点击下拉到底查看最新
    </n-text>
  </n-divider>

  <div class="logs hidden p-0 md:block">
    <el-table
      :data="store.curDice.logs"
      :row-class-name="getLogRowClassName"
      :header-cell-style="{ backgroundColor: '#f3f5f7' }">
      <el-table-column label="时间" width="90">
        <template #default="scope">
          <div style="display: flex; align-items: center">
            <n-icon v-if="scope.row.msg.startsWith('onebot | ')" color="var(--el-color-warning)">
              <i-carbon-time />
            </n-icon>
            <n-icon v-else-if="scope.row.msg.startsWith('发给')" color="var(--el-color-primary)">
              <i-carbon-time />
            </n-icon>
            <n-icon v-else-if="scope.row.level === 'warn'" color="var(--el-color-warning)">
              <i-carbon-time />
            </n-icon>
            <n-icon v-else-if="scope.row.level === 'error'" color="var(--el-color-danger)">
              <i-carbon-time />
            </n-icon>
            <n-icon v-else><i-carbon-time /></n-icon>
            <span style="margin-left: 0.3rem">
              <span
                v-if="scope.row.msg.startsWith('onebot | ')"
                style="color: var(--el-color-warning)"
                >{{ dayjs.unix(scope.row.ts).format('HH:mm:ss') }}</span
              >
              <span
                v-else-if="scope.row.msg.startsWith('发给')"
                style="color: var(--el-color-primary)"
                >{{ dayjs.unix(scope.row.ts).format('HH:mm:ss') }}</span
              >
              <span v-else-if="scope.row.level === 'warn'" style="color: var(--el-color-warning)">{{
                dayjs.unix(scope.row.ts).format('HH:mm:ss')
              }}</span>
              <span v-else-if="scope.row.level === 'error'" style="color: var(--el-color-danger)">{{
                dayjs.unix(scope.row.ts).format('HH:mm:ss')
              }}</span>
              <span v-else>{{ dayjs.unix(scope.row.ts).format('HH:mm:ss') }}</span>
            </span>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="level" label="级别" width="55">
        <template #default="scope">
          <n-text v-if="scope.row.msg.startsWith('onebot | ')" type="warning">
            {{ scope.row.level }}
          </n-text>
          <n-text v-else-if="scope.row.msg.startsWith('发给')" type="primary">
            {{ scope.row.level }}
          </n-text>
          <n-text v-else-if="scope.row.level === 'warn'" type="warning">
            {{ scope.row.level }}
          </n-text>
          <n-text v-else-if="scope.row.level === 'error'" type="error">
            {{ scope.row.level }}
          </n-text>
          <n-text v-else>{{ scope.row.level }}</n-text>
        </template>
      </el-table-column>
      <el-table-column prop="msg" label="信息">
        <template #default="scope">
          <span v-if="scope.row.msg.startsWith('onebot | ')" style="color: var(--el-color-warning)">
            {{ scope.row.msg }}
          </span>
          <span v-else-if="scope.row.msg.startsWith('发给')" style="color: var(--el-color-primary)">
            {{ scope.row.msg }}
          </span>
          <span v-else-if="scope.row.level === 'warn'" style="color: var(--el-color-warning)">
            {{ scope.row.msg }}
          </span>
          <span v-else-if="scope.row.level === 'error'" style="color: var(--el-color-danger)">
            {{ scope.row.msg }}
          </span>
          <span v-else>{{ scope.row.msg }}</span>
        </template>
      </el-table-column>
    </el-table>
  </div>
  <el-table
    :data="store.curDice.logs"
    class="logs w-full md:hidden"
    :row-class-name="getLogRowClassName"
    :header-cell-style="{ backgroundColor: '#f3f5f7' }">
    <el-table-column label="时间" width="60">
      <template #default="scope">
        <div style="display: flex; align-items: center">
          <span
            v-if="scope.row.msg.startsWith('onebot | ')"
            style="color: var(--el-color-warning)"
            >{{ dayjs.unix(scope.row.ts).format('HH:mm') }}</span
          >
          <span
            v-else-if="scope.row.msg.startsWith('发给')"
            style="color: var(--el-color-primary)"
            >{{ dayjs.unix(scope.row.ts).format('HH:mm') }}</span
          >
          <span v-else-if="scope.row.level === 'warn'" style="color: var(--el-color-warning)">{{
            dayjs.unix(scope.row.ts).format('HH:mm')
          }}</span>
          <span v-else-if="scope.row.level === 'error'" style="color: var(--el-color-danger)">{{
            dayjs.unix(scope.row.ts).format('HH:mm')
          }}</span>
          <span v-else>{{ dayjs.unix(scope.row.ts).format('HH:mm') }}</span>
        </div>
      </template>
    </el-table-column>
    <el-table-column prop="msg" label="信息">
      <template #default="scope">
        <span v-if="scope.row.msg.startsWith('onebot | ')" style="color: var(--el-color-warning)">{{
          scope.row.msg
        }}</span>
        <span v-else-if="scope.row.msg.startsWith('发给')" style="color: var(--el-color-primary)">{{
          scope.row.msg
        }}</span>
        <span v-else-if="scope.row.level === 'warn'" style="color: var(--el-color-warning)">{{
          scope.row.msg
        }}</span>
        <span v-else-if="scope.row.level === 'error'" style="color: var(--el-color-danger)">{{
          scope.row.msg
        }}</span>
        <span v-else>{{ scope.row.msg }}</span>
      </template>
    </el-table-column>
  </el-table>

  <n-modal
    v-model:show="upgradeDialogVisible"
    preset="card"
    title="升级新版本"
    :mask-closable="false"
    :close-on-esc="false"
    :closable="false"
    class="the-dialog">
    <n-button
      text
      tag="a"
      style="font-size: 16px; font-weight: bolder"
      type="info"
      href="https://dice.weizaima.com/changelog"
      target="_blank">
      查看更新日志
    </n-button>

    <div>请及时更新海豹到最新版本，这意味着功能增加和 BUG 修复。</div>
    <div>当然，在更新前最好看看右上角的海豹新闻，通常会很有帮助。</div>
    <div>在操作之前，最好能确保你目前可以接触到服务器，以防万一需要人工干预。</div>
    <div>
      <b>如果升级后无法启动，请删除海豹目录中的"update"、"auto_update.exe"并手动进行升级</b>
    </div>
    <div><b>进一步的内容请查阅届时自动生成的“升级失败指引”或加群询问。</b></div>

    <n-button style="margin: 1rem 0" type="primary" @click="doUpgrade">
      确认升级到 {{ store.curDice.baseInfo.versionNew }}
    </n-button>

    <div>{{ store.curDice.baseInfo.versionNewNote }}</div>
    <div>注意：升级成功后界面不会自动刷新，请在重连完成后手动刷新</div>
    <div><b>当前 Win11 22H2 无法自动重启，建议此系统用户手动更新</b></div>
    <div>不要连续多次执行</div>

    <template #footer>
      <span class="dialog-footer">
        <!-- <el-button @click="dialogImportVisible = false">取消</el-button> -->
        <!-- <el-button @click="configForImport = ''">清空</el-button> -->
        <!-- <el-button data-clipboard-target="#import-edit" @click="copied" id="btnCopy1">复制</el-button> -->
        <!-- <el-button type="primary" @click="doImport" :disabled="configForImport === ''">导入并保存</el-button> -->
      </span>
    </template>
  </n-modal>

  <!-- <div v-for="i in store.curDice.logs">
    {{i}}
  </div> -->
</template>

<script lang="ts" setup>
import { useStore } from '~/store';
import dayjs from 'dayjs';
import { filesize } from 'filesize';
import { postUpgrade } from '~/api/v1/dice';
import { useDialog } from 'naive-ui';

const store = useStore();
const dialog = useDialog();

const upgradeDialogVisible = ref(false);
const autoRefresh = ref(true);
const now = ref<dayjs.Dayjs>(dayjs());

let timerId: number;

const doUpgrade = async () => {
  upgradeDialogVisible.value = false;
  dialog.info({
    title: '升级',
    content: '开始下载更新，请等待……\n完成后将自动重启海豹，并进入更新流程',
  });
  try {
    const ret = await postUpgrade();
    dialog.info({
      title: '升级',
      content: ret.text + '\n如果几分钟后服务没有恢复，检查一下海豹目录',
    });
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
  } catch (e) {
    // ElMessageBox.alert('升级失败', '升级')
  }
};

const scrollDown = () => {
  const panel = document.querySelector<HTMLElement>('.logs')?.parentElement;
  if (panel) {
    panel.scrollTop = panel.scrollHeight;
  }
};

const getLogRowClassName = ({ row }: { row: any }) => {
  switch (row.level) {
    case 'warn':
      return 'no-hover warning-row';
    case 'error':
      return 'no-hover danger-row';
    default:
      return 'no-hover normal-row';
  }
};

onBeforeMount(async () => {
  if (autoRefresh.value) {
    await store.logFetchAndClear();
  }

  timerId = setInterval(() => {
    if (autoRefresh.value) {
      store.logFetchAndClear();
    }
    now.value = dayjs();
  }, 5000) as any;
});

onBeforeUnmount(() => {
  clearInterval(timerId);
});
</script>

<style scoped>
.btn-scrolldown {
  position: absolute;
  right: 40px;
  bottom: 60px;
  width: 2rem;
  height: 2rem;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1rem;
  box-shadow: 0 0 6px rgba(0, 0, 0, 0.12);
  cursor: pointer;
  z-index: 5;
  opacity: 0.4;
}

.btn-scrolldown:hover {
  transition: all 0.3s;
  opacity: 1;
}

.latest-log-warn {
  margin-top: 0;
  margin-bottom: 1rem;
  :deep(.el-divider__text) {
    background: #f3f4f6;
  }
}
</style>

<style lang="css">
.el-table .warning-row {
  --el-table-tr-bg-color: var(--el-color-warning-light-8);
  &:hover {
    --el-table-tr-bg-color: var(--el-color-warning-light-9);
  }
}

.el-table .danger-row {
  --el-table-tr-bg-color: var(--el-color-danger-light-8);
  &:hover {
    --el-table-tr-bg-color: var(--el-color-danger-light-9);
  }
}

.el-table .normal-row {
  --el-table-tr-bg-color: #f3f5f7;
  &:hover {
    --el-table-tr-bg-color: var(--el-color-primary-light-9);
  }
}

.no-hover:hover > td {
  background-color: initial !important;
}
</style>

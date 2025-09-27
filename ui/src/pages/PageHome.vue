<template>
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
    <div class="flex items-center justify-end">
      <n-checkbox v-model:checked="displayReverse">最新日志展示在最上方</n-checkbox>
      <n-checkbox v-model:checked="autoRefresh">保持刷新</n-checkbox>
    </div>
  </div>

  <main class="logs pb-8">
    <n-data-table
      :data="logData"
      :class="isMobile ? 'w-full' : ''"
      :columns="columns"
    />

    <n-back-top :right="30" />
  </main>

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

<script lang="tsx" setup>
import { useStore } from '~/store';
import dayjs from 'dayjs';
import { filesize } from 'filesize';
import { postUpgrade } from '~/api/v1/dice';
import { useDialog, type DataTableColumns } from 'naive-ui';
import { useThemeVars } from 'naive-ui'
import { breakpointsTailwind } from '@vueuse/core'

const themeVars = useThemeVars()
const breakpoints = useBreakpoints(breakpointsTailwind)
const isMobile = breakpoints.smaller('md')

type LogRow = {
  ts: number;
  msg: string;
  level: 'info' | 'warn' | 'error' | string;
};

const store = useStore();
const dialog = useDialog();

const displayReverse = ref<boolean>(true);
const logData = computed(() => {
  if (displayReverse.value) {
    return [...store.curDice.logs].reverse()
  }
  return store.curDice.logs
})

const upgradeDialogVisible = ref(false);
const autoRefresh = ref(true);
const now = ref<dayjs.Dayjs>(dayjs());

let timerId: number;

const getMsgColor = (row: LogRow): string | undefined => {
  if (row.msg.startsWith('onebot | ')) return themeVars.value.warningColor;
  if (row.msg.startsWith('发给')) return themeVars.value.infoColor;
  if (row.level === 'warn') return themeVars.value.warningColor;
  if (row.level === 'error') return themeVars.value.errorColor;
  return undefined;
};

const columns: ComputedRef<DataTableColumns<LogRow>> = computed(() => {
  let data = []
  data.push({
    title: '时间',
    key: 'ts',
    width: isMobile.value ? 70 : 100,
    render: (row) => {
      const color = getMsgColor(row);
      return (
        <div class="flex items-center" style={{ color }}>
          { isMobile.value ? <></> : <n-icon><i-carbon-time/></n-icon> }
          <span class="ml-1">
            { dayjs.unix(row.ts).format(isMobile.value ? 'HH:mm' : 'HH:mm:ss') }
          </span>
        </div>
      );
    }
  })
  if (!isMobile.value) {
    data.push({
      title: '级别',
      key: 'level',
      width: 55,
      render: (row) => {
        const color = getMsgColor(row);
        return <span style={{ color }}>{row.level}</span>;
      },
    })
  }
  data.push({
    title: '信息',
    key: 'msg',
    render: (row) => {
      const color = getMsgColor(row);
      return <span style={{ color }}>{row.msg}</span>;
    },
  })
  return data;
});

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
  console.log('scrollDown: ', panel)
  if (panel) {
    panel.scrollTop = panel.scrollHeight;
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
  }, 5000);
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
</style>

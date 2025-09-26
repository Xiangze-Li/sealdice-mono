<template>
  <header class="page-header">
    <n-button type="primary" @click="doBackup">
      <template #icon>
        <n-icon><i-carbon-renew /></n-icon>
      </template>
      重载牌堆
    </n-button>
  </header>

  <n-tabs v-model:value="mode" justify-content="space-evenly">
    <n-tab-pane tab="牌堆列表" name="list">
      <header class="mb-4 flex flex-wrap-reverse items-center justify-between gap-y-2">
        <n-flex align="center">
          <span>
            <n-input v-model:value="filter" size="small" clearable>
              <template #prefix>
                <n-icon><i-carbon-search /></n-icon>
              </template>
            </n-input>
          </span>
        </n-flex>
        <n-flex size="small" align="center" class="ml-auto">
          <n-button
            type="info"
            size="tiny"
            text
            tag="a"
            target="_blank"
            href="https://github.com/sealdice/draw">
            <template #icon>
              <n-icon><i-carbon-link /></n-icon>
            </template>
            获取牌堆
          </n-button>
          <span>
            <n-upload
              class="upload"
              action=""
              multiple
              :show-file-list="false"
              @before-upload="beforeUpload"
              :file-list="fileList">
              <n-button type="info" secondary>
                <template #icon>
                  <n-icon><i-carbon-upload /></n-icon>
                </template>
                上传牌堆
              </n-button>
            </n-upload>
          </span>
        </n-flex>
      </header>
      <aside class="mb-4 flex items-center">
        <n-text v-if="filterCount > 0" class="text-xs" type="info">
          已过滤 {{ filterCount }} 条
        </n-text>
        <n-flex size="small" align="center" class="ml-auto">
          <n-text class="text-xs">目前支持 json/yaml/deck/toml 格式的牌堆</n-text>
          <n-tooltip>
            <template #trigger>
              <n-icon size="small"><i-carbon-help-filled /></n-icon>
            </template>
            deck 牌堆：一种单文件带图的牌堆格式<br />
            在牌堆文件中使用./images/xxx.png 的相对路径引用图片。并连同图片目录一起打包成
            zip，修改扩展名为 deck 即可制作<br />
            <br />
            toml 牌堆：海豹支持的新牌堆格式。格式更加友好，还提供了包括云牌组在内的更多功能支持。
          </n-tooltip>
        </n-flex>
      </aside>
      <main class="deck-list-main">
        <foldable-card
          v-for="(i, index) in filtered"
          :key="index"
          class="deck-item"
          :err-title="i.filename"
          :err-text="i.errText">
          <template #title>
            <n-flex size="small" align="center">
              <n-text class="text-base" tag="b">{{ i.name }}</n-text>
              <n-text>{{ i.version }}</n-text>
              <n-tag
                size="small"
                :type="i.fileFormat === 'toml' ? 'success' : 'info'"
                :bordered="false">
                {{ i.fileFormat }}
              </n-tag>
            </n-flex>
          </template>

          <template #title-extra>
            <n-flex>
              <n-popconfirm
                v-if="i.updateUrls && i.updateUrls.length > 0"
                @positive-click="doCheckUpdate(i)">
                <template #trigger>
                  <n-button type="info" size="small" secondary :loading="diffLoading">
                    <template #icon>
                      <n-icon><i-carbon-download /></n-icon>
                    </template>
                    更新
                  </n-button>
                </template>
                更新地址由牌堆作者提供，是否确认要检查该牌堆更新？
              </n-popconfirm>
              <n-button type="error" size="small" secondary @click="doDelete(i)">
                <template #icon>
                  <n-icon><i-carbon-row-delete /></n-icon>
                </template>
                删除
              </n-button>
            </n-flex>
          </template>

          <template #title-extra-error>
            <n-button type="error" size="small" secondary @click="doDelete(i)">
              <template #icon>
                <n-icon><i-carbon-row-delete /></n-icon>
              </template>
              删除
            </n-button>
          </template>

          <template #description>
            <n-flex size="small" vertical align="normal">
              <n-text v-if="i.cloud" type="info" class="text-xs">
                <n-icon><i-carbon-cloud /></n-icon>
                作者提供云端内容，请自行鉴别安全性
              </n-text>
              <n-text v-if="i.fileFormat === 'jsonc'" type="warning" class="text-xs">
                <n-icon><i-carbon-warning-filled /></n-icon>
                注意：该牌堆的格式并非标准 JSON，而是允许尾逗号与注释语法的扩展 JSON
              </n-text>
            </n-flex>
          </template>

          <n-descriptions content-class="whitespace-pre-line">
            <n-descriptions-item :span="3" label="作者">
              {{ i.author || '&lt;佚名>' }}
            </n-descriptions-item>
            <n-descriptions-item v-if="i.desc" :span="3" label="简介">
              {{ i.desc }}
            </n-descriptions-item>
            <n-descriptions-item :span="3" label="牌组列表">
              <n-flex size="small">
                <n-tag
                  v-for="(visible, c) of i.command"
                  :key="c"
                  size="small"
                  :type="visible ? 'info' : 'default'"
                  :bordered="false">
                  {{ c }}
                </n-tag>
              </n-flex>
            </n-descriptions-item>
            <n-descriptions-item v-if="i.license" label="许可协议">
              {{ i.license }}
            </n-descriptions-item>
            <n-descriptions-item v-if="i.date" label="发布时间">
              {{ i.date }}
            </n-descriptions-item>
            <n-descriptions-item v-if="i.updateDate" label="更新时间">
              {{ i.updateDate }}
            </n-descriptions-item>
          </n-descriptions>

          <template #unfolded-extra>
            <n-descriptions content-class="whitespace-pre-line">
              <n-descriptions-item :span="3" label="可见牌组列表">
                <n-flex size="small">
                  <n-tag
                    v-for="(visible, c) of i.command"
                    :key="c"
                    size="small"
                    :type="visible ? 'info' : 'default'"
                    :style="{
                      display: visible ? '' : 'none',
                    }"
                    :bordered="false">
                    {{ c }}
                  </n-tag>
                </n-flex>
              </n-descriptions-item>
            </n-descriptions>
          </template>
        </foldable-card>
      </main>

      <n-modal v-model:show="showDiff" preset="card" title="牌堆内容对比" class="diff-dialog">
        <diff-viewer :lang="deckCheck.format" :old="deckCheck.old" :new="deckCheck.new" />
        <template #footer>
          <n-flex wrap>
            <n-button @click="showDiff = false">取消</n-button>
            <n-button v-if="!(deckCheck.old === deckCheck.new)" type="success" @click="deckUpdate">
              <template #icon>
                <n-icon><i-carbon-save /></n-icon>
              </template>
              确认更新
            </n-button>
          </n-flex>
        </template>
      </n-modal>
    </n-tab-pane>
  </n-tabs>
</template>

<script lang="ts" setup>
import { getBackupConfig } from '~/api/v1/backup';
import {
  checkDeckUpdate,
  deleteDeck,
  getDeckList,
  reloadDeck,
  updateDeck,
  uploadDeck,
} from '~/api/v1/deck';
import { type UploadFileInfo, useDialog, useMessage } from 'naive-ui';

const message = useMessage();
const dialog = useDialog();
const mode = ref<string>('list');

const filter = ref<string>('');
const filterCount = computed(() => data.value.length - filtered.value.length);
const data = ref<any[]>([]);
const filtered = computed(() =>
  data.value.filter(deck => {
    if (filter.value === '') {
      return true;
    }
    const val = filter.value.toLowerCase();
    return (
      deck.name?.toLowerCase()?.includes(val) ||
      deck.desc?.toLowerCase()?.includes(val) ||
      deck.author?.toLowerCase()?.includes(val) ||
      Object.keys(deck.command)
        .map(tag => tag?.toLowerCase()?.includes(val))
        .includes(true)
    );
  }),
);

const cfg = ref<any>({});

const refreshList = async () => {
  const lst = await getDeckList();
  data.value = lst;
};

const configGet = async () => {
  const data = await getBackupConfig();
  cfg.value = data;
};

const fileList = ref<any[]>([]);

const doBackup = async () => {
  const ret = await reloadDeck();
  if (ret.testMode) {
    message.success('展示模式无法重载牌堆');
  } else {
    message.success('已重载');
    await refreshList();
  }
};

const doDelete = async (data: any) => {
  dialog.warning({
    title: '确认删除',
    content: `确认删除牌堆「${data.name}」，确定吗？`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      await deleteDeck(data.filename);
      await reloadDeck();
      await refreshList();
      message.success('牌堆已删除');
    },
  });
};

const beforeUpload = async (data: { file: UploadFileInfo }) => {
  const file = data.file.file as File;
  // UploadRawFile
  await uploadDeck(file);
  message.success('上传完成，即将自动重载牌堆');
  await reloadDeck();
  await refreshList();
};

onBeforeMount(async () => {
  await configGet();
  await refreshList();
});

const showDiff = ref<boolean>(false);
const diffLoading = ref<boolean>(false);

interface DeckCheckResult {
  old: string;
  new: string;
  format: 'json' | 'yaml' | 'toml';
  filename: string;
  tempFileName: string;
}

const deckCheck = ref<DeckCheckResult>({
  old: '',
  new: '',
  format: 'json',
  filename: '',
  tempFileName: '',
});

const doCheckUpdate = async (data: any) => {
  diffLoading.value = true;
  const checkResult = await checkDeckUpdate(data.filename);
  diffLoading.value = false;
  if (checkResult.result) {
    deckCheck.value = { ...checkResult, filename: data.filename };
    showDiff.value = true;
  } else {
    message.error('检查更新失败！' + checkResult.err);
  }
};

const deckUpdate = async () => {
  const res = await updateDeck(deckCheck.value.filename, deckCheck.value.tempFileName);
  if (res.result) {
    showDiff.value = false;
    message.success('更新成功，即将自动重载牌堆');
    await reloadDeck();
    await refreshList();
  } else {
    showDiff.value = false;
    message.error('更新失败！' + res.err);
  }
};
</script>

<style lang="css">
@media screen and (max-width: 700px) {
  .diff-dialog {
    width: 90% !important;
  }
}

@media screen and (min-width: 700px) and (max-width: 900px) {
  .diff-dialog {
    width: 80% !important;
  }
}

@media screen and (min-width: 900px) and (max-width: 1100px) {
  .diff-dialog {
    width: 65% !important;
  }
}

@media screen and (min-width: 1100px) {
  .diff-dialog {
    width: 50% !important;
  }
}

.deck-list-main {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
}

.deck-item {
  width: 100%;
}

.upload {
  > ul {
    display: none;
  }
}
</style>

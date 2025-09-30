<template>
  <h4>词库列表</h4>
  <header class="page-header">
    <n-upload
      action=""
      multiple
      accept="application/text,.txt,application/toml,.toml"
      @before-upload="beforeUpload"
      :show-file-list="false">
      <n-button type="info" secondary>
        <template #icon>
          <n-icon><i-carbon-upload /></n-icon>
        </template>
        导入
      </n-button>
    </n-upload>
    <n-flex>
      <n-button
        type="primary"
        size="tiny"
        text
        tag="a"
        target="_blank"
        :href="`${urlBase}/sd-api/v1/censor/files/template/toml`">
        <template #icon>
          <n-icon><i-carbon-download /></n-icon>
        </template>
        下载 toml 词库模板
      </n-button>
      <n-button
        type="primary"
        size="tiny"
        text
        tag="a"
        target="_blank"
        :href="`${urlBase}/sd-api/v1/censor/files/template/txt`">
        <template #icon>
          <n-icon><i-carbon-save /></n-icon>
        </template>
        下载 txt 词库模板
      </n-button>
    </n-flex>
  </header>
  <main style="margin-top: 1rem">
    <n-data-table :columns="columns" :data="files" />
  </main>
</template>

<script setup lang="tsx">
import { urlBase } from '~/backend';
import { useCensorStore } from '~/components/censor/censor';
import SensitiveTag from '~/components/censor/sensitive-tag.tsx';
import { deleteCensorFiles, getCensorFiles, uploadCensorFile } from '~/api/v1/censor';
import { useMessage, type UploadFileInfo, type DataTableColumns } from 'naive-ui';

const message = useMessage();

onBeforeMount(() => {
  refreshFiles();
});

const censorStore = useCensorStore();

interface SensitiveWordFile {
  key: string;
  path: string;
  counter: number[];
}

const columns: DataTableColumns<SensitiveWordFile> = [
  {
    title: '文件名',
    key: 'name',
  },
  {
    title: () => <SensitiveTag type="default" />,
    key: 'count[1]',
  },
  {
    title: () => <SensitiveTag type="info" />,
    key: 'count[2]',
  },
  {
    title: () => <SensitiveTag type="warning" />,
    key: 'count[3]',
  },
  {
    title: () => <SensitiveTag type="error" />,
    key: 'count[4]',
  },
];

const files = ref<SensitiveWordFile[]>();

censorStore.$subscribe(async (_, state) => {
  if (state.filesNeedRefresh === true) {
    await refreshFiles();
    state.filesNeedRefresh = false;
  }
});

const refreshFiles = async () => {
  const c:
    | { result: false }
    | {
        result: true;
        data: SensitiveWordFile[];
      } = await getCensorFiles();
  if (c.result) {
    files.value = c.data;
  }
};

const beforeUpload = async (data: { file: UploadFileInfo }) => {
  const file = data.file.file;
  const c = await uploadCensorFile(file as File);
  if (c.result) {
    await refreshFiles();
    message.success('上传完成，请在全部操作完成后，手动重载拦截');
    censorStore.markReload();
  } else {
    message.error('上传失败！' + c.err);
  }
};

const deleteFile = async (key: string) => {
  await ElMessageBox.confirm('是否删除此词库？', '删除', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  }).then(async () => {
    const c: { result: true } | { result: false; err: string } = await deleteCensorFiles([key]);
    if (c.result) {
      message.success('删除词库完成，请在全部操作完成后，手动重载拦截');
      censorStore.markReload();
    } else {
      message.error('删除词库失败！' + c.err);
    }
  });
};
</script>

<script setup lang="ts">
import type { Resource } from '~/store';
import { useStore } from '~/store';
import { filesize } from 'filesize';
import { urlBase } from '~/backend';
import ClipboardJS from 'clipboard';
import {
  createResource,
  deleteResource as deleteResourceApi,
  getResourcePage,
} from '~/api/v1/resource';
import { useDialog, useMessage, type UploadFileInfo } from 'naive-ui';

const dialog = useDialog();
const message = useMessage();

const store = useStore();

const loading = ref(true);
const images = ref<Resource[]>([]);

const refreshResources = async () => {
  loading.value = true;
  const imagesRes = await getResourcePage('image');
  if (imagesRes.result) {
    images.value = imagesRes.data;
  } else {
    message.error(imagesRes.err);
  }
  loading.value = false;
};

const deleteResource = async (resource: Resource) => {
  dialog.warning({
    title: '确认删除',
    content: `确认删除「${resource.name}（${resource.path}）」吗？删除后将无法找回`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      const res = await deleteResourceApi(resource.path);
      if (res.result) {
        message.success('删除成功');
      } else {
        message.error(res.err);
      }
      await nextTick(async () => {
        await refreshResources();
      });
    },
  });
};

const beforeUpload = async (data: { file: UploadFileInfo }) => {
  // UploadRawFile
  const file = data.file.file as File;
  if (file.type !== 'image/jpeg' && file.type !== 'image/png' && file.type !== 'image/gif') {
    message.error('上传的文件不是图片！');
    return false;
  }
  const fd = new FormData();
  fd.append('files', file);
  try {
    const resp = await createResource(file);
    if (resp.result) {
      message.success('上传完成');
    } else {
      message.error(resp.err);
    }
  } catch (e: any) {
    message.error(e.toString());
  } finally {
    await nextTick(async () => {
      await refreshResources();
    });
  }
};

const copySealCode = async () => {
  message.success('复制海豹码成功！');
};

onBeforeMount(async () => {
  loading.value = true;
  await refreshResources();
  new ClipboardJS('.resource-seal-code-copy-btn');
});
</script>

<template>
  <h2>资源管理</h2>
  <tip-box>
    <span>此处可以上传图片，方便引用。</span>
  </tip-box>

  <main>
    <h3 class="flex items-center justify-between">
      <span>图片列表</span>
      <span>
        <n-upload
          action=""
          multiple
          accept=".png, .jpg, jpeg, .gif"
          @before-upload="beforeUpload"
          :show-file-list="false">
          <n-button type="primary">
            <template #icon>
              <i-carbon-upload />
            </template>
            上传图片
          </n-button>
        </n-upload>
      </span>
    </h3>

    <n-list>
      <n-list-item v-for="image in images" :key="image.path">
        <template #prefix>
          <resource-render :key="image.path" class="ml-4 min-w-16" :data="image" mini />
        </template>

        <n-thing :title="image.name">
          <template #description>
            <n-tag :bordered="false" type="info" size="small">
              {{ filesize(image.size) }}
            </n-tag>
          </template>
          <n-text code class="break-all">{{ image.path }}</n-text>
        </n-thing>

        <template #suffix>
          <n-flex size="small" justify="end" class="mr-4">
            <n-button
              v-if="image.type === 'image'"
              type="info"
              quaternary
              size="tiny"
              class="resource-seal-code-copy-btn"
              :data-clipboard-text="`[图:${image.path}]`"
              @click="copySealCode()">
              <template #icon>
                <i-carbon-copy />
              </template>
              复制海豹码
            </n-button>
            <n-button
              type="success"
              size="tiny"
              quaternary
              tag="a"
              style="text-decoration: none"
              :href="`${urlBase}/sd-api/resource/download?path=${encodeURIComponent(image.path)}&token=${encodeURIComponent(store.token)}`">
              <template #icon>
                <i-carbon-download />
              </template>
              下载
            </n-button>
            <n-button type="error" quaternary size="tiny" @click="deleteResource(image)">
              <template #icon>
                <i-carbon-row-delete />
              </template>
              删除
            </n-button>
          </n-flex>
        </template>
      </n-list-item>
    </n-list>
  </main>
</template>

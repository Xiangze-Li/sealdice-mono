<template>
  <header class="page-header">
    <n-button type="primary" :loading="reloadLoading" :disabled="reloadLoading" @click="reload">
      <template #icon>
        <n-icon><i-carbon-renew /></n-icon>
      </template>
      重载帮助文档
    </n-button>
  </header>

  <n-affix v-if="needReload" :top="60">
    <tip-box type="error">
      <n-text type="error" class="text-base" tag="strong">存在修改，需要重载后生效！</n-text>
    </tip-box>
  </n-affix>

  <n-affix v-if="configNeedSave" :top="60">
    <tip-box type="error">
      <n-text type="error" class="text-base" tag="strong">设置存在修改，别忘记保存！</n-text>
    </tip-box>
  </n-affix>

  <n-tabs v-model:value="tab" justify-content="space-evenly" class="pb-8">
    <n-tab-pane tab="文件" name="file">
      <n-flex justify="end" class="mb-4 flex-1">
        <n-button v-show="showDeleteFile" type="error" secondary @click="deleteFiles">
          <template #icon>
            <n-icon><i-carbon-row-delete /></n-icon>
          </template>
          删除所选
        </n-button>
        <n-button type="info" secondary @click="uploadDialogVisible = true">
          <template #icon>
            <n-icon><i-carbon-upload /></n-icon>
          </template>
          上传
        </n-button>
        <n-button type="info" secondary @click="configDialogVisible = true">
          <template #icon>
            <n-icon><i-carbon-settings /></n-icon>
          </template>
          设置
        </n-button>
      </n-flex>

      <n-modal
        v-model:show="uploadDialogVisible"
        preset="card"
        title="上传帮助文档"
        :closable="false">
        <n-alert v-show="uploadForm.group === 'default'" type="warning" class="mb-4">
          更具体的分组能提供组内搜索命令
          <n-tag size="small" :bordered="false">.find#&lt;分组&gt; &lt;搜索内容&gt;</n-tag>
          ，是否一定要选择默认分组？
        </n-alert>
        <n-form ref="uploadFormRef" v-model:model="uploadForm" :rules="uploadRules">
          <n-form-item label="分组" prop="group">
            <n-select
              v-model:value="uploadForm.group"
              placeholder="选择分组"
              filterable
              clearable
              :options="docGroups"
              tag />
          </n-form-item>
          <n-form-item label="帮助文档" prop="files">
            <n-upload
              @on-change="fileChange"
              :file-list="uploadForm.files"
              :show-file-list="false"
              multiple
              accept=".json, .xlsx">
              <n-button type="primary">
                <template #icon>
                  <n-icon><i-carbon-upload /></n-icon>
                </template>
                选择文件
              </n-button>
            </n-upload>
          </n-form-item>
        </n-form>
        <template #footer>
          <n-flex>
            <n-button @click="uploadDialogVisible = false">取消</n-button>
            <n-button type="primary" @click="submitUpload(uploadFormRef)">上传</n-button>
          </n-flex>
        </template>
      </n-modal>

      <n-modal
        v-model:show="configDialogVisible"
        preset="card"
        title="设置帮助文档"
        :closable="false">
        <n-alert v-show="configNeedSave" type="warning" class="mb-4">
          设置存在修改，别忘记保存！
        </n-alert>
        <h4>分组别名</h4>
        <n-form label-placement="left">
          <n-form-item v-for="group in docGroups" :key="group.value" :label="group.label">
            <HelpConfigTags
              :group="group"
              :aliases="helpAliases"
              @add-alias="addAlias"
              @remove-alias="removeAlias" />
          </n-form-item>
        </n-form>
        <template #footer>
          <n-flex justify="end">
            <n-button @click="configDialogClose">取消</n-button>
            <n-button type="primary" @click="summitConfig">保存</n-button>
          </n-flex>
        </template>
      </n-modal>

      <main>
        <header class="file-tree-title">
          <n-text size="large" tag="b">文件名</n-text>
          <n-text size="large" tag="b">分组</n-text>
        </header>
        <n-tree
          ref="fileTreeRef"
          :data="docTree"
          label-field="name"
          block-line
          default-expand-all
          :show-line="true"
          cascade
          checkable />
      </main>
    </n-tab-pane>

    <n-tab-pane tab="词条" name="item">
      <main class="item-list-container">
        <header>
          <n-form
            v-model:model="textItemQuery"
            size="small"
            class="flex flex-wrap"
            label-width="auto"
            label-placement="left"
            inline>
            <n-form-item label="序号">
              <n-input-number v-model:value="textItemQuery.id" placeholder="" clearable />
            </n-form-item>
            <n-form-item label="分组" j>
              <n-select
                v-model:value="textItemQuery.group"
                placeholder="选择分组"
                filterable
                clearable
                :options="[{ value: 'builtin', label: '内置' }, ...docGroups]">
              </n-select>
            </n-form-item>
            <n-form-item label="来源文件">
              <n-input v-model:value="textItemQuery.from" placeholder="" clearable />
            </n-form-item>
            <n-form-item label="词条名">
              <n-input v-model:value="textItemQuery.title" placeholder="" clearable />
            </n-form-item>
            <n-form-item>
              <n-button type="info" secondary @click="refreshTextItems">搜索</n-button>
            </n-form-item>
          </n-form>
        </header>
        <n-data-table class="item-list" :columns="columns" :data="textItems" size="small" />
        <footer>
          <n-pagination
            class="item-list-pagination"
            v-model:page="textItemQuery.pageNum"
            v-model:page-size="textItemQuery.pageSize"
            show-size-picker
            :page-sizes="[10, 20, 30, 50]"
            :page-slot="5"
            :item-count="textItemQuery.total"
            @update:page="handleCurrentPageChange"
            @update:page-size="handleCurrentPageChange" />
        </footer>
      </main>
    </n-tab-pane>
  </n-tabs>
</template>

<script lang="tsx" setup>
import type { FormRules, DataTableColumns, TreeOption, NTree } from 'naive-ui';
import { useMessage, useDialog } from 'naive-ui';
import { trim } from 'es-toolkit/compat';
import type { HelpDoc, HelpTextItem, HelpTextItemQuery } from '#';
import {
  deleteHelpDoc,
  getHelpDocConfig,
  getHelpDocPage,
  getHelpDocTree,
  postHelpDocConfig,
  reloadHelpDoc,
  uploadHelpDoc,
} from '~/api/v1/helpdoc';

interface Group {
  label: string;
  value: string;
}

const tab = ref('file');

const message = useMessage();
const dialog = useDialog();

const needReload = ref<boolean>(false);

interface UploadForm {
  group: string;
  files: any[];
}

const docTree = ref<TreeOption[]>([]);
const docGroups = ref<Group[]>([]);
const uploadDialogVisible = ref<boolean>(false);
const uploadFormRef = ref<FormInstance>();
const uploadForm = reactive<UploadForm>({
  group: '',
  files: [] as any[],
});
const uploadRules = reactive<FormRules>({
  group: [
    { required: true, message: '请选择分组', trigger: 'blur-sm' },
    { pattern: '^(?!builtin).*', message: '不能为内置分组', trigger: 'blur-sm' },
  ],
  files: [{ required: true, message: '请选择文件', trigger: 'blur-sm' }],
});

const fileChange = (_f: any, newFiles: any) => {
  uploadForm.files = newFiles;
};

const submitUpload = async (formData: FormInstance | undefined) => {
  if (!formData) return;
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  await formData.validate(async (valid, _) => {
    if (valid) {
      // const fd = new FormData();
      // fd.append("group", uploadForm.group)
      // uploadForm.files.forEach(f => {
      //   fd.append("files", f.raw)
      // })

      const res = await uploadHelpDoc(uploadForm);
      if (res.result) {
        message.success('上传完成，请在全部操作完成后，手动重载帮助文件');
      } else {
        message.error(res.err ?? '上传失败');
      }
      formData.resetFields();
      needReload.value = true;
      await refreshFileTree();
      uploadDialogVisible.value = false;
    }
  });
};

const fileTreeRef = ref<InstanceType<typeof NTree>>();
const checkedFileKeys = computed(() => {
  return (fileTreeRef.value?.getCheckedData().keys as string[]) ?? [];
});
const showDeleteFile = computed(() => {
  return checkedFileKeys.value.length > 0;
});
const deleteFiles = async () => {
  if (checkedFileKeys && checkedFileKeys.value.length !== 0) {
    dialog.warning({
      title: '删除',
      content: '确认删除选择的文件吗？',
      positiveText: '确定',
      negativeText: '取消',
      onPositiveClick: () => {
        deleteHelpDoc(checkedFileKeys.value).then(res => {
          if (res.result) {
            message.success('删除文件成功');
            refreshFileTree();
          } else {
            message.error(res.err ?? '删除文件失败');
          }
        });
        needReload.value = true;
      },
    });
  } else {
    message.error('未选择文件');
  }
};

const getHelpDocTag = (
  loadStatus: number,
  deleted: boolean,
  group: string,
):
  | { type: 'primary'; label: string }
  | { type: 'success'; label: string }
  | { type: 'warning'; label: string }
  | { type: 'info'; label: string }
  | { type: 'error'; label: string } => {
  if (loadStatus === 0) {
    return {
      type: 'warning',
      label: '未加载',
    };
  } else if (loadStatus === 2) {
    return {
      type: 'error',
      label: '格式有误',
    };
  } else if (deleted) {
    return {
      type: 'warning',
      label: group,
    };
  } else {
    return {
      type: 'success',
      label: group,
    };
  }
};

const convertHelpDocInfo = (doc: HelpDoc): TreeOption => {
  const option = {
    ...doc,
    children: doc.children?.map(e => {
      return convertHelpDocInfo(e);
    }),
    prefix: () => {
      if (doc.isDir) {
        return <i-bi-folder2 color="var(--color-stone-600)" />;
      } else if (doc.type === '.json') {
        return <i-bi-filetype-json color="var(--color-amber-600)" />;
      } else if (doc.type === '.xlsx') {
        return <i-bi-filetype-xlsx color="var(--color-green-600)" />;
      } else {
        return <i-bi-file-break />;
      }
    },
    suffix: () => {
      if (doc.isDir) {
        return <></>;
      }
      const tagInfo = getHelpDocTag(doc.loadStatus, doc.deleted, doc.group);
      return (
        <n-tag class="ml-auto" size="small" type={tagInfo.type} bordered={false}>
          {tagInfo.label}
        </n-tag>
      );
    },
  };

  return option;
};

const refreshFileTree = async () => {
  const resp = await getHelpDocTree();
  if (resp?.result) {
    docTree.value = resp.data.map(entry => {
      return convertHelpDocInfo(entry);
    });
    docGroups.value = [
      { label: '默认', value: 'default' },
      ...resp.data
        .filter(entry => {
          return entry.isDir;
        })
        .map(entry => {
          return { label: entry.name, value: entry.name };
        }),
    ];
  } else {
    message.error('无法获取帮助文件信息');
  }
};

const textItems = ref<HelpTextItem[]>([]);
const textItemQuery = ref<HelpTextItemQuery>({
  pageNum: 1,
  pageSize: 20,
  total: 20,
});
const getHelpText = (row: string) => {
  const temp = trim(row);
  if (temp.length <= 200) {
    return temp;
  } else {
    return temp.substring(0, 151) + '...';
  }
};
const getHelpTextTooltip = (row: string) => {
  return trim(row).replaceAll('\n', '<br />');
};

const columns: DataTableColumns<HelpTextItem> = [
  { title: '序号', key: 'id', align: 'center' },
  {
    title: '分组',
    key: 'group',
    align: 'center',
    render: row => (
      <n-tag type="success" size="small" bordered={false}>
        {row.group}
      </n-tag>
    ),
  },
  { title: '来源文件', key: 'from' },
  { title: '词条名', key: 'title' },
  {
    title: '内容',
    key: 'content',
    render: row => (
      <n-tooltip content-class="">
        {{
          trigger: () => getHelpText(row.content),
          default: () => getHelpTextTooltip(row.content),
        }}
      </n-tooltip>
    ),
  },
  { title: '分类', key: 'packageName' },
];
const handleCurrentPageChange = async (val: number) => {
  textItemQuery.value.pageNum = val;
  await refreshTextItems();
};
const refreshTextItems = async () => {
  const resp = await getHelpDocPage(textItemQuery.value);
  if (resp?.result) {
    textItems.value = resp.data;
    textItemQuery.value.total = resp.total;
  } else {
    message.error('无法获取帮助词条信息');
  }
};

const reloadLoading = ref<boolean>(false);
const reload = async () => {
  reloadLoading.value = true;
  const resp = await reloadHelpDoc();
  if (resp?.result) {
    message.success('重载帮助文件成功');
    needReload.value = false;
    await refreshFileTree();
    await refreshTextItems();
    await refreshConfig();
  } else {
    message.error(resp.err ?? '重载帮助文件失败');
  }
  reloadLoading.value = false;
};

const configDialogVisible = ref<boolean>(false);
const configNeedSave = ref<boolean>(false);
const helpAliases = ref<Map<string, string[]>>(new Map());

const refreshConfig = async () => {
  const config = await getHelpDocConfig();
  helpAliases.value = new Map(Object.entries(config.aliases));
  configNeedSave.value = false;
};

const configDialogClose = () => {
  configDialogVisible.value = false;
};

const addAlias = (groupKey: string, alias: string) => {
  // 别名不能重复
  for (const aliases of helpAliases.value.values()) {
    if (aliases.includes(alias)) {
      message.error('别名 ' + alias + ' 已被使用');
      return;
    }
  }
  if (helpAliases.value.has(groupKey)) {
    helpAliases.value?.get?.(groupKey)?.push(alias);
  } else {
    helpAliases.value.set(groupKey, [alias]);
  }
  configNeedSave.value = true;
};

const removeAlias = (groupKey: string, alias: string) => {
  if (helpAliases.value.has(groupKey)) {
    const lst = helpAliases.value.get(groupKey) ?? [];
    helpAliases.value.set(
      groupKey,
      lst.filter(v => v !== alias),
    );
    configNeedSave.value = true;
  }
};

const summitConfig = async () => {
  if (helpAliases.value) {
    console.log('aliases=', helpAliases.value);
    const res = await postHelpDocConfig({
      aliases: Object.fromEntries(helpAliases.value),
    });
    if (res.result) {
      message.success('保存设置成功');
      configDialogClose();
      nextTick(async () => {
        await refreshConfig();
      });
    } else {
      message.error('保存设置失败！' + res.err);
    }
  }
};

onBeforeMount(async () => {
  await refreshFileTree();
  await refreshTextItems();
  await refreshConfig();
});
</script>

<style lang="css">
@media screen and (max-width: 700px) {
  .bak-item {
    flex-direction: column;

    & > span {
      overflow: hidden;
      white-space: nowrap;
      text-overflow: ellipsis;
    }
  }
}

.file-tree-title {
  padding: 0 23px 6px 50px;
  margin-bottom: 10px;

  flex: 1;
  display: flex;
  justify-content: space-between;
}

.file-line {
  flex: auto;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.file-info {
  flex: auto;
  width: 0;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.item-list-container {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.item-list-pagination {
  margin-top: 10px;
}

.del-line {
  text-decoration: line-through;
}
</style>

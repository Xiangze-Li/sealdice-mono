<template>
  <n-affix :top="60" v-if="modified">
    <tip-box type="error">
      <n-flex>
        <n-text type="error" class="text-base" tag="strong">内容已修改，不要忘记保存！</n-text>
        <n-button class="button" type="primary" @click="save" :disabled="!modified">
          <template #icon>
            <n-icon><i-carbon-save /></n-icon>
          </template>
          点我保存
        </n-button>
      </n-flex>
    </tip-box>
  </n-affix>

  <tip-box type="info">
    <n-collapse>
      <n-collapse-item name="1">
        <template #header>
          <n-text tag="strong">查看帮助</n-text>
        </template>

        <n-text tag="p">
          <div>此处可以对骰子返回的文本进行修改。最终返回的文本将为多个条目中随机抽取的一个。</div>
          <div>
            随机文本：默认一种显示结果，如果需要多种反馈结果，使用＋添加条目，使用 - 删除条目
          </div>
          <!-- 权重选择：默认 1，权重—致则没有优先级。数字越小，优先级越高 -->
          <!-- <div>文件备份：已修改的指令统一存在于路径/路径 1/路径 2/文件名，如有需要替换文件即可</div> -->
          <div>
            遇到有此标记 (<n-icon><i-carbon-paint-brush /></n-icon>)
            的条目，说明和默认值不同，是一个自定义条目
          </div>
          <div style="margin-top: 1rem">
            文本下方的
            <n-tag size="small" :bordered="false">标签</n-tag>
            代表了被默认文本所使用的特殊变量，你可以使用 {变量名} 来插入他们，例如 {$t判定值}
          </div>
          <div>
            除此之外，这些变量可以在所有文本中使用：
            <n-flex size="small" wrap>
              <n-tag
                :key="i"
                size="small"
                :bordered="false"
                v-for="i in [
                  '$t玩家',
                  '$tQQ昵称',
                  '$t个人骰子面数',
                  '$tQQ',
                  '$t骰子帐号',
                  '$t骰子昵称',
                  '$t群号',
                  '$t群名',
                ]">
                {{ i }}
              </n-tag>
            </n-flex>
          </div>
          <div>
            <span>以及，所有的自定义文本都可以嵌套使用，例如：</span>
            <div>
              <b>这里是{核心：骰子名字}，我是一个示例</b>
            </div>
            <div>默认会被解析为：</div>
            <div>
              <b>这里是海豹 bot，我是一个示例</b>
            </div>
            <div>注意！千万不要递归嵌套，会发生很糟糕的事情。</div>
          </div>

          <div style="margin-top: 1rem">
            <div>此外，支持插入图片，将图片放在骰子的适当目录，再写这样一句话即可：</div>
            <div><n-text code>[图:data/images/sealdice.png]</n-text></div>
            <div>可以参考 核心：骰子进群 词条</div>
            <div>同样的，可以使用 CQ 码插入图片和其他内容，关于 CQ 码，请参阅 onebot 项目文档</div>
          </div>

          <div style="margin-top: 1rem">
            <b>
              COC 的“判定 - 常规”和“判定 - 简短”主要区别是，多重检定会默认使用简短版本 (.ra 3#射击)
            </b>
            <b>进行调整后，可以在左侧面板“指令测试”中进行测试！</b>
          </div>
        </n-text>
      </n-collapse-item>
    </n-collapse>
  </tip-box>

  <div style="margin-top: 1rem; display: flex; justify-content: space-between; align-items: center">
    <n-flex align="center" size="small">
      <n-text>搜索：</n-text>
      <span>
        <n-input size="small" v-model:value="currentFilterName" clearable>
          <template #prefix>
            <n-icon><i-carbon-search /> </n-icon>
          </template>
        </n-input>
      </span>
    </n-flex>
    <n-button type="info" secondary @click="dialogImportVisible = true">导入/导出</n-button>
  </div>

  <n-flex class="mb-8 mt-4" align="center" wrap>
    <n-radio-group v-model:value="filterMode" @update:value="handleFilterModeChange">
      <n-radio
        v-for="mode of filterModes"
        :key="mode.value"
        :value="mode.value"
        :label="mode.desc" />
    </n-radio-group>
    <n-flex v-if="filterMode === 'group'" align="center">
      <n-text>分组：</n-text>
      <span>
        <n-select
          v-model:value="currentFilterGroup"
          filterable
          tag
          :options="
            filterGroups.map(g => {
              return { label: g, value: g };
            })
          " />
      </span>
    </n-flex>
  </n-flex>

  <n-collapse class="text-collapse" :default-expanded-names="['__others__']">
    <custom-text-box
      :key="group"
      v-for="[group, values] in reactive(doSort(category))"
      :group="group">
      <template #values>
        <n-grid x-gap="24" :cols="2">
          <n-grid-item v-for="[k, v] in values" :key="k">
            <n-form ref="form" label-width="auto" label-position="top">
              <n-form-item class="w-full">
                <template #label>
                  <div>
                    <n-tag
                      type="default"
                      size="small"
                      style="margin-right: 0.5rem"
                      :bordered="false">
                      {{
                        store.curDice.customTextsHelpInfo[category][k.toString()].subType ||
                        (store.curDice.customTextsHelpInfo[category][k.toString()].notBuiltin
                          ? '旧版文本'
                          : '其它')
                      }}
                    </n-tag>

                    <span>
                      <span>{{ k.toString() }}</span>
                      <n-tooltip
                        v-if="store.curDice.customTextsHelpInfo[category][k.toString()].extraText">
                        <template #trigger>
                          <n-icon><i-carbon-help-filled /></n-icon>
                        </template>
                        {{ store.curDice.customTextsHelpInfo[category][k.toString()].extraText }}
                      </n-tooltip>
                    </span>

                    <template
                      v-if="store.curDice.customTextsHelpInfo[category][k.toString()].notBuiltin">
                      <n-tooltip placement="bottom-end">
                        <template #trigger>
                          <n-icon
                            style="float: right; margin-left: 1rem"
                            @click="askDeleteValue(category, k.toString())">
                            <i-carbon-row-delete />
                          </n-icon>
                        </template>
                        移除 - 这个文本在新版的默认配置中不被使用，<br />
                        但升级而来时仍可能被使用，请确认无用后删除
                      </n-tooltip>
                    </template>

                    <template
                      v-if="store.curDice.customTextsHelpInfo[category][k.toString()].modified">
                      <n-tooltip content="" placement="bottom-end">
                        <template #trigger>
                          <n-icon
                            style="float: right; margin-left: 1rem"
                            @click="askResetValue(category, k.toString())">
                            <i-carbon-paint-brush />
                          </n-icon>
                        </template>
                        重置为初始值
                      </n-tooltip>
                    </template>
                    <!-- <el-tooltip content="效果预览" placement="bottom-end">
                    <el-icon @click="askResetValue(category, k.toString())" style="float: right;"><video-play /></el-icon>
                  </el-tooltip> -->
                  </div>
                </template>

                <n-flex vertical class="w-full">
                  <n-text v-if="k.toString() === '戳一戳'" type="warning" class="mb-1 text-xs">
                    请确认你使用的 QQ
                    连接方式支持该功能，若不支持请于「基本设置」中关闭戳一戳来避免日志中出现相关报错。
                  </n-text>

                  <div
                    v-for="(k2, index) in v"
                    :key="index"
                    style="width: 100%; margin-bottom: 0.5rem">
                    <!-- 这里面是单条修改项 -->
                    <n-flex align="center">
                      <div>
                        <n-tooltip placement="bottom-start">
                          <template #trigger>
                            <n-icon>
                              <i-carbon-add-filled v-if="index === 0" @click="addItem(k)" />
                              <i-carbon-close-outline v-else @click="removeItem(v, index)" />
                            </n-icon>
                          </template>
                          {{
                            index === 0
                              ? '点击添加一个回复语，SealDice将会随机抽取一个回复'
                              : '点击删除你不想要的回复语'
                          }}
                        </n-tooltip>
                      </div>
                      <div class="relative flex-auto">
                        <n-input
                          class="w-full"
                          type="textarea"
                          :autosize="{ minRows: 3 }"
                          v-model:value="k2[0]"
                          @update:value="doChanged(category, k.toString())" />

                        <div class="absolute bottom-0 right-1" v-if="getPreview(k, k2[0])">
                          <n-popover placement="bottom-start">
                            <template #trigger>
                              <span
                                v-if="getPreviewCheckErr(k, k2[0])"
                                class="text-red-500"
                                style="margin-left: 0.1rem; margin-top: 0.1rem">
                                <n-icon><i-carbon-close-filled /></n-icon>
                              </span>
                              <n-flex v-else>
                                <span
                                  v-if="getPreview(k, k2[0]).version === 'v2'"
                                  class="text-blue-500"
                                  style="margin-left: 0.1rem; margin-top: 0.1rem">
                                  <n-icon><i-carbon-checkmark-filled /></n-icon>
                                </span>
                                <span
                                  v-if="getPreview(k, k2[0]).version === 'v1'"
                                  class="text-yellow-500"
                                  style="margin-left: 0.1rem; margin-top: 0.1rem">
                                  <n-icon><i-carbon-checkmark-filled /></n-icon>
                                </span>
                              </n-flex>
                            </template>

                            <component :is="getPreviewInfo(k, k2[0])" />
                          </n-popover>
                        </div>
                      </div>
                    </n-flex>
                  </div>
                  <n-flex size="small" wrap>
                    <n-tag
                      size="small"
                      type="info"
                      :bordered="false"
                      v-for="i in store.curDice.customTextsHelpInfo[category][k.toString()].vars"
                      :key="i">
                      {{ i }}
                    </n-tag>
                    <!-- {{ store.curDice.customTextsHelpInfo[category][k.toString()] }} -->
                  </n-flex>
                </n-flex>
              </n-form-item>
            </n-form>
          </n-grid-item>
        </n-grid>
      </template>
    </custom-text-box>
  </n-collapse>

  <n-modal
    v-model:show="dialogImportVisible"
    preset="card"
    title="导入导出"
    :mask-closable="false"
    :close-on-esc="false"
    :closable="false"
    class="the-dialog">
    <template #header-extra>
      <n-flex>
        <n-switch v-model:value="importOnlyCurrent">
          <template #checked>仅当前页面</template>
          <template #unchecked>全部文案</template>
        </n-switch>
        <n-checkbox v-model:checked="importImpact">紧凑</n-checkbox>
      </n-flex>
    </template>
    <n-flex vertical>
      <!-- <template > -->
      <n-text tag="strong">以下为导出内容，可以复制给别人</n-text>
      <n-input
        placeholder="填入数据"
        type="textarea"
        :autosize="{ minRows: 4 }"
        class="import-edit"
        id="import-edit"
        v-model:value="configForImport" />
      <!-- </template> -->
    </n-flex>

    <template #footer>
      <n-flex>
        <n-button @click="dialogImportVisible = false">返回</n-button>
        <n-button type="warning" @click="configForImport = ''">清空</n-button>
        <n-button type="info" data-clipboard-target="#import-edit" @click="copied" id="btnCopy1">
          复制
        </n-button>
        <n-button type="primary" @click="doImport" :disabled="configForImport === ''">
          导入并保存
        </n-button>
      </n-flex>
    </template>
  </n-modal>
</template>

<script setup lang="tsx">
import { useStore } from '~/store';
import {
  cloneDeep,
  filter,
  groupBy,
  map,
  sortBy,
  startsWith,
  trim,
  uniq,
  mapValues,
} from 'es-toolkit/compat';
import entries from 'es-toolkit/compat';
import ClipboardJS from 'clipboard';
import { useDialog, useMessage } from 'naive-ui';

const store = useStore();
const message = useMessage();
const dialog = useDialog();
const props = defineProps<{ category: string }>();

const configForImport = ref('');
const importOnlyCurrent = ref(true);
const importImpact = ref(true);
const dialogImportVisible = ref(false);

const doSort = (category: string) => {
  let items = Object.entries(store.curDice?.customTexts?.[category] ?? {});
  const helpInfo = store.curDice.customTextsHelpInfo[category];

  if (currentFilterName.value != '') {
    items = items.filter(
      item =>
        item[0].includes(currentFilterName.value) ||
        helpInfo[item[0]].subType.includes(currentFilterName.value),
    );
  }

  switch (filterMode.value) {
    case 'all':
      break;
    case 'unmodified':
      items = items.filter(item => !helpInfo[item[0]].modified);
      break;
    case 'modified':
      items = items.filter(item => helpInfo[item[0]].modified);
      break;
    case 'deprecated':
      items = items.filter(item => helpInfo[item[0]].notBuiltin);
      break;
    case 'group':
      filterGroups.value = sortBy(
        uniq(
          Object.values(helpInfo)
            .map(info => trim(info.subType))
            .filter(subType => subType !== ''),
        ),
      );
      items = items.filter(item =>
        startsWith(trim(helpInfo[item[0]].subType), currentFilterGroup.value),
      );
      break;
  }

  // 子类别元素超过 4 个的，展示上需要打包
  const boxedGroups = map(
    filter(
      entries(
        groupBy(
          map(items, item => {
            const subType = helpInfo[item[0]].subType;
            return subType.split(' ')[0];
          }),
          group => group,
        ),
      ),
      group => group[1].length >= 4,
    ),
    group => group[0],
  );

  const outMap = entries(
    mapValues(
      groupBy(items, item => {
        const group = helpInfo[item[0]].subType.split(' ')[0];
        if (boxedGroups.includes(group)) {
          return group;
        } else {
          return '__others__';
        }
      }),
      items =>
        items.sort((a, b) => {
          const ia = helpInfo[a[0]];
          const ib = helpInfo[b[0]];

          if (ia.topOrder !== ib.topOrder) {
            return ib.topOrder - ia.topOrder;
          }

          if (ia.subType !== ib.subType) {
            return ib.subType.localeCompare(ia.subType);
          }

          return 0;
        }),
    ),
  ).sort((a, b) => {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const [aGroup, _] = a;
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const [bGroup, _2] = b;
    if (aGroup === '__others__') {
      return -1;
    } else if (bGroup === '__others__') {
      return 1;
    } else {
      return 0;
    }
  });

  return outMap;
};

const copied = () => {
  message.success('进行了复制！');
};

const importRefresh = () => {
  const indent = !importImpact.value ? 2 : 0;
  if (importOnlyCurrent.value) {
    configForImport.value = JSON.stringify(
      {
        title: '某人的自定义配置',
        items: {
          [props.category]: store.curDice.customTexts[props.category],
        },
      },
      null,
      indent,
    );
  } else {
    configForImport.value = JSON.stringify(store.curDice.customTexts, null, indent);
  }
};

const doImport = async () => {
  try {
    const data = JSON.parse(configForImport.value);
    if (!(data.title && data.items)) {
      message.error('格式不正确');
      return;
    }

    for (const [k, v] of Object.entries(data.items)) {
      if (store.curDice.customTexts[k]) {
        store.curDice.customTexts[k] = v as any;
      }
      await store.customTextSave(k);
    }
    await store.getCustomText();
    modified.value = false;
    message.success('已保存');
    dialogImportVisible.value = false;
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
  } catch (e) {
    message.error('格式不正确');
  }
};

watch(
  () => dialogImportVisible.value,
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  (newValue, oldValue) => {
    if (newValue) {
      importRefresh();
      nextTick(() => {
        new ClipboardJS(document.getElementById('btnCopy1') as any);
      });
    }
  },
);

watch(
  () => [importImpact.value, importOnlyCurrent.value],
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  newValue => {
    importRefresh();
  },
);

watch(
  () => props.category,
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  (newValue, oldValue) => {
    //直接监听
    modified.value = false;
  },
);

const addItem = (k: any) => {
  store.curDice.customTexts[props.category][k].push(['', 1 as any]);
  modified.value = true;
};

const doChanged = (category: string, keyName: string) => {
  modified.value = true;
  const helpInfo = store.curDice.customTextsHelpInfo[category][keyName];
  helpInfo.modified = true;
};

const removeItem = (v: any[], index: number) => {
  v.splice(index, 1);
  modified.value = true;
};

const save = async () => {
  modified.value = false;
  await store.customTextSave(props.category);
  await store.getCustomText();
  modified.value = false;
  message.success('已保存');
};

const getPreview = (k: string, text: string) => {
  const x = store.curDice.previewInfo as any;
  if (x) {
    return x[`${props.category}:${k}`][text] as any;
  }
};

const getPreviewCheckErr = (k: string, text: string) => {
  const info = getPreview(k, text);
  if (info) {
    if (info.version == 'v2') return Boolean(info.errV2);
    if (info.version == 'v1') return Boolean(info.errV1);
  }
  return false;
};

const getPreviewInfo = (k: string, text: string) => {
  const info = getPreview(k, text);
  if (info) {
    let version = info.version;

    if (version === 'v1') {
      version = 'v1 [建议修改]';
    }
    const exists = info.presetExists ? '是' : '否';

    return (
      <div>
        <n-descriptions
          label-placement="left"
          label-align="left"
          separator=" "
          column={1}
          content-class="whitespace-nowrap break-words">
          <n-descriptions-item>
            {{
              label: () => (
                <n-tag type="success" size="small" bordered={false}>
                  引擎版本
                </n-tag>
              ),
              default: () => version,
            }}
          </n-descriptions-item>
          <n-descriptions-item>
            {{
              label: () => (
                <n-tag type="info" size="small" bordered={false}>
                  V2 预览
                </n-tag>
              ),
              default: () => info.textV2 || info.errV2,
            }}
          </n-descriptions-item>
          <n-descriptions-item>
            {{
              label: () => (
                <n-tag type="warning" size="small" bordered={false}>
                  V1 预览
                </n-tag>
              ),
              default: () => info.textV1 || info.errV1,
            }}
          </n-descriptions-item>
          <n-descriptions-item>
            {{
              label: () => (
                <n-tag type="success" size="small" bordered={false}>
                  存在预设
                </n-tag>
              ),
              default: () => exists + ' [存在时预览较为可靠]',
            }}
          </n-descriptions-item>
        </n-descriptions>
      </div>
    );
  }
};

const deleteValue = async (category: string, keyName: string) => {
  // const helpInfo = store.curDice.customTextsHelpInfo[category][keyName];
  delete store.curDice.customTexts[category][keyName];
  modified.value = true;
};

const askDeleteValue = async (category: string, keyName: string) => {
  dialog.warning({
    title: '警告',
    content: '删除这条文本，确定吗？',
    positiveText: '确定',
    negativeText: '取消',
    closable: false,
    onPositiveClick: async () => {
      await deleteValue(category, keyName);
      message.success('成功');
    },
  });
};

const resetValue = async (category: string, keyName: string) => {
  const helpInfo = store.curDice.customTextsHelpInfo[category][keyName];
  store.curDice.customTexts[category][keyName] = cloneDeep(helpInfo.origin);
  helpInfo.modified = false;
  modified.value = true;
};

const askResetValue = async (category: string, keyName: string) => {
  dialog.warning({
    title: '警告',
    content: '重置这条文本回默认状态，确定吗？',
    positiveText: '确定',
    negativeText: '取消',
    closable: false,
    onPositiveClick: async () => {
      await resetValue(category, keyName);
      message.success('成功');
    },
  });
};

const modified = ref(false);

interface FilterMode {
  value: string;
  desc: string;
}

const filterModes: FilterMode[] = [
  { value: 'all', desc: '全部' },
  { value: 'unmodified', desc: '默认文案' },
  { value: 'modified', desc: '修改过' },
  { value: 'group', desc: '指定分组' },
  { value: 'deprecated', desc: '旧版文本' },
];
const filterMode = ref<string>('all');
const filterGroups = ref<string[]>([]);
const currentFilterGroup = ref<string>('');
const currentFilterName = ref<string>('');

const handleFilterModeChange = (newMode: any) => {
  if (newMode === 'group') {
    nextTick(() => {
      currentFilterGroup.value = filterGroups.value[0];
      currentFilterName.value = '';
    });
  } else {
    currentFilterGroup.value = '';
    currentFilterName.value = '';
  }
};

onBeforeMount(async () => {
  modified.value = false;
});

watch(props, () => {
  filterMode.value = 'all';
});
</script>

<style scoped>
@reference '../styles/index.css';

.import-edit > textarea {
  max-height: 65vh;
}

.text-collapse {
  width: 100%;
}

.preview-tip-label {
  @apply rounded-xs px-1 mr-1 my-1 inline-block w-14 text-center border border-gray-700;
}
</style>

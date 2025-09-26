<template>
  <header class="page-header mb-4">
    <n-switch v-model:value="replyEnable" @update:value="switchClick">
      <template #checked>启用</template>
      <template #unchecked>关闭</template>
      总开关
    </n-switch>
    <n-button v-if="store.curDice.config.customReplyConfigEnable" type="primary" @click="doSave">
      <template #icon>
        <n-icon><i-carbon-save /></n-icon>
      </template>
      保存
    </n-button>
  </header>

  <n-affix v-if="modified" :top="60">
    <tip-box type="error">
      <n-text type="error" class="text-base" tag="strong">内容已修改，不要忘记保存！</n-text>
    </tip-box>
  </n-affix>

  <div class="tip">
    <n-collapse v-model="activeTip">
      <n-collapse-item name="basic">
        <template #header>
          <n-text tag="strong">基础帮助</n-text>
        </template>
        <n-text tag="p">
          每项自定义回复由一个&lt;条件>触发，产生一系列&lt;结果><br />
          一旦一个&lt;条件>被满足，将立即停止匹配，并执行&lt;结果>
        </n-text>
      </n-collapse-item>
      <n-collapse-item name="advanced">
        <template #header>
          <n-text tag="strong">进阶内容</n-text>
        </template>
        <n-text tag="p">
          越靠前的项具有越高的优先级，可以拖动来调整优先顺序！<br />
          为了避免滥用和无限互答，自定义回复的响应间隔最低为 5s<br />
          匹配到的文本将被存入变量$t0，正则匹配的组将被存入$t1 $t2 ....<br />
          若存在组名，如 (?P&lt;A&gt;cc)，将额外存入$tA
        </n-text>
      </n-collapse-item>
    </n-collapse>
  </div>

  <n-divider />

  <main class="reply-main">
    <header
      v-if="store.curDice.config.customReplyConfigEnable"
      style="display: flex; flex-wrap: wrap; justify-content: space-between">
      <n-flex class="current-reply" vertical>
        <n-flex align="center" wrap>
          <strong>当前文件</strong>
          <n-select
            v-model:value="curFilename"
            :options="
              fileItems.map(item => {
                return { value: item?.filename, label: item?.filename };
              })
            "
            class="w-48">
          </n-select>
          <n-checkbox v-model:checked="cr.enable" style="margin-left: 1rem">
            {{ cr.enable ? '已启用' : '未启用' }}
          </n-checkbox>
        </n-flex>
        <n-flex class="mt-2" warp>
          <n-button type="error" size="small" tertiary @click="customReplyFileDelete">
            <template #icon>
              <i-carbon-row-delete />
            </template>
            删除
          </n-button>
          <n-button
            type="info"
            size="small"
            tertiary
            tag="a"
            style="text-decoration: none"
            :href="`${urlBase}/sd-api/configs/custom_reply/file_download?name=${encodeURIComponent(curFilename)}&token=${encodeURIComponent(store.token)}`">
            <template #icon>
              <n-icon><i-carbon-download /></n-icon>
            </template>
            下载
          </n-button>
        </n-flex>
        <n-text v-if="!cr.enable" class="mt-2" type="warning">
          注意：启用后该文件中的自定义回复才会生效
        </n-text>
      </n-flex>
      <div class="reply-operation mt-4 sm:mt-0">
        <div>
          <n-tooltip>
            <template #trigger>
              <n-button type="success" secondary @click="customReplyFileNew">
                <template #icon>
                  <n-icon><i-carbon-document-add /></n-icon>
                </template>
                新建
              </n-button>
            </template>
            新建一个自定义回复文件
          </n-tooltip>
        </div>
        <div>
          <n-tooltip>
            <template #trigger>
              <n-upload
                action=""
                multiple
                accept=".yaml"
                @before-upload="beforeUpload"
                :show-file-list="false"
                :file-list="uploadFileList">
                <n-button type="info" secondary>
                  <template #icon>
                    <n-icon><i-carbon-upload /></n-icon>
                  </template>
                  上传
                </n-button>
              </n-upload>
            </template>
            上传自定义回复的 .yaml 文件
          </n-tooltip>
        </div>
        <div>
          <n-tooltip>
            <template #trigger>
              <n-button type="default" secondary @click="dialogFormVisible = true">
                <template #icon>
                  <n-icon><i-carbon-document /></n-icon>
                </template>
                解析
              </n-button>
            </template>
            通过粘贴/编辑文本来导入自定义回复
          </n-tooltip>
        </div>
      </div>
    </header>

    <n-divider v-if="store.curDice.config.customReplyConfigEnable" />

    <template v-if="!store.curDice.config.customReplyConfigEnable">
      <n-text type="error" class="text-xl">请先启用总开关！</n-text>
    </template>
    <template v-else>
      <foldable-card
        ref="commonConditionsRef"
        type="div"
        :default-fold="true"
        :class="cr.enable ? '' : 'disabled'">
        <template #title>
          <n-flex size="large" align="center" wrap>
            <n-flex size="small" align="center">
              <strong>公共条件</strong>
              <n-button type="success" size="tiny" secondary @click="addOneCondition(conditions)">
                <template #icon>
                  <n-icon><i-carbon-add-large /></n-icon>
                </template>
                添加一项
              </n-button>
            </n-flex>
            <n-text type="info" class="text-xs">
              该文件下所有的回复的执行，都需要先满足以下公共条件（需同时满足，即 and）。
            </n-text>
          </n-flex>
        </template>

        <template v-if="conditions && conditions.length > 0">
          <custom-reply-conditions v-model="conditions" />
        </template>
        <template v-else>
          <n-text type="default">当前无公共条件</n-text>
        </template>

        <template #unfolded-extra>
          <template v-if="conditions && conditions.length > 0">
            <n-text type="info">公共条件数量：{{ conditions.length }}</n-text>
          </template>
          <template v-else>
            <n-text type="default">当前无公共条件</n-text>
          </template>
        </template>
      </foldable-card>

      <n-divider />

      <nested-draggable :tasks="list" :class="cr.enable ? '' : 'disabled'" />
      <n-flex align="center" justify="space-between">
        <n-button type="info" secondary @click="addOne(list)">
          <template #icon>
            <n-icon><i-carbon-add-large /> </n-icon>
          </template>
          添加一项
        </n-button>
        <n-button type="primary" @click="doSave">
          <template #icon>
            <n-icon><i-carbon-save /> </n-icon>
          </template>
          保存
        </n-button>
      </n-flex>
    </template>
  </main>

  <n-modal
    v-model:show="dialogFormVisible"
    preset="card"
    title="导入配置"
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    :show-close="false"
    class="the-dialog">
    <!-- <template > -->
    <n-input
      v-model:value="configForImport"
      placeholder="支持格式: 关键字/回复语"
      class="reply-text"
      type="textarea"
      :autosize="{ minRows: 4, maxRows: 10 }" />
    <template #footer>
      <n-flex>
        <n-button @click="dialogFormVisible = false">取消</n-button>
        <n-button type="primary" :disabled="configForImport === ''" @click="doImport">
          下一步
        </n-button>
      </n-flex>
    </template>
  </n-modal>

  <n-modal
    v-model:show="dialogLicenseVisible"
    preset="card"
    title="许可协议"
    :mask-closable="false"
    :close-on-esc="false"
    :closable="false"
    class="the-dialog">
    <pre style="white-space: pre-wrap">
尊敬的用户，欢迎您选择由木落等研发的海豹骰点核心（SealDice），在您使用自定义功能前，请务必仔细阅读使用须知，当您使用我们提供的服务时，即代表您已同意使用须知的内容。

  您需了解，海豹核心官方版只支持 TRPG 功能，娱乐功能定制化请自便，和海豹无关。
  您清楚并明白您对通过骰子提供的全部内容负责，包括自定义回复、非自带的插件、牌堆。海豹骰不对非自身提供以外的内容合法性负责。您不得在使用海豹骰服务时，导入包括但不限于以下情形的内容：
  (1) 反对中华人民共和国宪法所确定的基本原则的；
  (2) 危害国家安全，泄露国家秘密，颠覆国家政权，破坏国家统一的;
  (3) 损害国家荣誉和利益的;
  (4) 煽动民族仇恨、民族歧视、破坏民族团结的；
  (5) 破坏国家宗教政策，宣扬邪教和封建迷信的;
  (6) 散布谣言，扰乱社会秩序，破坏社会稳定的；
  (7) 散布淫秽、色情、赌博、暴力、凶杀、恐怖或者教唆犯罪的；
  (8) 侮辱或者诽谤他人，侵害他人合法权益的；
  (9) 宣扬、教唆使用外挂、私服、病毒、恶意代码、木马及其相关内容的；
  (10) 侵犯他人知识产权或涉及第三方商业秘密及其他专有权利的；
  (11) 散布任何贬损、诋毁、恶意攻击海豹骰及开发人员、海洋馆工作人员、mod 编写人员、关联合作者的；
  (12) 含有中华人民共和国法律、行政法规、政策、上级主管部门下发通知中所禁止的其他内容的。

  一旦查实您有以上禁止行为，请立即停用海豹骰。同时我们也会主动对你进行举报。</pre
    >
    <!-- 一旦查实您有以上禁止行为，我们有权进行核查、修改和/或删除您导入的内容，而不需事先通知。 -->

    <template #footer>
      <n-flex>
        <n-button type="primary" @click="dialogLicenseVisible = false">我同意</n-button>
        <n-button type="error" @click="licenseRefuse">我拒绝</n-button>
      </n-flex>
    </template>
  </n-modal>
</template>

<script lang="ts" setup>
import { urlBase } from '~/backend';
import { useStore } from '~/store';
import {
  getCustomReply,
  getCustomReplyFileList,
  postCustomReplyDel,
  postCustomReplyNew,
  saveCustomReply,
  uploadCustomReply,
} from '~/api/v1/configs';
import type { DiceConfig } from '~/api/v1/dice';
import { useDialog, useMessage } from 'naive-ui';

const store = useStore();
const message = useMessage();
const dialog = useDialog();

const dialogFormVisible = ref(false);
const dialogLicenseVisible = ref(false);
const configForImport = ref('');

const replyEnable = computed({
  get: () => store.curDice.config.customReplyConfigEnable,
  set: value => {
    store.diceConfigSet({ customReplyConfigEnable: value } as DiceConfig);
    if (!store.curDice.config.customReplyConfigEnable) {
      dialogLicenseVisible.value = true;
    }
  },
});

watch(replyEnable, async (newStatus, oldStatus) => {
  if (newStatus != oldStatus) {
  }
});

const activeTip = ref('basic');

const curFilename = ref('reply.yaml');

const conditions = ref<any>([]);

const commonConditionsRef = ref<any>(null);

const list = ref<any>([
  // {"enable":true,"condition":{"condType":"match","matchType":"match_exact","value":"asd"},"results":[{"resultType":"replyToSender","delay":0.3,"message":"text"}]},
  // {"enable":true,"condition":{"condType":"match","matchType":"match_fuzzy","value":"ccc"},"results":[{"resultType":"replyToSender","delay":0.3,"message":"text"}]},
]);

const fileItems = ref<any>([
  // {"enable":true,"condition":{"condType":"match","matchType":"match_exact","value":"asd"},"results":[{"resultType":"replyToSender","delay":0.3,"message":"text"}]},
  // {"enable":true,"condition":{"condType":"match","matchType":"match_fuzzy","value":"ccc"},"results":[{"resultType":"replyToSender","delay":0.3,"message":"text"}]},
]);

const uploadFileList = ref<any[]>([]);

const cr = ref<any>({ enable: true });

const switchClick = () => {
  if (!store.curDice.config.customReplyConfigEnable) {
    dialogLicenseVisible.value = true;
  }

  store.diceConfigSet({
    customReplyConfigEnable: !store.curDice.config.customReplyConfigEnable,
  } as DiceConfig);
};

const licenseRefuse = () => {
  dialogLicenseVisible.value = false;
  store.curDice.config.customReplyConfigEnable = false;
  store.diceConfigSet({ customReplyConfigEnable: false } as DiceConfig);
};

const modified = ref(false);

watch(
  () => cr.value,
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  (newValue, oldValue) => {
    //直接监听
    modified.value = true;
  },
  { deep: true },
);

watch(
  () => curFilename.value,
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  (newValue, oldValue) => {
    //直接监听
    nextTick(() => {
      refreshCurrent();
    });
  },
);

const customReplyFileNew = () => {
  ElMessageBox.prompt('创建一个新的回复文件', '', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'info',
    inputPlaceholder: 'reply2.yaml',
    inputValue: `reply${Math.ceil(Math.random() * 10000)}.yaml`,
  }).then(async data => {
    if (!data.value) {
      data.value = `reply${Math.ceil(Math.random() * 10000)}.yaml`;
    }
    const ret = await postCustomReplyNew(data.value);
    const ret2 = await getCustomReplyFileList();
    fileItems.value = ret2.items;
    curFilename.value = ret2.items[0].filename;

    if (ret.success) {
      message.success('创建成功');
    } else {
      message.error('创建失败');
    }
  });
};

const customReplyFileDelete = () => {
  dialog.warning({
    title: '删除文件',
    content: '是否删除此文件？',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      const ret = await postCustomReplyDel(curFilename.value);
      if (ret.success) {
        const ret = await getCustomReplyFileList();
        fileItems.value = ret.items;
        curFilename.value = ret.items[0].filename;
        message.success('删除成功');
        await refreshCurrent();
      } else {
        message.error('删除失败');
      }
    },
  });
};

const beforeUpload = async (file: any) => {
  // UploadRawFile
  try {
    await uploadCustomReply(file);
    message.success('上传完成');
    const ret = await getCustomReplyFileList();
    fileItems.value = ret.items;
    curFilename.value = ret.items[0].filename;
    await refreshCurrent();
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
  } catch (e) {
    message.error('上传失败，可能有同名文件！');
  }
};

const addOneCondition = (lst: any) => {
  lst.push({
    condType: 'textMatch',
    matchType: 'matchExact',
    value: '要匹配的文本',
  });
  commonConditionsRef.value.open();
};

const addOne = (lst: any) => {
  lst.push({
    enable: true,
    conditions: [{ condType: 'textMatch', matchType: 'matchExact', value: '要匹配的文本' }],
    results: [{ resultType: 'replyToSender', delay: 0, message: [['说点什么', 1]] }],
  });
};

// const deleteAnyItem = (lst: any[], index: number) => {
//   lst.splice(index, 1);
// };

const doSave = async () => {
  try {
    for (const i of cr.value.items) {
      for (const j of i.conditions) {
        if (j.condType === 'textLenLimit') {
          j.value = parseInt(j.value) || 0;
        }
      }

      for (const j of i.results) {
        if (j.delay) {
          j.delay = parseFloat(j.delay);
          if (j.delay < 0) {
            j.delay = 0;
          }
        }
        if (!j.delay) j.delay = 0;
      }
    }
    cr.value.filename = curFilename.value;
    cr.value.conditions = conditions.value;
    for (const cond of cr.value.conditions) {
      if (cond.condType === 'textLenLimit') {
        cond.value = parseInt(cond.value) || 0;
      }
    }
    await saveCustomReply(cr.value);
    message.success('已保存');
    modified.value = false;
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
  } catch (e) {
    message.error('保存失败！！');
  }
};

function parseString(str: string): [string[], string[], string] {
  const leftArr: string[] = [];
  const rightArr: string[] = [];
  let restIndex = 0;

  let currentStr = '';
  let isLeft = true;
  let isEscaped = false;

  for (let i = 0; i < str.length; i++) {
    const char = str[i];
    restIndex = i;
    if (isEscaped) {
      if (char !== '\r' && char !== '\n' && char !== '/') {
        currentStr += '\\'; // 如果是被转义的字符，将反斜杠保留
      }
      if (char === 'n' || char === 'r') {
        currentStr = currentStr.slice(0, -1) + '\n';
      } else {
        currentStr += char;
      }
      isEscaped = false;
      continue;
    }
    if (char == '\n') {
      break;
    }
    if (char === '\\') {
      isEscaped = true;
      continue;
    }
    if (char === '|') {
      if (isLeft) {
        leftArr.push(currentStr);
      } else {
        rightArr.push(currentStr);
      }
      currentStr = '';
    } else if (char === '/') {
      if (i < str.length - 1 && str[i + 1] === '|') {
        currentStr += char;
      } else {
        if (isLeft) {
          leftArr.push(currentStr);
        } else {
          rightArr.push(currentStr);
        }
        currentStr = '';
        isLeft = false;
      }
    } else {
      currentStr += char;
    }
  }
  // 处理最后一个字符串
  if (isLeft) {
    leftArr.push(currentStr);
  } else {
    rightArr.push(currentStr);
  }

  return [leftArr, rightArr, str.slice(restIndex + 1)];
}

const doImport = () => {
  // let count = 0;
  let text = configForImport.value;

  while (true) {
    const [a, b, rest] = parseString(text);
    if (a.length && b.length) {
      const replies = b.map(v => [v, 1]);
      list.value.push({
        enable: true,
        conditions: [
          {
            condType: 'textMatch',
            matchType: 'matchMulti',
            value: a.join('|'),
          },
        ],
        results: [{ resultType: 'replyToSender', delay: 0, message: replies }],
      });
      // count += 1;
    }
    text = rest;
    if (!rest) break;
  }

  message.success('导入成功！');
  configForImport.value = '';
};

onBeforeMount(async () => {
  const ret = await getCustomReplyFileList();
  fileItems.value = ret.items;
  curFilename.value = ret.items[0].filename;
  await store.diceConfigGet();
  await refreshCurrent();
});

const refreshCurrent = async () => {
  console.log('load', curFilename.value);
  const ret = await getCustomReply(curFilename.value);
  cr.value = ret;
  conditions.value = ret.conditions;
  list.value = ret.items;
  await nextTick(() => {
    modified.value = false;
  });
};

onBeforeUnmount(() => {
  // clearInterval(timerId)
});
</script>

<style scoped>
.reply-text > textarea {
  max-height: 65vh;
}
</style>

<style scoped lang="css">
@media screen and (max-width: 700px) {
  .bak-item {
    flex-direction: column;

    & > span {
      overflow: hidden;
      white-space: nowrap;
      text-overflow: ellipsis;
    }
  }

  .current-reply {
    width: 100%;
  }

  .reply-operation {
    width: 100%;
    justify-content: space-around;
  }

  .reply-operation div:not(:last-child) {
    margin-right: 1rem;
  }
}

@media screen and (min-width: 700px) {
  .reply-operation {
    flex-direction: column;
    margin-right: 0.5rem;
  }

  .reply-operation div:not(:last-child) {
    margin-bottom: 0.5rem;
  }
}

.reply-operation {
  display: flex;
  justify-content: center;
  align-items: self-start;
}

.disabled {
  filter: grayscale(1);
  /* pointer-events: none; */
  cursor: not-allowed;
  user-select: none;
}

.img-box {
  height: 250px;
  margin-right: 3rem;
  float: left;

  img {
    height: 200px;
  }
}

.about {
  background-color: #fff;
  padding: 2rem;
  line-height: 2rem;
  text-align: left;
  box-shadow:
    0 2px 4px rgba(0, 0, 0, 0.12),
    0 0 6px rgba(0, 0, 0, 0.04);
}

.subtitle {
  margin-bottom: 1rem;
  font-weight: bold;
}
</style>

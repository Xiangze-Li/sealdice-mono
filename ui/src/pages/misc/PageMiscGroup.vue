<template>
  <h2>群管理</h2>
  <n-flex vertical>
    <n-flex align="center">
      <n-text class="text-base">平台：</n-text>
      <n-checkbox-group v-model:value="checkList">
        <n-checkbox value="QQ-Group:">QQ 群</n-checkbox>
        <n-checkbox value="QQ-CH-Group:">QQ 频道</n-checkbox>
        <n-checkbox value="DISCORD-CH-Group:">Discord 频道</n-checkbox>
        <n-checkbox value="DODO-Group:">Dodo 频道</n-checkbox>
        <n-checkbox value="KOOK-CH-Group:">KOOK 频道</n-checkbox>
        <n-checkbox value="DINGTALK-Group:">钉钉群</n-checkbox>
        <n-checkbox value="SLACK-CH-Group:">Slack 频道</n-checkbox>
        <n-checkbox value="TG-Group:">TG 群</n-checkbox>
        <n-checkbox value="SEALCHAT-Group:">Sealchat 频道</n-checkbox>
      </n-checkbox-group>
    </n-flex>
    <n-flex align="center">
      <n-text class="text-base">其他：</n-text>
      <n-checkbox v-model:checked="orderByTimeDesc">按最后使用时间降序</n-checkbox>
      <n-checkbox v-model:checked="filter30daysUnused">30 天未使用</n-checkbox>
    </n-flex>
    <n-flex align="center">
      <n-text class="text-base">搜索：</n-text>
      <n-input
        v-model:value="searchBy"
        style="max-width: 15rem"
        placeholder="请输入帐号或群名的一部分" />
    </n-flex>
  </n-flex>

  <div style="margin-top: 2rem">
    <div v-bind="containerProps" style="height: calc(100vh - 22.5rem)">
      <div v-bind="wrapperProps">
        <div v-for="item in list" :key="item.index">
          <foldable-card style="margin-top: 1rem">
            <template #title>
              <n-flex class="item-header" size="large" align="center">
                <n-switch
                  v-model:value="item.data.active"
                  @update-value="item.data.changed = true" />
                <n-flex size="small" wrap>
                  <n-text
                    class="hidden max-w-96 text-ellipsis whitespace-nowrap text-base"
                    tag="strong">
                    {{ item.data.groupId }}
                  </n-text>
                  <n-text>「{{ item.data.groupName || '未获取到' }}」</n-text>
                </n-flex>
              </n-flex>
            </template>

            <template #title-extra>
              <n-flex>
                <n-button
                  v-if="item.data.changed"
                  type="success"
                  size="tiny"
                  secondary
                  @click="saveOne(item.data, item.index)">
                  <template #icon>
                    <n-icon> <i-carbon-save /></n-icon>
                  </template>
                  保存
                </n-button>
                <template v-if="item.data.groupId.startsWith('QQ-Group:')">
                  <n-tooltip v-for="(_, j) in item.data.diceIdExistsMap" :key="j">
                    <template #trigger>
                      <n-button
                        type="error"
                        size="tiny"
                        secondary
                        @click="quitGroup(item.data, item.index, j.toString())">
                        <template #icon>
                          <n-icon><i-carbon-close-large /></n-icon>
                        </template>
                        退出 {{ j.toString().slice(-4) }}
                      </n-button>
                    </template>
                    {{ j.toString() }}<br />有二次确认
                  </n-tooltip>
                </template>
              </n-flex>
            </template>

            <n-descriptions label-placement="left" :column="3" separator=" ">
              <n-descriptions-item label="上次使用">
                {{
                  item.data.recentDiceSendTime
                    ? dayjs.unix(item.data.recentDiceSendTime).fromNow()
                    : '从未'
                }}
              </n-descriptions-item>
              <n-descriptions-item label="入群时间">
                {{ item.data.enteredTime ? dayjs.unix(item.data.enteredTime).fromNow() : '未知' }}
              </n-descriptions-item>
              <n-descriptions-item label="邀请人">
                {{ item.data.inviteUserId || '未知' }}
              </n-descriptions-item>
              <n-descriptions-item label="Log状态">
                {{ item.data.logOn ? '开启' : '关闭' }}
              </n-descriptions-item>
              <n-descriptions-item label="迎新">
                {{ item.data.showGroupWelcome ? '开启' : '关闭' }}
              </n-descriptions-item>
              <n-descriptions-item />
              <n-descriptions-item :span="3" label="启用扩展">
                <n-flex size="small" v-if="item.data.tmpExtList">
                  <n-tag
                    :key="group"
                    v-for="group of item.data.tmpExtList"
                    size="small"
                    type="info"
                    :bordered="false">
                    {{ group }}
                  </n-tag>
                </n-flex>
                <n-text v-else>'未知'</n-text>
              </n-descriptions-item>
            </n-descriptions>

            <template #unfolded-extra>
              <n-descriptions label-placement="left" separator=" ">
                <n-descriptions-item :span="2" label="上次使用">
                  {{
                    item.data.recentDiceSendTime
                      ? dayjs.unix(item.data.recentDiceSendTime).fromNow()
                      : '从未'
                  }}
                </n-descriptions-item>
                <n-descriptions-item label="邀请人">
                  {{ item.data.inviteUserId || '未知' }}
                </n-descriptions-item>
              </n-descriptions>
            </template>
          </foldable-card>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import * as dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import { now, sortBy } from 'es-toolkit/compat';
import { getGroupList, postQuitGroup, setGroup } from '~/api/v1/group';
import { useDialog, useMessage } from 'naive-ui';

const message = useMessage();
const dialog = useDialog();

dayjs.extend(relativeTime);

const checkList = ref<string[]>(['QQ-Group:']);

const groupList = ref<any>({});

const orderByTimeDesc = ref(true);
const filter30daysUnused = ref(false);
const searchBy = ref('');

const groupItems = computed<any[]>(() => {
  if (groupList.value.items) {
    // const groupListItems = cloneDeep(groupList.value.items)
    let items = [];
    for (const i of groupList.value.items) {
      let ok = false;
      for (const f of checkList.value) {
        if (i.groupId.startsWith(f)) {
          ok = true;
        }
      }

      if (ok && searchBy.value !== '') {
        let a = false;
        let b = false;
        if (i.groupId.indexOf(searchBy.value) !== -1) {
          a = true;
        }
        if (i.groupName.indexOf(searchBy.value) !== -1) {
          b = true;
        }
        ok = a || b;
      }

      if (ok) {
        const t = Math.max(i.enteredTime || 0, i.recentCommandTime || 0, i.recentDiceSendTime || 0);
        if (filter30daysUnused.value) {
          if (now() / 1000 - t < 30 * 24 * 60 * 60) {
            ok = false;
          }
        }
      }

      if (ok) items.push(i);
    }

    items = sortBy(items, ['recentCommandTime']);
    if (orderByTimeDesc.value) {
      items = items.reverse();
    }
    return items;
  }
  return [];
});

const refreshList = async () => {
  const data = await getGroupList();
  groupList.value = data;
};

const { list, containerProps, wrapperProps } = useVirtualList(groupItems, {
  itemHeight: 230,
});

const quitTextSave = ref(false);

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const saveOne = async (i: any, index: number) => {
  // await store.backupConfigSave(cfg.value)
  // console.log(222, i, index)
  await setGroup(i);
  i.changed = false;
  message.success('已保存');
};

const quitGroup = async (i: any, index: number, diceId: string) => {
  const quitGroupText =
    localStorage.getItem('quitGroupText') || '因长期不使用等原因，骰主后台操作退出';
  ElMessageBox.prompt(
    '会进行退出留言“因长期不使用等原因，骰主后台操作退出”，输入英文大写 NO 则静默退出，写别的则为附加留言',
    '退出此群？',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
      inputValue: quitGroupText,
      message: h('div', null, [
        h(
          'p',
          null,
          '会进行退出留言“因长期不使用等原因，骰主后台操作退出”，输入英文大写 NO 则静默退出，写别的则为附加留言',
        ),
        h(
          'label',
          {
            onInput: (e: any) => {
              quitTextSave.value = e.target.checked;
            },
          },
          [
            h('input', {
              value: quitTextSave.value,
              type: 'checkbox',
            }),
            h('span', null, '设为默认'),
          ],
        ),
      ]),
    },
  ).then(async data => {
    await postQuitGroup({
      groupId: i.groupId,
      diceId,
      silence: data.value === 'NO',
      extraText: data.value,
    });
    if (quitTextSave.value) {
      localStorage.setItem('quitGroupText', data.value);
    }

    await refreshList();
    message.success('退出完成');
    message.success('成功！');
  });
};

onBeforeMount(async () => {
  await refreshList();
});
</script>

<style lang="css">
span.left {
  display: inline-block;
  min-width: 5rem;
}

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

.item-header {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
  justify-content: space-between;
}
</style>

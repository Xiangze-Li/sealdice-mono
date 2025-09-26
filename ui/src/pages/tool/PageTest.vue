<template>
  <div class="flex h-full flex-col">
    <div class="flex justify-end">
      <n-flex align="center">
        <n-text>测试模式：</n-text>
        <n-radio-group v-model:value="mode" size="small">
          <n-radio-button label="私聊" value="private" />
          <n-radio-button label="群聊" value="group" />
        </n-radio-group>
      </n-flex>
    </div>

    <n-scrollbar ref="chat" trigger="none" class="flex-1 p-4">
      <div
        :key="index"
        v-for="(i, index) in store.talkLogs"
        v-show="i.mode === mode"
        class="talk-item"
        :class="!i.isSeal ? 'mine' : ''">
        <div class="left">
          <n-avatar :round="i.isSeal" :size="60" :src="i.isSeal ? imgSeal : imgMe" />
        </div>
        <div class="right">
          <n-text>{{ i.isSeal ? '海豹核心' : i.name }}</n-text>
          <div class="content">
            <n-text>{{ i.content }}</n-text>
          </div>
        </div>
      </div>
    </n-scrollbar>

    <div class="flex items-center">
      <n-auto-complete
        ref="autocomplete"
        v-model:value="input"
        :options="querySearch"
        placeholder="来试一试，回车键发送"
        @on-select="inputChanged"
        @keyup.enter="doSend" />
      <n-button class="mx-2 min-w-12" type="primary" @click="doSend">发送</n-button>
      <n-popover placement="top">
        <template #trigger>
          <n-button secondary vertical circle>
            <i-carbon-add-large />
          </n-button>
        </template>
        <n-flex vertical justify="center" class="flex w-full flex-col justify-center">
          <n-button quaternary :disabled="deckReloading" @click="reloadDeck">重载牌堆</n-button>
          <n-button quaternary :disabled="jsReloading" @click="reloadJs">重载 JS</n-button>
          <n-button quaternary :disabled="helpdocReloading" @click="reloadHelpdoc">
            重载帮助文件
          </n-button>
        </n-flex>
      </n-popover>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { useStore } from '~/store';
import imgSeal from '~/assets/seal.png';
import imgMe from '~/assets/me.jpg';
import { getRecentMessage, postExec } from '~/api/v1/dice';
import { reloadDeck as postReloadDeck } from '~/api/v1/deck';
import { reloadHelpDoc } from '~/api/v1/helpdoc';
import { reloadJS } from '~/api/v1/js';
import { NAutoComplete, NScrollbar, useThemeVars } from 'naive-ui';
const store = useStore();
const themeVars = useThemeVars();

const mode = ref<'private' | 'group'>('private');

let timerMsg: number;
onBeforeMount(async () => {
  restaurants.value = loadAll();
  timerMsg = setInterval(async () => {
    try {
      const msg = await getRecentMessage();
      console.log('msg:', msg);
      for (const i of msg) {
        store.talkLogs.push({
          content: i.message,
          isSeal: true,
          mode: i.messageType,
        });
      }
      if (msg.length) {
        // 拉下滚动条
        nextTick(() => {
          const c = chat.value;
          if (c) {
            c.scrollTo({ position: 'bottom', silent: false } as any);
          }
        });
      }
      // eslint-disable-next-line @typescript-eslint/no-unused-vars
    } catch (e: any) {}
  }, 1000) as any;
});

onBeforeUnmount(() => {
  clearInterval(timerMsg);
});

const restaurants = ref<RestaurantItem[]>([]);

interface RestaurantItem {
  label: string;
  value: string;
}

const input = ref('');

const chat = ref<InstanceType<typeof NScrollbar>>();

let lastTime = 0;

const inputChanged = () => {
  lastTime = Date.now();
};

const doSend = async () => {
  if (input.value === '') {
    return;
  }
  // 我的机器上至少要 50ms，其实应该有更好的办法
  if (Date.now() - lastTime > 300) {
    const text = input.value;
    store.talkLogs.push({
      name: '',
      content: text,
      isSeal: false,
      mode: mode.value,
    });
    try {
      await postExec(text, mode.value);
      // eslint-disable-next-line @typescript-eslint/no-unused-vars
    } catch (e) {
      store.talkLogs.push({
        name: '',
        content: '消息过于频繁',
        isSeal: true,
        mode: mode.value,
      });
    }

    nextTick(() => {
      const c = chat.value;
      if (c) {
        c.scrollTo({ position: 'bottom', silent: false } as any);
      }
      input.value = '';
    });
  }
};

const querySearch = computed(() => {
  const results = input.value
    ? restaurants.value.filter(elem => elem.value.toLowerCase().indexOf(input.value) === 0)
    : [];
  return results;
});

const loadAll = () => {
  const raw =
    '死亡豁免 spellslots character dlongrest 法术位 longrest botlist 查询 setcoc 咕咕 master 长休 角色 dcast reply dbuff gugu roll buff send name char drcv jrrp help find text cast draw init deck drav dndx rch dst drc rah log dnd rhd coc rhx ext dss rcv set rav bot li st st en ti ri sc ra rc ds rh rd pc nn ch rx ss r';
  const ret = [];
  for (const i of raw.split(' ')) {
    ret.push({ label: '.' + i, value: '.' + i });
  }
  ret.reverse();
  return ret;
};

const deckReloading = ref<boolean>(false);
const reloadDeck = async () => {
  deckReloading.value = true;
  const ret = await postReloadDeck();
  if (ret.testMode) {
    ElMessage.success('展示模式无法重载牌堆');
  } else {
    ElMessage.success('已重载牌堆');
  }
  deckReloading.value = false;
};

const jsReloading = ref<boolean>(false);
const reloadJs = async () => {
  jsReloading.value = true;
  const ret = await reloadJS();
  if (ret && ret?.testMode) {
    ElMessage.success('展示模式无法重载 JS');
  } else {
    ElMessage.success('已重载 JS');
  }
  jsReloading.value = false;
};

const helpdocReloading = ref<boolean>(false);
const reloadHelpdoc = async () => {
  helpdocReloading.value = true;
  const ret = await reloadHelpDoc();
  if (ret && ret?.result) {
    ElMessage.success('已重载帮助文档');
  } else {
    ElMessage.success(ret.err || '无法重载帮助文档');
  }
  helpdocReloading.value = false;
};
</script>

<style scoped lang="css">
.talk-item {
  display: flex;
  margin-bottom: 2rem;

  &.mine {
    direction: rtl;
    & > .right > .content {
      background-color: v-bind('themeVars.infoColorSuppl');
      direction: ltr;
    }
  }

  & > .right {
    padding-left: 1rem;
    padding-right: 1rem;
    & > .name {
      font-size: smaller;
      line-height: 2rem;
      min-height: 2rem;
    }
    & > .content {
      background-color: v-bind('themeVars.cardColor');
      padding: 0.7rem;
      border-radius: 9px;
      white-space: pre-wrap;
      overflow-wrap: anywhere;
    }
  }
}
</style>

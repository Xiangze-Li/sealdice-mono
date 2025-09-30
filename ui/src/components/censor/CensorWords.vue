<template>
  <n-flex justify="space-between">
    <h4>敏感词列表</h4>
    <n-flex align="center">
      <n-text v-if="filterCount > 0" type="info" class="text-xs">
        已过滤 {{ filterCount }} 条
      </n-text>
      <span>
        <n-input v-model:value="filter" size="small" placeholder="" clearable>
          <template #prefix>
            <n-icon><i-carbon-search /></n-icon>
          </template>
        </n-input>
      </span>
    </n-flex>
  </n-flex>

  <main class="mt-2 mb-8">
    <n-data-table class="w-full" :columns="columns" :data="filteredWords" virtual-scroll>
    </n-data-table>
  </main>
</template>
<script setup lang="tsx">
import { useCensorStore } from '~/components/censor/censor';
import SensitiveTag from '~/components/censor/sensitive-tag.tsx';
import { getCensorWords } from '~/api/v1/censor';
import type { DataTableColumns } from 'naive-ui';

const columns: DataTableColumns<SensitiveWord> = [
  {
    key: 'level',
    title: '级别',
    render: ({ level }) => {
      switch (level) {
        case 1:
          return <SensitiveTag type="default" />;
        case 2:
          return <SensitiveTag type="info" />;
        case 3:
          return <SensitiveTag type="warning" />;
        case 4:
          return <SensitiveTag type="error" />;
        default:
          return <SensitiveTag type="default" message="未知" />;
      }
    },
  },
  {
    key: 'related',
    title: '匹配词汇',
    render: ({ related, main }) => {
      if (related) {
        return (
          <n-flex size="small">
            {related.map((word: { word: string }) => (
              <n-text key={word.word}>{word.word}</n-text>
            ))}
          </n-flex>
        );
      } else {
        return (
          <n-flex>
            <n-text>{main}</n-text>
          </n-flex>
        );
      }
    },
  },
];

onMounted(async () => {
  await refreshWords();
});

const censorStore = useCensorStore();

interface SensitiveRelatedWord {
  word: string;
  reason: string;
}

interface SensitiveWord {
  main: string;
  level: number;
  related: SensitiveRelatedWord[];
}

const words = ref<SensitiveWord[]>([]);
const filter = ref<string>('');
const filterCount = computed(() => words.value.length - filteredWords.value.length);
const filteredWords = computed(() =>
  words.value.filter(word => {
    if (filter.value === '') {
      return true;
    }
    const val = filter.value.toLowerCase();
    return (
      word.main.toLowerCase().includes(val) ||
      word.related.map(w => w.word.toLowerCase().includes(val)).includes(true)
    );
  }),
);

censorStore.$subscribe(async (_, state) => {
  if (state.wordsNeedRefresh === true) {
    await refreshWords();
    state.wordsNeedRefresh = false;
  }
});

const refreshWords = async () => {
  const c:
    | { result: false }
    | {
        result: true;
        data: SensitiveWord[];
      } = await getCensorWords();
  if (c.result) {
    words.value = c.data;
  }
};
</script>

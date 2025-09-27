<template>
  <div class="flex flex-wrap gap-4">
    <n-card title="公骰设置">
      <template #header-extra>
        <n-flex size="large" align="center">
          <n-switch v-model:value="config.publicDiceEnable" @on-update:value="enableChange">
            <template #checked>启用</template>
            <template #unchecked>关闭</template>
          </n-switch>
          <n-button :disabled="!config.publicDiceEnable" type="primary" @click="doSave">
            <template #icon>
              <n-icon><i-carbon-save /></n-icon>
            </template>
            保存
          </n-button>
        </n-flex>
      </template>

      <div class="w-full flex flex-col md:flex-row gap-4">
        <div v-if="isSmallWindow" class="w-full">
          <div :class="{ disabledOverlay: !config.publicDiceEnable }">
            <n-avatar
              shape="square"
              style="width: auto; height: auto; vertical-align: top"
              fit="contain"
              :src="imgSeal" />
          </div>
        </div>
        <div v-else class="w-full md:w-1/5">
          <div :class="{ disabledOverlay: !config.publicDiceEnable }">
            <n-avatar
              shape="square"
              style="width: auto; height: auto; vertical-align: top"
              fit="contain"
              :src="imgSeal" />
          </div>
        </div>

        <div class="flex-1">
          <n-form class="flex flex-col h-full">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-5 w-full">
              <n-form-item class="w-full">
                <template #label>
                  <div class="flex items-center gap-1">
                    <span>公骰 UID</span>
                    <n-tooltip placement="left">
                      <template #trigger>
                        <n-icon class="cursor-help"><i-carbon-help-filled /></n-icon>
                      </template>
                      <div style="width: 10rem">
                        公骰UID是上报公骰所使用的密钥， 通常情况下留空让系统自动生成，
                        请勿随意将公骰的UID展示给他人
                      </div>
                    </n-tooltip>
                  </div>
                </template>
                <n-input
                  type="password"
                  v-model:value="config.publicDiceId"
                  :disabled="!config.publicDiceEnable"
                  placeholder="留空则会自动注册" />
              </n-form-item>

              <n-form-item class="w-full">
                <template #label>
                  <div class="flex items-center gap-1">
                    <span>公骰昵称</span>
                    <n-tooltip placement="left">
                      <template #trigger>
                        <n-icon class="cursor-help"><i-carbon-help-filled /></n-icon>
                      </template>
                      <div style="width: 10rem">公骰昵称是展示在公骰列表的昵称</div>
                    </n-tooltip>
                  </div>
                </template>
                <n-input
                  v-model:value="config.publicDiceName"
                  :disabled="!config.publicDiceEnable"
                  placeholder="请输入公骰昵称" />
              </n-form-item>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-5 w-full mt-5">
              <n-form-item class="w-full">
                <template #label>
                  <div class="flex items-center gap-1">
                    <span>公骰头像</span>
                    <n-tooltip placement="left">
                      <template #trigger>
                        <n-icon class="cursor-help"><i-carbon-help-filled /></n-icon>
                      </template>
                      <div style="width: 10rem">公骰头像是展示在公骰列表的头像</div>
                    </n-tooltip>
                  </div>
                </template>
                <n-input
                  v-model:value="config.publicDiceAvatar"
                  :disabled="!config.publicDiceEnable"
                  placeholder="请输入公骰头像Url" />
              </n-form-item>

              <n-form-item class="w-full">
                <template #label>
                  <div class="flex items-center gap-1">
                    <span>骰主留言</span>
                    <n-tooltip placement="left">
                      <template #trigger>
                        <n-icon class="cursor-help"><i-carbon-help-filled /></n-icon>
                      </template>
                      <div style="width: 10rem">骰主留言是展示在公骰列表的留言</div>
                    </n-tooltip>
                  </div>
                </template>
                <n-input
                  v-model:value="config.publicDiceNote"
                  :disabled="!config.publicDiceEnable"
                  placeholder="请输入你的留言" />
              </n-form-item>
            </div>

            <div class="w-full mt-5 flex-1">
              <n-form-item class="w-full h-full">
                <template #label>
                  <div class="flex items-center gap-1">
                    <span>公骰简介</span>
                    <n-tooltip placement="left">
                      <template #trigger>
                        <n-icon class="cursor-help"><i-carbon-help-filled /></n-icon>
                      </template>
                      <div style="width: 10rem">公骰简介是展示在公骰列表的简介</div>
                    </n-tooltip>
                  </div>
                </template>
                <n-input
                  v-model:value="config.publicDiceBrief"
                  :disabled="!config.publicDiceEnable"
                  placeholder="请输入简介"
                  type="textarea" />
              </n-form-item>
            </div>
          </n-form>
        </div>
      </div>

      <template #footer>
        <div :class="{ disabledOverlay: !config.publicDiceEnable }">
          <span style="font-size: 1.2rem; font-weight: bold">选择要上报的终端</span>
          <n-data-table
            class="mt-4"
            :columns="columns"
            :data="tableData"
            :row-key="rowKey"
            checkable
            v-model:checked-row-keys="checkedRowKeys"
            @update:checked-rows="handleSelectionChange" />
        </div>
      </template>
    </n-card>
  </div>
</template>
<script lang="ts" setup>
import imgSeal from '~/assets/seal.png';
import { getDicePublicInfo, setDicePublicInfo } from '~/api/v1/public_dice';
import { breakpointsTailwind, useBreakpoints } from '@vueuse/core';
import { useMessage, type DataTableColumns } from 'naive-ui';

const breakpoints = useBreakpoints(breakpointsTailwind);
const isSmallWindow = breakpoints.smallerOrEqual('md');
const message = useMessage();

const config = ref<any>({});
const tableData = ref<any>([]);
let selected: any[] = [];

const rowKey = (row: any) => row.id as string | number;
let checkedRowKeys = ref<Array<string | number>>([]);

const columns: DataTableColumns<any> = [
  { type: 'selection' },
  { title: '账号', key: 'userId', sorter: 'default' },
  { title: '平台', key: 'platform', sorter: 'default' },
  { title: '协议', key: 'adapter', sorter: 'default' },
  { title: '状态', key: 'state', sorter: 'default' },
];

const handleSelectionChange = (rows: any[]) => {
  selected = rows;
};

const enableChange = async (value: string | number | boolean) => {
  config.value.publicDiceEnable = value;
  await setDicePublicInfo(config.value, selected);
  await refreshInfo();
};
const doSave = async () => {
  await setDicePublicInfo(config.value, selected);
  await refreshInfo();
  message.success('已保存');
};

const refreshInfo = async () => {
  tableData.value = [];
  const infos = await getDicePublicInfo();
  config.value = infos.config;
  if (infos.endpoints !== null) {
    infos.endpoints.forEach((dc: any) => {
      const state = () => {
        switch (dc.state) {
          case 0:
            return '断开';
          case 1:
            return '已连接';
          case 2:
            return '连接中';
          case 3:
            return '连接失败';
          default:
            return '未知';
        }
      };
      const adapter = () => {
        if (
          dc.platform === 'QQ' &&
          dc.protocolType === 'onebot' &&
          dc.adapter.builtinMode !== null
        ) {
          switch (dc.adapter.builtinMode) {
            case 'lagrange':
              return '内置客户端';
            case 'lagrange-gocq':
              return '内置gocq';
            case 'gocq':
              return '分离部署';
          }
        }
        return '-';
      };
      tableData.value.push({
        userId: dc.userId,
        platform: dc.platform,
        adapter: adapter(),
        state: state(),
        isPublic: dc.isPublic,
        id: dc.id,
      });
    });

    checkedRowKeys.value = tableData.value.filter((r: any) => r.isPublic).map((r: any) => r.id);
    selected = tableData.value.filter((r: any) => r.isPublic);
  }
};

onBeforeMount(async () => {
  await refreshInfo();
});
</script>

<style scoped lang="css">
.disabledOverlay {
  filter: grayscale(1);
  opacity: 0.6;
  pointer-events: none;
}
</style>

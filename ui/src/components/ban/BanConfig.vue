<script setup lang="ts">
import { getBanConfig, setBanConfig } from '~/api/v1/banconfig';
import type { BanConfig } from '#';
import { useMessage } from 'naive-ui';

const message = useMessage();

const banConfig = ref<BanConfig>({} as BanConfig);
const modified = ref<boolean>(false);

const banConfigSave = async () => {
  await setBanConfig(banConfig.value);
  await configGet();
  message.success('已保存');
  modified.value = false;
  await nextTick(() => {
    modified.value = false;
  });
};

const configGet = async () => {
  banConfig.value = await getBanConfig();
};

onBeforeMount(async () => {
  await configGet();
  modified.value = false;
});

watch(
  banConfig,
  () => {
    modified.value = true;
  },
  { deep: true },
);
</script>

<template>
  <header>
    <n-button type="info" secondary @click="banConfigSave">
      <template #icon>
        <n-icon><i-carbon-save /></n-icon>
      </template>
      保存设置
    </n-button>
    <tip-box v-if="modified" type="error" class="my-4">
      <n-text type="error" tag="strong" class="text-lg">内容已修改，不要忘记保存！</n-text>
    </tip-box>
  </header>

  <h4>基本设置</h4>
  <n-flex>
    <n-text>黑名单惩罚：</n-text>
    <n-checkbox v-model:checked="banConfig.banBehaviorRefuseReply">拒绝回复</n-checkbox>
    <n-checkbox v-model:checked="banConfig.banBehaviorRefuseInvite">拒绝邀请</n-checkbox>
    <n-checkbox v-model:checked="banConfig.banBehaviorQuitLastPlace">退出事发群</n-checkbox>
    <!-- <div>自动拉黑时长(分钟): <el-input style="max-width: 5rem;" type="number" v-model="banConfig.autoBanMinutes"></el-input></div> -->
    <n-checkbox v-model:checked="banConfig.banBehaviorQuitPlaceImmediately">
      使用时立即退出群
    </n-checkbox>
    <n-checkbox v-model:checked="banConfig.banBehaviorQuitIfAdmin">
      使用者为管理员立即退群，为普通群员进行通告
    </n-checkbox>
    <n-checkbox v-model:checked="banConfig.banBehaviorQuitIfAdminSilentIfNotAdmin">
      使用者为管理员立即退群，为普通群员仅拒绝回复
    </n-checkbox>
  </n-flex>

  <h4>怒气值设置</h4>
  <tip-box type="info" class="my-4">
    <n-text type="info">
      说明：海豹的黑名单使用积分制，每当用户做出恶意行为，其积分上涨一定数值，到达阈值后自动进入黑名单。会通知邀请者、通知列表、事发群（如果可能）。
    </n-text>
  </tip-box>
  <n-form size="small" label-placement="left" label-width="auto">
    <n-form-item label="警告阈值">
      <n-input-number v-model:value="banConfig.thresholdWarn" :min="0" :step="1" :precision="0" />
    </n-form-item>
    <n-form-item label="拉黑阈值">
      <n-input-number v-model:value="banConfig.thresholdBan" :min="0" :step="1" :precision="0" />
    </n-form-item>

    <n-form-item class="mt-10" label="禁言增加">
      <n-input-number v-model:value="banConfig.scoreGroupMuted" :min="0" :step="1" :precision="0" />
    </n-form-item>
    <n-form-item label="踢出增加">
      <n-input-number
        v-model:value="banConfig.scoreGroupKicked"
        :min="0"
        :step="1"
        :precision="0" />
    </n-form-item>
    <n-form-item label="刷屏增加">
      <n-input-number
        v-model:value="banConfig.scoreTooManyCommand"
        :min="0"
        :step="1"
        :precision="0" />
    </n-form-item>
    <n-form-item label="每分钟下降">
      <n-input-number
        v-model:value="banConfig.scoreReducePerMinute"
        :min="0"
        :step="1"
        :precision="0" />
    </n-form-item>

    <n-form-item class="mt-10" label="群组连带责任">
      <n-input-number
        v-model:value="banConfig.jointScorePercentOfGroup"
        :min="0"
        :max="1"
        :step="0.1" />
    </n-form-item>
    <n-form-item label="邀请人连带责任">
      <n-input-number
        v-model:value="banConfig.jointScorePercentOfInviter"
        :min="0"
        :max="1"
        :step="0.1" />
    </n-form-item>
  </n-form>
</template>

<style scoped lang="css"></style>

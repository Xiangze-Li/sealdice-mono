<template>
  <n-affix v-if="modified" :top="120" :trigger-top="60" class="w-full">
    <tip-box type="error">
      <n-text type="error" tag="strong" class="text-lg">内容已修改，不要忘记保存！</n-text>
    </tip-box>
  </n-affix>

  <n-form label-placement="left" label-width="auto">
    <h2>Master 管理</h2>
    <n-form-item label="">
      <template #label>
        <span>Master 解锁码</span>
        <n-tooltip>
          <template #trigger>
            <n-icon><i-carbon-help-filled /></n-icon>
          </template>
          用于误操作或被抢占 Master。<br />
          执行这条指令可以直接获取 Master 权限并踢掉其他所有人，指令有效期 20 分钟
        </n-tooltip>
      </template>
      <div>
        <n-button v-if="!isShowUnlockCode" @click="isShowUnlockCode = true">查看</n-button>
        <div v-else style="font-weight: bold">.master unlock {{ config.masterUnlockCode }}</div>
      </div>
    </n-form-item>

    <n-form-item label="Master列表">
      <template #label>
        <span>Master 列表</span>
        <n-tooltip>
          <template #trigger>
            <n-icon><i-carbon-help-filled /></n-icon>
          </template>
          单行格式：QQ:12345<br />也可添加群 ID, 在指定的群中所有人均视为有 Master 权限
        </n-tooltip>
      </template>
      <template v-if="config.diceMasters && config.diceMasters.length">
        <n-flex>
          <div
            :key="index"
            v-for="(k2, index) in config.diceMasters"
            style="width: 100%; margin-bottom: 0.5rem">
            <!-- 这里面是单条修改项 -->
            <div style="display: flex">
              <div>
                <!-- :suffix-icon="Management" -->
                <n-input
                  v-model:value="config.diceMasters[index]"
                  autosize
                  placeholder=""
                  style="min-width: 10rem" />
              </div>
              <div style="display: flex; align-items: center; width: 1.3rem; margin-left: 1rem">
                <n-tooltip placement="bottom-start">
                  <template #trigger>
                    <n-icon>
                      <i-carbon-add-filled
                        v-if="index == 0"
                        @click="config.diceMasters = addItem(config.diceMasters)" />
                      <i-carbon-close-outline
                        v-else
                        @click="removeItem(config.diceMasters, index)" />
                    </n-icon>
                  </template>
                  {{ index === 0 ? '点击添加项目' : '点击删除你不想要的项' }}
                </n-tooltip>
              </div>
            </div>
          </div>
        </n-flex>
      </template>
      <template v-else>
        <n-icon>
          <i-carbon-add-filled @click="config.diceMasters = addItem(config.diceMasters)" />
        </n-icon>
      </template>
    </n-form-item>

    <n-form-item label="消息通知列表">
      <template #label>
        <div>
          <span>消息通知列表</span>
          <n-tooltip>
            <template #trigger>
              <n-icon><i-carbon-help-filled /></n-icon>
            </template>
            会对以下消息进行通知：<br />
            加群邀请、好友邀请、进入群组、被踢出群、被禁言、自动激活、指令退群<br />
            单行格式：QQ:12345 QQ-Group:12345 Mail:abc@foo.bar<br />
            通知列表中的 QQ 号在发件时会自动转换成对应邮箱
          </n-tooltip>
        </div>
      </template>

      <template v-if="config.noticeIds && config.noticeIds.length">
        <n-flex>
          <div
            :key="index"
            v-for="(k2, index) in config.noticeIds"
            style="width: 100%; margin-bottom: 0.5rem">
            <!-- 这里面是单条修改项 -->
            <div style="display: flex">
              <div>
                <!-- :suffix-icon="Management" -->
                <n-input
                  v-model:value="config.noticeIds[index]"
                  autosize
                  placeholder=""
                  style="min-width: 10rem" />
              </div>
              <div style="display: flex; align-items: center; width: 1.3rem; margin-left: 1rem">
                <n-tooltip placement="bottom-start">
                  <template #trigger>
                    <n-icon>
                      <i-carbon-add-filled
                        v-if="index == 0"
                        @click="config.noticeIds = addItem(config.noticeIds)" />
                      <i-carbon-close-outline v-else @click="removeItem(config.noticeIds, index)" />
                    </n-icon>
                  </template>
                  {{ index === 0 ? '点击添加项目' : '点击删除你不想要的项' }}
                </n-tooltip>
              </div>
            </div>
          </div>
        </n-flex>
      </template>
      <template v-else>
        <n-icon>
          <i-carbon-add-filled @click="config.noticeIds = addItem(config.noticeIds)" />
        </n-icon>
      </template>
    </n-form-item>

    <n-form-item label="邮箱通知">
      <template #label>
        <span>邮箱通知</span>
        <n-tooltip>
          <template #trigger>
            <n-icon><i-carbon-help-filled /></n-icon>
          </template>
          为处理某些平台通知频繁引起的问题，启用后 <strong>只会通过邮件</strong> 进行通知
        </n-tooltip>
      </template>
      <n-checkbox v-model:checked="config.mailEnable" label="开启" />
    </n-form-item>
    <n-form-item label="发件邮箱">
      <template #label>
        <span>发件邮箱</span>
        <n-tooltip>
          <template #trigger>
            <n-icon><i-carbon-help-filled /></n-icon>
          </template>
          特别提醒：QQ 号被冻结后对应账号的邮箱也无法使用
        </n-tooltip>
      </template>
      <n-input v-model:value="config.mailFrom" placeholder="" style="width: 12rem" />
    </n-form-item>
    <n-form-item label="邮箱密钥">
      <template #label>
        <span>邮箱密钥</span>
        <n-tooltip>
          <template #trigger>
            <n-icon><i-carbon-help-filled /></n-icon>
          </template>
          请注意不一定是密码，例如 QQ 邮箱就需要在设置界面生成授权码
        </n-tooltip>
      </template>
      <n-input v-model:value="config.mailPassword" placeholder="" style="width: 12rem" />
    </n-form-item>
    <n-form-item label="SMTP 地址">
      <template #label>
        <span>SMTP 地址</span>
        <n-tooltip>
          <template #trigger>
            <n-icon><i-carbon-help-filled /></n-icon>
          </template>
          例如 smtp.qq.com
        </n-tooltip>
      </template>
      <n-input v-model:value="config.mailSmtp" placeholder="" style="width: 12rem" />
    </n-form-item>
    <n-form-item>
      <template #label>
        <span>发送测试邮件</span>
        <n-tooltip>
          <template #trigger>
            <n-icon><i-carbon-help-filled /></n-icon>
          </template>
          向通知列表发送测试邮件，便于检查邮件相关配置
        </n-tooltip>
      </template>
      <n-button type="primary" @click="mailTest">发送</n-button>
    </n-form-item>

    <n-form-item>
      <template #label>
        <span>私骰模式</span>
        <n-tooltip>
          <template #trigger>
            <n-icon><i-carbon-help-filled /></n-icon>
          </template>
          只允许信任用户拉入群聊
        </n-tooltip>
      </template>
      <n-checkbox v-model:checked="config.trustOnlyMode" label="开启" />
    </n-form-item>

    <n-form-item>
      <template #label>
        <span>允许自由开关</span>
        <n-tooltip>
          <template #trigger>
            <n-icon><i-carbon-help-filled /></n-icon>
          </template>
          默认开启。允许任何人执行 bot on/off和ext on/off，否则只有邀请者、管理员和master进行操作
        </n-tooltip>
      </template>
      <n-checkbox v-model:checked="config.botExtFreeSwitch" label="开启" />
    </n-form-item>
    <n-form-item>
      <template #label>
        <span>启用戳一戳</span>
        <n-tooltip>
          <template #trigger>
            <n-icon><i-carbon-help-filled /></n-icon>
          </template>
          默认开启。仅 QQ 平台，关闭后不回应戳一戳。
        </n-tooltip>
      </template>
      <n-flex vertical size="small">
        <n-checkbox v-model:checked="config.QQEnablePoke" label="开启" />
        <n-text type="warning" class="text-xs">
          启用前请确认你使用的 QQ 连接方式支持该功能，若不支持请关闭该功能来避免日志中出现相关报错。
        </n-text>
      </n-flex>
    </n-form-item>

    <n-form-item>
      <template #label>
        <span>限制.text 指令</span>
        <n-tooltip>
          <template #trigger>
            <n-icon><i-carbon-help-filled /></n-icon>
          </template>
          默认开启。开启后只有 Master 和信任用户可使用 .text 指令。<br />
          如果你不了解这个指令，一定要开启。
        </n-tooltip>
      </template>
      <n-checkbox v-model:checked="config.textCmdTrustOnly" label="开启" />
    </n-form-item>

    <n-form-item>
      <template #label>
        <span>忽略.bot 裸指令</span>
        <n-tooltip>
          <template #trigger>
            <n-icon><i-carbon-help-filled /></n-icon>
          </template>
          默认关闭。开启后在群聊中只有@骰子才会响应.bot
        </n-tooltip>
      </template>
      <n-checkbox v-model:checked="config.ignoreUnaddressedBotCmd" label="开启" />
    </n-form-item>

    <n-form-item>
      <template #label>
        <span>存活确认 (骰狗)</span>
        <n-tooltip>
          <template #trigger>
            <n-icon><i-carbon-help-filled /></n-icon>
          </template>
          定期向通知列表发送消息，以便于骰主知晓存活状态
        </n-tooltip>
      </template>
      <n-checkbox v-model:checked="config.aliveNoticeEnable" label="开启" />
    </n-form-item>

    <n-form-item>
      <template #label>
        <div>
          <span>存活消息间隔</span>
          <n-tooltip>
            <template #trigger>
              <n-icon><i-carbon-help-filled /></n-icon>
            </template>
            间隔写法请参阅
            <a href="https://pkg.go.dev/github.com/robfig/cron" target="_blank">cron 文档</a><br />
            注意：重启骰子后重新计时
          </n-tooltip>
        </div>
      </template>
      <n-input v-model:value="config.aliveNoticeValue" style="width: 12rem" />
    </n-form-item>

    <n-form-item>
      <template #label>
        <span>日志记录提示</span>
        <n-tooltip>
          <template #trigger>
            <n-icon><i-carbon-help-filled /></n-icon>
          </template>
          每记录 N 条文本后，主动发送一条提醒信息，避免忘记 log off
        </n-tooltip>
      </template>
      <n-flex align="center">
        <n-checkbox v-model:checked="config.logSizeNoticeEnable" label="开启" />
        <n-input-number v-model:value="config.logSizeNoticeCount" clearable style="width: 14rem" />
      </n-flex>
    </n-form-item>

    <h2>刷屏警告</h2>
    <n-form-item>
      <template #label>
        <span>刷屏警告开关</span>
        <n-tooltip>
          <template #trigger>
            <n-icon><i-carbon-help-filled /></n-icon>
          </template>
          默认关闭。开启后会对使用指令过快的用户和群组进行警告，警告后继续使用指令会增加怒气值，只对
          QQ 平台有效
        </n-tooltip>
      </template>
      <n-checkbox v-model:checked="config.rateLimitEnabled" label="开启" />
    </n-form-item>
    <h4>刷屏警告速率</h4>

    <tip-box type="info" class="my-4">
      <span>刷屏警告工作原理如下：</span>
      <div class="ml-4">
        <ul class="list-disc">
          <li>每群每用户独立有一个装令牌的桶，桶最多能装「上限」枚令牌</li>
          <li>每次指令视作拿走一枚令牌</li>
          <li>当桶里没有令牌时，试图拿走令牌将被阻止（触发警告）</li>
          <li>桶以「速率」自动补充令牌</li>
        </ul>
      </div>
      <div>处置时，优先惩罚个人刷屏，其次是群组内许多人一起刷屏。</div>
      <div>如果您感觉难以理解，为了稳定性还是不要更改比较好！</div>
      <div>对速率限制的所有更改 <strong>重启后生效</strong></div>
    </tip-box>

    <n-form-item>
      <template #label>
        <span>个人速率</span>
        <n-tooltip>
          <template #trigger>
            <n-icon><i-carbon-help-filled /></n-icon>
          </template>
          支持以下格式：<br />@every 3s 每 3 秒一个<br />3 每秒 3 个
        </n-tooltip>
      </template>
      <n-input
        v-model:value="config.personalReplenishRate"
        clearable
        placeholder=""
        style="width: 14rem" />
    </n-form-item>
    <n-form-item label="个人上限">
      <n-input-number v-model:value="config.personalBurst" :min="1" :step="1" clearable />
    </n-form-item>
    <n-form-item>
      <template #label>
        <span>群组速率</span>
        <n-tooltip>
          <template #trigger>
            <n-icon><i-carbon-help-filled /></n-icon>
          </template>
          支持以下格式：<br />@every 3s 每 3 秒一个<br />3 每秒 3 个
        </n-tooltip>
      </template>
      <n-input
        v-model:value="config.groupReplenishRate"
        clearable
        placeholder=""
        style="width: 14rem" />
    </n-form-item>
    <n-form-item label="群组上限">
      <n-input-number v-model:value="config.groupBurst" :min="1" :step="1" clearable />
    </n-form-item>

    <h2>访问控制</h2>
    <n-form-item label="UI界面地址">
      <template #label>
        <span>UI 界面地址</span>
        <n-tooltip>
          <template #trigger>
            <n-icon><i-carbon-help-filled /></n-icon>
          </template>
          0.0.0.0:3211 主要用于服务器，代表可以在公网中用你的手机和电脑访问 <br />
          127.0.0.1:3211 主要用于自己的电脑/手机，只能在当前设备上管理海豹<br />
          注意：重启骰子后生效！<br />
          另，想开多个海豹必须修改端口号！
        </n-tooltip>
      </template>
      <n-select
        v-model:value="config.serveAddress"
        :options="links"
        filterable
        clearable
        tag
        placeholder=""
        style="width: 14rem" />
    </n-form-item>

    <n-form-item label="UI界面密码">
      <template #label>
        <span>UI 界面密码</span>
        <n-tooltip>
          <template #trigger>
            <n-icon><i-carbon-help-filled /></n-icon>
          </template>
          公网用户一定要加，登录后会自动记住一段时间！
        </n-tooltip>
      </template>
      <n-input
        v-model:value="config.uiPassword"
        type="password"
        show-password-on="mousedown"
        placeholder=""
        style="width: 14rem" />
    </n-form-item>

    <h2>QQ 频道设置</h2>
    <n-form-item>
      <template #label>
        <span>总开关</span>
        <n-tooltip>
          <template #trigger>
            <n-icon><i-carbon-help-filled /></n-icon>
          </template>
          如果关闭，将忽略任何频道消息
        </n-tooltip>
      </template>
      <n-checkbox v-model:checked="config.workInQQChannel" label="开启" />
    </n-form-item>

    <n-form-item>
      <template #label>
        <span>自动 bot on</span>
        <n-tooltip>
          <template #trigger>
            <n-icon><i-carbon-help-filled /></n-icon>
          </template>
          如果开启，需要在每一个子频道手动 bot off，推荐关闭
        </n-tooltip>
      </template>
      <n-checkbox v-model:checked="config.QQChannelAutoOn" label="开启" />
    </n-form-item>

    <n-form-item>
      <template #label>
        <span>记录消息日志</span>
        <n-tooltip>
          <template #trigger>
            <n-icon><i-carbon-help-filled /></n-icon>
          </template>
          是否记录频道消息到日志，如果频道较多，可能造成严重刷屏。<br />
          若关闭则仅在日志记录指令，推荐关闭
        </n-tooltip>
      </template>
      <n-checkbox v-model:checked="config.QQChannelLogMessage" label="开启" />
    </n-form-item>

    <h2>游戏</h2>
    <n-form-item>
      <template #label>
        <span>COC 默认房规</span>
        <n-tooltip>
          <template #trigger>
            <n-icon><i-carbon-help-filled /></n-icon>
          </template>
          可设置为 0-5，以及 dg（DeltaGreen）
        </n-tooltip>
      </template>
      <n-input
        v-model:value="config.defaultCocRuleIndex"
        clearable
        placeholder=""
        style="width: 14rem" />
    </n-form-item>

    <n-form-item>
      <template #label>
        <span>COC 制卡上限</span>
        <n-tooltip>
          <template #trigger>
            <n-icon><i-carbon-help-filled /></n-icon>
          </template>
          .coc n 中 n 的最大值，1-12 之间，默认 5
        </n-tooltip>
      </template>
      <n-input v-model:value="config.maxCocCardGen" clearable placeholder="" style="width: 14rem" />
    </n-form-item>

    <n-form-item>
      <template #label>
        <span>骰点轮数上限</span>
        <n-tooltip>
          <template #trigger>
            <n-icon><i-carbon-help-filled /></n-icon>
          </template>
          .r n#中n的最大值，1-25之间，默认12
        </n-tooltip>
      </template>
      <n-input
        v-model:value="config.maxExecuteTime"
        clearable
        placeholder=""
        style="width: 14rem" />
    </n-form-item>

    <h2>海豹</h2>
    <tip-box type="info" class="my-2">
      <div>使用指定的压缩包对当前海豹进行覆盖，上传完成后会自动重启海豹。</div>
      <div>请注意尽量不要从高版本降低到低版本，数据库有可能不兼容。</div>
      <div>选择文件确定后请耐心等待反馈。</div>
    </tip-box>

    <n-form-item>
      <template #label>
        <span>保护开关</span>
      </template>
      <n-checkbox v-model:checked="isUploadEnable" label="我已阅读功能描述" />
    </n-form-item>

    <n-form-item>
      <template #label>
        <span>固件升级</span>
        <n-tooltip>
          <template #trigger>
            <n-icon><i-carbon-help-filled /></n-icon>
          </template>
          使用海豹安装包进行覆盖升级
        </n-tooltip>
      </template>
      <n-flex v-if="store.curDice.baseInfo.containerMode" wrap>
        <n-tooltip>
          <template #trigger>
            <n-button type="primary" disabled>
              <template #icon>
                <i-carbon-upload />
              </template>
              上传压缩包
            </n-button>
          </template>
          容器模式下固件升级被禁用，请手动拉取最新镜像
        </n-tooltip>
        <n-text type="warning" class="text-xs">容器模式下固件升级被禁用，请手动拉取最新镜像</n-text>
      </n-flex>
      <n-upload
        v-else
        class="upload"
        action=""
        @before-upload="beforeUpload"
        :file-list="fileList"
        :disabled="!isUploadEnable">
        <n-button type="primary" :disabled="!isUploadEnable">
          <template #icon>
            <i-carbon-upload />
          </template>
          上传压缩包
        </n-button>
      </n-upload>
    </n-form-item>

    <h2>其他</h2>
    <n-form-item label="加好友验证信息">
      <template #label>
        <span>加好友验证</span>
        <n-tooltip>
          <template #trigger>
            <n-icon><i-carbon-help-filled /></n-icon>
          </template>
          加好友时必须输入正确的验证信息才能通过<br />
          注意：若使用“回答问题并由我确认”，只写问题答案，有多个答案用空格隔开：<br />
          问题 1 答案 问题 2 答案<br />
          注意问题答案中本身不能有空格
        </n-tooltip>
      </template>
      <n-input
        v-model:value="config.friendAddComment"
        type="text"
        clearable
        style="width: auto"
        :placeholder="'为空则填写任意均可通过'" />
    </n-form-item>

    <n-form-item label="不活跃N天后自动退群">
      <template #label>
        <span>自动退群阈值</span>
        <n-tooltip>
          <template #trigger>
            <n-icon><i-carbon-help-filled /></n-icon>
          </template>
          每天凌晨 4 时检查所有群聊的不活跃天数，若超过设定天数则自动退群<br />
          设置为 0 以关闭该功能
        </n-tooltip>
      </template>
      <n-input-number
        v-model:value="config.quitInactiveThreshold"
        style="width: auto"
        :placeholder="'不活跃N天后自动退出'" />
    </n-form-item>

    <n-form-item label="退群批次大小">
      <n-input-number v-model:value="config.quitInactiveBatchSize" style="width: auto" />
    </n-form-item>

    <n-form-item label="退群批次间隔(分)">
      <n-input-number v-model:value="config.quitInactiveBatchWait" style="width: auto" />
    </n-form-item>

    <n-form-item label="QQ回复延迟(秒)">
      <n-input-number
        v-model:value="config.messageDelayRangeStart"
        :min="0"
        :max="config.messageDelayRangeEnd"
        style="width: 6rem" />
      <span style="margin: 0 1rem">-</span>
      <n-input-number
        v-model:value="config.messageDelayRangeEnd"
        :min="config.messageDelayRangeStart"
        style="width: 6rem" />
    </n-form-item>

    <n-form-item>
      <template #label>
        <span>{{ `<玩家名>` }}外框</span>
        <n-tooltip>
          <template #trigger>
            <n-icon><i-carbon-help-filled /></n-icon>
          </template>
          默认开启。建议开启，这可以防止用户改名为指令 (如.bot)，并利用播报去唤醒其他骰子造成刷屏。
        </n-tooltip>
      </template>
      <n-checkbox
        v-model:checked="config.playerNameWrapEnable"
        label="开启"
        @update:checked="nameWrapUncheck" />
    </n-form-item>

    <n-form-item label="日志仅记录指令">
      <n-checkbox v-model:checked="config.onlyLogCommandInGroup" label="在群聊中" />
      <n-checkbox v-model:checked="config.onlyLogCommandInPrivate" label="在私聊中" />
    </n-form-item>

    <n-form-item label="拒绝加入新群">
      <n-checkbox
        v-model:checked="config.refuseGroupInvite"
        label="拒绝加群(仅在非强制拉入时起效)" />
    </n-form-item>

    <n-form-item label="指令前缀">
      <template #label>
        <span>指令前缀</span>
        <n-tooltip>
          <template #trigger>
            <n-icon><i-carbon-help-filled /></n-icon>
          </template>
          举例：添加！作为指令前缀，运行 !r 可以骰点
        </n-tooltip>
      </template>
      <template v-if="config.commandPrefix && config.commandPrefix.length">
        <n-flex>
          <div
            :key="index"
            v-for="(k2, index) in config.commandPrefix"
            style="width: 100%; margin-bottom: 0.5rem">
            <!-- 这里面是单条修改项 -->
            <div style="display: flex">
              <div>
                <n-input
                  v-model:value="config.commandPrefix[index]"
                  autosize
                  placeholder=""
                  style="width: 14rem" />
              </div>
              <div style="display: flex; align-items: center; width: 1.3rem; margin-left: 1rem">
                <n-tooltip placement="bottom-start">
                  <template #trigger>
                    <n-icon>
                      <i-carbon-add-filled
                        v-if="index == 0"
                        @click="config.commandPrefix = addItem(config.commandPrefix)" />
                      <i-carbon-close-outline
                        v-else
                        @click="removeItem(config.commandPrefix, index)" />
                    </n-icon>
                  </template>
                  {{ index === 0 ? '点击添加项目' : '点击删除你不想要的项目' }}
                </n-tooltip>
              </div>
            </div>
          </div>
        </n-flex>
      </template>
      <template v-else>
        <n-icon>
          <i-carbon-add-filled @click="config.commandPrefix = addItem(config.commandPrefix)" />
        </n-icon>
      </template>
    </n-form-item>

    <h2>扩展及扩展指令</h2>
    <n-list>
      <n-list-item :key="i.name" v-for="i in config.extDefaultSettings">
        <n-flex vertical class="mx-4">
          <span>扩展：{{ i.name }}</span>
          <span>禁用指令</span>
          <n-flex size="small">
            <n-button
              :key="k"
              v-for="(v, k) in i.disabledCommand"
              :type="v ? 'error' : 'default'"
              size="tiny"
              @click="i.disabledCommand[k] = !v">
              {{ k }}
            </n-button>
          </n-flex>
          <div>
            <n-checkbox v-model:checked="i.autoActive">入群自动开启</n-checkbox>
          </div>
        </n-flex>
      </n-list-item>
    </n-list>

    <n-form-item v-if="modified" label="" style="margin-top: 3rem">
      <n-flex>
        <n-button type="error" @click="submitGiveup">放弃改动</n-button>
        <n-button type="success" @click="submit">保存设置</n-button>
      </n-flex>
    </n-form-item>
  </n-form>
</template>

<script lang="ts" setup>
import { cloneDeep, toNumber } from 'es-toolkit/compat';
import { useDialog, useMessage, type UploadFileInfo } from 'naive-ui';
import { postMailTest, postUploadToUpgrade } from '~/api/v1/dice';
import { useStore } from '~/store';
import { objDiff, passwordHash } from '~/utils';

const store = useStore();
const message = useMessage();
const dialog = useDialog();

const config = ref<any>({});
const fileList = ref<any[]>([]);

const isShowUnlockCode = ref(false);
const modified = ref(false);
const isUploadEnable = ref(false);

const beforeUpload = async (data: { file: UploadFileInfo }) => {
  // UploadRawFile
  const file = data.file.file as File;
  const fd = new FormData();
  fd.append('files', file);
  try {
    message.info('开始上传固件包，更新时程序可能离线');
    await postUploadToUpgrade(file);
  } catch (e: any) {
    message.error(e.toString());
    console.log(e);
  }
};

watch(
  () => config,
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  (newValue, oldValue) => {
    //直接监听
    modified.value = true;
  },
  {
    deep: true,
  },
);

const addItem = (k: any[]) => {
  if (!k) {
    k = [];
  }
  k.push('');
  return k;
};

const removeItem = (v: any[], index: number) => {
  v.splice(index, 1);
};

onBeforeMount(async () => {
  await store.diceConfigGet();
  const val = cloneDeep(store.curDice.config);
  if (!val.diceMasters || val.diceMasters.length === 0) {
    val.diceMasters = [''];
  }
  if (!val.commandPrefix || val.commandPrefix.length === 0) {
    val.commandPrefix = [''];
  }
  config.value = val;
  nextTick(() => {
    modified.value = false;
  });
});

const submitGiveup = async () => {
  config.value = cloneDeep(store.curDice.config);
  modified.value = false;
  nextTick(() => {
    modified.value = false;
  });
};

onBeforeUnmount(() => {
  // clearInterval(timerId)
});

const submit = async () => {
  const mod = objDiff(config.value, store.curDice.config);
  if (mod.uiPassword) {
    mod.uiPassword = await passwordHash(store.salt, mod.uiPassword);
  }

  if (mod.diceMasters) {
    mod.diceMasters = cloneDeep(config.value.diceMasters);
  }

  if (mod.commandPrefix) {
    mod.commandPrefix = cloneDeep(config.value.commandPrefix);
  }

  if (mod.noticeIds) {
    mod.noticeIds = cloneDeep(config.value.noticeIds);
  }

  if (mod.extDefaultSettings) {
    mod.extDefaultSettings = cloneDeep(config.value.extDefaultSettings);
  }

  if (mod.logSizeNoticeCount) {
    mod.logSizeNoticeCount = toNumber(mod.logSizeNoticeCount);
  }

  await store.diceConfigSet(mod);
  submitGiveup();
};

interface LinkItem {
  label: string;
  value: string;
}

// const state = ref('');
const links = ref<LinkItem[]>([
  { label: '0.0.0.0:3211', value: '0.0.0.0:3211' },
  { label: '127.0.0.1:3211', value: '127.0.0.1:3211' },
]);

const nameWrapUncheck = (v: boolean) => {
  if (!v) {
    // 取消
    dialog.warning({
      title: '取消<玩家名>外框',
      content:
        '不推荐：用户可能会改名为 .bot/.dismiss 等指令，并利用骰点播报让群内其他骰子刷屏，确定要取消吗？',
      positiveText: '确定',
      negativeText: '取消',
      closable: false,
      onNegativeClick() {
        config.value.playerNameWrapEnable = true;
      },
      onMaskClick() {
        config.value.playerNameWrapEnable = true;
      },
    });
  }
};

const mailTest = async () => {
  const res = await postMailTest();
  if (res.result) {
    message.success('已尝试发送测试邮件');
  } else {
    message.error('发送测试邮件失败！' + res.err);
  }
};
</script>

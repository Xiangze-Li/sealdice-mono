<template>
  <header class="page-header">
    <n-switch v-model:value="jsEnable">
      <template #checked>启用</template>
      <template #unchecked>关闭</template>
    </n-switch>
    <n-button v-show="jsEnable" type="primary" @click="jsReload">
      <template #icon>
        <n-icon><i-carbon-renew /></n-icon>
      </template>
      重载 JS
    </n-button>
  </header>

  <n-affix v-if="needReload" :top="60">
    <tip-box type="error">
      <n-text type="error" class="text-base" tag="strong">存在修改，需要重载后生效！</n-text>
    </tip-box>
  </n-affix>
  <n-affix v-if="jsConfigEdited" :top="70">
    <tip-box type="error">
      <n-text type="error" tag="strong" class="text-base">配置内容已修改，不要忘记保存！</n-text>
      <n-button type="info" secondary :disabled="!jsConfigEdited" @click="doJsConfigSave()">
        <template #icon>
          <n-icon><i-carbon-save /></n-icon>
        </template>
        点我保存
      </n-button>
    </tip-box>
  </n-affix>

  <n-row>
    <n-col :span="24">
      <n-tabs v-model:value="mode" class="demo-tabs" justify-content="space-evenly">
        <n-tab-pane tab="控制台" name="console">
          <div>
            <n-scrollbar>
              <code-mirror v-model="code" basic wrap :lang="javascript()" :dark="isDark" />
            </n-scrollbar>

            <div>
              <div style="margin-top: 1rem">
                <n-button type="info" secondary :disabled="!jsEnable" @click="doExecute">
                  <template #icon>
                    <n-icon><i-carbon-play /></n-icon>
                  </template>
                  执行代码
                </n-button>
              </div>
            </div>
            <n-text type="error" tag="p" style="padding: 1rem 0">
              注意：延迟执行的代码，其输出不会立即出现
            </n-text>
            <n-text code style="word-break: break-all; margin-bottom: 1rem; white-space: pre-line">
              <p :key="i" v-for="i in jsLines">{{ i }}</p>
            </n-text>
          </div>
        </n-tab-pane>

        <n-tab-pane tab="插件列表" name="list">
          <header class="mb-4">
            <div class="flex flex-wrap-reverse items-center justify-between gap-y-2">
              <n-flex align="center">
                <span>
                  <n-input v-model:value="jsFilter" size="small" placeholder="" clearable>
                    <template #prefix>
                      <n-icon><i-carbon-search /></n-icon>
                    </template>
                  </n-input>
                </span>
                <n-button
                  type="info"
                  size="tiny"
                  text
                  tag="a"
                  target="_blank"
                  style="text-decoration: none"
                  href="https://github.com/sealdice/javascript">
                  <template #icon>
                    <n-icon><i-carbon-link /></n-icon>
                  </template>
                  获取插件
                </n-button>
              </n-flex>
              <span class="ml-auto">
                <n-upload
                  action=""
                  multiple
                  accept="application/javascript,application/typescript,.js,.ts"
                  :show-file-list="false"
                  @before-upload="beforeUpload"
                  :file-list="uploadFileList">
                  <n-button type="info" secondary>
                    <template #icon>
                      <n-icon><i-carbon-upload /></n-icon>
                    </template>
                    上传插件
                  </n-button>
                </n-upload>
              </span>
            </div>
            <aside class="my-4">
              <n-text v-if="jsFilterCount > 0" type="info" class="text-xs">
                已过滤 {{ jsFilterCount }} 条
              </n-text>
            </aside>
          </header>

          <main class="js-list-main">
            <foldable-card
              v-for="(i, index) of filteredJsList"
              :key="index"
              class="js-item"
              :err-title="i.filename"
              :err-text="i.errText">
              <template #title>
                <el-space class="js-item-header">
                  <el-switch
                    v-model="i.enable"
                    :disabled="i.errText !== ''"
                    style="
                      --el-switch-on-color: var(--el-color-success);
                      --el-switch-off-color: var(--el-color-danger);
                    "
                    @change="changejsScriptStatus(i.name, i.enable)" />
                  <el-text size="large" tag="b">{{ i.name }}</el-text>
                  <el-text>{{ i.version || '&lt;未定义>' }}</el-text>
                  <el-tag v-if="i.official" size="small" type="success">官方</el-tag>

                  <el-tooltip content="该插件使用 TypeScript 编写">
                    <el-tag
                      v-if="i.filename.toLowerCase().endsWith('.ts')"
                      size="small"
                      type="primary"
                      >TS</el-tag
                    >
                  </el-tooltip>
                </el-space>
              </template>

              <template #title-extra>
                <el-button
                  v-if="i.official && i.updateUrls && i.updateUrls.length > 0"
                  type="success"
                  size="small"
                  plain
                  :loading="diffLoading">
                  <template #icon>
                    <i-carbon-download />
                  </template>
                  更新
                </el-button>
                <el-popconfirm
                  v-else-if="i.updateUrls && i.updateUrls.length > 0"
                  confirm-button-text="确认"
                  cancel-button-text="取消"
                  title="更新地址由插件作者提供，是否确认要检查该插件更新？"
                  @confirm="doCheckUpdate(i)">
                  <template #reference>
                    <el-button type="success" size="small" plain :loading="diffLoading">
                      <template #icon>
                        <i-carbon-download />
                      </template>
                      更新
                    </el-button>
                  </template>
                </el-popconfirm>
                <!--                    <el-button :icon="Setting" type="primary" size="small" plain @click="showSettingDialog = true">设置</el-button>-->
                <el-button
                  v-if="i.builtin && i.builtinUpdated"
                  type="danger"
                  size="small"
                  plain
                  @click="doDelete(i)">
                  <template #icon>
                    <i-carbon-row-delete />
                  </template>
                  卸载更新
                </el-button>
                <el-button v-if="!i.builtin" type="danger" size="small" plain @click="doDelete(i)"
                  ><template #icon> <i-carbon-row-delete /> </template>删除</el-button
                >
              </template>

              <template #title-extra-error>
                <el-space>
                  <el-button
                    v-if="i.builtin && i.builtinUpdated"
                    type="danger"
                    size="small"
                    plain
                    @click="doDelete(i)"
                    ><template #icon> <i-carbon-row-delete /> </template>卸载更新</el-button
                  >
                  <el-button
                    v-else-if="!i.builtin"
                    type="danger"
                    size="small"
                    plain
                    @click="doDelete(i)"
                    ><template #icon> <i-carbon-row-delete /> </template>删除</el-button
                  >
                </el-space>
              </template>

              <n-descriptions content-class="whitespace-pre-line">
                <n-descriptions-item v-if="!i.official" :span="3" label="作者">
                  {{ i.author || '&lt;佚名>' }}
                </n-descriptions-item>
                <n-descriptions-item :span="3" label="介绍">
                  {{ i.desc || '&lt;暂无>' }}
                </n-descriptions-item>
                <n-descriptions-item v-if="!i.official" :span="3" label="主页">
                  {{ i.homepage || '&lt;暂无>' }}
                </n-descriptions-item>
                <n-descriptions-item label="许可协议">
                  {{ i.license || '&lt;暂无>' }}
                </n-descriptions-item>
                <n-descriptions-item label="安装时间">
                  {{ dayjs.unix(i.installTime).fromNow() }}
                </n-descriptions-item>
                <n-descriptions-item label="更新时间">
                  {{ dayjs.unix(i.updateTime).fromNow() || '&lt;暂无>' }}
                </n-descriptions-item>
              </n-descriptions>

              <template #unfolded-extra>
                <el-text truncated>{{ i.desc || '&lt;暂无>' }}</el-text>
              </template>
            </foldable-card>

            <el-dialog v-model="showDiff" title="插件内容对比" class="diff-dialog">
              <diff-viewer lang="javascript" :old="jsCheck.old" :new="jsCheck.new" />
              <template #footer>
                <el-space wrap>
                  <el-button @click="showDiff = false">取消</el-button>
                  <el-button v-if="!(jsCheck.old === jsCheck.new)" type="success" @click="jsUpdate">
                    <template #icon>
                      <i-carbon-save />
                    </template>
                    确认更新
                  </el-button>
                </el-space>
              </template>
            </el-dialog>
          </main>
        </n-tab-pane>

        <n-tab-pane label="插件设置" name="config">
          <main>
            <div v-if="size(jsConfig) === 0" style="display: flex; justify-content: center">
              <el-text size="large" tag="strong">暂无设置项</el-text>
            </div>
            <el-collapse v-else class="js-list-main" style="margin-top: 0.5rem">
              <el-collapse-item v-for="(config, i) in jsConfig" :key="i" class="js-item">
                <template #title>
                  <div class="js-item-header">
                    <el-space>
                      <el-text size="large" tag="strong" style="margin-left: 1rem">{{
                        (config as unknown as JsPluginConfig)['pluginName']
                      }}</el-text>
                    </el-space>
                    <el-space>
                      <template
                        v-if="getDeprecatedKeys(config as unknown as JsPluginConfig).length > 0">
                        <el-tooltip
                          content="移除 - 当前插件存在暂未使用的配置项，<br />点击以移除此插件全部暂未使用的配置项。"
                          raw-content
                          placement="bottom-end">
                          <el-icon
                            style="float: right; margin-left: 1rem"
                            @click="
                              doDeleteUnusedConfigs(
                                (config as unknown as JsPluginConfig)['pluginName'],
                                getDeprecatedKeys(config as unknown as JsPluginConfig),
                              )
                            ">
                            <i-carbon-row-delete />
                          </el-icon>
                        </el-tooltip>
                      </template>
                    </el-space>
                  </div>
                </template>
                <el-card shadow="never" style="border: 0">
                  <el-form
                    v-for="(c, index) in (config as unknown as JsPluginConfig)['configs']"
                    :key="index">
                    <template #header>
                      <div class="js-item-header">
                        <el-space>
                          <el-text size="large">{{
                            (c as unknown as JsPluginConfigItem).key
                          }}</el-text>
                        </el-space>
                      </div>
                    </template>
                    <el-form-item
                      v-if="(c as unknown as JsPluginConfigItem).type == 'string'"
                      style="width: 100%; margin-bottom: 0.5rem">
                      <el-form-item label="字符串配置项:">{{
                        (c as unknown as JsPluginConfigItem).key
                      }}</el-form-item
                      ><br />
                      <div style="width: 100%">
                        <el-text>{{ (c as unknown as JsPluginConfigItem).description }}</el-text>
                      </div>
                      <div style="width: 100%; margin-bottom: 0.5rem">
                        <el-input
                          v-model="(c as unknown as JsPluginConfigItem).value"
                          type="textarea"
                          @change="doJsConfigChanged()"></el-input>
                      </div>
                      <template
                        v-if="
                          (c as unknown as JsPluginConfigItem).value !==
                          (c as unknown as JsPluginConfigItem).defaultValue
                        ">
                        <el-tooltip content="重置为初始值" placement="bottom-end">
                          <el-icon
                            style="float: right; margin-left: 1rem"
                            @click="
                              doResetJsConfig(
                                (config as unknown as JsPluginConfig)['pluginName'],
                                (c as unknown as JsPluginConfigItem).key,
                              )
                            ">
                            <i-carbon-paint-brush />
                          </el-icon>
                        </el-tooltip>
                      </template>
                      <template v-if="(c as unknown as JsPluginConfigItem).deprecated">
                        <el-tooltip
                          content="移除 - 这个配置在新版的默认配置中不被使用，<br />但升级而来时仍可能被使用，请确认无用后删除"
                          raw-content
                          placement="bottom-end">
                          <el-icon
                            style="float: right; margin-left: 1rem"
                            @click="
                              doDeleteUnusedConfig(
                                (config as unknown as JsPluginConfig)['pluginName'],
                                (c as unknown as JsPluginConfigItem).key,
                                false,
                              )
                            ">
                            <i-carbon-row-delete />
                          </el-icon>
                        </el-tooltip>
                      </template>
                    </el-form-item>
                    <el-form-item v-if="(c as unknown as JsPluginConfigItem).type == 'int'">
                      <el-form-item label="整数配置项:">{{
                        (c as unknown as JsPluginConfigItem).key
                      }}</el-form-item
                      ><br />
                      <div style="width: 100%">
                        <el-text>{{ (c as unknown as JsPluginConfigItem).description }}</el-text>
                      </div>
                      <el-form-item :span="30">
                        <div style="margin-left: 1rem">
                          <el-input-number
                            v-model="(c as unknown as JsPluginConfigItem).value"
                            type="number"
                            @change="doJsConfigChanged()"></el-input-number>
                        </div>
                      </el-form-item>
                      <template
                        v-if="
                          (c as unknown as JsPluginConfigItem).value !==
                          (c as unknown as JsPluginConfigItem).defaultValue
                        ">
                        <el-tooltip content="重置为初始值" placement="bottom-end">
                          <el-icon
                            style="float: right; margin-left: 1rem"
                            @click="
                              doResetJsConfig(
                                (config as unknown as JsPluginConfig)['pluginName'],
                                (c as unknown as JsPluginConfigItem).key,
                              )
                            ">
                            <i-carbon-paint-brush />
                          </el-icon>
                        </el-tooltip>
                      </template>
                      <template v-if="(c as unknown as JsPluginConfigItem).deprecated">
                        <el-tooltip
                          content="移除 - 这个配置在新版的默认配置中不被使用，<br />但升级而来时仍可能被使用，请确认无用后删除"
                          raw-content
                          placement="bottom-end">
                          <el-icon
                            style="float: right; margin-left: 1rem"
                            @click="
                              doDeleteUnusedConfig(
                                (config as unknown as JsPluginConfig)['pluginName'],
                                (c as unknown as JsPluginConfigItem).key,
                                false,
                              )
                            ">
                            <i-carbon-row-delete />
                          </el-icon>
                        </el-tooltip>
                      </template>
                    </el-form-item>
                    <el-form-item v-if="(c as unknown as JsPluginConfigItem).type == 'float'">
                      <el-form-item label="浮点数配置项:">{{
                        (c as unknown as JsPluginConfigItem).key
                      }}</el-form-item
                      ><br />
                      <div style="width: 100%">
                        <el-text>{{ (c as unknown as JsPluginConfigItem).description }}</el-text>
                      </div>
                      <el-form-item :span="30">
                        <div style="margin-left: 1rem">
                          <el-input-number
                            v-model="(c as unknown as JsPluginConfigItem).value"
                            type="number"
                            @change="doJsConfigChanged()"></el-input-number>
                        </div>
                      </el-form-item>
                      <template
                        v-if="
                          (c as unknown as JsPluginConfigItem).value !==
                          (c as unknown as JsPluginConfigItem).defaultValue
                        ">
                        <el-tooltip content="重置为初始值" placement="bottom-end">
                          <el-icon
                            style="float: right; margin-left: 1rem"
                            @click="
                              doResetJsConfig(
                                (config as unknown as JsPluginConfig)['pluginName'],
                                (c as unknown as JsPluginConfigItem).key,
                              )
                            ">
                            <i-carbon-paint-brush />
                          </el-icon>
                        </el-tooltip>
                      </template>
                      <template v-if="(c as unknown as JsPluginConfigItem).deprecated">
                        <el-tooltip
                          content="移除 - 这个配置在新版的默认配置中不被使用，<br />但升级而来时仍可能被使用，请确认无用后删除"
                          raw-content
                          placement="bottom-end">
                          <el-icon
                            style="float: right; margin-left: 1rem"
                            @click="
                              doDeleteUnusedConfig(
                                (config as unknown as JsPluginConfig)['pluginName'],
                                (c as unknown as JsPluginConfigItem).key,
                                false,
                              )
                            ">
                            <i-carbon-row-delete />
                          </el-icon>
                        </el-tooltip>
                      </template>
                    </el-form-item>
                    <el-form-item v-if="(c as unknown as JsPluginConfigItem).type == 'bool'">
                      <el-form-item label="布尔配置项:">{{
                        (c as unknown as JsPluginConfigItem).key
                      }}</el-form-item
                      ><br />
                      <div style="width: 100%">
                        <el-text>{{ (c as unknown as JsPluginConfigItem).description }}</el-text>
                      </div>
                      <el-form-item :span="30">
                        <div style="margin-left: 1rem">
                          <el-switch
                            v-model="(c as unknown as JsPluginConfigItem).value"
                            @change="doJsConfigChanged()"></el-switch>
                        </div>
                      </el-form-item>
                      <template
                        v-if="
                          (c as unknown as JsPluginConfigItem).value !==
                          (c as unknown as JsPluginConfigItem).defaultValue
                        ">
                        <el-tooltip content="重置为初始值" placement="bottom-end">
                          <el-icon
                            style="float: right; margin-left: 1rem"
                            @click="
                              doResetJsConfig(
                                (config as unknown as JsPluginConfig)['pluginName'],
                                (c as unknown as JsPluginConfigItem).key,
                              )
                            ">
                            <i-carbon-paint-brush />
                          </el-icon>
                        </el-tooltip>
                      </template>
                      <template v-if="(c as unknown as JsPluginConfigItem).deprecated">
                        <el-tooltip
                          content="移除 - 这个配置在新版的默认配置中不被使用，<br />但升级而来时仍可能被使用，请确认无用后删除"
                          raw-content
                          placement="bottom-end">
                          <el-icon
                            style="float: right; margin-left: 1rem"
                            @click="
                              doDeleteUnusedConfig(
                                (config as unknown as JsPluginConfig)['pluginName'],
                                (c as unknown as JsPluginConfigItem).key,
                                false,
                              )
                            ">
                            <i-carbon-row-delete />
                          </el-icon>
                        </el-tooltip>
                      </template>
                    </el-form-item>
                    <el-form-item
                      v-if="(c as unknown as JsPluginConfigItem).type == 'template'"
                      style="width: 100%; margin-bottom: 0.5rem">
                      <el-form-item
                        label="模板配置项:"
                        style="width: 100%; margin-bottom: 0.5rem"
                        >{{ (c as unknown as JsPluginConfigItem).key }}</el-form-item
                      ><br />
                      <div style="width: 100%">
                        <el-text>{{ (c as unknown as JsPluginConfigItem).description }}</el-text>
                      </div>
                      <template
                        v-if="
                          !isEqual(
                            (c as unknown as JsPluginConfigItem).value,
                            (c as unknown as JsPluginConfigItem).defaultValue,
                          )
                        ">
                        <el-tooltip content="重置为初始值" placement="bottom-end">
                          <el-icon
                            style="float: right; margin-left: 1rem"
                            @click="
                              doResetJsConfig(
                                (config as unknown as JsPluginConfig)['pluginName'],
                                (c as unknown as JsPluginConfigItem).key,
                              )
                            ">
                            <i-carbon-paint-brush />
                          </el-icon>
                        </el-tooltip>
                      </template>
                      <template v-if="(c as unknown as JsPluginConfigItem).deprecated">
                        <el-tooltip
                          content="移除 - 这个配置在新版的默认配置中不被使用，<br />但升级而来时仍可能被使用，请确认无用后删除"
                          raw-content
                          placement="bottom-end">
                          <el-icon
                            style="float: right; margin-left: 1rem"
                            @click="
                              doDeleteUnusedConfig(
                                (config as unknown as JsPluginConfig)['pluginName'],
                                (c as unknown as JsPluginConfigItem).key,
                                false,
                              )
                            ">
                            <i-carbon-row-delete />
                          </el-icon>
                        </el-tooltip>
                      </template>
                      <el-form-item style="width: 100%; margin-bottom: 0.5rem">
                        <div
                          v-for="(d, index) in (c as unknown as JsPluginConfigItem).value"
                          :key="index"
                          style="width: 100%; margin-bottom: 0.5rem">
                          <!-- 这里面是单条修改项 -->
                          <el-row>
                            <el-col style="width: 100%; margin-bottom: 0.5rem">
                              <span style="width: 100%">
                                <el-input
                                  v-model="(c as unknown as JsPluginConfigItem).value[index]"
                                  type="textarea"
                                  :autosize="true"
                                  @change="doJsConfigChanged()"></el-input>
                              </span>
                            </el-col>
                            <el-col :span="5">
                              <div
                                style="
                                  display: flex;
                                  align-items: center;
                                  width: 1.3rem;
                                  margin-left: 1rem;
                                  margin-top: 0.5rem;
                                ">
                                <el-tooltip
                                  :content="
                                    index === 0 ? '点击添加一项' : '点击删除你不想要的配置项'
                                  "
                                  placement="bottom-start">
                                  <el-icon>
                                    <i-carbon-add-filled
                                      v-if="index == 0"
                                      @click="
                                        doJsConfigAddItem(
                                          (c as unknown as JsPluginConfigItem).value,
                                        )
                                      " />
                                    <i-carbon-close-outline
                                      v-else
                                      @click="
                                        doJsConfigRemoveItemAt(
                                          (c as unknown as JsPluginConfigItem).value,
                                          index,
                                        )
                                      " />
                                  </el-icon>
                                </el-tooltip>
                              </div>
                            </el-col>
                          </el-row>
                        </div>
                      </el-form-item>
                    </el-form-item>
                    <el-form-item v-if="(c as unknown as JsPluginConfigItem).type == 'option'">
                      <el-form-item
                        label="选项配置项:"
                        style="width: 100%; margin-bottom: 0.5rem"
                        >{{ (c as unknown as JsPluginConfigItem).key }}</el-form-item
                      >
                      <div style="width: 100%">
                        <el-text>{{ (c as unknown as JsPluginConfigItem).description }}</el-text>
                      </div>
                      <template
                        v-if="
                          (c as unknown as JsPluginConfigItem).value !==
                          (c as unknown as JsPluginConfigItem).defaultValue
                        ">
                        <el-tooltip content="重置为初始值" placement="bottom-end">
                          <el-icon
                            style="float: right; margin-left: 1rem"
                            @click="
                              doResetJsConfig(
                                (config as unknown as JsPluginConfig)['pluginName'],
                                (c as unknown as JsPluginConfigItem).key,
                              )
                            ">
                            <i-carbon-paint-brush />
                          </el-icon>
                        </el-tooltip>
                      </template>
                      <template v-if="(c as unknown as JsPluginConfigItem).deprecated">
                        <el-tooltip
                          content="移除 - 这个配置在新版的默认配置中不被使用，<br />但升级而来时仍可能被使用，请确认无用后删除"
                          raw-content
                          placement="bottom-end">
                          <el-icon
                            style="float: right; margin-left: 1rem"
                            @click="
                              doDeleteUnusedConfig(
                                (config as unknown as JsPluginConfig)['pluginName'],
                                (c as unknown as JsPluginConfigItem).key,
                                false,
                              )
                            ">
                            <i-carbon-row-delete />
                          </el-icon>
                        </el-tooltip>
                      </template>
                      <div style="width: 100%; margin-bottom: 0.5rem">
                        <el-select
                          v-model="(c as unknown as JsPluginConfigItem).value"
                          @change="doJsConfigChanged()">
                          <el-option
                            v-for="s in (c as unknown as JsPluginConfigItem).option"
                            :key="s"
                            :value="s"
                            >{{ s }}</el-option
                          >
                        </el-select>
                      </div>
                    </el-form-item>
                    <el-form-item
                      v-if="(c as unknown as JsPluginConfigItem).type == 'task:cron'"
                      style="width: 100%; margin-bottom: 0.5rem">
                      <el-form-item label="Cron 型定时任务:">{{
                        (c as unknown as JsPluginConfigItem).key
                      }}</el-form-item
                      ><br />
                      <div style="width: 100%">
                        <el-text>{{ (c as unknown as JsPluginConfigItem).description }}</el-text>
                      </div>
                      <div style="width: 100%; margin-bottom: 0.5rem">
                        <el-input
                          v-model="(c as unknown as JsPluginConfigItem).value"
                          type="textarea"
                          @change="
                            doTaskCronFormatCheck(
                              (config as unknown as JsPluginConfig)['pluginName'],
                              (c as unknown as JsPluginConfigItem).key,
                              (c as unknown as JsPluginConfigItem).value,
                            )
                          "></el-input>
                        <el-text
                          v-if="
                            jsConfigFormatErrKeys.indexOf(
                              (config as unknown as JsPluginConfig)['pluginName'] +
                                '/' +
                                (c as unknown as JsPluginConfigItem).key,
                            ) !== -1
                          "
                          type="danger">
                          格式错误！
                        </el-text>
                      </div>
                      <template
                        v-if="
                          (c as unknown as JsPluginConfigItem).value !==
                          (c as unknown as JsPluginConfigItem).defaultValue
                        ">
                        <el-tooltip content="重置为初始值" placement="bottom-end">
                          <el-icon
                            style="float: right; margin-left: 1rem"
                            @click="
                              doResetJsConfig(
                                (config as unknown as JsPluginConfig)['pluginName'],
                                (c as unknown as JsPluginConfigItem).key,
                              )
                            ">
                            <i-carbon-paint-brush />
                          </el-icon>
                        </el-tooltip>
                      </template>
                      <template v-if="(c as unknown as JsPluginConfigItem).deprecated">
                        <el-tooltip
                          content="移除 - 这个定时任务在当前版本中不被使用，<br />但未来版本中仍可能被使用，请确认无用后删除"
                          raw-content
                          placement="bottom-end">
                          <el-icon
                            style="float: right; margin-left: 1rem"
                            @click="
                              doDeleteUnusedConfig(
                                (config as unknown as JsPluginConfig)['pluginName'],
                                (c as unknown as JsPluginConfigItem).key,
                                true,
                              )
                            ">
                            <i-carbon-row-delete />
                          </el-icon>
                        </el-tooltip>
                      </template>
                    </el-form-item>
                    <el-form-item
                      v-if="(c as unknown as JsPluginConfigItem).type == 'task:daily'"
                      style="width: 100%; margin-bottom: 0.5rem">
                      <el-form-item label="每日定时任务:">{{
                        (c as unknown as JsPluginConfigItem).key
                      }}</el-form-item
                      ><br />
                      <div style="width: 100%">
                        <el-text>{{ (c as unknown as JsPluginConfigItem).description }}</el-text>
                      </div>
                      <div style="width: 100%; margin-bottom: 0.5rem">
                        <el-input
                          v-model="(c as unknown as JsPluginConfigItem).value"
                          type="textarea"
                          @change="
                            doTaskDailyFormatCheck(
                              (config as unknown as JsPluginConfig)['pluginName'],
                              (c as unknown as JsPluginConfigItem).key,
                              (c as unknown as JsPluginConfigItem).value,
                            )
                          "></el-input>
                        <el-text
                          v-if="
                            jsConfigFormatErrKeys.indexOf(
                              (config as unknown as JsPluginConfig)['pluginName'] +
                                '/' +
                                (c as unknown as JsPluginConfigItem).key,
                            ) !== -1
                          "
                          type="danger">
                          格式错误！
                        </el-text>
                      </div>
                      <template
                        v-if="
                          (c as unknown as JsPluginConfigItem).value !==
                          (c as unknown as JsPluginConfigItem).defaultValue
                        ">
                        <el-tooltip content="重置为初始值" placement="bottom-end">
                          <el-icon
                            style="float: right; margin-left: 1rem"
                            @click="
                              doResetJsConfig(
                                (config as unknown as JsPluginConfig)['pluginName'],
                                (c as unknown as JsPluginConfigItem).key,
                              )
                            ">
                            <i-carbon-paint-brush />
                          </el-icon>
                        </el-tooltip>
                      </template>
                      <template v-if="(c as unknown as JsPluginConfigItem).deprecated">
                        <el-tooltip
                          content="移除 - 这个定时任务在当前版本中不被使用，<br />但未来版本中仍可能被使用，请确认无用后删除"
                          raw-content
                          placement="bottom-end">
                          <el-icon
                            style="float: right; margin-left: 1rem"
                            @click="
                              doDeleteUnusedConfig(
                                (config as unknown as JsPluginConfig)['pluginName'],
                                (c as unknown as JsPluginConfigItem).key,
                                true,
                              )
                            ">
                            <i-carbon-row-delete />
                          </el-icon>
                        </el-tooltip>
                      </template>
                    </el-form-item>
                  </el-form>
                </el-card>
              </el-collapse-item>
            </el-collapse>
          </main>
        </n-tab-pane>
      </n-tabs>
    </n-col>
  </n-row>
</template>

<script lang="ts" setup>
import * as dayjs from 'dayjs';
import { basicSetup, EditorView } from 'codemirror';
import { javascript } from '@codemirror/lang-javascript';
import CodeMirror from 'vue-codemirror6';
import { isEqual, size } from 'es-toolkit/compat';
import type { JsPluginConfig, JsPluginConfigItem, JsScriptInfo } from '#';
import { postUtilsCheckCronExpr } from '~/api/v1/utils';
import {
  checkJsUpdate,
  deleteJs,
  deleteUnusedJsConfigs,
  disableJS,
  enableJS,
  executeJS,
  getJsConfigs,
  getJsList,
  getJsRecord,
  getJsStatus,
  reloadJS,
  resetJsConfig,
  setJsConfigs,
  shutDownJS,
  updateJs,
  uploadJs,
} from '~/api/v1/js';
import type { UploadRawFile } from 'element-plus';
import { type UploadFileInfo, useDialog, useMessage } from 'naive-ui';

const message = useMessage();
const dialog = useDialog();
const isDark = useDark({ disableTransition: false });

const jsEnable = ref(false);
const editorBox = ref(null);
const mode = ref('console');
const needReload = ref(false);
let editor: EditorView;

const jsLines = ref([] as string[]);

const defaultText = [
  '// 学习制作可以看这里：https://github.com/sealdice/javascript/tree/main/examples',
  '// 下载插件可以看这里：https://github.com/sealdice/javascript/tree/main/scripts',
  '// 使用 TypeScript，编写更容易 https://github.com/sealdice/javascript/tree/main/examples_ts',
  '// 目前可用于：创建自定义指令，自定义 COC 房规，发送网络请求，读写本地数据',
  '',
  "console.log('这是测试控制台');",
  "console.log('可以这样来查看变量详情：');",
  'console.log(Object.keys(seal));',
  "console.log('更多内容正在制作中...')",
  "console.log('注意：测试版！API 仍然可能发生重大变化！')",
  '// 写在这里的所有变量都是临时变量，如果你希望全局变量，使用 globalThis',
  '// 但是注意，全局变量在进程关闭后失效，想保存状态请存入硬盘。',
  'globalThis._test = 123;',
  '',
  "let ext = seal.ext.find('test');",
  'if (!ext) {',
  "  ext = seal.ext.new('test', '木落', '1.0.0');",
  '  seal.ext.register(ext);',
  '}',
];
const code = ref<string>(defaultText.join('\n'));

/** 执行指令 */
const doExecute = async () => {
  jsLines.value = [];
  const data = await executeJS(code.value);

  // 优先填充 print 输出
  const lines = [] as string[];
  if (data.outputs) {
    lines.push(...data.outputs);
  }
  // 填充 err 或 ret
  if (data.err) {
    lines.push(data.err);
  } else {
    lines.push(data.ret as string);
    try {
      (window as any).lastJSValue = data.ret;
      (globalThis as any).lastJSValue = data.ret;
      // eslint-disable-next-line @typescript-eslint/no-unused-vars
    } catch (e) {}
  }
  jsLines.value = lines;
};

const jsConfigEdited = ref(false);
const doJsConfigChanged = () => {
  jsConfigEdited.value = true;
};

const jsConfigFormatErrKeys: Ref<string[]> = ref([]);
const doTaskCronFormatCheck = async (pluginName: string, key: string, expr: string) => {
  const index = jsConfigFormatErrKeys.value.indexOf(pluginName + '/' + key);
  try {
    await postUtilsCheckCronExpr(expr);
    if (index !== -1) {
      jsConfigFormatErrKeys.value.splice(index, 1);
    }
    jsConfigEdited.value = true;
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
  } catch (_err) {
    if (index === -1) {
      jsConfigFormatErrKeys.value.push(pluginName + '/' + key);
    }
    jsConfigEdited.value = true;
  }
};

const doTaskDailyFormatCheck = (pluginName: string, key: string, expr: string) => {
  const pattern = /^([0-1]?[0-9]|2[0-3]):([0-5][0-9])$/;
  const index = jsConfigFormatErrKeys.value.indexOf(pluginName + '/' + key);
  if (pattern.test(expr)) {
    if (index !== -1) {
      jsConfigFormatErrKeys.value.splice(index, 1);
    }
    jsConfigEdited.value = true;
  } else {
    if (index === -1) {
      jsConfigFormatErrKeys.value.push(pluginName + '/' + key);
    }
    jsConfigEdited.value = true;
  }
};

const doDeleteUnusedConfig = (pluginName: any, key: any, isTask: boolean) => {
  dialog.warning({
    title: '确认删除',
    content: isTask
      ? `确认删除插件 ${pluginName} 的定时任务 ${key}？`
      : `确认删除插件 ${pluginName} 的配置项 ${key}？`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      await deleteUnusedJsConfigs(pluginName, [key]);
      setTimeout(() => {
        // 稍等等再重载，以免出现没删掉
        refreshConfig();
      }, 1000);
      message.success('配置项已删除');
    },
  });
};

const getDeprecatedKeys = (config: JsPluginConfig): string[] => {
  const configs = config.configs || [];
  const result: string[] = [];
  for (const c of configs) {
    if (c.deprecated) {
      result.push(c.key);
    }
  }
  return result;
};

const doDeleteUnusedConfigs = (pluginName: string, keys: string[]) => {
  dialog.warning({
    title: '确认删除配置项',
    content: `确认删除插件 ${pluginName} 的共 ${keys.length} 个暂未使用的配置项/定时任务？`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      await deleteUnusedJsConfigs(pluginName, keys);
      setTimeout(() => {
        // 稍等等再重载，以免出现没删掉
        refreshConfig();
      }, 1000);
      message.success('配置项已删除');
    },
  });
};

const doResetJsConfig = (pluginName: string, key: string) => {
  dialog.warning({
    title: '确认重置',
    content: '确认重置这条配置项回默认状态？',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      await resetJsConfig(pluginName, key);
      setTimeout(() => {
        refreshConfig();
        jsConfigEdited.value = false;
        jsConfigFormatErrKeys.value = [];
      }, 1000);
      message.success('重置成功！');
    },
  });
};
const doJsConfigAddItem = (arr: any[]) => {
  arr.push('');
  doJsConfigChanged();
  return arr;
};
const doJsConfigRemoveItemAt = <T,>(arr: T[], index: number) => {
  if (index < 0 || index >= arr.length) {
    return arr;
  }
  arr.splice(index, 1);
  doJsConfigChanged();
  return arr;
};

const doJsConfigSave = async () => {
  await setJsConfigs(jsConfig.value);
  jsConfigEdited.value = false;
  message.success('已保存');
};

let timerId: number;

onMounted(async () => {
  jsEnable.value = await jsStatus();
  watch(jsEnable, async (newStatus, oldStatus) => {
    console.log('new:', newStatus, ' old:', oldStatus);
    if (oldStatus !== undefined) {
      if (newStatus) {
        console.log('reload');
        await jsReload();
      } else {
        console.log('shutdown');
        await jsShutdown();
      }
    }
  });

  const el = editorBox.value as any as HTMLElement;
  editor = new EditorView({
    extensions: [basicSetup, javascript(), EditorView.lineWrapping],
    parent: el,
    doc: defaultText.join('\n'),
  });
  try {
    (globalThis as any).editor = editor;
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
  } catch (e) {}

  await refreshList();
  if (jsList.value.length > 0) {
    mode.value = 'list';
  }
  await refreshConfig();

  timerId = setInterval(async () => {
    console.log('refresh');
    const data = await getJsRecord();

    if (data.outputs) {
      jsLines.value.push(...data.outputs);
    }
  }, 3000);
});

onBeforeUnmount(() => {
  clearInterval(timerId);
});

const jsList = ref<JsScriptInfo[]>([]);
const jsFilter = ref<string>('');
const jsFilterCount = computed(() => jsList.value.length - filteredJsList.value.length);
const filteredJsList = computed(() =>
  jsList.value.filter(js => {
    if (jsFilter.value === '') {
      return true;
    }
    const val = jsFilter.value.toLowerCase();
    return (
      js.name?.toLowerCase()?.includes(val) ||
      js.desc?.toLowerCase()?.includes(val) ||
      js.author?.toLowerCase()?.includes(val)
    );
  }),
);
const jsConfig = ref<{ [key: string]: JsPluginConfig }>({});
const uploadFileList = ref<any[]>([]);

// const jsVisitDir = async () => {
// 好像 webui 上没啥效果，先算了
// await store.jsVisitDir();
// };

const jsStatus = async () => {
  const res = await getJsStatus();
  return res.result ? res.status : false;
};

const refreshList = async () => {
  const lst = await getJsList();
  jsList.value = lst;
};

const refreshConfig = async () => {
  jsConfig.value = await getJsConfigs();
};

const jsReload = async () => {
  const ret = await reloadJS();
  if (ret && ret.testMode) {
    message.success('展示模式无法重载脚本');
  } else {
    message.success('已重载');
    await refreshList();
    needReload.value = false;
  }
  jsEnable.value = await jsStatus();
};

const jsShutdown = async () => {
  const ret = await shutDownJS();
  if (ret?.testMode) {
    message.success('展示模式无法关闭');
  } else if (ret?.result === true) {
    message.success('已关闭 JS 支持');
    jsLines.value = [];
    await refreshList();
  }
  jsEnable.value = await jsStatus();
};

const beforeUpload = async (file: UploadRawFile) => {
  // UploadRawFile
  const fd = new FormData();
  fd.append('file', file);
  await uploadJs(file);
  await refreshList();
  message.success('上传完成，请在全部操作完成后，手动重载插件');
  needReload.value = true;
};

const doDelete = async (data: JsScriptInfo) => {
  dialog.warning({
    title: data.official ? '确认卸载' : '确认删除',
    content: data.official
      ? `确认卸载官方插件「${data.name}」的更新，确定吗？`
      : `确认删除插件「${data.name}」，确定吗？`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      await deleteJs(data.filename);
      setTimeout(() => {
        // 稍等等再重载，以免出现没删掉
        refreshList();
      }, 1000);
      message.success('插件已删除，请手动重载后生效');
      needReload.value = true;
    },
  });
};

const changejsScriptStatus = async (name: string, status: boolean) => {
  if (status) {
    const ret = await enableJS(name);
    setTimeout(() => {
      refreshList();
    }, 1000);
    if (ret.result) {
      message.success('插件已删除，请手动重载后生效');
    }
  } else {
    const ret = await disableJS(name);
    setTimeout(() => {
      refreshList();
    }, 1000);
    if (ret.result) {
      message.success('插件已删除，请手动重载后生效');
    }
  }
  needReload.value = true;
  return true;
};

// const showSettingDialog = ref<boolean>(false);

// interface DeckProp {
//   key: string;
//   value: string;

//   name?: string;
//   desc?: string;
//   required?: boolean;
//   default?: string;
// }

// const settingForm = ref({
//   props: [{ key: 'name', value: 'test props' }] as DeckProp[],
// });

const showDiff = ref<boolean>(false);
const diffLoading = ref<boolean>(false);

interface JsCheckResult {
  old: string;
  new: string;
  filename: string;
  tempFileName: string;
}

const jsCheck = ref<JsCheckResult>({
  old: '',
  new: '',
  filename: '',
  tempFileName: '',
});

const doCheckUpdate = async (data: any) => {
  diffLoading.value = true;
  const checkResult = await checkJsUpdate(data.filename);
  diffLoading.value = false;
  if (checkResult.result) {
    jsCheck.value = { ...checkResult, filename: data.filename };
    showDiff.value = true;
  } else {
    message.error('检查更新失败！' + checkResult.err);
  }
};

const jsUpdate = async () => {
  const res = await updateJs(jsCheck.value.tempFileName, jsCheck.value.filename);
  if (res.result) {
    showDiff.value = false;
    needReload.value = true;
    setTimeout(() => {
      refreshList();
    }, 1000);
    message.success('更新成功，请手动重载后生效');
  } else {
    showDiff.value = false;
    message.error('更新失败！' + res.err);
  }
};
</script>

<style lang="css">
.cm-editor {
  /* height: v-bind("$props.initHeight"); */
  height: 20rem;
  /* font-size: 18px; */

  outline: 0 !important;
  /* height: 50rem; */
  box-shadow:
    0 2px 4px rgba(0, 0, 0, 0.12),
    0 0 6px rgba(0, 0, 0, 0.04);
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

.deck-keys {
  display: flex;
  flex-flow: wrap;

  & > span {
    margin-right: 1rem;
    /* width: fit-content; */
  }
}

.upload {
  > ul {
    display: none;
  }
}

.js-list-header {
  margin-bottom: 1rem;
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
  justify-content: space-between;
}

.js-list-main {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
}

.js-item-header {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
  justify-content: space-between;
}

.js-item {
  min-width: 100%;
}
</style>

<template>
  <div>
    <div class="float-left mr-12 h-[250px]">
      <img class="h-[200px]" :src="imgMascot" alt="Mascot" />
    </div>
    <div class="text-left text-sm leading-8">
      <div class="top-item">
        <span class="font-semibold">当前版本：{{ store.curDice.baseInfo.version }}</span>
      </div>
      <div class="top-item">
        <span class="font-semibold">最新版本：{{ store.curDice.baseInfo.versionNew }}</span>
      </div>
      <div class="top-item mt-4">
        官方网站：
        <n-button type="primary" text tag="a" target="_blank" href="https://github.com/sealdice-ce/sealdice-ce">
          https://github.com/sealdice-ce/sealdice-ce
        </n-button>
      </div>
      <div class="top-item">
        使用手册：
        <n-button
          type="primary"
          text
          tag="a"
          target="_blank"
          href="https://dice.weizaima.com/manual/">
          https://dice.weizaima.com/manual/
        </n-button>
      </div>
      <div class="top-item">
        投喂海豹：
        <n-button
          type="primary"
          text
          tag="a"
          target="_blank"
          href="https://dice.weizaima.com/feed/">
          https://dice.weizaima.com/feed/
        </n-button>
      </div>
      <div class="top-item">
        源码：
        <n-button
          type="primary"
          text
          tag="a"
          target="_blank"
          href="https://github.com/sealdice/sealdice-core/">
          https://github.com/sealdice/sealdice-core
        </n-button>
      </div>
    </div>
  </div>

  <div class="mt-16 text-center text-xl font-bold">感谢 ♪(・ω・)ノ</div>
  <div class="mb-8 mt-2 text-center text-sm">
    特别鸣谢参与测试、反馈问题，帮助完善海豹指令的各位！以下列出感谢名单（排名不分先后）
  </div>

  <n-flex vertical size="large" class="about">
    <div class="font-bold">社区协力</div>

    <div class="font-bold">V1.5 版本</div>
    <n-flex class="developers">
      <github-avatar username="木落" />
      <github-avatar username="PaienNate" />
      <github-avatar username="Fripine" />
      <github-avatar username="只是另一个 ID" />
      <github-avatar username="逐风" />
      <github-avatar username="暮星" />
      <github-avatar username="JohNSoN" />
      <github-avatar username="Bugtower100" />
      <github-avatar username="希望潇洒的风：Ceeling" info="UI" user="charyflys" />
    </n-flex>

    <div class="font-bold">V1.5.0 特别致谢</div>
    <n-flex class="developers">
      <github-avatar username="雪桃" user="LoranaAurelia" info="提供新的魔法阵代理" />
      <github-avatar
        username="[骰]诺扣提"
        src="https://d1.sealdice.com/images/bird.jpg"
        href="https://dice.weizaima.com/public-dice"
        info="公骰 (协助测试)" />
      <github-avatar
        username="山本健一"
        user="kenichiLyon"
        src="https://d1.sealdice.com/images/kenichiLyon.jpg"
        info="协助测试" />
      <github-avatar username="白鱼" user="baiyu-yu" info="协助测试" />
      <github-avatar
        username="[骰]Gaza(3658887052)"
        src="https://firehomework.top/img/6819CBB71F64D7BBFC90DCF4432B76BC.jpg"
        info="协助测试"
        only-name />
    </n-flex>

    <div class="font-bold">V1.4 版本</div>
    <n-flex class="developers">
      <github-avatar username="只是另一个 ID" />
      <github-avatar username="Szzrain" />
      <github-avatar username="JohNSoN" />
      <github-avatar username="檀轶步棋" />
      <github-avatar username="脑" />
      <github-avatar username="木落" />
      <github-avatar username="Director259" />
      <github-avatar username="浣熊旅記" user="VolEurr0Se" info="dnd5文档" />
      <github-avatar username="调零" user="zeroxilo" info="dnd5文档" />
      <github-avatar username="凤吹风雪" user="cherrybird7" info="7版怪锤" />
      <github-avatar username="奈亚猫猫汉化组" info="7版怪锤" only-name />
      <github-avatar username="稚鸢音" info="7版怪锤" only-name />
      <github-avatar username="Fripine" />
      <github-avatar username="逐风" />
      <github-avatar username="炽热" user="yichere" />
      <github-avatar username="暮星" />
      <github-avatar username="0.018" user="0018cha" />
      <github-avatar username="冷筱华" user="shakugannosaints" />
    </n-flex>

    <div class="font-bold">V1.4 安卓端</div>
    <n-flex class="developers">
      <github-avatar username="Szzrain" />
      <github-avatar username="PaienNate" />
      <github-avatar username="木落" />
    </n-flex>

    <div class="font-bold">V1.4.5 特别致谢</div>
    <n-flex class="developers">
      <github-avatar username="Linwenxuan" user="Linwenxuan04" />
    </n-flex>

    <div class="font-bold">V1.3 版本</div>
    <n-flex class="developers">
      <github-avatar username="只是另一个 ID" />
      <github-avatar username="Szzrain" />
      <github-avatar username="JohNSoN" />
      <github-avatar username="檀轶步棋" />
      <github-avatar username="脑" />
      <github-avatar username="木落" />
    </n-flex>

    <div class="font-bold">V1.3 安卓端</div>
    <n-flex class="developers">
      <github-avatar username="Szzrain" />
      <github-avatar username="PaienNate" />
    </n-flex>

    <div class="font-bold">V1.2 版本</div>
    <n-flex class="developers">
      <github-avatar username="木落" />
      <github-avatar username="Szzrain" />
      <github-avatar username="于言诺" user="yuyannuo" />
      <github-avatar username="檀轶步棋" />
      <github-avatar username="脑" />
      <github-avatar username="熊米" user="SunnyJoyce" />
      <github-avatar username="浣熊旅記" user="VolEurr0Se" />
      <github-avatar username="流溪" user="lxy071130" />
      <github-avatar username="病" user="nodisease" />
    </n-flex>

    <div class="font-bold">V1.2 安卓端</div>
    <n-flex class="developers">
      <github-avatar username="木末君" user="96368a" />
      <github-avatar username="Szzrain" />
      <github-avatar username="极夜幻想" user="JiYeHuanXiang" />
    </n-flex>

    <div class="font-bold">V1.1 版本</div>
    <span>Szzrain - 实现了Discord和Kook(开黑啦)两个平台的海豹接入</span>
    <span>
      于言诺 -
      制作了很多海豹扩展和牌堆，如养猫、踢海豹、赛博功德、风味月饼、万圣节糖果等等……协助撰写了一些海豹的文档和教程，并找出了众多海豹的bug
    </span>
    <span>云陌 - 海豹文档教程协力，同时也找了很多海豹的bug</span>
    <span>
      <n-button
        type="primary"
        text
        tag="a"
        target="_blank"
        href="https://github.com/kagangtuya-star">
        星尘
      </n-button>
      - 编写了海豹同网络登录的教程，友情提供了用于指令参考的fvtt。以及一些建议和bug反馈
    </span>

    <div class="font-bold">V1.0版本</div>
    <span>Ariel船长 - 早期测试参与者，协助解决登录流程问题</span>
    <span>Raycel - 早期测试参与者，协助解决登录流程问题</span>
    <span>kuma - 早期测试参与者，海豹的第一次全指令全流程测试</span>
    <span>卟啵 - 早期测试参与者，回报了中文路径和空格路径问题，协助解决了登录流程问题</span>
    <span>蜜瓜包 - 早期测试参与者，默认文档中“怪物之锤查询”的编纂者之一</span>
    <span>
      月森优姬 -
      早期测试参与者，提出了大量各种各样建议和BUG反馈，纠正了一些与规则书不统一的问题，COC同义词和默认技能点数的编纂者
    </span>
    <span>清茶 - 在4月7日的可靠性测试中，参与构造了让旧版海豹进程崩溃的指令</span>
    <span>脑 - 在4月7日的可靠性测试中，参与构造了让旧版海豹进程崩溃的指令</span>
    <span>Greed锦鲤 - 在4月7日的可靠性测试中，参与构造了让旧版海豹进程崩溃的指令</span>
    <span>格莱德 - 在4月7日的可靠性测试中，参与构造了让旧版海豹进程崩溃的指令</span>
    <span>我来逛街 - 提出很多建议；帮助改进了DND5E同义词列表，增加许多常用说法</span>

    <div class="font-bold">手册编写</div>
    <n-flex class="developers">
      <github-avatar username="只是另一个 ID" />
      <github-avatar username="JohNSoN" />
      <github-avatar username="暮星" />
      <github-avatar username="Szzrain" />
      <github-avatar username="木落" />
      <github-avatar
        username="山本健一"
        user="kenichiLyon"
        src="https://d1.sealdice.com/images/kenichiLyon.jpg" />
      <github-avatar username="逐风" />
      <github-avatar username="脑" />
      <github-avatar username="Bugtower" user="bugtower100" />
      <github-avatar username="Monad" user="Mitxoleta" />
      <github-avatar username="以炽热挥剑" user="yichirehuijian" />
      <github-avatar username="檀轶步棋" />
      <github-avatar username="流溪" user="lxy071130" />
      <github-avatar username="PaienNate" />
      <github-avatar username="病" user="nodisease" />
      <github-avatar username="綮小灰" user="Stanlty998" />
      <n-text>大概不全，下版本纠正……</n-text>
    </n-flex>

    <div class="font-bold">参考</div>
    <span>斯塔尼亚 - 塔系核心作者，指令实现过程中部分参考了塔系核心的指令表现</span>
    <span>赵喵喵 - ZhaoDice作者，主要指令参考之一</span>
    <span>Dice!核心的开发者们 - 同样的，在骰点格式和输出表现方面进行了参考</span>
    <span>FVTT - 经典的DND跑团平台，指令参考之一</span>
  </n-flex>
</template>

<script lang="ts" setup>
import { useStore } from '~/store';
import imgMascot from '~/assets/haibao1.png';

const store = useStore();

onBeforeMount(async () => {
  await store.logFetchAndClear();
});
</script>

<style scoped lang="css">
.about {
  padding: 2rem;
  line-height: 2rem;
  text-align: left;
  box-shadow:
    0 2px 4px rgba(0, 0, 0, 0.12),
    0 0 6px rgba(0, 0, 0, 0.04);
}

.developers {
  line-height: 4.5rem;
}

.developers > * {
  margin-right: 2rem;
}
</style>

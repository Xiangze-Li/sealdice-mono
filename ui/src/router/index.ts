import { createRouter, createWebHashHistory } from 'vue-router';
import PageAbout from '~/pages/PageAbout.vue';
import PageHome from '~/pages/PageHome.vue';
import PageConnectInfoItems from '~/pages/PageConnectInfoItems.vue';
import PageCustomText from '~/pages/PageCustomText.vue';
import PageCustomReply from '~/pages/mod/PageCustomReply.vue';
import PageJs from '~/pages/mod/PageJs.vue';
import PageMiscDeck from '~/pages/mod/PageMiscDeck.vue';
import PageHelpDoc from '~/pages/mod/PageHelpDoc.vue';
import PageStory from '~/pages/mod/PageStory.vue';
import PageCensor from '~/pages/mod/PageCensor.vue';
import PageTest from '~/pages/tool/PageTest.vue';
import PageResource from '~/pages/tool/PageResource.vue';
import PageMiscSettings from '~/pages/misc/PageMiscSettings.vue';
import PageMiscBackup from '~/pages/misc/PageMiscBackup.vue';
import PageMiscGroup from '~/pages/misc/PageMiscGroup.vue';
import PageMiscBan from '~/pages/misc/PageMiscBan.vue';
import PageMiscDicePublic from '~/pages/misc/PageMiscDicePublic.vue';
import PageMiscAdvancedSettings from '~/pages/misc/PageMiscAdvancedSettings.vue';

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/home', name: 'home', component: PageHome },
    { path: '/connect', component: PageConnectInfoItems },
    { path: '/custom-text/:category', component: PageCustomText, props: true },
    {
      path: '/mod',
      children: [
        { path: 'js', component: PageJs },
        { path: 'reply', component: PageCustomReply },
        { path: 'deck', component: PageMiscDeck },
        { path: 'helpdoc', component: PageHelpDoc },
        { path: 'story', component: PageStory },
        { path: 'censor', component: PageCensor },
      ],
    },
    {
      path: '/tool',
      children: [
        { path: 'test', component: PageTest },
        { path: 'resource', component: PageResource },
      ],
    },
    {
      path: '/misc',
      children: [
        { path: 'base-setting', component: PageMiscSettings },
        { path: 'backup', component: PageMiscBackup },
        { path: 'group', component: PageMiscGroup },
        { path: 'ban', component: PageMiscBan },
        { path: 'dice-public', component: PageMiscDicePublic },
        { path: 'advanced-setting', component: PageMiscAdvancedSettings },
      ],
    },
    { path: '/about', component: PageAbout },
    { path: '/:catchAll(.*)', name: 'default', redirect: { name: 'home' } },
  ],
});

export default router;

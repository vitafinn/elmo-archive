<script setup lang="ts">
import AppHeader from "./components/AppHeader.vue";
import AppMain from "./components/AppMain.vue";
import AppFooter from "./components/AppFooter.vue";
import {useLangStore} from "./store/langStore.ts";
import {useUserStore} from "./store/userStore.ts";
import {usePoolStore} from "./store/poolStore.ts";
import {useThemeStore} from "./store/themeStore.ts";
import {useImportStore} from './store/importStore';
import {computed, onBeforeMount, onMounted, onBeforeUnmount, watch} from "vue";
import zhCN from 'element-plus/es/locale/lang/zh-cn';
import zhTW from 'element-plus/es/locale/lang/zh-tw';
import en from 'element-plus/es/locale/lang/en';
import ImportSuccess from './components/ImportSuccess.vue';
import ImportConfirmation from './components/ImportConfirmation.vue';
import {EventsOn} from '../wailsjs/runtime/runtime';

const langStore = useLangStore()
const controlsLocale = computed(()=>langStore.lang==='zh-CN'?zhCN:langStore.lang==='zh-TW'?zhTW:en)
const userStore = useUserStore()
const poolStore = usePoolStore()
const flow = useImportStore()
useThemeStore()

onBeforeMount(async () => {
  await langStore.init()
  await userStore.init()
  await poolStore.init()
})

watch(() => [langStore.lang, userStore.userId, poolStore.poolType], async () => {
  await poolStore.updatePoolInfo()
})
let expiryTimer: ReturnType<typeof setInterval> | undefined
let stopProgress: (()=>void) | undefined
onMounted(() => {
  stopProgress=EventsOn('importProgress',(event:any)=>flow.addProgress(event))
  expiryTimer = setInterval(() => {
    if ((poolStore.poolInfo || poolStore.outfitInfo) && [6,7,-6,-7,9,10].includes(poolStore.poolType)) void poolStore.updatePoolInfo(true)
  }, 30000)
})
onBeforeUnmount(() => { if (expiryTimer) clearInterval(expiryTimer);stopProgress?.() })
</script>

<template>
  <el-config-provider :locale="controlsLocale">
  <div class="h-dvh w-dvw flex flex-col" 
       :style="{ background: `linear-gradient(to bottom, var(--theme-from), var(--theme-to))` }">
    <AppHeader/>
    <AppMain/>
    <AppFooter/>
  </div>
  <ImportConfirmation/>
  <ImportSuccess/>
  </el-config-provider>
</template>

<style>
:root {
  --theme-from: #0c1a32;
  --theme-to: #1c1917;
}
</style>

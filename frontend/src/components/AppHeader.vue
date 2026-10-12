<script setup lang="ts">
import {ref} from "vue";
import {WindowIsMaximised, WindowMaximise, WindowUnmaximise} from "../../wailsjs/runtime";
import ButtonVersion from "./ButtonVersion.vue";
import ButtonWindowMinimise from "./ButtonWindowMinimise.vue";
import ButtonWindowMaxmise from "./ButtonWindowMaxmise.vue";
import ButtonWindowClose from "./ButtonWindowClose.vue";
import SelectorLanguage from "./SelectorLanguage.vue";
import SelectorUser from "./SelectorUser.vue";
import ButtonColorPalette from "./ButtonColorPalette.vue";
import CaptureDialog from "./CaptureDialog.vue";
import {useImportStore} from '../store/importStore';
import HelpDialog from './HelpDialog.vue';

const flow = useImportStore();
const helpVisible = ref(false);

const handleWindowMaximise = async () => {
  await WindowIsMaximised().then(res => res ? WindowUnmaximise() : WindowMaximise())
}
</script>

<template>
  <div style="--wails-draggable:drag" class="app-header" @dblclick.self="handleWindowMaximise">
    <div class="app-brand">
      <span>{{ $t('window.title') }}</span>
      <ButtonVersion/>
    </div>
    
    <button class="header-control" @click="helpVisible = true">{{ $t('guide.help') }}</button>
    <button class="action-button primary-action" style="--wails-draggable:no-drag" @click="flow.visible = true">
      {{ $t('guide.import') }}
    </button>

    <SelectorUser/>
    <ButtonColorPalette/>
    <SelectorLanguage/>
    
    <ButtonWindowMinimise/>
    <ButtonWindowMaxmise/>
    <ButtonWindowClose/>
  </div>

  <CaptureDialog v-model="flow.visible" />
  <HelpDialog v-model="helpVisible" />
</template>

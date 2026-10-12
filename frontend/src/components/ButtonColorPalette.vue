<script setup lang="ts">
import {ref} from 'vue'
import {useThemeStore} from '../store/themeStore'
const theme=useThemeStore()
const visible=ref(false)
const mode=ref(theme.appearance),size=ref(theme.textSize)
const colors=ref({from:theme.gradientFrom,to:theme.gradientTo})
const open=()=>{mode.value=theme.appearance;size.value=theme.textSize;colors.value={from:theme.gradientFrom,to:theme.gradientTo};visible.value=true}
const apply=()=>{theme.updateAppearance(mode.value,size.value);theme.updateTheme(colors.value.from,colors.value.to);visible.value=false}
</script>
<template>
 <button class="header-control" style="--wails-draggable:no-drag" @click="open">{{ $t('appearance.title') }}</button>
 <el-dialog v-model="visible" class="import-dialog appearance-dialog" width="620" :title="$t('appearance.title')">
  <div class="import-guide">
   <p>{{ $t('appearance.description') }}</p>
   <fieldset class="appearance-options"><legend>{{ $t('appearance.palette') }}</legend>
    <label v-for="option in ['jade','slate','contrast']" :key="option" class="appearance-option" :class="{chosen:mode===option}">
     <input v-model="mode" type="radio" name="appearance-mode" :value="option">
     <span class="palette-preview" :class="option" aria-hidden="true"><i></i><i></i><i></i><i></i></span>
     <span><strong>{{ $t(`appearance.${option}`) }}</strong><small>{{ $t(`appearance.${option}Help`) }}</small></span>
    </label>
   </fieldset>
   <fieldset class="appearance-options"><legend>{{ $t('appearance.text') }}</legend>
    <label v-for="option in ['comfortable','large']" :key="option" class="appearance-option" :class="{chosen:size===option}"><input v-model="size" type="radio" name="appearance-size" :value="option"><span>{{ $t(`appearance.${option}`) }}</span></label>
   </fieldset>
   <div class="appearance-sample" :class="mode" :style="{fontSize:size==='large'?'18px':'16px'}"><span>{{ $t('appearance.sample') }}</span><strong>60 <small>{{ $t('character.pullsUnit') }}</small></strong><span>{{ $t('appearance.sampleNote') }}</span></div>
   <details><summary>{{ $t('appearance.frame') }}</summary><p class="detail-note mt-2">{{ $t('appearance.frameHelp') }}</p><div class="frame-colors"><label>{{ $t('theme.from') }}<input type="color" v-model="colors.from"></label><label>{{ $t('theme.to') }}<input type="color" v-model="colors.to"></label><button class="text-link" @click="colors={from:'#18212d',to:'#18212d'}">{{ $t('appearance.resetFrame') }}</button></div></details>
   <div class="editor-actions"><button class="action-button primary-action" @click="apply">{{ $t('theme.apply') }}</button><button class="action-button" @click="visible=false">{{ $t('theme.cancel') }}</button></div>
  </div>
 </el-dialog>
</template>

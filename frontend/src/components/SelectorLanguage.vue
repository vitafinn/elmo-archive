<script setup lang="ts">
import {useLangStore} from "../store/langStore.ts";

const langOptions = [
  {label: '简体中文', value: 'zh-CN'},
  {label: '繁體中文', value: 'zh-TW'},
  {label: 'English', value: 'en'},
  {label: '日本語', value: 'ja'},
  {label: '한국어', value: 'kr'},
]
const langStore = useLangStore()
</script>

<template>
  <el-dropdown style="--wails-draggable:no-drag">
    <button class="header-control" :aria-label="$t('appearance.language')">{{ langOptions.find(option=>option.value===langStore.lang)?.label || 'English' }} ⌄</button>
    <template #dropdown>
      <el-dropdown-menu>
        <el-dropdown-item v-for="lang in langOptions" :key="lang.value" @click="langStore.updateLang(lang.value)">
          <span>{{ lang.value===langStore.lang ? '✓ ' : '' }}{{ lang.label }}</span>
        </el-dropdown-item>
      </el-dropdown-menu>
    </template>
  </el-dropdown>
</template>

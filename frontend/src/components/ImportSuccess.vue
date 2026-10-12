<script setup lang="ts">
import {useImportStore} from '../store/importStore'
import {useCaptureSecurityStore} from '../store/captureSecurityStore'
const flow=useImportStore(),security=useCaptureSecurityStore()
const view=()=>{flow.successVisible=false;flow.visible=false}
const settings=()=>{flow.successVisible=false;flow.visible=true}
</script>
<template>
 <el-dialog v-model="flow.successVisible" class="import-dialog" width="560" :title="$t('syncSuccess.title')" :close-on-click-modal="false" append-to-body>
  <div class="sync-complete-content">
   <span class="sync-complete-icon" aria-hidden="true">✓</span>
   <p class="sync-complete-result" role="status">{{ flow.lastResult }}</p>
   <p>{{ $t('syncSuccess.done') }}</p>
   <p :class="security.error?'guide-error':'detail-note'">{{ $t(security.error?'syncSuccess.cleanupFailed':security.removed?'syncSuccess.removed':security.status?.Certificate?.Trusted?'syncSuccess.kept':'syncSuccess.none') }}</p>
   <p v-if="flow.error" class="guide-error" role="alert">{{ flow.error }}</p>
  </div>
  <template #footer><div class="editor-actions"><button class="action-button primary-action" @click="view">{{ $t('syncSuccess.view') }}</button><button class="action-button" @click="settings">{{ $t('syncSuccess.settings') }}</button></div></template>
 </el-dialog>
</template>

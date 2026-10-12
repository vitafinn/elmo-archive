<script setup lang="ts">
import ButtonSyncRecords from './ButtonSyncRecords.vue'
import {useImportStore} from '../store/importStore'
import {useUserStore} from '../store/userStore'
const flow = useImportStore()
const users = useUserStore()
</script>
<template>
  <footer class="app-status">
    <div role="status" class="text-sm min-w-0">
      <template v-if="flow.busy">{{ $t('sync.loading') }} <button class="text-link" @click="flow.visible=true">{{ $t('importReview.showLog') }}</button></template>
      <template v-else-if="flow.error"><span class="text-amber-200">{{ $t('guide.failed') }}</span> <button class="text-link" @click="flow.visible = true">{{ $t('guide.fixImport') }}</button></template>
      <template v-else-if="flow.result">{{ flow.result }} <span class="text-slate-400">· {{ flow.lastSync }}</span></template>
      <template v-else>{{ $t(users.userId ? 'guide.savedHistory' : 'guide.noAccount') }}</template>
    </div>
    <ButtonSyncRecords v-if="users.userId"/>
    <el-button v-else type="primary" @click="flow.visible = true">{{ $t('guide.import') }}</el-button>
  </footer>
</template>

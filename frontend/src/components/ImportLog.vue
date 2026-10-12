<script setup lang="ts">
import {useImportStore} from '../store/importStore'
const flow=useImportStore()
</script>
<template>
 <section class="import-log-section">
  <label class="appearance-option"><input type="checkbox" :checked="flow.verbose" @change="flow.setVerbose(($event.target as HTMLInputElement).checked)"><span>{{ $t('importReview.verbose') }}</span></label>
  <template v-if="flow.verbose"><p class="detail-note">{{ $t('importReview.logHelp') }}</p><ol class="import-log" role="log" aria-live="polite"><li v-for="(entry,index) in flow.entries" :key="index"><time>{{ entry.time }}</time> {{ $t(`importReview.stages.${entry.stage}`,{pool:entry.pool,page:entry.page,count:entry.count}) }}</li></ol></template>
 </section>
</template>

<script setup lang="ts">
import {useI18n} from 'vue-i18n'
import {handler} from '../../wailsjs/go/models'
import {outcomeKey,oddsFor} from '../utils/itemHistory'
import PullTiming from './PullTiming.vue'
const props=defineProps<{record:handler.DisplayRecord}>()
const emit=defineEmits<{select:[record:handler.DisplayRecord]}>()
const {t}=useI18n()
</script>
<template>
 <button type="button" class="pull-record" @click="emit('select',props.record)" :aria-label="t('character.open',{name:record.Name})">
  <span class="record-topline"><strong class="record-name">{{ record.Name }}</strong><strong class="record-count">{{ record.Count }} <small>{{ t('character.pullsUnit') }}</small></strong></span>
  <span class="record-date">{{ new Date(record.Timestamp*1000).toLocaleDateString($i18n.locale==='kr'?'ko':$i18n.locale) }}</span>
  <PullTiming :pool-type="record.PoolType" :count="record.Count"/>
  <span v-if="record.Outcome" class="outcome-pill" :class="`outcome-${record.Outcome}`">{{ t(outcomeKey(record.Outcome),{odds:oddsFor(record.PoolType)}) }}</span>
  <span v-else-if="record.IsMissing && [3,4].includes(record.PoolType)" class="outcome-pill outcome-loss">{{ t(record.PoolType === 3 ? 'guide.lostDoll' : 'guide.lostWeapon') }}</span>
  <span class="record-action">{{ t('character.view') }} <span aria-hidden="true">→</span></span>
 </button>
</template>


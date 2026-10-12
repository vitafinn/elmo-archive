<script setup lang="ts">
import {computed,ref} from 'vue'
import RecordSearch from './RecordSearch.vue'
import {matchesRecord} from '../utils/recordSearch'
import {useI18n} from 'vue-i18n'
import {handler} from '../../wailsjs/go/models'
import PoolStatistic from './PoolStatistic.vue'
import PoolRecord from './PoolRecord.vue'
const props = defineProps<{info: handler.PoolInfo}>()
const {t,locale} = useI18n()
const emit = defineEmits<{select:[record:handler.DisplayRecord]}>()
const query = ref('')
const matchingRecords = computed(() => (props.info.recordList || []).filter(record=>matchesRecord(query.value,record.Name,record.SearchNames)))
const matchingBanners = computed(() => (props.info.bannerPeriods || []).filter(banner=>!query.value.trim() || recordsByBanner.value.has(banner.id)))
const recordsByBanner = computed(() => {
  const groups = new Map<number, handler.DisplayRecord[]>()
  for (const record of matchingRecords.value) {
    const rows = groups.get(record.BannerID) || []
    rows.push(record); groups.set(record.BannerID, rows)
  }
  return groups
})
const date = (value: number) => value ? new Date(value * 1000).toLocaleString(locale.value==='kr'?'ko':locale.value) : t('selector.unknownDate')
const average = (value: number, count: number) => count ? value.toFixed(1) : '—'
const winRate = (wins: number, attempts: number) => attempts ? `${(100*wins/attempts).toFixed(1)}% (${wins}/${attempts})` : t('selector.noAttempts')
</script>

<template>
  <div class="flex flex-col gap-4 text-gray-100">
    <div class="stats-grid">
      <PoolStatistic :title="t('gacha.statistic.totalCount')" :value="info.totalCount"/>
      <PoolStatistic class="text-yellow-300" :title="t('gacha.statistic.rank5Data')" :value="info.rank5Count" :note="(info.rank5Rate*100).toFixed(2)+'%'"/>
      <PoolStatistic class="text-green-300" :title="t('selector.average')" :note="average(info.rank5Avg, info.rank5Count)"/>
      <PoolStatistic class="text-cyan-300" :title="t('selector.spent')" :note="average(info.rank5Spent, info.rank5Count)"/>
    </div>
    <details class="help-panel">
      <summary>{{ t('dashboard.details') }}</summary>
      <div class="stats-grid mt-4">
      <PoolStatistic class="text-purple-400" :title="t('gacha.statistic.rank4Data')" :value="info.rank4Count" :note="(info.rank4Rate*100).toFixed(2)+'%'"/>
      <PoolStatistic class="text-blue-300" :title="t('gacha.statistic.rank3Data')" :value="info.rank3Count" :note="(info.rank3Rate*100).toFixed(2)+'%'"/>
        <PoolStatistic :title="t('selector.expired')" :value="info.expiredPity"/>
        <PoolStatistic :title="t('selector.winRate')" :note="winRate(info.targetWins, info.targetAttempts)"/>
        <PoolStatistic v-if="info.unknownOutcomes" :title="t('selector.unknownOutcomes')" :value="info.unknownOutcomes"/>
        <PoolStatistic v-if="info.unknownRankCount" :title="t('selector.unknownRanks')" :value="info.unknownRankCount"/>
      </div>
      <p class="mt-4 text-sm text-slate-300">{{ t('selector.explanation') }}</p>
      <p class="mt-2 text-sm text-slate-400">{{ t('selector.historyNote') }}</p>
    </details>
    <div class="section-heading justify-between">
      <h2>{{ t('dashboard.bannerHistory') }}</h2>

    </div>
    <RecordSearch v-model="query" :count="matchingRecords.length"/>
    <p v-if="query && !matchingRecords.length" class="empty-state">{{ t('dashboard.noMatches') }}</p>
    <section v-for="banner in matchingBanners" :key="banner.id" class="history-panel">
      <header class="flex flex-wrap justify-between gap-2">
        <div><strong>{{ banner.name || t('selector.banner', {id: banner.id}) }}</strong><p v-if="banner.name" class="detail-note">{{ t('selector.banner', {id: banner.id}) }} · {{ t('luck.gameName') }}</p></div>
        <span class="status-badge" :class="banner.status === 'active' ? 'text-green-300' : 'text-gray-300'">{{ t(`selector.status.${banner.status}`) }}</span>
        <div class="w-full text-sm text-gray-300">{{ date(banner.startTime) }} — {{ date(banner.endTime) }}</div>
      </header>
      <p v-if="banner.status === 'active' && banner.tailCount" class="reward-highlight">{{ t('dashboard.expiringPity', {count: banner.tailCount}) }}</p>
      <dl class="period-metrics">
        <div><dt class="text-gray-400">{{ t('gacha.statistic.totalCount') }}</dt><dd>{{ banner.totalCount }}</dd></div>
        <div><dt class="text-gray-400">{{ t('gacha.statistic.rank5Data') }}</dt><dd>{{ banner.rank5Count }} · {{ banner.totalCount ? (100*banner.rank5Count/banner.totalCount).toFixed(2)+'%' : '—' }}</dd></div>
        <div><dt class="text-gray-400">{{ t('selector.average') }}</dt><dd>{{ average(banner.rank5Avg, banner.rank5Count) }}</dd></div>
        <div><dt class="text-gray-400">{{ t('selector.spent') }}</dt><dd>{{ average(banner.rank5Spent, banner.rank5Count) }}</dd></div>
        <div><dt class="text-gray-400">{{ t(banner.status === 'active' ? 'selector.active' : banner.status === 'expired' ? 'selector.expired' : 'selector.pending') }}</dt><dd>{{ banner.tailCount }}</dd></div>
        <div><dt class="text-gray-400">{{ t('selector.winRate') }}</dt><dd>{{ winRate(banner.targetWins, banner.targetAttempts) }}</dd></div>
      </dl>
      <p class="text-xs text-gray-400">{{ t('selector.lastRecorded') }}: {{ date(banner.lastTimestamp) }}. {{ t('selector.syncNote') }}</p>
      <details :open="!!query.trim() || banner.status === 'active'">
        <summary class="cursor-pointer mb-3">{{ t('gacha.records.title') }}</summary>
        <aside class="result-review-note"><strong>{{ t('character.reviewTitle') }}</strong><p>{{ t('character.reviewIntro') }}</p></aside>
        <div class="record-grid">
          <PoolRecord v-for="record in recordsByBanner.get(banner.id) || []" :key="record.ID" :record="record" @select="emit('select',$event)"/>
          <p v-if="!recordsByBanner.get(banner.id)?.length" class="text-gray-400 text-sm">{{ t('selector.noFiveStars') }}</p>
        </div>
      </details>
    </section>
  </div>
</template>


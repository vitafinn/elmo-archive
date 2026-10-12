<script setup lang="ts">
import {reactive} from 'vue'
import {useI18n} from 'vue-i18n'
import {handler} from '../../wailsjs/go/models'
import PoolStatistic from './PoolStatistic.vue'
import {matchesRecord} from '../utils/recordSearch'
const props = defineProps<{info: handler.OutfitInfo}>()
const {t,locale} = useI18n()
const limits = reactive<Record<number,number>>({})
const searches = reactive<Record<number,string>>({})
const categories = reactive<Record<number,string>>({})
const filtered = (season: handler.OutfitSeason) => season.records.filter(record =>
  (!categories[season.id] || record.category === categories[season.id]) &&
  matchesRecord(searches[season.id] || '',record.name,record.searchNames))
const resetLimit = (id: number) => { limits[id] = 50 }
const date = (value:number) => value ? new Date(value*1000).toLocaleString(locale.value==='kr'?'ko':locale.value) : t('selector.unknownDate')
const colors:Record<string,string> = {limited:'text-yellow-300',rare:'text-purple-300',normal:'text-gray-100',unknown:'text-gray-400'}
</script>
<template>
  <div class="flex flex-col gap-4 text-gray-100">
    
    <PoolStatistic :title="t('outfit.total')" :value="props.info.totalCount"/>
    <details class="help-panel"><summary>{{ t('dashboard.details') }}</summary><p class="mt-3 text-sm text-slate-300">{{ t('outfit.note') }}</p></details>
    <h2 class="section-heading">{{ t('dashboard.seasonHistory') }}</h2>
    <p v-if="!props.info.totalCount" class="text-gray-300">{{ t('outfit.empty') }}</p>
    <section v-for="season in props.info.seasons" :key="season.id" class="history-panel">
      <header class="flex flex-wrap justify-between gap-2">
        <strong>{{ season.name }}</strong><span class="status-badge" :class="{'text-emerald-300':season.status === 'active'}">{{ t(`outfit.status.${season.status}`) }}</span>
        <div class="w-full text-sm text-gray-300">{{ date(season.startTime) }} — {{ date(season.endTime) }}</div>
      </header>
      <div class="stats-grid">
        <PoolStatistic :title="t('outfit.total')" :value="season.totalCount"/>
        <PoolStatistic class="text-yellow-300" :title="t('outfit.limited')" :value="season.limitedCount"/>
        <PoolStatistic class="text-purple-300" :title="t('outfit.rare')" :value="season.rareCount"/>
        <PoolStatistic :title="t('outfit.normal')" :value="season.normalCount"/>
        <PoolStatistic v-if="season.unknownCount" :title="t('outfit.unknown')" :value="season.unknownCount"/>
      </div>
      <p v-if="season.mainOutfit" class="reward-highlight">{{ t('outfit.main') }}: {{ season.mainOutfit }} · {{ season.mainOutfitAt ? t('outfit.obtained',{count:season.mainOutfitAt}) : t('outfit.notObtained') }}</p>
      <details :open="season.totalCount>0">
        <summary class="cursor-pointer">{{ t('outfit.records') }} ({{ season.totalCount }})</summary>
        <div class="record-filters">
          <input v-model="searches[season.id]" type="search" :aria-label="t('recordSearch.label')" :placeholder="t('recordSearch.placeholder')" @input="resetLimit(season.id)">
          <select v-model="categories[season.id]" :aria-label="t('outfit.category')" @change="resetLimit(season.id)">
            <option :value="undefined">{{ t('dashboard.allRewards') }}</option>
            <option v-for="category in ['limited','rare','normal','unknown']" :key="category" :value="category">{{ t(`outfit.categories.${category}`) }}</option>
          </select>
        </div>
        <div class="overflow-x-auto mt-3">
          <table class="records-table">
            <thead class="text-gray-400"><tr><th class="p-2">{{ t('outfit.number') }}</th><th class="p-2">{{ t('outfit.reward') }}</th><th class="p-2">{{ t('outfit.category') }}</th><th class="p-2">{{ t('outfit.time') }}</th></tr></thead>
            <tbody><tr v-for="record in filtered(season).slice(0,limits[season.id] || 50)" :key="record.id + ':' + record.accessNumber" class="border-t border-white/10">
              <td class="p-2">{{ record.accessNumber }}</td><td class="p-2" :class="colors[record.category]">{{ record.name }}</td><td class="p-2">{{ t(`outfit.categories.${record.category}`) }}</td><td class="p-2 whitespace-nowrap">{{ date(record.timestamp) }}</td>
            </tr><tr v-if="!filtered(season).length"><td colspan="4" class="empty-state">{{ t('dashboard.noMatches') }}</td></tr></tbody>
          </table>
        </div>
        <button v-if="filtered(season).length > (limits[season.id] || 50)" class="mt-3 rounded bg-blue-500/30 px-4 py-2" @click="limits[season.id] = (limits[season.id] || 50)+50">{{ t('outfit.more') }}</button>
      </details>
    </section>
  </div>
</template>


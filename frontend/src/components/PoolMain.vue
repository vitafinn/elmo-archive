<script setup lang="ts">
import {ref,computed,watch} from 'vue';
import RecordSearch from './RecordSearch.vue';
import {matchesRecord} from '../utils/recordSearch';
import {handler} from '../../wailsjs/go/models';
import CharacterDetails from './CharacterDetails.vue';
import {usePoolStore} from "../store/poolStore.ts";
import PoolStatistic from "./PoolStatistic.vue";
import PoolRecord from "./PoolRecord.vue";
import SelectorStats from "./SelectorStats.vue";
import OutfitStats from "./OutfitStats.vue";
import {useUserStore} from '../store/userStore';
import {useImportStore} from '../store/importStore';

const poolStore = usePoolStore()
const users = useUserStore()
const flow = useImportStore()
const selectedRecord = ref<handler.DisplayRecord|null>(null)
const query = ref('')
const matchingRecords = computed(() => (poolStore.poolInfo?.recordList || []).filter(record => matchesRecord(query.value,record.Name,record.SearchNames)))
watch(()=>[users.userId,poolStore.poolType],()=>{query.value=''})

</script>

<template>
  <div class="flex flex-col gap-6">
    <section v-if="!users.userId" class="welcome-panel">
      <p class="eyebrow">Elmo Archive · Girls’ Frontline 2: Exilium</p>
      <h1>{{ $t('guide.welcome') }}</h1>
      <p>{{ $t('guide.welcomeBody') }}</p>
      <ol><li>{{ $t('guide.welcome1') }}</li><li>{{ $t('guide.welcome2') }}</li><li>{{ $t('guide.welcome3') }}</li></ol>
      <el-button type="primary" @click="flow.visible = true">{{ $t('guide.firstImport') }}</el-button>
    </section>
    <header v-else class="page-heading">
      <p class="eyebrow">{{ $t('dashboard.overview') }}</p>
      <h1>{{ $te(`gacha.type.${poolStore.poolType}`) ? $t(`gacha.type.${poolStore.poolType}`) : $t('gacha.unknownType', {id: poolStore.poolType}) }}</h1>
      <p class="text-sm text-slate-400">{{ $t('dashboard.allHistory') }}</p>
    </header>
    <OutfitStats v-if="poolStore.outfitInfo" :info="poolStore.outfitInfo"/>
    <template v-else-if="poolStore.poolInfo">
      <p v-if="!poolStore.poolInfo.totalCount" class="help-panel">{{ $t('guide.noRecords') }}</p>
      <SelectorStats v-if="[6,7,-6,-7].includes(poolStore.poolType)" :key="users.userId+':'+poolStore.poolType" :info="poolStore.poolInfo" @select="selectedRecord=$event"/>
      <template v-else>
      <div class="stats-grid">
        <PoolStatistic class="text-gray-50" :title="$t('gacha.statistic.totalCount')" :value="poolStore.poolInfo.totalCount"/>
        <PoolStatistic class="text-orange-400" :title="$t('gacha.statistic.pityCount')" :value="poolStore.poolInfo.storedCount"/>
        <PoolStatistic class="text-yellow-300" :title="$t('gacha.statistic.rank5Data')" :value="poolStore.poolInfo.rank5Count" :note="(poolStore.poolInfo.rank5Rate*100).toFixed(2)+'%'"/>
        <PoolStatistic class="text-purple-500" :title="$t('gacha.statistic.rank4Data')" :value="poolStore.poolInfo.rank4Count" :note="(poolStore.poolInfo.rank4Rate*100).toFixed(2)+'%'"/>
        <PoolStatistic class="text-blue-500" :title="$t('gacha.statistic.rank3Data')" :value="poolStore.poolInfo.rank3Count" :note="(poolStore.poolInfo.rank3Rate*100).toFixed(2)+'%'"/>
        <PoolStatistic class="text-green-400" :title="$t('gacha.statistic.rank5Avg')" :value="poolStore.poolInfo.rank5Avg"/>
        <PoolStatistic v-if="[3, 4, 6, 7, -6, -7].includes(poolStore.poolType)" class="text-pink-400" :title="$t('gacha.statistic.upRank5Avg')" :value="poolStore.poolInfo.rank5UpAvg"/>
        <PoolStatistic v-if="[3, 4, 6, 7, -6, -7].includes(poolStore.poolType)" class="text-red-500" :title="$t('gacha.statistic.nonUpRate')" :note="Math.round(poolStore.poolInfo.missingRate*1000) /10+'%'"/>
      </div>

      <div class="section-heading">
        <span>{{ $t('gacha.records.title') }}</span>
        <el-tooltip effect="dark" popper-class="records-help" :content="$t('gacha.records.tip')" placement="top">
          <button class="record-help-button" :aria-label="$t('gacha.records.tip')">?</button>
          <template #content>
            <div class="record-help-content">{{$t('gacha.records.tip')}}</div>
          </template>
        </el-tooltip>
      </div>

      <p class="text-sm text-slate-400">{{ $t('character.recordsHint') }}</p>
      <RecordSearch v-model="query" :count="matchingRecords.length"/>
      <p v-if="query && !matchingRecords.length" class="empty-state">{{ $t('dashboard.noMatches') }}</p>
      <el-scrollbar>
        <div class="record-grid">
          <PoolRecord v-for="record in matchingRecords" :key="record.ID" :record="record" @select="selectedRecord=$event"/>
        </div>
      </el-scrollbar>
      </template>
    </template>
    <div v-else-if="users.userId" class="empty-state" role="status">
      <template v-if="poolStore.loading">{{ $t('guide.loading') }}</template>
      <template v-else>{{ $t(poolStore.loadError ? 'guide.loadFailed' : 'guide.noRecords') }}
        <button class="text-link ml-2" @click="poolStore.updatePoolInfo()">{{ $t('guide.retry') }}</button>
      </template>
    </div>
    <CharacterDetails :record="selectedRecord" @close="selectedRecord=null"/>
  </div>
</template>

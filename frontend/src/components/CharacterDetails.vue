<script setup lang="ts">
import {computed,ref,watch,nextTick} from 'vue'
import {useI18n} from 'vue-i18n'
import {handler} from '../../wailsjs/go/models'
import {GetItemHistory,GetDollWeaponHistory,SetSelectorOutcome} from '../../wailsjs/go/handler/App'
import {useUserStore} from '../store/userStore'
import {usePoolStore} from '../store/poolStore'
import {summarizeCopies,withFirstWeapon,outcomeKey,oddsFor} from '../utils/itemHistory'
import PullTiming from './PullTiming.vue'
import MilestoneLuck from './MilestoneLuck.vue'
const props = defineProps<{record:handler.DisplayRecord | null}>()
const emit = defineEmits<{close:[]}>()
const users=useUserStore(), pools=usePoolStore()
const {t,te}=useI18n()
const records=ref<handler.DisplayRecord[]>([])
const content=ref<HTMLElement|null>(null)
const selectedIndex=computed(()=>records.value.findIndex(pull=>pull.ID===props.record?.ID))
const showSummary=()=>content.value?.closest('.el-dialog__body')?.scrollTo({top:0})
const locate=async()=>{
 await nextTick()
 const element=content.value?.querySelector<HTMLElement>('.selected-copy')
 element?.scrollIntoView({block:'center',behavior:'instant'})
 element?.focus({preventScroll:true})
}
const weapon=ref<handler.DollWeaponHistory|null>(null), weaponError=ref('')
const loading=ref(false), error=ref(''), saving=ref(false), saved=ref(false)
const editing=ref<string|null>(null), choice=ref('auto')
const savedCopy=ref<string|null>(null),failedCopy=ref<string|null>(null)
let generation=0
const load=async(background=false)=>{
  const request=++generation
  if(!props.record || !users.userId) return
  if(!background)loading.value=true; error.value=''
  try {
    const [history,related]=await Promise.allSettled([
      GetItemHistory(users.userId,props.record.ItemID),
      props.record.ItemType===10 ? GetDollWeaponHistory(users.userId,props.record.ItemID) : Promise.resolve(null)
    ])
    if(request!==generation)return
    if(history.status==='fulfilled') records.value=history.value || []
    else error.value=String(history.reason)
    if(related.status==='fulfilled') {weapon.value=related.value;weaponError.value=''}
    else {weapon.value=null;weaponError.value=String(related.reason)}
  }
  catch(err) { if(request===generation) error.value=String(err) }
  finally { if(request===generation) loading.value=false }
}
watch(()=>[props.record?.ID,users.userId],()=>{ records.value=[];weapon.value=null;weaponError.value='';editing.value=null;saved.value=false;savedCopy.value=null;failedCopy.value=null;void load() },{immediate:true})
const totals=computed(()=>summarizeCopies(records.value))
const paired=computed(()=>withFirstWeapon(records.value,weapon.value?.Records || []))
const weaponCoverage=computed(()=>summarizeCopies((weapon.value?.Records || []).slice(0,paired.value.v6r6!==null ? 6 : 1)))
const isDoll=computed(()=>props.record?.ItemType===10)
const hasSelector=computed(()=>records.value.some(pull=>[6,7,-6,-7].includes(pull.PoolType)))
const label=(record:handler.DisplayRecord)=>t(outcomeKey(record.Outcome),{odds:oddsFor(record.PoolType)})
const category=(type:number)=>te(`gacha.type.${type}`)?t(`gacha.type.${type}`):t('gacha.unknownType',{id:type})
const options=['auto','win','loss','initial','guaranteed','unknown']
const startEdit=(record:handler.DisplayRecord)=>{ editing.value=record.ID;choice.value=record.Annotation || 'auto';saved.value=false;savedCopy.value=null;failedCopy.value=null;error.value='' }
const save=async(record:handler.DisplayRecord)=>{
 if(!users.userId || saving.value) return
 const account=users.userId,item=props.record?.ID
 const scroller=content.value?.closest('.el-dialog__body')
 const row=Array.from(content.value?.querySelectorAll<HTMLElement>('[data-copy-id]') || []).find(el=>el.dataset.copyId===record.ID)
 const anchor=row?.querySelector<HTMLElement>('.review-result-button')
 const before=anchor?.getBoundingClientRect().top
 saving.value=true;error.value='';failedCopy.value=null
 try {
  await SetSelectorOutcome(account,record.ID,choice.value)
  await load(true);await pools.updatePoolInfo(true)
  if(account!==users.userId || item!==props.record?.ID)return
  saved.value=!error.value && !pools.loadError
  if(saved.value){editing.value=null;savedCopy.value=record.ID}
  else {failedCopy.value=record.ID;if(!error.value)error.value=t('character.failed')}
 } catch(err) {error.value=String(err);failedCopy.value=record.ID}
 finally {
  saving.value=false
  await nextTick()
  if(account===users.userId && item===props.record?.ID && anchor?.isConnected){
   if(scroller && before!==undefined)scroller.scrollTop+=anchor.getBoundingClientRect().top-before
   if(saved.value)anchor.focus({preventScroll:true})
  }
 }
}
const close=(done?:()=>void)=>{if(saving.value)return;generation++;emit('close');done?.()}
</script>
<template>
 <el-dialog :model-value="!!record" class="import-dialog character-dialog" width="850" :title="record?.Name" :before-close="close" :close-on-click-modal="!saving" @opened="showSummary" @update:model-value="(value:boolean)=>{if(!value)close()}">
  <div class="character-content" ref="content">
   <p v-if="selectedIndex>=0" class="detail-note">{{ t('copyLocation.selected',{number:selectedIndex+1,total:records.length}) }}</p>
   <p class="detail-scope">{{ t('character.scope') }}</p>
   <p v-if="loading" role="status">{{ t('guide.loading') }}</p>
   <template v-else-if="records.length">
    <div class="character-metrics">
     <div><span>{{ t('character.copies') }}</span><strong>{{ totals.copies }}</strong></div>
     <div v-if="isDoll"><span>{{ t('character.vertebrae') }}</span><strong>V{{ Math.min(totals.duplicates,6) }}</strong></div>
     <div class="acquisition-total"><span>{{ t('acquisition.total') }}</span><strong>{{ totals.acquisition }}</strong></div>
     <div><span>{{ t('acquisition.average') }}</span><strong>{{ totals.acquisitionAverage.toFixed(1) }}</strong></div>
    </div>
    <section class="cost-overview">
      <div class="cost-legend"><span><i class="direct-swatch"></i>{{ t('acquisition.direct') }} <strong>{{ totals.pulls }}</strong></span><span><i class="loss-swatch"></i>{{ t('acquisition.losses') }} <strong>{{ totals.lossPulls }}</strong></span></div>
      <div v-if="totals.acquisition" class="cost-bar" role="img" :aria-label="t('acquisition.split',{direct:totals.pulls,loss:totals.lossPulls,total:totals.acquisition})"><span class="direct-segment" :style="{width:100*totals.pulls/totals.acquisition+'%'}"></span><span class="loss-segment" :style="{width:100*totals.lossPulls/totals.acquisition+'%'}"></span></div>
      <div class="cost-extremes"><span>{{ t('acquisition.cheapest') }} <strong>{{ t('guide.pulls',{count:totals.cheapest}) }}</strong></span><span>{{ t('acquisition.highest') }} <strong>{{ t('guide.pulls',{count:totals.highest}) }}</strong></span></div>
      <p v-if="totals.highestRecord?.PreviousLoss" class="detail-note mt-3">{{ t('acquisition.highestBreakdown',{loss:totals.highestRecord.PreviousLoss.Count,name:totals.highestRecord.PreviousLoss.Name,sequence:totals.highestRecord.Count,target:record?.Name,total:totals.highest}) }}</p>
      <p v-if="totals.highestRecord?.PreviousLoss && !['linked','recordedLoss'].includes(totals.highestRecord.CostBasis)" class="detail-note">{{ t(`acquisition.basis.${totals.highestRecord.CostBasis}`) }}</p>
      <div class="cost-extremes mt-3"><span>{{ t('acquisition.shortestSequence') }} <strong>{{ t('guide.pulls',{count:totals.shortestSequence}) }}</strong></span><span>{{ t('acquisition.longestSequence') }} <strong>{{ t('guide.pulls',{count:totals.longestSequence}) }}</strong></span></div>
    </section>
    <section v-if="isDoll" class="weapon-panel">
      <header class="weapon-heading"><h3>{{ t('weapon.title') }}</h3><span class="detail-note" v-if="weapon?.ItemID">{{ t('weapon.source') }}</span></header>
      <template v-if="weapon?.ItemID">
        <p class="weapon-name">{{ weapon.Name }}</p>
        <template v-if="paired.first">
          <div class="character-metrics">
            <div><span>{{ t('weapon.first') }}</span><strong>{{ paired.weaponCost }} <small>{{ t('character.pullsUnit') }}</small></strong></div>
            <div class="acquisition-total"><span>{{ t('weapon.combined',{level:paired.level}) }}</span><strong>{{ paired.combined }} <small>{{ t('character.pullsUnit') }}</small></strong></div>
            <div v-if="paired.v6r6!==null" class="acquisition-total"><span>{{ t('weapon.v6r6') }}</span><strong>{{ paired.v6r6 }} <small>{{ t('character.pullsUnit') }}</small></strong></div>
          </div>
          <p class="weapon-equation">{{ t('weapon.formula',{doll:paired.dollCost,weapon:paired.weaponCost,total:paired.combined}) }}</p>
          <p v-if="paired.v6r6!==null" class="detail-note">{{ t('weapon.v6r6Note') }}</p>
          <p v-if="paired.first.CostBasis && paired.first.CostBasis!=='direct'" :class="['linked','recordedLoss'].includes(paired.first.CostBasis)?'detail-note':'cost-caveat'">{{ t(`acquisition.basis.${paired.first.CostBasis}`,{count:paired.first.PreviousLoss?.Count,odds:oddsFor(paired.first.PoolType)}) }}</p>
          <p v-if="weaponCoverage.estimates || weaponCoverage.unknown" class="cost-caveat">{{ t('acquisition.coverage',{estimates:weaponCoverage.estimates,unknown:weaponCoverage.unknown}) }}</p>
          <details><summary>{{ t('weapon.history') }} · {{ t('weapon.copies',{count:weapon.Records.length}) }}</summary>
            <article v-for="(copy,index) in weapon.Records" :key="copy.ID" class="weapon-copy">
              <header><strong>{{ t('character.copy',{count:index+1}) }}</strong><strong>{{ t('guide.pulls',{count:copy.AcquisitionCount || copy.Count}) }}</strong></header>
              <p class="detail-note">{{ new Date(copy.Timestamp*1000).toLocaleDateString($i18n.locale==='kr'?'ko':$i18n.locale) }} · {{ copy.BannerName || category(copy.PoolType) }}</p>
              <p v-if="copy.PreviousLoss" class="weapon-loss">{{ t('weapon.loss',{count:copy.PreviousLoss.Count,name:copy.PreviousLoss.Name}) }} +</p>
              <p>{{ t('weapon.sequence',{count:copy.Count}) }}</p>
              <p v-if="copy.CostBasis && copy.CostBasis!=='direct'" class="detail-note">{{ t(`acquisition.basis.${copy.CostBasis}`,{count:copy.PreviousLoss?.Count,odds:oddsFor(copy.PoolType)}) }}</p>
              <PullTiming :pool-type="copy.PoolType" :count="copy.Count"/>
            </article>
          </details>
          <p class="detail-note">{{ t('weapon.note') }}</p>
        </template>
        <p v-else class="detail-note">{{ t('weapon.noRecords') }}</p>
      </template>
      <template v-else-if="weaponError"><p>{{ t('weapon.failed') }}</p><button class="action-button" @click="load()">{{ t('guide.retry') }}</button></template>
      <p v-else class="detail-note">{{ t('weapon.noLink') }}</p>
    </section>
    <MilestoneLuck v-if="isDoll" :key="record?.ItemID" :dolls="records" :weapons="weapon?.Records || []"/>
    <p class="detail-note">{{ t('acquisition.note') }}</p>
    <p v-if="totals.estimates || totals.unknown" class="cost-caveat">{{ t('acquisition.coverage',{estimates:totals.estimates,unknown:totals.unknown}) }}</p>
    <details class="help-panel"><summary>{{ t('character.howCalculated') }}</summary>
      <p v-if="isDoll" class="detail-note mt-3">{{ t('character.vertebraeNote',{count:totals.duplicates}) }}</p>
      <p class="detail-note mt-3">{{ t('acquisition.method') }}</p>
      <p class="detail-note mt-3">{{ t('luck.assumption') }}</p>
      <p class="detail-note mt-3">{{ t('luck.model') }}</p>
      <p class="detail-note mt-3">{{ t('luck.bands') }}</p>
      <p class="detail-note mt-3">{{ t('luck.chart') }}</p>
    </details>
    <h3 class="section-heading">{{ t('character.history') }}</h3>
    <aside v-if="hasSelector" class="result-review-note"><strong>{{ t('character.reviewTitle') }}</strong><p>{{ t('character.editHint') }}</p></aside>
    <article v-for="(pull,index) in records" :key="pull.ID" :data-copy-id="pull.ID" tabindex="-1" class="copy-row" :class="{'selected-copy':pull.ID===record?.ID}">
     <span v-if="pull.ID===record?.ID" class="selected-copy-badge">{{ t('copyLocation.badge') }}</span>
     <header><div><strong>{{ t('character.copy',{count:index+1}) }}</strong><span class="copy-date">{{ new Date(pull.Timestamp*1000).toLocaleString($i18n.locale==='kr'?'ko':$i18n.locale) }}</span></div><div class="copy-total"><strong class="copy-cost">{{ t('guide.pulls',{count:pull.AcquisitionCount || pull.Count}) }}</strong><small>{{ t('acquisition.copyTotal') }}</small></div></header>
     <p class="copy-category">{{ pull.BannerName || category(pull.PoolType) }} · {{ t('selector.banner',{id:pull.BannerID}) }}</p>
     <div v-if="pull.PreviousLoss" class="acquisition-chain">
       <div class="loss-step"><span>{{ t('acquisition.previousLoss') }}</span><strong>{{ pull.PreviousLoss.Name }}</strong><b>{{ t('guide.pulls',{count:pull.PreviousLoss.Count}) }}</b><small>{{ t('selector.banner',{id:pull.PreviousLoss.BannerID}) }}</small></div>
       <span class="chain-plus" aria-hidden="true">+</span>
       <div class="obtained-step"><span>{{ t('acquisition.obtained') }}</span><strong>{{ pull.Name }}</strong><b>{{ t('guide.pulls',{count:pull.Count}) }}</b></div>
     </div>
     <p v-if="pull.CostBasis && pull.CostBasis !== 'direct'" :class="['detail-note','cost-basis',{'cost-caveat':!['linked','recordedLoss'].includes(pull.CostBasis)}]">{{ t(`acquisition.basis.${pull.CostBasis}`,{count:pull.PreviousLoss?.Count,odds:oddsFor(pull.PoolType)}) }}</p>
     <PullTiming :pool-type="pull.PoolType" :count="pull.Count" expanded/>
     <div v-if="[6,7,-6,-7].includes(pull.PoolType)" class="copy-outcome">
      <span class="outcome-pill" :class="`outcome-${pull.Outcome}`">{{ label(pull) }}</span>
      <span class="detail-note">{{ t(pull.Annotation && pull.Annotation!=='auto' ? 'character.manual' : 'selector.outcome.auto') }}</span>
      <button class="action-button review-result-button" :disabled="saving" :aria-expanded="editing===pull.ID" @click="startEdit(pull)">{{ t('character.correct') }} <span aria-hidden="true">✎</span></button>
     </div>
     <p v-if="[6,7,-6,-7].includes(pull.PoolType) && (!pull.Annotation || pull.Annotation==='auto')" class="detail-note">{{ t('character.autoReview') }}</p>
     <form v-if="editing===pull.ID" class="outcome-editor" @submit.prevent="save(pull)">
      <fieldset :disabled="saving"><legend>{{ t('character.question') }}</legend>
       <div class="outcome-choices"><label v-for="option in options" :key="option" class="outcome-option" :class="{chosen:choice===option}">
        <input v-model="choice" type="radio" :value="option" :name="`outcome-${pull.ID}`">
        <span><strong>{{ t(outcomeKey(option),{odds:oddsFor(pull.PoolType)}) }}</strong></span>
       </label></div>
       <p class="outcome-selection-help" aria-live="polite">{{ t(`character.optionHelp.${choice}`) }}</p>
      </fieldset>
      <details class="outcome-save-help"><summary>{{ t('character.saveChanges') }}</summary><p class="detail-note">{{ t('character.saveNote') }}</p></details>
      <div class="editor-actions"><button class="action-button primary-action" type="submit" :disabled="saving">{{ t(saving?'character.saving':'character.save') }}</button><button type="button" class="action-button" :disabled="saving" @click="editing=null">{{ t('character.cancel') }}</button></div>
     </form>
     <p v-if="savedCopy===pull.ID" class="save-confirmation" role="status">✓ {{ t('character.saved') }}</p>
     <p v-if="failedCopy===pull.ID" class="guide-error" role="alert">{{ t('character.failed') }} {{ error }}</p>
    </article>
   </template>
   <p v-else-if="!loading && !error">{{ t('guide.noRecords') }}</p>
   <div v-if="error && !failedCopy" class="guide-error" role="alert"><p>{{ t('character.failed') }}</p><button class="action-button" :disabled="saving" @click="load()">{{ t('guide.retry') }}</button><details><summary>{{ t('guide.technical') }}</summary>{{ error }}</details></div>

  </div>
  <template #footer><div v-if="selectedIndex>=0" class="copy-bottom-nav"><button class="text-link" @click="showSummary">{{ t('copyLocation.summary') }}</button><button class="copy-bubble" @click="locate"><span aria-hidden="true">↥</span>{{ t('copyLocation.jump') }} <small>#{{ selectedIndex+1 }}</small></button></div></template>
 </el-dialog>
</template>




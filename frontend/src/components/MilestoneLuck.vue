<script setup lang="ts">
import {computed,ref,useId} from 'vue'
import {useI18n} from 'vue-i18n'
import {handler} from '../../wailsjs/go/models'
import {chartPercentiles,milestoneComparison,milestoneScale} from '../utils/milestoneLuck'
const props=defineProps<{dolls:handler.DisplayRecord[],weapons:handler.DisplayRecord[]}>()
const {t}=useI18n()
const rank=ref(0),level=ref(0)
const explanationOpen=ref(false), explanationId=useId()
const result=computed(()=>milestoneComparison(props.dolls,props.weapons,level.value,rank.value))
const compare=(v:number)=>milestoneComparison(props.dolls,props.weapons,v,rank.value)
const scale=computed(()=>result.value ? milestoneScale(result.value.cost,result.value.row) : null)
const medianDifference=computed(()=>!scale.value?'':t(`milestoneLuck.${scale.value.delta<0?'belowMedian':scale.value.delta>0?'aboveMedian':'atMedian'}`,{count:Math.abs(scale.value.delta)}))
const graphicLabel=computed(()=>result.value ? t('milestoneLuck.graphicLabel',{cost:result.value.cost,median:result.value.row[2],low:result.value.lower,high:result.value.upper}) : '')
const percentileText=computed(()=>!result.value?'':result.value.percentile===null ? (result.value.cost<result.value.row[0]?'<10%':'>99%') : `≈${Math.round(result.value.percentile)}%`)
</script>
<template>
 <section class="milestone-luck">
  <h3>{{ t('milestoneLuck.title') }}</h3>
  <p class="detail-note">{{ t('milestoneLuck.intro') }}</p>
  <div class="milestone-toggle" role="group" :aria-label="t('milestoneLuck.weapon')">
   <button v-for="r in [0,1,6]" :key="r" class="action-button" :aria-pressed="rank===r" @click="rank=r">{{ t(`milestoneLuck.r${r}`) }}</button>
  </div>
  <div class="luck-milestones" role="group" :aria-label="t('milestoneLuck.title')">
   <button v-for="v in 7" :key="v" :aria-pressed="level===v-1" @click="level=v-1">
    <strong>V{{ v-1 }}{{ rank ? `R${rank}` : '' }}</strong>
    <small>{{ compare(v-1) ? t('guide.pulls',{count:compare(v-1)!.cost}) : '—' }}</small>
   </button>
  </div>
  <div class="milestone-result" aria-live="polite">
   <template v-if="result">
    <strong>V{{ level }}{{ rank ? `R${rank}` : '' }} · {{ t('guide.pulls',{count:result.cost}) }}</strong>
    <p v-if="result.unknown" class="detail-note">{{ t('milestoneLuck.unknown') }}</p>
    <template v-else>
     <p><span class="milestone-badge" :class="result.reference?'reference':result.band">{{ t(`milestoneLuck.${result.reference?'reference':result.band}`) }}</span> <span v-if="result.estimated" class="detail-note">{{ t('milestoneLuck.estimated') }}</span></p>
     <div class="percentile-summary">
      <div><span class="metric-label">{{ t('milestoneLuck.percentileTitle') }}</span>
       <div class="percentile-value"><strong>{{ percentileText }}</strong><button class="info-button" :aria-label="t('milestoneLuck.explain')" :aria-expanded="explanationOpen" :aria-controls="explanationId" @click="explanationOpen=!explanationOpen">ⓘ</button></div>
       <span class="metric-direction">{{ t('milestoneLuck.lowerBetter') }}</span>
      </div>
      <p class="median-difference">{{ medianDifference }}</p>
     </div>
     <aside v-show="explanationOpen" :id="explanationId" class="explanation-panel">
      <strong>{{ t('milestoneLuck.explain') }}</strong>
      <p>{{ t('milestoneLuck.interpolationHelp') }}</p>
      <p>{{ t('milestoneLuck.range',{low:result.lower,high:result.upper}) }}</p>
      <p>{{ t('milestoneLuck.medianHelp') }}</p>
      <p>{{ t('milestoneLuck.graphicHelp') }}</p>
     </aside>
     <div v-if="scale" class="comparison-graphic" tabindex="0" role="region" :aria-label="t('milestoneLuck.graphicTitle')">
      <svg viewBox="0 0 640 250" role="img" :aria-label="graphicLabel">
       <title>{{ graphicLabel }}</title>
       <rect x="32" y="72" width="576" height="16" rx="8" fill="var(--ui-input)"/>
       <rect x="32" y="72" :width="scale.position(result.row[1])-32" height="16" fill="var(--ui-success)" opacity=".5"/>
       <rect :x="scale.position(result.row[1])" y="72" :width="scale.position(result.row[3])-scale.position(result.row[1])" height="16" fill="var(--ui-muted)" opacity=".35"/>
       <rect :x="scale.position(result.row[3])" y="72" :width="scale.position(result.row[4])-scale.position(result.row[3])" height="16" fill="var(--ui-warning)" opacity=".5"/>
       <rect :x="scale.position(result.row[4])" y="72" :width="608-scale.position(result.row[4])" height="16" fill="var(--ui-loss)" opacity=".5"/>
       <g v-for="(percent,i) in chartPercentiles" :key="percent">
        <line :x1="scale.position(result.row[i])" y1="88" :x2="scale.position(result.row[i])" :y2="105+(i%4)*30" stroke="var(--ui-muted)" opacity=".5"/>
        <text :x="scale.position(result.row[i])" :y="119+(i%4)*30" text-anchor="middle">{{ percent }}% · {{ result.row[i] }}</text>
       </g>
       <line :x1="scale.position(result.row[2])" y1="51" :x2="scale.position(result.row[2])" y2="90" stroke="var(--ui-text)" stroke-width="2" stroke-dasharray="4 3"/>
       <text :x="Math.max(100,Math.min(540,scale.position(result.row[2])))" y="46" text-anchor="middle">{{ t('milestoneLuck.medianMarker',{count:result.row[2]}) }}</text>
       <line :x1="scale.position(result.cost)" y1="55" :x2="scale.position(result.cost)" y2="80" stroke="var(--ui-action)" stroke-width="2"/>
       <circle :cx="scale.position(result.cost)" cy="80" r="6" fill="var(--ui-action)" stroke="var(--ui-input)" stroke-width="2"/>
       <text class="your-marker" :x="Math.max(100,Math.min(540,scale.position(result.cost)))" y="17" text-anchor="middle">{{ t('milestoneLuck.youMarker',{count:result.cost}) }}</text>
       <text x="32" y="241">0</text><text x="608" y="241" text-anchor="end">{{ scale.maximum }} · {{ t('character.pullsUnit') }}</text>
      </svg>
     </div>
     <p v-if="result.reference" class="reference-note">{{ t('milestoneLuck.referenceHelp') }}</p>
    </template>
    <p class="detail-note">{{ t('milestoneLuck.expected',{mean:result.row[7].toFixed(1),median:result.row[2]}) }}</p>
    <details><summary>{{ t('milestoneLuck.thresholds') }}</summary>
     <table class="luck-table"><thead><tr><th>{{ t('milestoneLuck.probability') }}</th><th>{{ t('character.pullsUnit') }}</th></tr></thead><tbody><tr v-for="(p,i) in chartPercentiles" :key="p"><td>{{ p }}%</td><td>{{ result.row[i] }}</td></tr></tbody></table>
    </details>
   </template>
   <p v-else>{{ t('milestoneLuck.missing',{dolls:level+1,weapons:rank}) }}</p>
  </div>
  <details class="explanation-panel"><summary>{{ t('milestoneLuck.sourceTitle') }}</summary><p>{{ t('milestoneLuck.source') }}</p></details>
 </section>
</template>
<style scoped>
.milestone-luck {margin:20px 0;padding:18px;border:1px solid var(--ui-border);border-radius:12px;background:var(--ui-input)}
h3 {font-size:18px;font-weight:650;margin-bottom:8px}
.milestone-toggle {margin:14px 0;flex-wrap:wrap}
.luck-milestones {display:grid;grid-template-columns:repeat(auto-fit,minmax(80px,1fr));gap:8px}
.luck-milestones button {padding:12px 6px;border:1px solid var(--ui-border);border-radius:8px;background:var(--ui-surface);color:var(--ui-text);cursor:pointer}
button[aria-pressed=true] {background:var(--ui-selected);border-color:var(--ui-action);box-shadow:inset 0 -3px var(--ui-action)}
button:hover {border-color:var(--ui-action)}
button:focus-visible {outline:2px solid var(--ui-focus);outline-offset:3px}
small {display:block;margin-top:6px;color:var(--ui-muted)}
.milestone-result {margin:16px 0;padding:16px;border-radius:8px;background:var(--ui-surface)}
.milestone-result p {margin-top:10px}
.milestone-badge {display:inline-block;padding:4px 10px;border:1px solid currentColor;border-radius:6px;font-weight:650}
.veryLucky,.lucky {color:var(--ui-success)}
.typical,.reference {color:var(--ui-muted)}
.unlucky {color:var(--ui-warning)}
.veryUnlucky {color:var(--ui-loss)}
.median-difference {font-size:calc(18px * var(--text-scale));font-weight:650}
.percentile-summary {display:flex;align-items:center;gap:24px;flex-wrap:wrap;padding:18px 0;margin-top:12px;border-top:1px solid var(--ui-border);border-bottom:1px solid var(--ui-border)}
.metric-label {display:block;font-weight:600;color:var(--ui-text)}
.percentile-value {display:flex;align-items:flex-start;gap:4px;margin:4px 0}
.percentile-value strong {font-size:calc(36px * var(--text-scale));line-height:1.2;font-variant-numeric:tabular-nums;color:var(--ui-text)}
.info-button {display:inline-flex;align-items:center;justify-content:center;width:32px;height:32px;border:1px solid var(--ui-border);border-radius:50%;background:var(--ui-input);color:var(--ui-action);cursor:pointer;font-size:18px}
.info-button[aria-expanded=true] {background:var(--ui-selected);border-color:var(--ui-action)}
.metric-direction {color:var(--ui-muted)}
.explanation-panel {margin-top:14px;padding:12px 16px;border-left:3px solid var(--ui-border);border-radius:0 8px 8px 0;background:var(--ui-input);color:var(--ui-muted);line-height:1.6}
.explanation-panel strong {color:var(--ui-text)}
.reference-note {padding:12px;border:1px solid var(--ui-border);border-radius:8px;color:var(--ui-muted)}
.comparison-graphic {overflow-x:auto;margin-top:16px;border:1px solid var(--ui-border);border-radius:8px;padding:12px 4px;background:var(--ui-input)}
.comparison-graphic:focus-visible {outline:2px solid var(--ui-focus);outline-offset:3px}
svg {width:100%;min-width:560px;display:block}
svg text {fill:var(--ui-muted);font-size:13px;font-variant-numeric:tabular-nums}
svg .your-marker {fill:var(--ui-action);font-weight:700;font-size:15px}
details {margin-top:12px} summary {cursor:pointer;color:var(--ui-action);padding:8px 0}
.luck-table {width:100%;text-align:left;border-collapse:collapse;font-variant-numeric:tabular-nums}
th,td {padding:6px 10px;border-bottom:1px solid var(--ui-border)}
</style>

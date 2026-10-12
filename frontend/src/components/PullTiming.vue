<script setup lang="ts">
import {computed} from 'vue'
import {useI18n} from 'vue-i18n'
import {timingEstimate} from '../utils/luck'
const props=defineProps<{poolType:number,count:number,expanded?:boolean}>()
const {t}=useI18n()
const estimate=computed(()=>timingEstimate(props.poolType,props.count))
const scale=computed(()=>{
 const e=estimate.value;if(!e)return null
 const threshold=(p:number)=>{for(let n=1;n<=e.hard;n++){if((timingEstimate(props.poolType,n)?.percentile || 0)>=p)return 100*n/e.hard}return 100}
 const a=threshold(25),b=threshold(75),c=threshold(90)
 return {background: `linear-gradient(to right,#427f6d 0% ${a}%,#687586 ${a}% ${b}%,#9c8050 ${b}% ${c}%,#a3673d ${c}% 100%)`}
})
const explanation=computed(()=>estimate.value ? t('luck.percentile',{percent:estimate.value.percentile.toFixed(1),count:props.count}) : t('luck.unavailable'))
</script>
<template>
 <span v-if="!expanded" class="timing-badge" :class="estimate ? `timing-${estimate.band}` : 'timing-unknown'" :title="explanation + ' ' + t('luck.assumption')">
  <span aria-hidden="true" class="timing-dot"></span>{{ estimate ? t(`luck.${estimate.band}`) : t('luck.unrated') }}<span class="timing-context">{{ t('luck.timing') }}</span>
 </span>
 <div v-else class="timing-detail">
  <div class="timing-heading"><span class="timing-badge" :class="estimate ? `timing-${estimate.band}` : 'timing-unknown'"><span aria-hidden="true" class="timing-dot"></span>{{ estimate ? t(`luck.${estimate.band}`) : t('luck.unrated') }}</span><span>{{ t('luck.estimated') }}</span></div>
  <template v-if="estimate">
   <p class="detail-note">{{ t('luck.pullScale',{count}) }}</p>
   <div class="timing-scale" :style="scale" role="img" :aria-label="explanation"><span class="timing-marker" :style="{left:`${100*count/estimate.hard}%`}"></span></div>
   <div class="timing-axis"><span>0</span><span>{{ t('guide.pulls',{count:estimate.hard}) }}</span></div>
  </template>
  <p>{{ explanation }}</p><p v-if="estimate" class="detail-note">{{ t('luck.softExplanation') }}</p>
 </div>
</template>

<script setup lang="ts">
import {computed} from 'vue'
import {useI18n} from 'vue-i18n'
import {ConfirmImport} from '../../wailsjs/go/handler/App'
import {useImportStore} from '../store/importStore'
import {useUserStore} from '../store/userStore'
import {usePoolStore} from '../store/poolStore'
import {useCaptureSecurityStore} from '../store/captureSecurityStore'
const security=useCaptureSecurityStore()
const flow=useImportStore(),users=useUserStore(),pools=usePoolStore()
const {t}=useI18n()
const mismatch=computed(()=>{const user=users.userList.find(item=>item.id===users.userId);return user && flow.preview && (String(user.uid)!==flow.preview.UID || user.server!==flow.preview.Server)})
const confirm=async()=>{
 if(!flow.preview || !flow.acknowledged || flow.busy)return
 const preview=flow.preview
 const removeAfterImport=security.removeAfterImport
 flow.confirming=false;flow.busy=true;flow.error='';flow.result='';flow.entries=[]
 try {
  const result=await ConfirmImport(preview.Ticket,preview.UID)
  const count=(result.DiffList || []).reduce((sum,diff)=>sum+diff.Count,0)
  flow.result=t('guide.imported',{count,server:t(`server.${preview.Server}`),uid:preview.UID})
  flow.lastSync=new Date().toLocaleString()
  flow.lastResult=flow.result
  // Keep removal instructions reachable while Windows asks for confirmation.
  flow.visible=true
  // Removal failure is reported separately; imported records remain successful.
  if(removeAfterImport)await security.remove()
  else await security.refresh()
  flow.visible=true
  await users.updateUserList()
  const account=users.userList.find(user=>user.server===preview.Server && String(user.uid)===preview.UID)
  if(account)await users.updateUserId(account.id)
  await pools.updatePoolInfo()
 }catch(err){flow.error=String(err);flow.addProgress({Stage:'failed'});flow.visible=true}
 finally{flow.busy=false;flow.preview=null;if(flow.result)flow.successVisible=true}
}
</script>
<template>
 <el-dialog v-model="flow.confirming" class="import-dialog" width="640" :title="t('importReview.title')" :close-on-click-modal="false">
  <div v-if="flow.preview" class="import-guide">
   <p>{{ t('importReview.intro') }}</p>
   <section class="guide-step"><h3>{{ t(`server.${flow.preview.Server}`) }}</h3><p class="review-uid">{{ flow.preview.IdentitySource==='hash'?t('importReview.localId'):t('guide.uid') }}: <strong>{{ flow.preview.UID }}</strong></p><p>{{ t(`importReview.sources.${flow.preview.IdentitySource}`) }}</p></section>
   <p v-if="flow.preview.IdentitySource==='webshop'" class="cost-caveat">{{ t('importReview.webshop') }}</p>
   <p v-if="mismatch" class="cost-caveat">{{ t('importReview.mismatch') }}</p>
   <label class="appearance-option"><input type="checkbox" v-model="flow.acknowledged"><span>{{ t('importReview.check') }}</span></label>
   <div class="editor-actions"><button class="action-button primary-action" :disabled="!flow.acknowledged" @click="confirm">{{ t('importReview.confirm') }}</button><button class="action-button" @click="flow.confirming=false;flow.preview=null">{{ t('character.cancel') }}</button></div>
  </div>
 </el-dialog>
</template>

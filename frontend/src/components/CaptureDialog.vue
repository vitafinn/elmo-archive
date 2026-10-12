<script setup lang="ts">
import {ref, watch, onBeforeUnmount} from 'vue'
import {StartCaptureOS, StopCaptureOS,PrepareImport,PrepareCaptureCertificate,InstallCaptureCertificate,FocusCertificatePrompt} from '../../wailsjs/go/handler/App'
import {util} from '../../wailsjs/go/models'
import {EventsOnce} from '../../wailsjs/runtime/runtime'
import {useImportStore} from '../store/importStore'
import ButtonSyncRecords from './ButtonSyncRecords.vue'
import ImportLog from './ImportLog.vue'
import CaptureSecurity from './CaptureSecurity.vue'
import {useCaptureSecurityStore} from '../store/captureSecurityStore'
const visible = defineModel({type:Boolean, required:true})
const flow = useImportStore()
const security=useCaptureSecurityStore()
const mode = ref('manual')
const uid = ref('')
const active = ref(false)
const starting = ref(false)
const stopping = ref(false)
const error = ref('')
const invalidUid = ref(false)
const certificate=ref<util.CaptureCertificate|null>(null)
const checkingCertificate=ref(false), installingCertificate=ref(false)
const setCleanup=(event:Event)=>security.setRemoveAfterImport((event.target as HTMLInputElement).checked)
let unsubscribe: (()=>void) | undefined
const cleanup = () => { unsubscribe?.(); unsubscribe = undefined }
const stop = async () => {
  stopping.value = true
  try { await StopCaptureOS(); active.value = false; cleanup() }
  catch (err) { error.value = String(err) }
  finally { stopping.value = false;void security.refresh() }
}
const start = async () => {
  invalidUid.value = mode.value === 'manual' && !/^[1-9]\d*$/.test(uid.value.trim())
  if (invalidUid.value || starting.value || active.value || checkingCertificate.value) return
  checkingCertificate.value=true;error.value=''
  try {
    const details=await PrepareCaptureCertificate()
    if(!details.Trusted){certificate.value=details;return}
    certificate.value=null
  } catch(err){error.value=String(err);return}
  finally{checkingCertificate.value=false;void security.refresh()}
  starting.value = true; error.value = ''; flow.error = ''; flow.result = ''; flow.captured = false
  flow.preview=null;flow.entries=[];flow.addProgress({Stage:'captureStarted'})
  cleanup()
  unsubscribe = EventsOnce('captureSuccess', async () => {
    flow.captured = true
    await stop()
    try{flow.preview=await PrepareImport('os')}
    catch(err){error.value=String(err)}
    flow.addProgress({Stage:'captureReady'})
  })
  try { await StartCaptureOS(mode.value === 'manual' ? uid.value.trim() : mode.value); active.value = !flow.captured }
  catch (err) { error.value = String(err); await stop() }
  finally { starting.value = false;void security.refresh() }
}
const installCertificate=async()=>{
 if(!certificate.value || installingCertificate.value)return
 installingCertificate.value=true;error.value=''
 try{await InstallCaptureCertificate(certificate.value.Thumbprint);await start()}
 catch(err){error.value=String(err)}
 finally{installingCertificate.value=false;void security.refresh()}
}
const beforeClose = async (done:()=>void) => {
  if (security.busy || flow.busy || starting.value || stopping.value || checkingCertificate.value || installingCertificate.value) return
  if (active.value) await stop()
  if (!active.value) { cleanup(); done() }
}
watch(visible,value=>{if(value)void security.refresh()})
onBeforeUnmount(cleanup)
</script>
<template>
  <el-dialog v-model="visible" class="import-dialog" width="680" :title="$t('guide.import')" :close-on-click-modal="false" :before-close="beforeClose">
    <div class="import-guide">
      <section v-if="flow.lastResult" class="sync-summary" role="status"><strong>✓ {{ $t('syncSuccess.previous') }}</strong><p>{{ flow.lastResult }}</p><small>{{ flow.lastSync }}</small><button class="text-link" @click="flow.successVisible=true">{{ $t('syncSuccess.details') }}</button></section>
      <p>{{ $t('guide.intro') }}</p>
      <div class="capture-flow">
        <section v-if="!active && !flow.captured && !certificate" class="guide-step">
          <h3>1 · {{ $t('guide.identify') }}</h3>
          <label class="guide-field">{{ $t('guide.method') }}
            <select v-model="mode" :disabled="starting">
              <option value="manual">{{ $t('guide.manual') }}</option>
              <option value="webshop">{{ $t('guide.webshop') }}</option>
            </select>
          </label>
          <template v-if="mode === 'manual'">
            <label class="guide-field">{{ $t('guide.uid') }}<input v-model="uid" inputmode="numeric" :disabled="starting" :aria-invalid="invalidUid" :placeholder="$t('guide.uidPlaceholder')"></label>
            <p>{{ $t('guide.uidHelp') }}</p>
            <p v-if="invalidUid" role="alert" class="text-red-300">{{ $t('guide.uidInvalid') }}</p>
          </template>
          <p v-else>{{ $t('guide.webshopHelp') }}</p>
          <el-button type="primary" :loading="starting || checkingCertificate" :disabled="flow.busy || flow.preparing || flow.confirming || security.busy" @click="start">{{ $t('guide.connect') }}</el-button>
          <details><summary>{{ $t('guide.connectionDetails') }}</summary><p class="mt-2">{{ $t('guide.connectionNote') }}</p><p>{{ $t('guide.sameMethod') }}</p></details>
        </section>
        <section v-if="certificate" class="guide-step certificate-review" aria-live="polite">
          <h3>{{ $t('certificate.title') }}</h3>
          <p>{{ $t('certificate.purpose') }}</p>
          <p class="certificate-permission">{{ $t('certificate.permission') }}</p>
          <div class="certificate-identity"><span>{{ $t('certificate.expectedName') }}</span><strong>{{ certificate.Name }}</strong></div>
          <label class="cleanup-choice"><input type="checkbox" :checked="security.removeAfterImport" :disabled="installingCertificate" @change="setCleanup"><span>{{ $t('captureSecurity.autoRemove') }}</span></label>
          <template v-if="installingCertificate">
            <p role="status"><strong>{{ $t('certificate.waiting') }}</strong></p><p>{{ $t('certificate.hidden') }}</p>
            <button class="action-button" @click="FocusCertificatePrompt()">{{ $t('certificate.showPrompt') }}</button>
          </template>
          <div v-else class="certificate-actions"><el-button type="primary" @click="installCertificate">{{ $t('certificate.install') }}</el-button><button class="action-button" @click="certificate=null;error=''">{{ $t('character.cancel') }}</button></div>
          <details class="certificate-more"><summary>{{ $t('certificate.detailsTitle') }}</summary>
            <p>{{ $t('certificate.warning') }}</p>
            <div class="certificate-identity"><span>{{ $t('certificate.thumbprint') }}</span><code>{{ certificate.Thumbprint }}</code></div>
            <p>{{ $t('certificate.compare') }}</p>
            <p v-if="certificate.Legacy">{{ $t('certificate.legacy') }}</p>
            <p>{{ $t('certificate.files') }}<br><code class="certificate-path">{{ certificate.Directory }}</code></p>
            <p>{{ $t('certificate.removal') }}</p>
          </details>
        </section>
        <section v-if="active || starting" class="guide-step" aria-live="polite">
          <template v-if="mode==='webshop'">
           <div class="capture-step" :class="{'is-complete':flow.webshopReady}" role="status"><strong>{{ flow.webshopReady?'✓':'1' }} · {{ $t(flow.webshopReady?'guide.webshopCaptured':'guide.webshopStep') }}</strong></div>
           <div v-if="flow.webshopReady" class="capture-step is-current" role="status"><strong>{{ flow.gameReady?'✓':'2' }} · {{ $t(flow.gameReady?'guide.historyCaptured':'guide.moveToGame') }}</strong><p v-if="!flow.gameReady">{{ $t('guide.gameStep') }}</p></div>
          </template>
          <template v-else><h3>2 · {{ $t('guide.openHistory') }}</h3><p>{{ $t('guide.gameStep') }}</p></template>
          <p>{{ $t('guide.waiting') }}</p>
          <el-button :disabled="starting" :loading="stopping" @click="stop">{{ $t('guide.cancel') }}</el-button>
        </section>
        <section v-if="flow.captured && !active" class="guide-step" aria-live="polite">
          <h3>✓ {{ $t('guide.connected') }}</h3><p>{{ $t(flow.lastResult?'syncSuccess.readyAgain':'guide.ready') }}</p>
          <p v-if="flow.preview" class="review-uid">{{ $t(`server.${flow.preview.Server}`) }} · {{ flow.preview.IdentitySource==='hash'?$t('importReview.localId'):$t('guide.uid') }}: <strong>{{ flow.preview.UID }}</strong></p>
          <p v-if="flow.preview">{{ $t(`importReview.sources.${flow.preview.IdentitySource}`) }}</p>
          <ButtonSyncRecords/>
          <button class="text-link" :disabled="flow.busy" @click="flow.captured = false">{{ $t('guide.reconnect') }}</button>
        </section>
      </div>
      <CaptureSecurity v-if="!certificate" :locked="flow.busy || checkingCertificate || installingCertificate || starting || stopping" @removed="active=false;certificate=null;cleanup()"/>
      <ImportLog/>
      <section v-if="error || flow.error" class="guide-error" role="alert">
        <strong>{{ $t('guide.failed') }}</strong><p>{{ $t('guide.recovery') }}</p>
        <details><summary>{{ $t('guide.technical') }}</summary><p class="break-all">{{ error || flow.error }}</p></details>
      </section>
      
    </div>
  </el-dialog>
</template>
<style scoped>
.certificate-identity {display:grid;gap:8px;padding:16px;margin:14px 0;background:var(--ui-input);border:1px solid var(--ui-border);border-radius:8px}
.certificate-identity span {color:var(--ui-muted)}
.certificate-identity code {font-size:calc(16px * var(--text-scale));font-variant-numeric:tabular-nums;overflow-wrap:anywhere;user-select:text;color:var(--ui-text)}
.certificate-permission {color:var(--ui-muted)}
.cleanup-choice {display:flex;align-items:flex-start;gap:10px;cursor:pointer;line-height:1.5}
.cleanup-choice input {margin-top:5px;width:18px;height:18px;accent-color:var(--ui-action);flex:none}
.certificate-more {margin-top:22px;border-top:1px solid var(--ui-border);padding-top:12px;color:var(--ui-muted)}
.certificate-more summary {color:var(--ui-action);cursor:pointer}
.certificate-more p {margin-top:12px}
.certificate-actions {display:flex;gap:12px;flex-wrap:wrap;margin-top:16px}
</style>

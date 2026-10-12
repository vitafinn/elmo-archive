<script setup lang="ts">
import {FocusCertificatePrompt} from '../../wailsjs/go/handler/App'
import {onMounted} from 'vue'
import {useCaptureSecurityStore} from '../store/captureSecurityStore'
const props=defineProps<{locked?:boolean}>()
const emit=defineEmits<{removed:[]}>()
const security=useCaptureSecurityStore()
onMounted(()=>void security.refresh())
const remove=async()=>{if(await security.remove())emit('removed')}
</script>
<template>
 <section class="capture-security" :aria-label="$t('captureSecurity.title')">
  <header><h3>{{ $t('captureSecurity.title') }}</h3><span class="security-status" :class="{enabled:security.status?.Active}">{{ $t(!security.status?'captureSecurity.checking':security.status.Active?'captureSecurity.captureOn':'captureSecurity.captureOff') }}</span></header>
  <div class="trust-control">
   <div><strong>{{ $t(!security.status?'captureSecurity.checking':security.status.Certificate?.Trusted?'captureSecurity.installed':'captureSecurity.notInstalled') }}</strong></div>
   <button v-if="security.status?.Certificate?.Trusted" class="action-button" :disabled="props.locked || security.busy" @click="remove">{{ $t(security.busy?'captureSecurity.removing':security.status.Active?'captureSecurity.stopRemove':'captureSecurity.remove') }}</button>
  </div>
  <label class="cleanup-choice"><input type="checkbox" :checked="security.removeAfterImport" :disabled="props.locked || security.busy" @change="security.setRemoveAfterImport(($event.target as HTMLInputElement).checked)"><span><strong>{{ $t('captureSecurity.autoRemove') }}</strong><small>{{ $t('captureSecurity.autoHelp') }}</small></span></label>
  <div v-if="security.busy" role="status"><p>{{ $t('captureSecurity.removalPrompt') }}</p><button class="action-button" @click="FocusCertificatePrompt()">{{ $t('certificate.showPrompt') }}</button></div>
  <p v-if="security.removed" class="security-success" role="status">{{ $t('captureSecurity.removed') }}</p>
  <p v-if="security.error" class="guide-error" role="alert">{{ $t('captureSecurity.failed') }} {{ security.error }}</p>
  <details><summary>{{ $t('captureSecurity.details') }}</summary>
   <p>{{ $t('captureSecurity.scope') }}</p><p>{{ $t('captureSecurity.syncHelp') }}</p>
   <p>{{ $t('captureSecurity.trustLimit') }}</p><p>{{ $t('captureSecurity.removalScope') }}</p>
   <p v-if="security.status?.Certificate"><strong>{{ security.status.Certificate.Name }}</strong><br><code>{{ security.status.Certificate.Thumbprint }}</code></p>
   <p v-if="security.status?.Certificate">{{ $t('certificate.files') }}<br><code>{{ security.status.Certificate.Directory }}</code></p>
   <p>{{ $t('captureSecurity.hosts') }}</p><ul><li v-for="host in security.status?.Hosts" :key="host"><code>{{ host }}</code></li></ul>
  </details>
 </section>
</template>
<style scoped>
.capture-security {display:flex;flex-direction:column;gap:16px;padding:18px;border:1px solid var(--ui-border);border-radius:12px;background:var(--ui-input);margin:0;color:var(--ui-text)}
header {display:flex;align-items:center;justify-content:space-between;gap:12px;flex-wrap:wrap}
h3 {font-weight:650;font-size:17px;margin:0}
.security-status {padding:4px 10px;border-radius:20px;border:1px solid var(--ui-border);color:var(--ui-muted);font-size:13px}
.security-status.enabled {color:var(--ui-warning);border-color:var(--ui-warning)}
.security-summary {margin:12px 0;color:var(--ui-muted)}
.trust-control {display:flex;align-items:center;gap:16px;flex-wrap:wrap;padding:14px;background:var(--ui-surface);border-radius:8px}
.trust-control>div {flex:1;min-width:min(220px,100%)}
.trust-control .action-button {white-space:normal;text-align:center;max-width:100%}
p,small {line-height:1.5}
.trust-control p {margin-top:6px;color:var(--ui-muted)}
.cleanup-choice {display:flex;align-items:flex-start;gap:12px;padding:0;cursor:pointer}
.cleanup-choice input {margin-top:5px;width:18px;height:18px;flex:none;accent-color:var(--ui-action)}
.cleanup-choice small {display:block;margin-top:4px;color:var(--ui-muted)}
.security-success {color:var(--ui-success)}
details {border-top:1px solid var(--ui-border);padding-top:10px;color:var(--ui-muted)}
summary {cursor:pointer;color:var(--ui-action);padding:5px 0}
details p {margin-top:12px} code {overflow-wrap:anywhere;user-select:text;font-size:13px}
ul {margin-top:8px;padding-left:20px;list-style:disc}
</style>

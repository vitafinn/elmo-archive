import {defineStore} from 'pinia'
import {ref} from 'vue'
import {GetCaptureSecurityStatus,RemoveCaptureCertificate} from '../../wailsjs/go/handler/App'
import {service} from '../../wailsjs/go/models'
export const useCaptureSecurityStore=defineStore('captureSecurity',()=>{
 const status=ref<service.CaptureSecurityStatus|null>(null)
 const busy=ref(false),error=ref(''),removed=ref(false)
 const removeAfterImport=ref(localStorage.getItem('removeCertificateAfterImport')!=='false')
 let revision=0
 const refresh=async()=>{
  const request=++revision
  try{const next=await GetCaptureSecurityStatus();if(request===revision){status.value=next;if(next.Certificate?.Trusted)removed.value=false}return true}
  catch(err){if(request===revision)error.value=String(err);return false}
 }
 const setRemoveAfterImport=(enabled:boolean)=>{removeAfterImport.value=enabled;localStorage.setItem('removeCertificateAfterImport',String(enabled))}
 const remove=async()=>{
  if(busy.value)return false
  busy.value=true;error.value='';removed.value=false
  try{
   if(!await refresh())return false
   if(!status.value?.Certificate?.Trusted)return true
   await RemoveCaptureCertificate(status.value.Certificate.Thumbprint)
   removed.value=true
   return true
  }catch(err){error.value=String(err);return false}
  finally{busy.value=false;await refresh()}
 }
 return {status,busy,error,removed,removeAfterImport,setRemoveAfterImport,refresh,remove}
})

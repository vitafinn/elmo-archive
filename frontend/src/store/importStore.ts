import {defineStore} from 'pinia'
import {ref} from 'vue'
import {handler} from '../../wailsjs/go/models'
export const useImportStore = defineStore('import', () => {
  const visible = ref(false)
  const busy = ref(false)
  const captured = ref(false)
  const webshopReady=ref(false),gameReady=ref(false)
  const error = ref('')
  const result = ref('')
  const lastSync = ref('')
  const lastResult=ref(''),successVisible=ref(false)
  const preview=ref<handler.ImportPreview|null>(null)
  const confirming=ref(false),preparing=ref(false),acknowledged=ref(false)
  const verbose=ref(localStorage.getItem('verboseImport')==='true')
  const entries=ref<{stage:string,pool:number,page:number,count:number,time:string}[]>([])
  const addProgress=(event:{Stage:string,PoolType?:number,Page?:number,Count?:number})=>{
    if(event.Stage==='captureStarted'){webshopReady.value=false;gameReady.value=false}
    if(event.Stage==='webshopCaptured')webshopReady.value=true
    if(event.Stage==='gameCaptured')gameReady.value=true
    entries.value.push({stage:event.Stage,pool:event.PoolType||0,page:event.Page||0,count:event.Count||0,time:new Date().toLocaleTimeString()})
    if(entries.value.length>300)entries.value.shift()
  }
  const setVerbose=(value:boolean)=>{verbose.value=value;localStorage.setItem('verboseImport',String(value))}
  return {visible, busy, captured, webshopReady,gameReady,error, result, lastSync,lastResult,successVisible,preview,confirming,preparing,acknowledged,verbose,entries,addProgress,setVerbose}
})

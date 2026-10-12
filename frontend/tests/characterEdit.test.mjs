import test from 'node:test'
import assert from 'node:assert/strict'
import fs from 'node:fs'
import ts from 'typescript'
import {parse} from '@vue/compiler-sfc'
const script=parse(fs.readFileSync(new URL('../src/components/CharacterDetails.vue',import.meta.url),'utf8')).descriptor.scriptSetup.content.replace(/^import .*$/gm,'')
const js=ts.transpileModule(script,{compilerOptions:{target:ts.ScriptTarget.ES2022,module:ts.ModuleKind.None}}).outputText
function harness(fail=false){
 const record={ID:'copy-2',ItemID:10,ItemType:20,Annotation:'auto'}
 const scroller={scrollTop:500},anchor={isConnected:true,getBoundingClientRect:()=>({top:200}),focus:()=>{anchor.focused=true}}
 const row={dataset:{copyId:'copy-2'},querySelector:()=>anchor}
 let resolve;const response=new Promise(r=>resolve=r)
 const deps={ref:value=>({value}),computed:fn=>({get value(){return fn()}}),watch:()=>{},nextTick:async()=>{},defineProps:()=>({record}),defineEmits:()=>()=>{},useUserStore:()=>({userId:1}),usePoolStore:()=>({loadError:false,updatePoolInfo:async()=>{}}),useI18n:()=>({t:key=>key,te:()=>true}),GetItemHistory:()=>response,GetDollWeaponHistory:async()=>null,SetSelectorOutcome:async()=>{if(fail)throw Error('save failed')},summarizeCopies:()=>({}),withFirstWeapon:()=>({}),outcomeKey:x=>x,oddsFor:()=> '50/50'}
 const api=new Function(...Object.keys(deps),js+';return {records,loading,content,startEdit,save,editing,savedCopy,failedCopy};')(...Object.values(deps))
 api.records.value=[record];api.content.value={closest:()=>scroller,querySelectorAll:()=>[row]};api.startEdit(record)
 return {api,record,resolve,anchor,scroller}
}
test('saving keeps details mounted and restores focus to the same copy',async()=>{
 const h=harness();const pending=h.api.save(h.record);await Promise.resolve();await Promise.resolve()
 assert.equal(h.api.loading.value,false)
 assert.equal(h.api.records.value[0],h.record)
 h.resolve([h.record]);await pending
 assert.equal(h.api.savedCopy.value,'copy-2');assert.equal(h.api.editing.value,null)
 assert.equal(h.scroller.scrollTop,500);assert.equal(h.anchor.focused,true)
})
test('failed save preserves editor and places error beside the copy',async()=>{
 const h=harness(true);await h.api.save(h.record)
 assert.equal(h.api.editing.value,'copy-2');assert.equal(h.api.failedCopy.value,'copy-2');assert.equal(h.api.savedCopy.value,null)
})

import test from 'node:test'
import assert from 'node:assert/strict'
import fs from 'node:fs'
import {parse} from '@vue/compiler-sfc'
import {compile} from '@vue/compiler-dom'
import * as Vue from 'vue'
import {renderToString} from '@vue/server-renderer'
const source=fs.readFileSync(new URL('../src/components/CaptureSecurity.vue',import.meta.url),'utf8')
const template=parse(source).descriptor.template.content.replace(' as HTMLInputElement','')
const render=new Function('Vue',compile(template,{mode:'function',prefixIdentifiers:true}).code)(Vue)
async function card(state={}) {
 const app=Vue.createSSRApp({render,setup:()=>({props:{},remove:()=>{},security:{status:null,busy:false,error:'',removed:false,removeAfterImport:true,...state}})})
 app.config.globalProperties.$t=key=>key
 return renderToString(app)
}
test('unknown certificate status is not reported as capture off',async()=>{
 const html=await card();assert.match(html,/captureSecurity.checking/);assert.doesNotMatch(html,/captureSecurity.captureOff/)
})
test('installed certificate offers removal and defaults to cleanup after success',async()=>{
 const html=await card({status:{Active:false,Certificate:{Trusted:true,Name:'Sample CA',Thumbprint:'SAMPLE'},Hosts:['example.test']}})
 assert.match(html,/captureSecurity.remove</);assert.match(html,/<input[^>]*checked/);assert.match(html,/Sample CA/);assert.match(html,/captureSecurity.syncHelp/)
})
test('active capture makes the stop effect explicit',async()=>{
 const html=await card({status:{Active:true,Certificate:{Trusted:true}}});assert.match(html,/captureSecurity.stopRemove/)
})
test('cleanup failure remains visible separately from successful removal',async()=>{
 const html=await card({error:'Removal failed',status:{Active:false,Certificate:{Trusted:true}}})
 assert.match(html,/role="alert"/);assert.match(html,/Removal failed/);assert.doesNotMatch(html,/captureSecurity.removed/)
})

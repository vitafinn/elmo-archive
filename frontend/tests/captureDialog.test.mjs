import test from 'node:test'
import assert from 'node:assert/strict'
import fs from 'node:fs'
import {parse} from '@vue/compiler-sfc'
import {compile} from '@vue/compiler-dom'
import * as Vue from 'vue'
import {renderToString} from '@vue/server-renderer'

const source=fs.readFileSync(new URL('../src/components/CaptureDialog.vue',import.meta.url),'utf8')
const template=parse(source).descriptor.template.content
const render=new Function('Vue',compile(template,{mode:'function',prefixIdentifiers:true}).code)(Vue)
async function dialog(state={}) {
 const app=Vue.createSSRApp({render,setup:()=>({
  security:{busy:false,removeAfterImport:true},setCleanup:()=>{},visible:true,active:false,starting:false,stopping:false,error:'',invalidUid:false,
  certificate:null,checkingCertificate:false,installingCertificate:false,installCertificate:()=>{},FocusCertificatePrompt:()=>{},mode:'manual',uid:'',beforeClose:()=>{},start:()=>{},stop:()=>{},
  flow:{captured:false,busy:false,preparing:false,confirming:false,error:'',result:''},...state
 })})
 app.config.globalProperties.$t=key=>key
 for(const name of ['el-dialog','el-button','ImportLog','ButtonSyncRecords','CaptureSecurity']) {
  app.component(name,{render(){return Vue.h('div',this.$slots.default?.())}})
 }
 const html=await renderToString(app)
 // Native template contents are inert in WebView2, even if SSR includes them.
 return html.replace(/<template\b[^>]*>[\s\S]*?<\/template>/g,'')
}
test('new import visibly offers UID, supported methods and connect action',async()=>{
 const html=await dialog()
 assert.match(html,/guide\.identify/)
 assert.match(html,/<input[^>]*inputmode="numeric"/)
 assert.match(html,/value="manual"/)
 assert.match(html,/value="webshop"/)
 assert.match(html,/guide\.connect/)
 assert.doesNotMatch(html,/value="(?:cn|hash)"/)
})
test('active capture visibly shows game instructions and cancel action',async()=>{
 const html=await dialog({active:true})
 assert.match(html,/guide\.gameStep/)
 assert.match(html,/guide\.cancel/)
 assert.doesNotMatch(html,/guide\.identify/)
})
test('completed capture visibly shows the ready state',async()=>{
 const html=await dialog({flow:{captured:true,busy:false,error:'',result:''}})
 assert.match(html,/guide\.connected/)
 assert.match(html,/guide\.ready/)
 assert.match(html,/guide\.reconnect/)
})


test('certificate review shows identity before opening Windows and explains hidden prompt',async()=>{
 const cert={Name:'Elmo Archive Local Capture CA',Thumbprint:'12345678 ABCDEF00',Trusted:false,Legacy:false}
 const review=await dialog({certificate:cert})
 assert.match(review,/Elmo Archive Local Capture CA/)
 assert.match(review,/12345678 ABCDEF00/)
 assert.match(review,/certificate.compare/)
 assert.match(review,/certificate.install/)
 const main=review.replace(/<details[\s\S]*?<\/details>/g,'')
 assert.match(main,/certificate.permission/)
 assert.match(main,/certificate.install/)
 assert.match(main,/<input[^>]*checked/)
 assert.doesNotMatch(main,/12345678 ABCDEF00|certificate.files/)
 assert.doesNotMatch(review,/guide.identify/)
 const pending=await dialog({certificate:cert,installingCertificate:true})
 assert.match(pending,/certificate.showPrompt/)
 assert.match(pending,/certificate.hidden/)
 assert.doesNotMatch(pending,/certificate.install</)
})

test('previous sync result stays above connection instructions',async()=>{
 const html=await dialog({flow:{lastResult:'Imported 12 records',lastSync:'Today',captured:true,busy:false,error:'',result:''}})
 assert.match(html,/syncSuccess.previous/)
 assert.match(html,/syncSuccess.readyAgain/)
 assert.ok(html.indexOf('Imported 12 records')<html.indexOf('guide.connected'))
 assert.doesNotMatch(html,/guide.ready</)
})

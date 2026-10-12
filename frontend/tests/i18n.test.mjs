import test from 'node:test'
import assert from 'node:assert/strict'
import fs from 'node:fs'
import ts from 'typescript'
async function locale(name) {
 let src=fs.readFileSync(new URL(`../src/i18n/${name}.ts`,import.meta.url),'utf8')
 const extra=src.match(/import extra from '(.+)'/)
 if(extra) {
  const part=fs.readFileSync(new URL(`../src/i18n/${extra[1]}.ts`,import.meta.url),'utf8').replace('export default','const extra =')
  src=src.replace(extra[0],part)
 }
 const js=ts.transpileModule(src,{compilerOptions:{module:ts.ModuleKind.ESNext}}).outputText
 return (await import('data:text/javascript;base64,'+Buffer.from(js).toString('base64'))).default
}
function flatten(obj,prefix='') {return Object.fromEntries(Object.entries(obj).flatMap(([key,value])=>typeof value==='string'?[[prefix+key,value]]:Object.entries(flatten(value,prefix+key+'.'))))}
const english=flatten(await locale('en'))
for(const name of ['zh-CN','zh-TW']) {
 test(`${name} has every English message and matching interpolation parameters`,async()=>{
  const translated=flatten(await locale(name))
  for(const [key,value] of Object.entries(english)) {
   assert.equal(typeof translated[key],'string',`missing ${key}`)
   const params=text=>[...text.matchAll(/\{(\w+)\}/g)].map(match=>match[1]).sort()
   assert.deepEqual(params(translated[key]),params(value),`parameters differ in ${key}`)
   // Publisher brand names remain unchanged in every language.
   if(key==='window.title') assert.equal(translated[key],'Elmo Archive')
   else if(key==='server.us') assert.equal(translated[key],'Darkwinter')
   else assert.match(translated[key],/[\u3400-\u9fff]/,`untranslated ${key}`)
  }
 })
}

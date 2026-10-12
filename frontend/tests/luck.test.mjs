import test from 'node:test'
import assert from 'node:assert/strict'
import fs from 'node:fs'
import ts from 'typescript'
const source=fs.readFileSync(new URL('../src/utils/luck.ts',import.meta.url),'utf8')
const js=ts.transpileModule(source,{compilerOptions:{module:ts.ModuleKind.ESNext}}).outputText
const {timingEstimate}=await import('data:text/javascript;base64,'+Buffer.from(js).toString('base64'))
test('base-rate region and hard-pity normalization',()=>{
 assert.ok(Math.abs(timingEstimate(3,58).percentile-100*(1-.994**58))<1e-9)
 assert.ok(Math.abs(timingEstimate(4,50).percentile-100*(1-.993**50))<1e-9)
 assert.equal(timingEstimate(3,80).percentile,100)
 assert.equal(timingEstimate(4,70).percentile,100)
})
test('CDF is monotonic and bounded across every supported category',()=>{
 for(const type of [1,3,6,-6,4,7,-7]){
  const hard=[4,7,-7].includes(type)?70:80
  let previous=0
  for(let n=1;n<=hard;n++){const result=timingEstimate(type,n);assert.ok(result.percentile>=previous && result.percentile<=100);previous=result.percentile}
 }
})
test('same rarity counter rules, different target guarantees; no guessed special pool rules',()=>{
 assert.deepEqual(timingEstimate(3,60),timingEstimate(-6,60))
 assert.equal(timingEstimate(-6,60).band,'typical')
 assert.equal(timingEstimate(3,10).band,'early')
 assert.equal(timingEstimate(3,80).band,'veryLate')
 for(const type of [5,8,9,10,42])assert.equal(timingEstimate(type,10),null)
 for(const count of [0,-1,1.5,81,NaN])assert.equal(timingEstimate(3,count),null)
 assert.equal(timingEstimate(4,71),null)
})

import test from 'node:test'
import assert from 'node:assert/strict'
import fs from 'node:fs'
import ts from 'typescript'
const source=fs.readFileSync(new URL('../src/utils/milestoneLuck.ts',import.meta.url),'utf8')
const js=ts.transpileModule(source,{compilerOptions:{module:ts.ModuleKind.ESNext}}).outputText
const {milestoneComparison:compare,milestoneChart,milestoneScale,estimateChartPercentile}=await import('data:text/javascript;base64,'+Buffer.from(js).toString('base64'))
const doll=(cost,extra={})=>({Count:cost,PoolType:3,CostBasis:'direct',...extra})
test('chart boundaries use ranges, including beyond the 99th percentile',()=>{
 assert.equal(compare([doll(34)],[],0,0).band,'veryLucky')
 assert.deepEqual([compare([doll(61)],[],0,0).lower,compare([doll(61)],[],0,0).upper],[10,25])
 assert.equal(compare([doll(62)],[],0,0).band,'typical')
 assert.equal(compare([doll(121)],[],0,0).band,'unlucky')
 assert.equal(compare([doll(135)],[],0,0).lower,99)
 for(const rows of Object.values(milestoneChart)) for(const row of rows) assert.ok(row.slice(0,7).every((v,i,a)=>!i || v>a[i-1]))
})
test('milestones count losses and only the required copies; missing combinations stay unavailable',()=>{
 const dolls=Array.from({length:8},()=>doll(10,{AcquisitionCount:20}))
 const weapons=Array.from({length:7},()=>doll(5,{PoolType:4,AcquisitionCount:15}))
 assert.equal(compare(dolls,weapons,6,6).cost,230)
 assert.equal(compare(dolls,weapons,0,1).cost,35)
 assert.equal(compare(dolls,weapons.slice(0,5),6,6),null)
 assert.equal(compare(dolls,weapons,6,2),null)
 assert.equal(compare([],[],0,0),null)
})
test('special rules, incomplete costs and estimated attribution remain explicit',()=>{
 assert.equal(compare([doll(60,{PoolType:-6})],[],0,0).reference,true)
 assert.equal(compare([doll(60,{IsMissing:true})],[],0,0).reference,true)
 assert.equal(compare([doll(60,{CostBasis:'incomplete'})],[],0,0).unknown,true)
 assert.equal(compare([doll(60,{CostBasis:'carried'})],[],0,0).estimated,true)
})

test('graphic uses a linear pull scale and includes costs beyond chart thresholds',()=>{
 const row=milestoneChart[0][6]
 const typical=milestoneScale(562,row)
 assert.equal(typical.delta,6)
 assert.equal(typical.position(0),32)
 assert.equal(typical.position(typical.maximum),608)
 assert.ok(Math.abs((typical.position(200)-typical.position(100))-(typical.position(100)-typical.position(0)))<1e-10)
 const beyond=milestoneScale(1000,row)
 assert.ok(beyond.position(1000)<608)
 assert.ok(beyond.maximum>1000)
 assert.equal(milestoneScale(row[2],row).delta,0)
 assert.equal(milestoneScale(10,row).delta,10-row[2])
})

test('percentiles interpolate within published values without inventing tails',()=>{
 const row=milestoneChart[0][6]
 assert.ok(Math.abs(estimateChartPercentile(562,row)-52.38095238095238)<1e-10)
 for(const rows of Object.values(milestoneChart)) for(const values of rows) {
  const percentiles=[10,25,50,75,90,95,99]
  values.slice(0,7).forEach((cost,i)=>assert.equal(estimateChartPercentile(cost,values),percentiles[i]))
  assert.equal(estimateChartPercentile(values[0]-1,values),null)
  assert.equal(estimateChartPercentile(values[6]+1,values),null)
 }
})

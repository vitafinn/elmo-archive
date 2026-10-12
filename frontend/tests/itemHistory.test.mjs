import test from 'node:test'
import assert from 'node:assert/strict'
import fs from 'node:fs'
import ts from 'typescript'
const source=fs.readFileSync(new URL('../src/utils/itemHistory.ts',import.meta.url),'utf8')
const js=ts.transpileModule(source,{compilerOptions:{module:ts.ModuleKind.ESNext}}).outputText
const {summarizeCopies,withFirstWeapon}=await import('data:text/javascript;base64,'+Buffer.from(js).toString('base64'))
test('loss-inclusive totals and cumulative vertebra milestones',()=>{
 const stats=summarizeCopies([{Count:60,AcquisitionCount:125,CostBasis:'linked'},{Count:10,AcquisitionCount:10,CostBasis:'direct'}])
 assert.equal(stats.pulls,70);assert.equal(stats.acquisition,135);assert.equal(stats.lossPulls,65)
 assert.equal(stats.acquisitionAverage,67.5);assert.deepEqual(stats.cumulative,[125,135])
 assert.equal(stats.cheapest,10);assert.equal(stats.highest,125);assert.equal(stats.duplicates,1)
})
test('unknown and estimated attribution remain visible; empty data is safe',()=>{
 const stats=summarizeCopies([{Count:60,AcquisitionCount:125,CostBasis:'estimatedOutcome'},{Count:20,AcquisitionCount:20,CostBasis:'unknown'}])
 assert.equal(stats.estimates,1);assert.equal(stats.unknown,1)
 const empty=summarizeCopies([]);assert.equal(empty.acquisition,0);assert.equal(empty.acquisitionAverage,0);assert.equal(empty.cheapest,0)
})

test('first weapon cost is included once; missing weapon history stays unknown',()=>{
 const dolls=[{Count:60,AcquisitionCount:125},{Count:10,AcquisitionCount:10}]
 const paired=withFirstWeapon(dolls,[{Count:50,AcquisitionCount:90},{Count:20,AcquisitionCount:20}])
 assert.equal(paired.weaponCost,90);assert.equal(paired.combined,225)
 assert.deepEqual(paired.cumulative,[215,225])
 const absent=withFirstWeapon(dolls,[])
 assert.equal(absent.combined,null);assert.deepEqual(absent.cumulative,[])
})


test('V6R6 requires seven doll and six weapon copies; excludes excess copies',()=>{
 const dolls=Array.from({length:8},()=>({Count:10,AcquisitionCount:20}))
 const weapons=Array.from({length:7},()=>({Count:5,AcquisitionCount:15}))
 assert.equal(withFirstWeapon(dolls,weapons).v6r6,230)
 assert.equal(withFirstWeapon(dolls.slice(0,6),weapons).v6r6,null)
 assert.equal(withFirstWeapon(dolls,weapons.slice(0,5)).v6r6,null)
})

test('sequence extrema stay separate from loss-inclusive acquisition cost',()=>{
 const stats=summarizeCopies([{Count:66,AcquisitionCount:66},{Count:59,AcquisitionCount:124,PreviousLoss:{Count:65,Name:'Faye'}}])
 assert.equal(stats.longestSequence,66);assert.equal(stats.highest,124)
 assert.equal(stats.highestRecord.PreviousLoss.Count,65);assert.equal(stats.highestRecord.Count,59)
})

test('V6R1 excludes doll copies beyond V6',()=>{
 const pair=withFirstWeapon(Array.from({length:8},()=>({Count:20})),[{Count:15}])
 assert.equal(pair.level,6);assert.equal(pair.dollCost,140);assert.equal(pair.combined,155)
})

test('recorded selector loss is included without marking its cost estimated',()=>{
 const stats=summarizeCopies([{Count:34,AcquisitionCount:98,CostBasis:'recordedLoss'}])
 assert.equal(stats.acquisition,98)
 assert.equal(stats.lossPulls,64)
 assert.equal(stats.estimates,0)
 assert.equal(stats.unknown,0)
})

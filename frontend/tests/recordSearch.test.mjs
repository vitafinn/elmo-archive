import test from 'node:test'
import assert from 'node:assert/strict'
import fs from 'node:fs'
import ts from 'typescript'
const source=fs.readFileSync(new URL('../src/utils/recordSearch.ts',import.meta.url),'utf8')
const js=ts.transpileModule(source,{compilerOptions:{module:ts.ModuleKind.ESNext}}).outputText
const {matchesRecord}=await import('data:text/javascript;base64,'+Buffer.from(js).toString('base64'))

test('English and Chinese queries find either display language through table names',()=>{
 assert.ok(matchesRecord('qiong','瓊玖',['Qiongjiu','瓊玖']))
 assert.ok(matchesRecord('瓊','Qiongjiu',['Qiongjiu','瓊玖']))
 assert.ok(matchesRecord('OTS 14','OTs-14'))
 assert.ok(matchesRecord('ｏｔｓ１４','OTs-14'))
 assert.equal(matchesRecord('Vepley','瓊玖',['Qiongjiu']),false)
})
test('missing optional names still allow displayed-name search and clearing',()=>{
 assert.ok(matchesRecord('tol','Tololo',null))
 assert.ok(matchesRecord('  ','Tololo',[]))
 assert.equal(matchesRecord('Qiongjiu','Tololo',undefined),false)
})

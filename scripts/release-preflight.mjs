import fs from 'node:fs'
import path from 'node:path'
import {execFileSync} from 'node:child_process'
const root=process.cwd()
const candidates=[...new Set(execFileSync('git',['ls-files','--cached','--others','--exclude-standard','-z'],{encoding:'utf8'}).split('\0').filter(Boolean))]
const allowed=/^(?:\.github\/|docs\/|scripts\/|db\/|handler\/|logger\/|logic\/|model\/|pb\/|request\/|service\/|util\/|frontend\/(?:src\/|tests\/|public\/|wailsjs\/|[^/]+$)|build\/(?:appicon\.png|windows\/)|(?:\.gitignore|LICENSE|README\.md|CONTRIBUTING\.md|SECURITY\.md|POOL_SUPPORT\.md|go\.mod|go\.sum|main\.go|wails\.json)$)/
const forbidden=/(?:^|\/)(?:node_modules|dist|db_export|__pycache__)(?:\/|$)|(?:\.db(?:-|$)|\.log(?:\.|$)|\.key$|\.crt$|\.pem$|\.pfx$|\.exe$|\.zip$|gacha_(?:auth_details|session_info)\.json$)/i
const files=candidates.filter(file=>allowed.test(file)&&fs.existsSync(file))
// Reject private files even when they are outside the export allowlist.
const failures=candidates.filter(file=>fs.existsSync(file)&&forbidden.test(file))
for(const file of files) {
 if(forbidden.test(file)) {failures.push(file);continue}
 if(fs.lstatSync(file).isSymbolicLink()) {failures.push(file);continue}
 const bytes=fs.readFileSync(file)
 if(bytes.includes(0)) continue
 const text=bytes.toString('utf8')
 if(/-----BEGIN (?:RSA |EC |OPENSSH )?PRIVATE KEY-----/.test(text)||/\b(?:ghp_|github_pat_)[A-Za-z0-9_]{30,}\b/.test(text)||/\beyJ[A-Za-z0-9_-]{25,}\.[A-Za-z0-9_-]{25,}\.[A-Za-z0-9_-]{20,}/.test(text)) failures.push(file)
}
if(failures.length) {console.error('Review possible private files (contents are not printed):\n'+failures.join('\n'));process.exit(1)}
if(!files.includes('frontend/package-lock.json')) throw Error('A frontend lockfile is required')
console.log(`Preflight passed for ${files.length} source files. This checks common credential patterns, not every possible secret.`)
if(process.argv[2]==='--export') {
 const dest=path.resolve(process.argv[3]||'')
 const releaseRoot=path.join(root,'build','release')+path.sep
 if(!dest.startsWith(releaseRoot)||fs.existsSync(dest)) throw Error('Use a new directory inside build/release')
 for(const file of files) {const target=path.join(dest,file);fs.mkdirSync(path.dirname(target),{recursive:true});fs.copyFileSync(file,target)}
 console.log('Source snapshot exported without local data or Git history.')
}

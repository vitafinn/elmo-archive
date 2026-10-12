import {handler} from '../../wailsjs/go/models'
export function withFirstWeapon(dolls: handler.DisplayRecord[], weapons: handler.DisplayRecord[]) {
  const dollTotals = summarizeCopies(dolls)
  const first = weapons[0]
  const weaponCost = first ? first.AcquisitionCount || first.Count : null
  const weaponTotals = summarizeCopies(weapons)
  const level=Math.min(6,Math.max(0,dolls.length-1))
  const dollCost=dollTotals.cumulative[level] || 0
  const v6r6 = dolls.length >= 7 && weapons.length >= 6 ? dollTotals.cumulative[6] + weaponTotals.cumulative[5] : null
  return {first, weaponCost, level, dollCost, combined:weaponCost === null ? null : dollCost + weaponCost,
    v6r6,
    cumulative:weaponCost === null ? [] : dollTotals.cumulative.map(cost=>cost + weaponCost)}
}
export function summarizeCopies(records: handler.DisplayRecord[]) {
  const pulls = records.reduce((sum, record) => sum + record.Count, 0)
  const costs = records.map(record => record.AcquisitionCount || record.Count)
  const acquisition = costs.reduce((sum, cost) => sum + cost, 0)
  let running = 0
  const cumulative = costs.map(cost => running += cost)
  return {copies:records.length, duplicates:Math.max(0,records.length-1), pulls, average:records.length ? pulls/records.length : 0,
    acquisition, lossPulls:acquisition-pulls, acquisitionAverage:records.length ? acquisition/records.length : 0,
    cheapest: costs.length ? Math.min(...costs) : 0, highest: costs.length ? Math.max(...costs) : 0, cumulative,
    shortestSequence: records.length ? Math.min(...records.map(record=>record.Count)) : 0,
    longestSequence: records.length ? Math.max(...records.map(record=>record.Count)) : 0,
    highestRecord: records[costs.indexOf(Math.max(...costs))],
    estimates: records.filter(record=>['estimatedOutcome','carried'].includes(record.CostBasis)).length,
    unknown: records.filter(record=>['unknown','incomplete'].includes(record.CostBasis)).length}
}
export function outcomeKey(outcome: string) {
  return outcome === 'win' ? 'guide.won' : outcome === 'loss' ? 'guide.lost' : `selector.outcome.${outcome || 'unknown'}`
}
export function oddsFor(poolType: number) { return [4,7,-7].includes(poolType) ? '75/25' : '50/50' }

// Approximate linear soft-pity model, not official per-pull rate data.
// This models the next five-star only; target wins and guarantees are separate.
export function timingEstimate(poolType: number, count: number) {
  const config = [1,3,6,-6].includes(poolType) ? {base:.006,soft:58,hard:80}
    : [4,7,-7].includes(poolType) ? {base:.007,soft:50,hard:70} : null
  if (!config || !Number.isInteger(count) || count<1 || count>config.hard) return null
  let survival=1
  for(let pull=1;pull<=count;pull++) {
    const rate=pull<=config.soft ? config.base : Math.min(1,config.base+(pull-config.soft)*(1-config.base)/(config.hard-config.soft))
    survival*=1-rate
  }
  const percentile=100*(1-survival)
  const band=percentile<=25?'early':percentile<=75?'typical':percentile<=90?'late':'veryLate'
  return {...config,percentile,band}
}

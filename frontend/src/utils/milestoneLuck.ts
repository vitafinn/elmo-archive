// Transcribed from the user-supplied fortification chart (一棵平衡树;
// translation: DiggusDBickus). Quantiles are rounded pull counts, not a full CDF.
export const chartPercentiles = [10,25,50,75,90,95,99]
export const milestoneChart: Record<number, number[][]> = {
 0: [[34,61,67,120,127,130,134,79.3299],[93,125,158,191,226,249,259,158.66],[158,194,241,283,317,337,378,237.99],[226,268,317,368,408,436,477,317.32],[296,342,396,449,500,524,574,396.65],[366,417,476,535,586,618,673,475.98],[436,492,556,619,675,707,767,555.309]],
 1: [[78,112,129,174,187,211,239,136.636],[142,179,216,251,295,308,351,215.965],[209,249,297,341,379,407,443,295.295],[278,323,374,427,473,498,549,374.625],[347,397,453,510,561,591,643,453.955],[418,472,533,595,649,682,741,533.285],[489,547,613,679,737,771,834,612.615]],
 6: [[329,372,421,473,519,550,604,423.164],[398,446,501,558,608,641,698,502.494],[467,521,581,642,698,731,791,581.823],[538,595,660,726,785,821,887,661.153],[609,671,740,810,873,910,980,740.483],[681,746,819,893,959,995,1073,819.813],[753,822,899,976,1046,1081,1165,899.143]]
}
type Copy = {Count:number; AcquisitionCount?:number; PoolType:number; CostBasis?:string; IsMissing?:boolean}
export function milestoneScale(cost:number, row:number[]) {
 const maximum=Math.ceil(Math.max(cost,row[6])*1.08)
 const position=(pulls:number)=>32+576*pulls/maximum
 return {maximum,position,delta:cost-row[2]}
}
// Interpolate only between published thresholds. The chart gives no tail shape.
export function estimateChartPercentile(cost:number,row:number[]) {
 if(cost<row[0] || cost>row[6]) return null
 const index=row.slice(0,7).findIndex(limit=>cost<=limit)
 if(index===0) return chartPercentiles[0]
 return chartPercentiles[index-1]+(cost-row[index-1])/(row[index]-row[index-1])*(chartPercentiles[index]-chartPercentiles[index-1])
}
export function milestoneComparison(dolls:Copy[], weapons:Copy[], v:number, r:number) {
 const row=milestoneChart[r]?.[v]
 if(!row || dolls.length<v+1 || weapons.length<r) return null
 const copies=[...dolls.slice(0,v+1),...weapons.slice(0,r)]
 const cost=copies.reduce((sum,copy)=>sum+(copy.AcquisitionCount || copy.Count),0)
 const unknown=copies.some(copy=>['unknown','incomplete'].includes(copy.CostBasis || ''))
 const reference=copies.some(copy=>![3,4].includes(copy.PoolType) || copy.IsMissing)
 const estimated=copies.some(copy=>['estimatedOutcome','carried'].includes(copy.CostBasis || ''))
 const index=row.slice(0,7).findIndex(limit=>cost<=limit)
 const lower=index===-1 ? 99 : index===0 ? 0 : chartPercentiles[index-1]
 const upper=index===-1 ? 100 : chartPercentiles[index]
 const band=cost<=row[0]?'veryLucky':cost<=row[1]?'lucky':cost<=row[3]?'typical':cost<=row[4]?'unlucky':'veryUnlucky'
 return {cost,row,lower,upper,band,unknown,reference,estimated,percentile:estimateChartPercentile(cost,row)}
}

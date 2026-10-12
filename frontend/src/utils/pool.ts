// Only known pool mechanics get pity colouring; new categories remain neutral.
export function poolRecordColor(type: number, count: number): string {
    if ([1, 3, 6, -6].includes(type)) {
        return count <= 58 ? 'bg-green-500' : count < 66 ? 'bg-cyan-400' : 'bg-red-600'
    }
    if ([4, 7, -7].includes(type)) {
        return count <= 50 ? 'bg-green-500' : count < 58 ? 'bg-cyan-400' : 'bg-red-600'
    }
    if (type === 5) return count < 50 ? 'bg-green-500' : 'bg-red-600'
    if (type === 8) return count <= 100 ? 'bg-green-500' : count <= 400 ? 'bg-cyan-400' : 'bg-red-600'
    return 'bg-slate-600'
}

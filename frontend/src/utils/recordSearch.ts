// Match spelling, not UI language. Names come from the installed game tables.
const normalize = (text: string) => text.normalize('NFKC').toLowerCase()
  .replace(/[\s\p{P}\p{M}]+/gu, '')

export function matchesRecord(query: string, name: string, aliases?: string[] | null): boolean {
  const needle = normalize(query)
  return !needle || [name, ...(aliases || [])].some(value => normalize(value).includes(needle))
}

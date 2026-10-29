// Date helpers. All ranges are YYYY-MM-DD strings in Asia/Jakarta (WIB),
// matching how the backend interprets `from` / `to`.

const TZ = 'Asia/Jakarta'

// todayWIB returns today's date in WIB as YYYY-MM-DD.
export function todayWIB(): string {
  // en-CA formats as YYYY-MM-DD
  return new Intl.DateTimeFormat('en-CA', { timeZone: TZ }).format(new Date())
}

// addDays shifts a YYYY-MM-DD date by n days (calendar math, no timezone).
export function addDays(ymd: string, n: number): string {
  const d = new Date(`${ymd}T00:00:00Z`)
  d.setUTCDate(d.getUTCDate() + n)
  return d.toISOString().slice(0, 10)
}

export type Preset = 'last7' | 'last30' | 'thisMonth' | 'all' | 'custom'

export const PRESETS: { value: Preset; label: string }[] = [
  { value: 'last7', label: 'Last 7 days' },
  { value: 'last30', label: 'Last 30 days' },
  { value: 'thisMonth', label: 'This month' },
  { value: 'all', label: 'All time' },
]

export function presetRange(p: Preset): { from: string; to: string } {
  const today = todayWIB()
  switch (p) {
    case 'last7':
      return { from: addDays(today, -6), to: today }
    case 'last30':
      return { from: addDays(today, -29), to: today }
    case 'thisMonth':
      return { from: today.slice(0, 8) + '01', to: today }
    default:
      return { from: '', to: '' }
  }
}

// formatDate shows an ISO timestamp as a WIB date, e.g. "29 Oct 2026".
export function formatDate(iso: string | null): string {
  if (!iso) return '—'
  return new Intl.DateTimeFormat('en-GB', { timeZone: TZ, day: 'numeric', month: 'short', year: 'numeric' }).format(
    new Date(iso),
  )
}

// formatDateTime shows an ISO timestamp in WIB with time, e.g. "29 Oct 2026, 21:30 WIB".
export function formatDateTime(iso: string): string {
  const s = new Intl.DateTimeFormat('en-GB', {
    timeZone: TZ,
    day: 'numeric',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  }).format(new Date(iso))
  return `${s} WIB`
}

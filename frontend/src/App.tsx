import { useEffect, useState } from 'react'
import { fetchLeads, fetchMetrics, type Lead, type LeadsPage, type Metrics } from './api'
import { formatDateTime, presetRange, type Preset } from './dates'
import KpiCards, { KpiSkeleton } from './components/KpiCards'
import LeadsTable, { TableSkeleton } from './components/LeadsTable'
import Sidebar from './components/Sidebar'
import StatusChart from './components/StatusChart'
import Toolbar from './components/Toolbar'
import TopBar from './components/TopBar'
import { DownloadIcon, RefreshIcon } from './components/icons'

const DEFAULT_PRESET: Preset = 'last30'

// App owns all page state (filters, page number, fetched data) and passes it
// down to components as props. Components never fetch on their own.
export default function App() {
  const [menuOpen, setMenuOpen] = useState(false)

  // Filters
  const [preset, setPreset] = useState<Preset>(DEFAULT_PRESET)
  const [range, setRange] = useState(() => presetRange(DEFAULT_PRESET))
  const [search, setSearch] = useState('')
  const [debouncedSearch, setDebouncedSearch] = useState('')
  const [page, setPage] = useState(1)
  const [reloadKey, setReloadKey] = useState(0) // bump to re-fetch

  // Fetched data
  const [metrics, setMetrics] = useState<Metrics | null>(null)
  const [leads, setLeads] = useState<LeadsPage | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(false)

  // Wait 300ms after the last keystroke before searching, so we don't call the API on every key.
  useEffect(() => {
    const t = setTimeout(() => setDebouncedSearch(search.trim()), 300)
    return () => clearTimeout(t)
  }, [search])

  // Any filter change goes back to page 1.
  useEffect(() => setPage(1), [range, debouncedSearch])

  // Fetch both endpoints whenever a filter, the page, or reloadKey changes.
  useEffect(() => {
    let cancelled = false // ignore late responses from an older request
    setLoading(true)
    setError(false)
    Promise.all([fetchMetrics(range), fetchLeads({ ...range, q: debouncedSearch, page })])
      .then(([m, l]) => {
        if (cancelled) return
        setMetrics(m)
        setLeads(l)
      })
      .catch(() => !cancelled && setError(true))
      .finally(() => !cancelled && setLoading(false))
    return () => {
      cancelled = true
    }
  }, [range, debouncedSearch, page, reloadKey])

  const choosePreset = (p: Preset) => {
    setPreset(p)
    setRange(presetRange(p))
  }
  const setFrom = (from: string) => {
    setPreset('custom')
    setRange((r) => ({ ...r, from }))
  }
  const setTo = (to: string) => {
    setPreset('custom')
    setRange((r) => ({ ...r, to }))
  }
  const clearAll = () => {
    setPreset('all')
    setRange(presetRange('all'))
    setSearch('')
  }

  const firstLoad = loading && !metrics
  const noApplications = metrics !== null && metrics.totalApplications === 0

  return (
    <div className="min-h-screen md:pl-16">
      <Sidebar open={menuOpen} onClose={() => setMenuOpen(false)} />
      <TopBar onMenu={() => setMenuOpen(true)} />

      <main className="mx-auto flex max-w-7xl flex-col gap-4 px-4 py-5 md:px-8 md:py-6">
        {/* Title row */}
        <div className="flex flex-wrap items-start justify-between gap-3">
          <div>
            <div className="flex flex-wrap items-center gap-2">
              <h1 className="text-2xl font-bold md:text-3xl">Leads Dashboard</h1>
              <span className="rounded-full border border-amber-300 bg-amber-50 px-2 py-0.5 text-xs font-medium text-amber-700">
                Prototype · mock data
              </span>
            </div>
            <p className="mt-1 text-xs text-muted">
              {metrics ? `Last updated ${formatDateTime(metrics.lastUpdated)}` : ' '}
            </p>
          </div>
          <div className="flex items-center gap-2">
            <button
              onClick={() => leads && exportCSV(leads.data)}
              disabled={!leads?.data.length}
              title="Export this page as CSV"
              aria-label="Export CSV"
              className="flex h-9 w-9 items-center justify-center rounded-md border border-line bg-card text-muted hover:text-ink disabled:opacity-40"
            >
              <DownloadIcon className="h-5 w-5" />
            </button>
            <button
              onClick={() => setReloadKey((k) => k + 1)}
              className="flex h-9 items-center gap-2 rounded-md bg-primary px-4 text-sm font-medium text-white hover:bg-primary/90"
            >
              <RefreshIcon className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
              Refresh
            </button>
          </div>
        </div>

        <Toolbar
          preset={preset}
          from={range.from}
          to={range.to}
          search={search}
          onPreset={choosePreset}
          onFrom={setFrom}
          onTo={setTo}
          onSearch={setSearch}
          onClear={clearAll}
        />

        {error ? (
          <div className="rounded-lg border border-red-200 bg-red-50 p-6 text-center">
            <p className="font-medium text-red-700">Couldn't load data. Try again.</p>
            <button
              onClick={() => setReloadKey((k) => k + 1)}
              className="mt-3 rounded-md bg-primary px-4 py-2 text-sm font-medium text-white hover:bg-primary/90"
            >
              Try again
            </button>
          </div>
        ) : firstLoad ? (
          <>
            <KpiSkeleton />
            <TableSkeleton />
          </>
        ) : (
          // Dim (not hide) old numbers while a new request is in flight.
          <div className={`flex flex-col gap-4 transition-opacity ${loading ? 'opacity-60' : ''}`}>
            {metrics && <KpiCards m={metrics} />}
            {noApplications ? (
              <div className="rounded-lg border border-line bg-card p-10 text-center text-muted">
                No applications in this date range.
              </div>
            ) : (
              metrics && <StatusChart byStatus={metrics.byStatus} />
            )}
            {leads && <LeadsTable page={leads} onPage={setPage} />}
          </div>
        )}
      </main>
    </div>
  )
}

// exportCSV downloads the leads currently shown (one page) as a CSV file.
function exportCSV(rows: Lead[]) {
  const cols: (keyof Lead)[] = ['id', 'name', 'programme', 'homeLC', 'status', 'appliedAt', 'approvedAt', 'realizedAt']
  const esc = (v: unknown) => `"${String(v ?? '').replace(/"/g, '""')}"`
  const csv = [cols.join(','), ...rows.map((r) => cols.map((c) => esc(r[c])).join(','))].join('\n')
  const url = URL.createObjectURL(new Blob([csv], { type: 'text/csv' }))
  const a = document.createElement('a')
  a.href = url
  a.download = 'leads.csv'
  a.click()
  URL.revokeObjectURL(url)
}

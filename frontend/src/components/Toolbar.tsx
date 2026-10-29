import { PRESETS, type Preset } from '../dates'
import { SearchIcon } from './icons'

interface Props {
  preset: Preset
  from: string
  to: string
  search: string
  onPreset: (p: Preset) => void
  onFrom: (v: string) => void
  onTo: (v: string) => void
  onSearch: (v: string) => void
  onClear: () => void
}

const input =
  'h-9 rounded-md border border-line bg-card px-3 text-sm outline-none focus:border-primary focus:ring-2 focus:ring-primary/20'

// Toolbar only reports what the user picked; App owns the values (a "controlled" component).
export default function Toolbar(p: Props) {
  return (
    <div className="flex flex-col gap-3 rounded-lg border border-line bg-card p-3 sm:flex-row sm:flex-wrap sm:items-end">
      <label className="flex flex-col gap-1 text-xs font-medium text-muted">
        Date range
        <select className={input} value={p.preset} onChange={(e) => p.onPreset(e.target.value as Preset)}>
          {PRESETS.map((o) => (
            <option key={o.value} value={o.value}>
              {o.label}
            </option>
          ))}
          <option value="custom" disabled hidden>
            Custom
          </option>
        </select>
      </label>

      <div className="grid grid-cols-2 gap-3 sm:flex">
        <label className="flex flex-col gap-1 text-xs font-medium text-muted">
          Applied from
          <input type="date" className={input} value={p.from} max={p.to || undefined} onChange={(e) => p.onFrom(e.target.value)} />
        </label>
        <label className="flex flex-col gap-1 text-xs font-medium text-muted">
          Applied to
          <input type="date" className={input} value={p.to} min={p.from || undefined} onChange={(e) => p.onTo(e.target.value)} />
        </label>
      </div>

      <label className="flex flex-1 flex-col gap-1 text-xs font-medium text-muted sm:min-w-48">
        Search
        <span className="relative">
          <SearchIcon className="pointer-events-none absolute top-1/2 left-2.5 h-4 w-4 -translate-y-1/2 text-muted" />
          <input
            type="search"
            placeholder="Name or ID"
            className={`${input} w-full pl-8`}
            value={p.search}
            onChange={(e) => p.onSearch(e.target.value)}
          />
        </span>
      </label>

      <button onClick={p.onClear} className="h-9 self-end px-2 text-sm font-medium text-primary hover:underline sm:ml-auto">
        Clear all
      </button>
    </div>
  )
}

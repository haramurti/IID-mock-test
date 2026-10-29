import type { Metrics } from '../api'
import { CheckIcon, InfoIcon, PercentIcon, PlaneIcon, UserIcon } from './icons'

// Four KPI numbers. All values come straight from /api/metrics.
export default function KpiCards({ m }: { m: Metrics }) {
  const rate = m.conversionRate === null ? '—' : `${m.conversionRate.toFixed(2)}%`
  const cards = [
    { label: 'Total Applications', value: m.totalApplications.toLocaleString(), Icon: UserIcon, color: 'text-primary' },
    { label: 'Approved', value: m.approved.toLocaleString(), Icon: CheckIcon, color: 'text-green-600' },
    { label: 'Realized', value: m.realized.toLocaleString(), Icon: PlaneIcon, color: 'text-teal-600' },
    {
      label: 'Conversion Rate',
      value: rate,
      Icon: PercentIcon,
      color: 'text-indigo-600',
      hint: 'Approved ÷ Applied × 100. Approved includes realized, finished and completed.',
    },
  ]

  return (
    <div className="grid grid-cols-2 gap-3 lg:grid-cols-4 lg:gap-4">
      {cards.map(({ label, value, Icon, color, hint }) => (
        <div key={label} className="rounded-lg border border-line bg-card p-4">
          <div className="flex items-center gap-1.5 text-xs font-medium text-muted sm:text-sm">
            {label}
            {hint && (
              // "group" lets the tooltip show on hover/focus of this icon only
              <span className="group relative" tabIndex={0} aria-label={hint}>
                <InfoIcon className="h-3.5 w-3.5 cursor-help" />
                <span
                  role="tooltip"
                  className="invisible absolute right-0 bottom-full z-10 mb-2 w-52 rounded-md bg-ink px-3 py-2 text-xs font-normal text-white opacity-0 shadow-lg transition group-hover:visible group-hover:opacity-100 group-focus:visible group-focus:opacity-100 sm:right-auto sm:left-1/2 sm:-translate-x-1/2"
                >
                  {hint}
                </span>
              </span>
            )}
          </div>
          <div className="mt-2 flex items-center gap-2">
            <Icon className={`h-5 w-5 shrink-0 ${color}`} />
            <span className="text-2xl font-bold sm:text-3xl">{value}</span>
          </div>
        </div>
      ))}
    </div>
  )
}

export function KpiSkeleton() {
  return (
    <div className="grid grid-cols-2 gap-3 lg:grid-cols-4 lg:gap-4">
      {[0, 1, 2, 3].map((i) => (
        <div key={i} className="h-[92px] animate-pulse rounded-lg border border-line bg-card p-4">
          <div className="h-3 w-24 rounded bg-gray-200" />
          <div className="mt-4 h-6 w-16 rounded bg-gray-200" />
        </div>
      ))}
    </div>
  )
}

import type { LeadsPage } from '../api'
import { formatDate } from '../dates'
import StatusBadge from './StatusBadge'

interface Props {
  page: LeadsPage
  onPage: (p: number) => void
}

// Leads list: a table on desktop, stacked cards on mobile. One page at a time.
export default function LeadsTable({ page, onPage }: Props) {
  const { data, paging } = page
  const first = paging.totalItems === 0 ? 0 : (paging.page - 1) * paging.perPage + 1
  const last = first === 0 ? 0 : first + data.length - 1

  return (
    <section className="rounded-lg border border-line bg-card">
      <div className="flex items-center justify-between border-b border-line px-4 py-3">
        <h2 className="font-semibold">Leads</h2>
        <span className="text-sm text-muted">{paging.totalItems.toLocaleString()} total</span>
      </div>

      {data.length === 0 ? (
        <p className="px-4 py-10 text-center text-sm text-muted">No leads match these filters.</p>
      ) : (
        <>
          {/* Desktop table */}
          <table className="hidden w-full text-sm md:table">
            <thead className="text-left text-xs text-muted uppercase">
              <tr className="border-b border-line">
                <th className="px-4 py-3 font-semibold">ID</th>
                <th className="px-4 py-3 font-semibold">Name</th>
                <th className="px-4 py-3 font-semibold">Programme</th>
                <th className="px-4 py-3 font-semibold">Home LC</th>
                <th className="px-4 py-3 font-semibold">Status</th>
                <th className="px-4 py-3 font-semibold">Applied ↓</th>
              </tr>
            </thead>
            <tbody>
              {data.map((l) => (
                <tr key={l.id} className="border-b border-line last:border-0 hover:bg-page">
                  <td className="px-4 py-3 font-medium text-primary">{l.id}</td>
                  <td className="px-4 py-3">{l.name}</td>
                  <td className="px-4 py-3">{l.programme}</td>
                  <td className="px-4 py-3 text-muted">{l.homeLC}</td>
                  <td className="px-4 py-3">
                    <StatusBadge status={l.status} />
                  </td>
                  <td className="px-4 py-3 whitespace-nowrap">{formatDate(l.appliedAt)}</td>
                </tr>
              ))}
            </tbody>
          </table>

          {/* Mobile cards */}
          <ul className="divide-y divide-line md:hidden">
            {data.map((l) => (
              <li key={l.id} className="px-4 py-3">
                <div className="flex items-start justify-between gap-2">
                  <div>
                    <div className="font-medium">{l.name}</div>
                    <div className="text-xs text-muted">
                      <span className="font-medium text-primary">{l.id}</span> · {l.programme} · {l.homeLC}
                    </div>
                  </div>
                  <StatusBadge status={l.status} />
                </div>
                <div className="mt-1 text-xs text-muted">Applied {formatDate(l.appliedAt)}</div>
              </li>
            ))}
          </ul>
        </>
      )}

      {paging.totalPages > 1 && (
        <div className="flex items-center justify-between gap-2 border-t border-line px-4 py-3 text-sm">
          <span className="text-muted">
            {first}–{last} of {paging.totalItems}
          </span>
          <div className="flex items-center gap-2">
            <button
              className="rounded-md border border-line px-3 py-1.5 hover:bg-page disabled:opacity-40 disabled:hover:bg-transparent"
              disabled={paging.page <= 1}
              onClick={() => onPage(paging.page - 1)}
            >
              Prev
            </button>
            <span className="text-muted">
              {paging.page} / {paging.totalPages}
            </span>
            <button
              className="rounded-md border border-line px-3 py-1.5 hover:bg-page disabled:opacity-40 disabled:hover:bg-transparent"
              disabled={paging.page >= paging.totalPages}
              onClick={() => onPage(paging.page + 1)}
            >
              Next
            </button>
          </div>
        </div>
      )}
    </section>
  )
}

export function TableSkeleton() {
  return (
    <section className="animate-pulse rounded-lg border border-line bg-card p-4">
      <div className="mb-4 h-4 w-20 rounded bg-gray-200" />
      {[0, 1, 2, 3, 4].map((i) => (
        <div key={i} className="mb-3 h-8 rounded bg-gray-100" />
      ))}
    </section>
  )
}

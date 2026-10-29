// Fetch helpers + response types. These types mirror the Go JSON structs
// (backend/internal/httpapi/handlers.go). If the backend changes a field, change it here too.

export type LeadStatus =
  | 'open'
  | 'applied'
  | 'accepted'
  | 'approved'
  | 'realized'
  | 'finished'
  | 'completed'
  | 'rejected'
  | 'withdrawn'

export interface Lead {
  id: string
  name: string
  email: string
  programme: string
  homeLC: string
  status: LeadStatus
  appliedAt: string | null
  approvedAt: string | null
  realizedAt: string | null
  createdAt: string
}

export interface Metrics {
  range: { from: string; to: string; timezone: string }
  totalApplications: number
  byStatus: Partial<Record<LeadStatus, number>>
  approved: number
  realized: number
  conversionRate: number | null // null when there are no applications
  lastUpdated: string
}

export interface LeadsPage {
  data: Lead[]
  paging: { page: number; perPage: number; totalItems: number; totalPages: number }
}

export interface DateRange {
  from?: string // YYYY-MM-DD, WIB
  to?: string
}

export interface LeadsQuery extends DateRange {
  status?: LeadStatus
  q?: string
  page?: number
  perPage?: number
}

// getJSON calls our own backend (same origin) and throws on non-2xx,
// using the backend's {"error": "..."} message when there is one.
async function getJSON<T>(path: string, params: Record<string, string | number | undefined> = {}): Promise<T> {
  const qs = new URLSearchParams()
  for (const [k, v] of Object.entries(params)) {
    if (v !== undefined && v !== '') qs.set(k, String(v))
  }
  const url = qs.size ? `${path}?${qs}` : path
  const res = await fetch(url)
  if (!res.ok) {
    const body = await res.json().catch(() => ({}))
    throw new Error(body.error ?? `HTTP ${res.status}`)
  }
  return res.json() as Promise<T>
}

export const fetchHealth = () => getJSON<{ status: string }>('/api/health')

export const fetchMetrics = (r: DateRange) => getJSON<Metrics>('/api/metrics', { from: r.from, to: r.to })

export const fetchLeads = (q: LeadsQuery) =>
  getJSON<LeadsPage>('/api/leads', {
    from: q.from,
    to: q.to,
    status: q.status,
    q: q.q,
    page: q.page,
    per_page: q.perPage,
  })

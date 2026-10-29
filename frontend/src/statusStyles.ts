import type { LeadStatus } from './api'

// One color per status. Add a status here when the backend adds one
// (backend/internal/leads/types.go).
export const STATUS_STYLE: Record<LeadStatus, { badge: string; bar: string }> = {
  open: { badge: 'border-gray-300 bg-gray-50 text-gray-600', bar: '#9ca3af' },
  applied: { badge: 'border-blue-300 bg-blue-50 text-blue-700', bar: '#3b82f6' },
  accepted: { badge: 'border-indigo-300 bg-indigo-50 text-indigo-700', bar: '#6366f1' },
  approved: { badge: 'border-green-300 bg-green-50 text-green-700', bar: '#22c55e' },
  realized: { badge: 'border-teal-300 bg-teal-50 text-teal-700', bar: '#14b8a6' },
  finished: { badge: 'border-green-600 bg-green-100 text-green-900', bar: '#15803d' },
  completed: { badge: 'border-green-700 bg-green-100 text-green-900', bar: '#166534' },
  rejected: { badge: 'border-red-300 bg-red-50 text-red-700', bar: '#ef4444' },
  withdrawn: { badge: 'border-red-200 bg-red-50 text-red-600', bar: '#f87171' },
}

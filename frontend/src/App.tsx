import { useEffect, useState } from 'react'
import { fetchHealth } from './api'

// Phase 5 skeleton: proves the frontend can reach the Go backend.
// The real dashboard UI replaces this in Phase 6.
export default function App() {
  const [status, setStatus] = useState<'loading' | 'ok' | 'error'>('loading')

  useEffect(() => {
    fetchHealth()
      .then((h) => setStatus(h.status === 'ok' ? 'ok' : 'error'))
      .catch(() => setStatus('error'))
  }, [])

  const label = { loading: 'Checking…', ok: 'Backend OK', error: 'Backend unreachable' }[status]
  const color = { loading: 'bg-gray-400', ok: 'bg-green-500', error: 'bg-red-500' }[status]

  return (
    <main className="mx-auto max-w-md p-6">
      <div className="rounded-lg border border-line bg-card p-6 shadow-sm">
        <h1 className="text-2xl font-bold">Leads Dashboard</h1>
        <p className="mt-1 text-sm text-muted">Frontend skeleton</p>
        <p className="mt-4 flex items-center gap-2 text-sm">
          <span className={`h-2.5 w-2.5 rounded-full ${color}`} />
          {label}
        </p>
      </div>
    </main>
  )
}

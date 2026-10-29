import { Bar, BarChart, CartesianGrid, Cell, ResponsiveContainer, Tooltip, XAxis, YAxis } from 'recharts'
import type { LeadStatus, Metrics } from '../api'
import { STATUS_STYLE } from '../statusStyles'
import StatusBadge from './StatusBadge'

// Counts per status: bar chart on desktop, badge grid on mobile.
export default function StatusChart({ byStatus }: { byStatus: Metrics['byStatus'] }) {
  // Object -> array in a fixed order, because charts need a list.
  const order: LeadStatus[] = ['applied', 'accepted', 'approved', 'realized', 'finished', 'completed', 'rejected', 'withdrawn']
  const data = order.map((s) => ({ status: s, count: byStatus[s] ?? 0 }))

  return (
    <section className="rounded-lg border border-line bg-card p-4">
      <h2 className="mb-3 font-semibold">Applications per status</h2>

      {/* Mobile (< 768px): compact grid of badges + numbers */}
      <ul className="grid grid-cols-2 gap-2 md:hidden">
        {data.map((d) => (
          <li key={d.status} className="flex items-center justify-between rounded-md border border-line px-3 py-2">
            <StatusBadge status={d.status} />
            <span className="font-semibold">{d.count}</span>
          </li>
        ))}
      </ul>

      {/* Desktop: bar chart */}
      <div className="hidden h-64 md:block">
        <ResponsiveContainer width="100%" height="100%">
          <BarChart data={data} margin={{ top: 8, right: 8, bottom: 0, left: -16 }}>
            <CartesianGrid vertical={false} stroke="#e5e7eb" />
            <XAxis
              dataKey="status"
              tickLine={false}
              axisLine={false}
              tick={{ fontSize: 12, fill: '#6b7280' }}
              tickFormatter={(s: string) => s.toUpperCase()}
            />
            <YAxis allowDecimals={false} tickLine={false} axisLine={false} tick={{ fontSize: 12, fill: '#6b7280' }} />
            <Tooltip
              cursor={{ fill: '#f5f6f8' }}
              formatter={(v: number) => [v, 'Applications']}
              labelFormatter={(s: string) => s.toUpperCase()}
            />
            <Bar dataKey="count" radius={[4, 4, 0, 0]} maxBarSize={56}>
              {data.map((d) => (
                <Cell key={d.status} fill={STATUS_STYLE[d.status].bar} />
              ))}
            </Bar>
          </BarChart>
        </ResponsiveContainer>
      </div>
    </section>
  )
}

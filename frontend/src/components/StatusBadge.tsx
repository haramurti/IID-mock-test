import type { LeadStatus } from '../api'
import { STATUS_STYLE } from '../statusStyles'

export default function StatusBadge({ status }: { status: LeadStatus }) {
  const style = STATUS_STYLE[status]?.badge ?? STATUS_STYLE.open.badge
  return (
    <span className={`inline-block rounded-full border px-2 py-0.5 text-[11px] font-semibold tracking-wide uppercase ${style}`}>
      {status}
    </span>
  )
}

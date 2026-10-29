import { BriefcaseIcon, ChartIcon, CloseIcon, HomeIcon, SettingsIcon, UsersIcon } from './icons'

// Decorative navigation, styled like EXPA's dark sidebar.
const items = [
  { label: 'Home', Icon: HomeIcon },
  { label: 'Leads Dashboard', Icon: ChartIcon, active: true },
  { label: 'People', Icon: UsersIcon },
  { label: 'Opportunities', Icon: BriefcaseIcon },
  { label: 'Settings', Icon: SettingsIcon },
]

interface Props {
  open: boolean // mobile drawer state
  onClose: () => void
}

export default function Sidebar({ open, onClose }: Props) {
  return (
    <>
      {/* Desktop: slim icon bar, always visible from md (768px) up */}
      <aside className="fixed inset-y-0 left-0 z-20 hidden w-16 flex-col items-center gap-2 bg-sidebar py-4 md:flex">
        <div className="mb-4 flex h-9 w-9 items-center justify-center rounded-md bg-primary text-sm font-bold text-white">
          L
        </div>
        {items.map(({ label, Icon, active }) => (
          <button
            key={label}
            title={label}
            aria-label={label}
            className={`flex h-10 w-10 items-center justify-center rounded-md transition-colors ${
              active ? 'bg-white/15 text-white' : 'text-white/60 hover:bg-white/10 hover:text-white'
            }`}
          >
            <Icon className="h-5 w-5" />
          </button>
        ))}
      </aside>

      {/* Mobile: drawer that slides in when ☰ is tapped */}
      {open && (
        <div className="fixed inset-0 z-40 md:hidden">
          <div className="absolute inset-0 bg-black/40" onClick={onClose} />
          <nav className="absolute inset-y-0 left-0 w-64 bg-sidebar p-4 text-white shadow-xl">
            <div className="mb-6 flex items-center justify-between">
              <span className="font-semibold">Menu</span>
              <button onClick={onClose} aria-label="Close menu" className="rounded p-1 hover:bg-white/10">
                <CloseIcon className="h-5 w-5" />
              </button>
            </div>
            {items.map(({ label, Icon, active }) => (
              <a
                key={label}
                href="#"
                onClick={onClose}
                className={`mb-1 flex items-center gap-3 rounded-md px-3 py-2 text-sm ${
                  active ? 'bg-white/15' : 'text-white/70 hover:bg-white/10'
                }`}
              >
                <Icon className="h-5 w-5" />
                {label}
              </a>
            ))}
          </nav>
        </div>
      )}
    </>
  )
}

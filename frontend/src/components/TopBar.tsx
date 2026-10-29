import { HomeIcon, MenuIcon } from './icons'

// White top bar: ☰ (mobile only), breadcrumb, avatar.
export default function TopBar({ onMenu }: { onMenu: () => void }) {
  return (
    <header className="sticky top-0 z-10 flex h-14 items-center gap-3 border-b border-line bg-card px-4 md:px-8">
      <button onClick={onMenu} aria-label="Open menu" className="-ml-1 rounded p-1 text-muted hover:bg-page md:hidden">
        <MenuIcon className="h-6 w-6" />
      </button>
      <nav className="flex items-center gap-2 text-sm text-muted" aria-label="Breadcrumb">
        <HomeIcon className="h-4 w-4" />
        <span>/</span>
        <span className="font-medium text-ink">Leads Dashboard</span>
      </nav>
      <div
        className="ml-auto flex h-8 w-8 items-center justify-center rounded-full bg-primary/10 text-sm font-semibold text-primary"
        title="Signed in user"
      >
        N
      </div>
    </header>
  )
}

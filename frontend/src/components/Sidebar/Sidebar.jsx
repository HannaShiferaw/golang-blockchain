import { NavLink } from 'react-router-dom'
import styles from './Sidebar.module.css'

const items = [
  { to: '/workflow', label: 'Workflow' },
  { to: '/ledger', label: 'Ledger' },
  { to: '/identities', label: 'Identities' },
  { to: '/dashboard', label: 'Dashboard' },
]

export function Sidebar() {
  return (
    <aside className={styles.sidebar}>
      <div className={styles.title}>Demo navigation</div>
      <nav className={styles.nav}>
        {items.map((it) => (
          <NavLink
            key={it.to}
            to={it.to}
            className={({ isActive }) => (isActive ? `${styles.link} ${styles.active}` : styles.link)}
          >
            {it.label}
          </NavLink>
        ))}
      </nav>

      <div className={styles.note}>
        Each action you submit becomes a signed transaction, appended into blocks, and stored in CouchDB.
      </div>
    </aside>
  )
}


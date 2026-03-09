import { useMemo, useState } from 'react'
import { Link, useLocation } from 'react-router-dom'
import styles from './Header.module.css'
import { useIdentities } from '../../hooks/useIdentities'
import { getActorId, setActorId } from '../../state/actor'

export function Header() {
  const { items, loading, error, reload } = useIdentities()
  const [actorId, setActor] = useState(() => getActorId())
  const loc = useLocation()

  const active = useMemo(() => items.find((x) => x.id === actorId), [items, actorId])

  function onPick(e) {
    const id = e.target.value
    setActorId(id)
    setActor(id)
  }

  return (
    <header className={styles.header}>
      <div className={styles.left}>
        <Link to="/workflow" className={styles.brand}>
          Coffee Export Consortium
        </Link>
        <div className={styles.path}>{loc.pathname}</div>
      </div>

      <div className={styles.right}>
        <div className={styles.picker}>
          <div className={styles.label}>Active stakeholder</div>
          <div className={styles.row}>
            <select className={styles.select} value={actorId} onChange={onPick} disabled={loading}>
              <option value="">Select…</option>
              {items.map((it) => (
                <option key={it.id} value={it.id}>
                  {it.name} · {it.role}
                </option>
              ))}
            </select>
            <button className={styles.smallBtn} onClick={reload} type="button">
              Refresh
            </button>
          </div>
          {active ? (
            <div className={styles.hint}>
              Using <b>{active.name}</b> ({active.role})
            </div>
          ) : (
            <div className={styles.hint}>Pick an identity to submit verified steps.</div>
          )}
          {error ? <div className={styles.error}>{error}</div> : null}
        </div>
      </div>
    </header>
  )
}


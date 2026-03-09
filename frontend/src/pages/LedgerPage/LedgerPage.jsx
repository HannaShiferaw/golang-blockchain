import { useEffect, useState } from 'react'
import { Card } from '../../components/Card/Card'
import { Button } from '../../components/Button/Button'
import { CodeBlock } from '../../components/CodeBlock/CodeBlock'
import { apiFetch } from '../../api/client'

export function LedgerPage() {
  const [items, setItems] = useState([])
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(true)

  async function load() {
    setLoading(true)
    setError('')
    try {
      const res = await apiFetch('/api/v1/blocks')
      setItems(res.items || [])
    } catch (e) {
      setError(e.message || 'Failed to load blocks')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    load()
  }, [])

  return (
    <div style={{ display: 'grid', gap: 14 }}>
      <Card
        title="Ledger blocks (CouchDB)"
        right={
          <Button variant="secondary" onClick={load} disabled={loading}>
            Refresh
          </Button>
        }
      >
        {error ? <div style={{ color: 'var(--danger)' }}>{error}</div> : null}
        {!items.length ? (
          <div style={{ color: 'var(--muted)' }}>No blocks yet. Run the workflow to create transactions.</div>
        ) : (
          <div style={{ display: 'grid', gap: 12 }}>
            {items.map((b) => (
              <Card
                key={b.hash}
                title={`Block #${b.index}`}
                right={<span style={{ color: 'var(--muted)', fontSize: 12 }}>{b.hash?.slice(0, 14)}…</span>}
              >
                <CodeBlock value={b} />
              </Card>
            ))}
          </div>
        )}
      </Card>
    </div>
  )
}


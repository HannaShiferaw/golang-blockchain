import { useState } from 'react'
import { Card } from '../../components/Card/Card'
import { Button } from '../../components/Button/Button'
import { Input } from '../../components/Input/Input'
import { CodeBlock } from '../../components/CodeBlock/CodeBlock'
import { apiFetch } from '../../api/client'
import { useIdentities } from '../../hooks/useIdentities'

export function IdentitiesPage() {
  const { items, loading, error, reload } = useIdentities()
  const [name, setName] = useState('')
  const [role, setRole] = useState('EXPORTER')
  const [busy, setBusy] = useState(false)
  const [created, setCreated] = useState(null)
  const [localError, setLocalError] = useState('')

  async function create() {
    setBusy(true)
    setLocalError('')
    try {
      const it = await apiFetch('/api/v1/identities', { method: 'POST', body: { name, role } })
      setCreated(it)
      setName('')
      await reload()
    } catch (e) {
      setLocalError(e.message || 'Failed to create identity')
    } finally {
      setBusy(false)
    }
  }

  return (
    <div style={{ display: 'grid', gap: 14 }}>
      <Card title="Create identity (PKI-issued certificate)">
        <div style={{ display: 'grid', gap: 10 }}>
          <Input label="Name" value={name} onChange={(e) => setName(e.target.value)} placeholder="e.g. exporter2" />
          <label style={{ display: 'grid', gap: 6 }}>
            <div style={{ fontSize: 12, color: 'var(--muted)' }}>Role</div>
            <select
              style={{
                background: 'rgba(255,255,255,0.06)',
                border: '1px solid var(--border)',
                color: 'var(--text)',
                padding: '10px 12px',
                borderRadius: 12,
              }}
              value={role}
              onChange={(e) => setRole(e.target.value)}
            >
              <option value="EXPORTER">EXPORTER</option>
              <option value="BUYER">BUYER</option>
              <option value="BANK">BANK</option>
              <option value="CUSTOMS">CUSTOMS</option>
              <option value="SHIPMENT">SHIPMENT</option>
            </select>
          </label>
          <div style={{ display: 'flex', gap: 10, flexWrap: 'wrap' }}>
            <Button onClick={create} disabled={!name || busy}>
              Create
            </Button>
            <Button variant="secondary" onClick={reload} disabled={loading}>
              Refresh list
            </Button>
          </div>
          {localError ? <div style={{ color: 'var(--danger)' }}>{localError}</div> : null}
          {created ? (
            <div>
              <div style={{ fontSize: 12, color: 'var(--muted)', marginBottom: 6 }}>Created</div>
              <CodeBlock value={created} />
            </div>
          ) : null}
        </div>
      </Card>

      <Card title="All identities">
        {error ? <div style={{ color: 'var(--danger)' }}>{error}</div> : null}
        {!items.length ? (
          <div style={{ color: 'var(--muted)' }}>No identities loaded.</div>
        ) : (
          <CodeBlock value={items} />
        )}
      </Card>
    </div>
  )
}


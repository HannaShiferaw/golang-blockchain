import { useMemo } from 'react'
import { Card } from '../../components/Card/Card'
import { CodeBlock } from '../../components/CodeBlock/CodeBlock'
import { useActorId } from '../../hooks/useActorId'
import { useIdentities } from '../../hooks/useIdentities'

export function DashboardPage() {
  const actorId = useActorId()
  const { items } = useIdentities()
  const actor = useMemo(() => items.find((x) => x.id === actorId), [items, actorId])

  return (
    <div style={{ display: 'grid', gap: 14 }}>
      <Card title="Demo overview">
        <div style={{ color: 'var(--muted)' }}>
          This is a consortium workflow demo: every stakeholder action is a signed transaction, validated by the smart
          contract rules, then appended into the ledger (CouchDB) and indexed (Postgres).
        </div>
      </Card>

      <Card title="Active stakeholder (for verification)">
        <CodeBlock value={actor || { hint: 'Select an identity in the header.' }} />
      </Card>
    </div>
  )
}


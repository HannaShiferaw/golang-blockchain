import { Link } from 'react-router-dom'
import { Card } from '../../components/Card/Card'

export function NotFoundPage() {
  return (
    <Card title="Not found">
      <div style={{ color: 'var(--muted)', display: 'grid', gap: 10 }}>
        <div>This page doesn’t exist.</div>
        <Link to="/workflow" style={{ color: 'var(--brand)' }}>
          Go to workflow
        </Link>
      </div>
    </Card>
  )
}


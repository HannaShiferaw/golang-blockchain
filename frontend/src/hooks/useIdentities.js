import { useCallback, useEffect, useState } from 'react'
import { apiFetch } from '../api/client'

export function useIdentities() {
  const [items, setItems] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  const load = useCallback(async () => {
    setLoading(true)
    setError('')
    try {
      const res = await apiFetch('/api/v1/identities')
      setItems(res.items || [])
    } catch (e) {
      setError(e.message || 'Failed to load identities')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    load()
  }, [load])

  return { items, loading, error, reload: load }
}


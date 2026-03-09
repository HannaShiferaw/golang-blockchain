import { useEffect, useState } from 'react'
import { getActorId } from '../state/actor'

export function useActorId() {
  const [actorId, setActorId] = useState(() => getActorId())

  useEffect(() => {
    const t = setInterval(() => {
      const v = getActorId()
      setActorId((cur) => (cur === v ? cur : v))
    }, 400)
    return () => clearInterval(t)
  }, [])

  return actorId
}


const KEY = 'coffee_demo_actor_id'

export function getActorId() {
  return localStorage.getItem(KEY) || ''
}

export function setActorId(id) {
  if (!id) localStorage.removeItem(KEY)
  else localStorage.setItem(KEY, id)
}


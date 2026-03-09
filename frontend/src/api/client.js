const BASE_URL = import.meta.env.VITE_BACKEND_BASE_URL || 'http://localhost:8080'

export async function apiFetch(path, { method = 'GET', actorId, body } = {}) {
  const headers = {
    'Content-Type': 'application/json',
  }
  if (actorId) headers['X-Actor-Id'] = actorId

  const res = await fetch(`${BASE_URL}${path}`, {
    method,
    headers,
    body: body ? JSON.stringify(body) : undefined,
  })

  const text = await res.text()
  const data = text ? safeJson(text) : null
  if (!res.ok) {
    const msg = (data && data.error) || `${res.status} ${res.statusText}`
    throw new Error(msg)
  }
  return data
}

function safeJson(text) {
  try {
    return JSON.parse(text)
  } catch {
    return { raw: text }
  }
}


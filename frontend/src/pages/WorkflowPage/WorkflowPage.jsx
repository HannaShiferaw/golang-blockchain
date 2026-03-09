import { useMemo, useState } from 'react'
import styles from './WorkflowPage.module.css'

import { Card } from '../../components/Card/Card'
import { Button } from '../../components/Button/Button'
import { Input } from '../../components/Input/Input'
import { CodeBlock } from '../../components/CodeBlock/CodeBlock'
import { apiFetch } from '../../api/client'
import { useIdentities } from '../../hooks/useIdentities'
import { useActorId } from '../../hooks/useActorId'

export function WorkflowPage() {
  const actorId = useActorId()
  const { items: identities } = useIdentities()

  const actor = useMemo(() => identities.find((x) => x.id === actorId), [identities, actorId])

  const [orderId, setOrderId] = useState('')
  const [shipmentId, setShipmentId] = useState('')

  const [createBuyerId, setCreateBuyerId] = useState('')
  const [createGrade, setCreateGrade] = useState('YIRGACHEFFE_GRADE_1')
  const [createQty, setCreateQty] = useState(1000)
  const [createUnit, setCreateUnit] = useState(6.5)

  const [lcAmount, setLcAmount] = useState(0)
  const [customsNotes, setCustomsNotes] = useState('Verified export documentation.')
  const [trackingNo, setTrackingNo] = useState('ET-COFFEE-TRK-0001')
  const [shipStatus, setShipStatus] = useState('PICKED_UP')
  const [shipLocation, setShipLocation] = useState('ADDIS_ABABA')

  const [orderState, setOrderState] = useState(null)
  const [shipmentState, setShipmentState] = useState(null)
  const [lastTx, setLastTx] = useState(null)
  const [error, setError] = useState('')
  const [busy, setBusy] = useState(false)

  async function run(fn) {
    setBusy(true)
    setError('')
    try {
      await fn()
    } catch (e) {
      setError(e.message || 'Request failed')
    } finally {
      setBusy(false)
    }
  }

  async function loadOrder(id = orderId) {
    if (!id) return
    const data = await apiFetch(`/api/v1/state?key=${encodeURIComponent(`state:order:${id}`)}`)
    setOrderState(data)
    if (!lcAmount && data?.totalUsd) setLcAmount(Number(data.totalUsd))
  }

  async function loadShipment(id = shipmentId) {
    if (!id) return
    const data = await apiFetch(`/api/v1/state?key=${encodeURIComponent(`state:shipment:${id}`)}`)
    setShipmentState(data)
  }

  return (
    <div className={styles.grid}>
      <Card
        title="Stakeholder verification"
        right={
          <Button variant="secondary" onClick={() => run(async () => loadOrder())} disabled={!orderId || busy}>
            Refresh order
          </Button>
        }
      >
        <div className={styles.row}>
          <div>
            <div className={styles.k}>Active stakeholder</div>
            <div className={styles.v}>
              {actor ? (
                <>
                  <b>{actor.name}</b> ({actor.role})
                </>
              ) : (
                <span className={styles.muted}>Pick one in the header.</span>
              )}
            </div>
          </div>
          <div>
            <div className={styles.k}>Current order ID</div>
            <div className={styles.inline}>
              <input
                className={styles.inlineInput}
                value={orderId}
                onChange={(e) => setOrderId(e.target.value)}
                placeholder="Paste orderId"
              />
              <Button variant="secondary" onClick={() => run(async () => loadOrder())} disabled={!orderId || busy}>
                Load
              </Button>
            </div>
          </div>
          <div>
            <div className={styles.k}>Current shipment ID</div>
            <div className={styles.inline}>
              <input
                className={styles.inlineInput}
                value={shipmentId}
                onChange={(e) => setShipmentId(e.target.value)}
                placeholder="Paste shipmentId"
              />
              <Button
                variant="secondary"
                onClick={() => run(async () => loadShipment())}
                disabled={!shipmentId || busy}
              >
                Load
              </Button>
            </div>
          </div>
        </div>
        {error ? <div className={styles.error}>{error}</div> : null}
      </Card>

      <Card title="1) Exporter creates export order (signed)">
        <div className={styles.form}>
          <label className={styles.field}>
            <div className={styles.k}>Buyer</div>
            <select
              className={styles.inlineInput}
              value={createBuyerId}
              onChange={(e) => setCreateBuyerId(e.target.value)}
            >
              <option value="">Select buyer identity…</option>
              {identities
                .filter((x) => x.role === 'BUYER')
                .map((x) => (
                  <option key={x.id} value={x.id}>
                    {x.name} ({x.id.slice(0, 8)}…)
                  </option>
                ))}
            </select>
          </label>
          <Input label="Coffee grade" value={createGrade} onChange={(e) => setCreateGrade(e.target.value)} />
          <Input
            label="Quantity (kg)"
            type="number"
            value={createQty}
            onChange={(e) => setCreateQty(Number(e.target.value))}
          />
          <Input
            label="Unit price (USD)"
            type="number"
            value={createUnit}
            onChange={(e) => setCreateUnit(Number(e.target.value))}
          />
          <div className={styles.actions}>
            <Button
              onClick={() =>
                run(async () => {
                  const res = await apiFetch('/api/v1/orders', {
                    method: 'POST',
                    actorId,
                    body: {
                      buyerId: createBuyerId,
                      coffeeGrade: createGrade,
                      quantityKg: createQty,
                      unitPriceUsd: createUnit,
                    },
                  })
                  setOrderId(res.orderId)
                  setLastTx(res.tx)
                  await loadOrder(res.orderId)
                })
              }
              disabled={!actorId || !createBuyerId || busy}
            >
              Create order
            </Button>
          </div>
        </div>
      </Card>

      <Card title="2) Buyer accepts order">
        <div className={styles.actions}>
          <Button
            onClick={() =>
              run(async () => {
                const res = await apiFetch(`/api/v1/orders/${orderId}/accept`, { method: 'POST', actorId })
                setLastTx(res.tx)
                await loadOrder()
              })
            }
            disabled={!actorId || !orderId || busy}
          >
            Accept order
          </Button>
        </div>
      </Card>

      <Card title="3) CBE issues Letter of Credit (LC)">
        <div className={styles.form}>
          <Input
            label="LC amount (USD)"
            type="number"
            value={lcAmount}
            onChange={(e) => setLcAmount(Number(e.target.value))}
          />
          <div className={styles.actions}>
            <Button
              onClick={() =>
                run(async () => {
                  const res = await apiFetch(`/api/v1/orders/${orderId}/lc`, {
                    method: 'POST',
                    actorId,
                    body: { amountUsd: lcAmount },
                  })
                  setLastTx(res.tx)
                  await loadOrder()
                })
              }
              disabled={!actorId || !orderId || !lcAmount || busy}
            >
              Issue LC
            </Button>
          </div>
        </div>
      </Card>

      <Card title="4) Customs approves clearance">
        <div className={styles.form}>
          <Input label="Notes" value={customsNotes} onChange={(e) => setCustomsNotes(e.target.value)} />
          <div className={styles.actions}>
            <Button
              onClick={() =>
                run(async () => {
                  const res = await apiFetch(`/api/v1/orders/${orderId}/customs-approve`, {
                    method: 'POST',
                    actorId,
                    body: { notes: customsNotes },
                  })
                  setLastTx(res.tx)
                  await loadOrder()
                })
              }
              disabled={!actorId || !orderId || busy}
            >
              Approve customs
            </Button>
          </div>
        </div>
      </Card>

      <Card title="5) Shipment creates shipment + tracking">
        <div className={styles.form}>
          <Input label="Tracking number" value={trackingNo} onChange={(e) => setTrackingNo(e.target.value)} />
          <div className={styles.actions}>
            <Button
              onClick={() =>
                run(async () => {
                  const res = await apiFetch(`/api/v1/orders/${orderId}/shipments`, {
                    method: 'POST',
                    actorId,
                    body: { trackingNo },
                  })
                  setShipmentId(res.shipmentId)
                  setLastTx(res.tx)
                  await loadOrder()
                  await loadShipment(res.shipmentId)
                })
              }
              disabled={!actorId || !orderId || busy}
            >
              Create shipment
            </Button>
          </div>
        </div>
      </Card>

      <Card title="6) Shipment updates shipment status">
        <div className={styles.form}>
          <label className={styles.field}>
            <div className={styles.k}>Status</div>
            <select className={styles.inlineInput} value={shipStatus} onChange={(e) => setShipStatus(e.target.value)}>
              <option value="PICKED_UP">PICKED_UP</option>
              <option value="EXPORTED">EXPORTED</option>
              <option value="ARRIVED">ARRIVED</option>
              <option value="DELIVERED">DELIVERED</option>
            </select>
          </label>
          <Input label="Location" value={shipLocation} onChange={(e) => setShipLocation(e.target.value)} />
          <div className={styles.actions}>
            <Button
              onClick={() =>
                run(async () => {
                  const res = await apiFetch(`/api/v1/shipments/${shipmentId}/status`, {
                    method: 'POST',
                    actorId,
                    body: { status: shipStatus, location: shipLocation },
                  })
                  setLastTx(res.tx)
                  await loadShipment()
                })
              }
              disabled={!actorId || !shipmentId || busy}
            >
              Update shipment
            </Button>
          </div>
        </div>
      </Card>

      <Card title="7) Buyer confirms delivery">
        <div className={styles.actions}>
          <Button
            onClick={() =>
              run(async () => {
                const res = await apiFetch(`/api/v1/orders/${orderId}/confirm-delivery`, { method: 'POST', actorId })
                setLastTx(res.tx)
                await loadOrder()
              })
            }
            disabled={!actorId || !orderId || busy}
          >
            Confirm delivery
          </Button>
        </div>
      </Card>

      <Card title="8) Bank releases payment (settlement)">
        <div className={styles.actions}>
          <Button
            onClick={() =>
              run(async () => {
                const res = await apiFetch(`/api/v1/orders/${orderId}/release-payment`, { method: 'POST', actorId })
                setLastTx(res.tx)
                await loadOrder()
              })
            }
            disabled={!actorId || !orderId || busy}
          >
            Release payment
          </Button>
        </div>
      </Card>

      <Card title="Live state (CouchDB world-state docs)">
        <div className={styles.split}>
          <div>
            <div className={styles.k}>Order state</div>
            <CodeBlock value={orderState || { hint: 'Load or create an order to see state.' }} />
          </div>
          <div>
            <div className={styles.k}>Shipment state</div>
            <CodeBlock value={shipmentState || { hint: 'Create/load a shipment to see state.' }} />
          </div>
        </div>
      </Card>

      <Card title="Last submitted transaction (signed)">
        <CodeBlock value={lastTx || { hint: 'No tx yet.' }} />
      </Card>
    </div>
  )
}


# Coffee Export Consortium Blockchain (Demo)

End-to-end demo of a **consortium blockchain workflow** for nationwide coffee export in Ethiopia, with 5 stakeholders:

- Exporter
- Buyer
- Commercial Bank of Ethiopia (CBE)
- Customs
- Shipment / Logistics

This repo contains:

- `backend/`: Go backend API, consortium ledger + smart contract, PKI/CA (Go), CouchDB + Postgres integration
- `frontend/`: React (JavaScript) web app with role-based screens

## What this demo enforces

- **PKI & signatures**: every action is a signed transaction by the stakeholder certificate.
- **Smart contract state machine**: the export process advances only when the required stakeholder verifies/approves.
- **Immutable audit**: transactions are appended into a block chain (hash-linked) stored in CouchDB.
- **Queryable data**: key fields are indexed into Postgres for fast filtering/reporting.

## Quick start (Docker)

Prereqs: Docker Desktop + `docker compose`.

On Windows, make sure **Docker Desktop is running** (and you can run `docker ps` successfully). If you see errors about the `docker_engine` pipe, start Docker Desktop or run your terminal as Administrator.

1. Copy env file:

```bash
cp .env.example .env
```

2. Start everything:

```bash
docker compose up --build
```

3. Open the UI:

- Frontend: `http://localhost:5173`
- Backend: `http://localhost:8080`

## Quick start (No Docker)

This runs the demo in **in-memory mode** (no CouchDB/Postgres), which is still enough to demo the full verified workflow.

1. Start backend:

```bash
cd backend
SEED_ON_START=true BACKEND_PORT=8080 go run ./cmd/api
```

2. Start frontend:

```bash
cd frontend
npm install
npm run dev
```

3. Open `http://localhost:5173`

## Default demo accounts

On first boot, the backend seeds 5 demo identities and their certificates:

- exporter1 (Exporter)
- buyer1 (Buyer)
- cbe1 (Bank)
- customs1 (Customs)
- shipper1 (Shipment)

The UI lets you pick a role/account and run the full end-to-end flow.

## Folder structure

```text
.
├── backend/
│   ├── cmd/
│   │   └── api/
│   ├── internal/
│   │   ├── contract/
│   │   ├── ledger/
│   │   ├── pki/
│   │   ├── repo/
│   │   ├── service/
│   │   └── transport/
│   └── migrations/
└── frontend/
    └── src/
        ├── components/
        └── pages/
```



import { Navigate, Route, Routes } from 'react-router-dom'
import { Header } from './components/Header/Header'
import { Sidebar } from './components/Sidebar/Sidebar'
import styles from './App.module.css'

import { DashboardPage } from './pages/DashboardPage/DashboardPage'
import { WorkflowPage } from './pages/WorkflowPage/WorkflowPage'
import { LedgerPage } from './pages/LedgerPage/LedgerPage'
import { IdentitiesPage } from './pages/IdentitiesPage/IdentitiesPage'
import { NotFoundPage } from './pages/NotFoundPage/NotFoundPage'

export default function App() {
  return (
    <div className={styles.shell}>
      <Header />
      <div className={styles.body}>
        <Sidebar />
        <main className={styles.main}>
          <Routes>
            <Route path="/" element={<Navigate to="/workflow" replace />} />
            <Route path="/dashboard" element={<DashboardPage />} />
            <Route path="/workflow" element={<WorkflowPage />} />
            <Route path="/ledger" element={<LedgerPage />} />
            <Route path="/identities" element={<IdentitiesPage />} />
            <Route path="*" element={<NotFoundPage />} />
          </Routes>
        </main>
      </div>
    </div>
  )
}

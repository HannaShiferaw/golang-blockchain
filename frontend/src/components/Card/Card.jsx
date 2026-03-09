import styles from './Card.module.css'

export function Card({ title, children, right }) {
  return (
    <section className={styles.card}>
      {title ? (
        <header className={styles.header}>
          <div className={styles.title}>{title}</div>
          {right ? <div className={styles.right}>{right}</div> : null}
        </header>
      ) : null}
      <div className={styles.body}>{children}</div>
    </section>
  )
}


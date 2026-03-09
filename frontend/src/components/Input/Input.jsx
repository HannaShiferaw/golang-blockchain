import styles from './Input.module.css'

export function Input({ label, hint, ...props }) {
  return (
    <label className={styles.wrap}>
      {label ? <div className={styles.label}>{label}</div> : null}
      <input className={styles.input} {...props} />
      {hint ? <div className={styles.hint}>{hint}</div> : null}
    </label>
  )
}


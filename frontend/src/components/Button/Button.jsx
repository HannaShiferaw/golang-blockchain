import styles from './Button.module.css'

export function Button({ children, variant = 'primary', ...props }) {
  const cls = variant === 'secondary' ? `${styles.btn} ${styles.secondary}` : styles.btn
  return (
    <button className={cls} {...props}>
      {children}
    </button>
  )
}


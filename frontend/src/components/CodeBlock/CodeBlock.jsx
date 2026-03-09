import styles from './CodeBlock.module.css'

export function CodeBlock({ value }) {
  const text = typeof value === 'string' ? value : JSON.stringify(value, null, 2)
  return <pre className={styles.pre}>{text}</pre>
}


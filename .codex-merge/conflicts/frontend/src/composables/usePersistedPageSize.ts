import { getConfiguredTableDefaultPageSize, normalizeTablePageSize } from '@/utils/tablePreferences'

const STORAGE_KEY = 'table-page-size'

export function getPersistedPageSize(fallback = getConfiguredTableDefaultPageSize()): number {
<<<<<<< HEAD
  if (typeof window !== 'undefined' && window.__APP_CONFIG__?.table_default_page_size !== undefined) {
    return normalizeTablePageSize(getConfiguredTableDefaultPageSize())
=======
  const configuredDefault = normalizeTablePageSize(getConfiguredTableDefaultPageSize() || fallback)
  const fallbackDefault = normalizeTablePageSize(fallback)
  if (configuredDefault !== fallbackDefault || configuredDefault !== normalizeTablePageSize(20)) {
    return configuredDefault
>>>>>>> 1157e4cfe6271a81c9dc96df84af6c6a3b22f831
  }

  if (typeof window !== 'undefined') {
    try {
      const stored = window.localStorage.getItem(STORAGE_KEY)
      if (stored !== null) {
        const parsed = Number(stored)
        if (Number.isFinite(parsed)) {
          return normalizeTablePageSize(parsed)
        }
      }
    } catch (error) {
      console.warn('Failed to read persisted page size:', error)
    }
  }
  return configuredDefault
}

export function setPersistedPageSize(size: number): void {
  if (typeof window === 'undefined') return
  try {
    window.localStorage.setItem(STORAGE_KEY, String(size))
  } catch (error) {
    console.warn('Failed to persist page size:', error)
  }
}

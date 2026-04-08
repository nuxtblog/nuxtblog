import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'
import 'dayjs/locale/en'

dayjs.extend(relativeTime)

// ─── Internal helpers ────────────────────────────────────────────────────────

/**
 * Parse a date string from the backend. Handles:
 * - RFC3339 with timezone:  "2025-03-29T06:30:00Z"  → parsed as UTC
 * - GoFrame without TZ:     "2025-03-29 14:30:05"   → parsed as browser local time
 */
export function parseBackendDate(s: string | number | Date): Date {
  if (s instanceof Date) return s
  if (typeof s === 'number') return new Date(s)
  if (!s) return new Date(NaN)
  // Normalize GoFrame "YYYY-MM-DD HH:MM:SS" to ISO-like (browsers parse as local time)
  return new Date((s as string).replace(' ', 'T'))
}

function toDayjsLocale(locale: string): string {
  return locale.startsWith('zh') ? 'zh-cn' : 'en'
}

function toIntlLocale(locale: string): string {
  return locale.startsWith('zh') ? 'zh-CN' : 'en-US'
}

// ─── datetime-local input helpers ────────────────────────────────────────────

/**
 * Convert any date to the value format for <input type="datetime-local">.
 * Uses browser LOCAL time so the user sees their own timezone.
 * e.g.  "2025-03-29T06:30:00Z" (UTC+8) → "2025-03-29T14:30"
 */
export function toDatetimeInputValue(date: Date | string | number): string {
  const d = parseBackendDate(date)
  if (isNaN(d.getTime())) return ''
  const pad = (n: number) => String(n).padStart(2, '0')
  return (
    `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}` +
    `T${pad(d.getHours())}:${pad(d.getMinutes())}`
  )
}

/**
 * Convert a datetime-local string (local time, no TZ) to UTC ISO for the API.
 * e.g.  "2025-03-29T14:30" (UTC+8 browser) → "2025-03-29T06:30:00.000Z"
 */
export function fromDatetimeInputValue(localStr: string): string {
  if (!localStr) return ''
  return new Date(localStr).toISOString()
}

// ─── Relative time ────────────────────────────────────────────────────────────

/**
 * Relative time: "3 minutes ago" / "2天前" etc.
 * Uses dayjs for natural language output.
 * @param locale  app locale string, e.g. 'zh' | 'en'
 */
export function formatRelative(date: Date | string | number, locale = 'zh'): string {
  const d = parseBackendDate(date as string)
  if (isNaN(d.getTime())) return ''
  return dayjs(d).locale(toDayjsLocale(locale)).fromNow()
}

// ─── Absolute time ────────────────────────────────────────────────────────────

/**
 * Full absolute date: "2025年3月29日" / "March 29, 2025"
 */
export function formatAbsolute(
  date: Date | string | number,
  locale = 'zh',
  opts: Intl.DateTimeFormatOptions = { year: 'numeric', month: 'long', day: 'numeric' },
): string {
  const d = parseBackendDate(date as string)
  if (isNaN(d.getTime())) return ''
  return d.toLocaleDateString(toIntlLocale(locale), opts)
}

/**
 * Short absolute date: "2025/3/29" / "3/29/2025"
 */
export function formatAbsoluteShort(date: Date | string | number, locale = 'zh'): string {
  return formatAbsolute(date, locale, { year: 'numeric', month: 'numeric', day: 'numeric' })
}

/**
 * Date + time: "2025年3月29日 14:30" / "Mar 29, 2025, 2:30 PM"
 */
export function formatDateTime(date: Date | string | number, locale = 'zh'): string {
  const d = parseBackendDate(date as string)
  if (isNaN(d.getTime())) return ''
  return d.toLocaleString(toIntlLocale(locale), {
    year: 'numeric', month: 'short', day: 'numeric',
    hour: '2-digit', minute: '2-digit',
  })
}

/**
 * Month label: "三月" / "March"
 */
export function formatMonth(date: Date | string | number, locale = 'zh'): string {
  const d = parseBackendDate(date as string)
  if (isNaN(d.getTime())) return ''
  return d.toLocaleDateString(toIntlLocale(locale), { month: 'long' })
}

/**
 * Year: "2025"
 */
export function formatYear(date: Date | string | number): string {
  const d = parseBackendDate(date as string)
  if (isNaN(d.getTime())) return ''
  return String(d.getFullYear())
}

/**
 * Day of month: "05"
 */
export function formatDay(date: Date | string | number): string {
  const d = parseBackendDate(date as string)
  if (isNaN(d.getTime())) return ''
  return String(d.getDate()).padStart(2, '0')
}

// ─── Smart auto-fallback ──────────────────────────────────────────────────────

/**
 * Smart format: relative for recent dates (< 30 days), absolute for older.
 * @param locale  app locale string
 */
export function formatSmart(date: Date | string | number, locale = 'zh'): string {
  const d = parseBackendDate(date as string)
  if (isNaN(d.getTime())) return ''
  const days = (Date.now() - d.getTime()) / 86_400_000
  return days < 30 ? formatRelative(d, locale) : formatAbsolute(d, locale)
}

/**
 * Auto-importable formatDate — uses browser locale (navigator.language).
 * Outputs relative time for recent dates, absolute for older.
 * For full i18n control, use useFormatDate() composable instead.
 */
export function formatDate(input: string | number | Date): string {
  const locale = typeof navigator !== 'undefined'
    ? (navigator.language.startsWith('zh') ? 'zh' : 'en')
    : 'zh'
  return formatSmart(input, locale)
}

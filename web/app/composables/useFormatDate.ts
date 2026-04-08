import {
  formatRelative,
  formatAbsolute,
  formatAbsoluteShort,
  formatDateTime,
  formatSmart,
  formatMonth,
  formatYear,
  formatDay,
} from '~/utils/date'

/**
 * Locale-aware date formatting composable.
 * All methods automatically use the app's current i18n locale.
 */
export const useFormatDate = () => {
  const { locale } = useI18n()
  const l = () => locale.value

  return {
    /** "3分钟前" / "2 hours ago" — relative, falls back to absolute after 30 days */
    formatRelative: (date: string | Date | number) => formatRelative(date, l()),
    /** "2025年3月29日" / "March 29, 2025" */
    formatAbsolute: (date: string | Date | number, opts?: Intl.DateTimeFormatOptions) =>
      formatAbsolute(date, l(), opts),
    /** "2025/3/29" / "3/29/2025" */
    formatAbsoluteShort: (date: string | Date | number) => formatAbsoluteShort(date, l()),
    /** "2025年3月29日 14:30" / "Mar 29, 2025, 2:30 PM" */
    formatDateTime: (date: string | Date | number) => formatDateTime(date, l()),
    /** Smart: relative if < 30 days, else absolute */
    formatSmart: (date: string | Date | number) => formatSmart(date, l()),
    /** "三月" / "March" */
    formatMonth: (date: string | Date | number) => formatMonth(date, l()),
    /** "2025" */
    formatYear: (date: string | Date | number) => formatYear(date),
    /** "29" */
    formatDay: (date: string | Date | number) => formatDay(date),
  }
}

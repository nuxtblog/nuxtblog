/**
 * Locale-aware number formatting composable.
 * Usage: const { abbreviate } = useNumberFormat()
 */
export const useNumberFormat = () => {
  const { locale } = useI18n()
  const abbreviate = (num: number) => abbreviateNumber(num, locale.value)
  return { abbreviate }
}

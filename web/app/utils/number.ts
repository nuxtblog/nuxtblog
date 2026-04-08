/**
 * Abbreviate a number with locale-aware suffixes.
 * zh: 1k / 1万 / 1亿
 * en: 1k / 1M / 1B
 */
export const abbreviateNumber = (num: number, locale = 'zh'): string => {
  if (num < 0) return '0'
  if (num < 1000) return num.toString()

  if (locale.startsWith('zh')) {
    if (num < 10000) return (num / 1000).toFixed(num % 1000 === 0 ? 0 : 1) + 'k'
    if (num < 100000000) return (num / 10000).toFixed(num % 10000 === 0 ? 0 : 1) + '万'
    return (num / 100000000).toFixed(num % 100000000 === 0 ? 0 : 1) + '亿'
  }

  // Western: k / M / B
  if (num < 1000000) return (num / 1000).toFixed(num % 1000 === 0 ? 0 : 1) + 'k'
  if (num < 1000000000) return (num / 1000000).toFixed(num % 1000000 === 0 ? 0 : 1) + 'M'
  return (num / 1000000000).toFixed(num % 1000000000 === 0 ? 0 : 1) + 'B'
}

export function clampMax99(value: number) {
  return value > 99 ? 99 : value
}

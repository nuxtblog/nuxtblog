/**
 * Compare two semver-like strings. Returns 1 if a > b, -1 if a < b, 0 if equal.
 * Tolerates leading "v", extra whitespace, and missing minor/patch (treats as 0).
 * Pre-release suffixes (-beta, -rc1) are ignored for comparison — only the
 * numeric MAJOR.MINOR.PATCH portion is used. Good enough for plugin versions.
 */
export function semverCompare(a: string, b: string): number {
  const parse = (v: string): [number, number, number] => {
    const s = (v || '').trim().replace(/^v/i, '').split(/[-+]/)[0] ?? ''
    const parts = s.split('.').map(n => parseInt(n, 10) || 0)
    return [parts[0] ?? 0, parts[1] ?? 0, parts[2] ?? 0]
  }
  const [a1, a2, a3] = parse(a)
  const [b1, b2, b3] = parse(b)
  if (a1 !== b1) return a1 > b1 ? 1 : -1
  if (a2 !== b2) return a2 > b2 ? 1 : -1
  if (a3 !== b3) return a3 > b3 ? 1 : -1
  return 0
}

/** True iff `candidate` is strictly newer than `current`. */
export function isNewerVersion(candidate: string, current: string): boolean {
  return semverCompare(candidate, current) > 0
}

import { defineStore } from 'pinia'

/**
 * Plugin Context State Machine
 *
 * Maintains a global reactive context that drives `when` expression evaluation
 * for plugin contribution points. Plugins can also register custom context keys.
 *
 * Built-in keys:
 * - post.status, post.wordCount, post.hasContent, post.type
 * - editor.hasSelection, editor.selectionLength
 * - user.role
 */
export const usePluginContextStore = defineStore('plugin-context', () => {
  const context = reactive<Record<string, unknown>>({
    'post.status': 'draft',
    'post.wordCount': 0,
    'post.hasContent': false,
    'post.type': 'post',
    'editor.hasSelection': false,
    'editor.selectionLength': 0,
    'user.role': 'admin',
  })

  /** Set a context key. Plugins use this to register custom keys. */
  function set(key: string, value: unknown) {
    context[key] = value
  }

  /** Get a context key value. */
  function get(key: string): unknown {
    return context[key]
  }

  /**
   * Evaluate a `when` expression against the current context.
   * Only supports comparison operators: ==, !=, >, <, >=, <=, &&, ||, !
   * No function calls or assignments allowed (security sandbox).
   */
  function evaluateWhen(expr: string): boolean {
    if (!expr || !expr.trim()) return true
    try {
      return evaluateWhenExpr(expr.trim(), context)
    }
    catch {
      return true // show on parse error
    }
  }

  return { context, set, get, evaluateWhen }
})

/**
 * Safe `when` expression evaluator. Only allows:
 * - Variable references (dotted names like "post.wordCount")
 * - Comparison: ==, !=, >, <, >=, <=
 * - Logical: &&, ||, !
 * - Literals: numbers, strings, true, false, null
 */
function evaluateWhenExpr(expr: string, ctx: Record<string, unknown>): boolean {
  // Handle logical OR
  if (expr.includes('||')) {
    const parts = splitOutsideParens(expr, '||')
    if (parts.length > 1) {
      return parts.some(p => evaluateWhenExpr(p.trim(), ctx))
    }
  }

  // Handle logical AND
  if (expr.includes('&&')) {
    const parts = splitOutsideParens(expr, '&&')
    if (parts.length > 1) {
      return parts.every(p => evaluateWhenExpr(p.trim(), ctx))
    }
  }

  // Handle NOT
  if (expr.startsWith('!')) {
    return !evaluateWhenExpr(expr.slice(1).trim(), ctx)
  }

  // Handle parentheses
  if (expr.startsWith('(') && expr.endsWith(')')) {
    return evaluateWhenExpr(expr.slice(1, -1).trim(), ctx)
  }

  // Handle comparisons
  const compMatch = expr.match(/^(.+?)\s*(===|!==|==|!=|>=|<=|>|<)\s*(.+)$/)
  if (compMatch) {
    const [, left, op, right] = compMatch
    const lVal = resolveValue(left.trim(), ctx)
    const rVal = resolveValue(right.trim(), ctx)
    switch (op) {
      case '===': case '==': return lVal === rVal
      case '!==': case '!=': return lVal !== rVal
      case '>': return (lVal as number) > (rVal as number)
      case '<': return (lVal as number) < (rVal as number)
      case '>=': return (lVal as number) >= (rVal as number)
      case '<=': return (lVal as number) <= (rVal as number)
    }
  }

  // Bare truthy check
  return !!resolveValue(expr, ctx)
}

function resolveValue(token: string, ctx: Record<string, unknown>): unknown {
  if (token === 'true') return true
  if (token === 'false') return false
  if (token === 'null') return null
  if (/^['"].*['"]$/.test(token)) return token.slice(1, -1)
  if (!Number.isNaN(Number(token))) return Number(token)
  // Context variable reference
  return ctx[token]
}

function splitOutsideParens(str: string, sep: string): string[] {
  const parts: string[] = []
  let depth = 0
  let current = ''
  for (let i = 0; i < str.length; i++) {
    if (str[i] === '(') depth++
    else if (str[i] === ')') depth--
    if (depth === 0 && str.substring(i, i + sep.length) === sep) {
      parts.push(current)
      current = ''
      i += sep.length - 1
    }
    else {
      current += str[i]
    }
  }
  parts.push(current)
  return parts
}

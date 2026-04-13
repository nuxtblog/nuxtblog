import { marked } from 'marked'
import hljs from 'highlight.js'
import katex from 'katex'
import 'katex/dist/katex.min.css'
import markedAlert from 'marked-alert'

// ── KaTeX math extension ─────────────────────────────────────────────────
const mathExtension = {
  extensions: [{
    name: 'blockMath',
    level: 'block' as const,
    start(src: string) { return src.indexOf('$$') },
    tokenizer(src: string) {
      const match = src.match(/^\$\$([\s\S]+?)\$\$/)
      if (match) return { type: 'blockMath', raw: match[0], text: (match[1] ?? '').trim() }
    },
    renderer(token: { text: string }) {
      return `<div class="math-block">${katex.renderToString(token.text, { displayMode: true, throwOnError: false })}</div>`
    }
  }, {
    name: 'inlineMath',
    level: 'inline' as const,
    start(src: string) { return src.indexOf('$') },
    tokenizer(src: string) {
      const match = src.match(/^\$([^\s$](?:[^$]*[^\s$])?)\$/)
      if (match) return { type: 'inlineMath', raw: match[0], text: match[1] }
    },
    renderer(token: { text: string }) {
      return katex.renderToString(token.text, { throwOnError: false })
    }
  }]
}

// ── Code renderer: highlight.js + mermaid passthrough ────────────────────
const renderer = new marked.Renderer()
renderer.code = ({ text, lang }: { text: string; lang?: string }) => {
  if (lang === 'mermaid') {
    return `<pre><code class="language-mermaid">${text}</code></pre>`
  }
  const language = lang && hljs.getLanguage(lang) ? lang : 'plaintext'
  const highlighted = hljs.highlight(text, { language }).value
  return `<pre><code class="hljs language-${language}">${highlighted}</code></pre>`
}

// ── Register all extensions (runs once at module load) ───────────────────
marked.use(mathExtension)
marked.use(markedAlert())
marked.use({ renderer })

/**
 * Preprocess markdown before passing to marked.
 *
 * TipTap's blockquote serialization inserts blank lines between
 * `> [!TYPE]` and `> content`, breaking marked-alert parsing.
 * This merges them back into a single blockquote block.
 */
export function normalizeMarkdown(md: string): string {
  return md.replace(
    /^([ \t]*> \[!(NOTE|TIP|IMPORTANT|WARNING|CAUTION)\])\s*\n\s*\n((?:[ \t]*> .*(?:\n|$))+)/gm,
    '$1\n$3'
  )
}

/**
 * Render markdown to HTML with all extensions (KaTeX, marked-alert, highlight.js, mermaid).
 * Use this instead of calling `marked()` directly.
 */
export function renderMarkdown(md: string): string {
  return marked(normalizeMarkdown(md)) as string
}

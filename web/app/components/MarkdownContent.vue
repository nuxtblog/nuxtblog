<script setup lang="ts">
import { renderMarkdown } from '~/utils/markdown'

const props = defineProps<{
  content?: string
  html?: string
}>()

const rendered = computed(() =>
  props.html ?? (props.content ? renderMarkdown(props.content) : '')
)

useHighlightTheme()

const containerRef = ref<HTMLElement>()

// ── Copy buttons ─────────────────────────────────────────────────────────
function addCopyButtons() {
  if (!containerRef.value) return
  containerRef.value.querySelectorAll('pre').forEach((pre) => {
    if (pre.querySelector('.copy-btn')) return
    // Skip mermaid blocks (they will be replaced)
    if (pre.querySelector('code.language-mermaid')) return
    const btn = document.createElement('button')
    btn.className = 'copy-btn'
    btn.textContent = '复制'
    btn.addEventListener('click', async () => {
      const code = pre.querySelector('code')?.innerText ?? ''
      await navigator.clipboard.writeText(code)
      btn.textContent = '已复制'
      setTimeout(() => { btn.textContent = '复制' }, 1500)
    })
    pre.style.position = 'relative'
    pre.appendChild(btn)
  })
}

// ── Mermaid client-side rendering ────────────────────────────────────────
let mermaidInstance: typeof import('mermaid')['default'] | null = null

async function renderMermaidBlocks() {
  if (!import.meta.client || !containerRef.value) return
  const blocks = containerRef.value.querySelectorAll('code.language-mermaid')
  if (!blocks.length) return

  // Lazy-load mermaid only when needed
  if (!mermaidInstance) {
    const mod = await import('mermaid')
    mermaidInstance = mod.default
    const isDark = document.documentElement.classList.contains('dark')
    mermaidInstance.initialize({ startOnLoad: false, theme: isDark ? 'dark' : 'default' })
  }

  for (const code of blocks) {
    const pre = code.closest('pre')
    if (!pre || pre.dataset.mermaidRendered) continue
    pre.dataset.mermaidRendered = 'true'
    const source = code.textContent ?? ''
    try {
      const id = `mermaid-${Math.random().toString(36).slice(2, 8)}`
      const { svg } = await mermaidInstance.render(id, source)
      const wrapper = document.createElement('div')
      wrapper.className = 'mermaid-diagram'
      wrapper.innerHTML = svg
      pre.replaceWith(wrapper)
    } catch {
      // Render failed — keep code block visible
    }
  }
}

// ── Image lightbox ───────────────────────────────────────────────────────
function addImageLightbox() {
  if (!import.meta.client || !containerRef.value) return
  containerRef.value.querySelectorAll('img').forEach((img) => {
    if (img.dataset.lightbox) return
    img.dataset.lightbox = 'true'
    img.style.cursor = 'zoom-in'
    img.addEventListener('click', () => {
      const overlay = document.createElement('div')
      overlay.className = 'lightbox-overlay'
      const clone = document.createElement('img')
      clone.src = img.src
      clone.alt = img.alt
      clone.className = 'lightbox-img'
      overlay.appendChild(clone)
      overlay.addEventListener('click', () => {
        overlay.classList.add('lightbox-closing')
        overlay.addEventListener('animationend', () => overlay.remove())
      })
      document.body.appendChild(overlay)
    })
  })
}

function postRender() {
  addCopyButtons()
  addImageLightbox()
  renderMermaidBlocks()
}

watch(rendered, () => nextTick(postRender))
onMounted(() => nextTick(postRender))
</script>

<template>
  <div ref="containerRef" class="prose" v-html="rendered" />
</template>

<style>
/* highlight.js theme is dynamically injected by useHighlightTheme() */

/* Image lightbox (appended to body, must be global) */
.lightbox-overlay {
  position: fixed;
  inset: 0;
  z-index: 9999;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.85);
  cursor: zoom-out;
  animation: lightbox-in 0.2s ease;
}
.lightbox-overlay.lightbox-closing {
  animation: lightbox-out 0.15s ease forwards;
}
.lightbox-img {
  max-width: 90vw;
  max-height: 90vh;
  object-fit: contain;
  border-radius: 0.5rem;
  box-shadow: 0 8px 40px rgba(0, 0, 0, 0.4);
}
@keyframes lightbox-in {
  from { opacity: 0; }
  to   { opacity: 1; }
}
@keyframes lightbox-out {
  from { opacity: 1; }
  to   { opacity: 0; }
}
</style>

<style scoped>
.prose {
  color: var(--ui-text);
}
.prose :deep(h1),
.prose :deep(h2),
.prose :deep(h3),
.prose :deep(h4),
.prose :deep(h5),
.prose :deep(h6) {
  color: var(--ui-text-highlighted);
  font-weight: 700;
  margin-top: 2em;
  margin-bottom: 0.6em;
  line-height: 1.3;
}
.prose :deep(h1) { font-size: 2em; }
.prose :deep(h2) { font-size: 1.5em; border-bottom: 1px solid var(--ui-border); padding-bottom: 0.3em; }
.prose :deep(h3) { font-size: 1.25em; }
.prose :deep(p) {
  margin-top: 1em;
  margin-bottom: 1em;
  line-height: 1.8;
}
.prose :deep(a) {
  color: var(--ui-primary);
  text-decoration: underline;
  font-weight: 500;
}
.prose :deep(a:hover) { opacity: 0.8; }
.prose :deep(code) {
  background-color: var(--ui-bg-muted);
  padding: 0.2em 0.4em;
  border-radius: 0.25rem;
  font-size: 0.875em;
  font-family: "JetBrains Mono", "Courier New", monospace;
}
.prose :deep(pre) {
  padding: 1.25em;
  border-radius: 0.75rem;
  overflow-x: auto;
  margin: 1.5em 0;
  border: 1px solid var(--ui-border);
}
.prose :deep(pre code) {
  background-color: transparent;
  padding: 0;
  font-size: 0.875em;
}
.prose :deep(blockquote) {
  border-left: 3px solid var(--ui-primary);
  padding: 0.75em 1.25em;
  margin: 1.5em 0;
  background: var(--ui-bg-muted);
  border-radius: 0 0.5rem 0.5rem 0;
  color: var(--ui-text-muted);
}
.prose :deep(ul),
.prose :deep(ol) {
  margin-top: 1em;
  margin-bottom: 1em;
  padding-left: 1.5em;
}
.prose :deep(li) { margin-top: 0.4em; margin-bottom: 0.4em; }
.prose :deep(img) {
  border-radius: 0.75rem;
  margin: 2em auto;
  max-width: 100%;
  height: auto;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}
.prose :deep(table) {
  width: 100%;
  border-collapse: collapse;
  margin: 1.5em 0;
}
.prose :deep(th),
.prose :deep(td) {
  border: 1px solid var(--ui-border);
  padding: 0.75em 1em;
  text-align: left;
}
.prose :deep(th) {
  background-color: var(--ui-bg-elevated);
  font-weight: 600;
}
.prose :deep(hr) {
  border: none;
  border-top: 1px solid var(--ui-border);
  margin: 2em 0;
}

/* KaTeX — ensure formula color follows text in dark mode */
.prose :deep(.katex) {
  color: inherit;
}
/* KaTeX math blocks */
.prose :deep(.math-block) {
  text-align: center;
  margin: 1.5em 0;
  overflow-x: auto;
}

/* Mermaid diagrams */
.prose :deep(.mermaid-diagram) {
  text-align: center;
  margin: 1.5em 0;
  overflow-x: auto;
}
.prose :deep(.mermaid-diagram svg) {
  max-width: 100%;
  height: auto;
}

/* GitHub Alerts / Callouts */
.prose :deep(.markdown-alert) {
  padding: 0.75em 1.25em;
  margin: 1.5em 0;
  border-left: 4px solid;
  border-radius: 0 0.5rem 0.5rem 0;
}
.prose :deep(.markdown-alert-note)      { border-color: var(--ui-primary);    background: color-mix(in srgb, var(--ui-primary) 8%, transparent); }
.prose :deep(.markdown-alert-tip)       { border-color: var(--ui-success);    background: color-mix(in srgb, var(--ui-success) 8%, transparent); }
.prose :deep(.markdown-alert-important) { border-color: #a855f7;             background: color-mix(in srgb, #a855f7 8%, transparent); }
.prose :deep(.markdown-alert-warning)   { border-color: var(--ui-warning);    background: color-mix(in srgb, var(--ui-warning) 8%, transparent); }
.prose :deep(.markdown-alert-caution)   { border-color: var(--ui-error);      background: color-mix(in srgb, var(--ui-error) 8%, transparent); }
.prose :deep(.markdown-alert-title) {
  font-weight: 600;
  margin-bottom: 0.25em;
}
.prose :deep(.markdown-alert p:first-child) {
  margin-top: 0;
}
.prose :deep(.markdown-alert p:last-child) {
  margin-bottom: 0;
}

/* Copy button */
.prose :deep(.copy-btn) {
  position: absolute;
  top: 0.6em;
  right: 0.6em;
  padding: 0.2em 0.6em;
  font-size: 0.75em;
  border-radius: 0.25rem;
  border: 1px solid var(--ui-border);
  background: var(--ui-bg);
  color: var(--ui-text-muted);
  cursor: pointer;
  opacity: 0;
  transition: opacity 0.15s;
}
.prose :deep(pre:hover .copy-btn) {
  opacity: 1;
}
</style>

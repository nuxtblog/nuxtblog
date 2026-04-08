<script setup lang="ts">
import { marked } from 'marked'
import hljs from 'highlight.js'

const props = defineProps<{
  content?: string
  html?: string
}>()

// Configure marked to use highlight.js via custom renderer
const renderer = new marked.Renderer()
renderer.code = ({ text, lang }: { text: string; lang?: string }) => {
  const language = lang && hljs.getLanguage(lang) ? lang : 'plaintext'
  const highlighted = hljs.highlight(text, { language }).value
  return `<pre><code class="hljs language-${language}">${highlighted}</code></pre>`
}
marked.use({ renderer })

const rendered = computed(() =>
  props.html ?? (props.content ? marked(props.content) as string : '')
)

const containerRef = ref<HTMLElement>()

// Inject copy buttons after render
function addCopyButtons() {
  if (!containerRef.value) return
  containerRef.value.querySelectorAll('pre').forEach((pre) => {
    if (pre.querySelector('.copy-btn')) return
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

watch(rendered, () => nextTick(addCopyButtons))
onMounted(() => nextTick(addCopyButtons))
</script>

<template>
  <div ref="containerRef" class="prose" v-html="rendered" />
</template>

<style>
/* highlight.js theme — GitHub-style, adapts to dark mode */
@import 'highlight.js/styles/github.css' screen and (prefers-color-scheme: light);
@import 'highlight.js/styles/github-dark.css' screen and (prefers-color-scheme: dark);
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

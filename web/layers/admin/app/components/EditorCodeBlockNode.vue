<script setup lang="ts">
import { NodeViewWrapper, NodeViewContent } from '@tiptap/vue-3'

const props = defineProps<{
  node: { attrs: { language: string | null } }
  updateAttributes: (attrs: Record<string, unknown>) => void
  extension: { options: { languageClassPrefix: string } }
}>()

const languages = [
  { label: 'Plain', value: '' },
  { label: 'JavaScript', value: 'javascript' },
  { label: 'TypeScript', value: 'typescript' },
  { label: 'Python', value: 'python' },
  { label: 'Go', value: 'go' },
  { label: 'Rust', value: 'rust' },
  { label: 'HTML', value: 'html' },
  { label: 'CSS', value: 'css' },
  { label: 'JSON', value: 'json' },
  { label: 'YAML', value: 'yaml' },
  { label: 'Bash', value: 'bash' },
  { label: 'SQL', value: 'sql' },
  { label: 'Markdown', value: 'markdown' },
  { label: 'Mermaid', value: 'mermaid' },
  { label: 'Java', value: 'java' },
  { label: 'C++', value: 'cpp' },
  { label: 'C#', value: 'csharp' },
  { label: 'PHP', value: 'php' },
  { label: 'Ruby', value: 'ruby' },
  { label: 'Swift', value: 'swift' },
  { label: 'Kotlin', value: 'kotlin' },
  { label: 'Lua', value: 'lua' },
  { label: 'Dockerfile', value: 'dockerfile' },
  { label: 'XML', value: 'xml' },
  { label: 'TOML', value: 'toml' },
]

function onLangChange(event: Event) {
  const val = (event.target as HTMLSelectElement).value
  props.updateAttributes({ language: val || null })
}
</script>

<template>
  <NodeViewWrapper class="code-block-with-lang">
    <div class="code-block-header" contenteditable="false">
      <select
        :value="node.attrs.language || ''"
        class="lang-select"
        @change="onLangChange">
        <option
          v-for="lang in languages"
          :key="lang.value"
          :value="lang.value">
          {{ lang.label }}
        </option>
      </select>
    </div>
    <pre><NodeViewContent as="code" /></pre>
  </NodeViewWrapper>
</template>

<style scoped>
.code-block-with-lang {
  position: relative;
  margin: 1em 0;
  border: 1px solid var(--ui-border);
  border-radius: 0.75rem;
  overflow: hidden;
}
.code-block-header {
  display: flex;
  justify-content: flex-end;
  padding: 0.25rem 0.5rem;
  background: var(--ui-bg-elevated);
  border-bottom: 1px solid var(--ui-border);
}
.lang-select {
  font-size: 0.75rem;
  padding: 0.125rem 0.375rem;
  border-radius: 0.25rem;
  border: 1px solid var(--ui-border);
  background: var(--ui-bg);
  color: var(--ui-text-muted);
  cursor: pointer;
  outline: none;
}
.lang-select:focus {
  border-color: var(--ui-primary);
}
.code-block-with-lang pre {
  margin: 0;
  padding: 1em;
  overflow-x: auto;
  background: transparent;
  border: none;
  border-radius: 0;
}
.code-block-with-lang pre code {
  font-size: 0.875em;
  font-family: "JetBrains Mono", "Courier New", monospace;
}
</style>

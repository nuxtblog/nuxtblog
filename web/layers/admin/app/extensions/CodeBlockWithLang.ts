import { Node, mergeAttributes, textblockTypeInputRule } from '@tiptap/core'
import { VueNodeViewRenderer } from '@tiptap/vue-3'
import EditorCodeBlockNode from '../components/EditorCodeBlockNode.vue'

export const CodeBlockWithLang = Node.create({
  name: 'codeBlock',
  group: 'block',
  content: 'text*',
  marks: '',
  code: true,
  defining: true,

  addAttributes() {
    return {
      language: { default: null },
    }
  },

  parseHTML() {
    return [
      {
        tag: 'pre',
        preserveWhitespace: 'full' as const,
        getAttrs: (el: string | HTMLElement) => {
          if (typeof el === 'string') return {}
          const code = el.querySelector('code')
          if (!code) return {}
          const cls = [...code.classList].find(c => c.startsWith('language-'))
          return { language: cls ? cls.replace('language-', '') : null }
        },
      },
    ]
  },

  renderHTML({ node, HTMLAttributes }) {
    const lang = node.attrs.language
    const codeAttrs = lang ? { class: `language-${lang}` } : {}
    return ['pre', mergeAttributes(HTMLAttributes), ['code', codeAttrs, 0]]
  },

  // @tiptap/markdown integration
  markdownTokenName: 'code',
  parseMarkdown: (token: any, helpers: any) => {
    if (token.raw?.startsWith('```') === false && token.raw?.startsWith('~~~') === false && token.codeBlockStyle !== 'indented') {
      return []
    }
    return helpers.createNode(
      'codeBlock',
      { language: token.lang || null },
      token.text ? [helpers.createTextNode(token.text)] : [],
    )
  },
  renderMarkdown: (node: any, h: any) => {
    const language = node.attrs?.language || ''
    if (!node.content) {
      return `\`\`\`${language}\n\n\`\`\``
    }
    return [`\`\`\`${language}`, h.renderChildren(node.content), '```'].join('\n')
  },

  addNodeView() {
    return VueNodeViewRenderer(EditorCodeBlockNode as any)
  },

  addCommands() {
    return {
      setCodeBlock:
        (attrs?: { language?: string }) =>
        ({ commands }) => {
          return commands.setNode(this.name, attrs)
        },
      toggleCodeBlock:
        (attrs?: { language?: string }) =>
        ({ commands }) => {
          return commands.toggleNode(this.name, 'paragraph', attrs)
        },
    }
  },

  addKeyboardShortcuts() {
    return {
      'Mod-Alt-c': () => this.editor.commands.toggleCodeBlock(),
      // Exit code block on triple Enter
      Backspace: () => {
        const { empty, $anchor } = this.editor.state.selection
        const isAtStart = $anchor.pos === $anchor.start()
        if (!empty || $anchor.parent.type.name !== this.name) return false
        if (isAtStart && !$anchor.parent.textContent.length) {
          return this.editor.commands.clearNodes()
        }
        return false
      },
      Enter: ({ editor }) => {
        const { $from } = editor.state.selection
        if ($from.parent.type.name !== this.name) return false
        // Allow normal Enter in code blocks
        return false
      },
      'Mod-Enter': () => {
        // Exit code block on Cmd/Ctrl+Enter
        if (this.editor.state.selection.$from.parent.type.name !== this.name) return false
        return this.editor.commands.exitCode()
      },
    }
  },

  addInputRules() {
    return [
      textblockTypeInputRule({
        find: /^```([a-z]*)?[\s\n]$/,
        type: this.type,
        getAttributes: (match) => ({
          language: match[1] || null,
        }),
      }),
    ]
  },
})

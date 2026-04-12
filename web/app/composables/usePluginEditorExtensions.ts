// Using `any` for extension type to avoid requiring @tiptap/core in the app layer.
// The admin layer (which has tiptap) provides the actual typed extensions.
// eslint-disable-next-line @typescript-eslint/no-explicit-any
const pluginExtensions = ref<any[]>([])

export function registerPluginExtension(ext: unknown): () => void {
  pluginExtensions.value.push(ext)
  return () => {
    const idx = pluginExtensions.value.indexOf(ext)
    if (idx >= 0) pluginExtensions.value.splice(idx, 1)
  }
}

export function usePluginEditorExtensions() {
  return { pluginExtensions: readonly(pluginExtensions) }
}

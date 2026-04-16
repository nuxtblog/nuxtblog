/**
 * Plugin I18n — centralized i18n message store and template resolution.
 *
 * Plugins declare i18n messages in plugin.yaml under the `i18n` key.
 * Text fields use {{key}} templates resolved at render time by the frontend.
 */

// Global i18n message store: pluginId → { locale → { key → message } }
const pluginMessages = new Map<string, Record<string, Record<string, string>>>()

/** Register a plugin's i18n messages. */
export function registerPluginI18n(pluginId: string, messages: Record<string, Record<string, string>>) {
  pluginMessages.set(pluginId, messages)
}

/** Resolve {{key}} template strings using the plugin's i18n messages. */
export function resolveTemplate(text: string, pluginId: string, locale: string): string {
  if (!text || !text.includes('{{')) return text
  const msgs = pluginMessages.get(pluginId)
  if (!msgs) return text
  return text.replace(/\{\{(\w+)\}\}/g, (_, key) => {
    return msgs[locale]?.[key] ?? msgs['zh']?.[key] ?? `{{${key}}}`
  })
}

/**
 * Create a composable factory for plugin Vue components.
 * Injected at module load time as `usePluginT` — no import needed in plugins.
 */
export function createPluginTComposable(pluginId: string, localeRef: Ref<string>) {
  return function usePluginT() {
    return {
      t(key: string, fallback?: string): string {
        const msgs = pluginMessages.get(pluginId)
        if (!msgs) return fallback ?? key
        return msgs[localeRef.value]?.[key] ?? msgs['zh']?.[key] ?? fallback ?? key
      },
      locale: computed(() => localeRef.value),
    }
  }
}

/** Create a per-plugin runtime i18n API for activate() scripts. */
export function createPluginI18nApi(pluginId: string, localeRef: Ref<string>) {
  return {
    t(key: string, fallback?: string): string {
      const msgs = pluginMessages.get(pluginId)
      if (!msgs) return fallback ?? key
      return msgs[localeRef.value]?.[key] ?? msgs['zh']?.[key] ?? fallback ?? key
    },
    get locale(): string {
      return localeRef.value
    },
  }
}

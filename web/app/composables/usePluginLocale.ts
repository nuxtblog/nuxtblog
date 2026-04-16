/**
 * Plugin Locale — composable for resolving plugin i18n templates in Vue components.
 */
import { resolveTemplate } from './usePluginI18n'

export function usePluginLocale() {
  const { locale } = useI18n()

  /** Resolve a contribution item's title via i18n template. */
  function t(item: { pluginId: string; title?: string }): string {
    return resolveTemplate(item.title || '', item.pluginId, locale.value)
  }

  /** Resolve all text fields of a setting field. */
  function tField(pluginId: string, field: {
    label?: string
    description?: string
    group?: string
    placeholder?: string
  }) {
    return {
      label: resolveTemplate(field.label || '', pluginId, locale.value),
      description: resolveTemplate(field.description || '', pluginId, locale.value),
      group: resolveTemplate(field.group || '', pluginId, locale.value),
      placeholder: resolveTemplate(field.placeholder || '', pluginId, locale.value),
    }
  }

  return { t, tField, locale }
}

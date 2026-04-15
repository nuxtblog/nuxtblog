import { WIDGET_REGISTRY } from '~/config/widgets'

export type { WidgetDef } from '~/config/widgets'
export { WIDGET_REGISTRY } from '~/config/widgets'

export interface WidgetConfig {
  id: string
  label: string
  enabled: boolean
  title?: string
  showRecent?: boolean
  showHot?: boolean
  maxCount?: number
  pluginSettings?: Record<string, unknown>
  // Plugin widget extension fields
  isPlugin?: boolean
  pluginId?: string
  component?: string
  module?: string
}

/** Derived from WIDGET_REGISTRY — no duplication */
export const WIDGET_DEFAULTS: WidgetConfig[] = WIDGET_REGISTRY.map(def => ({
  id: def.id,
  label: def.label,
  enabled: def.defaultEnabled,
  title: def.defaultTitle,
  ...def.fieldDefaults,
}))

// Sidebar provides this key so child widgets read from the right option store entry
export const WIDGET_OPTION_KEY = Symbol('widgetOptionKey')

export const useWidgetConfig = () => {
  const optionStore = useOptionsStore()
  const optionKey = inject(WIDGET_OPTION_KEY, 'blog_sidebar_widgets')

  const getWidgetConfig = (id: string): WidgetConfig => {
    const saved = optionStore.getJSON<WidgetConfig[]>(optionKey, WIDGET_DEFAULTS)
    const def = WIDGET_DEFAULTS.find(w => w.id === id) ?? { id, label: id, enabled: false }
    const found = saved.find(w => w.id === id)
    // Always use label from def (i18n key). Saved data may contain old Chinese strings.
    return { ...def, ...found, label: def.label }
  }

  return { getWidgetConfig, WIDGET_DEFAULTS }
}

/** Merged widget defaults: built-in + plugin-registered sidebar widgets. */
export function useAllWidgetDefaults() {
  const contributionsStore = usePluginContributionsStore()
  const pluginViews = contributionsStore.getViewItems('public:sidebar-widget')

  return computed(() => {
    const pluginWidgets: WidgetConfig[] = pluginViews.value
      .filter((v: { component?: string; module?: string }) => v.component && v.module)
      .map((v: { pluginId: string; id: string; title: string; component?: string; module?: string; settings?: Array<{ key: string; default?: unknown }> }) => {
        const pluginSettings: Record<string, unknown> = {}
        if (v.settings) {
          for (const s of v.settings) {
            if (s.default !== undefined) pluginSettings[s.key] = s.default
          }
        }
        return {
          id: `plugin:${v.pluginId}:${v.id}`,
          label: v.title || v.id,
          enabled: false,
          isPlugin: true,
          pluginId: v.pluginId,
          component: v.component,
          module: v.module,
          maxCount: 5,
          ...(Object.keys(pluginSettings).length > 0 ? { pluginSettings } : {}),
        }
      })
    return [...WIDGET_DEFAULTS, ...pluginWidgets]
  })
}

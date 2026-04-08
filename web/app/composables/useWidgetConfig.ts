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

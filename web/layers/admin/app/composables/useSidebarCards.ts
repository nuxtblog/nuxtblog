interface SidebarCardPrefs {
  order: string[]
  hidden: string[]
  collapsed: string[]
}

const DEFAULT_ORDER = [
  'publish', 'display', 'layout', 'categories',
  'tags', 'downloads', 'featured-image', 'meta-fields', 'seo',
]

const DEFAULT_PREFS: SidebarCardPrefs = {
  order: [...DEFAULT_ORDER],
  hidden: [],
  collapsed: [],
}

export function useSidebarCards() {
  const prefs = useLocalStorage<SidebarCardPrefs>('nuxtblog:editor:sidebar-cards', structuredClone(DEFAULT_PREFS))

  const allCards = computed(() => {
    // Ensure any new default cards are included
    const known = new Set(prefs.value.order)
    const merged = [...prefs.value.order]
    for (const id of DEFAULT_ORDER) {
      if (!known.has(id)) merged.push(id)
    }
    return merged
  })

  const visibleCards = computed(() =>
    allCards.value.filter(id => !prefs.value.hidden.includes(id))
  )

  const isCollapsed = (id: string) => prefs.value.collapsed.includes(id)

  const isVisible = (id: string) => !prefs.value.hidden.includes(id)

  const toggleCollapsed = (id: string) => {
    const idx = prefs.value.collapsed.indexOf(id)
    if (idx >= 0) prefs.value.collapsed.splice(idx, 1)
    else prefs.value.collapsed.push(id)
  }

  const toggleVisibility = (id: string) => {
    const idx = prefs.value.hidden.indexOf(id)
    if (idx >= 0) prefs.value.hidden.splice(idx, 1)
    else prefs.value.hidden.push(id)
  }

  const setOrder = (ids: string[]) => {
    prefs.value.order = ids
  }

  const resetDefaults = () => {
    prefs.value = structuredClone(DEFAULT_PREFS)
  }

  return {
    prefs,
    allCards,
    visibleCards,
    isCollapsed,
    isVisible,
    toggleCollapsed,
    toggleVisibility,
    setOrder,
    resetDefaults,
  }
}

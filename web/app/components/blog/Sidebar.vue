<template>
  <aside class="lg:col-span-3">
    <div class="sticky top-20 space-y-6">
      <ClientOnly><ContributionSlot name="public:sidebar-top" /></ClientOnly>

      <template v-for="widget in activeWidgets" :key="widget.id">
        <!-- Plugin widget -->
        <ClientOnly v-if="widget.isPlugin">
          <component :is="getPluginWidgetComponent(widget)" v-if="getPluginWidgetComponent(widget)" :config="widget" :context="sidebarContext" />
        </ClientOnly>
        <!-- Built-in widget (needs ClientOnly) -->
        <ClientOnly v-else-if="CLIENT_ONLY_WIDGETS.has(widget.id)">
          <component :is="builtinWidgetMap[widget.id]" :config="widget" />
        </ClientOnly>
        <!-- Built-in widget (regular) -->
        <component v-else-if="builtinWidgetMap[widget.id]" :is="builtinWidgetMap[widget.id]" :config="widget" />
      </template>

      <ClientOnly><ContributionSlot name="public:sidebar-bottom" /></ClientOnly>
    </div>
  </aside>
</template>

<script setup lang="ts">
import type { Component } from 'vue'
import { WIDGET_DEFAULTS, type WidgetConfig } from '~/composables/useWidgetConfig'
import { getWidgetsByContext } from '~/config/widgets'
import { createPluginAsyncComponent } from '~/composables/usePluginComponents'
import { type SidebarContext, SIDEBAR_CONTEXT_KEY } from '~/composables/useSidebarContext'

const props = defineProps<{ optionKey?: string; context?: SidebarContext }>()
const optionStore = useOptionsStore()
const resolvedKey = props.optionKey ?? 'blog_sidebar_widgets'

// Provide page context to all child widgets (built-in + plugin)
function derivePageType(key: string): SidebarContext['pageType'] {
  if (key === 'homepage_sidebar_widgets') return 'homepage'
  return 'blog'
}
const sidebarContext = computed<SidebarContext>(() => props.context ?? { pageType: derivePageType(resolvedKey) })
provide(SIDEBAR_CONTEXT_KEY, sidebarContext)

// Built-in widget component map
const builtinWidgetMap: Record<string, string | Component> = {
  user_box: resolveComponent('WidgetUserBox'),
  search: resolveComponent('WidgetSearchBar'),
  author: resolveComponent('WidgetAuthorBox'),
  tags: resolveComponent('WidgetTagBox'),
  latest_posts: resolveComponent('WidgetLatestPosts'),
  latest_comments: resolveComponent('WidgetLatestComments'),
  recommend: resolveComponent('WidgetRecommend'),
  featured: resolveComponent('WidgetFeaturedContent'),
  random_posts: resolveComponent('WidgetRandomPosts'),
  toc: resolveComponent('WidgetToc'),
  downloads: resolveComponent('WidgetDownloads'),
}

// Widgets that need ClientOnly wrapping
const CLIENT_ONLY_WIDGETS = new Set(['random_posts'])

// Plugin widget component cache
const pluginComponentCache = new Map<string, Component>()
function getPluginWidgetComponent(widget: WidgetConfig) {
  if (!widget.pluginId || !widget.component || !widget.module) return null
  const key = `${widget.pluginId}:${widget.component}:${widget.module}`
  if (!pluginComponentCache.has(key)) {
    pluginComponentCache.set(key, createPluginAsyncComponent(widget.pluginId, widget.component, widget.module))
  }
  return pluginComponentCache.get(key)!
}

// Pick the right fallback defaults based on which option key is being used
const fallbackDefaults = computed((): WidgetConfig[] => {
  if (resolvedKey === 'homepage_sidebar_widgets') {
    return WIDGET_DEFAULTS.filter(w => getWidgetsByContext('homepage').some(d => d.id === w.id))
  }
  return WIDGET_DEFAULTS
})

const activeWidgets = computed(() =>
  optionStore.getJSON<WidgetConfig[]>(resolvedKey, fallbackDefaults.value).filter(w => w.enabled)
)
</script>

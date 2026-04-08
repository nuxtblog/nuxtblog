<template>
  <aside class="lg:col-span-3">
    <div class="sticky top-20 space-y-6">
      <template v-for="widget in activeWidgets" :key="widget.id">
        <WidgetUserBox v-if="widget.id === 'user_box'" :config="widget" />
        <WidgetSearchBar v-else-if="widget.id === 'search'" :config="widget" />
        <WidgetAuthorBox v-else-if="widget.id === 'author'" :config="widget" />
        <WidgetTagBox v-else-if="widget.id === 'tags'" :config="widget" />
        <WidgetLatestPosts v-else-if="widget.id === 'latest_posts'" :config="widget" />
        <WidgetLatestComments v-else-if="widget.id === 'latest_comments'" :config="widget" />
        <WidgetRecommend v-else-if="widget.id === 'recommend'" :config="widget" />
        <WidgetFeaturedContent v-else-if="widget.id === 'featured'" :config="widget" />
        <ClientOnly v-else-if="widget.id === 'random_posts'">
          <WidgetRandomPosts :config="widget" />
        </ClientOnly>
        <WidgetToc v-else-if="widget.id === 'toc'" :config="widget" />
        <WidgetDownloads v-else-if="widget.id === 'downloads'" :config="widget" />
      </template>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { WIDGET_DEFAULTS, type WidgetConfig } from '~/composables/useWidgetConfig'
import { getWidgetsByContext } from '~/config/widgets'

const props = defineProps<{ optionKey?: string }>()
const optionStore = useOptionsStore()

const resolvedKey = props.optionKey ?? 'blog_sidebar_widgets'

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

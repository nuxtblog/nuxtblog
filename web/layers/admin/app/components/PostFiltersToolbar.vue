<script setup lang="ts">
interface TaxonomyItem {
  id: number
  term: { id: number; name: string; slug: string }
  taxonomy: string
  post_count: number
  parent_id?: number
}

const props = defineProps<{
  categories: TaxonomyItem[]
  authors: { id: number; username: string; display_name: string }[]
  filterStatus: string
  totalPosts: number
}>()

const searchKeyword = defineModel<string>('searchKeyword', { required: true })
const filterCategory = defineModel<string | undefined>('filterCategory', { required: true })
const filterAuthor = defineModel<string | undefined>('filterAuthor', { required: true })
const sortBy = defineModel<string>('sortBy', { required: true })
const pageSize = defineModel<number>('pageSize', { required: true })
const viewMode = defineModel<'list' | 'grid'>('viewMode', { required: true })

const emit = defineEmits<{
  (e: 'clearTrash'): void
}>()

const { t } = useI18n()

const catSearch = ref('')
const authorSearch = ref('')
const catPopoverOpen = ref(false)
const authorPopoverOpen = ref(false)

const categoryOptions = computed(() => [
  { label: t('admin.posts.all_categories'), value: undefined as string | undefined },
  ...props.categories.map(c => ({ label: c.term.name, value: String(c.id) })),
])

const authorOptions = computed(() => [
  { label: t('admin.posts.all_authors'), value: undefined as string | undefined },
  ...props.authors.map(u => ({ label: u.display_name || u.username, value: String(u.id) })),
])

const filteredCategoryOptions = computed(() => {
  if (!catSearch.value) return categoryOptions.value
  const q = catSearch.value.toLowerCase()
  return categoryOptions.value.filter(o => o.label.toLowerCase().includes(q))
})

const filteredAuthorOptions = computed(() => {
  if (!authorSearch.value) return authorOptions.value
  const q = authorSearch.value.toLowerCase()
  return authorOptions.value.filter(o => o.label.toLowerCase().includes(q))
})

const selectedCategoryLabel = computed(() =>
  categoryOptions.value.find(o => o.value === filterCategory.value)?.label || t('admin.posts.all_categories'),
)

const selectedAuthorLabel = computed(() =>
  authorOptions.value.find(o => o.value === filterAuthor.value)?.label || t('admin.posts.all_authors'),
)
</script>

<template>
  <div class="flex flex-col gap-3 mb-4">
    <div class="flex flex-wrap items-center gap-3">
      <!-- Search -->
      <UInput
        v-model="searchKeyword"
        :placeholder="$t('admin.posts.search_placeholder')"
        leading-icon="i-tabler-search"
        class="w-56"
        size="sm">
        <template v-if="searchKeyword" #trailing>
          <UButton icon="i-tabler-x" color="neutral" variant="ghost" size="xs" @click="searchKeyword = ''" />
        </template>
      </UInput>

      <!-- Category filter -->
      <UPopover v-model:open="catPopoverOpen" :popper="{ placement: 'bottom-start' }">
        <UButton
          color="neutral" variant="outline" size="sm"
          trailing-icon="i-tabler-chevron-down"
          :class="['w-40 justify-between font-normal', filterCategory ? 'text-highlighted' : 'text-muted']">
          {{ selectedCategoryLabel }}
        </UButton>
        <template #content>
          <div class="p-2 w-48">
            <UInput v-model="catSearch" :placeholder="$t('admin.posts.search_categories')" size="sm" leading-icon="i-tabler-search" autofocus class="mb-2" />
            <div class="max-h-52 overflow-y-auto space-y-0.5">
              <button
                v-for="opt in filteredCategoryOptions"
                :key="String(opt.value)"
                class="w-full text-left px-2 py-1.5 text-sm rounded-md transition-colors hover:bg-elevated"
                :class="filterCategory === opt.value ? 'text-primary font-medium bg-primary/5' : 'text-default'"
                @click="filterCategory = opt.value; catPopoverOpen = false; catSearch = ''">
                {{ opt.label }}
              </button>
            </div>
          </div>
        </template>
      </UPopover>

      <!-- Author filter -->
      <UPopover v-model:open="authorPopoverOpen" :popper="{ placement: 'bottom-start' }">
        <UButton
          color="neutral" variant="outline" size="sm"
          trailing-icon="i-tabler-chevron-down"
          :class="['w-40 justify-between font-normal', filterAuthor ? 'text-highlighted' : 'text-muted']">
          {{ selectedAuthorLabel }}
        </UButton>
        <template #content>
          <div class="p-2 w-48">
            <UInput v-model="authorSearch" :placeholder="$t('admin.posts.search_authors')" size="sm" leading-icon="i-tabler-search" autofocus class="mb-2" />
            <div class="max-h-52 overflow-y-auto space-y-0.5">
              <button
                v-for="opt in filteredAuthorOptions"
                :key="String(opt.value)"
                class="w-full text-left px-2 py-1.5 text-sm rounded-md transition-colors hover:bg-elevated"
                :class="filterAuthor === opt.value ? 'text-primary font-medium bg-primary/5' : 'text-default'"
                @click="filterAuthor = opt.value; authorPopoverOpen = false; authorSearch = ''">
                {{ opt.label }}
              </button>
            </div>
          </div>
        </template>
      </UPopover>

      <!-- Sort -->
      <USelect
        v-model="sortBy"
        :items="[
          { label: $t('admin.posts.sort_created_desc'), value: 'created_desc' },
          { label: $t('admin.posts.sort_created_asc'), value: 'created_asc' },
          { label: $t('admin.posts.sort_updated_desc'), value: 'updated_desc' },
          { label: $t('admin.posts.sort_views_desc'), value: 'views_desc' },
        ]"
        class="w-32" size="sm" />

      <!-- Page size -->
      <USelect
        v-model="pageSize"
        :items="[
          { label: $t('common.items_per_page', { n: 10 }), value: 10 },
          { label: $t('common.items_per_page', { n: 20 }), value: 20 },
          { label: $t('common.items_per_page', { n: 50 }), value: 50 },
        ]"
        class="w-28" size="sm" />

      <!-- View switch + clear trash -->
      <div class="ml-auto flex items-center gap-2">
        <UButton
          v-if="filterStatus === 'trash' && totalPosts > 0"
          color="error" variant="soft" size="sm" icon="i-tabler-trash-x"
          @click="emit('clearTrash')">
          {{ $t('admin.posts.clear_trash') }}
        </UButton>
        <ViewModeSwitcher
          v-model="viewMode"
          :modes="[
            { value: 'grid', title: $t('admin.posts.view_mode_grid') },
            { value: 'list', title: $t('admin.posts.view_mode_list') },
          ]" />
      </div>
    </div>
  </div>
</template>

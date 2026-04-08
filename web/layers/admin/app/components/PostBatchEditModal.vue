<template>
  <UModal v-model:open="open" :ui="{ content: 'max-w-lg' }">
    <template #content>
      <div class="flex flex-col max-h-[90vh]">

        <!-- Modal header -->
        <div class="px-5 pt-5 pb-4 border-b border-default shrink-0">
          <h3 class="font-semibold text-highlighted">{{ $t('admin.posts.batch_edit_title') }}</h3>
          <p class="text-xs text-muted mt-0.5">{{ $t('admin.posts.batch_edit_hint', { n: postIds.length }) }}</p>
        </div>

        <!-- Scrollable fields -->
        <div class="overflow-y-auto flex-1 divide-y divide-default">

          <!-- ── 状态 ── -->
          <div class="px-5 py-4 flex items-start gap-3">
            <UCheckbox v-model="enableStatus" class="mt-0.5 shrink-0" />
            <div class="flex-1 min-w-0" :class="!enableStatus && 'opacity-40 pointer-events-none'">
              <p class="text-sm font-medium text-highlighted mb-2">{{ $t('admin.posts.field_status') }}</p>
              <USelect
                v-model="statusVal"
                :items="statusOptions"
                size="sm"
                class="w-full" />
            </div>
          </div>

          <!-- ── 封面图 ── -->
          <div class="px-5 py-4 flex items-start gap-3">
            <UCheckbox v-model="enableCover" class="mt-0.5 shrink-0" />
            <div class="flex-1 min-w-0">
              <div class="flex items-center justify-between mb-2">
                <p class="text-sm font-medium text-highlighted" :class="!enableCover && 'opacity-40'">
                  {{ $t('admin.posts.field_cover') }}
                </p>
                <UButton
                  v-if="enableCover && coverImgId"
                  color="error"
                  variant="ghost"
                  size="xs"
                  icon="i-tabler-x"
                  @click="clearCover">
                  {{ $t('admin.posts.batch_clear_cover') }}
                </UButton>
              </div>
              <div :class="!enableCover && 'opacity-40 pointer-events-none'">
                <FeaturedImagePicker v-model:img-id="coverImgId" v-model:img-url="coverImgUrl" />
              </div>
              <p v-if="enableCover && clearFeaturedImg && !coverImgId" class="text-xs text-warning mt-2 flex items-center gap-1">
                <UIcon name="i-tabler-alert-triangle" class="size-3.5" />
                {{ $t('admin.posts.batch_cover_will_clear') }}
              </p>
            </div>
          </div>

          <!-- ── 分类 ── -->
          <div class="px-5 py-4 flex items-start gap-3">
            <UCheckbox v-model="enableCategories" class="mt-0.5 shrink-0" />
            <div class="flex-1 min-w-0">
              <p class="text-sm font-medium text-highlighted mb-1" :class="!enableCategories && 'opacity-40'">
                {{ $t('admin.posts.field_categories') }}
              </p>
              <p class="text-xs text-muted mb-2" :class="!enableCategories && 'opacity-40'">
                {{ $t('admin.posts.batch_taxonomy_hint') }}
              </p>
              <div :class="!enableCategories && 'opacity-40 pointer-events-none'">
                <div v-if="categories.length === 0" class="text-xs text-muted italic">
                  {{ $t('admin.posts.no_categories') }}
                </div>
                <div v-else class="max-h-36 overflow-y-auto rounded-md border border-default bg-default">
                  <label
                    v-for="cat in categories"
                    :key="cat.id"
                    class="flex items-center gap-2.5 px-3 py-2 cursor-pointer hover:bg-muted/50 transition-colors border-b border-default/50 last:border-0">
                    <UCheckbox
                      :model-value="selectedCatIds.includes(cat.id)"
                      @update:model-value="toggleId(selectedCatIds, cat.id)" />
                    <span class="text-sm truncate">{{ cat.term.name }}</span>
                  </label>
                </div>
              </div>
            </div>
          </div>

          <!-- ── 标签 ── -->
          <div class="px-5 py-4 flex items-start gap-3">
            <UCheckbox v-model="enableTags" class="mt-0.5 shrink-0" />
            <div class="flex-1 min-w-0">
              <p class="text-sm font-medium text-highlighted mb-1" :class="!enableTags && 'opacity-40'">
                {{ $t('admin.posts.field_tags') }}
              </p>
              <p class="text-xs text-muted mb-2" :class="!enableTags && 'opacity-40'">
                {{ $t('admin.posts.batch_taxonomy_hint') }}
              </p>
              <div :class="!enableTags && 'opacity-40 pointer-events-none'">
                <div v-if="tags.length === 0" class="text-xs text-muted italic">
                  {{ $t('admin.posts.no_tags') }}
                </div>
                <template v-else>
                  <!-- Search -->
                  <UInput
                    v-model="tagSearch"
                    :placeholder="$t('common.search')"
                    leading-icon="i-tabler-search"
                    size="xs"
                    class="w-full mb-2" />
                  <!-- Tag chips -->
                  <div class="flex flex-wrap gap-1.5 max-h-28 overflow-y-auto p-1">
                    <button
                      v-for="tag in filteredTags"
                      :key="tag.id"
                      type="button"
                      class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-medium border transition-all"
                      :class="selectedTagIds.includes(tag.id)
                        ? 'bg-primary text-white border-primary'
                        : 'bg-default border-default text-muted hover:border-primary/50 hover:text-highlighted'"
                      @click="toggleId(selectedTagIds, tag.id)">
                      <UIcon
                        :name="selectedTagIds.includes(tag.id) ? 'i-tabler-check' : 'i-tabler-tag'"
                        class="size-3" />
                      {{ tag.term.name }}
                    </button>
                  </div>
                </template>
              </div>
            </div>
          </div>

          <!-- ── 作者（仅管理员）── -->
          <div v-if="isAdmin" class="px-5 py-4 flex items-start gap-3">
            <UCheckbox v-model="enableAuthor" class="mt-0.5 shrink-0" />
            <div class="flex-1 min-w-0" :class="!enableAuthor && 'opacity-40 pointer-events-none'">
              <p class="text-sm font-medium text-highlighted mb-2">{{ $t('admin.posts.field_author') }}</p>
              <USelect
                v-if="authors.length > 0"
                v-model="authorVal"
                :items="authorItems"
                size="sm"
                class="w-full" />
              <p v-else class="text-xs text-muted italic">{{ $t('admin.posts.no_authors') }}</p>
            </div>
          </div>

        </div>

        <!-- Footer -->
        <div class="px-5 py-4 border-t border-default shrink-0 flex items-center justify-between gap-3">
          <p v-if="!hasAnyEnabled" class="text-xs text-muted">
            {{ $t('admin.posts.batch_edit_select_hint') }}
          </p>
          <div v-else class="flex flex-wrap gap-1">
            <UBadge v-if="enableStatus" color="neutral" variant="soft" size="xs" icon="i-tabler-toggle-right">
              {{ $t('admin.posts.field_status') }}
            </UBadge>
            <UBadge v-if="enableCover" color="neutral" variant="soft" size="xs" icon="i-tabler-photo">
              {{ $t('admin.posts.field_cover') }}
            </UBadge>
            <UBadge v-if="enableCategories" color="neutral" variant="soft" size="xs" icon="i-tabler-folder">
              {{ $t('admin.posts.field_categories') }}
            </UBadge>
            <UBadge v-if="enableTags" color="neutral" variant="soft" size="xs" icon="i-tabler-tag">
              {{ $t('admin.posts.field_tags') }}
            </UBadge>
            <UBadge v-if="isAdmin && enableAuthor" color="neutral" variant="soft" size="xs" icon="i-tabler-user">
              {{ $t('admin.posts.field_author') }}
            </UBadge>
          </div>
          <div class="flex items-center gap-2 ml-auto shrink-0">
            <UButton color="neutral" variant="ghost" size="sm" @click="open = false">{{ $t('common.cancel') }}</UButton>
            <UButton
              color="primary"
              size="sm"
              :loading="saving"
              :disabled="!hasAnyEnabled"
              icon="i-tabler-check"
              @click="apply">
              {{ $t('admin.posts.batch_apply', { n: postIds.length }) }}
            </UButton>
          </div>
        </div>

      </div>
    </template>
  </UModal>
</template>

<script setup lang="ts">
interface TaxonomyItem {
  id: number
  term: { id: number; name: string; slug: string }
  taxonomy: string
}

interface UserItem {
  id: number
  username: string
  display_name: string
}

const props = defineProps<{
  postIds: number[]
  categories: TaxonomyItem[]
  tags: TaxonomyItem[]
  authors: UserItem[]
  isAdmin?: boolean
}>()

const emit = defineEmits<{
  applied: []
}>()

const open = defineModel<boolean>('open', { required: true })

const { apiFetch } = useApiFetch()
const toast = useToast()
const { t } = useI18n()

// ── Status ────────────────────────────────────────────────────────────────────
const enableStatus = ref(false)
const statusVal = ref<number>(1)

const statusOptions = computed(() => [
  { label: t('admin.posts.status_draft'),     value: 1 },
  { label: t('admin.posts.status_published'), value: 2 },
  { label: t('admin.posts.status_private'),   value: 3 },
  { label: t('admin.posts.status_archived'),  value: 4 },
])

// ── Cover ─────────────────────────────────────────────────────────────────────
const enableCover = ref(false)
const coverImgId = ref<number | undefined>(undefined)
const coverImgUrl = ref('')
const clearFeaturedImg = ref(false)

const clearCover = () => {
  coverImgId.value = undefined
  coverImgUrl.value = ''
  clearFeaturedImg.value = true
}

watch(coverImgId, (v) => { if (v) clearFeaturedImg.value = false })

// ── Categories ────────────────────────────────────────────────────────────────
const enableCategories = ref(false)
const selectedCatIds = ref<number[]>([])

// ── Tags ──────────────────────────────────────────────────────────────────────
const enableTags = ref(false)
const selectedTagIds = ref<number[]>([])
const tagSearch = ref('')

const filteredTags = computed(() => {
  if (!tagSearch.value.trim()) return props.tags
  const q = tagSearch.value.toLowerCase()
  return props.tags.filter(t => t.term.name.toLowerCase().includes(q))
})

// ── Author ────────────────────────────────────────────────────────────────────
const enableAuthor = ref(false)
const authorVal = ref<number | undefined>(undefined)

const authorItems = computed(() =>
  props.authors.map(u => ({ label: u.display_name || u.username, value: u.id }))
)

// ── Helpers ───────────────────────────────────────────────────────────────────
const toggleId = (list: number[], id: number) => {
  const idx = list.indexOf(id)
  if (idx > -1) list.splice(idx, 1)
  else list.push(id)
}

const hasAnyEnabled = computed(() =>
  enableStatus.value ||
  enableCover.value ||
  enableCategories.value ||
  enableTags.value ||
  (props.isAdmin && enableAuthor.value)
)

// ── Reset on open ─────────────────────────────────────────────────────────────
watch(open, (v) => {
  if (!v) return
  enableStatus.value = false
  enableCover.value = false
  enableCategories.value = false
  enableTags.value = false
  enableAuthor.value = false
  statusVal.value = 1
  coverImgId.value = undefined
  coverImgUrl.value = ''
  clearFeaturedImg.value = false
  selectedCatIds.value = []
  selectedTagIds.value = []
  tagSearch.value = ''
  authorVal.value = undefined
})

// ── Apply ─────────────────────────────────────────────────────────────────────
const saving = ref(false)

const apply = async () => {
  if (!hasAnyEnabled.value || props.postIds.length === 0) return
  saving.value = true
  try {
    const body: Record<string, any> = { ids: props.postIds }

    if (enableStatus.value) {
      body.status = statusVal.value
    }
    if (enableCover.value) {
      body.featured_img_id = clearFeaturedImg.value ? 0 : (coverImgId.value ?? null)
    }
    // Merge selected category and tag IDs into one term_taxonomy_ids array.
    // Each enabled field independently contributes its IDs; if both are enabled,
    // the union is sent. If one is disabled, only the enabled set is sent.
    if (enableCategories.value || enableTags.value) {
      const ids: number[] = []
      if (enableCategories.value) ids.push(...selectedCatIds.value)
      if (enableTags.value) ids.push(...selectedTagIds.value)
      body.term_taxonomy_ids = ids
    }
    if (props.isAdmin && enableAuthor.value && authorVal.value) {
      body.author_id = authorVal.value
    }

    const res = await apiFetch<{ affected: number }>('/posts/batch', {
      method: 'PATCH',
      body,
    })
    toast.add({
      title: t('admin.posts.batch_edit_success', { n: res.affected }),
      icon: 'i-tabler-circle-check',
      color: 'success',
    })
    open.value = false
    emit('applied')
  } catch (err: any) {
    toast.add({
      title: err?.message || t('admin.posts.batch_edit_failed'),
      color: 'error',
      icon: 'i-tabler-alert-circle',
    })
  } finally {
    saving.value = false
  }
}
</script>

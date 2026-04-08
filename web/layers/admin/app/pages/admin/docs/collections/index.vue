<template>
  <AdminPageContainer>
    <AdminPageHeader :title="$t('admin.docs.collections_title')" :subtitle="$t('admin.docs.collections_subtitle')">
      <template #actions>
        <UButton as="NuxtLink" to="/admin/docs" variant="outline" color="neutral" icon="i-tabler-arrow-left" size="sm">
          {{ $t('admin.docs.title') }}
        </UButton>
      </template>
    </AdminPageHeader>

    <AdminPageContent>
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <!-- 左侧：新建/编辑表单 -->
        <div class="lg:col-span-1">
          <UCard>
            <template #header>
              <h2 class="text-base font-semibold text-highlighted">
                {{ isEditing ? $t('admin.docs.edit_collection') : $t('admin.docs.new_collection') }}
              </h2>
            </template>

            <form class="space-y-4" @submit.prevent="submitForm">
              <UFormField :label="$t('admin.docs.collection_name')" required>
                <UInput
                  v-model="formTitle"
                  required
                  maxlength="100"
                  :placeholder="$t('admin.docs.collection_name')"
                  class="w-full"
                />
              </UFormField>

              <UFormField :label="$t('admin.docs.collection_slug')" required>
                <UInput
                  v-model="formSlug"
                  maxlength="100"
                  :placeholder="$t('admin.docs.collection_slug')"
                  class="w-full"
                />
              </UFormField>

              <UFormField :label="$t('admin.docs.collection_description')">
                <UTextarea
                  v-model="formDescription"
                  :rows="3"
                  maxlength="255"
                  :placeholder="$t('admin.docs.collection_description')"
                  class="w-full"
                />
                <template #hint>
                  {{ formDescription?.length || 0 }} / 255
                </template>
              </UFormField>

              <UFormField :label="$t('common.status')">
                <USelect v-model="formStatus" :items="statusOptions" class="w-full" />
              </UFormField>

              <UFormField label="Sort Order">
                <UInput v-model.number="formSortOrder" type="number" class="w-full" />
              </UFormField>

              <UFormField label="Locale">
                <UInput v-model="formLocale" placeholder="zh" class="w-full" />
              </UFormField>

              <div class="flex gap-3 pt-2">
                <UButton
                  type="submit"
                  color="primary"
                  :loading="formSubmitting"
                  class="flex-1"
                >
                  {{ isEditing ? $t('common.save') : $t('admin.docs.new_collection') }}
                </UButton>
                <UButton
                  v-if="isEditing"
                  type="button"
                  color="neutral"
                  variant="outline"
                  @click="resetForm"
                >
                  {{ $t('common.cancel') }}
                </UButton>
              </div>
            </form>
          </UCard>
        </div>

        <!-- 右侧：合集列表 -->
        <div class="lg:col-span-2">
          <UCard :ui="{ body: 'p-0' }">
            <!-- 工具栏 -->
            <div class="p-4 border-b border-default flex items-center justify-between gap-4">
              <UInput
                v-model="searchQuery"
                :placeholder="$t('admin.docs.search_placeholder')"
                leading-icon="i-tabler-search"
                class="flex-1"
                size="sm"
              />
              <span class="text-sm text-muted shrink-0">
                {{ $t('common.total', { count: filteredCollections.length }) }}
              </span>
            </div>

            <!-- 加载状态 -->
            <div v-if="loading" class="p-4 space-y-3">
              <div
                v-for="i in 5" :key="i"
                class="flex items-center gap-4 p-4 bg-default border border-default rounded-md">
                <div class="flex-1 space-y-1.5">
                  <div class="flex items-center gap-3">
                    <USkeleton :class="`h-4 w-${i % 2 === 0 ? '32' : '24'}`" />
                    <USkeleton class="h-3 w-24" />
                    <USkeleton class="h-5 w-16 rounded-full" />
                  </div>
                  <USkeleton class="h-3 w-1/2" />
                </div>
              </div>
            </div>

            <!-- 列表 -->
            <div v-else-if="filteredCollections.length > 0" class="p-4 space-y-3">
              <div
                v-for="item in filteredCollections" :key="item.id"
                class="flex items-center group gap-4 p-4 bg-default border border-default rounded-md hover:shadow-sm transition-all"
                :class="editingId === item.id ? 'border-primary/50 bg-primary/5' : ''">
                <div class="flex-1 min-w-0">
                  <div class="flex items-center gap-3">
                    <h3 class="text-sm font-semibold text-highlighted">{{ item.title }}</h3>
                    <span class="text-xs text-muted font-mono">({{ item.slug }})</span>
                    <UBadge
                      v-if="item.doc_count != null"
                      :label="`${item.doc_count} 篇`"
                      size="sm"
                      color="primary"
                      variant="soft" />
                    <UBadge
                      :label="item.status === 2 ? $t('admin.docs.status_published') : $t('admin.docs.status_draft')"
                      :color="item.status === 2 ? 'success' : 'neutral'"
                      size="sm"
                      variant="soft" />
                    <span class="text-xs text-muted">排序: {{ item.sort_order }}</span>
                  </div>
                  <p v-if="item.description" class="text-sm text-muted mt-1.5 line-clamp-1">
                    {{ item.description }}
                  </p>
                </div>

                <div class="shrink-0 opacity-0 group-hover:opacity-100 transition-opacity flex items-center gap-1">
                  <UButton
                    icon="i-tabler-pencil"
                    color="neutral"
                    variant="ghost"
                    size="xs"
                    square
                    @click="startEdit(item)" />
                  <UButton
                    icon="i-tabler-trash"
                    color="error"
                    variant="ghost"
                    size="xs"
                    square
                    @click="openDeleteModal(item)" />
                </div>
              </div>
            </div>

            <!-- 空状态 -->
            <div v-else class="flex flex-col items-center justify-center py-16">
              <UIcon name="i-tabler-folders-off" class="size-16 text-muted mb-4" />
              <h3 class="text-lg font-medium text-highlighted mb-1">{{ $t('admin.docs.no_collections') }}</h3>
            </div>
          </UCard>

          <!-- 统计信息 -->
          <UCard v-if="!loading && collections.length > 0" class="mt-4">
            <div class="grid grid-cols-3 gap-4 text-center">
              <div>
                <div class="text-2xl font-semibold text-highlighted">{{ collections.length }}</div>
                <div class="text-sm text-muted">合集总数</div>
              </div>
              <div>
                <div class="text-2xl font-semibold text-highlighted">
                  {{ collections.filter(c => c.status === 2).length }}
                </div>
                <div class="text-sm text-muted">已发布</div>
              </div>
              <div>
                <div class="text-2xl font-semibold text-highlighted">
                  {{ collections.filter(c => c.status === 1).length }}
                </div>
                <div class="text-sm text-muted">草稿</div>
              </div>
            </div>
          </UCard>
        </div>
      </div>

      <!-- 删除确认弹窗 -->
      <UModal v-model:open="showDeleteModal">
        <template #content>
          <div class="p-6">
            <div class="flex items-center gap-3 mb-4">
              <div class="size-10 rounded-full bg-error/10 flex items-center justify-center shrink-0">
                <UIcon name="i-tabler-alert-triangle" class="size-5 text-error" />
              </div>
              <div>
                <h3 class="font-semibold text-highlighted">{{ $t('admin.docs.delete_collection_title') }}</h3>
                <p class="text-sm text-muted mt-0.5">{{ $t('admin.docs.delete_collection_desc') }}</p>
              </div>
            </div>
            <div class="flex justify-end gap-2 mt-6">
              <UButton color="neutral" variant="outline" @click="showDeleteModal = false">
                {{ $t('common.cancel') }}
              </UButton>
              <UButton color="error" :loading="deleteLoading" @click="confirmDelete">
                {{ $t('common.delete') }}
              </UButton>
            </div>
          </div>
        </template>
      </UModal>
    </AdminPageContent>
  </AdminPageContainer>
</template>

<script setup lang="ts">
import type { DocCollectionItem } from '~/types/api/doc'

definePageMeta({ layout: 'admin' })

const { t } = useI18n()
const toast = useToast()
const docApi = useDocApi()

// loading
const rawLoading = ref(true)
const loading = useMinLoading(rawLoading)

// state
const collections = ref<DocCollectionItem[]>([])
const searchQuery  = ref('')

// form state
const editingId      = ref<number | null>(null)
const formSlug       = ref('')
const formTitle      = ref('')
const formDescription = ref('')
const formStatus     = ref<1 | 2>(2)
const formSortOrder  = ref(0)
const formLocale     = ref('zh')
const formSubmitting = ref(false)

// delete modal
const showDeleteModal = ref(false)
const deleteTarget    = ref<DocCollectionItem | null>(null)
const deleteLoading   = ref(false)

const statusOptions = computed(() => [
  { label: t('admin.docs.status_published'), value: 2 },
  { label: t('admin.docs.status_draft'), value: 1 },
])

const isEditing = computed(() => editingId.value !== null)

const filteredCollections = computed(() => {
  if (!searchQuery.value) return collections.value
  const q = searchQuery.value.toLowerCase()
  return collections.value.filter(
    c =>
      c.title.toLowerCase().includes(q) ||
      c.slug.toLowerCase().includes(q) ||
      (c.description ?? '').toLowerCase().includes(q),
  )
})

async function fetchCollections() {
  rawLoading.value = true
  try {
    const res = await docApi.getCollections({ page_size: 100 })
    collections.value = res.data ?? []
  } catch (e: any) {
    toast.add({ title: e?.message ?? t('common.load_failed'), color: 'error' })
  } finally {
    rawLoading.value = false
  }
}

function resetForm() {
  editingId.value      = null
  formSlug.value       = ''
  formTitle.value      = ''
  formDescription.value = ''
  formStatus.value     = 2
  formSortOrder.value  = 0
  formLocale.value     = 'zh'
}

function startEdit(item: DocCollectionItem) {
  editingId.value       = item.id
  formSlug.value        = item.slug
  formTitle.value       = item.title
  formDescription.value = item.description ?? ''
  formStatus.value      = item.status
  formSortOrder.value   = item.sort_order ?? 0
  formLocale.value      = item.locale ?? 'zh'
}

// Auto-generate slug from title
watch(formTitle, (val) => {
  if (!formSlug.value || (!isEditing.value && !formSlug.value)) {
    formSlug.value = val
      .toLowerCase()
      .trim()
      .replace(/[^\w\s-]/g, '')
      .replace(/\s+/g, '-')
      .replace(/-+/g, '-')
  }
})

async function submitForm() {
  if (!formTitle.value.trim() || !formSlug.value.trim()) {
    toast.add({ title: t('validation.required'), color: 'warning' })
    return
  }
  formSubmitting.value = true
  try {
    const data = {
      slug:        formSlug.value.trim(),
      title:       formTitle.value.trim(),
      description: formDescription.value.trim() || undefined,
      status:      formStatus.value,
      sort_order:  formSortOrder.value,
      locale:      formLocale.value || undefined,
    }
    if (isEditing.value) {
      await docApi.updateCollection(editingId.value!, data)
      toast.add({ title: t('common.save_success'), color: 'success' })
    } else {
      await docApi.createCollection(data)
      toast.add({ title: t('common.created'), color: 'success' })
    }
    resetForm()
    fetchCollections()
  } catch (e: any) {
    toast.add({ title: e?.message ?? t('common.save_failed'), color: 'error' })
  } finally {
    formSubmitting.value = false
  }
}

function openDeleteModal(item: DocCollectionItem) {
  deleteTarget.value = item
  showDeleteModal.value = true
}

async function confirmDelete() {
  if (!deleteTarget.value) return
  deleteLoading.value = true
  try {
    await docApi.deleteCollection(deleteTarget.value.id)
    toast.add({ title: t('common.deleted'), color: 'success' })
    showDeleteModal.value = false
    deleteTarget.value = null
    fetchCollections()
  } catch (e: any) {
    toast.add({ title: e?.message ?? t('common.delete_failed'), color: 'error' })
  } finally {
    deleteLoading.value = false
  }
}

onMounted(fetchCollections)
</script>

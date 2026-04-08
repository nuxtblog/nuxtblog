<template>
  <div v-if="pending" class="flex items-center justify-center min-h-screen">
    <UIcon name="i-tabler-loader-2" class="size-8 text-muted animate-spin" />
  </div>
  <div v-else-if="!initialData">
    <AdminPageContainer>
      <AdminPageHeader title="文档不存在" subtitle="未找到该文档" />
    </AdminPageContainer>
  </div>
  <AdminDocEditor
    v-else
    ref="editorRef"
    mode="edit"
    :initial-data="initialData"
    :submitting="submitting"
    @save="handleSave" />
</template>

<script setup lang="ts">
import type { CreateDocRequest, UpdateDocRequest } from '~/types/api/doc'

interface DocEditorInitialData {
  id?: number
  collectionId?: number
  parentId?: number | null
  title?: string
  slug?: string
  content?: string
  excerpt?: string
  status?: number
  commentStatus?: number
  locale?: string
  sortOrder?: number
  publishedAt?: string
  seo?: { meta_title?: string; meta_desc?: string; og_title?: string; og_image?: string; canonical_url?: string; robots?: string }
}

definePageMeta({ layout: 'admin' })

const route = useRoute()
const toast = useToast()
const { t } = useI18n()
const docApi = useDocApi()
const { apiFetch } = useApiFetch()

const submitting = ref(false)

const editorRef = useTemplateRef<{
  isDirty: Ref<boolean>
  getIsDirty: () => boolean
  markSaved: () => void
  seoData: Ref<Record<string, string>>
}>('editorRef')

onBeforeRouteLeave(() => {
  if (editorRef.value?.getIsDirty()) {
    return window.confirm(t('admin.docs.editor.unsaved_warning'))
  }
})

const handleBeforeUnload = (e: BeforeUnloadEvent) => {
  if (editorRef.value?.getIsDirty()) e.preventDefault()
}
onMounted(() => window.addEventListener('beforeunload', handleBeforeUnload))
onUnmounted(() => window.removeEventListener('beforeunload', handleBeforeUnload))

const { data: initialData, pending } = await useAsyncData(
  `doc-edit-${route.params.id}`,
  async (): Promise<DocEditorInitialData | null> => {
    const doc = await docApi.getDoc(Number(route.params.id)).catch(() => null)
    if (!doc) return null
    return {
      id: doc.id,
      collectionId: doc.collection_id,
      parentId: doc.parent_id,
      title: doc.title,
      slug: doc.slug,
      content: doc.content,
      excerpt: doc.excerpt,
      status: doc.status,
      commentStatus: doc.comment_status,
      locale: doc.locale,
      sortOrder: doc.sort_order,
      publishedAt: doc.published_at ?? undefined,
      seo: doc.seo
        ? {
            meta_title:    doc.seo.meta_title,
            meta_desc:     doc.seo.meta_desc,
            og_title:      doc.seo.og_title,
            og_image:      doc.seo.og_image,
            canonical_url: doc.seo.canonical_url,
            robots:        doc.seo.robots,
          }
        : undefined,
    }
  },
)

const handleSave = async (payload: CreateDocRequest | UpdateDocRequest) => {
  if (!initialData.value?.id) return
  submitting.value = true
  try {
    await docApi.updateDoc(initialData.value.id, payload as UpdateDocRequest)

    // 同步保存 SEO 数据
    const seo = editorRef.value?.seoData?.value
    if (seo) {
      await apiFetch(`/docs/${initialData.value.id}/seo`, {
        method: 'PUT',
        body: seo,
      }).catch(() => {})
    }

    editorRef.value?.markSaved()
    toast.add({ title: t('admin.docs.editor.updated'), color: 'success' })
  } catch (error: any) {
    const msg = error?.data?.message || error?.message || t('common.unknown_error')
    toast.add({ title: t('admin.docs.editor.save_failed'), description: msg, color: 'error' })
  } finally {
    submitting.value = false
  }
}
</script>

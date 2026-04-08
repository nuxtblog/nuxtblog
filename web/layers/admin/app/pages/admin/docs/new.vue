<template>
  <AdminDocEditor
    ref="editorRef"
    mode="create"
    :submitting="submitting"
    @save="handleSave" />
</template>

<script setup lang="ts">
import type { CreateDocRequest, UpdateDocRequest } from '~/types/api/doc'

definePageMeta({ layout: 'admin' })

const submitting = ref(false)
const toast = useToast()
const { t } = useI18n()
const docApi = useDocApi()

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

const handleSave = async (payload: CreateDocRequest | UpdateDocRequest) => {
  submitting.value = true
  try {
    const result = await docApi.createDoc(payload as CreateDocRequest)
    editorRef.value?.markSaved()
    toast.add({ title: t('admin.docs.editor.created'), color: 'success' })
    await navigateTo(`/admin/docs/${result.id}/edit`)
  } catch (error: any) {
    const msg = error?.data?.message || error?.message || t('common.unknown_error')
    toast.add({ title: t('admin.docs.editor.save_failed'), description: msg, color: 'error' })
  } finally {
    submitting.value = false
  }
}
</script>

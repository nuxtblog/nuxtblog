<template>
  <AdminPostEditor
    ref="editorRef"
    mode="create"
    :simple="true"
    :submitting="submitting"
    @save="handleSave"
  />
</template>

<script setup lang="ts">
import type { CreatePostRequest, UpdatePostRequest } from '~/types/api/post'

const submitting = ref(false)
const toast = useToast()
const { t } = useI18n()
const postStore = usePostStore()
const editorRef = useTemplateRef<{ reset: () => void }>('editorRef')

const handleSave = async (payload: CreatePostRequest | UpdatePostRequest) => {
  submitting.value = true
  try {
    const req = payload as CreatePostRequest
    req.post_type = 2 // force page type
    await postStore.addPost(req)
    toast.add({ title: t('admin.pages.saved'), color: 'success' })
    editorRef.value?.reset()
    await navigateTo('/admin/pages')
  } catch (error: any) {
    const msg = error?.data?.message || error?.message || t('common.unknown_error')
    toast.add({ title: t('admin.pages.save_failed'), description: msg, color: 'error' })
  } finally {
    submitting.value = false
  }
}
</script>

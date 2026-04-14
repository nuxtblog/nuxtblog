<script setup lang="ts">
import type { TermDetailResponse } from "~/types/api/term";

const props = defineProps<{
  title: string
  collapsed: boolean
  sidebarLoading?: boolean
}>()

const emit = defineEmits<{
  (e: 'toggle'): void
}>()

const selectedTags = defineModel<TermDetailResponse[]>('tags', { required: true })

const { t } = useI18n()
const toast = useToast()
const tagStore = useTagStore()

const { tags } = storeToRefs(tagStore)
const availableTags = computed(() => {
  const ids = selectedTags.value.map(t => t.term_taxonomy_id)
  return tags.value.filter(t => !ids.includes(t.term_taxonomy_id)).slice(0, 10)
})

const showAddTagModal = ref(false)
const creatingTag = ref(false)
const newTag = ref({ name: '', slug: '' })

const addExistingTag = (tag: TermDetailResponse) => {
  if (!selectedTags.value.find(t => t.term_taxonomy_id === tag.term_taxonomy_id)) {
    selectedTags.value.push(tag)
  }
}
const removeTag = (id: number) => {
  selectedTags.value = selectedTags.value.filter(t => t.term_taxonomy_id !== id)
}
const handleAddTag = async () => {
  if (!newTag.value.name.trim()) return
  creatingTag.value = true
  try {
    const created = await tagStore.addNewTag({
      name: newTag.value.name,
      slug: newTag.value.slug || undefined,
    })
    selectedTags.value.push(created)
    newTag.value = { name: '', slug: '' }
    showAddTagModal.value = false
  } catch (err: any) {
    toast.add({ title: t('admin.posts.editor.create_tag_failed'), description: err?.message, color: 'error' })
  } finally {
    creatingTag.value = false
  }
}
</script>

<template>
  <SidebarCard :title="title" :collapsed="collapsed" @toggle="emit('toggle')">
    <template #header-actions>
      <UButton color="primary" variant="link" size="xs" @click.stop="showAddTagModal = true">
        {{ t('admin.posts.editor.new_tag_btn') }}
      </UButton>
    </template>
    <template v-if="sidebarLoading">
      <div class="flex flex-wrap gap-2 mb-3">
        <USkeleton v-for="i in 3" :key="i" class="h-5 w-16 rounded-full" />
      </div>
      <div>
        <USkeleton class="h-3 w-16 mb-2" />
        <div class="flex flex-wrap gap-2">
          <USkeleton v-for="i in 6" :key="i" class="h-6 w-14 rounded-md" />
        </div>
      </div>
    </template>
    <template v-else>
      <div class="flex flex-wrap gap-2 mb-3">
        <span
          v-for="tag in selectedTags"
          :key="tag.term_taxonomy_id"
          class="inline-flex items-center gap-1 px-2 py-0.5 text-xs rounded-full bg-primary/10 text-primary">
          {{ tag.name }}
          <button type="button" class="hover:text-primary/60 transition-colors" @click="removeTag(tag.term_taxonomy_id)">
            <UIcon name="i-tabler-x" class="size-3" />
          </button>
        </span>
      </div>
      <div v-if="availableTags.length > 0">
        <p class="text-xs text-muted mb-2">{{ t('admin.posts.editor.popular_tags') }}</p>
        <div class="flex flex-wrap gap-2">
          <UButton
            v-for="tag in availableTags"
            :key="tag.term_taxonomy_id"
            color="neutral" variant="outline" size="xs"
            @click="addExistingTag(tag)">
            {{ tag.name }}
          </UButton>
        </div>
      </div>
    </template>
  </SidebarCard>

  <!-- New tag modal -->
  <UModal v-model:open="showAddTagModal" :title="t('admin.posts.editor.new_tag_modal')">
    <template #content>
      <div class="p-6">
        <form class="space-y-4" @submit.prevent="handleAddTag">
          <UFormField :label="t('common.name')" required>
            <UInput v-model="newTag.name" required maxlength="100" :placeholder="t('admin.posts.editor.tag_name_placeholder')" class="w-full" />
          </UFormField>
          <UFormField :label="t('admin.posts.editor.slug_label')" :hint="t('admin.posts.editor.slug_hint')">
            <UInput v-model="newTag.slug" maxlength="100" :placeholder="t('admin.posts.editor.tag_slug_placeholder')" class="w-full" />
          </UFormField>
          <div class="flex gap-3 justify-end">
            <UButton color="neutral" variant="soft" type="button" @click="showAddTagModal = false">{{ t('common.cancel') }}</UButton>
            <UButton color="primary" type="submit" :loading="creatingTag">{{ t('common.create') }}</UButton>
          </div>
        </form>
      </div>
    </template>
  </UModal>
</template>

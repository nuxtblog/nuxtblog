<script setup lang="ts">
const props = defineProps<{
  title: string
  collapsed: boolean
  sidebarLoading?: boolean
}>()

const emit = defineEmits<{
  (e: 'toggle'): void
}>()

const selectedCategories = defineModel<number[]>('categories', { required: true })

const { t } = useI18n()
const toast = useToast()
const categoryStore = useCategoryStore()

const showAddCategoryModal = ref(false)
const creatingCategory = ref(false)
const newCategory = ref({
  name: '',
  slug: '',
  taxonomy: 'category',
  description: '',
  parent_id: undefined as number | undefined,
})

const handleCreateCategory = async () => {
  if (!newCategory.value.name.trim()) return
  creatingCategory.value = true
  try {
    const created = await categoryStore.addNewCategory({
      name: newCategory.value.name,
      slug: newCategory.value.slug || undefined,
      description: newCategory.value.description || undefined,
      parent_id: newCategory.value.parent_id,
    })
    if (created) {
      selectedCategories.value.push(created.term_taxonomy_id)
      showAddCategoryModal.value = false
      newCategory.value = { name: '', slug: '', taxonomy: 'category', description: '', parent_id: undefined }
    }
  } catch (error: any) {
    toast.add({ title: t('admin.posts.editor.create_category_failed'), description: error?.message, color: 'error' })
  } finally {
    creatingCategory.value = false
  }
}
</script>

<template>
  <SidebarCard :title="title" :collapsed="collapsed" @toggle="emit('toggle')">
    <template #header-actions>
      <UButton color="primary" variant="link" size="xs" @click.stop="showAddCategoryModal = true">
        {{ t('admin.posts.editor.new_category_btn') }}
      </UButton>
    </template>
    <div v-if="sidebarLoading" class="space-y-2">
      <div v-for="i in 5" :key="i" class="flex items-center gap-2">
        <USkeleton class="size-4 rounded" />
        <USkeleton :class="`h-4 w-${i % 2 === 0 ? '24' : '32'}`" />
      </div>
    </div>
    <ParentCategoryMultiSelector v-else v-model="selectedCategories" />
  </SidebarCard>

  <!-- New category modal -->
  <UModal v-model:open="showAddCategoryModal" :title="t('admin.posts.editor.new_category_modal')">
    <template #content>
      <div class="p-6">
        <form class="space-y-4" @submit.prevent="handleCreateCategory">
          <UFormField :label="t('common.name')" required>
            <UInput v-model="newCategory.name" required maxlength="100" :placeholder="t('admin.posts.editor.category_name_placeholder')" class="w-full" />
          </UFormField>
          <UFormField :label="t('admin.posts.editor.slug_label')">
            <UInput v-model="newCategory.slug" maxlength="100" :placeholder="t('admin.posts.editor.slug_auto')" class="w-full" />
          </UFormField>
          <ParentCategorySelector :label="t('admin.posts.categories.parent_category')" />
          <UFormField :label="t('common.description')">
            <UTextarea v-model="newCategory.description" :rows="3" maxlength="255" :placeholder="t('admin.posts.editor.category_desc_placeholder')" class="w-full" />
          </UFormField>
          <div class="flex gap-3 justify-end">
            <UButton color="neutral" variant="soft" type="button" @click="showAddCategoryModal = false">{{ t('common.cancel') }}</UButton>
            <UButton color="primary" type="submit" :loading="creatingCategory">{{ t('common.create') }}</UButton>
          </div>
        </form>
      </div>
    </template>
  </UModal>
</template>

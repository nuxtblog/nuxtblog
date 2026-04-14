<script setup lang="ts">
const props = defineProps<{
  open: boolean
  mode: 'single' | 'clearTrash'
  trashCount?: number
}>()

const emit = defineEmits<{
  (e: 'update:open', val: boolean): void
  (e: 'confirm'): void
}>()

const { t } = useI18n()
const loading = ref(false)

watch(() => props.open, (val) => {
  if (!val) loading.value = false
})

const handleConfirm = async () => {
  loading.value = true
  emit('confirm')
}

defineExpose({ setLoading: (val: boolean) => { loading.value = val } })
</script>

<template>
  <UModal :open="open" @update:open="emit('update:open', $event)">
    <template #content>
      <div class="p-6">
        <div class="flex items-center gap-3 mb-4">
          <div class="size-10 rounded-full bg-error/10 flex items-center justify-center shrink-0">
            <UIcon :name="mode === 'clearTrash' ? 'i-tabler-trash-x' : 'i-tabler-alert-triangle'" class="size-5 text-error" />
          </div>
          <div>
            <h3 class="font-semibold text-highlighted">
              {{ mode === 'clearTrash' ? $t('admin.posts.clear_trash_title') : $t('admin.posts.delete_title') }}
            </h3>
            <p class="text-sm text-muted mt-0.5">
              {{ mode === 'clearTrash' ? $t('admin.posts.clear_trash_desc', { n: trashCount }) : $t('admin.posts.delete_desc') }}
            </p>
          </div>
        </div>
        <div class="flex justify-end gap-2 mt-6">
          <UButton color="neutral" variant="outline" @click="emit('update:open', false)">{{ $t('common.cancel') }}</UButton>
          <UButton color="error" :loading="loading" @click="handleConfirm">
            {{ mode === 'clearTrash' ? $t('admin.posts.confirm_clear') : $t('admin.posts.hard_delete') }}
          </UButton>
        </div>
      </div>
    </template>
  </UModal>
</template>

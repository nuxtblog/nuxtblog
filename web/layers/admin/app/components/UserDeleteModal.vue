<template>
  <UModal :open="!!user" @update:open="v => { if (!v) emit('update:user', null) }">
    <template #content>
      <div class="p-6">
        <div class="flex items-center gap-3 mb-4">
          <div class="size-10 rounded-full bg-error/10 flex items-center justify-center shrink-0">
            <UIcon name="i-tabler-alert-triangle" class="size-5 text-error" />
          </div>
          <div>
            <h3 class="font-semibold text-highlighted">{{ $t('admin.users.delete_title') }}</h3>
            <p class="text-sm text-muted mt-0.5">{{ $t('admin.users.delete_desc', { username: user?.username }) }}</p>
          </div>
        </div>
        <p class="text-sm text-warning mb-2">{{ $t('admin.users.delete_warning') }}</p>
        <div class="flex justify-end gap-2 mt-6">
          <UButton color="neutral" variant="outline" @click="emit('update:user', null)">{{ $t('common.cancel') }}</UButton>
          <UButton color="error" :loading="deleting" @click="handleDelete">
            {{ deleting ? $t('admin.users.deleting') : $t('common.confirm_delete') }}
          </UButton>
        </div>
      </div>
    </template>
  </UModal>
</template>

<script setup lang="ts">
import type { UserListResponse } from '~/types/api/user'

const props = defineProps<{
  user: UserListResponse | null
  currentUserId?: number
}>()

const emit = defineEmits<{ 'update:user': [null]; 'deleted': [] }>()

const userApi = useUserApi()
const toast = useToast()
const { t } = useI18n()
const deleting = ref(false)

const handleDelete = async () => {
  if (!props.user) return
  if (props.user.id === props.currentUserId) { emit('update:user', null); return }
  deleting.value = true
  try {
    await userApi.deleteUser(props.user.id)
    emit('update:user', null)
    emit('deleted')
  } catch (e: any) {
    toast.add({ title: t('admin.users.delete_failed'), description: e?.message, color: 'error' })
  } finally {
    deleting.value = false
  }
}
</script>

<template>
  <UModal :open="!!user" @update:open="v => { if (!v) emit('update:user', null) }">
    <template #content>
      <div class="p-6">
        <h3 class="text-lg font-semibold text-highlighted mb-4">{{ $t('admin.users.reset_password') }}</h3>
        <form @submit.prevent="handleSubmit" class="space-y-4">
          <p class="text-sm text-muted">{{ $t('admin.users.reset_pw_desc', { username: user?.username }) }}</p>
          <UFormField :label="$t('admin.users.new_password_label')" required>
            <UInput
              v-model="password"
              type="password"
              :placeholder="$t('admin.users.field_password_placeholder')"
              class="w-full" />
          </UFormField>
          <div class="flex gap-3 justify-end">
            <UButton type="button" color="neutral" variant="outline" @click="emit('update:user', null)">
              {{ $t('common.cancel') }}
            </UButton>
            <UButton type="submit" color="primary" :loading="resetting">
              {{ resetting ? $t('admin.users.resetting') : $t('admin.users.confirm_reset') }}
            </UButton>
          </div>
        </form>
      </div>
    </template>
  </UModal>
</template>

<script setup lang="ts">
import type { UserListResponse } from '~/types/api/user'

const props = defineProps<{ user: UserListResponse | null }>()
const emit = defineEmits<{ 'update:user': [null] }>()

const userApi = useUserApi()
const toast = useToast()
const { t } = useI18n()

const password = ref('')
const resetting = ref(false)

watch(() => props.user, (u) => { if (u) password.value = '' })

const handleSubmit = async () => {
  if (!props.user) return
  resetting.value = true
  try {
    await userApi.resetUserPassword(props.user.id, { new_password: password.value })
    toast.add({ title: t('admin.users.pw_reset_success'), color: 'success' })
    emit('update:user', null)
  } catch (e: any) {
    toast.add({ title: t('admin.users.pw_reset_failed'), description: e?.message, color: 'error' })
  } finally {
    resetting.value = false
  }
}
</script>

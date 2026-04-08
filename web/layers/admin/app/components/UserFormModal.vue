<template>
  <UModal :open="open" @update:open="emit('update:open', $event)">
    <template #content>
      <div class="p-6 max-h-[90vh] overflow-y-auto">
        <h3 class="text-lg font-semibold text-highlighted mb-4">
          {{ user ? $t('admin.users.edit_modal_title') : $t('admin.users.add_modal_title') }}
        </h3>
        <form @submit.prevent="handleSubmit" class="space-y-4">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <UFormField :label="$t('admin.users.field_username')" required>
              <UInput
                v-model="form.username"
                :disabled="!!user"
                :placeholder="$t('admin.users.field_username_placeholder')"
                class="w-full" />
            </UFormField>
            <UFormField :label="$t('admin.users.field_display_name')">
              <UInput
                v-model="form.display_name"
                :placeholder="$t('admin.users.field_display_name_placeholder')"
                class="w-full" />
            </UFormField>
            <UFormField :label="$t('admin.users.field_email')" required>
              <UInput v-model="form.email" type="email" placeholder="user@example.com" class="w-full" />
            </UFormField>
            <UFormField v-if="!user" :label="$t('admin.users.field_password')" required>
              <UInput
                v-model="form.password"
                type="password"
                :placeholder="$t('admin.users.field_password_placeholder')"
                class="w-full" />
            </UFormField>
            <UFormField :label="$t('admin.users.field_role')" required>
              <USelect
                v-model="form.role"
                :items="[{ label: $t('admin.users.field_role_placeholder'), value: 0 }, ...roleList.map(r => ({ label: r.name, value: r.id }))]"
                class="w-full" />
            </UFormField>
            <UFormField :label="$t('admin.users.field_status')">
              <USelect
                v-model="form.status"
                :items="[
                  { label: $t('admin.users.status_active'), value: 1 },
                  { label: $t('admin.users.status_pending'), value: 3 },
                  { label: $t('admin.users.status_banned'), value: 2 },
                ]"
                class="w-full" />
            </UFormField>
          </div>
          <UFormField :label="$t('admin.users.field_bio')">
            <UTextarea
              v-model="form.bio"
              :rows="3"
              :placeholder="$t('admin.users.field_bio_placeholder')"
              class="w-full resize-none" />
            <p class="text-xs text-muted mt-1">{{ form.bio?.length || 0 }} / 500</p>
          </UFormField>
          <div class="flex gap-3 justify-end pt-2">
            <UButton type="button" color="neutral" variant="outline" @click="emit('update:open', false)">
              {{ $t('common.cancel') }}
            </UButton>
            <UButton type="submit" color="primary" :loading="submitting">
              {{ submitting ? $t('common.saving') : user ? $t('admin.users.save_changes') : $t('admin.users.create_user') }}
            </UButton>
          </div>
        </form>
      </div>
    </template>
  </UModal>
</template>

<script setup lang="ts">
import type { UserListResponse, UserStatus, UserRole } from '~/types/api/user'

interface RoleItem { id: number; name: string }

const props = defineProps<{
  open: boolean
  user: UserListResponse | null
  roleList: RoleItem[]
}>()

const emit = defineEmits<{
  'update:open': [boolean]
  'saved': []
}>()

const userApi = useUserApi()
const toast = useToast()
const { t } = useI18n()
const submitting = ref(false)

const form = ref({
  username: '',
  password: '',
  email: '',
  display_name: '',
  bio: '',
  status: 1 as UserStatus,
  role: 1 as UserRole,
})

watch(() => props.open, (open) => {
  if (!open) return
  if (props.user) {
    form.value = {
      username: props.user.username,
      password: '',
      email: props.user.email,
      display_name: props.user.display_name || '',
      bio: props.user.bio || '',
      status: props.user.status,
      role: props.user.role,
    }
  } else {
    form.value = { username: '', password: '', email: '', display_name: '', bio: '', status: 1, role: 1 }
  }
})

const handleSubmit = async () => {
  submitting.value = true
  try {
    if (props.user) {
      await userApi.updateUser(props.user.id, {
        display_name: form.value.display_name || undefined,
        bio: form.value.bio || undefined,
        status: form.value.status,
        role: form.value.role,
      })
      toast.add({ title: t('admin.users.updated'), color: 'success' })
    } else {
      await userApi.createUser({
        username: form.value.username,
        password: form.value.password,
        email: form.value.email,
        display_name: form.value.display_name || form.value.username,
        role: form.value.role,
      })
      toast.add({ title: t('admin.users.created'), color: 'success' })
    }
    emit('update:open', false)
    emit('saved')
  } catch (e: any) {
    toast.add({
      title: props.user ? t('admin.users.save_failed') : t('admin.users.create_error'),
      description: e?.message,
      color: 'error',
    })
  } finally {
    submitting.value = false
  }
}
</script>

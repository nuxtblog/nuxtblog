<template>
  <AdminPageContainer>
    <AdminPageHeader :title="$t('admin.users_new.title')" :subtitle="$t('admin.users_new.subtitle')">
      <template #actions>
        <NuxtLink to="/admin/users">
          <UButton color="neutral" variant="outline" leading-icon="i-tabler-corner-up-left">
            {{ $t('admin.users_new.back_to_list') }}
          </UButton>
        </NuxtLink>
      </template>
    </AdminPageHeader>

    <AdminPageContent>
      <form @submit.prevent="handleSubmit" class="space-y-4 md:space-y-6">
        <!-- 基本信息 -->
        <UCard>
          <template #header>
            <h2 class="text-base font-semibold text-highlighted">{{ $t('admin.users_new.basic_info') }}</h2>
          </template>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <!-- 用户名 -->
            <UFormField :label="$t('admin.users_new.username_label')">
              <UInput
                v-model="formData.username"
                required
                :placeholder="$t('admin.users_new.username_placeholder')"
                class="w-full" />
              <template #hint>
                <p class="text-xs text-muted mt-1.5">{{ $t('admin.users_new.username_hint') }}</p>
              </template>
            </UFormField>

            <!-- 显示名称 -->
            <UFormField :label="$t('admin.users_new.display_name_label')">
              <UInput
                v-model="formData.display_name"
                :placeholder="$t('admin.users_new.display_name_placeholder')"
                class="w-full" />
              <template #hint>
                <p class="text-xs text-muted mt-1.5">{{ $t('admin.users_new.display_name_hint') }}</p>
              </template>
            </UFormField>

            <!-- 邮箱 -->
            <UFormField :label="$t('admin.users_new.email_label')">
              <UInput
                v-model="formData.email"
                type="email"
                required
                placeholder="user@example.com"
                class="w-full" />
            </UFormField>

            <!-- 密码 -->
            <UFormField :label="$t('admin.users_new.password_label')">
              <UInput
                v-model="formData.password"
                :type="showPassword ? 'text' : 'password'"
                required
                :placeholder="$t('admin.users_new.password_placeholder')"
                class="w-full"
                :trailing-icon="showPassword ? 'i-tabler-eye-off' : 'i-tabler-eye'"
                @click:trailing="showPassword = !showPassword" />
              <template #hint>
                <p class="text-xs text-muted mt-1.5">{{ $t('admin.users_new.password_hint') }}</p>
              </template>
            </UFormField>
          </div>
        </UCard>

        <!-- 角色和权限 -->
        <UCard>
          <template #header>
            <h2 class="text-base font-semibold text-highlighted">{{ $t('admin.users_new.role_permissions') }}</h2>
          </template>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <!-- 角色选择 -->
            <UFormField :label="$t('admin.users_new.role_label')">
              <USelect
                v-model="formData.role"
                :items="roleItems"
                class="w-full" />
              <template #hint>
                <p class="text-xs text-muted mt-1.5">{{ $t('admin.users_new.role_hint') }}</p>
              </template>
            </UFormField>

            <!-- 语言/地区 -->
            <UFormField :label="$t('admin.users_new.locale_label')">
              <USelect
                v-model="formData.locale"
                :items="[
                  { label: '简体中文', value: 'zh-CN' },
                  { label: 'English', value: 'en' },
                ]"
                class="w-full" />
            </UFormField>
          </div>

          <!-- 角色说明 -->
          <div v-if="selectedRole" class="mt-4 p-4 bg-elevated/50 rounded-md">
            <div class="flex items-start gap-3">
              <UIcon name="i-tabler-info-circle" class="text-primary shrink-0 mt-0.5 size-5" />
              <div>
                <div class="text-sm font-medium text-highlighted">{{ selectedRole.name }}</div>
                <div class="text-xs text-muted mt-1">{{ selectedRole.description }}</div>
              </div>
            </div>
          </div>
        </UCard>

        <!-- 操作按钮 -->
        <div class="flex items-center justify-end gap-3 pt-4">
          <UButton type="button" color="neutral" variant="outline" @click="navigateTo('/admin/users')">
            {{ $t('common.cancel') }}
          </UButton>
          <UButton
            type="button"
            color="neutral"
            variant="outline"
            :disabled="submitting"
            :loading="submitting"
            @click="handleSaveAndNew">
            {{ $t('admin.users_new.save_continue') }}
          </UButton>
          <UButton type="submit" color="primary" :disabled="submitting" :loading="submitting">
            {{ submitting ? $t('admin.users_new.creating') : $t('admin.users_new.create_btn') }}
          </UButton>
        </div>
      </form>
    </AdminPageContent>
  </AdminPageContainer>
</template>

<script setup lang="ts">
import type { UserRole } from '~/types/api/user'

const toast = useToast()
const { t } = useI18n()
const userApi = useUserApi()

// ── State ──────────────────────────────────────────────────────────────────────

const submitting  = ref(false)
const showPassword = ref(false)

const formData = ref({
  username:     '',
  password:     '',
  email:        '',
  display_name: '',
  role:         1 as UserRole,
  locale:       'zh-CN',
})

const initialForm = { username: '', password: '', email: '', display_name: '', role: 1 as UserRole, locale: 'zh-CN' }

// ── Role config ────────────────────────────────────────────────────────────────

const getRoles = () => [
  { id: 1 as UserRole, name: t('admin.users_new.role_subscriber'),  description: t('admin.users_new.role_subscriber_desc') },
  { id: 2 as UserRole, name: t('admin.users_new.role_editor'),      description: t('admin.users_new.role_editor_desc') },
  { id: 3 as UserRole, name: t('admin.users_new.role_admin'),       description: t('admin.users_new.role_admin_desc') },
  { id: 4 as UserRole, name: t('admin.users_new.role_super_admin'), description: t('admin.users_new.role_super_admin_desc') },
]

const roleItems = computed(() => getRoles().map((r) => ({ label: r.name, value: r.id })))

const selectedRole = computed(() => getRoles().find((r) => r.id === formData.value.role))

// ── Helpers ────────────────────────────────────────────────────────────────────

const buildPayload = () => ({
  username:     formData.value.username,
  password:     formData.value.password,
  email:        formData.value.email,
  display_name: formData.value.display_name || formData.value.username,
  role:         formData.value.role,
  locale:       formData.value.locale || undefined,
})

// ── Submit ─────────────────────────────────────────────────────────────────────

const handleSubmit = async () => {
  submitting.value = true
  try {
    await userApi.createUser(buildPayload())
    toast.add({ title: t('admin.users.created'), color: 'success' })
    navigateTo('/admin/users')
  } catch (error: any) {
    toast.add({ title: t('admin.users.create_error'), description: error?.message, color: 'error' })
  } finally {
    submitting.value = false
  }
}

const handleSaveAndNew = async () => {
  submitting.value = true
  try {
    await userApi.createUser(buildPayload())
    toast.add({ title: t('admin.users_new.created_continue'), color: 'success' })
    // Reset form, keep role + locale
    formData.value = { ...initialForm, role: formData.value.role, locale: formData.value.locale }
  } catch (error: any) {
    toast.add({ title: t('admin.users.create_error'), description: error?.message, color: 'error' })
  } finally {
    submitting.value = false
  }
}
</script>

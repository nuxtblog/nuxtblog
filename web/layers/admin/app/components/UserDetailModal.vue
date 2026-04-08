<template>
  <UModal :open="!!user" @update:open="v => { if (!v) emit('update:user', null) }">
    <template #content>
      <div class="p-6 max-h-[90vh] overflow-y-auto">
        <div class="flex items-start justify-between mb-6">
          <h3 class="text-lg font-semibold text-highlighted">{{ $t('admin.users.detail_title') }}</h3>
          <UButton icon="i-tabler-x" color="neutral" variant="ghost" size="sm" square @click="emit('update:user', null)" />
        </div>

        <div v-if="user" class="space-y-6">
          <!-- 头像 + 基本信息 -->
          <div class="flex items-start gap-4">
            <BaseAvatar :src="user.avatar" :alt="user.display_name || user.username" size="xl" />
            <div class="flex-1">
              <h4 class="text-xl font-semibold text-highlighted">{{ user.display_name || user.username }}</h4>
              <p class="text-sm text-muted">@{{ user.username }}</p>
              <div class="flex items-center gap-2 mt-2">
                <UBadge
                  :label="roleMap[user.role]?.name || String(user.role)"
                  :color="(roleMap[user.role]?.color as any) || 'neutral'"
                  variant="soft" />
                <UBadge
                  :label="statusConfig[user.status]?.label || String(user.status)"
                  :color="(statusConfig[user.status]?.color as any) || 'neutral'"
                  variant="soft" />
              </div>
            </div>
          </div>

          <!-- 详细字段 -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div class="space-y-1">
              <div class="text-sm font-medium text-muted">{{ $t('admin.users.field_email_label') }}</div>
              <div class="text-sm text-highlighted">{{ user.email }}</div>
            </div>
            <div class="space-y-1">
              <div class="text-sm font-medium text-muted">{{ $t('admin.users.field_user_id') }}</div>
              <div class="text-sm text-highlighted">{{ user.id }}</div>
            </div>
            <div class="space-y-1">
              <div class="text-sm font-medium text-muted">{{ $t('admin.users.field_register_time') }}</div>
              <div class="text-sm text-highlighted">{{ formatFullDate(user.created_at) }}</div>
            </div>
            <div class="space-y-1">
              <div class="text-sm font-medium text-muted">{{ $t('admin.users.field_last_updated') }}</div>
              <div class="text-sm text-highlighted">{{ formatFullDate(user.updated_at || user.created_at) }}</div>
            </div>
          </div>

          <!-- 个人简介 -->
          <div v-if="user.bio" class="space-y-2">
            <div class="text-sm font-medium text-muted">{{ $t('admin.users.field_bio_label') }}</div>
            <div class="text-sm text-highlighted whitespace-pre-wrap">{{ user.bio }}</div>
          </div>

          <!-- 操作按钮 -->
          <div class="flex gap-3 pt-4 border-t border-default">
            <UButton color="primary" class="flex-1" @click="emit('edit', user); emit('update:user', null)">
              {{ $t('admin.users.edit_user') }}
            </UButton>
            <UButton color="neutral" variant="outline" class="flex-1" @click="emit('resetPassword', user); emit('update:user', null)">
              {{ $t('admin.users.reset_password') }}
            </UButton>
          </div>
        </div>
      </div>
    </template>
  </UModal>
</template>

<script setup lang="ts">
import type { UserListResponse } from '~/types/api/user'

defineProps<{
  user: UserListResponse | null
  roleMap: Record<number, { name: string; color: string }>
  statusConfig: Record<number, { label: string; color: string }>
}>()

const emit = defineEmits<{
  'update:user': [null]
  'edit': [UserListResponse]
  'resetPassword': [UserListResponse]
}>()

const parseDate = (s: string) => new Date(s.includes('T') ? s : s.replace(' ', 'T'))
const formatFullDate = (s: string) =>
  parseDate(s).toLocaleString(undefined, {
    year: 'numeric', month: 'long', day: 'numeric', hour: '2-digit', minute: '2-digit',
  })
</script>

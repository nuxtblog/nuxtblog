<template>
  <div class="space-y-5">
    <!-- 创建 Token 弹窗 -->
    <UModal v-model:open="createModalOpen" :ui="{ content: 'max-w-md' }">
      <template #content>
        <div class="p-5 space-y-4">
          <div>
            <h3 class="font-semibold text-highlighted">
              {{ createdToken ? $t('site.settings.token_created_title') : $t('site.settings.token_create_title') }}
            </h3>
            <p class="text-xs text-muted mt-1">
              {{ createdToken ? $t('site.settings.token_created_hint') : $t('site.settings.token_create_hint') }}
            </p>
          </div>

          <!-- 创建成功：显示 token -->
          <template v-if="createdToken">
            <div class="rounded-md bg-muted p-3 space-y-1">
              <p class="text-xs text-muted">{{ $t('site.settings.token_display_once') }}</p>
              <div class="flex items-center gap-2">
                <code class="flex-1 text-sm font-mono text-highlighted break-all select-all">{{ createdToken.token }}</code>
                <UButton
                  :icon="copiedToken ? 'i-tabler-check' : 'i-tabler-copy'"
                  color="neutral" variant="ghost" size="xs" square
                  @click="copyToken" />
              </div>
            </div>
            <UAlert color="warning" icon="i-tabler-alert-triangle" :title="$t('site.settings.token_warning')" />
            <div class="flex justify-end">
              <UButton color="primary" @click="createModalOpen = false; createdToken = null">{{ $t('site.settings.token_done') }}</UButton>
            </div>
          </template>

          <!-- 创建表单 -->
          <template v-else>
            <UFormField :label="$t('site.settings.token_name_label')" :hint="$t('site.settings.token_name_hint')">
              <UInput v-model="newTokenName" :placeholder="$t('site.settings.token_name_placeholder')" class="w-full" />
            </UFormField>
            <UFormField :label="$t('site.settings.token_expiry_label')">
              <div class="flex flex-wrap gap-2">
                <UButton
                  v-for="opt in expiryOptions"
                  :key="String(opt.value)"
                  size="sm" color="neutral"
                  :variant="newTokenExpiry === opt.value ? 'solid' : 'outline'"
                  @click="newTokenExpiry = opt.value">
                  {{ opt.label }}
                </UButton>
              </div>
            </UFormField>
            <div class="flex gap-2 justify-end pt-2">
              <UButton color="neutral" variant="ghost" @click="createModalOpen = false">{{ $t('site.settings.token_cancel') }}</UButton>
              <UButton color="primary" :loading="tokenCreating" :disabled="!newTokenName.trim()" @click="createToken">
                {{ $t('site.settings.token_generate') }}
              </UButton>
            </div>
          </template>
        </div>
      </template>
    </UModal>

    <!-- Token 列表 Card -->
    <UCard>
      <template #header>
        <div class="flex items-center justify-between">
          <div>
            <h2 class="font-semibold text-highlighted flex items-center gap-2">
              <UIcon name="i-tabler-key" class="size-5 text-primary" />
              {{ $t('site.settings.tokens_title') }}
            </h2>
            <p class="text-xs text-muted mt-0.5">{{ $t('site.settings.tokens_desc') }}</p>
          </div>
          <UButton color="primary" size="sm" icon="i-tabler-plus" @click="openCreateModal">
            {{ $t('site.settings.new_token_btn') }}
          </UButton>
        </div>
      </template>

      <div v-if="tokensLoading" class="space-y-3">
        <div v-for="i in 2" :key="i" class="flex items-center gap-3 py-2">
          <USkeleton class="size-8 rounded-md shrink-0" />
          <div class="flex-1 space-y-1.5">
            <USkeleton class="h-4 w-36" />
            <USkeleton class="h-3 w-48" />
          </div>
          <USkeleton class="h-7 w-14 rounded" />
        </div>
      </div>

      <div v-else-if="tokens.length === 0" class="py-10 text-center">
        <div class="size-12 rounded-md bg-muted flex items-center justify-center mx-auto mb-3">
          <UIcon name="i-tabler-key-off" class="size-6 text-muted" />
        </div>
        <p class="text-sm text-highlighted font-medium mb-1">{{ $t('site.settings.token_no_tokens') }}</p>
        <p class="text-xs text-muted">{{ $t('site.settings.token_no_tokens_desc') }}</p>
      </div>

      <div v-else class="divide-y divide-default -my-1">
        <div v-for="token in tokens" :key="token.id" class="flex items-start gap-3 py-4">
          <div
            class="size-9 rounded-md flex items-center justify-center shrink-0 mt-0.5"
            :class="isExpired(token) ? 'bg-error/10' : 'bg-primary/10'">
            <UIcon name="i-tabler-key" :class="isExpired(token) ? 'text-error' : 'text-primary'" class="size-4" />
          </div>
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2 flex-wrap">
              <span class="text-sm font-medium text-highlighted">{{ token.name }}</span>
              <UBadge v-if="isExpired(token)" :label="$t('site.settings.token_expired')" color="error" variant="soft" size="xs" />
            </div>
            <p class="text-xs text-muted mt-0.5 font-mono">{{ token.prefix }}…</p>
            <div class="flex flex-wrap gap-x-3 text-xs text-muted mt-1">
              <span>{{ $t('site.settings.token_created_at', { date: formatDate(token.created_at) }) }}</span>
              <span v-if="token.expires_at">{{ $t('site.settings.token_expires_at', { date: formatDate(token.expires_at) }) }}</span>
              <span v-else>{{ $t('site.settings.token_never_expires') }}</span>
              <span v-if="token.last_used_at">{{ $t('site.settings.token_last_used', { date: formatDate(token.last_used_at) }) }}</span>
              <span v-else class="italic">{{ $t('site.settings.token_never_used') }}</span>
            </div>
          </div>
          <UButton color="error" variant="soft" size="xs" icon="i-tabler-trash" @click="revokeToken(token.id)">
            {{ $t('site.settings.token_revoke') }}
          </UButton>
        </div>
      </div>
    </UCard>

    <!-- API 使用说明 -->
    <UCard>
      <template #header>
        <h2 class="font-semibold text-highlighted flex items-center gap-2">
          <UIcon name="i-tabler-terminal" class="size-5 text-muted" />
          {{ $t('site.settings.usage_title') }}
        </h2>
      </template>
      <div class="space-y-3">
        <p class="text-sm text-muted">{{ $t('site.settings.usage_desc') }}</p>
        <div class="rounded-md bg-muted p-3">
          <code class="text-sm font-mono text-highlighted whitespace-pre-wrap break-all">Authorization: Bearer yblog_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx</code>
        </div>
        <p class="text-xs text-muted">{{ $t('site.settings.usage_note') }}</p>
      </div>
    </UCard>
  </div>
</template>

<script setup lang="ts">
import type { TokenItem, TokenCreateResult } from '~/composables/useTokenApi'

const tokenApi = useTokenApi()
const toast = useToast()
const { t } = useI18n()

const tokens = ref<TokenItem[]>([])
const tokensLoading = ref(false)
const tokensLoaded = ref(false)

const createModalOpen = ref(false)
const newTokenName = ref('')
const newTokenExpiry = ref<number | null>(null)
const tokenCreating = ref(false)
const createdToken = ref<TokenCreateResult | null>(null)
const copiedToken = ref(false)

const expiryOptions = computed(() => [
  { label: t('site.settings.token_expiry_never'), value: null },
  { label: t('site.settings.token_expiry_7d'), value: 7 },
  { label: t('site.settings.token_expiry_30d'), value: 30 },
  { label: t('site.settings.token_expiry_90d'), value: 90 },
  { label: t('site.settings.token_expiry_1y'), value: 365 },
])

async function loadTokens() {
  if (tokensLoaded.value) return
  tokensLoading.value = true
  try {
    tokens.value = await tokenApi.list()
    tokensLoaded.value = true
  } catch {
    toast.add({ title: t('site.settings.token_load_failed'), color: 'error' })
  } finally {
    tokensLoading.value = false
  }
}

onMounted(loadTokens)

function openCreateModal() {
  newTokenName.value = ''
  newTokenExpiry.value = null
  createdToken.value = null
  createModalOpen.value = true
}

async function createToken() {
  if (!newTokenName.value.trim()) {
    toast.add({ title: t('site.settings.token_no_name'), color: 'error' })
    return
  }
  tokenCreating.value = true
  try {
    const result = await tokenApi.create(newTokenName.value.trim(), newTokenExpiry.value ?? undefined)
    createdToken.value = result
    tokens.value.unshift(result)
  } catch (e: any) {
    toast.add({ title: t('site.settings.token_create_failed'), description: e.message, color: 'error' })
  } finally {
    tokenCreating.value = false
  }
}

async function copyToken() {
  if (!createdToken.value?.token) return
  await navigator.clipboard.writeText(createdToken.value.token)
  copiedToken.value = true
  setTimeout(() => (copiedToken.value = false), 2000)
}

async function revokeToken(id: number) {
  if (!confirm(t('site.settings.token_revoke_confirm'))) return
  try {
    await tokenApi.revoke(id)
    tokens.value = tokens.value.filter(tok => tok.id !== id)
    toast.add({ title: t('site.settings.token_revoked'), color: 'success' })
  } catch (e: any) {
    toast.add({ title: t('site.settings.token_revoke_failed'), description: e.message, color: 'error' })
  }
}

function formatDate(s?: string) {
  if (!s) return '—'
  return new Date(s).toLocaleDateString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit' })
}

function isExpired(tok: TokenItem) {
  return !!tok.expires_at && new Date(tok.expires_at) < new Date()
}
</script>

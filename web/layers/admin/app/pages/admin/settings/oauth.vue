<script setup lang="ts">
import type { OAuthProviderConfig, GenericOAuthProviderConfig } from '~/config/options'

const { apiFetch } = useApiFetch()
const optionsStore = useOptionsStore()
const toast        = useToast()

// ── Registration setting ──────────────────────────────────────────────────
const allowRegistration = ref(false)
const regSaving = ref(false)

const saveRegistration = async () => {
  regSaving.value = true
  try {
    await apiFetch('/options/allow_registration', {
      method: 'PUT',
      body: { value: JSON.stringify(allowRegistration.value), autoload: 1 },
    })
    await optionsStore.reload()
    toast.add({ title: '保存成功', color: 'success' })
  } catch (e: any) {
    toast.add({ title: '保存失败', description: e?.message, color: 'error' })
  } finally {
    regSaving.value = false
  }
}

// ── Builtin providers (backed by Go code) ─────────────────────────────────
const BUILTIN = [
  { slug: 'github', label: 'GitHub', icon: 'i-tabler-brand-github',
    docsUrl: 'https://github.com/settings/developers', docsLabel: 'GitHub Developer Settings' },
  { slug: 'google', label: 'Google', icon: 'i-tabler-brand-google',
    docsUrl: 'https://console.cloud.google.com/apis/credentials', docsLabel: 'Google Cloud Console' },
  { slug: 'qq',     label: 'QQ',     icon: 'i-tabler-brand-qq',
    docsUrl: 'https://connect.qq.com/manage.html', docsLabel: 'QQ 互联管理中心' },
]

type BuiltinForms = Record<string, OAuthProviderConfig>
const builtinForms  = ref<BuiltinForms>({})
const builtinSaving = ref<Record<string, boolean>>({})
const isLoading     = ref(true)

const defaultBuiltin = (slug: string): OAuthProviderConfig => ({
  enabled: false, clientId: '', clientSecret: '',
  callbackUrl: `http://localhost:9000/api/v1/auth/oauth/${slug}/callback`,
})

// ── Generic (frontend-added) providers ───────────────────────────────────
const genericList    = ref<GenericOAuthProviderConfig[]>([])
const genericSaving  = ref<Record<string, boolean>>({})
const showAddForm    = ref(false)
const editingGeneric = ref<GenericOAuthProviderConfig | null>(null)

const defaultGeneric = (): GenericOAuthProviderConfig => ({
  enabled: false, clientId: '', clientSecret: '',
  callbackUrl: 'http://localhost:9000/api/v1/auth/oauth/__slug__/callback',
  slug: '', label: '', icon: 'i-tabler-login',
  authUrl: '', tokenUrl: '', userInfoUrl: '',
  scopes: [],
  fields: { id: 'id', email: 'email', name: 'name', avatar: 'avatar' },
})

// Common icon options for the icon picker (Tabler brand icons)
const ICON_OPTIONS = [
  'i-tabler-brand-discord', 'i-tabler-brand-twitter', 'i-tabler-brand-facebook',
  'i-tabler-brand-instagram', 'i-tabler-brand-linkedin', 'i-tabler-brand-microsoft',
  'i-tabler-brand-apple', 'i-tabler-brand-steam', 'i-tabler-brand-twitch',
  'i-tabler-brand-gitlab', 'i-tabler-brand-bitbucket', 'i-tabler-login',
]

// ── Scope input helper ────────────────────────────────────────────────────
const scopeInput = ref('')
const addScope = (cfg: GenericOAuthProviderConfig) => {
  const s = scopeInput.value.trim()
  if (s && !cfg.scopes.includes(s)) cfg.scopes.push(s)
  scopeInput.value = ''
}
const removeScope = (cfg: GenericOAuthProviderConfig, i: number) => cfg.scopes.splice(i, 1)

// ── Slug auto-sync callback URL ───────────────────────────────────────────
const syncCallbackUrl = (cfg: GenericOAuthProviderConfig) => {
  if (cfg.callbackUrl.includes('__slug__') || cfg.callbackUrl === '') {
    cfg.callbackUrl = `http://localhost:9000/api/v1/auth/oauth/${cfg.slug}/callback`
  }
}

// ── Load ─────────────────────────────────────────────────────────────────
const loadAll = async () => {
  isLoading.value = true
  try {
    // Load registration setting
    try {
      const regRes = await apiFetch<{ value: string }>('/options/allow_registration')
      allowRegistration.value = JSON.parse(regRes.value) ?? false
    } catch { allowRegistration.value = false }

    // Load builtin configs
    await Promise.all(BUILTIN.map(async (p) => {
      try {
        const res = await apiFetch<{ value: string }>(`/options/oauth_${p.slug}`)
        builtinForms.value[p.slug] = JSON.parse(res.value) as OAuthProviderConfig
      } catch {
        builtinForms.value[p.slug] = defaultBuiltin(p.slug)
      }
    }))

    // Load generic provider list
    const slugsRes = await apiFetch<{ value: string }>('/options/oauth_providers').catch(() => null)
    const slugs: string[] = slugsRes ? JSON.parse(slugsRes.value) : []
    const generics = await Promise.all(slugs.map(async (slug) => {
      try {
        const res = await apiFetch<{ value: string }>(`/options/oauth_${slug}`)
        return JSON.parse(res.value) as GenericOAuthProviderConfig
      } catch {
        return null
      }
    }))
    genericList.value = generics.filter(Boolean) as GenericOAuthProviderConfig[]
  } finally {
    isLoading.value = false
  }
}

// ── Save builtin ─────────────────────────────────────────────────────────
const saveBuiltin = async (slug: string) => {
  builtinSaving.value[slug] = true
  try {
    await apiFetch(`/options/oauth_${slug}`, {
      method: 'PUT',
      body: { value: JSON.stringify(builtinForms.value[slug]), autoload: 0 },
    })
    await optionsStore.reload()
    toast.add({ title: '保存成功', color: 'success' })
  } catch (e: any) {
    toast.add({ title: '保存失败', description: e?.message, color: 'error' })
  } finally {
    builtinSaving.value[slug] = false
  }
}

// ── Save generic ──────────────────────────────────────────────────────────
const saveGeneric = async (cfg: GenericOAuthProviderConfig) => {
  if (!cfg.slug) return toast.add({ title: 'Slug 不能为空', color: 'error' })
  if (BUILTIN.some(b => b.slug === cfg.slug)) {
    return toast.add({ title: `"${cfg.slug}" 是内置提供商，请换一个 slug`, color: 'error' })
  }

  genericSaving.value[cfg.slug] = true
  try {
    // Save provider config
    await apiFetch(`/options/oauth_${cfg.slug}`, {
      method: 'PUT',
      body: { value: JSON.stringify(cfg), autoload: 0 },
    })

    // Update the providers list
    const currentSlugs = genericList.value.map(g => g.slug)
    if (!currentSlugs.includes(cfg.slug)) currentSlugs.push(cfg.slug)
    await apiFetch('/options/oauth_providers', {
      method: 'PUT',
      body: { value: JSON.stringify(currentSlugs), autoload: 0 },
    })

    await loadAll()
    showAddForm.value = false
    editingGeneric.value = null
    toast.add({ title: `${cfg.label || cfg.slug} 已保存`, color: 'success' })
  } catch (e: any) {
    toast.add({ title: '保存失败', description: e?.message, color: 'error' })
  } finally {
    genericSaving.value[cfg.slug] = false
  }
}

// ── Delete generic ────────────────────────────────────────────────────────
const deleteGeneric = async (slug: string) => {
  try {
    await apiFetch(`/options/oauth_${slug}`, { method: 'DELETE' })
    const newSlugs = genericList.value.map(g => g.slug).filter(s => s !== slug)
    await apiFetch('/options/oauth_providers', {
      method: 'PUT',
      body: { value: JSON.stringify(newSlugs), autoload: 0 },
    })
    await loadAll()
    toast.add({ title: '已删除', color: 'success' })
  } catch (e: any) {
    toast.add({ title: '删除失败', description: e?.message, color: 'error' })
  }
}

onMounted(loadAll)
</script>

<template>
  <AdminPageContainer>
    <AdminPageHeader :title="$t('admin.settings.oauth.title')" :subtitle="$t('admin.settings.oauth.subtitle')" />
    <AdminPageContent>
    <div v-if="isLoading" class="flex justify-center py-12">
      <UIcon name="i-tabler-loader-2" class="animate-spin text-2xl text-muted" />
    </div>

    <div v-else class="space-y-8">
      <!-- ── Registration ───────────────────────────────────────────── -->
      <div class="space-y-4">
        <h2 class="text-sm font-semibold text-muted uppercase tracking-wide">成员资格</h2>
        <UCard>
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm font-medium text-highlighted">开放注册</p>
              <p class="text-xs text-muted mt-0.5">允许任何人注册账户</p>
            </div>
            <USwitch v-model="allowRegistration" />
          </div>
          <template #footer>
            <div class="flex justify-end">
              <UButton color="primary" size="sm" :loading="regSaving" @click="saveRegistration">保存</UButton>
            </div>
          </template>
        </UCard>
      </div>

      <!-- ── Builtin providers ──────────────────────────────────────── -->
      <div class="space-y-4">
        <h2 class="text-sm font-semibold text-muted uppercase tracking-wide">内置提供商</h2>

        <UAlert
          icon="i-tabler-info-circle"
          color="neutral"
          variant="subtle"
          description="在对应平台创建 OAuth 应用后填入凭证。Callback URL 需与平台后台完全一致。"
        />

        <UCard v-for="p in BUILTIN" :key="p.slug">
          <template #header>
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2">
                <UIcon :name="p.icon" class="text-xl" />
                <span class="font-semibold">{{ p.label }}</span>
                <UBadge
                  :label="builtinForms[p.slug]?.enabled ? '已启用' : '未启用'"
                  :color="builtinForms[p.slug]?.enabled ? 'success' : 'neutral'"
                  variant="soft" size="sm"
                />
              </div>
              <UButton :href="p.docsUrl" target="_blank" rel="noopener"
                color="neutral" variant="ghost" size="xs" trailing-icon="i-tabler-external-link">
                {{ p.docsLabel }}
              </UButton>
            </div>
          </template>

          <div v-if="builtinForms[p.slug]" class="space-y-4">
            <div class="flex items-center justify-between rounded-md border border-default px-4 py-3">
              <div>
                <p class="text-sm font-medium">启用 {{ p.label }} 登录</p>
                <p class="text-xs text-muted">启用后登录页将显示该按钮</p>
              </div>
              <USwitch v-model="builtinForms[p.slug].enabled" />
            </div>
            <UFormField label="Client ID" required>
              <UInput v-model="builtinForms[p.slug].clientId" placeholder="粘贴 Client ID" class="w-full font-mono" />
            </UFormField>
            <UFormField label="Client Secret" required>
              <UInput v-model="builtinForms[p.slug].clientSecret" type="password" placeholder="粘贴 Client Secret" class="w-full font-mono" />
            </UFormField>
            <UFormField label="Callback URL">
              <div class="flex gap-2">
                <UInput v-model="builtinForms[p.slug].callbackUrl" readonly class="w-full font-mono text-xs" />
                <UButton icon="i-tabler-copy" color="neutral" variant="ghost" square
                  @click="navigator.clipboard.writeText(builtinForms[p.slug].callbackUrl); toast.add({ title: '已复制', color: 'success' })" />
              </div>
            </UFormField>
          </div>

          <template #footer>
            <div class="flex justify-end gap-2">
              <UButton color="neutral" variant="ghost" @click="builtinForms[p.slug] = defaultBuiltin(p.slug)">重置</UButton>
              <UButton color="primary" :loading="builtinSaving[p.slug]"
                :disabled="!builtinForms[p.slug]?.clientId || !builtinForms[p.slug]?.clientSecret"
                @click="saveBuiltin(p.slug)">保存</UButton>
            </div>
          </template>
        </UCard>
      </div>

      <!-- ── Custom providers ───────────────────────────────────────── -->
      <div class="space-y-4">
        <div class="flex items-center justify-between">
          <h2 class="text-sm font-semibold text-muted uppercase tracking-wide">自定义提供商</h2>
          <UButton icon="i-tabler-plus" size="sm" @click="editingGeneric = defaultGeneric(); showAddForm = true">
            添加提供商
          </UButton>
        </div>

        <UAlert
          icon="i-tabler-plug"
          color="neutral"
          variant="subtle"
          description="支持任意标准 OAuth 2.0（授权码模式）提供商，如 Discord、微博、抖音等。无需修改后端代码。"
        />

        <!-- Existing generic providers -->
        <UCard v-for="g in genericList" :key="g.slug">
          <template #header>
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2">
                <UIcon :name="g.icon || 'i-tabler-login'" class="text-xl" />
                <span class="font-semibold">{{ g.label || g.slug }}</span>
                <UBadge :label="g.enabled ? '已启用' : '未启用'"
                  :color="g.enabled ? 'success' : 'neutral'" variant="soft" size="sm" />
              </div>
              <div class="flex gap-1">
                <UButton icon="i-tabler-edit" color="neutral" variant="ghost" size="xs"
                  @click="editingGeneric = { ...g }; showAddForm = true" />
                <UButton icon="i-tabler-trash" color="error" variant="ghost" size="xs"
                  @click="deleteGeneric(g.slug)" />
              </div>
            </div>
          </template>
          <div class="grid grid-cols-2 gap-x-6 gap-y-1 text-xs text-muted">
            <div><span class="font-medium">Auth URL：</span>{{ g.authUrl }}</div>
            <div><span class="font-medium">Token URL：</span>{{ g.tokenUrl }}</div>
            <div><span class="font-medium">Scopes：</span>{{ g.scopes.join(' ') }}</div>
            <div><span class="font-medium">Client ID：</span>{{ g.clientId ? '已配置' : '未配置' }}</div>
          </div>
        </UCard>

        <div v-if="genericList.length === 0 && !showAddForm"
          class="rounded-md border border-dashed border-default p-8 text-center text-sm text-muted">
          暂无自定义提供商，点击「添加提供商」开始配置
        </div>
      </div>

      <!-- ── Add / Edit generic provider form ──────────────────────── -->
      <UCard v-if="showAddForm && editingGeneric" class="border-primary/50">
        <template #header>
          <div class="flex items-center justify-between">
            <h3 class="font-semibold">{{ editingGeneric.slug ? `编辑 ${editingGeneric.label || editingGeneric.slug}` : '添加自定义提供商' }}</h3>
            <UButton icon="i-tabler-x" color="neutral" variant="ghost" square
              @click="showAddForm = false; editingGeneric = null" />
          </div>
        </template>

        <div class="space-y-5">
          <!-- Identity -->
          <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
            <UFormField label="Slug（唯一标识）" required>
              <UInput v-model="editingGeneric.slug" placeholder="如 discord、weibo"
                class="font-mono" @input="syncCallbackUrl(editingGeneric!)" />
            </UFormField>
            <UFormField label="显示名称" required>
              <UInput v-model="editingGeneric.label" placeholder="如 Discord" />
            </UFormField>
            <UFormField label="图标">
              <USelect
                v-model="editingGeneric.icon"
                :items="ICON_OPTIONS.map(i => ({ label: i.replace('i-tabler-brand-','').replace('i-tabler-',''), value: i }))"
              />
            </UFormField>
          </div>

          <!-- Enable -->
          <div class="flex items-center justify-between rounded-md border border-default px-4 py-3">
            <p class="text-sm font-medium">启用此提供商</p>
            <USwitch v-model="editingGeneric.enabled" />
          </div>

          <!-- Credentials -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <UFormField label="Client ID" required>
              <UInput v-model="editingGeneric.clientId" class="font-mono" />
            </UFormField>
            <UFormField label="Client Secret" required>
              <UInput v-model="editingGeneric.clientSecret" type="password" class="font-mono" />
            </UFormField>
          </div>

          <!-- Endpoints -->
          <div class="space-y-3">
            <p class="text-sm font-medium text-highlighted">OAuth2 端点</p>
            <UFormField label="Authorization URL" required>
              <UInput v-model="editingGeneric.authUrl" placeholder="https://example.com/oauth/authorize" class="w-full font-mono text-xs" />
            </UFormField>
            <UFormField label="Token URL" required>
              <UInput v-model="editingGeneric.tokenUrl" placeholder="https://example.com/oauth/token" class="w-full font-mono text-xs" />
            </UFormField>
            <UFormField label="User Info URL" required>
              <UInput v-model="editingGeneric.userInfoUrl" placeholder="https://example.com/api/user" class="w-full font-mono text-xs" />
            </UFormField>
            <UFormField label="Callback URL">
              <div class="flex gap-2">
                <UInput v-model="editingGeneric.callbackUrl" class="w-full font-mono text-xs" />
                <UButton icon="i-tabler-copy" color="neutral" variant="ghost" square
                  @click="navigator.clipboard.writeText(editingGeneric!.callbackUrl); toast.add({ title: '已复制', color: 'success' })" />
              </div>
              <template #hint><span class="text-xs text-muted">在平台后台填写此地址</span></template>
            </UFormField>
          </div>

          <!-- Scopes -->
          <UFormField label="Scopes">
            <div class="space-y-2">
              <div class="flex gap-2">
                <UInput v-model="scopeInput" placeholder="输入 scope 后按 Enter 添加"
                  class="flex-1" @keydown.enter.prevent="addScope(editingGeneric!)" />
                <UButton icon="i-tabler-plus" color="neutral" variant="outline" square
                  @click="addScope(editingGeneric!)" />
              </div>
              <div class="flex flex-wrap gap-1">
                <UBadge v-for="(s, i) in editingGeneric.scopes" :key="i"
                  :label="s" color="neutral" variant="soft"
                  class="cursor-pointer" @click="removeScope(editingGeneric!, i)">
                  <template #trailing><UIcon name="i-tabler-x" class="size-3" /></template>
                </UBadge>
              </div>
            </div>
          </UFormField>

          <!-- Field mappings -->
          <div class="space-y-3">
            <p class="text-sm font-medium text-highlighted">用户信息字段映射</p>
            <p class="text-xs text-muted">填写 User Info 接口返回的 JSON 字段名</p>
            <div class="grid grid-cols-2 md:grid-cols-4 gap-3">
              <UFormField label="ID 字段">
                <UInput v-model="editingGeneric.fields.id" placeholder="id" class="font-mono" />
              </UFormField>
              <UFormField label="邮箱字段">
                <UInput v-model="editingGeneric.fields.email" placeholder="email" class="font-mono" />
              </UFormField>
              <UFormField label="名称字段">
                <UInput v-model="editingGeneric.fields.name" placeholder="name" class="font-mono" />
              </UFormField>
              <UFormField label="头像字段">
                <UInput v-model="editingGeneric.fields.avatar" placeholder="avatar_url" class="font-mono" />
              </UFormField>
            </div>
          </div>
        </div>

        <template #footer>
          <div class="flex justify-end gap-2">
            <UButton color="neutral" variant="outline" @click="showAddForm = false; editingGeneric = null">取消</UButton>
            <UButton color="primary"
              :loading="editingGeneric?.slug ? genericSaving[editingGeneric.slug] : false"
              :disabled="!editingGeneric?.slug || !editingGeneric?.authUrl || !editingGeneric?.tokenUrl || !editingGeneric?.userInfoUrl"
              @click="saveGeneric(editingGeneric!)">
              保存提供商
            </UButton>
          </div>
        </template>
      </UCard>
    </div>
    </AdminPageContent>
  </AdminPageContainer>
</template>

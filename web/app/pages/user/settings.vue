<script setup lang="ts">
definePageMeta({ middleware: 'auth' })

const { t } = useI18n()
const { containerClass } = useContainerWidth()
const authStore = useAuthStore()
const { apiFetch } = useApiFetch()
const toast = useToast()
const route = useRoute()
const router = useRouter()

type Section = 'profile' | 'account' | 'password' | 'tokens' | 'history'
const validSections: Section[] = ['profile', 'account', 'password', 'tokens', 'history']
const activeSection = ref<Section>(
  validSections.includes(route.query.section as Section)
    ? (route.query.section as Section)
    : 'profile'
)

const navItems = computed(() => [
  { label: t('site.settings.section_profile'), value: 'profile' as Section, icon: 'i-tabler-user' },
  { label: t('site.settings.section_account'), value: 'account' as Section, icon: 'i-tabler-id' },
  { label: t('site.settings.section_password'), value: 'password' as Section, icon: 'i-tabler-lock' },
  { label: t('site.settings.section_tokens'), value: 'tokens' as Section, icon: 'i-tabler-key' },
  { label: '浏览历史', value: 'history' as Section, icon: 'i-tabler-history' },
])

watch(activeSection, (s) => { router.replace({ query: { section: s } }) })

// ── 修改密码 ──────────────────────────────────────────────────────────────
const pw = reactive({ current: '', next: '', confirm: '' })
const pwSaving = ref(false)
const pwError = ref('')

async function savePassword() {
  pwError.value = ''
  if (pw.next.length < 6) { pwError.value = t('site.settings.pw_min_length'); return }
  if (pw.next !== pw.confirm) { pwError.value = t('site.settings.pw_mismatch'); return }
  if (!authStore.user?.id) return
  pwSaving.value = true
  try {
    const body: Record<string, string> = { new_password: pw.next }
    if (authStore.user.has_password) body.old_password = pw.current
    await apiFetch(`/users/${authStore.user.id}/password`, { method: 'PUT', body })
    pw.current = ''; pw.next = ''; pw.confirm = ''
    authStore.setUser({ ...authStore.user!, has_password: true })
    toast.add({ title: authStore.user.has_password ? t('site.settings.pw_changed') : t('site.settings.pw_set'), color: 'success' })
  } catch (e: any) {
    pwError.value = e.message || t('site.settings.pw_change_failed')
  } finally {
    pwSaving.value = false
  }
}

function formatDate(s?: string) {
  if (!s) return '—'
  return new Date(s).toLocaleDateString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit' })
}
</script>

<template>
  <div :class="[containerClass, 'mx-auto px-4 md:px-6 py-6 md:py-10 pb-16']">
    <h1 class="text-xl font-bold text-highlighted mb-6">{{ $t('site.settings.title') }}</h1>

    <div class="flex gap-6 items-start">
      <!-- 侧边栏 -->
      <aside class="hidden md:block w-52 shrink-0 sticky top-6">
        <UCard :ui="{ body: 'p-0 sm:p-0' }">
          <div class="p-4 flex items-center gap-3 border-b border-default">
            <BaseAvatar
              :src="authStore.user?.avatar"
              :alt="authStore.user?.display_name || authStore.user?.username"
              size="md" />
            <div class="min-w-0">
              <p class="text-sm font-semibold text-highlighted truncate">
                {{ authStore.user?.display_name || authStore.user?.username }}
              </p>
              <p class="text-xs text-muted truncate">@{{ authStore.user?.username }}</p>
            </div>
          </div>
          <div class="p-1.5 space-y-0.5">
            <button
              v-for="item in navItems"
              :key="item.value"
              class="flex items-center gap-2.5 px-3 py-2 rounded-md text-sm font-medium transition-all text-left w-full"
              :class="activeSection === item.value ? 'bg-primary/10 text-primary' : 'text-muted hover:text-highlighted hover:bg-muted'"
              @click="activeSection = item.value">
              <UIcon :name="item.icon" class="size-4 shrink-0" />
              {{ item.label }}
            </button>
            <div class="border-t border-default my-1.5" />
            <NuxtLink
              :to="`/user/${authStore.user?.id}`"
              class="flex items-center gap-2.5 px-3 py-2 rounded-md text-sm font-medium text-muted hover:text-highlighted hover:bg-muted transition-all">
              <UIcon name="i-tabler-user" class="size-4 shrink-0" />
              {{ $t('site.settings.view_my_profile') }}
            </NuxtLink>
          </div>
        </UCard>
      </aside>

      <!-- 主内容 -->
      <div class="flex-1 min-w-0">
        <!-- 移动端 Tab -->
        <div class="md:hidden flex gap-1 mb-4 overflow-x-auto pb-1">
          <UButton
            v-for="item in navItems"
            :key="item.value"
            :icon="item.icon"
            color="neutral"
            :variant="activeSection === item.value ? 'outline' : 'ghost'"
            size="sm"
            :class="activeSection === item.value ? 'bg-default shadow text-highlighted shrink-0' : 'text-muted shrink-0'"
            @click="activeSection = item.value">
            {{ item.label }}
          </UButton>
        </div>

        <div class="space-y-5">

          <!-- ── 公开资料 ── -->
          <SettingsProfileSection v-if="activeSection === 'profile'" />

          <!-- ── 账号信息 ── -->
          <template v-else-if="activeSection === 'account'">
            <UCard>
              <template #header>
                <h2 class="font-semibold text-highlighted flex items-center gap-2">
                  <UIcon name="i-tabler-id" class="size-5 text-primary" />
                  {{ $t('site.settings.account_title') }}
                </h2>
              </template>
              <div class="space-y-4">
                <div class="flex items-center justify-between py-3 border-b border-default">
                  <div>
                    <p class="text-sm font-medium text-highlighted">{{ $t('site.settings.username_label') }}</p>
                    <p class="text-sm text-muted mt-0.5">@{{ authStore.user?.username }}</p>
                  </div>
                  <UBadge :label="$t('site.settings.username_immutable_badge')" color="neutral" variant="soft" size="sm" />
                </div>
                <div class="flex items-center justify-between py-3 border-b border-default">
                  <div>
                    <p class="text-sm font-medium text-highlighted">{{ $t('site.settings.email_label') }}</p>
                    <p class="text-sm text-muted mt-0.5">{{ authStore.user?.email }}</p>
                  </div>
                  <UBadge :label="$t('site.settings.email_verified')" color="success" variant="soft" size="sm" />
                </div>
                <div class="flex items-center justify-between py-3 border-b border-default">
                  <div>
                    <p class="text-sm font-medium text-highlighted">{{ $t('site.settings.role_label') }}</p>
                    <p class="text-sm text-muted mt-0.5">
                      {{ ['', $t('site.settings.roles.user'), $t('site.settings.roles.editor'), $t('site.settings.roles.admin'), $t('site.settings.roles.super_admin')][authStore.user?.role ?? 1] }}
                    </p>
                  </div>
                  <UBadge
                    :label="['', $t('site.settings.role_badges.subscriber'), $t('site.settings.role_badges.editor'), $t('site.settings.role_badges.admin'), $t('site.settings.role_badges.super')][authStore.user?.role ?? 1]"
                    color="primary" variant="soft" size="sm" />
                </div>
                <div class="flex items-center justify-between py-3">
                  <div>
                    <p class="text-sm font-medium text-highlighted">{{ $t('site.settings.register_time_label') }}</p>
                    <p class="text-sm text-muted mt-0.5">{{ formatDate(authStore.user?.created_at) }}</p>
                  </div>
                </div>
              </div>
            </UCard>

            <UCard>
              <template #header>
                <h2 class="font-semibold text-error flex items-center gap-2">
                  <UIcon name="i-tabler-alert-triangle" class="size-5" />
                  {{ $t('site.settings.danger_zone') }}
                </h2>
              </template>
              <div class="flex items-center justify-between">
                <div>
                  <p class="text-sm font-medium text-highlighted">{{ $t('site.settings.logout_label') }}</p>
                  <p class="text-xs text-muted mt-0.5">{{ $t('site.settings.logout_desc') }}</p>
                </div>
                <UButton color="error" variant="soft" size="sm" icon="i-tabler-logout"
                  @click="authStore.logout().then(() => navigateTo('/'))">
                  {{ $t('site.settings.logout_btn') }}
                </UButton>
              </div>
            </UCard>
          </template>

          <!-- ── 修改密码 ── -->
          <template v-else-if="activeSection === 'password'">
            <UCard>
              <template #header>
                <h2 class="font-semibold text-highlighted flex items-center gap-2">
                  <UIcon name="i-tabler-lock" class="size-5 text-primary" />
                  {{ authStore.user?.has_password ? $t('site.settings.password_title') : $t('site.settings.set_password_title') }}
                </h2>
              </template>
              <form class="space-y-4" @submit.prevent="savePassword">
                <UAlert v-if="!authStore.user?.has_password" color="neutral" variant="subtle"
                  icon="i-tabler-info-circle" title="你当前通过第三方登录，设置密码后也可以用账号 + 密码直接登录" />
                <UAlert v-if="pwError" color="error" :title="pwError" icon="i-tabler-alert-circle" />
                <UFormField v-if="authStore.user?.has_password" :label="$t('site.settings.current_pw')">
                  <UInput v-model="pw.current" type="password" :placeholder="$t('site.settings.current_pw_placeholder')" class="w-full" autocomplete="current-password" />
                </UFormField>
                <UFormField :label="$t('site.settings.new_pw')">
                  <UInput v-model="pw.next" type="password" :placeholder="$t('site.settings.new_pw_placeholder')" class="w-full" autocomplete="new-password" />
                </UFormField>
                <UFormField :label="$t('site.settings.confirm_pw')">
                  <UInput v-model="pw.confirm" type="password" :placeholder="$t('site.settings.confirm_pw_placeholder')" class="w-full" autocomplete="new-password" />
                </UFormField>
                <div class="flex justify-end pt-2">
                  <UButton type="submit" color="primary" :loading="pwSaving"
                    :disabled="(!!authStore.user?.has_password && !pw.current) || !pw.next || !pw.confirm"
                    icon="i-tabler-lock-check">
                    {{ authStore.user?.has_password ? $t('site.settings.update_pw') : $t('site.settings.set_pw') }}
                  </UButton>
                </div>
              </form>
            </UCard>
          </template>

          <!-- ── 浏览历史 ── -->
          <SettingsHistorySection v-else-if="activeSection === 'history'" />

          <!-- ── 访问令牌 ── -->
          <SettingsTokensSection v-else-if="activeSection === 'tokens'" />

        </div>
      </div>
    </div>
  </div>
</template>

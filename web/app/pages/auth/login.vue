<script setup lang="ts">
definePageMeta({ layout: 'default' })

const authStore = useAuthStore()
const router = useRouter()
const route = useRoute()
const { t } = useI18n()
const config = useRuntimeConfig()

// Redirect if already logged in
if (authStore.isLoggedIn) {
  await navigateTo(route.query.redirect as string || '/')
}

const form = reactive({ login: '', password: '' })
const loading = ref(false)
const error = ref('')

const submit = async () => {
  error.value = ''
  loading.value = true
  try {
    await authStore.login(form.login, form.password)
    await navigateTo(route.query.redirect as string || '/')
  } catch (e: any) {
    error.value = e?.message || t('auth.login.failed')
  } finally {
    loading.value = false
  }
}

// Registration allowed?
const { apiFetch } = useApiFetch()
const allowRegistration = ref(true)
onMounted(async () => {
  try {
    const res = await apiFetch<{ value: string }>('/options/allow_registration')
    allowRegistration.value = JSON.parse(res.value) !== false
  } catch { /* default: allow */ }
})

// OAuth
const authApi = useAuthApi()
const { data: oauthProviders } = await useAsyncData('oauth-providers', () =>
  authApi.getOAuthProviders().catch(() => [] as string[])
)

// Provider display metadata — add a new entry here when registering a new OAuth provider.
// Slug must match the provider's Name() return value in server/internal/oauth/{slug}.go.
const providerMeta: Record<string, { label: string; icon: string }> = {
  github: { label: 'GitHub', icon: 'i-tabler-brand-github' },
  google: { label: 'Google', icon: 'i-tabler-brand-google' },
  qq:     { label: 'QQ',     icon: 'i-tabler-brand-qq'     },
  // ↓ Add new providers here
  // weibo:  { label: '微博',   icon: 'i-tabler-brand-weibo'  },
}

const redirectToOAuth = (provider: string) => {
  const apiBase = (config.public.apiBase as string).replace(/\/api\/v1\/?$/, '')
  const redirect = (route.query.redirect as string) || '/'
  const url = `${apiBase}/api/v1/auth/oauth/${provider}/redirect?redirect=${encodeURIComponent(redirect)}`
  window.location.href = url
}
</script>

<template>
  <div class="min-h-[60vh] flex items-center justify-center px-4 py-12">
    <div class="w-full max-w-sm">
      <div class="text-center mb-8">
        <div class="size-14 rounded-md bg-primary/10 flex items-center justify-center mx-auto mb-4">
          <UIcon name="i-tabler-pencil" class="size-7 text-primary" />
        </div>
        <h1 class="text-2xl font-bold text-highlighted">{{ $t('auth.login.title') }}</h1>
        <p class="text-sm text-muted mt-1">{{ $t('auth.login.subtitle') }}</p>
      </div>

      <UCard>
        <form class="space-y-4" @submit.prevent="submit">
          <UAlert
            v-if="error"
            color="error"
            variant="subtle"
            :title="error"
            icon="i-tabler-alert-circle"
          />

          <UFormField :label="$t('auth.login.username_label')">
            <UInput
              v-model="form.login"
              :placeholder="$t('auth.login.username_placeholder')"
              icon="i-tabler-user"
              autofocus
              class="w-full"
            />
          </UFormField>

          <UFormField :label="$t('auth.login.password_label')">
            <UInput
              v-model="form.password"
              type="password"
              :placeholder="$t('auth.login.password_placeholder')"
              icon="i-tabler-lock"
              class="w-full"
            />
          </UFormField>

          <UButton
            type="submit"
            color="primary"
            block
            :loading="loading"
            :disabled="!form.login || !form.password"
          >
            {{ $t('auth.login.submit') }}
          </UButton>
        </form>

        <!-- OAuth providers -->
        <template v-if="oauthProviders && oauthProviders.length > 0">
          <div class="flex items-center gap-3 my-4">
            <div class="flex-1 border-t border-default" />
            <span class="text-xs text-muted shrink-0">{{ $t('auth.login.or_continue_with') }}</span>
            <div class="flex-1 border-t border-default" />
          </div>
          <div class="flex flex-col gap-2">
            <UButton
              v-for="provider in oauthProviders"
              :key="provider"
              color="neutral"
              variant="outline"
              block
              :icon="providerMeta[provider]?.icon ?? 'i-tabler-login'"
              @click="redirectToOAuth(provider)"
            >
              {{ $t('auth.login.login_with', { provider: providerMeta[provider]?.label ?? provider }) }}
            </UButton>
          </div>
        </template>
      </UCard>

      <p v-if="allowRegistration" class="text-center text-sm text-muted mt-4">
        {{ $t('auth.login.no_account') }}
        <NuxtLink to="/auth/register" class="text-primary hover:underline font-medium">{{ $t('auth.login.register_link') }}</NuxtLink>
      </p>
    </div>
  </div>
</template>

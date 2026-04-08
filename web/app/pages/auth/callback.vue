<script setup lang="ts">
definePageMeta({ layout: 'default' })

const route = useRoute()
const authStore = useAuthStore()
const authApi = useAuthApi()
const { t } = useI18n()

const error = ref('')

onMounted(async () => {
  const { access_token, refresh_token, expires_in, redirect, error: oauthError } = route.query as Record<string, string>

  if (oauthError) {
    error.value = oauthError
    return
  }

  if (!access_token || !refresh_token) {
    error.value = t('auth.callback.missing_tokens')
    return
  }

  // Store tokens
  authStore.setTokens(access_token, refresh_token)

  // Fetch user info
  try {
    const res = await authApi.me()
    authStore.setUser(res.user)
    await navigateTo(redirect || '/', { replace: true })
  } catch {
    error.value = t('auth.callback.fetch_user_failed')
    authStore.setTokens('', '')
  }
})
</script>

<template>
  <div class="min-h-[60vh] flex items-center justify-center px-4">
    <div class="text-center">
      <template v-if="!error">
        <UIcon name="i-tabler-loader-2" class="animate-spin size-10 text-primary mx-auto mb-4" />
        <p class="text-muted">{{ $t('auth.callback.processing') }}</p>
      </template>
      <template v-else>
        <UIcon name="i-tabler-alert-circle" class="size-10 text-error mx-auto mb-4" />
        <h2 class="text-lg font-semibold text-highlighted mb-2">{{ $t('auth.callback.failed') }}</h2>
        <p class="text-sm text-muted mb-4">{{ error }}</p>
        <NuxtLink to="/auth/login">
          <UButton color="primary">{{ $t('auth.login.submit') }}</UButton>
        </NuxtLink>
      </template>
    </div>
  </div>
</template>

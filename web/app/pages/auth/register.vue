<script setup lang="ts">
definePageMeta({ layout: 'default' })

const authStore = useAuthStore()
if (authStore.isLoggedIn) await navigateTo('/')

const { apiFetch } = useApiFetch()
const { t } = useI18n()

const form = reactive({ username: '', email: '', password: '', display_name: '', code: '' })
const loading = ref(false)
const error = ref('')

// ── Registration allowed ───────────────────────────────────────────────────
const registrationAllowed = ref<boolean | null>(null) // null = loading

// ── Verify mode ───────────────────────────────────────────────────────────
const verifyMode = ref<'none' | 'email' | 'sms'>('none')
onMounted(async () => {
  try {
    const regRes = await apiFetch<{ value: string }>('/options/allow_registration')
    registrationAllowed.value = JSON.parse(regRes.value) !== false
  } catch {
    registrationAllowed.value = true
  }

  if (!registrationAllowed.value) return

  try {
    const res = await apiFetch<{ key: string; value: string }>('/options/auth_register_verify')
    const parsed = JSON.parse(res.value)
    verifyMode.value = parsed.mode ?? 'none'
  } catch {
    verifyMode.value = 'none'
  }
})

// ── Send code ─────────────────────────────────────────────────────────────
const codeSending = ref(false)
const countdown = ref(0)
let countdownTimer: ReturnType<typeof setInterval> | null = null

const codeTarget = computed(() =>
  verifyMode.value === 'email' ? form.email : form.username,
)

const canSendCode = computed(() =>
  countdown.value === 0 &&
  !codeSending.value &&
  (verifyMode.value === 'email' ? !!form.email : !!form.username),
)

async function sendCode() {
  if (!canSendCode.value) return
  codeSending.value = true
  try {
    await apiFetch('/auth/send-code', {
      method: 'POST',
      body: { type: verifyMode.value, target: codeTarget.value },
    })
    countdown.value = 60
    countdownTimer = setInterval(() => {
      countdown.value--
      if (countdown.value <= 0) {
        clearInterval(countdownTimer!)
        countdownTimer = null
      }
    }, 1000)
  } catch (e: any) {
    error.value = e?.message || t('auth.register.code_failed')
  } finally {
    codeSending.value = false
  }
}

onUnmounted(() => {
  if (countdownTimer) clearInterval(countdownTimer)
})

// ── Submit ────────────────────────────────────────────────────────────────
const submit = async () => {
  error.value = ''
  loading.value = true
  try {
    await authStore.register(form)
    await navigateTo('/')
  } catch (e: any) {
    error.value = e?.message || t('auth.register.failed')
  } finally {
    loading.value = false
  }
}

const submitDisabled = computed(() => {
  if (!form.username || !form.email || !form.password) return true
  if (verifyMode.value !== 'none' && !form.code) return true
  return false
})
</script>

<template>
  <div class="min-h-[60vh] flex items-center justify-center px-4 py-12">
    <div class="w-full max-w-sm">

      <!-- Registration disabled -->
      <template v-if="registrationAllowed === false">
        <div class="text-center mb-8">
          <div class="size-14 rounded-md bg-muted flex items-center justify-center mx-auto mb-4">
            <UIcon name="i-tabler-user-off" class="size-7 text-muted" />
          </div>
          <h1 class="text-2xl font-bold text-highlighted">暂不开放注册</h1>
          <p class="text-sm text-muted mt-1">该站点当前未开放注册，请联系管理员</p>
        </div>
        <UCard>
          <div class="text-center py-2 space-y-4">
            <p class="text-sm text-muted">已有账号？</p>
            <UButton to="/auth/login" color="primary" block>前往登录</UButton>
          </div>
        </UCard>
      </template>

      <!-- Registration form -->
      <template v-else-if="registrationAllowed === true">
      <div class="text-center mb-8">
        <div class="size-14 rounded-md bg-primary/10 flex items-center justify-center mx-auto mb-4">
          <UIcon name="i-tabler-user-plus" class="size-7 text-primary" />
        </div>
        <h1 class="text-2xl font-bold text-highlighted">{{ $t('auth.register.title') }}</h1>
        <p class="text-sm text-muted mt-1">{{ $t('auth.register.subtitle') }}</p>
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

          <UFormField :label="$t('auth.register.username_label')">
            <UInput v-model="form.username" :placeholder="$t('auth.register.username_placeholder')" icon="i-tabler-at" class="w-full" />
          </UFormField>

          <UFormField :label="$t('auth.register.display_name_label')">
            <UInput v-model="form.display_name" :placeholder="$t('auth.register.display_name_placeholder')" icon="i-tabler-user" class="w-full" />
          </UFormField>

          <UFormField :label="$t('auth.register.email_label')">
            <UInput v-model="form.email" type="email" placeholder="your@email.com" icon="i-tabler-mail" class="w-full" />
          </UFormField>

          <UFormField :label="$t('auth.register.password_label')">
            <UInput v-model="form.password" type="password" :placeholder="$t('auth.register.password_placeholder')" icon="i-tabler-lock" class="w-full" />
          </UFormField>

          <!-- Verification code (shown when mode != none) -->
          <template v-if="verifyMode !== 'none'">
            <UFormField :label="verifyMode === 'email' ? $t('auth.register.email_code_label') : $t('auth.register.sms_code_label')">
              <div class="flex gap-2">
                <UInput
                  v-model="form.code"
                  :placeholder="$t('auth.register.code_placeholder')"
                  icon="i-tabler-shield-check"
                  class="flex-1"
                  maxlength="6"
                />
                <UButton
                  color="neutral"
                  variant="outline"
                  :loading="codeSending"
                  :disabled="!canSendCode"
                  @click="sendCode"
                >
                  {{ countdown > 0 ? `${countdown}s` : $t('auth.register.send_code') }}
                </UButton>
              </div>
              <p class="text-xs text-muted mt-1">
                {{ verifyMode === 'email' ? $t('auth.register.email_tip') : $t('auth.register.sms_tip') }}
              </p>
            </UFormField>
          </template>

          <UButton
            type="submit"
            color="primary"
            block
            :loading="loading"
            :disabled="submitDisabled"
          >
            {{ $t('auth.register.submit') }}
          </UButton>
        </form>
      </UCard>

      <p class="text-center text-sm text-muted mt-4">
        {{ $t('auth.register.has_account') }}
        <NuxtLink to="/auth/login" class="text-primary hover:underline font-medium">{{ $t('auth.register.login_link') }}</NuxtLink>
      </p>
      </template>

    </div>
  </div>
</template>

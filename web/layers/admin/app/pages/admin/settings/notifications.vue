<script setup lang="ts">
const optionApi = useOptionApi()
const toast = useToast()
const { t } = useI18n()

// ── State ──────────────────────────────────────────────────────────────────
const isLoading = ref(true)
const emailSaving = ref(false)
const smsSaving = ref(false)
const verifySaving = ref(false)
const testEmailSending = ref(false)
const showEmailPassword = ref(false)
const showSmsSecret = ref(false)

const verifyMode = ref('none')
const verifyModeSelected = ref('none')
watch(verifyModeSelected, (v) => {
  verifyMode.value = typeof v === 'object' ? (v as any)?.value ?? 'none' : v
})

const email = reactive({
  host: '',
  port: 587,
  username: '',
  password: '',
  from: 'Blog <noreply@example.com>',
  site_name: 'My Blog',
  site_url: 'http://localhost:3000',
})

const sms = reactive({
  provider: '',
  access_key_id: '',
  access_key_secret: '',
  sign_name: '',
  template_code: '',
})

// USelect does not support empty-string item values (Reka UI SelectItem quirk).
// Use 'disabled' as the sentinel for "not enabled" and map to '' when saving.
const smsProviderSelected = ref('disabled')
watch(smsProviderSelected, (v) => {
  const raw = typeof v === 'object' ? (v as any)?.value ?? 'disabled' : (v || 'disabled')
  sms.provider = raw === 'disabled' ? '' : raw
})

const smsProviders = computed(() => [
  { label: t('admin.settings.notifications.sms_provider_disabled'), value: 'disabled' },
  { label: t('admin.settings.notifications.sms_provider_aliyun'), value: 'aliyun' },
  { label: t('admin.settings.notifications.sms_provider_tencent'), value: 'tencent' },
  { label: 'Twilio', value: 'twilio' },
])

// ── Load ───────────────────────────────────────────────────────────────────
onMounted(async () => {
  try {
    const [emailData, smsData, verifyData] = await Promise.allSettled([
      optionApi.getOption('notify_email'),
      optionApi.getOption('notify_sms'),
      optionApi.getOption('auth_register_verify'),
    ])
    if (emailData.status === 'fulfilled' && emailData.value) {
      Object.assign(email, emailData.value)
    }
    if (smsData.status === 'fulfilled' && smsData.value) {
      Object.assign(sms, smsData.value)
      smsProviderSelected.value = sms.provider || 'disabled'
    }
    if (verifyData.status === 'fulfilled' && verifyData.value) {
      verifyMode.value = verifyData.value.mode ?? 'none'
      verifyModeSelected.value = verifyMode.value
    }
  } finally {
    isLoading.value = false
  }
})

// ── Save ───────────────────────────────────────────────────────────────────
async function saveEmail() {
  emailSaving.value = true
  try {
    await optionApi.setOption('notify_email', { ...email, port: Number(email.port) })
    toast.add({ title: t('admin.settings.notifications.email_saved'), color: 'success' })
  } catch (e: any) {
    toast.add({ title: t('common.save_failed'), description: e.message, color: 'error' })
  } finally {
    emailSaving.value = false
  }
}

async function saveSms() {
  smsSaving.value = true
  try {
    await optionApi.setOption('notify_sms', { ...sms })
    toast.add({ title: t('admin.settings.notifications.sms_saved'), color: 'success' })
  } catch (e: any) {
    toast.add({ title: t('common.save_failed'), description: e.message, color: 'error' })
  } finally {
    smsSaving.value = false
  }
}

async function saveVerify() {
  verifySaving.value = true
  try {
    await optionApi.setOption('auth_register_verify', { mode: verifyMode.value })
    toast.add({ title: t('admin.settings.notifications.verify_saved'), color: 'success' })
  } catch (e: any) {
    toast.add({ title: t('common.save_failed'), description: e.message, color: 'error' })
  } finally {
    verifySaving.value = false
  }
}

// ── Derived ────────────────────────────────────────────────────────────────
const emailEnabled = computed(() => !!email.host)
const smsEnabled = computed(() => !!sms.provider)

const verifyModeOptions = computed(() => [
  { label: t('admin.settings.notifications.verify_none'), value: 'none' },
  { label: t('admin.settings.notifications.verify_email'), value: 'email' },
  { label: t('admin.settings.notifications.verify_sms'), value: 'sms' },
])
</script>

<template>
  <AdminPageContainer>
    <AdminPageHeader :title="$t('admin.settings.notifications.title')" :subtitle="$t('admin.settings.notifications.subtitle')" />
    <AdminPageContent>

      <!-- Skeleton -->
      <div v-if="isLoading" class="space-y-4">
        <UCard v-for="i in 2" :key="i">
          <template #header><USkeleton class="h-5 w-40" /></template>
          <div class="space-y-4">
            <div v-for="j in 4" :key="j" class="space-y-2">
              <USkeleton class="h-4 w-24" />
              <USkeleton class="h-9 w-full rounded-md" />
            </div>
          </div>
        </UCard>
      </div>

      <template v-else>
        <!-- ── Email ──────────────────────────────────────────────── -->
        <UCard>
          <template #header>
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2">
                <UIcon name="i-tabler-mail" class="size-5 text-primary" />
                <h3 class="text-lg font-semibold text-highlighted">{{ $t('admin.settings.notifications.email_title') }}</h3>
              </div>
              <UBadge
                :label="emailEnabled ? $t('admin.settings.notifications.email_enabled') : $t('admin.settings.notifications.email_disabled')"
                :color="emailEnabled ? 'success' : 'neutral'"
                variant="soft" />
            </div>
          </template>

          <div class="space-y-4">
            <UAlert
              icon="i-tabler-info-circle"
              color="primary"
              variant="soft"
              :title="$t('admin.settings.notifications.email_tip_title')"
              :description="$t('admin.settings.notifications.email_tip_desc')" />

            <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
              <UFormField :label="$t('admin.settings.notifications.smtp_host_label')" class="md:col-span-2">
                <UInput v-model="email.host" placeholder="smtp.qq.com" class="w-full" />
                <p class="text-xs text-muted mt-1">{{ $t('admin.settings.notifications.smtp_host_hint') }}</p>
              </UFormField>
              <UFormField :label="$t('admin.settings.notifications.smtp_port_label')">
                <UInput v-model="email.port" type="number" placeholder="587" class="w-full" />
                <p class="text-xs text-muted mt-1">{{ $t('admin.settings.notifications.smtp_port_hint') }}</p>
              </UFormField>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <UFormField :label="$t('admin.settings.notifications.smtp_user_label')">
                <UInput v-model="email.username" placeholder="noreply@example.com" autocomplete="off" class="w-full" />
              </UFormField>
              <UFormField :label="$t('admin.settings.notifications.smtp_password_label')">
                <div class="flex gap-2">
                  <UInput
                    v-model="email.password"
                    :type="showEmailPassword ? 'text' : 'password'"
                    placeholder="••••••••"
                    autocomplete="new-password"
                    class="w-full" />
                  <UButton
                    color="neutral" variant="outline" square
                    :icon="showEmailPassword ? 'i-tabler-eye-off' : 'i-tabler-eye'"
                    @click="showEmailPassword = !showEmailPassword" />
                </div>
              </UFormField>
            </div>

            <UFormField :label="$t('admin.settings.notifications.from_label')">
              <UInput v-model="email.from" placeholder="Blog <noreply@example.com>" class="w-full" />
              <p class="text-xs text-muted mt-1">{{ $t('admin.settings.notifications.from_hint') }}</p>
            </UFormField>

            <USeparator />

            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <UFormField :label="$t('admin.settings.notifications.site_name_label')">
                <UInput v-model="email.site_name" placeholder="My Blog" class="w-full" />
                <p class="text-xs text-muted mt-1">{{ $t('admin.settings.notifications.site_name_hint') }}</p>
              </UFormField>
              <UFormField :label="$t('admin.settings.notifications.site_url_label')">
                <UInput v-model="email.site_url" placeholder="https://myblog.com" class="w-full" />
                <p class="text-xs text-muted mt-1">{{ $t('admin.settings.notifications.site_url_hint') }}</p>
              </UFormField>
            </div>
          </div>

          <template #footer>
            <div class="flex justify-end">
              <UButton color="primary" :loading="emailSaving" icon="i-tabler-device-floppy" @click="saveEmail">
                {{ $t('admin.settings.notifications.save_email') }}
              </UButton>
            </div>
          </template>
        </UCard>

        <!-- ── SMS ───────────────────────────────────────────────── -->
        <UCard>
          <template #header>
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2">
                <UIcon name="i-tabler-message" class="size-5 text-primary" />
                <h3 class="text-lg font-semibold text-highlighted">{{ $t('admin.settings.notifications.sms_title') }}</h3>
              </div>
              <UBadge
                :label="smsEnabled ? $t('admin.settings.notifications.email_enabled') : $t('admin.settings.notifications.email_disabled')"
                :color="smsEnabled ? 'success' : 'neutral'"
                variant="soft" />
            </div>
          </template>

          <div class="space-y-4">
            <UAlert
              icon="i-tabler-info-circle"
              color="primary"
              variant="soft"
              :title="$t('admin.settings.notifications.sms_tip_title')"
              :description="$t('admin.settings.notifications.sms_tip_desc')" />

            <UFormField :label="$t('admin.settings.notifications.sms_provider_label')">
              <USelect v-model="smsProviderSelected" :items="smsProviders" class="w-full md:w-60" />
              <p class="text-xs text-muted mt-1">{{ $t('admin.settings.notifications.sms_provider_hint') }}</p>
            </UFormField>

            <template v-if="sms.provider">
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <UFormField label="AccessKey ID">
                  <UInput v-model="sms.access_key_id" placeholder="LTAI5t..." autocomplete="off" class="w-full" />
                </UFormField>
                <UFormField label="AccessKey Secret">
                  <div class="flex gap-2">
                    <UInput
                      v-model="sms.access_key_secret"
                      :type="showSmsSecret ? 'text' : 'password'"
                      placeholder="••••••••"
                      autocomplete="new-password"
                      class="w-full" />
                    <UButton
                      color="neutral" variant="outline" square
                      :icon="showSmsSecret ? 'i-tabler-eye-off' : 'i-tabler-eye'"
                      @click="showSmsSecret = !showSmsSecret" />
                  </div>
                </UFormField>
              </div>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <UFormField :label="$t('admin.settings.notifications.sms_sign_label')">
                  <UInput v-model="sms.sign_name" :placeholder="$t('admin.settings.notifications.sms_sign_label')" class="w-full" />
                </UFormField>
                <UFormField :label="$t('admin.settings.notifications.sms_template_label')">
                  <UInput v-model="sms.template_code" placeholder="SMS_123456789" class="w-full" />
                </UFormField>
              </div>
            </template>
          </div>

          <template #footer>
            <div class="flex justify-end">
              <UButton color="primary" :loading="smsSaving" icon="i-tabler-device-floppy" @click="saveSms">
                {{ $t('admin.settings.notifications.save_sms') }}
              </UButton>
            </div>
          </template>
        </UCard>

        <!-- ── Register Verify ────────────────────────────────── -->
        <UCard>
          <template #header>
            <div class="flex items-center gap-2">
              <UIcon name="i-tabler-shield-check" class="size-5 text-primary" />
              <h3 class="text-lg font-semibold text-highlighted">{{ $t('admin.settings.notifications.verify_title') }}</h3>
            </div>
          </template>

          <div class="space-y-4">
            <UAlert
              icon="i-tabler-info-circle"
              color="primary"
              variant="soft"
              :title="$t('admin.settings.notifications.verify_tip_title')"
              :description="$t('admin.settings.notifications.verify_tip_desc')" />

            <UFormField :label="$t('admin.settings.notifications.verify_mode_label')">
              <USelect v-model="verifyModeSelected" :items="verifyModeOptions" class="w-full md:w-60" />
              <p class="text-xs text-muted mt-1">
                {{ $t('admin.settings.notifications.verify_mode_hint') }}
              </p>
            </UFormField>
          </div>

          <template #footer>
            <div class="flex justify-end">
              <UButton color="primary" :loading="verifySaving" icon="i-tabler-device-floppy" @click="saveVerify">
                {{ $t('admin.settings.notifications.save_verify') }}
              </UButton>
            </div>
          </template>
        </UCard>
      </template>
    </AdminPageContent>
  </AdminPageContainer>
</template>

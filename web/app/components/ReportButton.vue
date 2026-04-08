<script setup lang="ts">
const props = defineProps<{
  targetType: 'post' | 'comment' | 'user'
  targetId: number
}>()

const { t } = useI18n()
const authStore = useAuthStore()
const router = useRouter()
const route = useRoute()
const toast = useToast()
const reportApi = useReportApi()

const open = ref(false)
const reason = ref('')
const detail = ref('')
const submitting = ref(false)

const REASONS = computed(() => [
  t('site.report.reason_spam'),
  t('site.report.reason_illegal'),
  t('site.report.reason_porn'),
  t('site.report.reason_attack'),
  t('site.report.reason_false'),
  t('site.report.reason_copyright'),
  t('site.report.reason_other'),
])

const submit = async () => {
  if (!authStore.isLoggedIn) {
    router.push(`/auth/login?redirect=${encodeURIComponent(route.fullPath)}`)
    return
  }
  if (!reason.value) return
  submitting.value = true
  try {
    await reportApi.create(props.targetType, props.targetId, reason.value, detail.value)
    toast.add({ title: t('site.report.submitted'), color: 'success' })
    open.value = false
    reason.value = ''
    detail.value = ''
  } catch (e: any) {
    toast.add({ title: e?.message || t('site.report.failed'), color: 'error' })
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <UModal v-model:open="open" :ui="{ content: 'max-w-md' }">
    <UButton
      color="neutral" variant="ghost" size="xs"
      icon="i-tabler-flag"
      @click="open = true">
      {{ $t('site.report.button') }}
    </UButton>

    <template #content>
      <div class="p-5 space-y-4">
        <div>
          <h3 class="font-semibold text-highlighted">{{ $t('site.report.title') }}</h3>
          <p class="text-xs text-muted mt-1">{{ $t('site.report.subtitle') }}</p>
        </div>

        <div class="space-y-2">
          <p class="text-sm font-medium">{{ $t('site.report.reason_label') }}</p>
          <div class="flex flex-wrap gap-2">
            <button
              v-for="r in REASONS" :key="r"
              class="px-3 py-1.5 rounded-full text-sm border transition-colors"
              :class="reason === r
                ? 'border-primary bg-primary/10 text-primary'
                : 'border-default text-muted hover:border-primary hover:text-primary'"
              @click="reason = r">
              {{ r }}
            </button>
          </div>
        </div>

        <UFormField :label="$t('site.report.detail_label')">
          <UTextarea v-model="detail" :placeholder="$t('site.report.detail_placeholder')" :rows="3" class="w-full" />
        </UFormField>

        <div class="flex gap-2 justify-end">
          <UButton color="neutral" variant="ghost" @click="open = false">{{ $t('common.cancel') }}</UButton>
          <UButton color="error" :loading="submitting" :disabled="!reason" @click="submit">
            {{ $t('site.report.submit') }}
          </UButton>
        </div>
      </div>
    </template>
  </UModal>
</template>

<template>
  <AdminPageContainer>
    <AdminPageHeader :title="$t('admin.appearance.customize.title')" :subtitle="$t('admin.appearance.customize.subtitle')">
      <template #actions>
        <UButton color="primary" icon="i-tabler-device-floppy" :loading="saving" @click="save">
          {{ $t('common.save_changes') }}
        </UButton>
      </template>
    </AdminPageHeader>

    <AdminPageContent>
      <template v-if="loading">
        <div class="space-y-4">
          <UCard v-for="i in 3" :key="i">
            <template #header><USkeleton class="h-5 w-32" /></template>
            <div class="space-y-3">
              <USkeleton class="h-32 w-full" />
            </div>
          </UCard>
        </div>
      </template>

      <template v-else>
        <!-- 自定义 CSS -->
        <UCard>
          <template #header>
            <h3 class="text-base font-semibold text-highlighted">{{ $t('admin.appearance.customize.custom_css_title') }}</h3>
          </template>
          <UFormField :label="$t('admin.appearance.customize.custom_css_label')">
            <UTextarea
              v-model="form.customCss"
              :rows="10"
              placeholder="/* 直接写 CSS，无需 &lt;style&gt; 标签 */&#10;&#10;.prose { line-height: 1.8; }&#10;&#10;.prose h2 { font-size: 1.5rem; }"
              class="w-full font-mono text-sm"
            />
          </UFormField>
          <UAlert
            class="mt-3"
            icon="i-tabler-alert-triangle"
            color="warning"
            variant="soft"
            :title="$t('admin.appearance.customize.custom_css_warning_title')"
            :description="$t('admin.appearance.customize.custom_css_warning_desc')" />
        </UCard>

        <!-- 页头代码 -->
        <UCard>
          <template #header>
            <h3 class="text-base font-semibold text-highlighted">{{ $t('admin.appearance.customize.head_code_title') }}</h3>
          </template>
          <div class="space-y-4">
            <UFormField :label="$t('admin.appearance.customize.head_code_label')" :hint="$t('admin.appearance.customize.head_code_hint')">
              <UTextarea
                v-model="form.headCode"
                :rows="6"
                placeholder="<!-- 需要完整的 HTML 标签，例如：-->&#10;<script async src=&quot;https://www.googletagmanager.com/gtag/js&quot;></script>&#10;<script>&#10;  window.dataLayer = window.dataLayer || [];&#10;</script>&#10;<link rel=&quot;preconnect&quot; href=&quot;https://fonts.googleapis.com&quot;>"
                class="w-full font-mono text-sm"
              />
            </UFormField>
            <UAlert
              icon="i-tabler-info-circle"
              color="info"
              variant="soft"
              :title="$t('admin.appearance.customize.head_code_tip_title')"
              :description="$t('admin.appearance.customize.head_code_tip_desc')"
            />
          </div>
        </UCard>

        <!-- 页脚代码 -->
        <UCard>
          <template #header>
            <h3 class="text-base font-semibold text-highlighted">{{ $t('admin.appearance.customize.body_code_title') }}</h3>
          </template>
          <UFormField :label="$t('admin.appearance.customize.body_code_label')" :hint="$t('admin.appearance.customize.body_code_hint')">
            <UTextarea
              v-model="form.bodyCode"
              :rows="6"
              placeholder="<!-- 需要完整的 HTML 标签，例如：-->&#10;<script>&#10;  console.log('页面加载完成');&#10;</script>"
              class="w-full font-mono text-sm"
            />
          </UFormField>
        </UCard>
      </template>
    </AdminPageContent>
  </AdminPageContainer>
</template>

<script setup lang="ts">
const toast        = useToast()
const { t }        = useI18n()
const optionApi    = useOptionApi()
const optionsStore = useOptionsStore()

const rawLoading = ref(true)
const loading    = useMinLoading(rawLoading)
const saving     = ref(false)

const form = ref({
  customCss: '',
  headCode:  '',
  bodyCode:  '',
})

onMounted(async () => {
  try {
    await optionsStore.load()
    form.value = {
      customCss: optionsStore.get('theme_custom_css', ''),
      headCode:  optionsStore.get('theme_head_code',  ''),
      bodyCode:  optionsStore.get('theme_body_code',  ''),
    }
  } finally {
    rawLoading.value = false
  }
})

const save = async () => {
  saving.value = true
  try {
    await Promise.all([
      optionApi.setOption('theme_custom_css', form.value.customCss),
      optionApi.setOption('theme_head_code',  form.value.headCode),
      optionApi.setOption('theme_body_code',  form.value.bodyCode),
    ])
    toast.add({ title: t('admin.appearance.customize.saved'), color: 'success' })
  } catch (error: any) {
    toast.add({ title: t('common.save_failed'), description: error?.message, color: 'error' })
  } finally {
    saving.value = false
  }
}
</script>

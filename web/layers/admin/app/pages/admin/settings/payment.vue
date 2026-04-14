<script setup lang="ts">
import type { PaymentProvider } from '~/composables/usePaymentApi'

const { listProviders, setProviderConfig } = usePaymentApi()
const toast = useToast()

const isLoading = ref(true)
const savingSlug = ref<string | null>(null)
const providers = ref<PaymentProvider[]>([])
// Each provider's form data, keyed by slug
const forms = ref<Record<string, Record<string, any>>>({})

const load = async () => {
  isLoading.value = true
  try {
    const res = await listProviders()
    providers.value = res.items
    for (const p of res.items) {
      forms.value[p.slug] = { ...p.config }
    }
  } catch {
    // keep empty
  } finally {
    isLoading.value = false
  }
}

const save = async (slug: string) => {
  savingSlug.value = slug
  try {
    const res = await setProviderConfig(slug, forms.value[slug])
    // Update provider info from response
    const idx = providers.value.findIndex(p => p.slug === slug)
    if (idx >= 0) {
      providers.value[idx] = res
      forms.value[slug] = { ...res.config }
    }
    toast.add({ title: '保存成功', color: 'success' })
  } catch (e: any) {
    toast.add({ title: '保存失败', description: e?.message, color: 'error' })
  } finally {
    savingSlug.value = null
  }
}

const reset = (slug: string) => {
  const p = providers.value.find(p => p.slug === slug)
  if (p) forms.value[slug] = { ...p.config }
}

onMounted(load)
</script>

<template>
  <AdminPageContainer>
    <AdminPageHeader title="支付设置" subtitle="管理支付方式，启用后可用于文章付费等场景" />
    <AdminPageContent>
      <div v-if="isLoading" class="flex justify-center py-12">
        <UIcon name="i-tabler-loader-2" class="animate-spin text-2xl text-muted" />
      </div>

      <div v-else class="space-y-6">
        <UCard v-for="provider in providers" :key="provider.slug">
          <template #header>
            <div class="flex items-center gap-2">
              <UIcon :name="provider.icon" class="text-xl" />
              <span class="font-semibold">{{ provider.label }}</span>
              <UBadge
                :label="forms[provider.slug]?.enabled ? '已启用' : '未启用'"
                :color="forms[provider.slug]?.enabled ? 'success' : 'neutral'"
                variant="soft" size="sm"
              />
            </div>
          </template>

          <div v-if="forms[provider.slug]" class="space-y-4">
            <template v-for="field in provider.fields" :key="field.key">
              <!-- switch -->
              <div v-if="field.type === 'switch'"
                class="flex items-center justify-between rounded-md border border-default px-4 py-3">
                <p class="text-sm font-medium text-highlighted">{{ field.label }}</p>
                <USwitch v-model="forms[provider.slug][field.key]" />
              </div>

              <!-- text / password -->
              <UFormField v-else-if="field.type === 'text' || field.type === 'password'"
                :label="field.label" :required="field.required">
                <form v-if="field.type === 'password'" @submit.prevent>
                  <input type="text" autocomplete="username" class="hidden" aria-hidden="true" tabindex="-1" />
                  <UInput
                    v-model="forms[provider.slug][field.key]"
                    type="password"
                    :placeholder="field.placeholder"
                    class="w-full font-mono"
                  />
                </form>
                <UInput
                  v-else
                  v-model="forms[provider.slug][field.key]"
                  type="text"
                  :placeholder="field.placeholder"
                  class="w-full font-mono"
                />
              </UFormField>

              <!-- select -->
              <UFormField v-else-if="field.type === 'select'" :label="field.label">
                <USelect
                  v-model="forms[provider.slug][field.key]"
                  :items="field.options || []"
                  class="w-full"
                />
              </UFormField>
            </template>
          </div>

          <template #footer>
            <div class="flex justify-end gap-2">
              <UButton color="neutral" variant="ghost" @click="reset(provider.slug)">重置</UButton>
              <UButton color="primary"
                :loading="savingSlug === provider.slug"
                @click="save(provider.slug)">
                保存
              </UButton>
            </div>
          </template>
        </UCard>

        <div v-if="providers.length === 0"
          class="rounded-md border border-dashed border-default p-8 text-center text-sm text-muted">
          暂无可用的支付方式
        </div>
      </div>
    </AdminPageContent>
  </AdminPageContainer>
</template>

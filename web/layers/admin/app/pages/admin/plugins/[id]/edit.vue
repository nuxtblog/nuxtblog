<script setup lang="ts">
import type { PluginManifest, PipelineDef } from '~/composables/usePluginApi'

definePageMeta({ layout: 'admin' })
const { t } = useI18n()
const route = useRoute()
const pluginId = route.params.id as string
const pluginApi = usePluginApi()
const toast = useToast()

const loading = ref(true)
const saving = ref(false)
const pipelines = ref<PipelineDef[]>([])
let rawManifest: PluginManifest = {}

async function load() {
  loading.value = true
  try {
    const res = await pluginApi.getManifest(pluginId)
    try {
      rawManifest = JSON.parse(res.manifest || '{}')
    } catch {
      rawManifest = {}
    }
    pipelines.value = (rawManifest.pipelines ?? []) as PipelineDef[]
  } catch (e: any) {
    toast.add({ title: e?.message ?? t('common.load_failed'), color: 'error' })
  } finally {
    loading.value = false
  }
}

async function save() {
  saving.value = true
  try {
    const updated: PluginManifest = { ...rawManifest, pipelines: pipelines.value }
    await pluginApi.updateManifest(pluginId, JSON.stringify(updated))
    rawManifest = updated
    toast.add({ title: t('admin.plugins.pipeline_saved'), color: 'success' })
  } catch (e: any) {
    toast.add({ title: e?.message ?? t('common.save_failed'), color: 'error' })
  } finally {
    saving.value = false
  }
}

onMounted(load)
</script>

<template>
  <AdminPageContainer>
    <AdminPageHeader
      :title="t('admin.plugins.pipeline_edit_title')"
      :subtitle="pluginId"
    >
      <template #actions>
        <UButton
          variant="outline"
          icon="i-tabler-arrow-left"
          to="/admin/plugins"
        >
          {{ t('common.back') }}
        </UButton>
        <UButton
          icon="i-tabler-device-floppy"
          :loading="saving"
          @click="save"
        >
          {{ t('common.save') }}
        </UButton>
      </template>
    </AdminPageHeader>

    <AdminPageContent>
      <div v-if="loading" class="flex justify-center py-16">
        <UIcon name="i-tabler-loader-2" class="animate-spin text-3xl text-gray-400" />
      </div>

      <template v-else>
        <div class="max-w-3xl space-y-6">
          <div>
            <div class="text-sm font-semibold text-gray-700 dark:text-gray-300 mb-1">
              {{ t('admin.plugins.pipeline_section_title') }}
            </div>
            <p class="text-xs text-gray-500 mb-4">
              {{ t('admin.plugins.pipeline_section_desc') }}
            </p>
            <PluginPipelineEditor v-model="pipelines" />
          </div>
        </div>
      </template>
    </AdminPageContent>
  </AdminPageContainer>
</template>

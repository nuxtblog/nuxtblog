<script setup lang="ts">
import type { ExtensionGroup, FormatPolicy } from '~/types/api/media'

const props = defineProps<{
  open: boolean
  /** null = new policy */
  policy: FormatPolicy | null
  extensionGroups: ExtensionGroup[]
}>()

const emit = defineEmits<{
  (e: 'update:open', val: boolean): void
  (e: 'saved'): void
}>()

const { t, locale } = useI18n()
const isZh = computed(() => locale.value.startsWith('zh'))
const mediaApi = useMediaApi()
const toast = useToast()

const isNew = computed(() => !props.policy)
const originalName = ref('')

const policyForm = ref({
  name: '',
  label_zh: '',
  label_en: '',
  is_system: false,
  groups: [] as string[],
})

watch(() => props.open, (val) => {
  if (!val) return
  if (props.policy) {
    originalName.value = props.policy.name
    policyForm.value = {
      name: props.policy.name,
      label_zh: props.policy.label_zh,
      label_en: props.policy.label_en,
      is_system: props.policy.is_system,
      groups: [...props.policy.groups],
    }
  } else {
    originalName.value = ''
    policyForm.value = { name: '', label_zh: '', label_en: '', is_system: false, groups: [] }
  }
})

const toggleGroup = (name: string) => {
  const idx = policyForm.value.groups.indexOf(name)
  if (idx >= 0) policyForm.value.groups.splice(idx, 1)
  else policyForm.value.groups.push(name)
}

const save = async () => {
  const f = policyForm.value
  if (!f.name.trim()) return
  try {
    if (isNew.value) {
      await mediaApi.createFormatPolicy({
        name: f.name.trim(),
        label_zh: f.label_zh.trim(),
        label_en: f.label_en.trim(),
        groups: f.groups,
      })
    } else {
      await mediaApi.updateFormatPolicy(originalName.value, {
        name: f.name.trim(),
        label_zh: f.label_zh.trim(),
        label_en: f.label_en.trim(),
        groups: f.groups,
      })
    }
    emit('saved')
    emit('update:open', false)
  } catch (err) {
    toast.add({ title: t('admin.settings.media.save_failed'), description: err instanceof Error ? err.message : undefined, color: 'error' })
  }
}
</script>

<template>
  <UModal :open="open" :title="isNew ? $t('admin.settings.media.add_policy') : $t('admin.settings.media.edit_policy')" @update:open="emit('update:open', $event)">
    <template #content>
      <div class="p-6 space-y-4">
        <UFormField :label="$t('admin.settings.media.group_name')">
          <UInput v-model="policyForm.name" :disabled="!isNew" class="w-full" />
        </UFormField>
        <div class="grid grid-cols-2 gap-3">
          <UFormField :label="$t('admin.settings.media.category_label_zh')">
            <UInput v-model="policyForm.label_zh" :disabled="policyForm.is_system" class="w-full" />
          </UFormField>
          <UFormField :label="$t('admin.settings.media.category_label_en')">
            <UInput v-model="policyForm.label_en" :disabled="policyForm.is_system" class="w-full" />
          </UFormField>
        </div>
        <UFormField :label="$t('admin.settings.media.policy_groups')">
          <div class="space-y-2">
            <label v-for="grp in extensionGroups" :key="grp.name" class="flex items-center gap-2 cursor-pointer">
              <input
                type="checkbox"
                :checked="policyForm.groups.includes(grp.name)"
                class="rounded border-default"
                @change="toggleGroup(grp.name)" />
              <span class="text-sm">{{ isZh ? grp.label_zh : grp.label_en }} <span class="text-muted font-mono text-xs">({{ grp.name }})</span></span>
            </label>
          </div>
        </UFormField>
        <div class="flex justify-end gap-2">
          <UButton color="neutral" variant="outline" @click="emit('update:open', false)">{{ $t('common.cancel') }}</UButton>
          <UButton color="primary" @click="save">{{ $t('common.confirm') }}</UButton>
        </div>
      </div>
    </template>
  </UModal>
</template>

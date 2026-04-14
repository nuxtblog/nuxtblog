<script setup lang="ts">
import type { ExtensionGroup } from '~/types/api/media'

const props = defineProps<{
  open: boolean
  /** Index in parent array; -1 = new group */
  editIndex: number
  /** The group to edit (when editIndex >= 0) */
  group: ExtensionGroup | null
}>()

const emit = defineEmits<{
  (e: 'update:open', val: boolean): void
  (e: 'save', group: ExtensionGroup, index: number): void
}>()

const { t } = useI18n()

const groupForm = ref({
  name: '',
  label_zh: '',
  label_en: '',
  extensionsStr: '',
  max_size_mb: 10,
})

watch(() => props.open, (val) => {
  if (!val) return
  if (props.editIndex >= 0 && props.group) {
    const g = props.group
    groupForm.value = {
      name: g.name,
      label_zh: g.label_zh,
      label_en: g.label_en,
      extensionsStr: g.extensions.join(', '),
      max_size_mb: g.max_size_mb,
    }
  } else {
    groupForm.value = { name: '', label_zh: '', label_en: '', extensionsStr: '', max_size_mb: 10 }
  }
})

const save = () => {
  const f = groupForm.value
  if (!f.name.trim()) return
  const extensions = f.extensionsStr.split(/[,\s]+/).map(s => s.trim().replace(/^\./, '').toLowerCase()).filter(Boolean)
  const group: ExtensionGroup = {
    name: f.name.trim(),
    label_zh: f.label_zh.trim(),
    label_en: f.label_en.trim(),
    extensions,
    max_size_mb: f.max_size_mb,
  }
  emit('save', group, props.editIndex)
  emit('update:open', false)
}
</script>

<template>
  <UModal :open="open" :title="editIndex < 0 ? $t('admin.settings.media.add_group') : $t('admin.settings.media.edit_group')" @update:open="emit('update:open', $event)">
    <template #content>
      <div class="p-6 space-y-4">
        <UFormField :label="$t('admin.settings.media.group_name')">
          <UInput v-model="groupForm.name" :disabled="editIndex >= 0" class="w-full" />
        </UFormField>
        <div class="grid grid-cols-2 gap-3">
          <UFormField :label="$t('admin.settings.media.category_label_zh')">
            <UInput v-model="groupForm.label_zh" class="w-full" />
          </UFormField>
          <UFormField :label="$t('admin.settings.media.category_label_en')">
            <UInput v-model="groupForm.label_en" class="w-full" />
          </UFormField>
        </div>
        <UFormField :label="$t('admin.settings.media.group_extensions')">
          <UInput v-model="groupForm.extensionsStr" placeholder="jpg,jpeg,png,webp,gif" class="w-full" />
          <template #hint>{{ $t('admin.settings.media.group_extensions_hint') }}</template>
        </UFormField>
        <UFormField :label="$t('admin.settings.media.group_max_size')">
          <div class="flex items-center gap-2">
            <UInput v-model.number="groupForm.max_size_mb" type="number" min="1" class="w-32" />
            <span class="text-sm text-muted">MB</span>
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

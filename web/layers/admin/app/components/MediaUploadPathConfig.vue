<script setup lang="ts">
const props = defineProps<{
  modelValue: string
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', val: string): void
}>()

const { t } = useI18n()

const PRESET_VALUES = [
  '{year}/{month}',
  '{year}/{month}/{day}',
  '{category}',
  '{category}/{year}/{month}',
  'files',
]

const uploadPathPresets = computed(() => [
  { value: '{year}/{month}', label: t('admin.settings.media.upload_path_preset_year_month') },
  { value: '{year}/{month}/{day}', label: t('admin.settings.media.upload_path_preset_year_month_day') },
  { value: '{category}', label: t('admin.settings.media.upload_path_preset_category') },
  { value: '{category}/{year}/{month}', label: t('admin.settings.media.upload_path_preset_category_year_month') },
  { value: 'files', label: t('admin.settings.media.upload_path_preset_flat') },
  { value: '__custom__', label: t('admin.settings.media.upload_path_preset_custom') },
])

const customPathActive = ref(false)

const uploadPathPreset = computed({
  get: () => {
    if (customPathActive.value) return '__custom__'
    return PRESET_VALUES.includes(props.modelValue) ? props.modelValue : '__custom__'
  },
  set: (val: string) => {
    if (val === '__custom__') {
      customPathActive.value = true
    } else {
      customPathActive.value = false
      emit('update:modelValue', val)
    }
  },
})

const uploadPathExample = computed(() => {
  const now = new Date()
  const yyyy = String(now.getFullYear())
  const mm = String(now.getMonth() + 1).padStart(2, '0')
  const dd = String(now.getDate()).padStart(2, '0')
  const cat = 'post_content'
  return (props.modelValue || '{year}/{month}')
    .replace('{year}', yyyy)
    .replace('{month}', mm)
    .replace('{day}', dd)
    .replace('{category}', cat)
    + '/abc123.jpg'
})

// Sync customPathActive when parent changes the value
watch(() => props.modelValue, (val) => {
  if (PRESET_VALUES.includes(val)) customPathActive.value = false
})
</script>

<template>
  <UCard>
    <template #header>
      <h3 class="text-lg font-semibold text-highlighted">{{ $t('admin.settings.media.upload_path') }}</h3>
    </template>
    <div class="space-y-4">
      <p class="text-xs text-muted">{{ $t('admin.settings.media.upload_path_hint') }}</p>
      <UFormField :label="$t('admin.settings.media.upload_path_label')">
        <div class="space-y-2">
          <USelect
            v-model="uploadPathPreset"
            :items="uploadPathPresets"
            class="w-80" />
          <template v-if="uploadPathPreset === '__custom__'">
            <UInput
              :model-value="modelValue"
              placeholder="{year}/{month}"
              class="w-80 font-mono text-sm"
              @update:model-value="emit('update:modelValue', String($event))" />
            <p class="text-xs text-muted">
              {{ $t('admin.settings.media.path_vars_hint') }}:
              <code class="font-mono text-highlighted">{year}</code>
              <code class="font-mono text-highlighted">{month}</code>
              <code class="font-mono text-highlighted">{day}</code>
              <code class="font-mono text-highlighted">{category}</code>
            </p>
          </template>
        </div>
      </UFormField>
      <div class="rounded-md bg-muted px-4 py-3 text-xs font-mono">
        <span class="text-muted">{{ $t('admin.settings.media.upload_path_preview') }}: </span>
        <span class="text-highlighted">{{ uploadPathExample }}</span>
      </div>
    </div>
  </UCard>
</template>

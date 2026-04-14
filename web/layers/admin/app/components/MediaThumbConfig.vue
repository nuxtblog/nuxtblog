<script setup lang="ts">
interface ThumbSize { width: number; height: number }

const props = defineProps<{
  thumbnailThumb: ThumbSize
  coverThumb: ThumbSize
  contentThumb: ThumbSize
}>()

const emit = defineEmits<{
  (e: 'update:thumbnailThumb', val: ThumbSize): void
  (e: 'update:coverThumb', val: ThumbSize): void
  (e: 'update:contentThumb', val: ThumbSize): void
}>()

const { t } = useI18n()

const thumbPresets = computed(() => [
  { key: 'thumbnailThumb' as const, variantKey: 'thumbnail', label: t('admin.settings.media.thumb_general'), desc: t('admin.settings.media.thumb_general_desc') },
  { key: 'coverThumb' as const, variantKey: 'cover', label: t('admin.settings.media.thumb_cover'), desc: t('admin.settings.media.thumb_cover_desc') },
  { key: 'contentThumb' as const, variantKey: 'content', label: t('admin.settings.media.thumb_content'), desc: t('admin.settings.media.thumb_content_desc') },
])

const values = computed(() => ({
  thumbnailThumb: props.thumbnailThumb,
  coverThumb: props.coverThumb,
  contentThumb: props.contentThumb,
}))

const thumbMode = (s: ThumbSize) => {
  if (s.width > 0 && s.height > 0) return t('admin.settings.media.mode_crop')
  if (s.width > 0) return t('admin.settings.media.mode_width')
  if (s.height > 0) return t('admin.settings.media.mode_height')
  return t('admin.settings.media.mode_none')
}
const thumbBadgeColor = (s: ThumbSize) =>
  (s.width === 0 && s.height === 0) ? 'neutral' : 'primary'

const emitMap = {
  thumbnailThumb: 'update:thumbnailThumb',
  coverThumb: 'update:coverThumb',
  contentThumb: 'update:contentThumb',
} as const

const updateWidth = (key: keyof typeof emitMap, val: number) => {
  emit(emitMap[key], { ...values.value[key], width: val })
}
const updateHeight = (key: keyof typeof emitMap, val: number) => {
  emit(emitMap[key], { ...values.value[key], height: val })
}
</script>

<template>
  <UCard>
    <template #header>
      <h3 class="text-lg font-semibold text-highlighted">{{ $t('admin.settings.media.thumb_gen') }}</h3>
    </template>
    <div class="space-y-6">
      <UAlert
        icon="i-tabler-info-circle"
        color="info"
        variant="soft"
        :title="$t('admin.settings.media.thumb_info_title')"
        :description="$t('admin.settings.media.thumb_info_desc')"
      />
      <div v-for="preset in thumbPresets" :key="preset.key" class="space-y-3">
        <div>
          <div class="text-sm font-medium text-highlighted">{{ preset.label }}</div>
          <div class="text-xs text-muted mt-0.5">{{ preset.desc }}</div>
        </div>
        <div class="flex items-center gap-3 flex-wrap">
          <UFormField :label="$t('admin.settings.media.width_label')">
            <UInput :model-value="values[preset.key].width" type="number" min="0" class="w-28" @update:model-value="updateWidth(preset.key, Number($event))" />
          </UFormField>
          <span class="text-muted mt-5">×</span>
          <UFormField :label="$t('admin.settings.media.height_label')">
            <UInput :model-value="values[preset.key].height" type="number" min="0" class="w-28" @update:model-value="updateHeight(preset.key, Number($event))" />
          </UFormField>
          <UBadge
            class="mt-5"
            :color="thumbBadgeColor(values[preset.key])"
            variant="soft"
            :label="thumbMode(values[preset.key])" />
        </div>
        <USeparator v-if="preset.key !== 'contentThumb'" />
      </div>

      <!-- variants preview -->
      <div class="rounded-md bg-muted p-4 text-xs space-y-1.5 font-mono">
        <p class="text-highlighted font-semibold mb-2 font-sans text-xs">{{ $t('admin.settings.media.variants_preview') }}</p>
        <template v-for="preset in thumbPresets" :key="preset.key">
          <p v-if="values[preset.key].width > 0 || values[preset.key].height > 0" class="text-muted">
            <span class="text-highlighted">"{{ preset.variantKey }}"</span>: ".../xxx_{{ preset.variantKey }}.jpg"
            <span class="ml-2 opacity-60">
              {{ values[preset.key].width || 'auto' }}×{{ values[preset.key].height || 'auto' }}px
            </span>
          </p>
        </template>
      </div>
    </div>
  </UCard>
</template>

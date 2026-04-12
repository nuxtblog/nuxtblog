<template>
  <AdminPageContainer>
    <AdminPageHeader :title="$t('admin.appearance.themes.title')" :subtitle="$t('admin.appearance.themes.subtitle')">
      <template #actions>
        <UButton color="neutral" variant="outline" :disabled="saving" @click="resetToDefault">
          {{ $t('admin.appearance.themes.reset_default') }}
        </UButton>
        <UButton color="primary" icon="i-tabler-device-floppy" :loading="saving" @click="save">
          {{ $t('common.save_changes') }}
        </UButton>
      </template>
    </AdminPageHeader>

    <AdminPageContent>
      <template v-if="loading">
        <div class="space-y-4">
          <UCard v-for="i in 4" :key="i">
            <template #header><USkeleton class="h-5 w-32" /></template>
            <div class="space-y-3">
              <USkeleton class="h-10 w-full" />
              <USkeleton class="h-10 w-3/4" />
            </div>
          </UCard>
        </div>
      </template>

      <template v-else>
        <!-- ── 颜色模式 ───────────────────────────────────────────────────── -->
        <UCard>
          <template #header>
            <h3 class="text-base font-semibold text-highlighted">{{ $t('admin.appearance.themes.color_mode') }}</h3>
          </template>

          <div class="grid grid-cols-3 gap-3">
            <button
              v-for="mode in COLOR_MODES"
              :key="mode.value"
              class="p-4 border-2 rounded-md transition-all cursor-pointer flex flex-col items-center gap-2 hover:border-primary/50"
              :class="form.colorMode === mode.value ? 'border-primary bg-primary/5' : 'border-default'"
              @click="setColorMode(mode.value)">
              <UIcon :name="mode.icon" class="size-7 text-highlighted" />
              <span class="text-sm font-medium text-highlighted">{{ mode.label }}</span>
            </button>
          </div>
        </UCard>

        <!-- ── 主色调 ─────────────────────────────────────────────────────── -->
        <UCard>
          <template #header>
            <div class="flex items-center justify-between">
              <h3 class="text-base font-semibold text-highlighted">{{ $t('admin.appearance.themes.primary_color') }}</h3>
              <div class="flex items-center gap-2 text-sm text-muted">
                <span
                  class="size-4 rounded-full inline-block ring-1 ring-default"
                  :style="{ backgroundColor: primaryHex }" />
                {{ primaryLabel }}
              </div>
            </div>
          </template>

          <div class="flex flex-wrap gap-2">
            <button
              v-for="color in PRIMARY_COLORS"
              :key="color.name"
              class="size-8 rounded-md border-2 transition-all relative shrink-0 hover:scale-110"
              :class="form.primary === color.name
                ? 'border-highlighted ring-2 ring-offset-2 ring-offset-default'
                : 'border-transparent'"
              :style="{ backgroundColor: color.hex, '--tw-ring-color': color.hex }"
              :title="color.label"
              @click="setPrimary(color.name)">
              <UIcon
                v-if="form.primary === color.name"
                name="i-tabler-check"
                class="size-4 text-white absolute inset-0 m-auto drop-shadow" />
            </button>
          </div>

          <!-- Live preview strip -->
          <div class="mt-4 flex items-center gap-3 p-3 bg-elevated rounded-md">
            <UButton color="primary" size="sm">{{ $t('admin.appearance.themes.preview_primary_btn') }}</UButton>
            <UButton color="primary" variant="soft" size="sm">{{ $t('admin.appearance.themes.preview_soft') }}</UButton>
            <UButton color="primary" variant="outline" size="sm">{{ $t('admin.appearance.themes.preview_outline') }}</UButton>
            <UBadge color="primary" :label="$t('admin.appearance.themes.preview_badge')" variant="soft" />
            <UProgress :value="60" class="flex-1" />
          </div>
        </UCard>

        <!-- ── 中性色 ─────────────────────────────────────────────────────── -->
        <UCard>
          <template #header>
            <div class="flex items-center justify-between">
              <h3 class="text-base font-semibold text-highlighted">{{ $t('admin.appearance.themes.neutral_color') }}</h3>
              <div class="flex items-center gap-2 text-sm text-muted">
                <span
                  class="size-4 rounded-full inline-block ring-1 ring-default"
                  :style="{ backgroundColor: neutralHex }" />
                {{ neutralLabel }}
              </div>
            </div>
          </template>

          <div class="flex flex-wrap gap-2">
            <button
              v-for="color in NEUTRAL_COLORS"
              :key="color.name"
              class="size-8 rounded-md border-2 transition-all relative shrink-0 hover:scale-110"
              :class="form.neutral === color.name
                ? 'border-highlighted ring-2 ring-offset-2 ring-offset-default'
                : 'border-transparent'"
              :style="{ backgroundColor: color.hex, '--tw-ring-color': color.hex }"
              :title="color.label"
              @click="setNeutral(color.name)">
              <UIcon
                v-if="form.neutral === color.name"
                name="i-tabler-check"
                class="size-4 text-white absolute inset-0 m-auto drop-shadow" />
            </button>
          </div>

          <p class="text-xs text-muted mt-3">{{ $t('admin.appearance.themes.neutral_desc') }}</p>
        </UCard>

        <!-- ── 圆角 ───────────────────────────────────────────────────────── -->
        <UCard>
          <template #header>
            <h3 class="text-base font-semibold text-highlighted">{{ $t('admin.appearance.themes.radius') }}</h3>
          </template>

          <div class="flex items-center gap-3 flex-wrap">
            <button
              v-for="r in RADIUS_OPTIONS"
              :key="r.value"
              class="flex flex-col items-center gap-2 p-3 border-2 rounded-md transition-all cursor-pointer min-w-16 hover:border-primary/50"
              :class="form.radius === r.value ? 'border-primary bg-primary/5' : 'border-default'"
              @click="setRadius(r.value)">
              <!-- Visual preview: a small square with the radius applied -->
              <div
                class="size-8 bg-primary/20 border-2 border-primary/40"
                :style="{ borderRadius: r.preview }" />
              <span class="text-xs font-medium text-highlighted">{{ r.label }}</span>
            </button>
          </div>

          <!-- Live preview -->
          <div class="mt-4 flex items-center gap-2 flex-wrap p-3 bg-elevated rounded-md">
            <UButton color="primary" size="sm">{{ $t('admin.appearance.themes.preview_btn') }}</UButton>
            <UInput :placeholder="$t('admin.appearance.themes.preview_input_placeholder')" class="w-40" />
            <UBadge color="primary" :label="$t('admin.appearance.themes.preview_tag')" variant="soft" />
            <UBadge color="neutral" :label="$t('admin.appearance.themes.preview_badge')" variant="outline" />
          </div>
        </UCard>

        <!-- ── 字体 ───────────────────────────────────────────────────────── -->
        <UCard>
          <template #header>
            <h3 class="text-base font-semibold text-highlighted">{{ $t('admin.appearance.themes.font') }}</h3>
          </template>

          <div class="grid grid-cols-2 md:grid-cols-4 gap-3 mb-4">
            <button
              v-for="font in FONT_OPTIONS"
              :key="font.value"
              class="p-3 border-2 rounded-md transition-all cursor-pointer text-left hover:border-primary/50"
              :class="form.font === font.value ? 'border-primary bg-primary/5' : 'border-default'"
              @click="setFont(font.value)">
              <div class="text-base font-medium text-highlighted mb-1" :style="{ fontFamily: font.stack }">
                Aa 博客
              </div>
              <div class="text-xs text-muted">{{ font.label }}</div>
            </button>
          </div>

          <UFormField :label="$t('admin.appearance.themes.font_size_label')" :hint="`${form.fontSize}px`">
            <input
              v-model.number="form.fontSize"
              type="range" min="13" max="20" step="1"
              class="w-full max-w-xs accent-primary"
              @input="applyFont(form.font, form.fontSize)" />
          </UFormField>
        </UCard>

        <!-- ── 布局 ───────────────────────────────────────────────────────── -->
        <UCard>
          <template #header>
            <h3 class="text-base font-semibold text-highlighted">{{ $t('admin.settings.general.layout') }}</h3>
          </template>

          <UFormField :label="$t('admin.settings.general.max_width')">
            <USelect
              v-model="form.containerWidth"
              :items="[
                { label: $t('admin.settings.general.width_narrow'), value: '5xl' },
                { label: $t('admin.settings.general.width_medium'), value: '6xl' },
                { label: $t('admin.settings.general.width_wide'), value: '7xl' },
              ]"
              class="w-full" />
            <p class="text-xs text-muted mt-1">{{ $t('admin.settings.general.max_width_hint') }}</p>
          </UFormField>
        </UCard>

      </template>
    </AdminPageContent>
  </AdminPageContainer>
</template>

<script setup lang="ts">
import {
  PRIMARY_COLORS,
  NEUTRAL_COLORS,
  RADIUS_OPTIONS,
  FONT_OPTIONS,
  DEFAULT_THEME,
  type ThemeSettings,
} from '~/config/theme'

const toast        = useToast()
const { t }        = useI18n()
const optionApi    = useOptionApi()
const optionsStore = useOptionsStore()
const { applyPrimary, applyNeutral, applyRadius, applyColorMode, applyFont } = useTheme()

const rawLoading = ref(true)
const loading    = useMinLoading(rawLoading)
const saving     = ref(false)

// ── Color mode options ─────────────────────────────────────────────────────────

const COLOR_MODES = computed(() => [
  { value: 'light',  label: t('admin.appearance.themes.mode_light'), icon: 'i-tabler-sun' },
  { value: 'dark',   label: t('admin.appearance.themes.mode_dark'),  icon: 'i-tabler-moon' },
  { value: 'system', label: t('admin.appearance.themes.mode_system'),  icon: 'i-tabler-device-laptop' },
])

// ── Form state ─────────────────────────────────────────────────────────────────

type ThemeFormState = Omit<ThemeSettings, 'customCss'> & { containerWidth: string }
const form = ref<ThemeFormState>({
  primary:   DEFAULT_THEME.primary,
  neutral:   DEFAULT_THEME.neutral,
  radius:    DEFAULT_THEME.radius,
  colorMode: DEFAULT_THEME.colorMode,
  font:      DEFAULT_THEME.font,
  fontSize:  DEFAULT_THEME.fontSize,
  containerWidth: '7xl',
})

const primaryHex   = computed(() => PRIMARY_COLORS.find((c) => c.name === form.value.primary)?.hex ?? '#8b5cf6')
const primaryLabel = computed(() => PRIMARY_COLORS.find((c) => c.name === form.value.primary)?.label ?? form.value.primary)
const neutralHex   = computed(() => NEUTRAL_COLORS.find((c) => c.name === form.value.neutral)?.hex ?? '#71717a')
const neutralLabel = computed(() => NEUTRAL_COLORS.find((c) => c.name === form.value.neutral)?.label ?? form.value.neutral)

// ── Load settings ──────────────────────────────────────────────────────────────

onMounted(async () => {
  try {
    await optionsStore.load()
    form.value = {
      primary:   optionsStore.get('theme_primary',    DEFAULT_THEME.primary),
      neutral:   optionsStore.get('theme_neutral',    DEFAULT_THEME.neutral),
      radius:    optionsStore.get('theme_radius',     DEFAULT_THEME.radius),
      colorMode: optionsStore.get('theme_color_mode', DEFAULT_THEME.colorMode),
      font:      optionsStore.get('theme_font',       DEFAULT_THEME.font),
      fontSize:  Number(optionsStore.get('theme_font_size', String(DEFAULT_THEME.fontSize))) || DEFAULT_THEME.fontSize,
      containerWidth: optionsStore.get('site_container_width', '7xl'),
    }
  } finally {
    rawLoading.value = false
  }
})

// ── Apply in real-time ─────────────────────────────────────────────────────────

const setPrimary = (color: string) => {
  form.value.primary = color
  applyPrimary(color)
}

const setNeutral = (color: string) => {
  form.value.neutral = color
  applyNeutral(color)
}

const setRadius = (radius: string) => {
  form.value.radius = radius
  applyRadius(radius)
}

const setColorMode = (mode: string) => {
  form.value.colorMode = mode
  applyColorMode(mode)
}

const setFont = (fontValue: string) => {
  form.value.font = fontValue
  applyFont(fontValue, form.value.fontSize)
}

// ── Save ───────────────────────────────────────────────────────────────────────

const OPTION_KEYS: { key: string; field: keyof ThemeFormState }[] = [
  { key: 'theme_primary',        field: 'primary' },
  { key: 'theme_neutral',        field: 'neutral' },
  { key: 'theme_radius',         field: 'radius' },
  { key: 'theme_color_mode',     field: 'colorMode' },
  { key: 'theme_font',           field: 'font' },
  { key: 'theme_font_size',      field: 'fontSize' },
  { key: 'site_container_width', field: 'containerWidth' },
]

const save = async () => {
  saving.value = true
  try {
    await Promise.all(
      OPTION_KEYS.map(({ key, field }) =>
        optionApi.setOption(key, form.value[field]),
      ),
    )
    await optionsStore.reload()
    toast.add({ title: t('admin.appearance.themes.saved'), color: 'success' })
  } catch (error: any) {
    toast.add({ title: t('common.save_failed'), description: error?.message, color: 'error' })
  } finally {
    saving.value = false
  }
}

const resetToDefault = () => {
  form.value = {
    primary:   DEFAULT_THEME.primary,
    neutral:   DEFAULT_THEME.neutral,
    radius:    DEFAULT_THEME.radius,
    colorMode: DEFAULT_THEME.colorMode,
    font:      DEFAULT_THEME.font,
    fontSize:  DEFAULT_THEME.fontSize,
    containerWidth: '7xl',
  }
  setPrimary(DEFAULT_THEME.primary)
  setNeutral(DEFAULT_THEME.neutral)
  setRadius(DEFAULT_THEME.radius)
  setColorMode(DEFAULT_THEME.colorMode)
  setFont(DEFAULT_THEME.font)
  toast.add({ title: t('admin.appearance.themes.reset_notice'), color: 'warning' })
}
</script>

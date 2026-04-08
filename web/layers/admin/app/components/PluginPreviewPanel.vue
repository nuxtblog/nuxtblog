<script setup lang="ts">
import type { PluginPreviewInfo } from '~/composables/usePluginApi'

const { t } = useI18n()

const props = defineProps<{ info: PluginPreviewInfo }>()

const caps = computed(() => props.info.capabilities)
const hasAnyCap = computed(() => caps.value.http || caps.value.store || caps.value.events)

const settings = computed(() => props.info.settings ?? [])
const webhooks = computed(() => props.info.webhooks ?? [])
const pipelines = computed(() => props.info.pipelines ?? [])
const requiredCount = computed(() => settings.value.filter(f => f.required).length)

const httpDomains = computed(() => {
  const allow = caps.value.http?.allow ?? []
  return allow.length ? allow.join(', ') : t('admin.plugins.preview_http_any')
})
</script>

<template>
  <!-- Plugin header -->
  <div class="flex items-start gap-3 mb-4">
    <div class="h-12 w-12 rounded-md bg-elevated flex items-center justify-center shrink-0">
      <UIcon :name="info.icon || 'i-tabler-plug'" class="size-6 text-primary" />
    </div>
    <div class="flex-1 min-w-0">
      <div class="flex items-center gap-2 flex-wrap">
        <span class="text-sm font-semibold text-highlighted">{{ info.title || info.name }}</span>
        <UBadge v-if="info.version" :label="`v${info.version}`" color="neutral" variant="soft" size="sm" />
      </div>
      <p class="text-xs text-muted mt-0.5">{{ info.author }}</p>
      <p v-if="info.description" class="text-xs text-muted mt-1 line-clamp-2">{{ info.description }}</p>
    </div>
  </div>

  <!-- Capabilities -->
  <div class="rounded-md border border-default overflow-hidden mb-3">
    <div class="px-3 py-2 bg-elevated border-b border-default">
      <span class="text-xs font-semibold text-muted uppercase tracking-wide">{{ $t('admin.plugins.cap_section') }}</span>
    </div>
    <div class="divide-y divide-default">
      <!-- No capabilities -->
      <div v-if="!hasAnyCap" class="flex items-center gap-2.5 px-3 py-2.5">
        <UIcon name="i-tabler-shield-check" class="size-4 text-success shrink-0" />
        <span class="text-xs text-muted">{{ $t('admin.plugins.preview_no_caps') }}</span>
      </div>

      <!-- HTTP -->
      <div v-if="caps.http" class="flex items-start gap-2.5 px-3 py-2.5">
        <UIcon name="i-tabler-world" class="size-4 text-info shrink-0 mt-0.5" />
        <div>
          <p class="text-xs font-medium text-highlighted">{{ $t('admin.plugins.cap_http') }}</p>
          <p class="text-xs text-muted mt-0.5">{{ httpDomains }}</p>
          <p v-if="caps.http.timeout_ms" class="text-xs text-muted">
            {{ $t('admin.plugins.preview_http_timeout', { ms: caps.http.timeout_ms }) }}
          </p>
        </div>
      </div>

      <!-- Store -->
      <div v-if="caps.store" class="flex items-start gap-2.5 px-3 py-2.5">
        <UIcon name="i-tabler-database" class="size-4 text-primary shrink-0 mt-0.5" />
        <div>
          <p class="text-xs font-medium text-highlighted">
            {{ caps.store.write ? $t('admin.plugins.cap_store_rw') : $t('admin.plugins.cap_store_r') }}
          </p>
          <p class="text-xs text-muted mt-0.5">{{ $t('admin.plugins.preview_store_desc') }}</p>
        </div>
      </div>

      <!-- Events -->
      <div v-if="caps.events" class="flex items-start gap-2.5 px-3 py-2.5">
        <UIcon name="i-tabler-bell" class="size-4 text-success shrink-0 mt-0.5" />
        <div>
          <p class="text-xs font-medium text-highlighted">{{ $t('admin.plugins.cap_events') }}</p>
          <p class="text-xs text-muted mt-0.5 font-mono">
            {{ (caps.events.subscribe ?? []).join(', ') || '*' }}
          </p>
        </div>
      </div>

      <!-- CSS injection -->
      <div v-if="info.has_css" class="flex items-center gap-2.5 px-3 py-2.5">
        <UIcon name="i-tabler-brush" class="size-4 text-warning shrink-0" />
        <span class="text-xs text-muted">{{ $t('admin.plugins.preview_has_css') }}</span>
      </div>
    </div>
  </div>

  <!-- Summary row: settings / webhooks / pipelines -->
  <div class="grid grid-cols-3 gap-2 text-center">
    <div class="rounded-md bg-elevated border border-default px-2 py-2">
      <p class="text-base font-bold text-highlighted">{{ settings.length }}</p>
      <p class="text-xs text-muted leading-tight mt-0.5">{{ $t('admin.plugins.preview_settings_label') }}</p>
      <p v-if="requiredCount > 0" class="text-xs text-warning mt-0.5">
        {{ $t('admin.plugins.preview_settings_required', { n: requiredCount }) }}
      </p>
    </div>
    <div class="rounded-md bg-elevated border border-default px-2 py-2">
      <p class="text-base font-bold text-highlighted">{{ webhooks.length }}</p>
      <p class="text-xs text-muted leading-tight mt-0.5">{{ $t('admin.plugins.preview_webhooks_label') }}</p>
    </div>
    <div class="rounded-md bg-elevated border border-default px-2 py-2">
      <p class="text-base font-bold text-highlighted">{{ pipelines.length }}</p>
      <p class="text-xs text-muted leading-tight mt-0.5">{{ $t('admin.plugins.preview_pipelines_label') }}</p>
    </div>
  </div>

  <!-- Pipeline detail -->
  <div v-if="pipelines.length" class="mt-2 space-y-1">
    <div
      v-for="p in pipelines" :key="p.name"
      class="flex items-center gap-2 px-2.5 py-1.5 rounded-md bg-elevated border border-default text-xs">
      <UIcon name="i-tabler-git-branch" class="size-3.5 text-muted shrink-0" />
      <span class="font-mono text-highlighted font-medium">{{ p.name }}</span>
      <span class="text-muted">← {{ p.trigger }}</span>
      <span class="ml-auto text-muted">{{ $t('admin.plugins.preview_pipeline_steps', { n: p.step_count }) }}</span>
    </div>
  </div>

  <!-- Priority note -->
  <p v-if="info.priority && info.priority !== 10" class="text-xs text-muted mt-2">
    {{ $t('admin.plugins.preview_priority', { n: info.priority }) }}
  </p>
</template>

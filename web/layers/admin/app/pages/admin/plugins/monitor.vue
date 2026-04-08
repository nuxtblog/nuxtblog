<script setup lang="ts">
import type { PluginItem, PluginStatsRes, PluginErrorEntry, PluginWindowBucket } from '~/composables/usePluginApi'

const { t } = useI18n()
const toast = useToast()
const pluginApi = usePluginApi()

// ── State ───────────────────────────────────────────────────────────────────
const plugins = ref<PluginItem[]>([])
const statsMap = ref<Record<string, PluginStatsRes>>({})
const rawLoading = ref(true)
const loading = useMinLoading(rawLoading)
const lastRefreshed = ref<Date | null>(null)

const expandedId = ref<string | null>(null)
const detailErrors = ref<PluginErrorEntry[]>([])
const detailLoading = ref(false)
const expandedDiff = ref<Set<string>>(new Set())

// ── Load all plugin stats ────────────────────────────────────────────────────
const load = async () => {
  rawLoading.value = true
  try {
    const res = await pluginApi.list()
    plugins.value = res.items ?? []
    // Fetch stats for all loaded plugins in parallel (best-effort; missing = not loaded)
    const settled = await Promise.allSettled(
      plugins.value.map(p => pluginApi.getStats(p.id).then(s => ({ id: p.id, s })))
    )
    const map: Record<string, PluginStatsRes> = {}
    for (const r of settled) {
      if (r.status === 'fulfilled') map[r.value.id] = r.value.s
    }
    statsMap.value = map
    lastRefreshed.value = new Date()
  } catch (e: any) {
    toast.add({ title: e?.message, color: 'error' })
  } finally {
    rawLoading.value = false
  }
}

onMounted(load)

// ── Auto-poll every 30 s ─────────────────────────────────────────────────────
const { pause: pausePoll } = useIntervalFn(load, 30_000)
onUnmounted(pausePoll)

// ── Expand plugin detail ─────────────────────────────────────────────────────
const toggleExpand = async (id: string) => {
  if (expandedId.value === id) {
    expandedId.value = null
    return
  }
  expandedId.value = id
  expandedDiff.value = new Set()
  detailLoading.value = true
  try {
    const res = await pluginApi.getErrors(id)
    detailErrors.value = (res.items ?? []).slice().reverse() // newest first
  } catch (e: any) {
    toast.add({ title: e?.message, color: 'error' })
  } finally {
    detailLoading.value = false
  }
}

// ── Helpers ──────────────────────────────────────────────────────────────────
const errorRate = (s: PluginStatsRes): string => {
  if (s.invocations === 0) return '0%'
  return ((s.errors / s.invocations) * 100).toFixed(1) + '%'
}

const formatRelative = (iso?: string): string => {
  if (!iso) return '—'
  const d = new Date(iso)
  const diffMs = Date.now() - d.getTime()
  const diffMin = Math.floor(diffMs / 60000)
  if (diffMin < 1) return t('admin.plugins.monitor.just_now')
  if (diffMin < 60) return t('admin.plugins.monitor.n_min_ago', { n: diffMin })
  const diffH = Math.floor(diffMin / 60)
  if (diffH < 24) return t('admin.plugins.monitor.n_hour_ago', { n: diffH })
  return d.toLocaleDateString()
}

const lastRefreshedText = computed(() => {
  if (!lastRefreshed.value) return ''
  return new Intl.DateTimeFormat(undefined, { timeStyle: 'medium' }).format(lastRefreshed.value)
})

// ── SVG sparkline ────────────────────────────────────────────────────────────
const CHART_W = 160
const CHART_H = 36

function sparklinePath(buckets: PluginWindowBucket[], field: 'invocations' | 'errors'): string {
  if (!buckets.length) return ''
  const vals = buckets.map(b => b[field])
  const max = Math.max(...vals, 1)
  const pts = vals.map((v, i) => {
    const x = (i / (vals.length - 1)) * CHART_W
    const y = CHART_H - (v / max) * (CHART_H - 4)
    return `${x.toFixed(1)},${y.toFixed(1)}`
  })
  return 'M' + pts.join('L')
}

// ── JSON diff formatter ───────────────────────────────────────────────────────
const formatDiff = (raw?: string): string => {
  if (!raw) return ''
  try {
    return JSON.stringify(JSON.parse(raw), null, 2)
  } catch {
    return raw
  }
}

const toggleDiff = (key: string) => {
  const next = new Set(expandedDiff.value)
  if (next.has(key)) next.delete(key)
  else next.add(key)
  expandedDiff.value = next
}
</script>

<template>
  <AdminPageContainer>
    <AdminPageHeader
      :title="$t('admin.plugins.monitor.title')"
      :subtitle="$t('admin.plugins.monitor.subtitle')">
      <template #actions>
        <div class="flex items-center gap-3">
          <span v-if="lastRefreshedText" class="text-xs text-muted hidden sm:inline">
            {{ $t('admin.plugins.monitor.last_refresh', { time: lastRefreshedText }) }}
          </span>
          <UButton
            icon="i-tabler-refresh"
            color="neutral"
            variant="outline"
            size="sm"
            :loading="loading"
            @click="load">
            {{ $t('admin.plugins.monitor.refresh') }}
          </UButton>
        </div>
      </template>
    </AdminPageHeader>

    <AdminPageContent>
      <!-- Loading -->
      <div v-if="loading" class="space-y-3">
        <div v-for="i in 4" :key="i" class="flex items-center gap-4 p-4 border border-default rounded-lg">
          <USkeleton class="h-10 w-10 rounded-md shrink-0" />
          <div class="flex-1 space-y-2">
            <USkeleton class="h-4 w-40" />
            <USkeleton class="h-3 w-24" />
          </div>
          <USkeleton class="h-9 w-40 rounded shrink-0" />
          <USkeleton class="h-5 w-16 rounded-full shrink-0" />
          <USkeleton class="h-5 w-20 rounded-full shrink-0" />
          <USkeleton class="h-5 w-24 shrink-0" />
        </div>
      </div>

      <!-- Empty -->
      <div v-else-if="plugins.length === 0" class="flex flex-col items-center justify-center py-16">
        <UIcon name="i-tabler-plug-x" class="size-16 text-muted mb-4" />
        <h3 class="text-lg font-medium text-highlighted mb-1">{{ $t('admin.plugins.no_plugins') }}</h3>
        <p class="text-sm text-muted">{{ $t('admin.plugins.no_plugins_desc') }}</p>
      </div>

      <!-- Plugin rows -->
      <div v-else class="space-y-2">
        <div v-for="item in plugins" :key="item.id">
          <!-- Row -->
          <div
            class="group flex items-center gap-4 p-4 bg-default border border-default rounded-lg cursor-pointer hover:shadow-sm transition-all"
            :class="expandedId === item.id ? 'border-primary/40 bg-primary/5' : ''"
            @click="toggleExpand(item.id)">

            <!-- Icon -->
            <div class="h-10 w-10 rounded-md bg-elevated flex items-center justify-center shrink-0">
              <UIcon :name="item.icon || 'i-tabler-plug'" class="size-5 text-primary" />
            </div>

            <!-- Name & author -->
            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-2">
                <span class="text-sm font-semibold text-highlighted">{{ item.title || item.id }}</span>
                <UBadge
                  :label="item.enabled ? $t('admin.plugins.status_active') : $t('admin.plugins.status_inactive')"
                  :color="item.enabled ? 'success' : 'neutral'"
                  variant="soft" size="xs" />
              </div>
              <p class="text-xs text-muted">{{ item.id }}</p>
            </div>

            <!-- Mini sparkline (only when stats available) -->
            <div v-if="statsMap[item.id]?.history?.length" class="shrink-0 hidden md:block">
              <svg :width="CHART_W" :height="CHART_H" class="overflow-visible">
                <!-- Invocations line (blue) -->
                <path
                  :d="sparklinePath(statsMap[item.id].history, 'invocations')"
                  fill="none"
                  stroke="var(--color-primary-400)"
                  stroke-width="1.5"
                  stroke-linejoin="round"
                  stroke-linecap="round" />
                <!-- Errors line (red) -->
                <path
                  v-if="statsMap[item.id].errors > 0"
                  :d="sparklinePath(statsMap[item.id].history, 'errors')"
                  fill="none"
                  stroke="var(--color-error-400)"
                  stroke-width="1.5"
                  stroke-linejoin="round"
                  stroke-linecap="round" />
              </svg>
            </div>
            <div v-else class="shrink-0 hidden md:block w-40" />

            <!-- Stats badges -->
            <div class="flex items-center gap-3 shrink-0 text-right">
              <div v-if="statsMap[item.id]" class="text-center px-3 py-1.5 bg-default rounded-md border border-default min-w-16">
                <div class="text-sm font-semibold text-highlighted">{{ statsMap[item.id].invocations }}</div>
                <div class="text-xs text-muted">{{ $t('admin.plugins.monitor.calls') }}</div>
              </div>
              <div
                v-if="statsMap[item.id]"
                class="text-center px-3 py-1.5 rounded-md border min-w-14"
                :class="statsMap[item.id].errors > 0 ? 'bg-error/10 border-error/20' : 'bg-default border-default'">
                <div
                  class="text-sm font-semibold"
                  :class="statsMap[item.id].errors > 0 ? 'text-error' : 'text-highlighted'">
                  {{ errorRate(statsMap[item.id]) }}
                </div>
                <div class="text-xs text-muted">{{ $t('admin.plugins.monitor.error_rate') }}</div>
              </div>
              <div v-if="statsMap[item.id]" class="text-center px-3 py-1.5 bg-default rounded-md border border-default min-w-24">
                <div class="text-xs font-medium text-highlighted truncate max-w-24" :title="statsMap[item.id].last_error_at">
                  {{ formatRelative(statsMap[item.id].last_error_at) }}
                </div>
                <div class="text-xs text-muted">{{ $t('admin.plugins.monitor.last_error') }}</div>
              </div>
              <div v-if="!statsMap[item.id]" class="text-xs text-muted italic">
                {{ $t('admin.plugins.monitor.not_loaded') }}
              </div>
              <UIcon
                :name="expandedId === item.id ? 'i-tabler-chevron-up' : 'i-tabler-chevron-down'"
                class="size-4 text-muted transition-transform" />
            </div>
          </div>

          <!-- Expanded detail panel -->
          <div v-if="expandedId === item.id" class="border border-t-0 border-primary/20 rounded-b-lg bg-default px-4 pb-4">

            <!-- Detail loading -->
            <div v-if="detailLoading" class="pt-4 space-y-2">
              <USkeleton class="h-4 w-48" />
              <USkeleton class="h-3 w-full" />
              <USkeleton class="h-3 w-3/4" />
            </div>

            <template v-else>
              <!-- 4.5-F2: Full-width line chart -->
              <div v-if="statsMap[item.id]?.history?.length" class="pt-4 pb-2">
                <p class="text-xs text-muted mb-2">{{ $t('admin.plugins.monitor.chart_title') }}</p>
                <div class="relative h-20 w-full">
                  <svg class="w-full h-full overflow-visible" preserveAspectRatio="none" viewBox="0 0 160 36">
                    <!-- Invocations -->
                    <path
                      :d="sparklinePath(statsMap[item.id].history, 'invocations')"
                      fill="none"
                      stroke="var(--color-primary-400)"
                      stroke-width="1.5"
                      stroke-linejoin="round"
                      stroke-linecap="round" />
                    <!-- Errors -->
                    <path
                      :d="sparklinePath(statsMap[item.id].history, 'errors')"
                      fill="none"
                      stroke="var(--color-error-400)"
                      stroke-width="1.5"
                      stroke-linejoin="round"
                      stroke-linecap="round" />
                  </svg>
                </div>
                <div class="flex items-center gap-4 mt-1">
                  <span class="flex items-center gap-1 text-xs text-muted">
                    <span class="inline-block w-3 h-0.5 bg-primary-400 rounded" />
                    {{ $t('admin.plugins.monitor.legend_calls') }}
                  </span>
                  <span class="flex items-center gap-1 text-xs text-muted">
                    <span class="inline-block w-3 h-0.5 bg-error-400 rounded" />
                    {{ $t('admin.plugins.monitor.legend_errors') }}
                  </span>
                </div>
              </div>

              <!-- 4.5-F3: Error list -->
              <div class="pt-3">
                <p class="text-xs font-medium text-highlighted mb-2">
                  {{ $t('admin.plugins.monitor.recent_errors') }}
                  <span class="text-muted font-normal">({{ detailErrors.length }})</span>
                </p>

                <div v-if="detailErrors.length === 0" class="text-xs text-muted py-3 text-center">
                  {{ $t('admin.plugins.monitor.no_errors') }}
                </div>

                <div v-else class="rounded-md border border-default overflow-hidden divide-y divide-default">
                  <div
                    v-for="(entry, idx) in detailErrors"
                    :key="idx"
                    class="px-3 py-2 text-xs">
                    <div class="flex items-start justify-between gap-4">
                      <div class="flex-1 min-w-0">
                        <div class="flex items-center gap-2 mb-0.5">
                          <span class="text-muted shrink-0">{{ new Date(entry.at).toLocaleString() }}</span>
                          <UBadge :label="entry.event" color="neutral" variant="soft" size="xs" />
                        </div>
                        <p class="text-error font-medium truncate" :title="entry.message">{{ entry.message }}</p>
                      </div>
                      <UButton
                        v-if="entry.input_diff"
                        icon="i-tabler-code"
                        color="neutral"
                        variant="ghost"
                        size="xs"
                        square
                        :class="expandedDiff.has(entry.at + idx) ? 'text-primary' : ''"
                        @click.stop="toggleDiff(entry.at + idx)" />
                    </div>
                    <!-- InputDiff -->
                    <div v-if="entry.input_diff && expandedDiff.has(entry.at + idx)" class="mt-2">
                      <pre class="text-xs bg-muted rounded p-2 overflow-x-auto text-highlighted leading-relaxed">{{ formatDiff(entry.input_diff) }}</pre>
                    </div>
                  </div>
                </div>
              </div>
            </template>
          </div>
        </div>
      </div>
    </AdminPageContent>
  </AdminPageContainer>
</template>

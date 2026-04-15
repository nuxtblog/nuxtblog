<template>
  <VueDraggable
    tag="div"
    class="space-y-2"
    :model-value="modelValue"
    :animation="200"
    handle=".widget-drag-handle"
    @update:model-value="$emit('update:modelValue', $event)"
  >
    <div
      v-for="widget in modelValue"
      :key="widget.id"
      class="rounded-md border border-default overflow-hidden"
    >
      <!-- 行 -->
      <div
        class="flex items-center gap-3 px-3 py-2.5"
        :class="expanded[widget.id] ? 'bg-elevated' : 'bg-default'"
      >
        <UIcon
          name="i-tabler-grip-vertical"
          class="widget-drag-handle size-4 shrink-0 cursor-grab text-muted hover:text-primary active:cursor-grabbing"
        />
        <UCheckbox v-model="widget.enabled" />
        <span
          class="flex-1 text-sm select-none"
          :class="widget.enabled ? 'text-highlighted' : 'text-muted'"
        >
          {{ widget.isPlugin ? widget.label : $t(widget.label) }}
          <UIcon v-if="widget.isPlugin" name="i-tabler-puzzle" class="size-3.5 text-muted ml-1 inline-block align-text-bottom" />
        </span>
        <div class="flex items-center gap-0.5">
          <UButton
            :icon="expanded[widget.id] ? 'i-tabler-chevron-up' : 'i-tabler-settings'"
            variant="ghost"
            :color="expanded[widget.id] ? 'primary' : 'neutral'"
            size="xs"
            square
            @click="handleToggle(widget)"
          />
        </div>
      </div>

      <!-- 展开的配置面板 -->
      <div
        v-if="expanded[widget.id]"
        class="px-4 py-3 bg-muted/40 border-t border-default space-y-3"
      >
        <div class="flex items-center gap-3">
          <span class="text-sm text-default flex-1">{{ $t('admin.widgets.widget_title') }}</span>
          <UInput v-model="widget.title" :placeholder="widget.isPlugin ? widget.label : $t(widget.label)" size="sm" class="w-40" />
        </div>

        <!-- Plugin widgets: custom settings only (no built-in maxCount) -->
        <template v-if="widget.isPlugin">
          <PluginSettingFields
            v-if="getWidgetSettings(widget.id)?.length && widget.pluginSettings"
            :schema="getWidgetSettings(widget.id)!"
            :model-value="widget.pluginSettings"
            @update:model-value="widget.pluginSettings = $event"
          />
          <p v-else class="text-xs text-muted">{{ $t('admin.widgets.no_settings') }}</p>
        </template>
        <!-- Built-in widgets: use configFields from registry -->
        <template v-else>
          <template v-if="getWidgetDef(widget.id)?.configFields.includes('showRecent')">
            <div class="flex items-center justify-between">
              <span class="text-sm text-default">{{ $t('admin.widgets.show_recent') }}</span>
              <UCheckbox v-model="widget.showRecent" />
            </div>
          </template>
          <template v-if="getWidgetDef(widget.id)?.configFields.includes('showHot')">
            <div class="flex items-center justify-between">
              <span class="text-sm text-default">{{ $t('admin.widgets.show_hot') }}</span>
              <UCheckbox v-model="widget.showHot" />
            </div>
          </template>
          <template v-if="getWidgetDef(widget.id)?.configFields.includes('maxCount')">
            <div class="flex items-center gap-3">
              <span class="text-sm text-default flex-1">{{ $t('admin.widgets.max_count') }}</span>
              <UInput v-model.number="widget.maxCount" type="number" :min="1" :max="20" size="sm" class="w-20" />
              <span class="text-xs text-muted">{{ $t('admin.widgets.items') }}</span>
            </div>
          </template>
        </template>
      </div>
    </div>
  </VueDraggable>
</template>

<script setup lang="ts">
import { VueDraggable } from 'vue-draggable-plus'
import type { WidgetConfig } from '~/composables/useWidgetConfig'
import type { PluginSettingField } from '~/composables/usePluginApi'
import { getWidgetDef } from '~/config/widgets'
import { usePluginContributionsStore } from '~/stores/plugin-contributions'

const props = defineProps<{
  modelValue: WidgetConfig[]
  expanded: Record<string, boolean>
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', v: WidgetConfig[]): void
  (e: 'toggle', id: string): void
}>()

const contributionsStore = usePluginContributionsStore()
const pluginWidgetViews = contributionsStore.getViewItems('public:sidebar-widget')

/** Look up widget settings schema from the plugin contributions store. */
function getWidgetSettings(widgetId: string): PluginSettingField[] | undefined {
  const parts = widgetId.split(':')
  if (parts.length < 3 || parts[0] !== 'plugin') return undefined
  const pluginId = parts[1]
  const viewId = parts.slice(2).join(':')
  const view = pluginWidgetViews.value.find(v => v.pluginId === pluginId && v.id === viewId)
  return view?.settings?.length ? view.settings : undefined
}

/** Initialize pluginSettings on a widget if missing. Call OUTSIDE render (e.g. on toggle). */
function initPluginSettings(widget: WidgetConfig) {
  if (widget.pluginSettings) return
  const schema = getWidgetSettings(widget.id)
  if (!schema) return
  const defaults: Record<string, unknown> = {}
  for (const field of schema) {
    if (field.default !== undefined) defaults[field.key] = field.default
  }
  widget.pluginSettings = defaults
}

/** Handle toggle — init pluginSettings before expanding (avoids render-phase mutation). */
function handleToggle(widget: WidgetConfig) {
  initPluginSettings(widget)
  emit('toggle', widget.id)
}

</script>

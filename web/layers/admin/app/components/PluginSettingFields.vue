<script setup lang="ts">
import type { PluginSettingField } from '~/composables/usePluginApi'
import { usePluginLocale } from '~/composables/usePluginLocale'

const props = defineProps<{
  schema: PluginSettingField[]
  modelValue: Record<string, unknown>
  pluginId?: string
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', val: Record<string, unknown>): void
}>()

const values = computed({
  get: () => props.modelValue,
  set: (v) => emit('update:modelValue', v),
})

const { tField } = usePluginLocale()

function updateField(key: string, value: unknown) {
  emit('update:modelValue', { ...props.modelValue, [key]: value })
}

/** Resolve i18n template for a field's text props. */
function fieldText(field: PluginSettingField) {
  if (!props.pluginId) return { label: field.label, description: field.description, group: field.group, placeholder: field.placeholder }
  return tField(props.pluginId, field)
}

/** Group settings fields by their `group` property, preserving order. */
const groups = computed(() => {
  const result: Array<{ name: string | undefined; fields: PluginSettingField[] }> = []
  const groupMap = new Map<string | undefined, PluginSettingField[]>()
  for (const field of props.schema) {
    const key = field.group || undefined
    if (!groupMap.has(key)) {
      const fields: PluginSettingField[] = []
      groupMap.set(key, fields)
      result.push({ name: key, fields })
    }
    groupMap.get(key)!.push(field)
  }
  return result
})

/** Evaluate a field's showIf condition against current values. */
function isFieldVisible(field: PluginSettingField): boolean {
  if (!field.showIf) return true
  try {
    const expr = field.showIf.trim()
    const match = expr.match(/^(\w+)\s*(===|!==|==|!=)\s*(.+)$/)
    if (!match) return true
    const [, key, op, rawVal] = match
    const actual = values.value[key]
    let expected: unknown = rawVal.trim()
    if (expected === 'true') expected = true
    else if (expected === 'false') expected = false
    else if (expected === 'null') expected = null
    else if (/^['"].*['"]$/.test(expected as string)) expected = (expected as string).slice(1, -1)
    else if (!isNaN(Number(expected))) expected = Number(expected)
    if (op === '===' || op === '==') return actual === expected
    if (op === '!==' || op === '!=') return actual !== expected
  } catch { /* show field on parse error */ }
  return true
}
</script>

<template>
  <template v-for="(group, groupIdx) in groups" :key="group.name ?? '__default__'">
    <div v-if="group.name" :class="{ 'pt-3 border-t border-default': groupIdx > 0 }">
      <p class="text-xs font-medium text-muted mb-3 uppercase tracking-wide">{{ fieldText(group.fields[0]).group || group.name }}</p>
    </div>
    <template v-for="field in group.fields" :key="field.key">
      <UFormField v-if="isFieldVisible(field)" :label="fieldText(field).label || field.key" :description="fieldText(field).description" :required="field.required">
        <USwitch v-if="field.type === 'boolean'" :model-value="!!values[field.key]" @update:model-value="updateField(field.key, $event)" />
        <USelect v-else-if="field.type === 'select'" :model-value="String(values[field.key] ?? '')" :items="(field.options ?? []).map(o => ({ label: o, value: o }))" class="w-full" @update:model-value="updateField(field.key, $event)" />
        <UTextarea v-else-if="field.type === 'textarea'" :model-value="String(values[field.key] ?? '')" :placeholder="fieldText(field).placeholder" class="w-full" :rows="3" @update:model-value="updateField(field.key, $event)" />
        <UInput v-else-if="field.type === 'password'" type="password" :model-value="String(values[field.key] ?? '')" :placeholder="fieldText(field).placeholder" class="w-full" @update:model-value="updateField(field.key, $event)" />
        <UInput v-else-if="field.type === 'number'" type="number" :model-value="String(values[field.key] ?? '')" :placeholder="fieldText(field).placeholder" class="w-full" @update:model-value="updateField(field.key, Number($event))" />
        <UInput v-else :model-value="String(values[field.key] ?? '')" :placeholder="fieldText(field).placeholder" class="w-full" @update:model-value="updateField(field.key, $event)" />
      </UFormField>
    </template>
  </template>
</template>

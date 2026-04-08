<script setup lang="ts">
import type { StepDef } from '~/composables/usePluginApi'

const props = defineProps<{
  modelValue: StepDef
  depth?: number
}>()
const emit = defineEmits<{ 'update:modelValue': [StepDef] }>()
const { t } = useI18n()

const typeOptions = [
  { label: 'JS 函数', value: 'js' },
  { label: 'Webhook', value: 'webhook' },
  { label: '条件分支', value: 'condition' },
]

function updateField<K extends keyof StepDef>(key: K, value: StepDef[K]) {
  emit('update:modelValue', { ...props.modelValue, [key]: value })
}

function addBranchStep(branch: 'then' | 'else') {
  const current = [...(props.modelValue[branch] ?? [])]
  current.push({ type: 'js', name: '' } as StepDef)
  updateField(branch, current)
}

function removeBranchStep(branch: 'then' | 'else', i: number) {
  const current = (props.modelValue[branch] ?? []).filter((_, idx) => idx !== i)
  updateField(branch, current)
}

function updateBranchStep(branch: 'then' | 'else', i: number, v: StepDef) {
  const current = [...(props.modelValue[branch] ?? [])]
  current[i] = v
  updateField(branch, current)
}
</script>

<template>
  <div class="border border-gray-200 dark:border-gray-700 rounded-lg p-3 space-y-3 bg-gray-50 dark:bg-gray-800/50">
    <!-- Name + Type row -->
    <div class="flex gap-2">
      <UFormField :label="t('admin.plugins.pipeline_step_name')" class="flex-1">
        <UInput
          :model-value="modelValue.name"
          :placeholder="t('admin.plugins.pipeline_step_name_placeholder')"
          @update:model-value="updateField('name', String($event))"
        />
      </UFormField>
      <UFormField :label="t('admin.plugins.pipeline_step_type')" class="w-36">
        <USelect
          :model-value="modelValue.type"
          :options="typeOptions"
          @update:model-value="updateField('type', $event as StepDef['type'])"
        />
      </UFormField>
    </div>

    <!-- JS fields -->
    <template v-if="modelValue.type === 'js'">
      <UFormField :label="t('admin.plugins.pipeline_step_fn')">
        <UInput
          :model-value="modelValue.fn ?? ''"
          :placeholder="t('admin.plugins.pipeline_step_fn_placeholder')"
          @update:model-value="updateField('fn', String($event))"
        />
      </UFormField>
    </template>

    <!-- Webhook fields -->
    <template v-else-if="modelValue.type === 'webhook'">
      <UFormField :label="t('admin.plugins.pipeline_step_url')">
        <UInput
          :model-value="modelValue.url ?? ''"
          placeholder="https://..."
          @update:model-value="updateField('url', String($event))"
        />
      </UFormField>
    </template>

    <!-- Condition fields -->
    <template v-else-if="modelValue.type === 'condition'">
      <UFormField :label="t('admin.plugins.pipeline_step_if')">
        <UInput
          :model-value="modelValue.if ?? ''"
          :placeholder="t('admin.plugins.pipeline_step_if_placeholder')"
          @update:model-value="updateField('if', String($event))"
        />
      </UFormField>

      <!-- Then branch -->
      <div v-if="(depth ?? 0) < 3" class="pl-3 border-l-2 border-green-400 space-y-2">
        <div class="text-xs font-medium text-green-600 dark:text-green-400">{{ t('admin.plugins.pipeline_step_then') }}</div>
        <PluginStepEditor
          v-for="(sub, i) in modelValue.then ?? []"
          :key="i"
          :model-value="sub"
          :depth="(depth ?? 0) + 1"
          @update:model-value="updateBranchStep('then', i, $event)"
        />
        <UButton size="xs" variant="ghost" icon="i-tabler-plus" @click="addBranchStep('then')">
          {{ t('admin.plugins.pipeline_add_step') }}
        </UButton>
      </div>

      <!-- Else branch -->
      <div v-if="(depth ?? 0) < 3" class="pl-3 border-l-2 border-orange-400 space-y-2">
        <div class="text-xs font-medium text-orange-600 dark:text-orange-400">{{ t('admin.plugins.pipeline_step_else') }}</div>
        <PluginStepEditor
          v-for="(sub, i) in modelValue.else ?? []"
          :key="i"
          :model-value="sub"
          :depth="(depth ?? 0) + 1"
          @update:model-value="updateBranchStep('else', i, $event)"
        />
        <UButton size="xs" variant="ghost" icon="i-tabler-plus" @click="addBranchStep('else')">
          {{ t('admin.plugins.pipeline_add_step') }}
        </UButton>
      </div>
    </template>

    <!-- Timeout + Retry -->
    <div class="flex gap-2">
      <UFormField :label="t('admin.plugins.pipeline_step_timeout')" class="flex-1">
        <UInput
          type="number"
          :model-value="modelValue.timeout_ms ?? 5000"
          @update:model-value="updateField('timeout_ms', Number($event))"
        />
      </UFormField>
      <UFormField :label="t('admin.plugins.pipeline_step_retry')" class="w-28">
        <UInput
          type="number"
          :model-value="modelValue.retry ?? 0"
          @update:model-value="updateField('retry', Number($event))"
        />
      </UFormField>
    </div>
  </div>
</template>

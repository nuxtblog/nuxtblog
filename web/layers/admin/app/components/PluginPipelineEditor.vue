<script setup lang="ts">
import type { PipelineDef, StepDef } from '~/composables/usePluginApi'

const props = defineProps<{
  modelValue: PipelineDef[]
}>()
const emit = defineEmits<{ 'update:modelValue': [PipelineDef[]] }>()
const { t } = useI18n()

function update(i: number, v: PipelineDef) {
  const updated = [...props.modelValue]
  updated[i] = v
  emit('update:modelValue', updated)
}

function remove(i: number) {
  emit('update:modelValue', props.modelValue.filter((_, idx) => idx !== i))
}

function add() {
  emit('update:modelValue', [
    ...props.modelValue,
    { name: '', trigger: '', steps: [] },
  ])
}

function addStep(pIdx: number) {
  const p = { ...props.modelValue[pIdx], steps: [...props.modelValue[pIdx].steps, { type: 'js', name: '' } as StepDef] }
  update(pIdx, p)
}

function removeStep(pIdx: number, sIdx: number) {
  const p = { ...props.modelValue[pIdx], steps: props.modelValue[pIdx].steps.filter((_, i) => i !== sIdx) }
  update(pIdx, p)
}

function updateStep(pIdx: number, sIdx: number, v: StepDef) {
  const steps = [...props.modelValue[pIdx].steps]
  steps[sIdx] = v
  update(pIdx, { ...props.modelValue[pIdx], steps })
}

function updateField(pIdx: number, key: keyof PipelineDef, value: unknown) {
  update(pIdx, { ...props.modelValue[pIdx], [key]: value })
}
</script>

<template>
  <div class="space-y-4">
    <div
      v-for="(pipeline, pIdx) in modelValue"
      :key="pIdx"
      class="border border-gray-200 dark:border-gray-700 rounded-xl p-4 space-y-3"
    >
      <!-- Pipeline header: name + trigger -->
      <div class="flex items-start gap-2">
        <UIcon name="i-tabler-route" class="text-primary-500 mt-2 shrink-0" />
        <div class="flex-1 flex gap-2">
          <UFormField :label="t('admin.plugins.pipeline_name_placeholder')" class="flex-1">
            <UInput
              :model-value="pipeline.name"
              :placeholder="t('admin.plugins.pipeline_name_placeholder')"
              @update:model-value="updateField(pIdx, 'name', String($event))"
            />
          </UFormField>
          <UFormField :label="t('admin.plugins.pipeline_trigger_placeholder')" class="flex-1">
            <UInput
              :model-value="pipeline.trigger"
              :placeholder="t('admin.plugins.pipeline_trigger_placeholder')"
              @update:model-value="updateField(pIdx, 'trigger', String($event))"
            />
          </UFormField>
        </div>
        <UButton
          size="sm"
          variant="ghost"
          color="error"
          icon="i-tabler-trash"
          class="mt-5"
          @click="remove(pIdx)"
        />
      </div>

      <!-- Steps list -->
      <div class="pl-6 border-l border-gray-200 dark:border-gray-700 space-y-2">
        <div class="text-xs text-gray-500 font-medium">{{ t('admin.plugins.pipeline_steps') }}</div>
        <div
          v-for="(step, sIdx) in pipeline.steps"
          :key="sIdx"
          class="relative"
        >
          <PluginStepEditor
            :model-value="step"
            @update:model-value="updateStep(pIdx, sIdx, $event)"
          />
          <UButton
            size="xs"
            variant="ghost"
            color="error"
            icon="i-tabler-x"
            class="absolute top-2 right-2"
            @click="removeStep(pIdx, sIdx)"
          />
        </div>
        <UButton
          size="sm"
          variant="outline"
          icon="i-tabler-plus"
          @click="addStep(pIdx)"
        >
          {{ t('admin.plugins.pipeline_add_step') }}
        </UButton>
      </div>
    </div>

    <!-- Add new pipeline -->
    <UButton
      variant="outline"
      icon="i-tabler-plus"
      @click="add"
    >
      {{ t('admin.plugins.pipeline_add') }}
    </UButton>
  </div>
</template>

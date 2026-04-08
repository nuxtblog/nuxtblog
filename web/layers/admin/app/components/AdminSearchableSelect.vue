<template>
  <UPopover v-model:open="open" :disabled="disabled">
    <UButton
      color="neutral"
      variant="outline"
      size="sm"
      trailing-icon="i-tabler-chevron-down"
      :disabled="disabled"
      :class="triggerClass">
      {{ selectedLabel || placeholder }}
    </UButton>
    <template #content>
      <div class="p-2 w-52">
        <UInput
          v-model="search"
          :placeholder="searchPlaceholder || t('common.search')"
          leading-icon="i-tabler-search"
          size="sm"
          class="mb-2">
          <template v-if="search" #trailing>
            <UButton icon="i-tabler-x" color="neutral" variant="ghost" size="xs" @click="search = ''" />
          </template>
        </UInput>
        <div class="max-h-48 overflow-y-auto space-y-0.5">
          <div v-if="filteredItems.length === 0" class="px-2 py-3 text-sm text-muted text-center">
            {{ t('common.no_results') }}
          </div>
          <button
            v-for="item in filteredItems"
            :key="String(item.value)"
            class="w-full text-left px-2 py-1.5 text-sm rounded hover:bg-elevated transition-colors"
            :class="{ 'text-primary font-medium': isSelected(item.value) }"
            @click="select(item.value)">
            {{ item.label }}
          </button>
        </div>
      </div>
    </template>
  </UPopover>
</template>

<script setup lang="ts">
const { t } = useI18n()

const props = defineProps<{
  modelValue?: any
  items: Array<{ label: string; value: any }>
  placeholder?: string
  searchPlaceholder?: string
  triggerClass?: string
  disabled?: boolean
}>()

const emit = defineEmits<{ 'update:modelValue': [value: any] }>()

const open = ref(false)
const search = ref('')

const filteredItems = computed(() => {
  if (!search.value.trim()) return props.items
  const q = search.value.toLowerCase()
  return props.items.filter(i => i.label.toLowerCase().includes(q))
})

const selectedLabel = computed(() =>
  props.items.find(i => i.value === props.modelValue)?.label ?? ''
)

const isSelected = (val: any) => props.modelValue === val

const select = (val: any) => {
  emit('update:modelValue', val)
  open.value = false
  search.value = ''
}
</script>

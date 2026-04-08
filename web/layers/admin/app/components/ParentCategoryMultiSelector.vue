<template>
  <UFormField :label="label || t('admin.posts.categories.parent_category')">
    <div class="space-y-2">
      <!-- 已选择的标签 -->
      <div v-if="selectedItems.length > 0" class="flex flex-wrap gap-1.5">
        <UBadge
          v-for="item in selectedItems"
          :key="item.term_id"
          :label="item.name"
          color="primary"
          variant="soft"
          size="sm"
          class="cursor-pointer"
          @click="removeItem(item.term_id)">
          <template #trailing>
            <UIcon name="i-tabler-x" class="size-3" />
          </template>
        </UBadge>
      </div>

      <!-- 搜索框 -->
      <UInput
        v-model="catSearch"
        :placeholder="t('admin.posts.editor.search_categories')"
        leading-icon="i-tabler-search"
        size="sm"
        class="w-full">
        <template v-if="catSearch" #trailing>
          <UButton icon="i-tabler-x" color="neutral" variant="ghost" size="xs" @click="catSearch = ''" />
        </template>
      </UInput>

      <!-- 选择列表 -->
      <div
        class="border border-default rounded-md max-h-48 overflow-y-auto divide-y divide-default">
        <div
          v-if="filteredParents.length === 0"
          class="p-3 text-sm text-muted text-center">
          {{ catSearch ? t('common.no_results') : '没有可选项' }}
        </div>
        <label
          v-for="item in filteredParents"
          :key="item.term_id"
          class="flex items-center gap-2 px-3 py-2 cursor-pointer hover:bg-elevated transition-colors"
          :class="{ 'pl-8': item.level && item.level > 0 }">
          <UCheckbox
            :model-value="isSelected(item.term_id)"
            :disabled="isDisabled(item.term_id)"
            @update:model-value="handleChange(item.term_id)" />
          <span class="text-sm text-highlighted">
            {{ "　".repeat(item.level || 0) }}{{ item.name }}
          </span>
        </label>
      </div>
    </div>
  </UFormField>
</template>

<script setup lang="ts">
const { t } = useI18n();

const props = withDefaults(
  defineProps<{
    modelValue?: number[];
    excludeId?: number;
    label?: string;
  }>(),
  { modelValue: () => [] },
);

const emit = defineEmits<{ "update:modelValue": [value: number[]] }>();

const categoryStore = useCategoryStore();
const { getFlattenedParents } = categoryStore;

const selectedIds = computed({
  get: () => new Set(props.modelValue || []),
  set: (val) => emit("update:modelValue", Array.from(val)),
});

const isSelected = (termId: number) => selectedIds.value.has(termId);
const isDisabled = (termId: number) =>
  !!(props.excludeId && termId === props.excludeId);

const handleChange = (termId: number) => {
  const newSet = new Set(selectedIds.value);
  newSet.has(termId) ? newSet.delete(termId) : newSet.add(termId);
  selectedIds.value = newSet;
};

const removeItem = (termId: number) => {
  const newSet = new Set(selectedIds.value);
  newSet.delete(termId);
  selectedIds.value = newSet;
};

const catSearch = ref('')

const availableParents = computed(() => getFlattenedParents(props.excludeId))

const filteredParents = computed(() => {
  if (!catSearch.value.trim()) return availableParents.value
  const q = catSearch.value.toLowerCase()
  return availableParents.value.filter(item => item.name.toLowerCase().includes(q))
})

const selectedItems = computed(() =>
  availableParents.value.filter((item) => isSelected(item.term_id)),
);
</script>

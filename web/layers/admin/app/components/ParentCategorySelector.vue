<template>
  <UFormField :label="computedLabel">
    <AdminSearchableSelect
      :model-value="modelValue ?? undefined"
      :items="selectItems"
      :placeholder="computedEmptyLabel"
      :search-placeholder="t('admin.posts.editor.search_categories')"
      trigger-class="w-full justify-between"
      @update:model-value="$emit('update:modelValue', $event)"
    />
  </UFormField>
</template>

<script setup lang="ts">
const props = defineProps<{
  modelValue?: number | null;
  excludeId?: number;
  label?: string;
  emptyLabel?: string;
}>();

defineEmits<{ "update:modelValue": [value: number | undefined] }>();

const { t } = useI18n();
const categoryStore = useCategoryStore();
const { getFlattenedParents } = categoryStore;

const computedLabel = computed(() => props.label ?? t('admin.posts.categories.parent_category'));
const computedEmptyLabel = computed(() => props.emptyLabel ?? t('admin.posts.categories.no_parent'));

const availableParents = computed(() => getFlattenedParents(props.excludeId));

const selectItems = computed(() => [
  { label: computedEmptyLabel.value, value: undefined },
  ...availableParents.value.map((item) => ({
    label: `${"　".repeat(item.level || 0)}${item.name}`,
    value: item.term_id,
  })),
]);
</script>

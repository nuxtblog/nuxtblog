<script setup lang="ts">
defineProps<{
  localValue: string;
  sortBy: string;
  sortOrder: "asc" | "desc";
  viewMode: string;
  sortIcon: Record<string, string>;
}>();

const emit = defineEmits([
  "input",
  "selectSort",
  "update:viewMode",
  "toggleSort",
]);

const sortDropdown = ref<HTMLDetailsElement | null>(null);

const selectSort = (v: string) => {
  emit("selectSort", v);
  if (sortDropdown.value) sortDropdown.value.open = false;
};
</script>

<template>
  <!-- 搜索 + 排序 -->
  <div class="flex items-center gap-2">
    <!-- 搜索 -->
    <div class="relative flex-1">
      <UInput
        :model-value="localValue"
        @input="$emit('input', $event)"
        :placeholder="$t('admin.media.search_placeholder')"
        size="sm"
        class="w-48"
        :ui="{ base: 'pl-8' }">
        <template #leading>
          <UIcon name="i-tabler-search" class="size-4 text-muted" />
        </template>
      </UInput>
    </div>

    <!-- 排序字段 -->
    <UDropdownMenu
      :items="[
        [
          { label: $t('admin.media.sort_created'), icon: 'i-tabler-calendar-plus', onSelect: () => selectSort('created_at') },
          { label: $t('admin.media.sort_name'), icon: 'i-tabler-file-text', onSelect: () => selectSort('file_name') },
          { label: $t('admin.media.sort_size'), icon: 'i-tabler-scan', onSelect: () => selectSort('size') },
          { label: $t('admin.media.sort_updated'), icon: 'i-tabler-calendar', onSelect: () => selectSort('updated_at') },
        ],
      ]">
      <UButton variant="ghost" color="neutral" size="sm">
        <UIcon :name="sortIcon[sortBy]" class="w-4 h-4" />
      </UButton>
    </UDropdownMenu>

    <!-- 升降序 -->
    <UButton
      variant="ghost"
      color="neutral"
      size="sm"
      @click="$emit('toggleSort')">
      <UIcon
        :name="sortOrder === 'asc' ? 'i-tabler-arrow-up-narrow-wide' : 'i-tabler-arrow-down-wide-narrow'"
        class="w-4 h-4" />
    </UButton>

    <!-- 视图切换 -->
    <ViewModeSwitcher
      :model-value="viewMode"
      @update:model-value="$emit('update:viewMode', $event)"
      :modes="[
        { value: 'grid', title: $t('admin.media.grid_view') },
        { value: 'list', title: $t('admin.media.list_view') },
      ]" />
  </div>
</template>

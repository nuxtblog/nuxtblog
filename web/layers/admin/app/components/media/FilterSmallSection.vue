<script setup lang="ts">
interface MediaTypeStat {
  type: string;
  name: string;
  icon: string;
  count?: number;
}

defineProps<{
  statsData: MediaTypeStat[];
  filterType: string;
  selectedMedia: number[];
  sortBy: string;
  sortOrder: "asc" | "desc";
  viewMode: string;
  sortIcon: Record<string, string>;
}>();

const emit = defineEmits([
  "update:filterType",
  "batchDelete",
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
  <div class="flex items-center gap-2">
    <!-- 类型筛选 dropdown -->
    <UDropdownMenu
      :items="[
        statsData.map((data) => ({
          label: data.name,
          slot: data.type,
          onSelect: () => emit('update:filterType', data.type),
        })),
      ]"
      class="flex-1">
      <UButton variant="ghost" color="neutral" size="sm" class="w-full justify-start">
        <UIcon name="i-tabler-stack-2" class="w-4 h-4" />
        {{ $t('admin.media.filter_type') }}
      </UButton>
      <template v-for="data in statsData" #[data.type]="{ item }">
        <div class="flex items-center gap-2 w-full">
          <UIcon :name="data.icon" class="size-4" />
          <span :class="filterType === data.type ? 'font-semibold' : ''">{{ data.name }}</span>
          <UBadge v-if="data.count" color="neutral" variant="subtle" size="sm" class="ml-auto">
            {{ data.count }}
          </UBadge>
        </div>
      </template>
    </UDropdownMenu>

    <!-- 批量删除 -->
    <UButton
      v-if="selectedMedia.length > 0"
      color="error"
      size="xs"
      :title="$t('admin.media.selected_n', { n: selectedMedia.length })"
      @click="emit('batchDelete')">
      {{ $t('common.delete') }}
    </UButton>

    <!-- 排序 -->
    <div class="flex items-center gap-2">
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

      <UButton
        variant="ghost"
        color="neutral"
        size="sm"
        @click="emit('toggleSort')">
        <UIcon
          :name="sortOrder === 'asc' ? 'i-tabler-arrow-up-narrow-wide' : 'i-tabler-arrow-down-wide-narrow'"
          class="w-4 h-4" />
      </UButton>

      <ViewModeSwitcher
        :model-value="viewMode"
        @update:model-value="emit('update:viewMode', $event)"
        :modes="[
          { value: 'grid', title: $t('admin.media.grid_view') },
          { value: 'list', title: $t('admin.media.list_view') },
        ]" />
    </div>
  </div>
</template>

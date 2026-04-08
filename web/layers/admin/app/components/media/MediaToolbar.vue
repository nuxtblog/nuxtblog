<script setup lang="ts">
import FilterMiddleSection from "./FilterMiddleSection.vue";
import FilterRightSection from "./FilterRightSection.vue";
import FilterSearchInput from "./FilterSearchInput.vue";
import FilterSmallSection from "./FilterSmallSection.vue";
import FilterTypeButtons from "./FilterTypeButtons.vue";

interface MediaTypeStat {
  type: string;
  name: string;
  icon: string;
  count?: number;
}

const props = withDefaults(
  defineProps<{
    statsData: MediaTypeStat[];
    filterType: string;
    sortBy: string;
    sortOrder: "asc" | "desc";
    searchQuery: string;
    selectedMedia: number[];
    viewMode: string;
    debounce?: number;
  }>(),
  {
    debounce: 300,
  }
);

const emit = defineEmits([
  "update:filterType",
  "update:sortBy",
  "update:sortOrder",
  "update:searchQuery",
  "update:viewMode",
  "batchDelete",
]);

const localValue = ref(props.searchQuery);
let timer: number | null = null;

const onInput = (e: Event) => {
  const value = (e.target as HTMLInputElement).value;
  localValue.value = value;

  if (timer) {
    clearTimeout(timer);
  }

  timer = window.setTimeout(() => {
    emit("update:searchQuery", value);
  }, props.debounce);
};

const sortDropdown = ref<HTMLDetailsElement | null>(null);

const sortIcon: Record<string, string> = {
  created_at: "i-tabler-calendar-plus",
  file_name: "i-tabler-file-text",
  size: "i-tabler-scan",
  updated_at: "i-tabler-calendar",
};

const isAsc = computed({
  get: () => props.sortOrder === "asc",
  set: (v: boolean) => emit("update:sortOrder", v ? "asc" : "desc"),
});

const selectSort = (v: string) => {
  emit("update:sortBy", v);
  if (sortDropdown.value) sortDropdown.value.open = false;
};

const updateViewMode = (v: string) => emit("update:viewMode", v);
</script>

<template>
  <div class="pb-4 bg-background border-b border-default">
    <!-- lg 以上：左右布局 -->
    <div class="hidden lg:flex flex-row gap-4">
      <FilterTypeButtons
        :stats-data="props.statsData"
        :filter-type="props.filterType"
        :selected-media="props.selectedMedia"
        @update:filter-type="emit('update:filterType', $event)"
        @batch-delete="emit('batchDelete')" />

      <FilterRightSection
        :local-value="localValue"
        :sort-by="props.sortBy"
        :sort-order="props.sortOrder"
        :view-mode="props.viewMode"
        :sort-icon="sortIcon"
        @input="onInput"
        @select-sort="selectSort"
        @update:view-mode="updateViewMode"
        @toggle-sort="isAsc = !isAsc" />
    </div>

    <!-- md ~ lg：上下布局 -->
    <div class="hidden md:flex lg:hidden flex-col gap-4">
      <FilterTypeButtons
        :stats-data="props.statsData"
        :filter-type="props.filterType"
        :selected-media="props.selectedMedia"
        @update:filter-type="emit('update:filterType', $event)"
        @batch-delete="emit('batchDelete')" />

      <FilterMiddleSection
        :local-value="localValue"
        :sort-by="props.sortBy"
        :sort-order="props.sortOrder"
        :view-mode="props.viewMode"
        :sort-icon="sortIcon"
        @input="onInput"
        @select-sort="selectSort"
        @update:view-mode="updateViewMode"
        @toggle-sort="isAsc = !isAsc" />
    </div>

    <!-- sm ~ md：移动平板 -->
    <div class="hidden sm:flex md:hidden flex-col gap-3">
      <FilterSmallSection
        :stats-data="props.statsData"
        :filter-type="props.filterType"
        :selected-media="props.selectedMedia"
        :sort-by="props.sortBy"
        :sort-order="props.sortOrder"
        :view-mode="props.viewMode"
        :sort-icon="sortIcon"
        @update:filter-type="emit('update:filterType', $event)"
        @batch-delete="emit('batchDelete')"
        @select-sort="selectSort"
        @update:view-mode="updateViewMode"
        @toggle-sort="isAsc = !isAsc" />

      <FilterSearchInput :local-value="localValue" @input="onInput" />
    </div>

    <!-- < sm：超小屏 -->
    <div class="flex sm:hidden flex-col gap-3">
      <FilterSmallSection
        :stats-data="props.statsData"
        :filter-type="props.filterType"
        :selected-media="props.selectedMedia"
        :sort-by="props.sortBy"
        :sort-order="props.sortOrder"
        :view-mode="props.viewMode"
        :sort-icon="sortIcon"
        @update:filter-type="emit('update:filterType', $event)"
        @batch-delete="emit('batchDelete')"
        @select-sort="selectSort"
        @update:view-mode="updateViewMode"
        @toggle-sort="isAsc = !isAsc" />

      <FilterSearchInput :local-value="localValue" @input="onInput" />
    </div>
  </div>
</template>

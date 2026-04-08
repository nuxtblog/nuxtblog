<template>
  <UPagination
    v-if="totalPages > 1"
    v-model:page="currentPage"
    :total="total"
    :items-per-page="pageSize"
  />
</template>

<script setup lang="ts">
export interface PageNavigatorProps {
  total: number;
  pageSize?: number;
}

const props = withDefaults(defineProps<PageNavigatorProps>(), {
  pageSize: 12,
});

const currentPage = defineModel<number>("currentPage", { default: 1 });

const totalPages = computed(() => Math.ceil(props.total / props.pageSize));

const resetPage = () => {
  currentPage.value = 1;
};

const goToPage = (page: number) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page;
  }
};

defineExpose({ resetPage, goToPage, totalPages });
</script>

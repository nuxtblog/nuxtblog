<!-- FilterTypeButtons.vue -->
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
}>();

const emit = defineEmits(["update:filterType", "batchDelete"]);
</script>

<template>
  <div class="flex flex-wrap items-center gap-2 flex-1">
    <!-- 文件类型 -->
    <div class="flex items-center gap-2 text-sm">
      <UButton
        v-for="data in statsData"
        :key="data.type"
        size="sm"
        :variant="filterType === data.type ? 'solid' : 'ghost'"
        :color="filterType === data.type ? 'primary' : 'neutral'"
        @click="emit('update:filterType', data.type)">
        <UIcon :name="data.icon" class="size-4" />
        {{ data.name }}
        <span v-if="data.count" class="text-xs">({{ data.count }})</span>
      </UButton>
    </div>

    <!-- 批量操作 -->
    <div
      v-if="selectedMedia.length > 0"
      class="flex items-center gap-2 ml-4 pl-4 border-l border-default">
      <span class="text-xs text-muted">
        {{ $t('admin.media.selected_n', { n: selectedMedia.length }) }}
      </span>
      <UButton size="sm" color="error" @click="emit('batchDelete')">
        {{ $t('common.delete') }}
      </UButton>
    </div>
  </div>
</template>

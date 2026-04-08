<template>
  <VueDraggable
    class="drag-area"
    tag="ul"
    :model-value="modelValue"
    :group="{ name: 'categories', pull: true, put: true }"
    :data-parent-id="parentId"
    @update:model-value="$emit('update:modelValue', $event)"
    @end="onEnd"
  >
    <li v-for="item in modelValue" :key="item.term_id" :data-term-id="item.term_id">
      <!-- 节点内容 -->
      <div class="group rounded-md border border-default bg-default hover:border-primary/40 hover:shadow-sm px-3 py-2.5 flex items-center gap-3 transition-all">
        <!-- 拖拽手柄 -->
        <UIcon
          name="i-tabler-grip-vertical"
          class="size-4 shrink-0 cursor-grab text-muted hover:text-primary active:cursor-grabbing"
        />

        <!-- 展开/折叠按钮 -->
        <button
          v-if="item.children?.length"
          @click="toggleExpanded(item.term_id)"
          class="flex h-6 w-6 items-center justify-center shrink-0 rounded hover:bg-elevated transition-colors"
        >
          <UIcon
            name="i-tabler-chevron-right"
            class="size-4 text-muted transition-transform duration-200"
            :class="{ 'rotate-90': expanded.has(item.term_id) }"
          />
        </button>
        <div v-else class="h-6 w-6 shrink-0" />

        <!-- 内容区 -->
        <div class="min-w-0 flex-1">
          <div class="flex items-baseline gap-2 mb-0.5">
            <h3 class="truncate text-sm font-medium text-highlighted">{{ item.name }}</h3>
            <UBadge :label="String(item.term_id)" color="primary" variant="soft" size="xs" class="font-mono shrink-0" />
          </div>
          <div class="flex flex-wrap items-center gap-2 text-xs text-muted">
            <span>{{ item.slug }}</span>
            <span>{{ item.count }} 篇</span>
            <UBadge v-if="item.children?.length" :label="`${item.children.length} 个子分类`" color="neutral" variant="soft" size="xs" />
          </div>
        </div>

        <!-- 操作按钮 -->
        <div class="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity shrink-0">
          <UButton
            color="neutral"
            variant="ghost"
            icon="i-tabler-pencil"
            square
            size="xs"
            @click="$emit('edit', item)"
          />
          <UButton
            color="error"
            variant="ghost"
            icon="i-tabler-trash"
            square
            size="xs"
            @click="$emit('delete', item)"
          />
        </div>
      </div>

      <!-- 子节点 -->
      <CategoryNode
        v-if="expanded.has(item.term_id)"
        :model-value="item.children"
        :parent-id="item.term_id"
        @update:modelValue="item.children = $event"
        @edit="$emit('edit', $event)"
        @delete="$emit('delete', $event)"
      />
    </li>
  </VueDraggable>
</template>

<script setup lang="ts">
import { VueDraggable, type SortableEvent } from "vue-draggable-plus";
import type { TermDetailWithChildren } from "~/types/models/term";

const { updateCategoryParent } = useCategoryStore();

interface Props {
  modelValue: TermDetailWithChildren[];
  parentId?: number;
}

const props = withDefaults(defineProps<Props>(), { parentId: 0 });

const emit = defineEmits<{
  "update:modelValue": [value: TermDetailWithChildren[]];
  edit: [item: TermDetailWithChildren];
  delete: [item: TermDetailWithChildren];
}>();

const expanded = ref<Set<number>>(new Set());

const initializeExpanded = () => {
  const newExpanded = new Set<number>();
  const traverse = (items: TermDetailWithChildren[]) => {
    items.forEach((item) => {
      newExpanded.add(item.term_id);
      if (item.children?.length) traverse(item.children);
    });
  };
  traverse(props.modelValue);
  expanded.value = newExpanded;
};

watch(() => props.modelValue, initializeExpanded, { immediate: true, deep: true });

const toggleExpanded = (termId: number) => {
  if (expanded.value.has(termId)) {
    expanded.value.delete(termId);
  } else {
    expanded.value.add(termId);
  }
};

const onEnd = (event: SortableEvent) => {
  const categoryId = Number((event.item as HTMLElement).getAttribute("data-term-id"));
  const newParentId = Number((event.to as HTMLElement).getAttribute("data-parent-id"));
  if (!categoryId || isNaN(categoryId)) return;
  // newParentId === 0 表示拖到顶级（无父级），用 null 通知后端清除父级
  updateCategoryParent(categoryId, isNaN(newParentId) ? undefined : newParentId);
};
</script>

<style scoped>
.drag-area {
  list-style: none;
  margin: 0;
  padding: 0;
  border-left: 2px dashed var(--ui-border);
  margin-left: 16px;
  padding-left: 8px;
}

.drag-area > li {
  margin-bottom: 6px;
}
</style>

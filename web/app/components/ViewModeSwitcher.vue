<template>
  <div class="flex">
    <UButton
      v-for="mode in modes"
      :key="mode.value"
      :icon="getIconName(mode.value)"
      :title="mode.title"
      size="xs"
      :color="mode.value === modelValue ? 'primary' : 'neutral'"
      :variant="mode.value === modelValue ? 'solid' : 'ghost'"
      @click="selectMode(mode.value)"
    />
  </div>
</template>

<script setup lang="ts">
interface ModeOption {
  value: string;
  title: string;
  icon?: string;
}

const props = defineProps<{
  modes: ModeOption[];
  modelValue: string;
}>();

const emit = defineEmits<{
  "update:modelValue": [value: string];
}>();

const defaultIcons: Record<string, string> = {
  tree: "i-tabler-folder",
  list: "i-tabler-list",
  grid: "i-tabler-layout-grid",
  table: "i-tabler-table",
  kanban: "i-tabler-layout-columns",
};

const getIconName = (value: string): string => {
  const mode = props.modes.find((m) => m.value === value);
  return mode?.icon || defaultIcons[value] || defaultIcons.list;
};

const selectMode = (value: string) => {
  if (value !== props.modelValue) {
    emit("update:modelValue", value);
  }
};
</script>

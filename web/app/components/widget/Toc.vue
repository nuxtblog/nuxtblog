<script setup lang="ts">
import { useCurrentPost, type TocItem } from "~/composables/useCurrentPost";

import type { WidgetConfig } from "~/composables/useWidgetConfig";
const props = defineProps<{ config: WidgetConfig }>();
const { t } = useI18n()
const title = computed(() => props.config.title || t('site.widget.toc.default_title'));

const { toc } = useCurrentPost();
const activeId = ref("");

useEventListener("scroll", () => {
  const els = document.querySelectorAll<HTMLElement>("h2[id], h3[id], h4[id]");
  let current = "";
  for (const el of els) {
    if (el.offsetTop - 120 <= window.scrollY) current = el.id;
  }
  activeId.value = current;
});

const scrollTo = (id: string) => {
  const el = document.getElementById(id);
  if (el) el.scrollIntoView({ behavior: "smooth", block: "start" });
};

const indentClass = (item: TocItem) =>
  item.level === 2 ? "pl-2" : item.level === 3 ? "pl-6" : "pl-10";
</script>

<template>
  <UCard
    v-if="toc.length > 0"
    :ui="{ header: 'p-2.5 sm:px-4 flex', body: 'p-0 sm:p-0' }">
    <template #header>
      <div class="flex items-center gap-2">
        <h3 class="font-semibold text-highlighted">{{ title }}</h3>
      </div>
    </template>
    <nav class="px-2 py-2 space-y-0.5 max-h-80 overflow-y-auto">
      <button
        v-for="item in toc"
        :key="item.id"
        class="group flex items-start gap-1.5 w-full text-left py-1 px-2 rounded-md transition-colors"
        :class="[
          indentClass(item),
          activeId === item.id
            ? 'text-primary bg-primary/10 font-medium'
            : 'text-muted hover:text-highlighted hover:bg-muted',
        ]"
        @click="scrollTo(item.id)">
        <span
          class="mt-1.5 size-1.5 rounded-full shrink-0 transition-colors"
          :class="
            activeId === item.id
              ? 'bg-primary'
              : 'bg-muted group-hover:bg-highlighted'
          " />
        <span class="text-sm leading-snug">{{ item.text }}</span>
      </button>
    </nav>
  </UCard>
</template>

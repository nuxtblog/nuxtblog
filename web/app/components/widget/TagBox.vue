<script setup lang="ts">
import type { WidgetConfig } from '~/composables/useWidgetConfig'
const props = defineProps<{ config: WidgetConfig }>()
const { t } = useI18n()
const { abbreviate } = useNumberFormat()
const maxCount = computed(() => props.config.maxCount ?? 15)
const title = computed(() => props.config.title || t('site.widget.tags.default_title'))

const termApi = useTermApi();
const { data } = await useAsyncData("widget-tags", () =>
  termApi.getTerms({ taxonomy: "tag" }).catch(() => []),
);
const tags = computed(() => (data.value ?? []).slice(0, maxCount.value));

</script>

<template>
  <UCard
    :ui="{ header: 'p-2.5 sm:px-4', footer: 'p-2.5 sm:px-4 flex justify-end' }">
    <template #header>
      <h3 class="font-semibold text-highlighted">{{ title }}</h3>
    </template>

    <div class="flex flex-wrap gap-2">
      <UBadge
        v-for="tag in tags"
        :key="tag.term_id"
        color="primary"
        variant="soft"
        :as="'a'"
        :href="`/posts?tag=${tag.term_id}`"
        class="cursor-pointer">
        {{ tag.name }}
        <span class="text-xs opacity-70 ml-1">{{
          abbreviate(tag.count)
        }}</span>
      </UBadge>
    </div>

    <template #footer>
      <UButton
        variant="link"
        color="primary"
        to="/tags"
        trailing-icon="i-tabler-arrow-right"
        size="sm">
        {{ $t('site.widget.tags.view_all') }}
      </UButton>
    </template>
  </UCard>
</template>

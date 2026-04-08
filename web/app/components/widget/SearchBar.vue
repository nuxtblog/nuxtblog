<script setup lang="ts">
const { abbreviate } = useNumberFormat()
interface SearchResult {
  id: number;
  title: string;
  excerpt: string;
  slug: string;
}

interface HotKeyword {
  id: number;
  keyword: string;
  count: number;
}

import type { WidgetConfig } from "~/composables/useWidgetConfig";
const props = defineProps<{ config: WidgetConfig }>();
const { t } = useI18n()
const searchTitle = computed(() => props.config.title || t('site.widget.search.default_title'));
const showRecentEnabled = computed(() => props.config.showRecent ?? true);
const showHotEnabled = computed(() => props.config.showHot ?? true);

const { apiFetch } = useApiFetch();

const searchQuery = ref("");
const isSearching = ref(false);
const isFocused = ref(false);
const searchResults = ref<SearchResult[]>([]);

const hotKeywords = ref<HotKeyword[]>([
  { id: 1, keyword: "Vue 3", count: 128 },
  { id: 2, keyword: "TypeScript", count: 96 },
  { id: 3, keyword: "性能优化", count: 87 },
  { id: 4, keyword: "Tailwind CSS", count: 74 },
  { id: 5, keyword: "Vite", count: 69 },
]);

const maxRecent = 5;
const recentSearches = useLocalStorage<string[]>("blog:search:recent", []);

const showResults = computed(
  () => isFocused.value && searchQuery.value.length > 0,
);

const debouncedQuery = refDebounced(searchQuery, 350);
watch(debouncedQuery, (q) => {
  if (q.trim()) handleSearch();
  else searchResults.value = [];
});

const handleSearch = async () => {
  const q = searchQuery.value.trim();
  if (!q) return;
  isSearching.value = true;
  try {
    const data = await apiFetch<{ data: SearchResult[] }>("/posts", {
      params: {
        keyword: q,
        status: "published",
        page_size: 8,
        post_type: "post",
      },
    });
    searchResults.value = (data.data || []).map((p) => ({
      id: p.id,
      title: p.title,
      excerpt: p.excerpt,
      slug: p.slug,
    }));
  } catch {
    searchResults.value = [];
  } finally {
    isSearching.value = false;
  }
  if (!recentSearches.value.includes(q)) {
    recentSearches.value.unshift(q);
    if (recentSearches.value.length > maxRecent) recentSearches.value.pop();
  }
};

const handleKeywordClick = (keyword: string) => {
  searchQuery.value = keyword;
  handleSearch();
};

const clearSearch = () => {
  searchQuery.value = "";
  searchResults.value = [];
};

const removeRecentSearch = (index: number) => {
  recentSearches.value.splice(index, 1);
};


// ─── Teleport 定位 ────────────────────────────────────────────────────────────
const inputRef = useTemplateRef<HTMLElement>("inputWrapper");
const { bottom, left, width } = useElementBounding(inputRef);

const dropdownStyle = computed(() => ({
  top: `${bottom.value + 8}px`,
  left: `${left.value}px`,
  width: `${width.value}px`,
}));
</script>

<template>
  <UCard :ui="{ header: 'p-2.5 sm:px-4' }">
    <template #header>
      <h3 class="font-semibold text-highlighted">{{ searchTitle }}</h3>
    </template>

    <div class="space-y-4">
      <!-- 搜索框 -->
      <div ref="inputWrapper" class="relative">
        <UInput
          v-model="searchQuery"
          :placeholder="$t('site.widget.search.placeholder')"
          leading-icon="i-tabler-search"
          :loading="isSearching"
          class="w-full"
          @focus="isFocused = true"
          @blur="isFocused = false"
          @keyup.enter="handleSearch">
          <template v-if="searchQuery" #trailing>
            <UButton icon="i-tabler-x" color="neutral" variant="ghost" size="xs" @mousedown.prevent="clearSearch" />
          </template>
        </UInput>

        <Teleport to="body">
          <div
            v-if="showResults"
            class="fixed z-50 bg-default ring ring-default rounded-md shadow-lg overflow-hidden"
            :style="dropdownStyle"
            @mousedown.prevent>
            <div
              v-if="searchResults.length > 0"
              class="max-h-80 overflow-y-auto divide-y divide-default">
              <a
                v-for="result in searchResults"
                :key="result.id"
                :href="`/post/${result.slug}`"
                class="block px-4 py-3 hover:bg-muted transition-colors">
                <h4
                  class="text-sm font-medium text-highlighted mb-1 line-clamp-1">
                  {{ result.title }}
                </h4>
                <p
                  v-if="result.excerpt"
                  class="text-xs text-muted line-clamp-2">
                  {{ result.excerpt }}
                </p>
              </a>
            </div>
            <div v-else class="px-4 py-8 text-center">
              <p class="text-sm text-muted">{{ $t('site.widget.search.no_results') }}</p>
            </div>
          </div>
        </Teleport>
      </div>

      <!-- 最近搜索 -->
      <div
        v-if="showRecentEnabled && recentSearches.length > 0 && !showResults"
        class="space-y-2">
        <div class="flex items-center justify-between">
          <h4 class="text-xs font-medium text-muted">{{ $t('site.widget.search.recent') }}</h4>
          <UButton
            variant="link"
            color="neutral"
            size="xs"
            @click="recentSearches = []">
            {{ $t('site.widget.search.clear') }}
          </UButton>
        </div>
        <div class="space-y-1">
          <div
            v-for="(search, index) in recentSearches"
            :key="index"
            class="group flex items-center justify-between px-3 rounded-md hover:bg-muted transition-colors cursor-pointer"
            @click="handleKeywordClick(search)">
            <div class="flex items-center gap-2 flex-1 min-w-0">
              <UIcon name="i-tabler-clock" class="size-3 text-muted shrink-0" />
              <span class="text-xs text-default truncate">{{ search }}</span>
            </div>
            <UButton
              icon="i-tabler-x"
              variant="ghost"
              color="neutral"
              size="xs"
              class="opacity-0 group-hover:opacity-100 transition-opacity"
              @click.stop="removeRecentSearch(index)" />
          </div>
        </div>
      </div>

      <!-- 热门搜索 -->
      <div v-if="showHotEnabled" class="space-y-2">
        <h4 class="text-xs font-medium text-muted">{{ $t('site.widget.search.hot') }}</h4>
        <div class="flex flex-wrap gap-2">
          <UBadge
            v-for="item in hotKeywords"
            :key="item.id"
            color="primary"
            variant="soft"
            class="cursor-pointer"
            @click="handleKeywordClick(item.keyword)">
            {{ item.keyword }}
            <span class="text-xs opacity-70 ml-1">{{
              abbreviate(item.count)
            }}</span>
          </UBadge>
        </div>
      </div>
    </div>
  </UCard>
</template>

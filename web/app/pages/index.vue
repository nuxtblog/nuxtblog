<template>
  <!-- Hero Banner -->
  <div class="container mx-auto px-px md:px-4">
    <HomeBanner :posts="bannerPosts" />
  </div>

  <!-- 主内容区 -->
  <main class="container mx-auto px-px md:px-4 py-4 md:py-12">
    <div :class="sidebarEnabled ? 'grid grid-cols-1 lg:grid-cols-12 gap-8' : 'space-y-12'">
      <div :class="sidebarEnabled ? 'lg:col-span-9 space-y-12' : 'space-y-12'">

        <ClientOnly><ContributionSlot name="public:home-top" /></ClientOnly>

        <section v-for="section in activeSections" :key="section.id">
          <!-- Section Header -->
          <div class="flex items-center gap-3 mb-4 md:mb-6">
            <UIcon :name="sectionIcon(section.id)" class="size-5 text-primary" />
            <h3 class="md:text-2xl font-bold text-highlighted">{{ section.title }}</h3>
            <UBadge v-if="section.id === 'latest'" color="primary" variant="soft" size="sm">New</UBadge>
            <!-- Action button -->
            <div v-if="section.action?.enabled" class="ml-auto flex items-center">
              <UButton
                v-if="section.id === 'random'"
                variant="ghost"
                color="neutral"
                size="sm"
                icon="i-tabler-refresh"
                :loading="randomLoading"
                @click="refreshRandom">
                {{ section.action.label || $t('common.refresh') }}
              </UButton>
              <NuxtLink
                v-else
                :to="getActionHref(section)"
                class="flex items-center gap-0.5 text-sm text-primary hover:opacity-75 transition-opacity font-medium">
                {{ section.action.label || $t('common.view_more') }}
                <UIcon name="i-tabler-chevron-right" class="size-4" />
              </NuxtLink>
            </div>
          </div>

          <HomePostGridWithImage   v-if="section.layout === 'grid'"     :posts="sectionPosts(section.id)" :max-length="section.count" />
          <HomePostListWithImage   v-else-if="section.layout === 'list'"    :posts="sectionPosts(section.id)" :max-length="section.count" />
          <HomePostSimpleList      v-else-if="section.layout === 'simple'"  :posts="sectionPosts(section.id)" :max-length="section.count" />
          <HomePostMasonry         v-else-if="section.layout === 'masonry'" :posts="sectionPosts(section.id)" :max-length="section.count" />
          <HomePostFeaturedHero    v-else-if="section.layout === 'hero'"    :posts="sectionPosts(section.id)" :max-length="section.count" />
          <HomePostTimeline        v-else-if="section.layout === 'timeline'":posts="sectionPosts(section.id)" :max-length="section.count" />
          <HomePostRanking         v-else-if="section.layout === 'ranking'" :posts="sectionPosts(section.id)" :max-length="section.count" />
          <HomePostGridWithImage   v-else                                    :posts="sectionPosts(section.id)" :max-length="section.count" />

          <!-- Masonry load more -->
          <div v-if="section.id === 'masonry' && section.loadMoreEnabled && masonryHasMore" class="flex justify-center mt-6">
            <UButton
              variant="outline"
              color="neutral"
              :loading="masonryLoading"
              icon="i-tabler-loader-2"
              @click="loadMoreMasonry">
              {{ $t('common.load_more') }}
            </UButton>
          </div>
        </section>

        <ClientOnly><ContributionSlot name="public:home-bottom" /></ClientOnly>

      </div>

      <!-- 侧栏 -->
      <BlogSidebar v-if="sidebarEnabled" option-key="homepage_sidebar_widgets" />
    </div>
  </main>
</template>

<script setup lang="ts">
import type { PostCard } from "~/types/models/post";
import type { SectionConfig } from "~/composables/useHomepageSections";
import { transformPostList } from "~/utils/transformers/post";
import { SECTION_DEFAULTS } from "~/composables/useHomepageSections";

const { defaultCover } = usePostCover();
const postApi = usePostApi();
const optionsStore = useOptionsStore();

await optionsStore.load();

const sidebarEnabled = computed(() =>
  optionsStore.getJSON<boolean>("homepage_sidebar_enabled", false),
);

const { getSectionConfig } = useHomepageSections();

// Load all section configs
const allSections = computed(() =>
  SECTION_DEFAULTS.map((def) => getSectionConfig(def.id)),
);
const activeSections = computed(() =>
  allSections.value.filter((s) => s.enabled),
);

// Section icon map
const sectionIcon = (id: string) => ({
  latest:   "i-tabler-file-text",
  hot:      "i-tabler-trending-up",
  featured: "i-tabler-star",
  random:   "i-tabler-dice",
  timeline: "i-tabler-timeline",
  masonry:  "i-tabler-layout-columns",
}[id] ?? "i-tabler-layout-grid");

// Action href resolver
const DEFAULT_ACTION_HREF: Record<string, string> = {
  hot:      '/posts?sort=views',
  featured: '/posts?featured=1',
  timeline: '/archive',
  masonry:  '/posts',
  latest:   '/posts',
}
const getActionHref = (section: SectionConfig): string => {
  if (section.action?.href) return section.action.href
  if (section.id === 'latest' && section.action?.categorySlug)
    return `/category/${section.action.categorySlug}`
  return DEFAULT_ACTION_HREF[section.id] ?? '/posts'
}

// Hot posts date range
const hotMonths = parseInt(optionsStore.get("hot_posts_months", "3"), 10) || 3;
const hotAfterDate = new Date();
hotAfterDate.setMonth(hotAfterDate.getMonth() - hotMonths);
const publishedAfter = hotAfterDate.toISOString().slice(0, 10);

const latestCfg  = getSectionConfig("latest");
const hotCfg     = getSectionConfig("hot");
const featuredCfg = getSectionConfig("featured");
const randomCfg  = getSectionConfig("random");
const timelineCfg = getSectionConfig("timeline");
const masonryCfg  = getSectionConfig("masonry");

const latestParams: Record<string, unknown> = {
  page: 1, page_size: latestCfg.count, status: "published",
};
if (latestCfg.includeCategoryIds?.length)
  latestParams.include_category_ids = latestCfg.includeCategoryIds.join(",");
if (latestCfg.excludeCategoryIds?.length)
  latestParams.exclude_category_ids = latestCfg.excludeCategoryIds.join(",");

const bannerPosts    = ref<PostCard[]>([]);
const postsBySection = ref<Record<string, PostCard[]>>({});

// Masonry pagination state
const masonryPage    = ref(1);
const masonryTotal   = ref(0);
const masonryLoading = ref(false);
const masonryHasMore = computed(() =>
  (postsBySection.value.masonry?.length ?? 0) < masonryTotal.value
);

// Random refresh state
const randomLoading = ref(false);

const { data } = await useAsyncData("home-posts", () =>
  Promise.all([
    postApi.getPosts(latestParams as any).catch(() => null),
    postApi.getPosts({ page: 1, page_size: 5, status: "published", meta_key: "is_banner", meta_value: "1" }).catch(() => null),
    postApi.getPosts({ page: 1, page_size: featuredCfg.count, status: "published", meta_key: "is_featured", meta_value: "1" }).catch(() => null),
    postApi.getPosts({ page: 1, page_size: hotCfg.count, status: "published", sort_by: "view_count", published_after: publishedAfter }).catch(() => null),
    postApi.getPosts({ page: 1, page_size: timelineCfg.count, status: "published" }).catch(() => null),
    postApi.getPosts({ page: 1, page_size: masonryCfg.count, status: "published", sort_by: "random" }).catch(() => null),
    postApi.getPosts({ page: 1, page_size: randomCfg.count, status: "published", sort_by: "random" }).catch(() => null),
  ]),
);

if (data.value) {
  const [latestData, bannerData, featuredData, hotData, timelineData, masonryData, randomData] = data.value;
  bannerPosts.value = transformPostList(bannerData?.data ?? [], defaultCover.value);
  postsBySection.value = {
    latest:   transformPostList(latestData?.data   ?? [], defaultCover.value),
    featured: transformPostList(featuredData?.data ?? [], defaultCover.value),
    hot:      transformPostList(hotData?.data      ?? [], defaultCover.value),
    timeline: transformPostList(timelineData?.data ?? [], defaultCover.value),
    masonry:  transformPostList(masonryData?.data  ?? [], defaultCover.value),
    random:   transformPostList(randomData?.data   ?? [], defaultCover.value),
  };
  masonryTotal.value = masonryData?.total ?? 0;
}

const sectionPosts = (id: string) => postsBySection.value[id] ?? [];

// Refresh random section (client-side only)
const refreshRandom = async () => {
  if (randomLoading.value) return;
  randomLoading.value = true;
  try {
    const res = await postApi.getPosts({ page: 1, page_size: randomCfg.count, status: "published", sort_by: "random" });
    postsBySection.value = { ...postsBySection.value, random: transformPostList(res?.data ?? [], defaultCover.value) };
  } catch {}
  finally { randomLoading.value = false; }
};

// Masonry load more
const loadMoreMasonry = async () => {
  if (masonryLoading.value) return;
  masonryLoading.value = true;
  try {
    masonryPage.value++;
    const res = await postApi.getPosts({ page: masonryPage.value, page_size: masonryCfg.count, status: "published", sort_by: "random" });
    const more = transformPostList(res?.data ?? [], defaultCover.value);
    postsBySection.value = { ...postsBySection.value, masonry: [...(postsBySection.value.masonry ?? []), ...more] };
    masonryTotal.value = res?.total ?? masonryTotal.value;
  } catch {}
  finally { masonryLoading.value = false; }
};
</script>

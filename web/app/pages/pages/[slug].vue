<script setup lang="ts">
import { renderMarkdown } from "~/utils/markdown";
import { useCurrentPost, type TocItem } from "~/composables/useCurrentPost";

definePageMeta({ layout: "default" });

const { t } = useI18n()
const route = useRoute();
const slug = route.params.slug as string;
const postApi = usePostApi();
const { defaultCover } = usePostCover();
const optionsStore = useOptionsStore();

await optionsStore.load();

const {
  data: page,
  error,
  pending,
} = await useAsyncData(
  `page-${slug}`,
  () =>
    postApi.getPostBySlug(slug).catch((err) => {
      throw err;
    }),
  { lazy: false, server: true },
);

// ── Layout & sidebar ────────────────────────────────────────────────────────
const defaultLayout = computed(
  () =>
    optionsStore.get("page_default_layout", "none") as "hero" | "half" | "none",
);
const defaultSidebar = computed(() =>
  optionsStore.getJSON<boolean>("page_sidebar_enabled", false),
);
const effectiveLayout = computed(() => {
  const m = page.value?.metas?.post_layout;
  if (m === "hero" || m === "half" || m === "none") return m;
  return defaultLayout.value;
});
const effectiveSidebar = computed(() => {
  const m = page.value?.metas?.post_sidebar;
  if (m === "1") return true;
  if (m === "0") return false;
  return defaultSidebar.value;
});

// ── Derived data ────────────────────────────────────────────────────────────
const coverUrl = computed(
  () => page.value?.featured_img?.url || defaultCover.value,
);

// ── TOC ──────────────────────────────────────────────────────────────────────
const { toc, clearPost } = useCurrentPost();

function slugId(text: string) {
  return (
    text
      .toLowerCase()
      .replace(/\s+/g, "-")
      .replace(/[^\w\u4e00-\u9fa5-]/g, "")
      .replace(/^-+|-+$/g, "") || "h"
  );
}

const contentHtml = computed(() => {
  if (!page.value?.content) return "";
  const markdown = page.value.content;

  const items: TocItem[] = [];
  const idCounts: Record<string, number> = {};
  for (const line of markdown.split("\n")) {
    const m = /^(#{2,4})\s+(.+)$/.exec(line.trim());
    if (!m) continue;
    const level = m[1].length;
    const text = m[2]
      .trim()
      .replace(/`([^`]+)`/g, "$1")
      .replace(/\*+([^*]+)\*+/g, "$1")
      .replace(/\[([^\]]+)\]\([^)]+\)/g, "$1");
    const base = slugId(text);
    const n = (idCounts[base] = (idCounts[base] ?? 0) + 1);
    items.push({ id: n === 1 ? base : `${base}-${n}`, text, level });
  }
  toc.value = items;

  let html = renderMarkdown(markdown);
  let ptr = 0;
  html = html.replace(/<h([234])>([\s\S]*?)<\/h\1>/g, (match, lvl) => {
    if (ptr < items.length && items[ptr].level === parseInt(lvl)) {
      return match.replace(`<h${lvl}>`, `<h${lvl} id="${items[ptr++].id}">`);
    }
    return match;
  });
  return html;
});

onUnmounted(() => clearPost());

const fmtDate = (d?: string) =>
  d
    ? new Date(d).toLocaleDateString("zh-CN", {
        year: "numeric",
        month: "long",
        day: "numeric",
      })
    : "";

// ── SEO ─────────────────────────────────────────────────────────────────────
useHead({
  title: page.value?.title || t('site.post.page_title'),
  meta: [
    { name: "description", content: page.value?.excerpt || "" },
    { property: "og:title", content: page.value?.title || "" },
    { property: "og:description", content: page.value?.excerpt || "" },
    { property: "og:image", content: coverUrl.value },
  ],
  bodyAttrs: {
    'data-page-id': page.value?.id ? String(page.value.id) : '',
    'data-page-slug': page.value?.slug || '',
  },
});

// ── View count ───────────────────────────────────────────────────────────────
onMounted(() => {
  if (page.value?.id)
    postApi.incrementViewCount?.(page.value.id).catch(() => {});
});
</script>

<template>
  <!-- Loading -->
  <div v-if="pending" class="space-y-6 py-8">
    <USkeleton class="w-full h-[40vh] rounded-md" />
    <div class="max-w-3xl mx-auto space-y-4">
      <USkeleton class="h-10 w-3/4" />
      <USkeleton class="h-4 w-full" />
      <USkeleton class="h-4 w-5/6" />
    </div>
  </div>

  <!-- Error -->
  <UCard v-else-if="error" class="my-8">
    <div class="flex flex-col items-center py-20 text-center">
      <UIcon name="i-tabler-circle-x" class="size-12 text-error mb-4" />
      <p class="text-error text-lg mb-4">{{ $t('site.post.load_failed', { msg: error.message }) }}</p>
      <UButton color="primary" @click="navigateTo('/')">{{ $t('site.post.back_to_home') }}</UButton>
    </div>
  </UCard>

  <!-- Page -->
  <div v-else-if="page" class="-mt-8">
    <!-- Plugin slot: page top -->
    <ClientOnly><ContributionSlot name="public:page-before-content" /></ClientOnly>

    <PostPageActionBar :comment-count="page.comment_count" />

    <!-- ╔═══════════════════════════════════════╗ -->
    <!-- ║  HERO LAYOUT  (顶部大图)               ║ -->
    <!-- ╚═══════════════════════════════════════╝ -->
    <template v-if="effectiveLayout === 'hero'">
      <div
        class="relative -mx-4 sm:-mx-6 lg:-mx-8 h-[50vh] min-h-[320px] overflow-hidden">
        <BaseImg
          :src="coverUrl"
          :alt="page.title"
          class="w-full h-full object-cover" />
        <div
          class="absolute inset-0 bg-linear-to-t from-black/80 via-black/30 to-black/10" />
        <div
          class="absolute bottom-0 left-0 right-0 px-4 sm:px-6 lg:px-8 pb-10">
          <div
            :class="
              effectiveSidebar
                ? 'max-w-none lg:max-w-[66%]'
                : 'max-w-3xl mx-auto'
            ">
            <h1
              class="text-3xl md:text-5xl font-bold text-white leading-tight mb-3 drop-shadow-lg">
              {{ page.title }}
            </h1>
            <p
              v-if="page.excerpt"
              class="text-white/70 text-base md:text-lg line-clamp-2">
              {{ page.excerpt }}
            </p>
            <p v-if="page.published_at" class="text-white/50 text-sm mt-2">
              {{ fmtDate(page.published_at) }}
            </p>
          </div>
        </div>
      </div>

      <div
        :class="
          effectiveSidebar
            ? 'grid grid-cols-1 lg:grid-cols-12 gap-8 pt-10'
            : 'max-w-3xl mx-auto pt-10'
        ">
        <article :class="effectiveSidebar ? 'lg:col-span-8' : ''">
          <UCard>
            <MarkdownContent class="prose-lg" :html="contentHtml" />
          </UCard>
          <UCard class="mt-4">
            <PostCommentSection
              :object-id="page.id"
              object-type="post"
              :comment-meta="page.metas?.comment_open" />
          </UCard>
        </article>
        <BlogSidebar
          v-if="effectiveSidebar"
          option-key="page_sidebar_widgets"
          class="lg:col-span-4" />
      </div>
    </template>

    <!-- ╔═══════════════════════════════════════╗ -->
    <!-- ║  HALF LAYOUT  (半图)                   ║ -->
    <!-- ╚═══════════════════════════════════════╝ -->
    <template v-else-if="effectiveLayout === 'half'">
      <div class="grid grid-cols-1 md:grid-cols-5 gap-6 pt-8 mb-8">
        <div
          class="md:col-span-2 overflow-hidden rounded-md aspect-[4/3] shadow-xl">
          <BaseImg
            :src="coverUrl"
            :alt="page.title"
            class="w-full h-full object-cover" />
        </div>
        <div class="md:col-span-3 flex flex-col justify-center py-4">
          <h1
            class="text-3xl md:text-4xl font-bold text-highlighted leading-tight mb-4">
            {{ page.title }}
          </h1>
          <p v-if="page.excerpt" class="text-muted text-base mb-4 line-clamp-3">
            {{ page.excerpt }}
          </p>
          <p
            v-if="page.published_at"
            class="text-muted text-sm flex items-center gap-1">
            <UIcon name="i-tabler-calendar" class="size-3.5" />{{
              fmtDate(page.published_at)
            }}
          </p>
        </div>
      </div>

      <USeparator class="mb-8" />

      <div
        :class="
          effectiveSidebar
            ? 'grid grid-cols-1 lg:grid-cols-12 gap-8'
            : 'max-w-3xl mx-auto'
        ">
        <article :class="effectiveSidebar ? 'lg:col-span-8' : ''">
          <UCard>
            <MarkdownContent class="prose-lg" :html="contentHtml" />
          </UCard>
          <UCard class="mt-4">
            <PostCommentSection
              :object-id="page.id"
              object-type="post"
              :comment-meta="page.metas?.comment_open" />
          </UCard>
        </article>
        <BlogSidebar
          v-if="effectiveSidebar"
          option-key="page_sidebar_widgets"
          class="lg:col-span-4" />
      </div>
    </template>

    <!-- ╔═══════════════════════════════════════╗ -->
    <!-- ║  NONE LAYOUT  (无图)                   ║ -->
    <!-- ╚═══════════════════════════════════════╝ -->
    <template v-else>
      <div class="pt-8 mb-8">
        <div :class="effectiveSidebar ? '' : 'max-w-3xl mx-auto'">
          <h1
            class="text-4xl md:text-5xl font-bold text-highlighted leading-tight mb-4">
            {{ page.title }}
          </h1>
          <p v-if="page.excerpt" class="text-muted text-lg mb-4">
            {{ page.excerpt }}
          </p>
          <p
            v-if="page.published_at"
            class="text-muted text-sm flex items-center gap-1 pb-8 border-b border-default">
            <UIcon name="i-tabler-calendar" class="size-3.5" />{{
              fmtDate(page.published_at)
            }}
          </p>
        </div>
      </div>

      <div
        :class="
          effectiveSidebar
            ? 'grid grid-cols-1 lg:grid-cols-12 gap-8'
            : 'max-w-3xl mx-auto'
        ">
        <article :class="effectiveSidebar ? 'lg:col-span-8' : ''">
          <UCard>
            <MarkdownContent class="prose-lg" :html="contentHtml" />
          </UCard>
          <UCard class="mt-4">
            <PostCommentSection
              :object-id="page.id"
              object-type="post"
              :comment-meta="page.metas?.comment_open" />
          </UCard>
        </article>
        <BlogSidebar
          v-if="effectiveSidebar"
          option-key="page_sidebar_widgets"
          class="lg:col-span-4" />
      </div>
    </template>

    <!-- Plugin slot: page bottom -->
    <ClientOnly><ContributionSlot name="public:page-after-content" /></ClientOnly>
  </div>
</template>

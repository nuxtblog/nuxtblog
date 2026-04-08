<script setup lang="ts">
import { marked } from "marked";
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
  data: post,
  error,
  pending,
} = await useAsyncData(
  `post-${slug}`,
  () =>
    postApi.getPostBySlug(slug).catch((err) => {
      throw err;
    }),
  { lazy: false, server: true },
);

// ── Layout & sidebar ────────────────────────────────────────────────────────
const defaultLayout = computed(
  () =>
    optionsStore.get("post_default_layout", "hero") as "hero" | "half" | "none",
);
const defaultSidebar = computed(() =>
  optionsStore.getJSON<boolean>("post_sidebar_enabled", false),
);
const effectiveLayout = computed(() => {
  const m = post.value?.metas?.post_layout;
  if (m === "hero" || m === "half" || m === "none") return m;
  return defaultLayout.value;
});
const effectiveSidebar = computed(() => {
  const m = post.value?.metas?.post_sidebar;
  if (m === "1") return true;
  if (m === "0") return false;
  return defaultSidebar.value;
});

// ── Derived data ────────────────────────────────────────────────────────────
const coverUrl = computed(
  () => post.value?.featured_img?.url || defaultCover.value,
);
const categories = computed(
  () => post.value?.terms?.filter((t) => t.taxonomy === "category") ?? [],
);
const tags = computed(
  () => post.value?.terms?.filter((t) => t.taxonomy === "tag") ?? [],
);

// ── TOC + Downloads ──────────────────────────────────────────────────────────
const { toc, downloads, clearPost } = useCurrentPost();

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
  if (!post.value?.content) return "";
  const markdown = post.value.content;

  // Extract headings from markdown source
  const items: TocItem[] = [];
  const idCounts: Record<string, number> = {};
  for (const line of markdown.split("\n")) {
    const m = /^(#{2,4})\s+(.+)$/.exec(line.trim());
    if (!m) continue;
    const level = m[1]!.length;
    const text = m[2]!
      .trim()
      .replace(/`([^`]+)`/g, "$1")
      .replace(/\*+([^*]+)\*+/g, "$1")
      .replace(/\[([^\]]+)\]\([^)]+\)/g, "$1");
    const base = slugId(text);
    const n = (idCounts[base] = (idCounts[base] ?? 0) + 1);
    items.push({ id: n === 1 ? base : `${base}-${n}`, text, level });
  }
  toc.value = items;

  // Render markdown then inject IDs into heading tags
  let html = marked(markdown) as string;
  let ptr = 0;
  html = html.replace(/<h([234])>([\s\S]*?)<\/h\1>/g, (match, lvl) => {
    const item = items[ptr];
    if (ptr < items.length && item && item.level === parseInt(lvl)) {
      return match.replace(`<h${lvl}>`, `<h${lvl} id="${items[ptr++]!.id}">`);
    }
    return match;
  });
  return html;
});

// Populate downloads from post metas
watch(
  () => post.value?.metas?.post_downloads,
  (raw) => {
    try {
      downloads.value = raw ? JSON.parse(raw) : [];
    } catch {
      downloads.value = [];
    }
  },
  { immediate: true },
);

onUnmounted(() => clearPost());

const fmtDate = (d?: string) =>
  d
    ? new Date(d).toLocaleDateString("zh-CN", {
        year: "numeric",
        month: "long",
        day: "numeric",
      })
    : "";
const fmtNum = (n: number) =>
  n >= 10000
    ? (n / 10000).toFixed(1) + "w"
    : n >= 1000
      ? (n / 1000).toFixed(1) + "k"
      : String(n);

// ── SEO ─────────────────────────────────────────────────────────────────────
useHead({
  title: post.value?.title || t('site.post.detail_title'),
  meta: [
    { name: "description", content: post.value?.excerpt || "" },
    { property: "og:title", content: post.value?.title || "" },
    { property: "og:description", content: post.value?.excerpt || "" },
    { property: "og:image", content: coverUrl.value },
  ],
  bodyAttrs: {
    'data-page-id': post.value?.id ? String(post.value.id) : '',
    'data-page-slug': post.value?.slug || '',
  },
});

// ── View count ───────────────────────────────────────────────────────────────
onMounted(() => {
  if (post.value?.id)
    postApi.incrementViewCount?.(post.value.id).catch(() => {});
});
</script>

<template>
  <!-- Loading -->
  <div v-if="pending" class="space-y-6 py-8">
    <USkeleton class="w-full h-[55vh] rounded-md" />
    <div class="max-w-3xl mx-auto space-y-4">
      <USkeleton class="h-10 w-3/4" />
      <USkeleton class="h-5 w-1/2" />
      <USkeleton class="h-4 w-full" />
      <USkeleton class="h-4 w-5/6" />
    </div>
  </div>

  <!-- Error -->
  <UCard v-else-if="error" class="my-8">
    <div class="flex flex-col items-center py-20 text-center">
      <UIcon name="i-tabler-circle-x" class="size-12 text-error mb-4" />
      <p class="text-error text-lg mb-4">{{ $t('site.post.load_failed', { msg: error.message }) }}</p>
      <UButton color="primary" @click="navigateTo('/posts')"
        >{{ $t('site.post.back_to_articles') }}</UButton
      >
    </div>
  </UCard>

  <!-- Post -->
  <div v-else-if="post" class="-mt-8">
    <!-- Left action bar -->
    <PostActionBar
      :post-id="post.id"
      :like-count="post.like_count"
      :comment-count="post.comment_count" />

    <!-- ╔═══════════════════════════════════════╗ -->
    <!-- ║  HERO LAYOUT  (顶部大图)               ║ -->
    <!-- ╚═══════════════════════════════════════╝ -->
    <template v-if="effectiveLayout === 'hero'">
      <!-- Full-bleed hero cover -->
      <div class="relative h-[62vh] min-h-[420px] overflow-hidden">
        <BaseImg
          :src="coverUrl"
          :alt="post.title"
          class="w-full h-full object-cover" />
        <div
          class="absolute inset-0 bg-linear-to-t from-black/85 via-black/35 to-black/10" />

        <!-- Header content on image -->
        <div
          class="absolute bottom-0 left-0 right-0 px-4 sm:px-6 lg:px-8 pb-10">
          <div
            :class="
              effectiveSidebar
                ? 'max-w-none lg:max-w-[66%]'
                : 'max-w-3xl mx-auto'
            ">
            <div class="flex flex-wrap gap-2 mb-3">
              <UBadge
                v-for="cat in categories"
                :key="cat.id"
                color="primary"
                size="sm"
                >{{ cat.name }}</UBadge
              >
            </div>
            <h1
              class="text-3xl md:text-5xl font-bold text-white leading-tight mb-4 drop-shadow-lg">
              {{ post.title }}
            </h1>
            <p
              v-if="post.excerpt"
              class="text-white/70 text-base md:text-lg mb-5 line-clamp-2">
              {{ post.excerpt }}
            </p>
            <div
              class="flex flex-wrap items-center gap-4 text-sm text-white/60">
              <div v-if="post.author" class="flex items-center gap-2">
                <BaseAvatar
                  :src="post.author.avatar"
                  :alt="post.author.nickname"
                  size="sm" />
                <span class="text-white/80 font-medium">{{
                  post.author.nickname
                }}</span>
              </div>
              <span v-if="post.published_at" class="flex items-center gap-1">
                <UIcon name="i-tabler-calendar" class="size-3.5" />{{
                  fmtDate(post.published_at)
                }}
              </span>
              <span class="flex items-center gap-1">
                <UIcon name="i-tabler-eye" class="size-3.5" />{{
                  fmtNum(post.view_count)
                }}
              </span>
              <span class="flex items-center gap-1">
                <UIcon name="i-tabler-message-circle" class="size-3.5" />{{
                  post.comment_count
                }}
              </span>
              <ReportButton target-type="post" :target-id="post.id" />
            </div>
          </div>
        </div>
      </div>

      <!-- Body -->
      <div
        :class="
          effectiveSidebar
            ? 'grid grid-cols-1 lg:grid-cols-12 gap-8 pt-10'
            : 'max-w-3xl mx-auto pt-10'
        ">
        <article :class="effectiveSidebar ? 'lg:col-span-8' : ''">
          <UCard>
            <MarkdownContent class="prose-lg" :html="contentHtml" />
            <PostTags :tags="tags" class="mt-10" />
          </UCard>
          <UCard id="post-comments" class="mt-4">
            <PostCommentSection
              :object-id="post.id"
              object-type="post"
              :comment-meta="post.metas?.comment_open" />
          </UCard>
        </article>
        <BlogSidebar
          v-if="effectiveSidebar"
          option-key="blog_sidebar_widgets"
          class="lg:col-span-4" />
      </div>
    </template>

    <!-- ╔═══════════════════════════════════════╗ -->
    <!-- ║  HALF LAYOUT  (半图)                   ║ -->
    <!-- ╚═══════════════════════════════════════╝ -->
    <template v-else-if="effectiveLayout === 'half'">
      <!-- Outer wrapper: same grid logic as body so header + content align -->
      <div
        :class="
          effectiveSidebar
            ? 'grid grid-cols-1 lg:grid-cols-12 gap-8 pt-8'
            : 'max-w-4xl mx-auto pt-8'
        ">
        <!-- Main column (or full width when no sidebar) -->
        <div :class="effectiveSidebar ? 'lg:col-span-8' : ''">
          <!-- Split header -->
          <div class="grid grid-cols-1 md:grid-cols-5 gap-6 mb-8">
            <div
              class="md:col-span-2 overflow-hidden rounded-md aspect-[4/3] shadow-xl">
              <BaseImg
                :src="coverUrl"
                :alt="post.title"
                class="w-full h-full object-cover" />
            </div>
            <div class="md:col-span-3 flex flex-col justify-center py-2">
              <div class="flex flex-wrap gap-2 mb-4">
                <UBadge
                  v-for="cat in categories"
                  :key="cat.id"
                  color="primary">{{ cat.name }}</UBadge>
              </div>
              <h1
                class="text-3xl md:text-4xl font-bold text-highlighted leading-tight mb-4">
                {{ post.title }}
              </h1>
              <p
                v-if="post.excerpt"
                class="text-muted text-base mb-6 line-clamp-3">
                {{ post.excerpt }}
              </p>
              <div class="flex flex-wrap items-center gap-4 text-sm text-muted">
                <div v-if="post.author" class="flex items-center gap-2">
                  <BaseAvatar
                    :src="post.author.avatar"
                    :alt="post.author.nickname"
                    size="sm" />
                  <span class="text-highlighted font-medium">{{
                    post.author.nickname
                  }}</span>
                </div>
                <span v-if="post.published_at" class="flex items-center gap-1">
                  <UIcon name="i-tabler-calendar" class="size-3.5" />{{
                    fmtDate(post.published_at)
                  }}
                </span>
                <span class="flex items-center gap-1">
                  <UIcon name="i-tabler-eye" class="size-3.5" />{{
                    fmtNum(post.view_count)
                  }}
                </span>
                <span class="flex items-center gap-1">
                  <UIcon name="i-tabler-message-circle" class="size-3.5" />{{
                    post.comment_count
                  }}
                </span>
                <ReportButton target-type="post" :target-id="post.id" />
              </div>
            </div>
          </div>

          <USeparator class="mb-8" />

          <UCard>
            <MarkdownContent class="prose-lg" :html="contentHtml" />
            <PostTags :tags="tags" class="mt-10" />
          </UCard>
          <UCard id="post-comments" class="mt-4">
            <PostCommentSection
              :object-id="post.id"
              object-type="post"
              :comment-meta="post.metas?.comment_open" />
          </UCard>
        </div>

        <BlogSidebar
          v-if="effectiveSidebar"
          option-key="blog_sidebar_widgets"
          class="lg:col-span-4" />
      </div>
    </template>

    <!-- ╔═══════════════════════════════════════╗ -->
    <!-- ║  NONE LAYOUT  (无图)                   ║ -->
    <!-- ╚═══════════════════════════════════════╝ -->
    <template v-else>
      <div class="pt-8 mb-8">
        <div :class="effectiveSidebar ? '' : 'max-w-3xl mx-auto'">
          <div class="flex flex-wrap gap-2 mb-4">
            <UBadge v-for="cat in categories" :key="cat.id" color="primary">{{
              cat.name
            }}</UBadge>
          </div>
          <h1
            class="text-4xl md:text-5xl font-bold text-highlighted leading-tight mb-4">
            {{ post.title }}
          </h1>
          <p v-if="post.excerpt" class="text-muted text-lg mb-6">
            {{ post.excerpt }}
          </p>
          <div
            class="flex flex-wrap items-center gap-4 text-sm text-muted pb-8 border-b border-default">
            <div v-if="post.author" class="flex items-center gap-2">
              <BaseAvatar
                :src="post.author.avatar"
                :alt="post.author.nickname"
                size="xs" />
              <span class="text-highlighted font-medium">{{
                post.author.nickname
              }}</span>
            </div>
            <span v-if="post.published_at" class="flex items-center gap-1">
              <UIcon name="i-tabler-calendar" class="size-3.5" />{{
                fmtDate(post.published_at)
              }}
            </span>
            <span class="flex items-center gap-1">
              <UIcon name="i-tabler-eye" class="size-3.5" />{{
                fmtNum(post.view_count)
              }}
            </span>
            <span class="flex items-center gap-1">
              <UIcon name="i-tabler-message-circle" class="size-3.5" />{{
                post.comment_count
              }}
            </span>
            <ReportButton target-type="post" :target-id="post.id" />
          </div>
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
            <PostTags :tags="tags" class="mt-10" />
          </UCard>
          <UCard id="post-comments" class="mt-4">
            <PostCommentSection
              :object-id="post.id"
              object-type="post"
              :comment-meta="post.metas?.comment_open" />
          </UCard>
        </article>
        <BlogSidebar
          v-if="effectiveSidebar"
          option-key="blog_sidebar_widgets"
          class="lg:col-span-4" />
      </div>
    </template>
  </div>
</template>

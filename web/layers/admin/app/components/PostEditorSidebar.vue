<template>
  <div
    class="md:w-80 shrink-0 md:border-l border-default bg-default overflow-y-auto md:sticky md:top-0 md:self-start md:max-h-[100dvh]">
    <!-- Screen Options -->
    <div class="flex items-center justify-end px-4 pt-3 pb-1">
      <UPopover>
        <UButton
          variant="ghost"
          color="neutral"
          size="xs"
          icon="i-tabler-adjustments">
          {{ t("admin.posts.editor.screen_options") }}
        </UButton>
        <template #content>
          <div class="p-4 w-64 space-y-2">
            <p class="text-xs font-semibold text-muted uppercase mb-2">
              {{ t("admin.posts.editor.visible_sections") }}
            </p>
            <label
              v-for="id in filteredAllCards"
              :key="id"
              class="flex items-center gap-2 text-sm">
              <UCheckbox
                :model-value="isVisible(id)"
                @update:model-value="toggleVisibility(id)" />
              {{ cardLabels[id] }}
            </label>
            <UButton
              variant="link"
              size="xs"
              color="neutral"
              @click="resetDefaults">
              {{ t("common.reset") }}
            </UButton>
          </div>
        </template>
      </UPopover>
    </div>

    <!-- Draggable cards -->
    <ClientOnly>
    <VueDraggable
      :model-value="filteredVisibleCards"
      handle=".sidebar-drag-handle"
      :animation="200"
      class="p-3 space-y-3"
      @update:model-value="handleReorder">
      <div v-for="cardId in filteredVisibleCards" :key="cardId">
        <!-- Publish -->
        <SidebarCard
          v-if="cardId === 'publish'"
          :title="cardLabels.publish"
          :collapsed="isCollapsed('publish')"
          @toggle="toggleCollapsed('publish')">
          <div class="space-y-3">
            <UFormField :label="t('common.status')">
              <USelect
                v-model="formData.status"
                :items="[
                  { label: t('admin.posts.editor.status_draft'), value: 1 },
                  { label: t('admin.posts.editor.status_published'), value: 2 },
                  { label: t('admin.posts.editor.status_private'), value: 3 },
                  { label: t('admin.posts.editor.status_archived'), value: 4 },
                ]"
                class="w-full" />
            </UFormField>

            <UFormField :label="t('admin.posts.editor.published_at')">
              <input
                v-model="publishedAtLocal"
                type="datetime-local"
                class="w-full h-9 rounded-md border border-default bg-default px-3 text-sm text-highlighted focus:outline-none focus:ring-2 focus:ring-primary" />
              <p
                v-if="isScheduled"
                class="mt-1.5 flex items-center gap-1 text-xs text-primary">
                <UIcon name="i-tabler-clock" class="size-3 shrink-0" />
                {{ t("admin.posts.editor.scheduled_hint", { time: scheduledTimeLabel }) }}
              </p>
            </UFormField>

            <UFormField :label="t('common.type')">
              <USelect
                v-model="formData.post_type"
                :items="[
                  { label: t('admin.posts.editor.type_post'), value: 1 },
                  { label: t('admin.posts.editor.type_page'), value: 2 },
                  { label: t('admin.posts.editor.type_custom'), value: 3 },
                ]"
                class="w-full" />
            </UFormField>

            <UFormField :label="t('admin.posts.editor.comments')">
              <USelect
                v-model="formData.comment_status"
                :items="[
                  { label: t('admin.posts.editor.comments_open'), value: 1 },
                  { label: t('admin.posts.editor.comments_closed'), value: 0 },
                ]"
                class="w-full" />
            </UFormField>

            <UFormField :label="t('admin.posts.editor.password')">
              <form @submit.prevent>
                <input type="text" autocomplete="username" class="hidden" aria-hidden="true" tabindex="-1" />
                <UInput
                  v-model="formData.password"
                  type="password"
                  autocomplete="off"
                  :placeholder="t('admin.posts.editor.password_placeholder')"
                  class="w-full" />
              </form>
            </UFormField>

            <UFormField
              v-if="(authStore.user?.role ?? 0) >= 3"
              :label="t('admin.posts.editor.author_label')">
              <AdminSearchableSelect
                v-model="authorId"
                :items="authorOptions"
                :placeholder="t('admin.posts.editor.author_label')"
                :search-placeholder="t('admin.posts.search_authors')"
                trigger-class="w-full justify-between" />
            </UFormField>
          </div>
        </SidebarCard>

        <!-- Display (blog only) -->
        <SidebarCard
          v-else-if="cardId === 'display'"
          :title="cardLabels.display"
          :collapsed="isCollapsed('display')"
          @toggle="toggleCollapsed('display')">
          <div class="space-y-3">
            <div class="flex items-center justify-between">
              <div>
                <p class="text-sm text-default">{{ t("admin.posts.editor.banner") }}</p>
                <p class="text-xs text-muted">{{ t("admin.posts.editor.banner_hint") }}</p>
              </div>
              <USwitch v-model="isBanner" />
            </div>
            <div class="flex items-center justify-between">
              <div>
                <p class="text-sm text-default">{{ t("admin.posts.editor.featured") }}</p>
                <p class="text-xs text-muted">{{ t("admin.posts.editor.featured_hint") }}</p>
              </div>
              <USwitch v-model="isFeatured" />
            </div>
          </div>
        </SidebarCard>

        <!-- Layout -->
        <SidebarCard
          v-else-if="cardId === 'layout'"
          :title="cardLabels.layout"
          :collapsed="isCollapsed('layout')"
          @toggle="toggleCollapsed('layout')">
          <div class="space-y-3">
            <UFormField :label="t('admin.posts.editor.cover_layout')">
              <USelect
                v-model="postLayout"
                :items="[
                  { label: t('admin.posts.editor.layout_auto'), value: 'auto' },
                  { label: t('admin.posts.editor.layout_hero'), value: 'hero' },
                  { label: t('admin.posts.editor.layout_half'), value: 'half' },
                  { label: t('admin.posts.editor.layout_none'), value: 'none' },
                ]"
                class="w-full" />
              <p class="text-xs text-muted mt-1">{{ t("admin.posts.editor.cover_layout_hint") }}</p>
            </UFormField>
            <UFormField :label="t('admin.posts.editor.sidebar_label')">
              <USelect
                v-model="postSidebar"
                :items="[
                  { label: t('admin.posts.editor.sidebar_auto'), value: 'auto' },
                  { label: t('admin.posts.editor.sidebar_on'), value: '1' },
                  { label: t('admin.posts.editor.sidebar_off'), value: '0' },
                ]"
                class="w-full" />
            </UFormField>
          </div>
        </SidebarCard>

        <!-- Categories -->
        <PostCategoriesCard
          v-else-if="cardId === 'categories'"
          v-model:categories="selectedCategories"
          :title="cardLabels.categories"
          :collapsed="isCollapsed('categories')"
          :sidebar-loading="sidebarLoading"
          @toggle="toggleCollapsed('categories')"
        />

        <!-- Tags -->
        <PostTagsCard
          v-else-if="cardId === 'tags'"
          v-model:tags="selectedTags"
          :title="cardLabels.tags"
          :collapsed="isCollapsed('tags')"
          :sidebar-loading="sidebarLoading"
          @toggle="toggleCollapsed('tags')"
        />

        <!-- Downloads -->
        <PostDownloadsCard
          v-else-if="cardId === 'downloads'"
          v-model:downloads="postDownloads"
          :title="cardLabels.downloads"
          :collapsed="isCollapsed('downloads')"
          @toggle="toggleCollapsed('downloads')"
        />

        <!-- Featured Image -->
        <SidebarCard
          v-else-if="cardId === 'featured-image'"
          :title="cardLabels['featured-image']"
          :collapsed="isCollapsed('featured-image')"
          @toggle="toggleCollapsed('featured-image')">
          <FeaturedImagePicker
            v-model:imgId="formData.featured_img_id"
            v-model:imgUrl="featuredImageUrl" />
        </SidebarCard>

        <!-- Meta Fields -->
        <SidebarCard
          v-else-if="cardId === 'meta-fields'"
          :title="cardLabels['meta-fields']"
          :collapsed="isCollapsed('meta-fields')"
          @toggle="toggleCollapsed('meta-fields')">
          <template #header-actions>
            <UButton color="primary" variant="link" size="xs" @click.stop="addMetaField">
              {{ t("admin.posts.editor.add_meta") }}
            </UButton>
          </template>
          <div class="space-y-2">
            <div v-for="(meta, index) in metaFields" :key="index" class="flex gap-2">
              <UInput v-model="meta.key" :placeholder="t('admin.posts.editor.meta_key')" size="sm" class="flex-1" />
              <UInput v-model="meta.value" :placeholder="t('admin.posts.editor.meta_value')" size="sm" class="flex-1" />
              <UButton color="error" variant="ghost" icon="i-tabler-x" square size="sm" @click="metaFields.splice(index, 1)" />
            </div>
          </div>
        </SidebarCard>

        <!-- SEO -->
        <PostSEOCard
          v-else-if="cardId === 'seo'"
          v-model:seo="seoData"
          :title="cardLabels.seo"
          :collapsed="isCollapsed('seo')"
          @toggle="toggleCollapsed('seo')"
        />
      </div>
    </VueDraggable>
    </ClientOnly>
  </div>
</template>

<script setup lang="ts">
import { VueDraggable } from "vue-draggable-plus";
import type { CreatePostRequest } from "~/types/api/post";
import type { TermDetailResponse } from "~/types/api/term";

interface MetaField {
  key: string;
  value: string;
}
interface DownloadItem {
  name: string;
  url: string;
  size?: string;
  desc?: string;
}
interface SeoData {
  meta_title: string;
  meta_desc: string;
  og_title: string;
  og_image: string;
  canonical_url: string;
  robots: string;
}

const BLOG_ONLY_CARDS = new Set(["display", "categories", "tags", "downloads"]);

const props = defineProps<{
  simple?: boolean;
  sidebarLoading?: boolean;
}>();

// ── Two-way bindings with parent ──────────────────────────────────────────
const formData = defineModel<CreatePostRequest>("formData", { required: true });
const seoData = defineModel<SeoData>("seo", { required: true });
const selectedCategories = defineModel<number[]>("categories", { required: true });
const selectedTags = defineModel<TermDetailResponse[]>("tags", { required: true });
const metaFields = defineModel<MetaField[]>("metaFields", { required: true });
const postDownloads = defineModel<DownloadItem[]>("downloads", { required: true });
const featuredImageUrl = defineModel<string>("featuredImageUrl", { required: true });
const publishedAtLocal = defineModel<string>("publishedAt", { required: true });
const isBanner = defineModel<boolean>("isBanner", { required: true });
const isFeatured = defineModel<boolean>("isFeatured", { required: true });
const postLayout = defineModel<string>("layout", { required: true });
const postSidebar = defineModel<string>("sidebar", { required: true });
const authorId = defineModel<number>("authorId", { required: false });

// ── Stores & utils ────────────────────────────────────────────────────────
const { t } = useI18n();
const authStore = useAuthStore();
const userApi = useUserApi();

// ── Sidebar cards ─────────────────────────────────────────────────────────
const {
  allCards,
  visibleCards,
  isCollapsed,
  isVisible,
  toggleCollapsed,
  toggleVisibility,
  setOrder,
  resetDefaults,
  prefs,
} = useSidebarCards();

const cardLabels = computed<Record<string, string>>(() => ({
  publish: t("admin.posts.editor.publish_settings"),
  display: t("admin.posts.editor.display_settings"),
  layout: t("admin.posts.editor.layout_settings"),
  categories: t("admin.posts.editor.categories"),
  tags: t("admin.posts.editor.tags"),
  downloads: t("admin.posts.editor.downloads"),
  "featured-image": t("admin.posts.editor.featured_image"),
  "meta-fields": t("admin.posts.editor.meta_fields"),
  seo: t("admin.posts.editor.seo_settings"),
}));

const filteredVisibleCards = computed(() =>
  visibleCards.value.filter((id) => !props.simple || !BLOG_ONLY_CARDS.has(id)),
);

const filteredAllCards = computed(() =>
  allCards.value.filter((id) => !props.simple || !BLOG_ONLY_CARDS.has(id)),
);

function handleReorder(newVisible: string[]) {
  const hidden = new Set(prefs.value.hidden);
  const full = [...newVisible];
  for (const id of prefs.value.order) {
    if (hidden.has(id) && !full.includes(id)) full.push(id);
  }
  setOrder(full);
}

// ── Scheduled publish hint ────────────────────────────────────────────────
const isScheduled = computed(() => {
  if (!publishedAtLocal.value) return false;
  return new Date(publishedAtLocal.value).getTime() > Date.now() + 60_000;
});
const { locale } = useI18n();
const scheduledTimeLabel = computed(() => {
  const intlLocale = locale.value.startsWith("zh") ? "zh-CN" : "en-US";
  return new Date(publishedAtLocal.value).toLocaleString(intlLocale, {
    month: "numeric",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
});

// ── Author selector (admin only) ──────────────────────────────────────────
const authorsLoading = ref(false);
const authorOptions = ref<{ label: string; value: number }[]>([]);

if (import.meta.client) {
  onMounted(async () => {
    if ((authStore.user?.role ?? 0) < 3) return;
    authorsLoading.value = true;
    try {
      const res = await userApi.getUsers({ page: 1, size: 100 });
      authorOptions.value = (res?.list ?? []).map((u) => ({
        label: u.display_name || u.username,
        value: u.id,
      }));
    } catch {
    } finally {
      authorsLoading.value = false;
    }
  });
}

// ── Meta fields ───────────────────────────────────────────────────────────
const addMetaField = () => metaFields.value.push({ key: "", value: "" });
</script>

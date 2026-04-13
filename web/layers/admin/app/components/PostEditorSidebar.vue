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
                {{
                  t("admin.posts.editor.scheduled_hint", {
                    time: scheduledTimeLabel,
                  })
                }}
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
              <UInput
                v-model="formData.password"
                type="password"
                :placeholder="t('admin.posts.editor.password_placeholder')"
                class="w-full" />
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
                <p class="text-sm text-default">
                  {{ t("admin.posts.editor.banner") }}
                </p>
                <p class="text-xs text-muted">
                  {{ t("admin.posts.editor.banner_hint") }}
                </p>
              </div>
              <USwitch v-model="isBanner" />
            </div>
            <div class="flex items-center justify-between">
              <div>
                <p class="text-sm text-default">
                  {{ t("admin.posts.editor.featured") }}
                </p>
                <p class="text-xs text-muted">
                  {{ t("admin.posts.editor.featured_hint") }}
                </p>
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
              <p class="text-xs text-muted mt-1">
                {{ t("admin.posts.editor.cover_layout_hint") }}
              </p>
            </UFormField>

            <UFormField :label="t('admin.posts.editor.sidebar_label')">
              <USelect
                v-model="postSidebar"
                :items="[
                  {
                    label: t('admin.posts.editor.sidebar_auto'),
                    value: 'auto',
                  },
                  { label: t('admin.posts.editor.sidebar_on'), value: '1' },
                  { label: t('admin.posts.editor.sidebar_off'), value: '0' },
                ]"
                class="w-full" />
            </UFormField>
          </div>
        </SidebarCard>

        <!-- Categories -->
        <SidebarCard
          v-else-if="cardId === 'categories'"
          :title="cardLabels.categories"
          :collapsed="isCollapsed('categories')"
          @toggle="toggleCollapsed('categories')">
          <template #header-actions>
            <UButton
              color="primary"
              variant="link"
              size="xs"
              @click.stop="showAddCategoryModal = true">
              {{ t("admin.posts.editor.new_category_btn") }}
            </UButton>
          </template>
          <div v-if="sidebarLoading" class="space-y-2">
            <div v-for="i in 5" :key="i" class="flex items-center gap-2">
              <USkeleton class="size-4 rounded" />
              <USkeleton :class="`h-4 w-${i % 2 === 0 ? '24' : '32'}`" />
            </div>
          </div>
          <ParentCategoryMultiSelector v-else v-model="selectedCategories" />
        </SidebarCard>

        <!-- Tags -->
        <SidebarCard
          v-else-if="cardId === 'tags'"
          :title="cardLabels.tags"
          :collapsed="isCollapsed('tags')"
          @toggle="toggleCollapsed('tags')">
          <template #header-actions>
            <UButton
              color="primary"
              variant="link"
              size="xs"
              @click.stop="showAddTagModal = true">
              {{ t("admin.posts.editor.new_tag_btn") }}
            </UButton>
          </template>
          <template v-if="sidebarLoading">
            <div class="flex flex-wrap gap-2 mb-3">
              <USkeleton
                v-for="i in 3"
                :key="i"
                class="h-5 w-16 rounded-full" />
            </div>
            <div>
              <USkeleton class="h-3 w-16 mb-2" />
              <div class="flex flex-wrap gap-2">
                <USkeleton
                  v-for="i in 6"
                  :key="i"
                  class="h-6 w-14 rounded-md" />
              </div>
            </div>
          </template>
          <template v-else>
            <div class="flex flex-wrap gap-2 mb-3">
              <span
                v-for="tag in selectedTags"
                :key="tag.term_taxonomy_id"
                class="inline-flex items-center gap-1 px-2 py-0.5 text-xs rounded-full bg-primary/10 text-primary">
                {{ tag.name }}
                <button
                  type="button"
                  class="hover:text-primary/60 transition-colors"
                  @click="removeTag(tag.term_taxonomy_id)">
                  <UIcon name="i-tabler-x" class="size-3" />
                </button>
              </span>
            </div>
            <div v-if="availableTags.length > 0">
              <p class="text-xs text-muted mb-2">
                {{ t("admin.posts.editor.popular_tags") }}
              </p>
              <div class="flex flex-wrap gap-2">
                <UButton
                  v-for="tag in availableTags"
                  :key="tag.term_taxonomy_id"
                  color="neutral"
                  variant="outline"
                  size="xs"
                  @click="addExistingTag(tag)">
                  {{ tag.name }}
                </UButton>
              </div>
            </div>
          </template>
        </SidebarCard>

        <!-- Downloads -->
        <SidebarCard
          v-else-if="cardId === 'downloads'"
          :title="cardLabels.downloads"
          :collapsed="isCollapsed('downloads')"
          @toggle="toggleCollapsed('downloads')">
          <template #header-actions>
            <UButton
              color="primary"
              variant="link"
              size="xs"
              icon="i-tabler-plus"
              @click.stop="addDownload">
              {{ t("common.add") }}
            </UButton>
          </template>
          <div class="space-y-2">
            <div
              v-for="(dl, i) in postDownloads"
              :key="i"
              class="rounded-md border border-default p-3 space-y-2">
              <div class="flex items-center gap-2">
                <UInput
                  v-model="dl.name"
                  :placeholder="t('admin.posts.editor.download_name')"
                  size="xs"
                  class="flex-1" />
                <UButton
                  icon="i-tabler-trash"
                  color="error"
                  variant="ghost"
                  size="xs"
                  square
                  @click="postDownloads.splice(i, 1)" />
              </div>
              <UInput
                v-model="dl.url"
                :placeholder="t('admin.posts.editor.download_url')"
                size="xs"
                class="w-full" />
              <div class="flex gap-2">
                <UInput
                  v-model="dl.size"
                  :placeholder="t('admin.posts.editor.download_size')"
                  size="xs"
                  class="flex-1" />
                <UInput
                  v-model="dl.desc"
                  :placeholder="t('admin.posts.editor.download_desc')"
                  size="xs"
                  class="flex-1" />
              </div>
            </div>
            <p
              v-if="postDownloads.length === 0"
              class="text-xs text-muted text-center py-2">
              {{ t("admin.posts.editor.no_downloads") }}
            </p>
          </div>
        </SidebarCard>

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
            <UButton
              color="primary"
              variant="link"
              size="xs"
              @click.stop="addMetaField">
              {{ t("admin.posts.editor.add_meta") }}
            </UButton>
          </template>
          <div class="space-y-2">
            <div
              v-for="(meta, index) in metaFields"
              :key="index"
              class="flex gap-2">
              <UInput
                v-model="meta.key"
                :placeholder="t('admin.posts.editor.meta_key')"
                size="sm"
                class="flex-1" />
              <UInput
                v-model="meta.value"
                :placeholder="t('admin.posts.editor.meta_value')"
                size="sm"
                class="flex-1" />
              <UButton
                color="error"
                variant="ghost"
                icon="i-tabler-x"
                square
                size="sm"
                @click="metaFields.splice(index, 1)" />
            </div>
          </div>
        </SidebarCard>

        <!-- SEO -->
        <SidebarCard
          v-else-if="cardId === 'seo'"
          :title="cardLabels.seo"
          :collapsed="isCollapsed('seo')"
          @toggle="toggleCollapsed('seo')">
          <div class="space-y-3">
            <UFormField :label="t('admin.posts.editor.seo_title')">
              <UInput
                v-model="seoData.meta_title"
                :placeholder="t('admin.posts.editor.seo_title_placeholder')"
                class="w-full" />
            </UFormField>
            <UFormField :label="t('admin.posts.editor.seo_desc')">
              <UTextarea
                v-model="seoData.meta_desc"
                :rows="3"
                :placeholder="t('admin.posts.editor.seo_desc_placeholder')"
                class="w-full" />
            </UFormField>
            <UFormField :label="t('admin.posts.editor.og_title')">
              <UInput
                v-model="seoData.og_title"
                :placeholder="t('admin.posts.editor.og_title_placeholder')"
                class="w-full" />
            </UFormField>
            <UFormField :label="t('admin.posts.editor.og_image')">
              <UInput
                v-model="seoData.og_image"
                placeholder="https://..."
                class="w-full" />
            </UFormField>
            <UFormField label="Canonical URL">
              <UInput
                v-model="seoData.canonical_url"
                placeholder="https://..."
                class="w-full" />
            </UFormField>
            <UFormField label="Robots">
              <USelect
                v-model="seoData.robots"
                :items="[
                  { label: 'index, follow', value: 'index,follow' },
                  { label: 'noindex, follow', value: 'noindex,follow' },
                  { label: 'index, nofollow', value: 'index,nofollow' },
                  { label: 'noindex, nofollow', value: 'noindex,nofollow' },
                ]"
                class="w-full" />
            </UFormField>
          </div>
        </SidebarCard>
      </div>
    </VueDraggable>
  </div>

  <!-- 新建标签弹窗 -->
  <UModal
    v-model:open="showAddTagModal"
    :title="t('admin.posts.editor.new_tag_modal')">
    <template #content>
      <div class="p-6">
        <form class="space-y-4" @submit.prevent="handleAddTag">
          <UFormField :label="t('common.name')" required>
            <UInput
              v-model="newTag.name"
              required
              maxlength="100"
              :placeholder="t('admin.posts.editor.tag_name_placeholder')"
              class="w-full" />
          </UFormField>
          <UFormField
            :label="t('admin.posts.editor.slug_label')"
            :hint="t('admin.posts.editor.slug_hint')">
            <UInput
              v-model="newTag.slug"
              maxlength="100"
              :placeholder="t('admin.posts.editor.tag_slug_placeholder')"
              class="w-full" />
          </UFormField>
          <div class="flex gap-3 justify-end">
            <UButton
              color="neutral"
              variant="soft"
              type="button"
              @click="showAddTagModal = false">
              {{ t("common.cancel") }}
            </UButton>
            <UButton color="primary" type="submit" :loading="creatingTag">
              {{ t("common.create") }}
            </UButton>
          </div>
        </form>
      </div>
    </template>
  </UModal>

  <!-- 新建分类弹窗 -->
  <UModal
    v-model:open="showAddCategoryModal"
    :title="t('admin.posts.editor.new_category_modal')">
    <template #content>
      <div class="p-6">
        <form class="space-y-4" @submit.prevent="handleCreateCategory">
          <UFormField :label="t('common.name')" required>
            <UInput
              v-model="newCategory.name"
              required
              maxlength="100"
              :placeholder="t('admin.posts.editor.category_name_placeholder')"
              class="w-full" />
          </UFormField>
          <UFormField :label="t('admin.posts.editor.slug_label')">
            <UInput
              v-model="newCategory.slug"
              maxlength="100"
              :placeholder="t('admin.posts.editor.slug_auto')"
              class="w-full" />
          </UFormField>
          <ParentCategorySelector
            :label="t('admin.posts.categories.parent_category')" />
          <UFormField :label="t('common.description')">
            <UTextarea
              v-model="newCategory.description"
              :rows="3"
              maxlength="255"
              :placeholder="t('admin.posts.editor.category_desc_placeholder')"
              class="w-full" />
          </UFormField>
          <div class="flex gap-3 justify-end">
            <UButton
              color="neutral"
              variant="soft"
              type="button"
              @click="showAddCategoryModal = false">
              {{ t("common.cancel") }}
            </UButton>
            <UButton color="primary" type="submit" :loading="creatingCategory">
              {{ t("common.create") }}
            </UButton>
          </div>
        </form>
      </div>
    </template>
  </UModal>
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
const selectedCategories = defineModel<number[]>("categories", {
  required: true,
});
const selectedTags = defineModel<TermDetailResponse[]>("tags", {
  required: true,
});
const metaFields = defineModel<MetaField[]>("metaFields", { required: true });
const postDownloads = defineModel<DownloadItem[]>("downloads", {
  required: true,
});
const featuredImageUrl = defineModel<string>("featuredImageUrl", {
  required: true,
});
const publishedAtLocal = defineModel<string>("publishedAt", { required: true });
const isBanner = defineModel<boolean>("isBanner", { required: true });
const isFeatured = defineModel<boolean>("isFeatured", { required: true });
const postLayout = defineModel<string>("layout", { required: true });
const postSidebar = defineModel<string>("sidebar", { required: true });
const authorId = defineModel<number>("authorId", { required: false });

// ── Stores & utils ────────────────────────────────────────────────────────
const { t } = useI18n();
const toast = useToast();
const tagStore = useTagStore();
const categoryStore = useCategoryStore();
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

// ── Tags ──────────────────────────────────────────────────────────────────
const { tags } = storeToRefs(tagStore);
const availableTags = computed(() => {
  const ids = selectedTags.value.map((t) => t.term_taxonomy_id);
  return tags.value
    .filter((t) => !ids.includes(t.term_taxonomy_id))
    .slice(0, 10);
});

const showAddTagModal = ref(false);
const creatingTag = ref(false);
const newTag = ref({ name: "", slug: "" });

const addExistingTag = (tag: TermDetailResponse) => {
  if (
    !selectedTags.value.find((t) => t.term_taxonomy_id === tag.term_taxonomy_id)
  ) {
    selectedTags.value.push(tag);
  }
};
const removeTag = (id: number) => {
  selectedTags.value = selectedTags.value.filter(
    (t) => t.term_taxonomy_id !== id,
  );
};
const handleAddTag = async () => {
  if (!newTag.value.name.trim()) return;
  creatingTag.value = true;
  try {
    const created = await tagStore.addNewTag({
      name: newTag.value.name,
      slug: newTag.value.slug || undefined,
    });
    selectedTags.value.push(created);
    newTag.value = { name: "", slug: "" };
    showAddTagModal.value = false;
  } catch (err: any) {
    toast.add({
      title: t("admin.posts.editor.create_tag_failed"),
      description: err?.message,
      color: "error",
    });
  } finally {
    creatingTag.value = false;
  }
};

// ── Categories ────────────────────────────────────────────────────────────
const showAddCategoryModal = ref(false);
const creatingCategory = ref(false);
const newCategory = ref({
  name: "",
  slug: "",
  taxonomy: "category",
  description: "",
  parent_id: undefined as number | undefined,
});

const handleCreateCategory = async () => {
  if (!newCategory.value.name.trim()) return;
  creatingCategory.value = true;
  try {
    const created = await categoryStore.addNewCategory({
      name: newCategory.value.name,
      slug: newCategory.value.slug || undefined,
      description: newCategory.value.description || undefined,
      parent_id: newCategory.value.parent_id,
    });
    if (created) {
      selectedCategories.value.push(created.term_taxonomy_id);
      showAddCategoryModal.value = false;
      newCategory.value = {
        name: "",
        slug: "",
        taxonomy: "category",
        description: "",
        parent_id: undefined,
      };
    }
  } catch (error: any) {
    toast.add({
      title: t("admin.posts.editor.create_category_failed"),
      description: error?.message,
      color: "error",
    });
  } finally {
    creatingCategory.value = false;
  }
};

// ── Downloads ─────────────────────────────────────────────────────────────
const addDownload = () =>
  postDownloads.value.push({ name: "", url: "", size: "", desc: "" });

// ── Meta fields ───────────────────────────────────────────────────────────
const addMetaField = () => metaFields.value.push({ key: "", value: "" });
</script>

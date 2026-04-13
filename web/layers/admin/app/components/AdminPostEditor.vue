<template>
  <AdminPageContainer>
    <AdminPageHeader
      :title="
        mode === 'create'
          ? t('admin.posts.editor.create_title')
          : t('admin.posts.editor.edit_title')
      "
      :subtitle="
        mode === 'create'
          ? t('admin.posts.editor.create_subtitle')
          : t('admin.posts.editor.edit_subtitle')
      ">
      <template #actions>
        <div class="flex items-center gap-3">
          <UButton
            v-if="mode === 'edit' && init?.slug"
            color="neutral"
            variant="ghost"
            icon="i-tabler-eye"
            @click="
              navigateTo(`/posts/${init!.slug}`, { open: { target: '_blank' } })
            ">
            {{ t("admin.posts.editor.view_post") }}
          </UButton>
          <UButton
            v-if="mode === 'create'"
            color="neutral"
            variant="ghost"
            icon="i-tabler-eye"
            @click="showPreview = true">
            {{ t("admin.posts.editor.preview_content") }}
          </UButton>
          <UButton
            v-if="mode === 'edit'"
            color="neutral"
            variant="ghost"
            icon="i-tabler-history"
            @click="showRevisionModal = true">
            {{ t("admin.posts.editor.revisions") }}
          </UButton>
          <UButton
            color="neutral"
            variant="soft"
            :disabled="submitting"
            @click="handleSaveDraft">
            {{ t("admin.posts.editor.save_draft") }}
          </UButton>
          <UButton
            color="primary"
            :loading="submitting"
            :icon="isFuturePublish ? 'i-tabler-clock' : undefined"
            @click="handlePublish">
            {{
              isFuturePublish
                ? t("admin.posts.editor.schedule_btn")
                : mode === "create"
                  ? t("admin.posts.editor.publish")
                  : t("common.save")
            }}
          </UButton>
        </div>
      </template>
    </AdminPageHeader>

    <AdminPageContent class="p-0 flex flex-col md:flex-row">
      <!-- 主内容区 -->
      <div class="flex-1 min-w-0 max-w-5xl mx-auto">
        <div class="flex-1 bg-default px-2">
          <!-- 标题 -->
          <div class="pt-4 pb-3">
            <input
              v-model="formData.title"
              type="text"
              :placeholder="t('admin.posts.editor.title_placeholder')"
              :aria-label="t('admin.posts.editor.title_placeholder')"
              class="w-full text-3xl font-bold bg-transparent border-none outline-none placeholder:text-muted focus:border-b-2 focus:border-primary" />
          </div>

          <!-- 别名 -->
          <div class="pb-4">
            <input
              v-model="formData.slug"
              type="text"
              :placeholder="t('admin.posts.editor.slug_placeholder')"
              :aria-label="t('admin.posts.editor.slug_placeholder')"
              class="w-full text-sm bg-transparent border-b border-default pb-2 outline-none placeholder:text-muted focus:border-primary transition-colors" />
          </div>

          <!-- 编辑器 -->
          <AdminRichEditor
            ref="richEditorRef"
            v-model="formData.content"
            image-category="post"
            draft-key-prefix="blog:draft"
            :draft-entity-id="init?.id"
            :draft-mode="mode"
            :has-initial-content="!!init?.content"
            :placeholder="t('admin.posts.editor.content_placeholder')"
            enable-plugin-toolbar>
            <template #toolbar-extra="{ editor }">
              <ContributionSlot
                name="post-editor:toolbar"
                :ctx="{ editor }"
                class="contents"
                @command="
                  (cmdId: string) => handlePluginCommand(cmdId, editor)
                ">
                <template #menu="{ item, execute }">
                  <UTooltip :text="item.title ?? ''">
                    <UButton
                      variant="ghost"
                      color="neutral"
                      size="xs"
                      :icon="item.icon"
                      @click="execute" />
                  </UTooltip>
                </template>
              </ContributionSlot>
            </template>
            <template #footer>
              <ContributionSlot
                name="admin:post-editor-footer"
                :ctx="{ formData: formData }"
                class="empty:hidden" />
            </template>
          </AdminRichEditor>

          <!-- 摘要 -->
          <div class="pb-8">
            <UFormField :label="t('admin.posts.editor.excerpt')">
              <UTextarea
                v-model="formData.excerpt"
                :rows="3"
                :placeholder="t('admin.posts.editor.excerpt_placeholder')"
                class="w-full" />
            </UFormField>
          </div>
        </div>
      </div>

      <!-- 插件上下文插槽（摘要、标签推荐等） -->
      <ContributionSlot
        name="post-editor:context"
        :ctx="{ formData: formData }"
        class="empty:hidden" />

      <!-- 右侧边栏 -->
      <PostEditorSidebar
        :simple="props.simple"
        :sidebar-loading="sidebarLoading"
        v-model:formData="formData"
        v-model:seo="seoData"
        v-model:categories="selectedCategories"
        v-model:tags="selectedTags"
        v-model:metaFields="metaFields"
        v-model:downloads="postDownloads"
        v-model:featuredImageUrl="featuredImageUrl"
        v-model:publishedAt="publishedAtLocal"
        v-model:isBanner="isBanner"
        v-model:isFeatured="isFeatured"
        v-model:layout="postLayout"
        v-model:sidebar="postSidebar"
        v-model:authorId="formData.author_id" />
    </AdminPageContent>

    <!-- 版本历史弹窗 -->
    <PostEditorRevisionModal v-model="showRevisionModal" :post-id="init?.id" />

    <!-- 预览侧滑面板 -->
    <USlideover v-model:open="showPreview">
      <template #content>
        <div class="p-6 overflow-y-auto h-full">
          <h1 class="text-2xl font-bold mb-4">
            {{ formData.title || t("admin.posts.editor.untitled") }}
          </h1>
          <MarkdownContent :content="formData.content" />
          <div
            v-if="formData.excerpt"
            class="mt-6 pt-4 border-t border-default text-sm text-muted">
            {{ formData.excerpt }}
          </div>
        </div>
      </template>
    </USlideover>
  </AdminPageContainer>
</template>

<script setup lang="ts">
import type { Editor } from "@tiptap/vue-3";
import { dispatchCommand } from "~/composables/useNuxtblogAdmin";
import { usePluginContextStore } from "~/stores/plugin-context";
import type { TermDetailResponse } from "~/types/api/term";
import type { CreatePostRequest, UpdatePostRequest } from "~/types/api/post";

const { t } = useI18n();

// ── Props & Emits ─────────────────────────────────────────────────────────
export interface PostEditorInitialData {
  id?: number;
  title?: string;
  slug?: string;
  content?: string;
  excerpt?: string;
  status?: number;
  postType?: number;
  commentStatus?: number;
  featuredImgId?: number;
  featuredImgUrl?: string;
  publishedAt?: string;
  authorId?: number;
  categoryTaxonomyIds?: number[];
  selectedTagObjects?: TermDetailResponse[];
  isBanner?: boolean;
  isFeatured?: boolean;
  customMetaFields?: { key: string; value: string }[];
  seo?: {
    meta_title?: string;
    meta_desc?: string;
    og_title?: string;
    og_image?: string;
    canonical_url?: string;
    robots?: string;
  };
}

const props = defineProps<{
  mode: "create" | "edit";
  initialData?: PostEditorInitialData;
  submitting?: boolean;
  simple?: boolean;
}>();

const emit = defineEmits<{
  save: [payload: CreatePostRequest | UpdatePostRequest];
}>();

// ── Rich editor ref ───────────────────────────────────────────────────────
const richEditorRef = ref<any>(null);

// ── Plugin command dispatch ───────────────────────────────────────────────
const handlePluginCommand = (commandId: string, editor: Editor) => {
  const { state } = editor;
  const { from, to } = state.selection;
  const selectedText =
    from !== to ? state.doc.textBetween(from, to, " ") : null;

  dispatchCommand(commandId, {
    source: "editor",
    post: {
      title: formData.value.title,
      slug: formData.value.slug,
      content: formData.value.content ?? "",
      excerpt: formData.value.excerpt ?? "",
      status: formData.value.status === 2 ? "published" : "draft",
    },
    selection: selectedText,
    replace: (text: string) => {
      editor
        .chain()
        .focus()
        .deleteRange({ from, to })
        .insertContent(text)
        .run();
    },
    insert: (text: string) => {
      editor.chain().focus().insertContent(text).run();
    },
    setContent: (html: string) => {
      editor.commands.setContent(html);
    },
    setExcerpt: (text: string) => {
      formData.value.excerpt = text;
    },
    setSlug: (text: string) => {
      formData.value.slug = text;
    },
    addTags: async (tags: Array<string | { name: string; slug?: string }>) => {
      for (const item of tags) {
        const name = (typeof item === "string" ? item : item.name).trim();
        const slug = typeof item === "string" ? undefined : item.slug;
        if (!name) continue;
        if (
          selectedTags.value.some(
            (t) => t.name.toLowerCase() === name.toLowerCase(),
          )
        )
          continue;
        try {
          const tag = await tagStore.addNewTag({ name, slug });
          selectedTags.value.push(tag);
        } catch {
          // silently skip failed tags
        }
      }
    },
  });
};

// ── Form state ────────────────────────────────────────────────────────────
const categoryStore = useCategoryStore();
const tagStore = useTagStore();
const authStore = useAuthStore();

const init = props.initialData;

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

const sidebarLoading = ref(true);

const publishedAtLocal = ref(
  init?.publishedAt
    ? toDatetimeInputValue(init.publishedAt)
    : toDatetimeInputValue(new Date()),
);
const featuredImageUrl = ref(init?.featuredImgUrl ?? "");
const selectedCategories = ref<number[]>(init?.categoryTaxonomyIds ?? []);
const selectedTags = ref<TermDetailResponse[]>(init?.selectedTagObjects ?? []);
const metaFields = ref<MetaField[]>(init?.customMetaFields ?? []);
const isBanner = ref(init?.isBanner ?? false);
const isFeatured = ref(init?.isFeatured ?? false);
const postLayout = ref<string>(
  init?.customMetaFields?.find((f) => f.key === "post_layout")?.value || "auto",
);
const postSidebar = ref<string>(
  init?.customMetaFields?.find((f) => f.key === "post_sidebar")?.value ||
    "auto",
);

const parseDownloads = (): DownloadItem[] => {
  try {
    return JSON.parse(
      init?.customMetaFields?.find((f) => f.key === "post_downloads")?.value ??
        "[]",
    );
  } catch {
    return [];
  }
};
const postDownloads = ref<DownloadItem[]>(parseDownloads());

const formData = ref<CreatePostRequest>({
  post_type: init?.postType ?? 1,
  title: init?.title ?? "",
  slug: init?.slug ?? "",
  content: init?.content ?? "",
  excerpt: init?.excerpt ?? "",
  featured_img_id: init?.featuredImgId,
  status: init?.status ?? 1,
  published_at: undefined,
  password: "",
  comment_status: init?.commentStatus ?? 1,
  author_id: init?.authorId ?? authStore.user?.id ?? 1,
  term_taxonomy_ids: [],
});

const seoData = ref({
  meta_title: init?.seo?.meta_title ?? "",
  meta_desc: init?.seo?.meta_desc ?? "",
  og_title: init?.seo?.og_title ?? "",
  og_image: init?.seo?.og_image ?? "",
  canonical_url: init?.seo?.canonical_url ?? "",
  robots: init?.seo?.robots ?? "index,follow",
});

// ── Unsaved changes (tracked via rich editor + form watchers) ─────────────
const toast = useToast();
const hasUnsavedChanges = computed(
  () => richEditorRef.value?.hasUnsavedChanges ?? false,
);

watch(
  [
    formData,
    selectedCategories,
    selectedTags,
    isBanner,
    isFeatured,
    seoData,
    metaFields,
  ],
  () => {
    if (richEditorRef.value) richEditorRef.value.hasUnsavedChanges = true;
  },
  { deep: true },
);

// ── Revision modal & preview ──────────────────────────────────────────────
const showRevisionModal = ref(false);
const showPreview = ref(false);

// ── Build & save ──────────────────────────────────────────────────────────
const buildPayload = (
  status: number,
): CreatePostRequest | UpdatePostRequest => {
  formData.value.term_taxonomy_ids = [
    ...selectedCategories.value,
    ...selectedTags.value.map((t) => t.term_taxonomy_id),
  ];
  if (publishedAtLocal.value) {
    formData.value.published_at = new Date(
      publishedAtLocal.value,
    ).toISOString();
  }
  const metas: Record<string, string> = {};
  metaFields.value.forEach((field) => {
    if (field.key) metas[field.key] = field.value;
  });
  metas["is_banner"] = isBanner.value ? "1" : "";
  metas["is_featured"] = isFeatured.value ? "1" : "";
  if (postLayout.value && postLayout.value !== "auto")
    metas["post_layout"] = postLayout.value;
  if (postSidebar.value && postSidebar.value !== "auto")
    metas["post_sidebar"] = postSidebar.value;
  const validDownloads = postDownloads.value.filter((d) => d.name && d.url);
  if (validDownloads.length > 0)
    metas["post_downloads"] = JSON.stringify(validDownloads);
  else delete metas["post_downloads"];
  return {
    ...formData.value,
    post_type: formData.value.post_type ?? 1,
    status,
    metas,
  };
};

const isFormValid = computed(
  () =>
    formData.value.title.trim() !== "" &&
    (formData.value.content ?? "").trim() !== "",
);

// True when the selected publishedAt is more than 1 minute in the future.
const isFuturePublish = computed(() => {
  if (!publishedAtLocal.value) return false;
  return new Date(publishedAtLocal.value).getTime() > Date.now() + 60_000;
});

const handleSaveDraft = async () => {
  if (!isFormValid.value) {
    toast.add({
      title: t("admin.posts.editor.fill_required"),
      color: "warning",
    });
    return;
  }
  if (richEditorRef.value?.hasPendingUploads()) {
    toast.add({
      title: t("admin.posts.editor.pending_upload_warning"),
      color: "warning",
    });
    return;
  }
  await richEditorRef.value?.uploadPendingImages();
  emit("save", buildPayload(1));
};
const handlePublish = async () => {
  if (!isFormValid.value) {
    toast.add({
      title: t("admin.posts.editor.fill_required"),
      color: "warning",
    });
    return;
  }
  if (richEditorRef.value?.hasPendingUploads()) {
    toast.add({
      title: t("admin.posts.editor.pending_upload_warning"),
      color: "warning",
    });
    return;
  }
  await richEditorRef.value?.uploadPendingImages();
  // If publish time is in the future, save as draft — server cron will auto-publish at that time.
  emit("save", buildPayload(isFuturePublish.value ? 1 : 2));
};

const reset = () => {
  formData.value = {
    post_type: 1,
    title: "",
    slug: "",
    content: "",
    excerpt: "",
    featured_img_id: undefined,
    status: 1,
    published_at: undefined,
    password: "",
    comment_status: 1,
    author_id: authStore.user?.id ?? 1,
    term_taxonomy_ids: [],
  };
  selectedCategories.value = [];
  selectedTags.value = [];
  metaFields.value = [];
  isBanner.value = false;
  isFeatured.value = false;
  featuredImageUrl.value = "";
  publishedAtLocal.value = toDatetimeInputValue(new Date());
  seoData.value = {
    meta_title: "",
    meta_desc: "",
    og_title: "",
    og_image: "",
    canonical_url: "",
    robots: "index,follow",
  };
  if (richEditorRef.value) {
    richEditorRef.value.hasUnsavedChanges = false;
    richEditorRef.value.markSaved();
  }
};

const markSaved = () => {
  richEditorRef.value?.markSaved();
};

defineExpose({
  reset,
  isDirty: hasUnsavedChanges,
  getIsDirty: () => hasUnsavedChanges.value,
  markSaved,
  seoData,
});

// ── Plugin context: set post.type based on editor mode ────────────────────
const pluginCtx = usePluginContextStore();
pluginCtx.set('post.type', props.simple ? 'page' : 'post');

// ── Lifecycle ─────────────────────────────────────────────────────────────
onMounted(async () => {
  await Promise.all([categoryStore.loadCategories(), tagStore.loadTags()]);
  sidebarLoading.value = false;
});
</script>

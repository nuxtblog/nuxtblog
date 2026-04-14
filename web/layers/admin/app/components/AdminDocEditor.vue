<template>
  <AdminPageContainer>
    <AdminPageHeader
      :title="
        mode === 'create'
          ? t('admin.docs.editor.create_title')
          : t('admin.docs.editor.edit_title')
      "
      :subtitle="
        mode === 'create'
          ? t('admin.docs.editor.create_subtitle')
          : t('admin.docs.editor.edit_subtitle')
      ">
      <template #actions>
        <div class="flex items-center gap-3">
          <UButton
            v-if="mode === 'edit'"
            color="neutral"
            variant="ghost"
            icon="i-tabler-history"
            @click="showRevisionModal = true">
            {{ t("admin.docs.editor.revisions") }}
          </UButton>
          <UButton
            color="neutral"
            variant="soft"
            :disabled="submitting"
            @click="handleSaveDraft">
            {{ t("admin.docs.editor.save_draft") }}
          </UButton>
          <UButton color="primary" :loading="submitting" @click="handlePublish">
            {{
              mode === "create"
                ? t("admin.docs.editor.publish")
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
          <div class="px-8 sm:px-16 pt-4 pb-3">
            <input
              v-model="formData.title"
              type="text"
              :placeholder="t('admin.docs.editor.title_placeholder')"
              class="w-full text-3xl font-bold bg-transparent border-none outline-none placeholder:text-muted" />
          </div>

          <!-- 别名 -->
          <div class="px-8 sm:px-16 pb-4">
            <input
              v-model="formData.slug"
              type="text"
              :placeholder="t('admin.docs.editor.slug_placeholder')"
              class="w-full text-sm bg-transparent border-b border-default pb-2 outline-none placeholder:text-muted focus:border-primary transition-colors" />
          </div>

          <!-- 编辑器 -->
          <AdminRichEditor
            ref="richEditorRef"
            v-model="formData.content"
            image-category="doc"
            draft-key-prefix="blog:doc:draft"
            :draft-entity-id="init?.id"
            :draft-mode="mode"
            :has-initial-content="!!init?.content"
            :placeholder="t('admin.docs.editor.content_placeholder')"
            editor-class="px-8 sm:px-16"
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
          </AdminRichEditor>
        </div>
      </div>

      <!-- 右侧边栏 -->
      <div
        class="md:w-80 shrink-0 border-l border-default bg-default overflow-y-auto md:sticky md:top-0 md:self-start md:max-h-[100dvh]">
        <div class="p-4 space-y-6">
          <!-- 所属合集 -->
          <div class="space-y-3">
            <h4 class="text-sm font-semibold text-highlighted">
              {{ t("admin.docs.editor.collection_label") }}
            </h4>
            <AdminSearchableSelect
              v-model="formData.collection_id"
              :items="collectionOptions"
              :placeholder="t('admin.docs.editor.collection_label')"
              :search-placeholder="t('common.search')"
              trigger-class="w-full justify-between"
              @update:model-value="onCollectionChange" />
          </div>

          <!-- 父文档 -->
          <div class="space-y-3">
            <h4 class="text-sm font-semibold text-highlighted">
              {{ t("admin.docs.editor.parent_label") }}
            </h4>
            <AdminSearchableSelect
              v-model="formData.parent_id"
              :items="parentDocOptions"
              :placeholder="t('admin.docs.editor.no_parent')"
              :search-placeholder="t('common.search')"
              trigger-class="w-full justify-between"
              :disabled="!formData.collection_id" />
          </div>

          <!-- 发布设置 -->
          <div class="space-y-3">
            <h4 class="text-sm font-semibold text-highlighted">
              {{ t("admin.docs.editor.publish_settings") }}
            </h4>

            <UFormField :label="t('common.status')">
              <USelect
                v-model="formData.status"
                :items="statusOptions"
                class="w-full" />
            </UFormField>

            <UFormField :label="t('admin.docs.editor.published_at')">
              <UInput
                v-model="publishedAtLocal"
                type="datetime-local"
                class="w-full" />
            </UFormField>

            <UFormField :label="t('admin.docs.editor.comment_status')">
              <div class="flex items-center gap-2">
                <USwitch v-model="commentStatusBool" />
                <span class="text-sm text-muted">{{
                  formData.comment_status === 1
                    ? t("admin.docs.editor.comments_open")
                    : t("admin.docs.editor.comments_closed")
                }}</span>
              </div>
            </UFormField>

            <UFormField :label="t('admin.docs.editor.locale')">
              <UInput
                v-model="formData.locale"
                placeholder="zh"
                class="w-full" />
            </UFormField>

            <UFormField :label="t('admin.docs.editor.sort_order')">
              <UInput
                v-model.number="formData.sort_order"
                type="number"
                class="w-full" />
            </UFormField>
          </div>

          <!-- 摘要 -->
          <div class="space-y-3">
            <h4 class="text-sm font-semibold text-highlighted">
              {{ t("admin.posts.editor.excerpt") }}
            </h4>
            <UTextarea
              v-model="formData.excerpt"
              :rows="3"
              :placeholder="t('admin.posts.editor.excerpt_placeholder')"
              class="w-full" />
          </div>

          <!-- SEO 设置 -->
          <DocEditorSEOSettings v-model="seoData" />
        </div>
      </div>
    </AdminPageContent>

    <!-- 版本历史弹窗 -->
    <DocRevisionModal
      v-if="mode === 'edit' && init?.id"
      v-model:open="showRevisionModal"
      :doc-id="init.id"
      @restored="() => window.location.reload()" />
  </AdminPageContainer>
</template>

<script setup lang="ts">
import type { Editor } from "@tiptap/vue-3";
import { dispatchCommand } from "~/composables/useNuxtblogAdmin";
import { usePluginContextStore } from "~/stores/plugin-context";
import type {
  CreateDocRequest,
  UpdateDocRequest,
} from "~/types/api/doc";

const { t } = useI18n();

// ── Props & Emits ─────────────────────────────────────────────────────────
export interface DocEditorInitialData {
  id?: number;
  collectionId?: number;
  parentId?: number | null;
  title?: string;
  slug?: string;
  content?: string;
  excerpt?: string;
  status?: number;
  commentStatus?: number;
  locale?: string;
  sortOrder?: number;
  publishedAt?: string;
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
  initialData?: DocEditorInitialData;
  submitting?: boolean;
}>();

const emit = defineEmits<{
  save: [payload: CreateDocRequest | UpdateDocRequest];
}>();

// ── Rich editor ref ───────────────────────────────────────────────────────
const richEditorRef = ref<any>(null);
const toast = useToast();

// ── Plugin context: set post.type for doc editor ──────────────────────────
const pluginCtx = usePluginContextStore();
pluginCtx.set('post.type', 'doc');

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
    addTags: async () => {
      // Doc editor does not support tags — noop
    },
  });
};

// ── Form state ────────────────────────────────────────────────────────────
const docApi = useDocApi();

const init = props.initialData;

const DEFAULT_PUBLISHED_AT = () => new Date().toISOString().slice(0, 16);

const publishedAtLocal = ref(
  init?.publishedAt
    ? new Date(init.publishedAt).toISOString().slice(0, 16)
    : DEFAULT_PUBLISHED_AT(),
);

const formData = ref<{
  collection_id: number | undefined;
  parent_id: number | null | undefined;
  title: string;
  slug: string;
  content: string;
  excerpt: string;
  status: number;
  comment_status: number;
  locale: string;
  sort_order: number;
}>({
  collection_id: init?.collectionId,
  parent_id: init?.parentId ?? null,
  title: init?.title ?? "",
  slug: init?.slug ?? "",
  content: init?.content ?? "",
  excerpt: init?.excerpt ?? "",
  status: init?.status ?? 1,
  comment_status: init?.commentStatus ?? 1,
  locale: init?.locale ?? "zh",
  sort_order: init?.sortOrder ?? 0,
});

const seoData = ref({
  meta_title: init?.seo?.meta_title ?? "",
  meta_desc: init?.seo?.meta_desc ?? "",
  og_title: init?.seo?.og_title ?? "",
  og_image: init?.seo?.og_image ?? "",
  canonical_url: init?.seo?.canonical_url ?? "",
  robots: init?.seo?.robots ?? "index,follow",
});

const commentStatusBool = computed({
  get: () => formData.value.comment_status === 1,
  set: (val: boolean) => {
    formData.value.comment_status = val ? 1 : 0;
  },
});

// ── Status/options ────────────────────────────────────────────────────────
const statusOptions = computed(() => [
  { label: t("admin.docs.editor.status_draft"), value: 1 },
  { label: t("admin.docs.editor.status_published"), value: 2 },
  { label: t("admin.docs.editor.status_archived"), value: 3 },
]);

// ── Collections ───────────────────────────────────────────────────────────
const collections = ref<{ id: number; title: string; slug: string }[]>([]);
const parentDocs = ref<{ id: number; title: string }[]>([]);

const collectionOptions = computed(() => [
  { label: t("admin.docs.all_collections"), value: undefined },
  ...collections.value.map((c) => ({ label: c.title, value: c.id })),
]);

const parentDocOptions = computed(() => [
  { label: t("admin.docs.editor.no_parent"), value: null },
  ...parentDocs.value
    .filter((d) => d.id !== init?.id)
    .map((d) => ({ label: d.title, value: d.id })),
]);

async function fetchCollections() {
  try {
    const res = await docApi.getCollections({ page_size: 100, status: 2 });
    collections.value = res.data ?? [];
  } catch {}
}

async function fetchParentDocs(collectionId: number) {
  try {
    const res = await docApi.getDocs({
      collection_id: collectionId,
      page_size: 100,
    });
    parentDocs.value = res.data ?? [];
  } catch {}
}

function onCollectionChange(val: number | undefined) {
  formData.value.parent_id = null;
  if (val) fetchParentDocs(val);
  else parentDocs.value = [];
}

// ── Auto slug from title ──────────────────────────────────────────────────
watch(
  () => formData.value.title,
  (val) => {
    if (!formData.value.slug) {
      formData.value.slug = val
        .toLowerCase()
        .trim()
        .replace(/[^\w\s-]/g, "")
        .replace(/\s+/g, "-")
        .replace(/-+/g, "-");
    }
  },
);

// ── Unsaved changes (tracked via rich editor) ────────────────────────────
const hasUnsavedChanges = computed(
  () => richEditorRef.value?.hasUnsavedChanges ?? false,
);

watch(
  [formData, seoData],
  () => {
    if (richEditorRef.value) richEditorRef.value.hasUnsavedChanges = true;
  },
  { deep: true },
);

const markSaved = () => {
  richEditorRef.value?.markSaved();
};

// ── Revision modal ────────────────────────────────────────────────────────
const showRevisionModal = ref(false);

// ── Build & save ──────────────────────────────────────────────────────────
const isFormValid = computed(
  () => formData.value.title.trim() !== "" && !!formData.value.collection_id,
);

const buildPayload = (status: number): CreateDocRequest | UpdateDocRequest => {
  const base = {
    collection_id: formData.value.collection_id as number,
    parent_id: formData.value.parent_id ?? undefined,
    title: formData.value.title,
    slug: formData.value.slug || formData.value.title,
    content: formData.value.content,
    excerpt: formData.value.excerpt || undefined,
    status: status as 1 | 2 | 3,
    comment_status: formData.value.comment_status,
    locale: formData.value.locale,
    sort_order: formData.value.sort_order,
    published_at: publishedAtLocal.value
      ? new Date(publishedAtLocal.value).toISOString()
      : undefined,
  };
  return base;
};

const handleSaveDraft = async () => {
  if (!isFormValid.value) {
    toast.add({
      title: t("admin.posts.editor.fill_required"),
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
  await richEditorRef.value?.uploadPendingImages();
  const status = props.mode === "create" ? 2 : formData.value.status;
  emit("save", buildPayload(status));
};

defineExpose({
  isDirty: hasUnsavedChanges,
  getIsDirty: () => hasUnsavedChanges.value,
  markSaved,
  seoData,
});

// ── Lifecycle ─────────────────────────────────────────────────────────────
onMounted(async () => {
  await fetchCollections();

  if (init?.collectionId) {
    await fetchParentDocs(init.collectionId);
  }

  // Revision history is loaded by DocRevisionModal on open
});
</script>

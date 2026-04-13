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
            editor-class="px-8 sm:px-16" />
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
          <details class="space-y-3">
            <summary
              class="text-sm font-semibold text-highlighted cursor-pointer select-none flex items-center gap-2 py-1">
              <UIcon name="i-tabler-search" class="size-4" />
              {{ t("admin.docs.editor.seo_settings") }}
            </summary>
            <div class="space-y-3 pt-2">
              <UFormField label="Meta Title">
                <UInput
                  v-model="seoData.meta_title"
                  placeholder="留空使用文档标题"
                  class="w-full" />
              </UFormField>
              <UFormField label="Meta Description">
                <UTextarea
                  v-model="seoData.meta_desc"
                  placeholder="页面描述摘要"
                  :rows="3"
                  class="w-full" />
              </UFormField>
              <UFormField label="OG Title">
                <UInput
                  v-model="seoData.og_title"
                  placeholder="社交分享标题"
                  class="w-full" />
              </UFormField>
              <UFormField label="OG Image">
                <UInput
                  v-model="seoData.og_image"
                  placeholder="社交分享图片 URL"
                  class="w-full" />
              </UFormField>
              <UFormField label="Canonical URL">
                <UInput
                  v-model="seoData.canonical_url"
                  placeholder="规范链接（可选）"
                  class="w-full" />
              </UFormField>
              <UFormField label="Robots">
                <UInput
                  v-model="seoData.robots"
                  placeholder="index,follow"
                  class="w-full" />
              </UFormField>
            </div>
          </details>
        </div>
      </div>
    </AdminPageContent>

    <!-- 版本历史弹窗 -->
    <UModal
      v-if="mode === 'edit'"
      v-model:open="showRevisionModal"
      :ui="{ content: 'max-w-2xl' }">
      <template #content>
        <div class="p-6">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-lg font-semibold text-highlighted">
              {{ t("admin.docs.editor.revisions") }}
            </h3>
            <UButton
              icon="i-tabler-x"
              color="neutral"
              variant="ghost"
              size="sm"
              square
              @click="showRevisionModal = false" />
          </div>

          <div
            v-if="revisions.length === 0"
            class="flex flex-col items-center justify-center py-12">
            <UIcon name="i-tabler-history" class="size-12 text-muted mb-2" />
            <p class="text-sm text-muted">暂无修订历史</p>
          </div>

          <div v-else class="space-y-2 max-h-96 overflow-y-auto">
            <div
              v-for="rev in revisions"
              :key="rev.id"
              class="flex items-center gap-3 p-3 border border-default rounded-lg group hover:bg-elevated transition-colors">
              <div class="flex-1 min-w-0">
                <p class="text-sm font-medium text-highlighted truncate">
                  {{ rev.title || "(无标题)" }}
                </p>
                <div class="flex items-center gap-2 mt-0.5">
                  <span class="text-xs text-muted">{{
                    new Date(rev.created_at).toLocaleString("zh-CN")
                  }}</span>
                  <span v-if="rev.rev_note" class="text-xs text-muted"
                    >· {{ rev.rev_note }}</span
                  >
                </div>
              </div>
              <UButton
                size="xs"
                color="neutral"
                variant="outline"
                icon="i-tabler-restore"
                class="opacity-0 group-hover:opacity-100 transition-opacity"
                :loading="restoringRevisionId === rev.id"
                @click="handleRestoreRevision(rev)">
                {{ t("admin.docs.restore_revision") }}
              </UButton>
            </div>
          </div>
        </div>
      </template>
    </UModal>
  </AdminPageContainer>
</template>

<script setup lang="ts">
import type {
  CreateDocRequest,
  UpdateDocRequest,
  DocRevisionItem,
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
const revisions = ref<DocRevisionItem[]>([]);
const restoringRevisionId = ref<number | null>(null);

async function loadRevisions() {
  if (!init?.id) return;
  try {
    const res = await docApi.getRevisions(init.id);
    revisions.value = res.list ?? [];
  } catch {}
}

async function handleRestoreRevision(rev: DocRevisionItem) {
  if (!init?.id) return;
  restoringRevisionId.value = rev.id;
  try {
    await docApi.restoreRevision(init.id, rev.id);
    toast.add({ title: t("admin.docs.restore_revision"), color: "success" });
    showRevisionModal.value = false;
    // Reload page to reflect restored content
    window.location.reload();
  } catch (e: any) {
    toast.add({
      title: e?.message ?? t("common.operation_failed"),
      color: "error",
    });
  } finally {
    restoringRevisionId.value = null;
  }
}

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

  if (props.mode === "edit" && init?.id) {
    loadRevisions();
  }
});
</script>

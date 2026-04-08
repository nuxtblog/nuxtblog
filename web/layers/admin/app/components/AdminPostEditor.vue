<template>
  <AdminPageContainer>
    <AdminPageHeader
      :title="mode === 'create' ? t('admin.posts.editor.create_title') : t('admin.posts.editor.edit_title')"
      :subtitle="mode === 'create' ? t('admin.posts.editor.create_subtitle') : t('admin.posts.editor.edit_subtitle')">
      <template #actions>
        <div class="flex items-center gap-3">
          <UButton
            v-if="mode === 'edit'"
            color="neutral"
            variant="ghost"
            icon="i-tabler-history"
            @click="showRevisionModal = true">
            {{ t("admin.posts.editor.revisions") }}
          </UButton>
          <UButton color="neutral" variant="soft" :disabled="submitting" @click="handleSaveDraft">
            {{ t("admin.posts.editor.save_draft") }}
          </UButton>
          <UButton color="primary" :loading="submitting" :icon="isFuturePublish ? 'i-tabler-clock' : undefined" @click="handlePublish">
            {{ isFuturePublish
              ? t("admin.posts.editor.schedule_btn")
              : mode === "create" ? t("admin.posts.editor.publish") : t("common.save") }}
          </UButton>
        </div>
      </template>
    </AdminPageHeader>

    <AdminPageContent class="p-0 flex flex-col md:flex-row">
      <!-- 主内容区 -->
      <div class="flex-1 min-w-0 max-w-5xl mx-auto">
        <div class="flex-1 overflow-y-auto bg-default">

          <!-- 草稿恢复提示 -->
          <div v-if="showDraftRestore" class="px-8 sm:px-16 pt-8 pb-3">
            <UAlert
              icon="i-tabler-device-floppy"
              color="warning"
              variant="subtle"
              :title="t('admin.posts.editor.draft_found')">
              <template #description>
                <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-2">
                  <span class="text-sm text-gray-500">
                    {{ savedDraft?.savedAt ? t("admin.posts.editor.draft_saved_at", { time: new Date(savedDraft.savedAt).toLocaleString() }) : "" }}
                  </span>
                  <div class="flex gap-2">
                    <UButton size="xs" color="primary"  variant="soft"  @click="restoreDraft">{{ t("admin.posts.editor.restore_draft") }}</UButton>
                    <UButton size="xs" color="neutral"  variant="ghost" @click="discardDraft">{{ t("admin.posts.editor.discard") }}</UButton>
                  </div>
                </div>
              </template>
            </UAlert>
          </div>

          <!-- 标题 -->
          <div class="px-8 sm:px-16 pt-8 pb-3">
            <input
              v-model="formData.title"
              type="text"
              :placeholder="t('admin.posts.editor.title_placeholder')"
              class="w-full text-3xl font-bold bg-transparent border-none outline-none placeholder:text-muted" />
          </div>

          <!-- 别名 -->
          <div class="px-8 sm:px-16 pb-4">
            <input
              v-model="formData.slug"
              type="text"
              :placeholder="t('admin.posts.editor.slug_placeholder')"
              class="w-full text-sm bg-transparent border-b border-default pb-2 outline-none placeholder:text-muted focus:border-primary transition-colors" />
          </div>

          <!-- 编辑器骨架屏 -->
          <div v-if="editorLoading" class="px-8 sm:px-16 py-4">
            <div class="flex items-center gap-2 border-b border-default pb-2 mb-6">
              <USkeleton v-for="i in 8" :key="i" class="h-7 w-7 rounded" />
              <USkeleton class="h-7 w-px mx-1" />
              <USkeleton v-for="i in 5" :key="'b' + i" class="h-7 w-7 rounded" />
            </div>
            <div class="space-y-4 min-h-[500px]">
              <USkeleton class="h-8 w-2/3" />
              <USkeleton class="h-4 w-full" />
              <USkeleton class="h-4 w-5/6" />
              <USkeleton class="h-4 w-4/5" />
              <USkeleton class="h-4 w-full" />
              <USkeleton class="h-4 w-3/4" />
              <div class="pt-4 space-y-3">
                <USkeleton class="h-6 w-48" />
                <USkeleton class="h-4 w-full" />
                <USkeleton class="h-4 w-5/6" />
                <USkeleton class="h-4 w-full" />
              </div>
            </div>
          </div>

          <!-- 编辑器 -->
          <UEditor
            v-else
            ref="editorRef"
            v-slot="{ editor, handlers }"
            v-model="formData.content"
            content-type="markdown"
            :placeholder="t('admin.posts.editor.content_placeholder')"
            :extensions="editorExtensions"
            :handlers="editorHandlers"
            :ui="{ base: 'px-8 sm:px-16 py-6' }"
            class="min-h-[500px] px-8 sm:px-16 pb-4">
            <div class="border-b border-default sticky top-0 inset-x-0 z-10 bg-default/95 backdrop-blur">
              <UEditorToolbar
                :editor="editor"
                :items="toolbarItems"
                layout="fixed"
                class="px-4 py-1.5 overflow-x-auto" />
              <ContributionSlot name="post-editor:toolbar" :ctx="{ editor }" class="px-4 py-1 empty:hidden" />
            </div>

            <UEditorToolbar
              :editor="editor"
              :items="bubbleItems"
              class="z-50"
              layout="bubble"
              :should-show="({ editor: e, view, state }) => {
                const { selection } = state;
                return view.hasFocus() && !selection.empty && !e.isActive('image');
              }" />

            <UEditorDragHandle
              v-slot="{ ui, onClick }"
              :editor="editor"
              @node-change="selectedNode = $event">
              <UButton
                icon="i-tabler-plus"
                color="neutral"
                variant="ghost"
                size="sm"
                :class="ui.handle()"
                @click="(e) => {
                  e.stopPropagation();
                  const selected = onClick();
                  handlers.suggestion?.execute(editor, { pos: selected?.pos }).run();
                }" />
              <UDropdownMenu
                v-slot="{ open }"
                :modal="false"
                :items="dragHandleItems(editor)"
                :content="{ side: 'left' }"
                :ui="{ content: 'w-52', label: 'text-xs' }"
                @update:open="editor.chain().setMeta('lockDragHandle', $event).run()">
                <UButton
                  color="neutral"
                  variant="ghost"
                  active-variant="soft"
                  size="sm"
                  icon="i-tabler-grip-vertical"
                  :active="open"
                  :class="ui.handle()" />
              </UDropdownMenu>
            </UEditorDragHandle>

            <UEditorSuggestionMenu :editor="editor" :items="suggestionItems" :append-to="appendToBody" />
            <UEditorMentionMenu   :editor="editor" :items="[]"             :append-to="appendToBody" />
            <UEditorEmojiMenu     :editor="editor" :items="emojiItems"     :append-to="appendToBody" />
          </UEditor>

          <!-- 字数统计 -->
          <div class="px-8 sm:px-16 py-2 flex items-center gap-4 text-xs text-muted border-t border-default">
            <span>{{ t("admin.posts.editor.char_count",      { n: charCount }) }}</span>
            <span>{{ t("admin.posts.editor.reading_minutes", { n: readingMinutes }) }}</span>
            <span v-if="autoSavedLabel" class="ml-auto flex items-center gap-1">
              <UIcon name="i-tabler-circle-check" class="size-3 text-success" />{{ autoSavedLabel }}
            </span>
          </div>

          <!-- 摘要 -->
          <div class="px-8 sm:px-16 pb-8">
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
      <ContributionSlot name="post-editor:context" :ctx="{ formData: formData }" class="empty:hidden" />

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

    <!-- 图片上传（隐藏） -->
    <input
      ref="imageFileInput"
      type="file"
      accept="image/*"
      class="hidden"
      @change="handleImageFileSelect" />

  </AdminPageContainer>
</template>

<script setup lang="ts">
import type { EditorToolbarItem } from "@nuxt/ui";
import type { JSONContent, Editor } from "@tiptap/vue-3";
import { Emoji, gitHubEmojis } from "@tiptap/extension-emoji";
import { TextAlign } from "@tiptap/extension-text-align";
import { Markdown } from "tiptap-markdown";
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
    meta_title?: string; meta_desc?: string; og_title?: string;
    og_image?: string;   canonical_url?: string; robots?: string;
  };
}

const props = defineProps<{
  mode: "create" | "edit";
  initialData?: PostEditorInitialData;
  submitting?: boolean;
  simple?: boolean;
}>();

const emit = defineEmits<{ save: [payload: CreatePostRequest | UpdatePostRequest] }>();

// ── Editor extensions & emoji ─────────────────────────────────────────────
const editorExtensions = [
  Emoji,
  TextAlign.configure({ types: ["heading", "paragraph"] }),
  Markdown.configure({ transformPastedText: true, transformCopiedText: true }),
];
const appendToBody = import.meta.client ? () => document.body : undefined;
const emojiItems   = gitHubEmojis.filter((e) => !e.name.startsWith("regional_indicator_"));

// ── Toolbar config (composable) ───────────────────────────────────────────
const { toolbarItems, bubbleItems, suggestionItems, selectedNode, dragHandleItems } =
  usePostEditorToolbar();

// ── Image upload ──────────────────────────────────────────────────────────
// Strategy: insert a compressed data URL immediately (persisted in auto-save /
// localStorage), then upload all data-URL images to the server right before
// the post is saved/published and replace them with CDN URLs.
const toast          = useToast();
const mediaStore     = useMediaStore();
const imageFileInput = ref<HTMLInputElement | null>(null);
let   pendingEditor: Editor | null = null;

// Compress + resize to JPEG so the data URL stays small enough for localStorage
const compressImage = (file: File, maxWidth = 1920, quality = 0.85): Promise<string> =>
  new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = (e) => {
      const img = new Image();
      img.onload = () => {
        const canvas = document.createElement("canvas");
        const scale  = Math.min(1, maxWidth / img.width);
        canvas.width  = Math.round(img.width  * scale);
        canvas.height = Math.round(img.height * scale);
        canvas.getContext("2d")!.drawImage(img, 0, 0, canvas.width, canvas.height);
        resolve(canvas.toDataURL("image/jpeg", quality));
      };
      img.onerror = reject;
      img.src = e.target?.result as string;
    };
    reader.onerror = reject;
    reader.readAsDataURL(file);
  });

const editorHandlers = computed(() => ({
  image: {
    canExecute: (editor: Editor) => editor.isEditable,
    execute: (editor: Editor) => {
      pendingEditor = editor;
      nextTick(() => imageFileInput.value?.click());
      return editor.chain();
    },
    isActive: (_editor: Editor) => false,
  },
}));

const handleImageFileSelect = async (event: Event) => {
  const file = (event.target as HTMLInputElement).files?.[0];
  if (!file || !pendingEditor) return;
  const editor = pendingEditor;
  pendingEditor = null;
  try {
    const dataUrl = await compressImage(file);
    editor.chain().focus().setImage({ src: dataUrl, alt: file.name.replace(/\.[^.]+$/, "") }).run();
  } catch {
    toast.add({ title: t("admin.posts.editor.image_upload_failed"), color: "error" });
  } finally {
    if (imageFileInput.value) imageFileInput.value.value = "";
  }
};

// Upload every data-URL image embedded in the markdown content, replace with CDN URL
const dataUrlToFile = (dataUrl: string, name: string): File => {
  const [header, b64] = dataUrl.split(",");
  const mime = header.match(/:(.*?);/)?.[1] ?? "image/jpeg";
  const ext  = mime.split("/")[1] ?? "jpg";
  const bin  = atob(b64);
  const arr  = new Uint8Array(bin.length);
  for (let i = 0; i < bin.length; i++) arr[i] = bin.charCodeAt(i);
  return new File([arr], name.includes(".") ? name : `${name}.${ext}`, { type: mime });
};

const uploadPendingImages = async (): Promise<void> => {
  const content = formData.value.content ?? "";
  const regex   = /!\[([^\]]*)\]\((data:[^)]+)\)/g;
  const matches: { alt: string; src: string }[] = [];
  let m: RegExpExecArray | null;
  while ((m = regex.exec(content)) !== null) matches.push({ alt: m[1], src: m[2] });
  if (!matches.length) return;

  toast.add({ title: t("admin.posts.editor.image_uploading"), color: "neutral", duration: 0, id: "img-upload" });
  try {
    const results = await Promise.all(
      matches.map(async ({ alt, src }) => {
        const name   = alt || `image-${Date.now()}`;
        const result = await mediaStore.uploadMedia(dataUrlToFile(src, name), { title: name, category: "post" });
        return { src, cdnUrl: result?.cdn_url };
      }),
    );
    let updated = content;
    for (const { src, cdnUrl } of results) {
      if (cdnUrl) updated = updated.replaceAll(src, cdnUrl);
    }
    formData.value.content = updated;
    await nextTick();
  } finally {
    toast.remove("img-upload");
  }
};

// ── Form state ────────────────────────────────────────────────────────────
const optionsStore  = useOptionsStore();
const categoryStore = useCategoryStore();
const tagStore      = useTagStore();
const authStore     = useAuthStore();

const init = props.initialData;

interface MetaField    { key: string; value: string }
interface DownloadItem { name: string; url: string; size?: string; desc?: string }

const editorLoading  = ref(true);
const sidebarLoading = ref(true);

const publishedAtLocal = ref(
  init?.publishedAt ? toDatetimeInputValue(init.publishedAt) : toDatetimeInputValue(new Date()),
);
const featuredImageUrl  = ref(init?.featuredImgUrl ?? "");
const selectedCategories = ref<number[]>(init?.categoryTaxonomyIds ?? []);
const selectedTags      = ref<TermDetailResponse[]>(init?.selectedTagObjects ?? []);
const metaFields        = ref<MetaField[]>(init?.customMetaFields ?? []);
const isBanner          = ref(init?.isBanner   ?? false);
const isFeatured        = ref(init?.isFeatured  ?? false);
const postLayout        = ref<string>(init?.customMetaFields?.find((f) => f.key === "post_layout")?.value   || "auto");
const postSidebar       = ref<string>(init?.customMetaFields?.find((f) => f.key === "post_sidebar")?.value  || "auto");

const parseDownloads = (): DownloadItem[] => {
  try { return JSON.parse(init?.customMetaFields?.find((f) => f.key === "post_downloads")?.value ?? "[]"); }
  catch { return []; }
};
const postDownloads = ref<DownloadItem[]>(parseDownloads());

const formData = ref<CreatePostRequest>({
  post_type:       init?.postType      ?? 1,
  title:           init?.title         ?? "",
  slug:            init?.slug          ?? "",
  content:         init?.content       ?? "",
  excerpt:         init?.excerpt       ?? "",
  featured_img_id: init?.featuredImgId,
  status:          init?.status        ?? 1,
  published_at:    undefined,
  password:        "",
  comment_status:  init?.commentStatus ?? 1,
  author_id:       init?.authorId ?? authStore.user?.id ?? 1,
  term_taxonomy_ids: [],
});

const seoData = ref({
  meta_title:   init?.seo?.meta_title   ?? "",
  meta_desc:    init?.seo?.meta_desc    ?? "",
  og_title:     init?.seo?.og_title     ?? "",
  og_image:     init?.seo?.og_image     ?? "",
  canonical_url: init?.seo?.canonical_url ?? "",
  robots:       init?.seo?.robots       ?? "index,follow",
});

// ── Unsaved changes ───────────────────────────────────────────────────────
const hasUnsavedChanges = ref(false);
watch(
  [formData, selectedCategories, selectedTags, isBanner, isFeatured, seoData, metaFields],
  () => { hasUnsavedChanges.value = true; },
  { deep: true },
);
const markSaved = () => {
  hasUnsavedChanges.value = false;
  try { localStorage.removeItem(autoSaveKey.value); } catch {}
};

// ── Auto-save ─────────────────────────────────────────────────────────────
const autoSaveKey = computed(() =>
  props.mode === "edit" && init?.id ? `blog:draft:${init.id}` : "blog:draft:new",
);
const lastAutoSaved  = ref<Date | null>(null);
const autoSavedLabel = computed(() => {
  if (!lastAutoSaved.value) return "";
  return t("admin.posts.editor.auto_saved", {
    time: lastAutoSaved.value.toLocaleTimeString(undefined, { hour: "2-digit", minute: "2-digit" }),
  });
});
const showDraftRestore = ref(false);
const savedDraft = ref<{ title?: string; content?: string; excerpt?: string; savedAt?: string } | null>(null);

const doAutoSave = () => {
  if (!formData.value.title && !formData.value.content) return;
  try {
    localStorage.setItem(autoSaveKey.value, JSON.stringify({
      title: formData.value.title, content: formData.value.content,
      excerpt: formData.value.excerpt, savedAt: new Date().toISOString(),
    }));
    lastAutoSaved.value = new Date();
  } catch {}
};
const restoreDraft = () => {
  if (!savedDraft.value) return;
  if (savedDraft.value.title)   formData.value.title   = savedDraft.value.title;
  if (savedDraft.value.content) formData.value.content = savedDraft.value.content;
  if (savedDraft.value.excerpt) formData.value.excerpt = savedDraft.value.excerpt;
  showDraftRestore.value = false;
  localStorage.removeItem(autoSaveKey.value);
};
const discardDraft = () => {
  localStorage.removeItem(autoSaveKey.value);
  showDraftRestore.value = false;
};

// ── Revision modal ────────────────────────────────────────────────────────
const showRevisionModal = ref(false);

// ── Word count ────────────────────────────────────────────────────────────
const charCount = computed(() =>
  (formData.value.content ?? "").replace(/[#*`\[\]()>_~\-|]/g, "").replace(/\s+/g, "").trim().length,
);
const readingMinutes = computed(() => Math.max(1, Math.ceil(charCount.value / 400)));

// ── Build & save ──────────────────────────────────────────────────────────
const buildPayload = (status: number): CreatePostRequest | UpdatePostRequest => {
  formData.value.term_taxonomy_ids = [
    ...selectedCategories.value,
    ...selectedTags.value.map((t) => t.term_taxonomy_id),
  ];
  if (publishedAtLocal.value) {
    formData.value.published_at = new Date(publishedAtLocal.value).toISOString();
  }
  const metas: Record<string, string> = {};
  metaFields.value.forEach((field) => { if (field.key) metas[field.key] = field.value; });
  metas["is_banner"]   = isBanner.value   ? "1" : "";
  metas["is_featured"] = isFeatured.value ? "1" : "";
  if (postLayout.value  && postLayout.value  !== "auto") metas["post_layout"]  = postLayout.value;
  if (postSidebar.value && postSidebar.value !== "auto") metas["post_sidebar"] = postSidebar.value;
  const validDownloads = postDownloads.value.filter((d) => d.name && d.url);
  if (validDownloads.length > 0) metas["post_downloads"] = JSON.stringify(validDownloads);
  else delete metas["post_downloads"];
  return { ...formData.value, post_type: formData.value.post_type ?? 1, status, metas };
};

const isFormValid = computed(
  () => formData.value.title.trim() !== "" && (formData.value.content ?? "").trim() !== "",
);

// True when the selected publishedAt is more than 1 minute in the future.
const isFuturePublish = computed(() => {
  if (!publishedAtLocal.value) return false;
  return new Date(publishedAtLocal.value).getTime() > Date.now() + 60_000;
});

const handleSaveDraft = async () => {
  if (!isFormValid.value) { toast.add({ title: t("admin.posts.editor.fill_required"), color: "warning" }); return; }
  await uploadPendingImages();
  emit("save", buildPayload(1));
};
const handlePublish = async () => {
  if (!isFormValid.value) { toast.add({ title: t("admin.posts.editor.fill_required"), color: "warning" }); return; }
  await uploadPendingImages();
  // If publish time is in the future, save as draft — server cron will auto-publish at that time.
  emit("save", buildPayload(isFuturePublish.value ? 1 : 2));
};

const reset = () => {
  formData.value = {
    post_type: 1, title: "", slug: "", content: "", excerpt: "",
    featured_img_id: undefined, status: 1, published_at: undefined,
    password: "", comment_status: 1,
    author_id: authStore.user?.id ?? 1, term_taxonomy_ids: [],
  };
  selectedCategories.value = [];
  selectedTags.value       = [];
  metaFields.value         = [];
  isBanner.value           = false;
  isFeatured.value         = false;
  featuredImageUrl.value   = "";
  publishedAtLocal.value   = toDatetimeInputValue(new Date());
  seoData.value = { meta_title: "", meta_desc: "", og_title: "", og_image: "", canonical_url: "", robots: "index,follow" };
  hasUnsavedChanges.value  = false;
  try { localStorage.removeItem(autoSaveKey.value); } catch {}
};

defineExpose({ reset, isDirty: hasUnsavedChanges, getIsDirty: () => hasUnsavedChanges.value, markSaved, seoData });

// ── Lifecycle ─────────────────────────────────────────────────────────────
let autoSaveTimer: ReturnType<typeof setInterval>;

onMounted(async () => {
  editorLoading.value = false;
  await Promise.all([categoryStore.loadCategories(), tagStore.loadTags()]);
  sidebarLoading.value = false;

  if (!init?.content) {
    try {
      const saved = localStorage.getItem(autoSaveKey.value);
      if (saved) {
        const draft = JSON.parse(saved);
        if (draft.title || draft.content) { savedDraft.value = draft; showDraftRestore.value = true; }
      }
    } catch {}
  }

  const autoSaveEnabled  = optionsStore.get("auto_save", "true") !== "false";
  const autoSaveInterval = Math.max(10, Number(optionsStore.get("auto_save_interval", "60"))) * 1000;
  if (autoSaveEnabled) autoSaveTimer = setInterval(doAutoSave, autoSaveInterval);
});

onUnmounted(() => clearInterval(autoSaveTimer));
</script>

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
          <!-- 草稿恢复提示 -->
          <div v-if="showDraftRestore" class="pt-8 pb-3">
            <UAlert
              icon="i-tabler-device-floppy"
              color="warning"
              variant="subtle"
              :title="t('admin.posts.editor.draft_found')">
              <template #description>
                <div
                  class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-2">
                  <span class="text-sm text-muted">
                    {{
                      savedDraft?.savedAt
                        ? t("admin.posts.editor.draft_saved_at", {
                            time: new Date(savedDraft.savedAt).toLocaleString(),
                          })
                        : ""
                    }}
                  </span>
                  <div class="flex gap-2">
                    <UButton
                      size="xs"
                      color="primary"
                      variant="soft"
                      @click="restoreDraft"
                      >{{ t("admin.posts.editor.restore_draft") }}</UButton
                    >
                    <UButton
                      size="xs"
                      color="neutral"
                      variant="ghost"
                      @click="discardDraft"
                      >{{ t("admin.posts.editor.discard") }}</UButton
                    >
                  </div>
                </div>
              </template>
            </UAlert>
          </div>

          <!-- 标题 -->
          <div class="pt-8 pb-3">
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

          <!-- 编辑器骨架屏 -->
          <div v-if="editorLoading" class="py-4">
            <div
              class="flex items-center gap-2 border-b border-default pb-2 mb-6">
              <USkeleton v-for="i in 8" :key="i" class="h-7 w-7 rounded" />
              <USkeleton class="h-7 w-px mx-1" />
              <USkeleton
                v-for="i in 5"
                :key="'b' + i"
                class="h-7 w-7 rounded" />
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
            :starter-kit="{ codeBlock: false, blockquote: false }"
            :handlers="editorHandlers"
            :ui="{ base: 'px-8 sm:px-16 pt-14 pb-6' }"
            class="min-h-[500px] pb-4">
            <div
              class="border-b border-default sticky top-0 inset-x-0 z-10 bg-default/95 backdrop-blur before:content-[''] before:absolute before:inset-x-0 before:bottom-full before:h-8 before:bg-default">
              <div class="flex items-center">
                <div
                  ref="toolbarScrollRef"
                  class="flex-1 min-w-0 overflow-hidden">
                  <UEditorToolbar
                    :editor="editor"
                    :items="fullToolbarItems"
                    layout="fixed"
                    class="px-4 py-1.5" />
                </div>
                <div class="shrink-0 flex items-center gap-0.5 mr-1.5">
                  <UPopover v-if="hasOverflow" :content="{ align: 'end' }">
                    <UButton
                      variant="ghost"
                      color="neutral"
                      size="xs"
                      icon="i-tabler-dots" />
                    <template #content>
                      <UEditorToolbar
                        :editor="editor"
                        :items="overflowItems"
                        layout="fixed"
                        class="p-2 flex-wrap max-w-xs" />
                    </template>
                  </UPopover>
                  <UTooltip :text="t('admin.posts.editor.plugin_toolbar')">
                    <UButton
                      variant="ghost"
                      color="neutral"
                      size="xs"
                      :icon="
                        pluginToolbarExpanded
                          ? 'i-tabler-puzzle'
                          : 'i-tabler-puzzle-off'
                      "
                      @click="pluginToolbarExpanded = !pluginToolbarExpanded" />
                  </UTooltip>
                </div>
              </div>
              <div
                v-if="pluginToolbarExpanded"
                role="toolbar"
                class="border-t border-default/50 has-[button]:flex hidden items-stretch gap-1.5 px-4 py-1.5">
                <div role="group" class="flex items-center gap-0.5">
                  <ContributionSlot
                    name="post-editor:toolbar"
                    :ctx="{ editor }"
                    class="contents"
                    @command="
                      (cmdId: string) => handlePluginCommand(cmdId, editor)
                    ">
                    <template #menu="{ item, execute }">
                      <UTooltip :text="item.title">
                        <UButton
                          variant="ghost"
                          color="neutral"
                          size="xs"
                          :icon="item.icon"
                          @click="execute" />
                      </UTooltip>
                    </template>
                  </ContributionSlot>
                </div>
              </div>
            </div>

            <UEditorToolbar
              :editor="editor"
              :items="bubbleItems"
              class="z-50"
              layout="bubble"
              :should-show="
                ({ editor: e, view, state }) => {
                  const { selection } = state;
                  return (
                    view.hasFocus() && !selection.empty && !e.isActive('image')
                  );
                }
              " />

            <UEditorToolbar
              :editor="editor"
              :items="imageBubbleItems"
              class="z-50"
              layout="bubble"
              :should-show="({ editor: e }) => e.isActive('image')" />

            <UEditorToolbar
              :editor="editor"
              :items="linkBubbleItems"
              class="z-50"
              layout="bubble"
              :should-show="({ editor: e, state }) => e.isActive('link') && state.selection.empty" />

            <EditorLinkPopover ref="linkPopoverRef" :editor="editor" />

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
                @click="
                  (e) => {
                    e.stopPropagation();
                    const selected = onClick();
                    handlers.suggestion
                      ?.execute(editor, { pos: selected?.pos })
                      .run();
                  }
                " />
              <UDropdownMenu
                v-slot="{ open }"
                :modal="false"
                :items="dragHandleItems(editor)"
                :content="{ side: 'left' }"
                :ui="{ content: 'w-52', label: 'text-xs' }"
                @update:open="
                  editor.chain().setMeta('lockDragHandle', $event).run()
                ">
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

            <UEditorSuggestionMenu
              :editor="editor"
              :items="suggestionItems"
              :append-to="appendToBody" />
            <UEditorMentionMenu
              :editor="editor"
              :items="[]"
              :append-to="appendToBody" />
            <UEditorEmojiMenu
              :editor="editor"
              :items="emojiItems"
              :append-to="appendToBody" />
          </UEditor>

          <!-- 字数统计 -->
          <div
            class="py-2 flex px-2 items-center gap-4 text-xs text-muted border-t border-default">
            <span>{{
              t("admin.posts.editor.char_count", { n: charCount })
            }}</span>
            <span>{{
              t("admin.posts.editor.reading_minutes", { n: readingMinutes })
            }}</span>
            <span v-if="autoSavedLabel" class="ml-auto flex items-center gap-1">
              <UIcon
                name="i-tabler-circle-check"
                class="size-3 text-success" />{{ autoSavedLabel }}
            </span>
          </div>

          <ContributionSlot
            name="admin:post-editor-footer"
            :ctx="{ formData: formData }"
            class="empty:hidden" />

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
import { Emoji, gitHubEmojis } from "@tiptap/extension-emoji";
import { TextAlign } from "@tiptap/extension-text-align";
import { InlineMath } from "@tiptap/extension-mathematics";
import "katex/dist/katex.min.css";
import { ImageUpload } from "../extensions/ImageUpload";
import { Callout } from "../extensions/Callout";
import { CodeBlockWithLang } from "../extensions/CodeBlockWithLang";
import Blockquote from "@tiptap/extension-blockquote";
import { MathBlock } from "../extensions/MathBlock";
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

// ── Editor extensions & emoji ─────────────────────────────────────────────
const { pluginExtensions } = usePluginEditorExtensions();
const editorExtensions = computed(() => [
  Emoji,
  TextAlign.configure({ types: ["heading", "paragraph"] }),
  Blockquote.extend({ parseMarkdown: null as any }),
  InlineMath,
  MathBlock,
  ImageUpload,
  Callout,
  CodeBlockWithLang,
  ...pluginExtensions.value,
]);
const appendToBody = import.meta.client ? () => document.body : undefined;
const emojiItems = gitHubEmojis.filter(
  (e) => !e.name.startsWith("regional_indicator_"),
);

// ── Toolbar config (composable) ───────────────────────────────────────────
const {
  toolbarItems: fullToolbarItems,
  bubbleItems,
  linkBubbleItems,
  imageBubbleItems,
  suggestionItems,
  selectedNode,
  dragHandleItems,
  toolbarScrollRef,
  hasOverflow,
  overflowItems,
} = usePostEditorToolbar();

// ── Link popover ──────────────────────────────────────────────────────────
const linkPopoverRef = ref<{ openForInsert: () => void; openForEdit: () => void } | null>(null);

// ── Plugin toolbar toggle ─────────────────────────────────────────────────
const pluginToolbarExpanded = useLocalStorage(
  "nuxtblog:editor:plugin-toolbar",
  true,
);

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

const editorLoading = ref(true);
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

// ── Draft & auto-save (composable) ────────────────────────────────────────
const {
  autoSaveKey,
  autoSavedLabel,
  showDraftRestore,
  savedDraft,
  hasUnsavedChanges,
  restoreDraft,
  discardDraft,
  markSaved,
  startAutoSave,
} = usePostEditorDraft(formData, {
  mode: props.mode,
  postId: init?.id,
  hasInitialContent: !!init?.content,
});

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
    hasUnsavedChanges.value = true;
  },
  { deep: true },
);

// ── Image upload (composable) ─────────────────────────────────────────────
const toast = useToast();
const editorRef = ref<any>(null);

const { editorHandlers: baseEditorHandlers, uploadPendingImages, hasPendingUploads } =
  usePostEditorImageUpload(formData, editorRef);

const editorHandlers = computed(() => ({
  ...baseEditorHandlers.value,
  link: {
    canExecute: (editor: Editor) => editor.can().setLink({ href: '' }) || editor.can().unsetLink(),
    execute: (editor: Editor) => {
      if (editor.isActive('link')) {
        linkPopoverRef.value?.openForEdit()
      } else {
        linkPopoverRef.value?.openForInsert()
      }
      return editor.chain()
    },
    isActive: (editor: Editor) => editor.isActive('link'),
  },
  "link-edit": {
    canExecute: (editor: Editor) => editor.isActive('link'),
    execute: (_editor: Editor) => linkPopoverRef.value?.openForEdit(),
    isActive: (_editor: Editor) => false,
  },
  "link-open": {
    canExecute: (editor: Editor) => editor.isActive('link'),
    execute: (editor: Editor) => {
      const href = editor.getAttributes('link').href
      if (href) window.open(href, '_blank')
    },
    isActive: (_editor: Editor) => false,
  },
  "link-unlink": {
    canExecute: (editor: Editor) => editor.isActive('link'),
    execute: (editor: Editor) => editor.chain().focus().extendMarkRange('link').unsetLink().run(),
    isActive: (_editor: Editor) => false,
  },
}));

// ── Revision modal & preview ──────────────────────────────────────────────
const showRevisionModal = ref(false);
const showPreview = ref(false);

// ── Word count ────────────────────────────────────────────────────────────
const charCount = computed(
  () =>
    (formData.value.content ?? "")
      .replace(/[#*`\[\]()>_~\-|]/g, "")
      .replace(/\s+/g, "")
      .trim().length,
);
const readingMinutes = computed(() =>
  Math.max(1, Math.ceil(charCount.value / 400)),
);

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
  if (hasPendingUploads()) {
    toast.add({
      title: t("admin.posts.editor.pending_upload_warning"),
      color: "warning",
    });
    return;
  }
  await uploadPendingImages();
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
  if (hasPendingUploads()) {
    toast.add({
      title: t("admin.posts.editor.pending_upload_warning"),
      color: "warning",
    });
    return;
  }
  await uploadPendingImages();
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
  hasUnsavedChanges.value = false;
  try {
    localStorage.removeItem(autoSaveKey.value);
  } catch {}
};

defineExpose({
  reset,
  isDirty: hasUnsavedChanges,
  getIsDirty: () => hasUnsavedChanges.value,
  markSaved,
  seoData,
});

// ── Lifecycle ─────────────────────────────────────────────────────────────
onMounted(async () => {
  // Small delay to allow plugin extensions to register before editor mounts
  await new Promise((r) => setTimeout(r, 150));
  editorLoading.value = false;
  await Promise.all([categoryStore.loadCategories(), tagStore.loadTags()]);
  sidebarLoading.value = false;

  startAutoSave();
});
</script>

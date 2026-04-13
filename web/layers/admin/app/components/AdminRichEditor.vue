<template>
  <!-- 编辑器骨架屏 -->
  <div v-if="editorLoading" class="py-4">
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

  <div v-else>
    <!-- 草稿恢复提示 -->
    <div v-if="showDraftRestore" class="pt-4 pb-3">
      <UAlert
        icon="i-tabler-device-floppy"
        color="warning"
        variant="subtle"
        :title="t('admin.editor.draft_found')">
        <template #description>
          <div
            class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-2">
            <span class="text-sm text-muted">
              {{
                savedDraft?.savedAt
                  ? t("admin.editor.draft_saved_at", {
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
                >{{ t("admin.editor.restore_draft") }}</UButton
              >
              <UButton
                size="xs"
                color="neutral"
                variant="ghost"
                @click="discardDraft"
                >{{ t("admin.editor.discard") }}</UButton
              >
            </div>
          </div>
        </template>
      </UAlert>
    </div>

    <!-- 编辑器 -->
    <div class="overflow-y-clip">
    <UEditor
      ref="editorRef"
      v-slot="{ editor, handlers }"
      :model-value="modelValue"
      content-type="markdown"
      :placeholder="placeholder"
      :extensions="editorExtensions"
      :starter-kit="{ codeBlock: false, blockquote: false }"
      :handlers="mergedHandlers"
      :ui="editorUi ?? { base: 'px-8 sm:px-16 pt-14 pb-6' }"
      :class="['min-h-[500px] pb-4', editorClass]"
      @update:model-value="$emit('update:modelValue', $event)">
      <div
        class="border-b border-default sticky top-0 inset-x-0 z-10 bg-default/95 backdrop-blur before:content-[''] before:absolute before:inset-x-0 before:bottom-full before:h-8 before:bg-default">
        <div class="flex items-center">
          <div ref="toolbarScrollRef" class="flex-1 min-w-0 overflow-hidden">
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
            <template v-if="enablePluginToolbar">
              <UTooltip :text="t('admin.editor.plugin_toolbar')">
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
            </template>
          </div>
        </div>
        <div
          v-if="enablePluginToolbar && pluginToolbarExpanded"
          role="toolbar"
          class="border-t border-default/50 has-[button]:flex hidden items-stretch gap-1.5 px-4 py-1.5">
          <div role="group" class="flex items-center gap-0.5">
            <slot name="toolbar-extra" :editor="editor" />
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
            return view.hasFocus() && !selection.empty && !e.isActive('image');
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
        :should-show="
          ({ editor: e, state }) => e.isActive('link') && state.selection.empty
        " />

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
    </div>

    <!-- 字数统计 -->
    <div
      class="py-2 flex px-2 items-center gap-4 text-xs text-muted border-t border-default">
      <span>{{ t("admin.editor.char_count", { n: charCount }) }}</span>
      <span>{{
        t("admin.editor.reading_minutes", { n: readingMinutes })
      }}</span>
      <span v-if="autoSavedLabel" class="ml-auto flex items-center gap-1">
        <UIcon name="i-tabler-circle-check" class="size-3 text-success" />{{
          autoSavedLabel
        }}
      </span>
    </div>

    <slot name="footer" />
  </div>
</template>

<script setup lang="ts">
import type { Editor } from "@tiptap/vue-3";
import { Emoji, gitHubEmojis } from "@tiptap/extension-emoji";
import { TextAlign } from "@tiptap/extension-text-align";
import { InlineMath } from "@tiptap/extension-mathematics";
import "katex/dist/katex.min.css";
import { ImageUpload } from "../extensions/ImageUpload";
import { Callout } from "../extensions/Callout";
import { CodeBlockWithLang } from "../extensions/CodeBlockWithLang";
import Blockquote from "@tiptap/extension-blockquote";
import { MathBlock } from "../extensions/MathBlock";

const { t } = useI18n();

// ── Props ──────────────────────────────────────────────────────────────────
const props = withDefaults(
  defineProps<{
    modelValue: string;
    placeholder?: string;
    // Image upload
    imageCategory?: string;
    imageUploader?: (file: File) => Promise<string>;
    // Draft
    draftKeyPrefix?: string;
    draftEntityId?: number;
    draftMode?: "create" | "edit";
    hasInitialContent?: boolean;
    // UI / extensions
    enablePluginToolbar?: boolean;
    editorClass?: string;
    editorUi?: Record<string, any>;
    extraExtensions?: any[];
  }>(),
  {
    placeholder: undefined,
    imageCategory: "post",
    imageUploader: undefined,
    draftKeyPrefix: "blog:draft",
    draftEntityId: undefined,
    draftMode: "create",
    hasInitialContent: false,
    enablePluginToolbar: false,
    editorClass: undefined,
    editorUi: undefined,
    extraExtensions: () => [],
  },
);

defineEmits<{
  "update:modelValue": [value: string];
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
  ...props.extraExtensions,
]);

const appendToBody = import.meta.client ? () => document.body : undefined;
const emojiItems = gitHubEmojis.filter(
  (e) => !e.name.startsWith("regional_indicator_"),
);

// ── Toolbar config ────────────────────────────────────────────────────────
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
} = useEditorToolbar();

// ── Link popover ──────────────────────────────────────────────────────────
const linkPopoverRef = ref<{
  openForInsert: () => void;
  openForEdit: () => void;
} | null>(null);

// ── Plugin toolbar toggle ────────────────────────────────────────────────
const pluginToolbarExpanded = useLocalStorage(
  "nuxtblog:editor:plugin-toolbar",
  true,
);

// ── Internal formData ref for composables ────────────────────────────────
const internalFormData = computed({
  get: () => ({ content: props.modelValue }),
  set: () => {},
});
const formDataRef = ref({ content: props.modelValue });
watch(
  () => props.modelValue,
  (val) => {
    formDataRef.value.content = val;
  },
);

// ── Image upload ──────────────────────────────────────────────────────────
const editorRef = ref<any>(null);

const {
  editorHandlers: baseEditorHandlers,
  uploadPendingImages,
  hasPendingUploads,
} = useEditorImageUpload(formDataRef, editorRef, {
  imageCategory: props.imageCategory,
  imageUploader: props.imageUploader,
});

const mergedHandlers = computed(() => ({
  ...baseEditorHandlers.value,
  link: {
    canExecute: (editor: Editor) =>
      editor.can().setLink({ href: "" }) || editor.can().unsetLink(),
    execute: (editor: Editor) => {
      if (editor.isActive("link")) {
        linkPopoverRef.value?.openForEdit();
      } else {
        linkPopoverRef.value?.openForInsert();
      }
      return editor.chain();
    },
    isActive: (editor: Editor) => editor.isActive("link"),
  },
  "link-edit": {
    canExecute: (editor: Editor) => editor.isActive("link"),
    execute: (_editor: Editor) => linkPopoverRef.value?.openForEdit(),
    isActive: (_editor: Editor) => false,
  },
  "link-open": {
    canExecute: (editor: Editor) => editor.isActive("link"),
    execute: (editor: Editor) => {
      const href = editor.getAttributes("link").href;
      if (href) window.open(href, "_blank");
    },
    isActive: (_editor: Editor) => false,
  },
  "link-unlink": {
    canExecute: (editor: Editor) => editor.isActive("link"),
    execute: (editor: Editor) =>
      editor.chain().focus().extendMarkRange("link").unsetLink().run(),
    isActive: (_editor: Editor) => false,
  },
}));

// ── Draft & auto-save ────────────────────────────────────────────────────
const draftFormData = ref<{
  title?: string;
  content?: string;
  excerpt?: string;
}>({
  content: props.modelValue,
});
watch(
  () => props.modelValue,
  (val) => {
    draftFormData.value.content = val;
  },
);

const {
  autoSavedLabel,
  showDraftRestore,
  savedDraft,
  hasUnsavedChanges,
  restoreDraft,
  discardDraft,
  markSaved,
  startAutoSave,
} = useEditorDraft(draftFormData, {
  mode: props.draftMode,
  entityId: props.draftEntityId,
  keyPrefix: props.draftKeyPrefix,
  hasInitialContent: props.hasInitialContent,
});

// ── Word count ────────────────────────────────────────────────────────────
const charCount = computed(
  () =>
    (props.modelValue ?? "")
      .replace(/[#*`\[\]()>_~\-|]/g, "")
      .replace(/\s+/g, "")
      .trim().length,
);
const readingMinutes = computed(() =>
  Math.max(1, Math.ceil(charCount.value / 400)),
);

// ── Loading state ─────────────────────────────────────────────────────────
const editorLoading = ref(true);

onMounted(async () => {
  await new Promise((r) => setTimeout(r, 150));
  editorLoading.value = false;
  startAutoSave();
});

// ── Expose ────────────────────────────────────────────────────────────────
defineExpose({
  editor: computed(() => editorRef.value?.editor ?? null),
  editorRef,
  uploadPendingImages,
  hasPendingUploads,
  showDraftRestore,
  savedDraft,
  autoSavedLabel,
  restoreDraft,
  discardDraft,
  hasUnsavedChanges,
  markSaved,
  startAutoSave,
});
</script>

import type {
  EditorToolbarItem,
  EditorSuggestionMenuItem,
  DropdownMenuItem,
} from "@nuxt/ui";
import type { JSONContent, Editor } from "@tiptap/vue-3";
import { mapEditorItems } from "@nuxt/ui/utils/editor";

export const usePostEditorToolbar = () => {
  const { t } = useI18n();

  const selectedNode = ref<{ node: JSONContent; pos: number }>();

  const toolbarItems = computed<EditorToolbarItem[][]>(() => [
    [
      {
        kind: "undo",
        icon: "i-tabler-arrow-back-up",
        tooltip: { text: t("admin.posts.editor.undo") },
      },
      {
        kind: "redo",
        icon: "i-tabler-arrow-forward-up",
        tooltip: { text: t("admin.posts.editor.redo") },
      },
    ],
    [
      {
        icon: "i-tabler-heading",
        tooltip: { text: t("admin.posts.editor.heading") },
        content: { align: "start" },
        items: [
          { kind: "heading", level: 1, icon: "i-tabler-h-1", label: t("admin.posts.editor.heading_1") },
          { kind: "heading", level: 2, icon: "i-tabler-h-2", label: t("admin.posts.editor.heading_2") },
          { kind: "heading", level: 3, icon: "i-tabler-h-3", label: t("admin.posts.editor.heading_3") },
          { kind: "heading", level: 4, icon: "i-tabler-h-4", label: t("admin.posts.editor.heading_4") },
        ],
      },
      {
        icon: "i-tabler-list",
        tooltip: { text: t("admin.posts.editor.list") },
        content: { align: "start" },
        items: [
          { kind: "bulletList",   icon: "i-tabler-list",         label: t("admin.posts.editor.bullet_list") },
          { kind: "orderedList",  icon: "i-tabler-list-numbers", label: t("admin.posts.editor.ordered_list") },
        ],
      },
      { kind: "blockquote",     icon: "i-tabler-blockquote",   tooltip: { text: t("admin.posts.editor.blockquote") } },
      { kind: "codeBlock",      icon: "i-tabler-code-dots",    tooltip: { text: t("admin.posts.editor.code_block") } },
    ],
    [
      { kind: "mark", mark: "bold",      icon: "i-tabler-bold",          tooltip: { text: t("admin.posts.editor.bold") } },
      { kind: "mark", mark: "italic",    icon: "i-tabler-italic",        tooltip: { text: t("admin.posts.editor.italic") } },
      { kind: "mark", mark: "underline", icon: "i-tabler-underline",     tooltip: { text: t("admin.posts.editor.underline") } },
      { kind: "mark", mark: "strike",    icon: "i-tabler-strikethrough", tooltip: { text: t("admin.posts.editor.strikethrough") } },
      { kind: "mark", mark: "code",      icon: "i-tabler-code",          tooltip: { text: t("admin.posts.editor.inline_code") } },
    ],
    [
      { kind: "link",           icon: "i-tabler-link",      tooltip: { text: t("admin.posts.editor.link") } },
      { kind: "image",          icon: "i-tabler-photo",     tooltip: { text: t("admin.posts.editor.image") } },
      { kind: "horizontalRule", icon: "i-tabler-separator", tooltip: { text: t("admin.posts.editor.horizontal_rule") } },
    ],
    [
      {
        icon: "i-tabler-align-justified",
        tooltip: { text: t("admin.posts.editor.alignment") },
        content: { align: "end" },
        items: [
          { kind: "textAlign", align: "left",    icon: "i-tabler-align-left",      label: t("admin.posts.editor.align_left") },
          { kind: "textAlign", align: "center",  icon: "i-tabler-align-center",    label: t("admin.posts.editor.align_center") },
          { kind: "textAlign", align: "right",   icon: "i-tabler-align-right",     label: t("admin.posts.editor.align_right") },
          { kind: "textAlign", align: "justify", icon: "i-tabler-align-justified", label: t("admin.posts.editor.align_justify") },
        ],
      },
    ],
    [
      { kind: "clearFormatting", icon: "i-tabler-clear-formatting", tooltip: { text: t("admin.posts.editor.clear_format") } },
    ],
  ]);

  const bubbleItems = computed<EditorToolbarItem[][]>(() => [
    [
      {
        label: t("admin.posts.editor.convert_to"),
        trailingIcon: "i-tabler-chevron-down",
        activeColor: "neutral",
        activeVariant: "ghost",
        content: { align: "start" },
        ui: { label: "text-xs" },
        items: [
          { type: "label", label: t("admin.posts.editor.convert_to") },
          { kind: "paragraph",   label: t("admin.posts.editor.body_text"),     icon: "i-tabler-text-size" },
          { kind: "heading", level: 1, icon: "i-tabler-h-1", label: t("admin.posts.editor.heading_1") },
          { kind: "heading", level: 2, icon: "i-tabler-h-2", label: t("admin.posts.editor.heading_2") },
          { kind: "heading", level: 3, icon: "i-tabler-h-3", label: t("admin.posts.editor.heading_3") },
          { kind: "bulletList",  label: t("admin.posts.editor.bullet_list"),   icon: "i-tabler-list" },
          { kind: "orderedList", label: t("admin.posts.editor.ordered_list"),  icon: "i-tabler-list-numbers" },
          { kind: "blockquote",  label: t("admin.posts.editor.blockquote"),    icon: "i-tabler-blockquote" },
          { kind: "codeBlock",   label: t("admin.posts.editor.code_block"),    icon: "i-tabler-code-dots" },
        ],
      },
    ],
    [
      { kind: "mark", mark: "bold",      icon: "i-tabler-bold" },
      { kind: "mark", mark: "italic",    icon: "i-tabler-italic" },
      { kind: "mark", mark: "underline", icon: "i-tabler-underline" },
      { kind: "mark", mark: "strike",    icon: "i-tabler-strikethrough" },
      { kind: "mark", mark: "code",      icon: "i-tabler-code" },
    ],
    [
      { kind: "link",             icon: "i-tabler-link" },
      { kind: "clearFormatting",  icon: "i-tabler-clear-formatting" },
    ],
  ]);

  const suggestionItems = computed<EditorSuggestionMenuItem[][]>(() => [
    [
      { type: "label", label: t("admin.posts.editor.styles") },
      { kind: "paragraph",   label: t("admin.posts.editor.body_text"),    icon: "i-tabler-text-size" },
      { kind: "heading", level: 1, label: t("admin.posts.editor.heading_1"), icon: "i-tabler-h-1" },
      { kind: "heading", level: 2, label: t("admin.posts.editor.heading_2"), icon: "i-tabler-h-2" },
      { kind: "heading", level: 3, label: t("admin.posts.editor.heading_3"), icon: "i-tabler-h-3" },
      { kind: "heading", level: 4, label: t("admin.posts.editor.heading_4"), icon: "i-tabler-h-4" },
      { kind: "bulletList",  label: t("admin.posts.editor.bullet_list"),   icon: "i-tabler-list" },
      { kind: "orderedList", label: t("admin.posts.editor.ordered_list"),  icon: "i-tabler-list-numbers" },
      { kind: "blockquote",  label: t("admin.posts.editor.blockquote_full"), icon: "i-tabler-blockquote" },
      { kind: "codeBlock",   label: t("admin.posts.editor.code_block"),    icon: "i-tabler-code-dots" },
    ],
    [
      { type: "label", label: t("admin.posts.editor.insert") },
      { kind: "mention",        label: t("admin.posts.editor.mention"),        icon: "i-tabler-at" },
      { kind: "emoji",          label: "Emoji",                                icon: "i-tabler-mood-smile" },
      { kind: "image",          label: t("admin.posts.editor.image"),          icon: "i-tabler-photo" },
      { kind: "horizontalRule", label: t("admin.posts.editor.horizontal_rule"), icon: "i-tabler-separator-horizontal" },
    ],
  ]);

  const dragHandleItems = (editor: Editor): DropdownMenuItem[][] => {
    if (!selectedNode.value?.node?.type) return [];
    return mapEditorItems(
      editor,
      [
        [
          { type: "label", label: String(selectedNode.value.node.type) },
          {
            label: t("admin.posts.editor.convert_to"),
            icon: "i-tabler-transform",
            children: [
              { kind: "paragraph",   label: t("admin.posts.editor.body_text"),    icon: "i-tabler-text-size" },
              { kind: "heading", level: 1, label: t("admin.posts.editor.heading_1"), icon: "i-tabler-h-1" },
              { kind: "heading", level: 2, label: t("admin.posts.editor.heading_2"), icon: "i-tabler-h-2" },
              { kind: "heading", level: 3, label: t("admin.posts.editor.heading_3"), icon: "i-tabler-h-3" },
              { kind: "bulletList",  label: t("admin.posts.editor.bullet_list"),   icon: "i-tabler-list" },
              { kind: "orderedList", label: t("admin.posts.editor.ordered_list"),  icon: "i-tabler-list-numbers" },
              { kind: "blockquote",  label: t("admin.posts.editor.blockquote"),    icon: "i-tabler-blockquote" },
              { kind: "codeBlock",   label: t("admin.posts.editor.code_block"),    icon: "i-tabler-code-dots" },
            ],
          },
          {
            kind: "clearFormatting",
            pos: selectedNode.value?.pos,
            label: t("admin.posts.editor.clear_format"),
            icon: "i-tabler-clear-formatting",
          },
        ],
        [
          { kind: "duplicate", pos: selectedNode.value?.pos, label: t("admin.posts.editor.copy_block"),     icon: "i-tabler-copy" },
          {
            label: t("admin.posts.editor.copy_clipboard"),
            icon: "i-tabler-clipboard",
            onSelect: async () => {
              if (!selectedNode.value) return;
              const node = editor.state.doc.nodeAt(selectedNode.value.pos);
              if (node) await navigator.clipboard.writeText(node.textContent);
            },
          },
        ],
        [
          { kind: "moveUp",   pos: selectedNode.value?.pos, label: t("admin.posts.editor.move_up"),   icon: "i-tabler-arrow-up" },
          { kind: "moveDown", pos: selectedNode.value?.pos, label: t("admin.posts.editor.move_down"), icon: "i-tabler-arrow-down" },
        ],
        [
          { kind: "delete", pos: selectedNode.value?.pos, label: t("admin.posts.editor.delete"), icon: "i-tabler-trash" },
        ],
      ],
      {},
    ) as DropdownMenuItem[][];
  };

  // ── Toolbar overflow detection ──────────────────────────────────────────
  const toolbarScrollRef = ref<HTMLElement>();
  const overflowFromIndex = ref<number>(Infinity);

  const checkToolbarOverflow = () => {
    const container = toolbarScrollRef.value;
    if (!container) return;

    const toolbar = container.querySelector('[role="toolbar"]');
    if (!toolbar) return;

    const containerRight = container.getBoundingClientRect().right;
    const groups = toolbar.querySelectorAll(':scope > [role="group"]');

    let cutoff = groups.length;
    for (let i = 0; i < groups.length; i++) {
      if (groups[i].getBoundingClientRect().right > containerRight + 1) {
        cutoff = i;
        break;
      }
    }
    overflowFromIndex.value = cutoff;
  };

  useResizeObserver(toolbarScrollRef, checkToolbarOverflow);

  const hasOverflow = computed(
    () => overflowFromIndex.value < toolbarItems.value.length,
  );

  const overflowItems = computed(() =>
    toolbarItems.value.slice(overflowFromIndex.value),
  );

  return {
    toolbarItems,
    bubbleItems,
    suggestionItems,
    selectedNode,
    dragHandleItems,
    toolbarScrollRef,
    hasOverflow,
    overflowItems,
  };
};

import type { Editor } from "@tiptap/vue-3";
import type { CreatePostRequest } from "~/types/api/post";

export function usePostEditorImageUpload(
  formData: Ref<CreatePostRequest>,
  editorRef: Ref<any>,
) {
  const { t } = useI18n();
  const toast = useToast();
  const mediaStore = useMediaStore();

  const editorHandlers = computed(() => ({
    image: {
      canExecute: (editor: Editor) => editor.isEditable,
      execute: (editor: Editor) => {
        editor.chain().focus().insertContent({ type: "imageUpload" }).run();
        return editor.chain();
      },
      isActive: (_editor: Editor) => false,
    },
    "callout-note": {
      canExecute: (editor: Editor) => editor.isEditable,
      execute: (editor: Editor) => editor.chain().focus().setCallout({ type: "note" }).run(),
      isActive: (editor: Editor) => editor.isActive("callout", { type: "note" }),
    },
    "callout-tip": {
      canExecute: (editor: Editor) => editor.isEditable,
      execute: (editor: Editor) => editor.chain().focus().setCallout({ type: "tip" }).run(),
      isActive: (editor: Editor) => editor.isActive("callout", { type: "tip" }),
    },
    "callout-important": {
      canExecute: (editor: Editor) => editor.isEditable,
      execute: (editor: Editor) => editor.chain().focus().setCallout({ type: "important" }).run(),
      isActive: (editor: Editor) => editor.isActive("callout", { type: "important" }),
    },
    "callout-warning": {
      canExecute: (editor: Editor) => editor.isEditable,
      execute: (editor: Editor) => editor.chain().focus().setCallout({ type: "warning" }).run(),
      isActive: (editor: Editor) => editor.isActive("callout", { type: "warning" }),
    },
    "callout-caution": {
      canExecute: (editor: Editor) => editor.isEditable,
      execute: (editor: Editor) => editor.chain().focus().setCallout({ type: "caution" }).run(),
      isActive: (editor: Editor) => editor.isActive("callout", { type: "caution" }),
    },
    "math-inline": {
      canExecute: (editor: Editor) => editor.isEditable,
      execute: (editor: Editor) => editor.chain().focus().insertContent("$E=mc^2$").run(),
      isActive: (_editor: Editor) => false,
    },
    "math-block": {
      canExecute: (editor: Editor) => editor.isEditable,
      execute: (editor: Editor) => editor.chain().focus().setMathBlock({ latex: 'E = mc^2' }).run(),
      isActive: (editor: Editor) => editor.isActive("blockMath"),
    },
    "mermaid-block": {
      canExecute: (editor: Editor) => editor.isEditable,
      execute: (editor: Editor) => editor.chain().focus().setCodeBlock({ language: "mermaid" }).run(),
      isActive: (_editor: Editor) => false,
    },
    "download-image": {
      canExecute: (editor: Editor) => editor.isActive("image"),
      execute: (editor: Editor) => {
        const attrs = editor.getAttributes("image");
        if (attrs.src) {
          const a = document.createElement("a");
          a.href = attrs.src;
          a.download = attrs.alt || "image";
          a.target = "_blank";
          a.click();
        }
      },
      isActive: (_editor: Editor) => false,
    },
    "remove-image": {
      canExecute: (editor: Editor) => editor.isActive("image"),
      execute: (editor: Editor) => editor.chain().focus().deleteSelection().run(),
      isActive: (_editor: Editor) => false,
    },
  }));

  const dataUrlToFile = (dataUrl: string, name: string): File => {
    const [header, b64] = dataUrl.split(",");
    const mime = header.match(/:(.*?);/)?.[1] ?? "image/jpeg";
    const ext = mime.split("/")[1] ?? "jpg";
    const bin = atob(b64);
    const arr = new Uint8Array(bin.length);
    for (let i = 0; i < bin.length; i++) arr[i] = bin.charCodeAt(i);
    return new File([arr], name.includes(".") ? name : `${name}.${ext}`, {
      type: mime,
    });
  };

  const uploadPendingImages = async (): Promise<void> => {
    const content = formData.value.content ?? "";
    const regex = /!\[([^\]]*)\]\((data:[^)]+)\)/g;
    const matches: { alt: string; src: string }[] = [];
    let m: RegExpExecArray | null;
    while ((m = regex.exec(content)) !== null)
      matches.push({ alt: m[1], src: m[2] });
    if (!matches.length) return;

    toast.add({
      title: t("admin.posts.editor.image_uploading"),
      color: "neutral",
      duration: 0,
      id: "img-upload",
    });
    try {
      const results = await Promise.all(
        matches.map(async ({ alt, src }) => {
          const name = alt || `image-${Date.now()}`;
          const result = await mediaStore.uploadMedia(
            dataUrlToFile(src, name),
            { title: name, category: "post" },
          );
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

  const hasPendingUploads = (): boolean => {
    const editor = editorRef.value?.editor;
    if (!editor) return false;
    let found = false;
    editor.state.doc.descendants((node: any) => {
      if (node.type.name === "imageUpload") found = true;
    });
    return found;
  };

  return {
    editorHandlers,
    uploadPendingImages,
    hasPendingUploads,
  };
}

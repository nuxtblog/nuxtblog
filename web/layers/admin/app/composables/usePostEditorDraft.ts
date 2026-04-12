import type { CreatePostRequest } from "~/types/api/post";

interface DraftData {
  title?: string;
  content?: string;
  excerpt?: string;
  savedAt?: string;
}

export function usePostEditorDraft(
  formData: Ref<CreatePostRequest>,
  opts: { mode: "create" | "edit"; postId?: number; hasInitialContent: boolean },
) {
  const { t } = useI18n();
  const optionsStore = useOptionsStore();

  const autoSaveKey = computed(() =>
    opts.mode === "edit" && opts.postId
      ? `blog:draft:${opts.postId}`
      : "blog:draft:new",
  );

  const lastAutoSaved = ref<Date | null>(null);
  const autoSavedLabel = computed(() => {
    if (!lastAutoSaved.value) return "";
    return t("admin.posts.editor.auto_saved", {
      time: lastAutoSaved.value.toLocaleTimeString(undefined, {
        hour: "2-digit",
        minute: "2-digit",
      }),
    });
  });

  const showDraftRestore = ref(false);
  const savedDraft = ref<DraftData | null>(null);
  const hasUnsavedChanges = ref(false);

  const doAutoSave = () => {
    if (!formData.value.title && !formData.value.content) return;
    try {
      localStorage.setItem(
        autoSaveKey.value,
        JSON.stringify({
          title: formData.value.title,
          content: formData.value.content,
          excerpt: formData.value.excerpt,
          savedAt: new Date().toISOString(),
        }),
      );
      lastAutoSaved.value = new Date();
    } catch {}
  };

  const restoreDraft = () => {
    if (!savedDraft.value) return;
    if (savedDraft.value.title) formData.value.title = savedDraft.value.title;
    if (savedDraft.value.content)
      formData.value.content = savedDraft.value.content;
    if (savedDraft.value.excerpt)
      formData.value.excerpt = savedDraft.value.excerpt;
    showDraftRestore.value = false;
    localStorage.removeItem(autoSaveKey.value);
  };

  const discardDraft = () => {
    localStorage.removeItem(autoSaveKey.value);
    showDraftRestore.value = false;
  };

  const markSaved = () => {
    hasUnsavedChanges.value = false;
    try {
      localStorage.removeItem(autoSaveKey.value);
    } catch {}
  };

  let autoSaveTimer: ReturnType<typeof setInterval>;

  const startAutoSave = () => {
    // Check for existing draft on mount
    if (!opts.hasInitialContent) {
      try {
        const saved = localStorage.getItem(autoSaveKey.value);
        if (saved) {
          const draft = JSON.parse(saved);
          if (draft.title || draft.content) {
            savedDraft.value = draft;
            showDraftRestore.value = true;
          }
        }
      } catch {}
    }

    const autoSaveEnabled = optionsStore.get("auto_save", "true") !== "false";
    const autoSaveInterval =
      Math.max(10, Number(optionsStore.get("auto_save_interval", "60"))) * 1000;
    if (autoSaveEnabled)
      autoSaveTimer = setInterval(doAutoSave, autoSaveInterval);
  };

  onUnmounted(() => clearInterval(autoSaveTimer));

  return {
    autoSaveKey,
    autoSavedLabel,
    showDraftRestore,
    savedDraft,
    hasUnsavedChanges,
    restoreDraft,
    discardDraft,
    markSaved,
    startAutoSave,
  };
}

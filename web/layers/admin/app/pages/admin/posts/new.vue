<template>
  <AdminPostEditor
    ref="editorRef"
    mode="create"
    :submitting="submitting"
    @save="handleSave" />
</template>

<script setup lang="ts">
import type { CreatePostRequest, UpdatePostRequest } from "~/types/api/post";

const submitting = ref(false);
const toast = useToast();
const { t } = useI18n();
const { apiFetch } = useApiFetch();

const editorRef = useTemplateRef<{
  reset: () => void;
  isDirty: Ref<boolean>;
  getIsDirty: () => boolean;
  markSaved: () => void;
  seoData: Ref<Record<string, string>>;
}>("editorRef");

// 离开页面未保存警告
onBeforeRouteLeave(() => {
  if (editorRef.value?.getIsDirty()) {
    return window.confirm(t("admin.posts.editor.unsaved_warning"));
  }
});

const handleBeforeUnload = (e: BeforeUnloadEvent) => {
  if (editorRef.value?.getIsDirty()) {
    e.preventDefault();
  }
};
onMounted(() => window.addEventListener("beforeunload", handleBeforeUnload));
onUnmounted(() =>
  window.removeEventListener("beforeunload", handleBeforeUnload),
);

const handleSave = async (payload: CreatePostRequest | UpdatePostRequest) => {
  submitting.value = true;
  try {
    const result = await apiFetch<{ id: number }>("/posts", {
      method: "POST",
      body: payload,
    });
    editorRef.value?.markSaved();
    const isScheduled = payload.status === 1 && !!(payload as any).published_at &&
      new Date((payload as any).published_at).getTime() > Date.now()
    toast.add({
      title: isScheduled ? t('admin.posts.editor.scheduled_saved') : t("admin.posts.editor.created"),
      icon: isScheduled ? 'i-tabler-clock' : 'i-tabler-circle-check',
      color: 'success',
    });
    await navigateTo(`/admin/posts/edit/${result.id}`);
  } catch (error: any) {
    const msg = error?.data?.message || error?.message || t("common.unknown_error");
    toast.add({ title: t("admin.posts.editor.save_failed"), description: msg, color: "error" });
  } finally {
    submitting.value = false;
  }
};
</script>

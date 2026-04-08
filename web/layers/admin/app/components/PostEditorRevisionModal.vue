<template>
  <UModal v-model:open="open" :title="t('admin.posts.editor.revisions_title')">
    <template #content>
      <div class="p-6">
        <div v-if="revisionsLoading" class="flex justify-center py-8">
          <UIcon name="i-tabler-loader-2" class="size-6 animate-spin text-muted" />
        </div>
        <div v-else-if="revisions.length === 0" class="text-center py-8 text-muted text-sm">
          {{ t("admin.posts.editor.no_revisions") }}
        </div>
        <div v-else class="space-y-2 max-h-96 overflow-y-auto">
          <div
            v-for="rev in revisions"
            :key="rev.id"
            class="flex items-center justify-between p-3 rounded-md border border-default hover:bg-muted transition-colors">
            <div class="flex-1 min-w-0">
              <p class="text-sm font-medium text-highlighted truncate">
                {{ rev.title || t("admin.posts.editor.untitled") }}
              </p>
              <p class="text-xs text-muted mt-0.5">{{ rev.created_at }}</p>
            </div>
            <UButton
              color="primary"
              variant="soft"
              size="xs"
              class="ml-3 shrink-0"
              :loading="restoringRevision"
              @click="restoreRevision(rev.id)">
              {{ t("admin.posts.editor.restore") }}
            </UButton>
          </div>
        </div>
      </div>
    </template>
  </UModal>
</template>

<script setup lang="ts">
interface Revision {
  id: number;
  post_id: number;
  author_id: number;
  title: string;
  rev_note: string;
  created_at: string;
}

const props = defineProps<{ postId?: number }>();
const open = defineModel<boolean>({ required: true });

const { t } = useI18n();
const { apiFetch } = useApiFetch();
const toast = useToast();

const revisions = ref<Revision[]>([]);
const revisionsLoading = ref(false);
const restoringRevision = ref(false);

const loadRevisions = async () => {
  if (!props.postId) return;
  revisionsLoading.value = true;
  try {
    const res = await apiFetch<{ list: Revision[] }>(
      `/posts/${props.postId}/revisions`,
      { params: { page: 1, size: 20 } },
    );
    revisions.value = res.list || [];
  } catch {
  } finally {
    revisionsLoading.value = false;
  }
};

watch(open, (val) => { if (val) loadRevisions(); });

const restoreRevision = async (revId: number) => {
  if (!props.postId) return;
  restoringRevision.value = true;
  try {
    await apiFetch(`/posts/${props.postId}/revisions/${revId}/restore`, { method: "POST" });
    toast.add({ title: t("admin.posts.editor.restore_success"), color: "success" });
    open.value = false;
    setTimeout(() => window.location.reload(), 1000);
  } catch {
    toast.add({ title: t("admin.posts.editor.restore_failed"), color: "error" });
  } finally {
    restoringRevision.value = false;
  }
};
</script>

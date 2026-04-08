<script setup lang="ts">
import type { ReportItem } from "~/composables/useReportApi";

const { t } = useI18n();
const reportApi = useReportApi();
const toast = useToast();

const statusFilter = ref("pending");
const page = ref(1);
const items = ref<ReportItem[]>([]);
const total = ref(0);
const rawLoading = ref(true);
const loading = useMinLoading(rawLoading);
const handlingId = ref<number | null>(null);

const searchKeyword = ref('')

const filteredItems = computed(() => {
  if (!searchKeyword.value.trim()) return items.value
  const q = searchKeyword.value.toLowerCase()
  return items.value.filter(i =>
    i.reason.toLowerCase().includes(q) ||
    i.detail?.toLowerCase().includes(q) ||
    i.reporter_name?.toLowerCase().includes(q) ||
    i.target_name?.toLowerCase().includes(q)
  )
})

const notesModal = ref(false);
const notesText = ref("");
const pendingHandle = ref<{
  id: number;
  status: "resolved" | "dismissed";
} | null>(null);

const STATUS_TABS = computed(() => [
  { label: t("admin.reports.tab_pending"), value: "pending" },
  { label: t("admin.reports.tab_resolved"), value: "resolved" },
  { label: t("admin.reports.tab_dismissed"), value: "dismissed" },
  { label: t("admin.reports.tab_all"), value: "all" },
]);

const statusColor = (s: string) =>
  s === "pending" ? "warning" : s === "resolved" ? "success" : "neutral";
const targetTypeLabel = (type: string) =>
  ({
    post: t("admin.reports.target_post"),
    comment: t("admin.reports.target_comment"),
    user: t("admin.reports.target_user"),
  })[type] ?? type;

const load = async () => {
  rawLoading.value = true;
  try {
    const res = await reportApi.list(statusFilter.value, page.value);
    items.value = res.items;
    total.value = res.total;
  } catch (e: any) {
    toast.add({ title: e?.message, color: "error" });
  } finally {
    rawLoading.value = false;
  }
};

watch([statusFilter, page], load, { immediate: true });
watch(statusFilter, () => {
  page.value = 1;
});

const openHandle = (id: number, status: "resolved" | "dismissed") => {
  pendingHandle.value = { id, status };
  notesText.value = "";
  notesModal.value = true;
};

const confirmHandle = async () => {
  if (!pendingHandle.value) return;
  handlingId.value = pendingHandle.value.id;
  try {
    await reportApi.handle(
      pendingHandle.value.id,
      pendingHandle.value.status,
      notesText.value,
    );
    toast.add({ title: t("admin.reports.handled"), color: "success" });
    notesModal.value = false;
    await load();
  } catch (e: any) {
    toast.add({ title: e?.message, color: "error" });
  } finally {
    handlingId.value = null;
  }
};

const formatDate = (s: string) => new Date(s).toLocaleString("zh-CN");

const targetUrl = (type: string, id: number) => {
  if (type === "post") return `/admin/posts/edit/${id}`
  if (type === "user") return `/user/${id}`
  return null // comment: no direct link
}
</script>

<template>
  <AdminPageContainer>
    <AdminPageHeader :title="$t('admin.reports.title')" :subtitle="$t('admin.reports.subtitle')" />

    <AdminPageContent>
      <!-- Status tabs -->
      <div class="flex items-center gap-1 border-b border-default pb-0 mb-4 overflow-x-auto">
        <button
          v-for="tab in STATUS_TABS"
          :key="tab.value"
          class="px-3 py-2 text-sm font-medium whitespace-nowrap transition-colors"
          :class="statusFilter === tab.value
            ? 'text-primary border-b-2 border-primary'
            : 'text-muted hover:text-highlighted'"
          @click="statusFilter = tab.value">
          {{ tab.label }}
        </button>
      </div>

      <!-- Search toolbar -->
      <div class="flex items-center gap-3 mb-4">
        <UInput
          v-model="searchKeyword"
          :placeholder="$t('admin.reports.search_placeholder')"
          leading-icon="i-tabler-search"
          class="w-56"
          size="sm">
          <template v-if="searchKeyword" #trailing>
            <UButton icon="i-tabler-x" color="neutral" variant="ghost" size="xs" @click="searchKeyword = ''" />
          </template>
        </UInput>
      </div>

      <!-- Handle notes modal -->
      <UModal v-model:open="notesModal" :ui="{ content: 'max-w-md' }">
        <template #content>
          <div class="p-6 max-h-[90vh] overflow-y-auto">
            <h3 class="text-lg font-semibold text-highlighted mb-4">
              {{
                pendingHandle?.status === "resolved"
                  ? $t("admin.reports.confirm_resolve")
                  : $t("admin.reports.confirm_dismiss")
              }}
            </h3>
            <UFormField :label="$t('admin.reports.handle_notes_label')">
              <UTextarea
                v-model="notesText"
                :rows="3"
                class="w-full"
                :placeholder="$t('admin.reports.handle_notes_placeholder')" />
            </UFormField>
            <div class="flex justify-end gap-2 mt-6">
              <UButton
                color="neutral"
                variant="ghost"
                @click="notesModal = false"
                >{{ $t("common.cancel") }}</UButton
              >
              <UButton
                :color="
                  pendingHandle?.status === 'resolved' ? 'success' : 'neutral'
                "
                :loading="!!handlingId"
                @click="confirmHandle">
                {{
                  pendingHandle?.status === "resolved"
                    ? $t("admin.reports.confirm_resolve")
                    : $t("admin.reports.confirm_dismiss")
                }}
              </UButton>
            </div>
          </div>
        </template>
      </UModal>

      <!-- Loading -->
      <div v-if="loading" class="space-y-3">
        <div v-for="i in 5" :key="i" class="flex items-center gap-4 p-4 border border-default rounded-md">
          <div class="flex-1 space-y-2">
            <USkeleton class="h-4 w-48" />
            <USkeleton class="h-3 w-full" />
            <USkeleton class="h-3 w-32" />
          </div>
          <div class="flex gap-2">
            <USkeleton class="h-8 w-20 rounded-md" />
            <USkeleton class="h-8 w-20 rounded-md" />
          </div>
        </div>
      </div>

      <!-- Empty -->
      <div v-else-if="filteredItems.length === 0" class="flex flex-col items-center justify-center py-16">
        <UIcon name="i-tabler-flag-off" class="size-16 text-muted mb-4" />
        <h3 class="text-lg font-medium text-highlighted mb-1">{{ $t("admin.reports.empty") }}</h3>
      </div>

      <!-- List -->
      <div v-else class="space-y-3">
        <div
          v-for="item in filteredItems"
          :key="item.id"
          class="flex items-start gap-4 p-4 border border-default rounded-md group hover:shadow-sm transition-all">
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2 flex-wrap mb-1">
              <UBadge
                :label="targetTypeLabel(item.target_type)"
                color="neutral"
                variant="soft"
                size="sm" />
              <UBadge
                :label="
                  item.status === 'pending'
                    ? $t('admin.reports.status_pending')
                    : item.status === 'resolved'
                      ? $t('admin.reports.status_resolved')
                      : $t('admin.reports.status_dismissed')
                "
                :color="statusColor(item.status)"
                variant="soft"
                size="sm" />
              <span class="text-xs text-muted">{{
                formatDate(item.created_at)
              }}</span>
            </div>
            <p class="text-sm font-medium text-highlighted">
              {{ $t("admin.reports.reason_label") }}{{ item.reason }}
            </p>
            <p v-if="item.detail" class="text-xs text-muted mt-0.5">
              {{ item.detail }}
            </p>
            <p class="text-xs text-muted mt-1 flex items-center gap-1 flex-wrap">
              {{ $t("admin.reports.reporter_label") }}
              <NuxtLink
                :to="`/user/${item.reporter_id}`"
                target="_blank"
                class="text-primary hover:underline font-medium">
                {{ item.reporter_name }}
              </NuxtLink>
              &nbsp;·&nbsp;
              {{ $t("admin.reports.target_id_label") }}
              <NuxtLink
                v-if="targetUrl(item.target_type, item.target_id)"
                :to="targetUrl(item.target_type, item.target_id)!"
                target="_blank"
                class="text-primary hover:underline font-medium inline-flex items-center gap-0.5">
                {{ item.target_name || item.target_id }}
                <UIcon name="i-tabler-external-link" class="size-3" />
              </NuxtLink>
              <span v-else class="font-medium text-default">{{ item.target_name || item.target_id }}</span>
            </p>
            <p v-if="item.notes" class="text-xs text-primary mt-1">
              {{ $t("admin.reports.notes_label") }}{{ item.notes }}
            </p>
          </div>
          <div v-if="item.status === 'pending'" class="flex gap-1 shrink-0">
            <UButton
              size="xs"
              color="success"
              variant="soft"
              icon="i-tabler-check"
              @click="openHandle(item.id, 'resolved')"
              >{{ $t("admin.reports.action_resolve") }}</UButton
            >
            <UButton
              size="xs"
              color="neutral"
              variant="soft"
              icon="i-tabler-x"
              @click="openHandle(item.id, 'dismissed')"
              >{{ $t("admin.reports.action_dismiss") }}</UButton
            >
          </div>
        </div>
      </div>
    </AdminPageContent>

    <AdminPageFooter v-if="total > 20">
      <UPagination v-model:page="page" :total="total" :items-per-page="20" />
    </AdminPageFooter>
  </AdminPageContainer>
</template>

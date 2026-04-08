<script setup lang="ts">
import { z } from "zod";
import type { AnnouncementItem } from "~/composables/useAnnouncementApi";

const { t } = useI18n();
const announcementApi = useAnnouncementApi();
const toast = useToast();

// ── Create modal ──────────────────────────────────────────────────────────────

const createModal = ref(false);

const schema = z.object({
  title: z.string().min(1).max(200),
  content: z.string().min(1).max(5000),
  type: z.enum(["info", "warning", "success", "danger"]),
});

type FormState = z.infer<typeof schema>;

const state = reactive<FormState>({ title: "", content: "", type: "info" });
const creating = ref(false);

const TYPE_OPTIONS = computed(() => [
  { label: t("admin.announcements.type_info"), value: "info" },
  { label: t("admin.announcements.type_warning"), value: "warning" },
  { label: t("admin.announcements.type_success"), value: "success" },
  { label: t("admin.announcements.type_danger"), value: "danger" },
]);

const typeColor = (type: string) =>
  ({ info: "info", warning: "warning", success: "success", danger: "error" })[type] ?? "neutral";

function openCreateModal() {
  state.title = "";
  state.content = "";
  state.type = "info";
  createModal.value = true;
}

async function handleCreate() {
  creating.value = true;
  try {
    await announcementApi.adminCreateAnnouncement({
      title: state.title,
      content: state.content,
      type: state.type,
    });
    toast.add({ title: t("admin.announcements.create_ok"), color: "success" });
    createModal.value = false;
    await load();
  } catch (e: any) {
    toast.add({ title: e?.message ?? t("admin.announcements.create_failed"), color: "error" });
  } finally {
    creating.value = false;
  }
}

// ── List ──────────────────────────────────────────────────────────────────────

const items = ref<AnnouncementItem[]>([]);
const total = ref(0);
const page = ref(1);
const pageSize = 20;
const rawLoading = ref(true);
const loading = useMinLoading(rawLoading);

// batch selection
const selected = ref<number[]>([]);
const batchAction = ref<string | undefined>(undefined);

const searchKeyword = ref('')

const filteredItems = computed(() => {
  if (!searchKeyword.value.trim()) return items.value
  const q = searchKeyword.value.toLowerCase()
  return items.value.filter(i =>
    i.title.toLowerCase().includes(q) ||
    i.content.toLowerCase().includes(q)
  )
})

const isAllSelected = computed(() =>
  items.value.length > 0 && items.value.every(i => selected.value.includes(i.id))
);
const isIndeterminate = computed(() =>
  selected.value.length > 0 && !isAllSelected.value
);

function toggleSelect(id: number) {
  const idx = selected.value.indexOf(id);
  if (idx > -1) selected.value.splice(idx, 1);
  else selected.value.push(id);
}

function toggleSelectAll() {
  if (isAllSelected.value) selected.value = [];
  else selected.value = items.value.map(i => i.id);
}

// edit modal
const editModal = ref(false);
const editTarget = ref<AnnouncementItem | null>(null);
const editLoading = ref(false);
const editState = reactive({ title: '', content: '', type: 'info' as string });

function openEditModal(item: AnnouncementItem) {
  editTarget.value = item;
  editState.title = item.title;
  editState.content = item.content;
  editState.type = item.type;
  editModal.value = true;
}

async function handleEdit() {
  if (!editTarget.value) return;
  editLoading.value = true;
  try {
    await announcementApi.adminUpdateAnnouncement(editTarget.value.id, {
      title: editState.title,
      content: editState.content,
      type: editState.type,
    });
    toast.add({ title: t('common.saved'), color: 'success' });
    editModal.value = false;
    await load();
  } catch (e: any) {
    toast.add({ title: e?.message ?? t('common.save_failed'), color: 'error' });
  } finally {
    editLoading.value = false;
  }
}

// delete modal
const deleteModal = ref(false);
const deleteTarget = ref<AnnouncementItem | null>(null);
const deleteLoading = ref(false);

const load = async () => {
  rawLoading.value = true;
  try {
    const res = await announcementApi.adminListAnnouncements({ page: page.value, size: pageSize });
    items.value = res.list ?? [];
    total.value = res.total;
  } catch (e: any) {
    toast.add({ title: e?.message ?? t("admin.announcements.load_failed"), color: "error" });
  } finally {
    rawLoading.value = false;
  }
};

watch(page, load, { immediate: true });

function openDeleteModal(item: AnnouncementItem) {
  deleteTarget.value = item;
  deleteModal.value = true;
}

async function confirmDelete() {
  if (!deleteTarget.value) return;
  deleteLoading.value = true;
  try {
    await announcementApi.adminDeleteAnnouncement(deleteTarget.value.id);
    toast.add({ title: t("admin.announcements.delete_ok"), color: "success" });
    deleteModal.value = false;
    selected.value = selected.value.filter(id => id !== deleteTarget.value!.id);
    await load();
  } catch (e: any) {
    toast.add({ title: e?.message ?? t("admin.announcements.delete_failed"), color: "error" });
  } finally {
    deleteLoading.value = false;
  }
}

async function applyBatch() {
  if (!batchAction.value || !selected.value.length) return;
  if (batchAction.value === 'delete') {
    try {
      await Promise.all(selected.value.map(id => announcementApi.adminDeleteAnnouncement(id)));
      toast.add({ title: t("admin.announcements.delete_ok"), color: "success" });
      selected.value = [];
      batchAction.value = undefined;
      await load();
    } catch (e: any) {
      toast.add({ title: e?.message ?? t("admin.announcements.delete_failed"), color: "error" });
    }
  }
}

function getItemActions(item: AnnouncementItem) {
  return [
    [
      { label: t('common.edit'), icon: 'i-tabler-pencil', onClick: () => openEditModal(item) },
    ],
    [
      { label: t('common.delete'), icon: 'i-tabler-trash', color: 'error' as const, onClick: () => openDeleteModal(item) },
    ],
  ];
}

const formatDate = (s: string) => new Date(s).toLocaleString("zh-CN");
</script>

<template>
  <AdminPageContainer>
    <AdminPageHeader
      :title="$t('admin.announcements.title')"
      :subtitle="$t('admin.announcements.subtitle')">
      <template #actions>
        <UButton color="primary" icon="i-tabler-plus" @click="openCreateModal">
          {{ $t("admin.announcements.btn_create") }}
        </UButton>
      </template>
    </AdminPageHeader>

    <AdminPageContent>
      <!-- Search toolbar -->
      <div class="flex items-center gap-3 mb-4">
        <UInput
          v-model="searchKeyword"
          :placeholder="$t('admin.announcements.search_placeholder')"
          leading-icon="i-tabler-search"
          class="w-56"
          size="sm">
          <template v-if="searchKeyword" #trailing>
            <UButton icon="i-tabler-x" color="neutral" variant="ghost" size="xs" @click="searchKeyword = ''" />
          </template>
        </UInput>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="space-y-3">
        <div v-for="i in 5" :key="i" class="flex items-center gap-4 p-4 border border-default rounded-md">
          <USkeleton class="size-4 rounded shrink-0" />
          <div class="flex-1 space-y-2">
            <USkeleton class="h-4 w-48" />
            <USkeleton class="h-3 w-full" />
          </div>
          <USkeleton class="h-5 w-16 rounded-full shrink-0" />
        </div>
      </div>

      <!-- Empty -->
      <div v-else-if="filteredItems.length === 0" class="flex flex-col items-center justify-center py-16">
        <UIcon name="i-tabler-speakerphone" class="size-16 text-muted mb-4" />
        <h3 class="text-lg font-medium text-highlighted mb-1">{{ $t("admin.announcements.empty") }}</h3>
        <p class="text-sm text-muted mb-4">{{ $t("admin.announcements.empty_desc", "发布一条公告让读者了解最新动态") }}</p>
        <UButton color="primary" icon="i-tabler-plus" @click="openCreateModal">
          {{ $t("admin.announcements.btn_create") }}
        </UButton>
      </div>

      <!-- List -->
      <div v-else class="space-y-3">
        <div
          v-for="item in filteredItems"
          :key="item.id"
          class="flex items-center group gap-4 p-4 bg-default border border-default rounded-md hover:shadow-sm transition-all">
          <UCheckbox
            :model-value="selected.includes(item.id)"
            @update:model-value="toggleSelect(item.id)" />
          <div class="flex-1 min-w-0">
            <div class="flex items-start justify-between gap-4">
              <div class="flex-1 min-w-0">
                <p class="text-sm font-semibold text-highlighted">{{ item.title }}</p>
                <div class="flex items-center gap-2 mt-1">
                  <span class="text-xs text-muted">{{ formatDate(item.created_at) }}</span>
                </div>
                <p v-if="item.content" class="text-xs text-muted mt-0.5 line-clamp-2">{{ item.content }}</p>
              </div>
              <div class="flex items-center gap-3 shrink-0">
                <UBadge
                  :label="$t('admin.announcements.type_' + item.type)"
                  :color="typeColor(item.type)"
                  variant="soft"
                  size="sm" />
                <UDropdownMenu :items="getItemActions(item)" :popper="{ placement: 'bottom-end' }">
                  <UButton
                    color="neutral"
                    variant="ghost"
                    icon="i-tabler-dots-vertical"
                    square
                    size="xs"
                    class="opacity-0 group-hover:opacity-100 transition-opacity" />
                </UDropdownMenu>
              </div>
            </div>
          </div>
        </div>
      </div>
    </AdminPageContent>

    <AdminPageFooter>
      <template #left>
        <template v-if="items.length > 0">
          <UCheckbox
            :model-value="isAllSelected"
            :indeterminate="isIndeterminate"
            @update:model-value="toggleSelectAll" />
          <template v-if="selected.length > 0">
            <span>{{ $t('common.selected_n', { n: selected.length }) }}</span>
            <USeparator orientation="vertical" class="h-4" />
            <USelect
              v-model="batchAction"
              :items="[{ label: $t('common.delete'), value: 'delete' }]"
              :placeholder="$t('admin.posts.batch_action')"
              class="w-36"
              size="sm" />
            <UButton color="primary" variant="soft" size="sm" :disabled="!batchAction" @click="applyBatch">{{ $t('common.apply') }}</UButton>
            <UButton color="neutral" variant="ghost" size="sm" @click="selected = []; batchAction = undefined">{{ $t('common.cancel') }}</UButton>
          </template>
          <span v-else class="text-xs">{{ $t('common.selectAll') }}</span>
        </template>
      </template>
      <template #right>
        <UPagination v-if="total > pageSize" v-model:page="page" :total="total" :items-per-page="pageSize" size="sm" />
      </template>
    </AdminPageFooter>

    <!-- Edit Modal -->
    <UModal v-model:open="editModal" :ui="{ content: 'max-w-lg' }">
      <template #content>
        <div class="p-6">
          <h3 class="text-lg font-semibold text-highlighted mb-5">{{ $t('common.edit') }}</h3>
          <div class="space-y-4">
            <UFormField :label="$t('admin.announcements.field_title')">
              <UInput v-model="editState.title" :placeholder="$t('admin.announcements.field_title')" class="w-full" />
            </UFormField>
            <UFormField :label="$t('admin.announcements.field_content')">
              <UTextarea v-model="editState.content" :rows="4" :placeholder="$t('admin.announcements.field_content')" class="w-full" />
            </UFormField>
            <UFormField :label="$t('admin.announcements.field_type')">
              <USelect v-model="editState.type" :items="TYPE_OPTIONS" class="w-full" />
            </UFormField>
          </div>
          <div class="flex justify-end gap-2 pt-6">
            <UButton color="neutral" variant="outline" @click="editModal = false">{{ $t('common.cancel') }}</UButton>
            <UButton color="primary" :loading="editLoading" @click="handleEdit">{{ $t('common.save') }}</UButton>
          </div>
        </div>
      </template>
    </UModal>

    <!-- Create Modal -->
    <UModal v-model:open="createModal" :ui="{ content: 'max-w-lg' }">
      <template #content>
        <div class="p-6">
          <h3 class="text-lg font-semibold text-highlighted mb-5">
            {{ $t("admin.announcements.create_title") }}
          </h3>
          <UForm :schema="schema" :state="state" class="space-y-4" @submit="handleCreate">
            <UFormField :label="$t('admin.announcements.field_title')" name="title">
              <UInput
                v-model="state.title"
                :placeholder="$t('admin.announcements.field_title')"
                class="w-full" />
            </UFormField>
            <UFormField :label="$t('admin.announcements.field_content')" name="content">
              <UTextarea
                v-model="state.content"
                :rows="4"
                :placeholder="$t('admin.announcements.field_content')"
                class="w-full" />
            </UFormField>
            <UFormField :label="$t('admin.announcements.field_type')" name="type">
              <USelect v-model="state.type" :items="TYPE_OPTIONS" class="w-full" />
            </UFormField>
            <div class="flex justify-end gap-2 pt-2">
              <UButton color="neutral" variant="outline" @click="createModal = false">{{ $t('common.cancel') }}</UButton>
              <UButton type="submit" color="primary" icon="i-tabler-send" :loading="creating">
                {{ $t("admin.announcements.btn_create") }}
              </UButton>
            </div>
          </UForm>
        </div>
      </template>
    </UModal>

    <!-- Delete Modal -->
    <UModal v-model:open="deleteModal">
      <template #content>
        <div class="p-6">
          <div class="flex items-center gap-3 mb-4">
            <div class="size-10 rounded-full bg-error/10 flex items-center justify-center shrink-0">
              <UIcon name="i-tabler-alert-triangle" class="size-5 text-error" />
            </div>
            <div>
              <h3 class="font-semibold text-highlighted">{{ $t('common.confirm_delete') }}</h3>
              <p class="text-sm text-muted mt-0.5">{{ deleteTarget?.title }}</p>
            </div>
          </div>
          <div class="flex justify-end gap-2 mt-6">
            <UButton color="neutral" variant="outline" @click="deleteModal = false">{{ $t('common.cancel') }}</UButton>
            <UButton color="error" :loading="deleteLoading" @click="confirmDelete">{{ $t('common.delete') }}</UButton>
          </div>
        </div>
      </template>
    </UModal>
  </AdminPageContainer>
</template>

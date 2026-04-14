<template>
  <AdminPageContainer>
    <AdminPageHeader :title="$t('admin.media.title')" :subtitle="$t('admin.media.subtitle')">
      <template #actions>
        <UButton color="primary" icon="i-tabler-upload" @click="showUploadModal = true">
          <span class="hidden sm:block">{{ $t('admin.media.upload_file') }}</span>
        </UButton>
      </template>
    </AdminPageHeader>

    <AdminPageContent>
      <!-- Type tabs -->
      <div class="flex items-center gap-1 border-b border-default pb-0 mb-4 overflow-x-auto">
        <button
          v-for="t in typeTabs"
          :key="t.value"
          class="px-3 py-2 text-sm font-medium rounded-t transition-colors whitespace-nowrap"
          :class="filterType === t.value
            ? 'text-primary border-b-2 border-primary'
            : 'text-muted hover:text-highlighted'"
          @click="onTypeTab(t.value)">
          {{ t.label }}
          <span v-if="t.count != null" class="ml-1 text-xs text-muted">({{ t.count }})</span>
        </button>
      </div>

      <!-- Filters toolbar -->
      <div class="flex flex-wrap items-center gap-3 mb-4">
        <div class="flex flex-col gap-0.5">
          <UInput
            v-model="searchKeyword"
            :placeholder="$t('admin.media.search_placeholder')"
            leading-icon="i-tabler-search"
            class="w-56" size="sm">
            <template v-if="searchKeyword" #trailing>
              <UButton icon="i-tabler-x" color="neutral" variant="ghost" size="xs" @click="searchKeyword = ''" />
            </template>
          </UInput>
          <span v-if="searchKeyword" class="text-xs text-muted pl-1">{{ $t('admin.media.search_current_page') }}</span>
        </div>

        <AdminSearchableSelect
          v-if="showCategoryTabs"
          v-model="filterCategory"
          :items="categorySelectItems"
          :placeholder="$t('admin.media.filter_category')"
          :search-placeholder="$t('common.search')" />

        <UButton
          :color="filterExternal ? 'warning' : 'neutral'"
          :variant="filterExternal ? 'soft' : 'outline'"
          size="sm" icon="i-tabler-link"
          @click="toggleFilterExternal">
          {{ $t('admin.media.filter_external') }}
        </UButton>

        <USelect
          v-model="pageSize"
          :items="[
            { label: $t('common.items_per_page', { n: 20 }), value: 20 },
            { label: $t('common.items_per_page', { n: 48 }), value: 48 },
            { label: $t('common.items_per_page', { n: 100 }), value: 100 },
          ]"
          class="w-28" size="sm" />

        <div class="ml-auto flex items-center gap-2">
          <ViewModeSwitcher
            v-model="viewMode"
            :modes="[
              { value: 'grid', title: $t('admin.media.grid_view') },
              { value: 'list', title: $t('admin.media.list_view') },
            ]" />
        </div>
      </div>

      <!-- Content area -->
      <div class="flex-1">
        <!-- Loading -->
        <div
          v-if="displayLoading || !initialized"
          class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 2xl:grid-cols-8 gap-4">
          <div v-for="i in 16" :key="i" class="rounded-md overflow-hidden border border-default">
            <USkeleton class="aspect-square w-full" />
            <div class="p-3 space-y-1.5">
              <USkeleton class="h-3 w-full" />
              <USkeleton class="h-3 w-2/3" />
            </div>
          </div>
        </div>

        <!-- Grid view -->
        <MediaGridView
          v-else-if="viewMode === 'grid' && filteredMedias.length > 0"
          :medias="filteredMedias"
          :selected-ids="selectedIds"
          :variants-cache="variantsCache"
          @view="viewMediaDetails"
          @copy-url="copyUrl"
          @delete="confirmDelete"
          @toggle-select="toggleSelect" />

        <!-- List view -->
        <MediaListView
          v-else-if="viewMode === 'list' && filteredMedias.length > 0"
          :medias="filteredMedias"
          :selected-ids="selectedIds"
          :variants-cache="variantsCache"
          :actions="getMediaActions"
          @toggle-select="toggleSelect" />

        <!-- Empty state -->
        <div v-else-if="initialized && !displayLoading" class="flex flex-col items-center justify-center py-16">
          <UIcon
            :name="searchKeyword ? 'i-tabler-search-off' : 'i-tabler-photo-off'"
            class="size-16 text-muted mb-4" />
          <h3 class="text-lg font-medium text-highlighted mb-1">
            {{ searchKeyword ? $t('admin.media.no_results') : $t('admin.media.no_files') }}
          </h3>
          <p class="text-sm text-muted mb-4">
            {{ searchKeyword ? $t('admin.media.search_placeholder') : $t('admin.media.upload_area_desc') }}
          </p>
          <UButton v-if="searchKeyword" color="neutral" variant="outline" icon="i-tabler-filter-off" @click="searchKeyword = ''">
            {{ $t('posts.clear_filters') }}
          </UButton>
          <UButton v-else color="primary" @click="showUploadModal = true">{{ $t('admin.media.upload_file') }}</UButton>
        </div>
      </div>
    </AdminPageContent>

    <AdminPageFooter>
      <template #left>
        <template v-if="medias.length > 0">
          <UCheckbox
            :model-value="isAllSelected"
            :indeterminate="isIndeterminate"
            @update:model-value="toggleSelectAll" />
          <template v-if="selectedIds.length > 0">
            <span>{{ $t('admin.media.selected_n', { n: selectedIds.length }) }}</span>
            <USeparator orientation="vertical" class="h-4" />
            <UButton
              v-if="selectedExternalIds.length > 0"
              color="primary" variant="soft" size="sm" icon="i-tabler-download"
              :loading="batchLocalizing"
              @click="doBatchLocalize">
              {{ $t('admin.media.localize_selected', { n: selectedExternalIds.length }) }}
            </UButton>
            <UButton color="error" variant="soft" size="sm" icon="i-tabler-trash" @click="showBatchDeleteModal = true">
              {{ $t('admin.media.delete_selected') }}
            </UButton>
            <UButton color="neutral" variant="ghost" size="sm" @click="selectedIds = []">{{ $t('common.cancel') }}</UButton>
          </template>
          <span v-else class="text-xs">{{ $t('common.selectAll') }}</span>
        </template>
        <span v-else-if="initialized && !displayLoading && total > 0">
          {{ $t('admin.media.selected_n', { n: total }) }} · {{ currentPage }}/{{ totalPages }}
        </span>
      </template>
      <template #right>
        <UPagination
          v-if="initialized && !displayLoading && total > pageSize"
          v-model:page="currentPage"
          :total="total"
          :items-per-page="pageSize"
          size="sm" />
      </template>
    </AdminPageFooter>

    <!-- Upload modal -->
    <MediaUploadModal
      v-model:open="showUploadModal"
      :category-options="categoryOptions"
      @uploaded="onUploaded" />

    <!-- Detail modal -->
    <MediaDetailModal
      :open="showDetailModal"
      :media="viewingMedia"
      @update:open="handleDetailModalUpdate"
      @delete="confirmDelete"
      @saved="onDetailSaved"
      @open-lightbox="openLightbox" />

    <!-- Single delete modal -->
    <UModal v-model:open="showDeleteModal">
      <template #content>
        <div class="p-6">
          <div class="flex items-center gap-3 mb-4">
            <div class="size-10 rounded-full bg-error/10 flex items-center justify-center shrink-0">
              <UIcon name="i-tabler-alert-triangle" class="size-5 text-error" />
            </div>
            <div>
              <h3 class="font-semibold text-highlighted">{{ $t('admin.media.delete_confirm_title') }}</h3>
              <p class="text-sm text-muted mt-0.5">{{ $t('common.cannot_undo') }}</p>
            </div>
          </div>
          <div class="flex justify-end gap-2 mt-6">
            <UButton color="neutral" variant="outline" @click="showDeleteModal = false">{{ $t('common.cancel') }}</UButton>
            <UButton color="error" :loading="deleteLoading" @click="doDelete">{{ $t('common.delete') }}</UButton>
          </div>
        </div>
      </template>
    </UModal>

    <!-- Lightbox -->
    <MediaLightbox
      v-model:open="lightboxOpen"
      :media="lightboxMedia"
      :images="galleryImages"
      :index="lightboxIndex"
      :has-prev="hasPrev"
      :has-next="hasNext"
      :page-loading="lightboxPageLoading"
      :current-page="currentPage"
      :total-pages="totalPages"
      @prev="lightboxPrev"
      @next="lightboxNext"
      @select-index="i => { lightboxIndex = i }" />

    <!-- Batch delete modal -->
    <UModal v-model:open="showBatchDeleteModal">
      <template #content>
        <div class="p-6">
          <div class="flex items-center gap-3 mb-4">
            <div class="size-10 rounded-full bg-error/10 flex items-center justify-center shrink-0">
              <UIcon name="i-tabler-trash-x" class="size-5 text-error" />
            </div>
            <div>
              <h3 class="font-semibold text-highlighted">{{ $t('admin.media.confirm_delete') }}</h3>
              <p class="text-sm text-muted mt-0.5">{{ $t('admin.media.confirm_delete_desc', { n: selectedIds.length }) }}</p>
            </div>
          </div>
          <div class="flex justify-end gap-2 mt-6">
            <UButton color="neutral" variant="outline" @click="showBatchDeleteModal = false">{{ $t('common.cancel') }}</UButton>
            <UButton color="error" :loading="batchDeleteLoading" @click="doBatchDelete">{{ $t('common.confirm_delete') }}</UButton>
          </div>
        </div>
      </template>
    </UModal>
  </AdminPageContainer>
</template>

<script setup lang="ts">
import type { MediaResponse } from "~/types/api/media";

const { apiFetch } = useApiFetch();
const toast = useToast();
const route = useRoute();
const router = useRouter();
const mediaStore = useMediaStore();
const mediaApi = useMediaApi();
const { t } = useI18n();

// ── State ──────────────────────────────────────────────────────────────────
const { loading } = storeToRefs(mediaStore);
const displayLoading = useMinLoading(loading);
const medias = computed(() => mediaStore.medias);
const total = computed(() => mediaStore.pagination.total);
const totalPages = computed(() => Math.ceil(total.value / pageSize.value));

const filterType = ref((route.query.type as string) || "all");
const filterCategory = ref((route.query.category as string) || "all");
const filterExternal = ref(route.query.external === '1');
const searchKeyword = ref((route.query.q as string) || "");
const viewMode = ref<"grid" | "list">("grid");
const currentPage = ref(Number(route.query.page) || 1);
const pageSize = ref(Number(route.query.size) || 20);
const selectedIds = ref<number[]>([]);

const initialized = ref(false);
const showUploadModal = ref(false);
const showDetailModal = ref(false);
const lightboxOpen = ref(false);
const lightboxIndex = ref(0);
const showDeleteModal = ref(false);
const showBatchDeleteModal = ref(false);
const viewingMedia = ref<MediaResponse | null>(null);
const deleteLoading = ref(false);
const batchDeleteLoading = ref(false);
const batchLocalizing = ref(false);
const pendingDeleteMedia = ref<MediaResponse | null>(null);

// ── Helpers ────────────────────────────────────────────────────────────────
const getMediaType = (mimeType: string): string => {
  if (mimeType.startsWith("image/")) return "image";
  if (mimeType.startsWith("video/")) return "video";
  if (mimeType.startsWith("audio/")) return "audio";
  return "other";
};

const variantsCache = computed(() => {
  const map = new Map<number, Record<string, string>>();
  for (const m of medias.value) {
    if (!m.variants) continue;
    try { map.set(m.id, JSON.parse(m.variants)); } catch {}
  }
  return map;
});

// ── Stats ──────────────────────────────────────────────────────────────────
interface MediaStats { total: number; image: number; video: number; audio: number; other: number }
const stats = ref<MediaStats>({ total: 0, image: 0, video: 0, audio: 0, other: 0 });
const fetchStats = async () => { try { stats.value = await apiFetch<MediaStats>("/medias/stats"); } catch {} };

const typeTabs = computed(() => [
  { value: "all",   label: t('admin.media.tab_all'),    count: stats.value.total },
  { value: "image", label: t('admin.media.tab_images'), count: stats.value.image },
  { value: "video", label: t('admin.media.tab_videos'), count: stats.value.video },
  { value: "audio", label: t('admin.media.tab_audio'),  count: stats.value.audio },
  { value: "other", label: t('admin.media.tab_others'), count: stats.value.other },
]);

const { categories: mediaCats, load: loadMediaCats, getCategoryLabel } = useMediaCategories();

const categorySelectItems = computed(() => [
  { value: "all", label: t('common.all') },
  ...mediaCats.value.map(c => ({ value: c.slug, label: getCategoryLabel(c.slug) })),
]);

const categoryOptions = computed(() => categorySelectItems.value.filter(i => i.value !== 'all'));

// ── Filtered list ──────────────────────────────────────────────────────────
const filteredMedias = computed(() => {
  let result = medias.value;
  if (filterType.value !== "all") result = result.filter(m => getMediaType(m.mime_type) === filterType.value);
  if (searchKeyword.value.trim()) {
    const q = searchKeyword.value.toLowerCase();
    result = result.filter(m => m.filename.toLowerCase().includes(q));
  }
  return result;
});

// ── Selection ──────────────────────────────────────────────────────────────
const isAllSelected = computed(() => filteredMedias.value.length > 0 && filteredMedias.value.every(m => selectedIds.value.includes(m.id)));
const isIndeterminate = computed(() => selectedIds.value.length > 0 && !isAllSelected.value);
const toggleSelectAll = () => { selectedIds.value = isAllSelected.value ? [] : filteredMedias.value.map(m => m.id) };
const toggleSelect = (id: number) => {
  const idx = selectedIds.value.indexOf(id);
  if (idx > -1) selectedIds.value.splice(idx, 1);
  else selectedIds.value.push(id);
};

const showCategoryTabs = computed(() => filterType.value === 'all' || filterType.value === 'image');

// ── Fetch ──────────────────────────────────────────────────────────────────
const fetchMedias = async () => {
  const mimePrefix = filterType.value !== "all" ? filterType.value + "/" : undefined;
  await mediaStore.fetchMedias({
    page: currentPage.value,
    size: pageSize.value,
    ...(mimePrefix && { mime_type: mimePrefix }),
    ...(filterCategory.value !== "all" && { category: filterCategory.value as any }),
    ...(filterExternal.value && { storage_type: 5 }),
  });
};

const onTypeTab = (val: string) => {
  filterType.value = val;
  if (val !== 'all' && val !== 'image') filterCategory.value = 'all';
  currentPage.value = 1;
  selectedIds.value = [];
};

const toggleFilterExternal = () => {
  filterExternal.value = !filterExternal.value;
  currentPage.value = 1;
  selectedIds.value = [];
};

const debouncedKeyword = refDebounced(searchKeyword, 350);

watch([filterType, filterCategory, filterExternal, pageSize, debouncedKeyword], () => {
  selectedIds.value = [];
  if (currentPage.value !== 1) currentPage.value = 1;
  else fetchMedias();
});
watch(currentPage, fetchMedias);

watch([filterType, filterCategory, filterExternal, debouncedKeyword, currentPage, pageSize], () => {
  const q: Record<string, string> = {};
  if (filterType.value !== "all") q.type = filterType.value;
  if (filterCategory.value !== "all") q.category = filterCategory.value;
  if (filterExternal.value) q.external = '1';
  if (debouncedKeyword.value) q.q = debouncedKeyword.value;
  if (currentPage.value > 1) q.page = String(currentPage.value);
  if (pageSize.value !== 20) q.size = String(pageSize.value);
  router.replace({ query: q });
});

onMounted(async () => { await fetchMedias(); initialized.value = true; fetchStats(); loadMediaCats(); });

const onUploaded = () => { fetchMedias(); fetchStats(); };

// ── Lightbox ───────────────────────────────────────────────────────────────
const galleryImages = computed(() => filteredMedias.value.filter(m => getMediaType(m.mime_type) === 'image'));
const lightboxPageLoading = ref(false);

const openLightbox = (media: MediaResponse) => {
  const idx = galleryImages.value.findIndex(m => m.id === media.id);
  lightboxIndex.value = idx >= 0 ? idx : 0;
  showDetailModal.value = false;
  lightboxOpen.value = true;
};

watch(lightboxOpen, (open) => {
  if (!open && viewingMedia.value) showDetailModal.value = true;
});

const lightboxMedia = computed(() => galleryImages.value[lightboxIndex.value] ?? null);

const lightboxPrev = async () => {
  if (lightboxIndex.value > 0) { lightboxIndex.value--; }
  else if (currentPage.value > 1) {
    lightboxPageLoading.value = true;
    currentPage.value--;
    await fetchMedias();
    lightboxPageLoading.value = false;
    lightboxIndex.value = galleryImages.value.length - 1;
  }
};

const lightboxNext = async () => {
  if (lightboxIndex.value < galleryImages.value.length - 1) { lightboxIndex.value++; }
  else if (currentPage.value < totalPages.value) {
    lightboxPageLoading.value = true;
    currentPage.value++;
    await fetchMedias();
    lightboxPageLoading.value = false;
    lightboxIndex.value = 0;
  }
};

const hasPrev = computed(() => lightboxIndex.value > 0 || currentPage.value > 1);
const hasNext = computed(() => lightboxIndex.value < galleryImages.value.length - 1 || currentPage.value < totalPages.value);

// ── Detail modal ───────────────────────────────────────────────────────────
const handleDetailModalUpdate = (isOpen: boolean) => {
  if (!isOpen && lightboxOpen.value) return;
  showDetailModal.value = isOpen;
  if (!isOpen) lightboxOpen.value = false;
};

const viewMediaDetails = (media: MediaResponse) => {
  viewingMedia.value = { ...media };
  lightboxOpen.value = false;
  showDetailModal.value = true;
};

const onDetailSaved = () => {
  showDetailModal.value = false;
  viewingMedia.value = null;
  fetchMedias();
};

// ── Delete ─────────────────────────────────────────────────────────────────
const confirmDelete = (media: MediaResponse) => {
  pendingDeleteMedia.value = media;
  showDeleteModal.value = true;
};

const doDelete = async () => {
  if (!pendingDeleteMedia.value) return;
  deleteLoading.value = true;
  try {
    await mediaStore.deleteMedia(pendingDeleteMedia.value.id);
    toast.add({ title: t('admin.media.deleted'), icon: "i-tabler-trash", color: "neutral" });
    showDeleteModal.value = false;
    showDetailModal.value = false;
    viewingMedia.value = null;
    pendingDeleteMedia.value = null;
    fetchMedias();
    fetchStats();
  } catch (err) {
    toast.add({ title: err instanceof Error ? err.message : t('admin.media.delete_failed'), color: "error", icon: "i-tabler-alert-circle" });
  } finally {
    deleteLoading.value = false;
  }
};

const doBatchDelete = async () => {
  batchDeleteLoading.value = true;
  try {
    await mediaStore.batchDeleteMedias([...selectedIds.value]);
    toast.add({ title: t('admin.media.deleted'), icon: "i-tabler-trash-x", color: "success" });
    selectedIds.value = [];
    showBatchDeleteModal.value = false;
    fetchStats();
  } catch (err) {
    toast.add({ title: err instanceof Error ? err.message : t('admin.media.delete_failed'), color: "error", icon: "i-tabler-alert-circle" });
  } finally {
    batchDeleteLoading.value = false;
  }
};

// ── Batch localize ─────────────────────────────────────────────────────────
const selectedExternalIds = computed(() =>
  selectedIds.value.filter(id => medias.value.find(m => m.id === id)?.storage_type === 5),
);

const doBatchLocalize = async () => {
  const ids = [...selectedExternalIds.value];
  if (ids.length === 0) return;
  batchLocalizing.value = true;
  let success = 0, failed = 0;
  for (const id of ids) {
    try { await mediaApi.localize(id); success++; } catch { failed++; }
  }
  batchLocalizing.value = false;
  if (success > 0) {
    toast.add({ title: t('admin.media.localize_batch_done', { n: success }), icon: 'i-tabler-circle-check', color: 'success' });
    selectedIds.value = selectedIds.value.filter(id => !ids.includes(id));
    fetchMedias();
    fetchStats();
  }
  if (failed > 0) {
    toast.add({ title: t('admin.media.localize_batch_failed', { n: failed }), color: 'error', icon: 'i-tabler-alert-circle' });
  }
};

// ── Localize single ────────────────────────────────────────────────────────
const localizingIds = ref<Set<number>>(new Set());

const doLocalize = async (media: MediaResponse) => {
  if (localizingIds.value.has(media.id)) return;
  localizingIds.value = new Set([...localizingIds.value, media.id]);
  try {
    await mediaApi.localize(media.id);
    toast.add({ title: t('admin.media.localize_success'), icon: 'i-tabler-circle-check', color: 'success' });
    fetchMedias();
    fetchStats();
  } catch (err) {
    toast.add({ title: err instanceof Error ? err.message : t('admin.media.localize_failed'), color: 'error', icon: 'i-tabler-alert-circle' });
  } finally {
    const next = new Set(localizingIds.value);
    next.delete(media.id);
    localizingIds.value = next;
  }
};

// ── Actions dropdown ───────────────────────────────────────────────────────
const getMediaActions = (media: MediaResponse) => [
  [
    { label: t("admin.media.view_detail"), icon: "i-tabler-eye", onSelect: () => viewMediaDetails(media) },
    { label: t("admin.media.copy_url"), icon: "i-tabler-copy", onSelect: () => copyUrl(media) },
    ...(media.storage_type === 5 ? [{
      label: localizingIds.value.has(media.id) ? t("admin.media.localizing") : t("admin.media.localize"),
      icon: "i-tabler-download",
      disabled: localizingIds.value.has(media.id),
      onSelect: () => doLocalize(media),
    }] : []),
  ],
  [{ label: t("common.delete"), icon: "i-tabler-trash", color: "error" as const, onSelect: () => confirmDelete(media) }],
];

// ── Copy URL ───────────────────────────────────────────────────────────────
const copyUrl = async (media: MediaResponse) => {
  try {
    await navigator.clipboard.writeText(media.cdn_url);
  } catch {
    const input = document.createElement("input");
    input.value = media.cdn_url;
    document.body.appendChild(input);
    input.select();
    document.execCommand("copy");
    document.body.removeChild(input);
  }
  toast.add({ title: t('admin.media.url_copied'), icon: "i-tabler-copy", color: "success" });
};
</script>

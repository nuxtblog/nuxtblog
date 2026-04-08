<script setup lang="ts">
import type { MarketplaceItem, PluginPreviewInfo } from "~/composables/usePluginApi";

const { t } = useI18n();
const pluginApi = usePluginApi();
const toast = useToast();

// ── State ──────────────────────────────────────────────────────────────────
const marketItems = ref<MarketplaceItem[]>([]);
const rawLoading = ref(true);
const loading = useMinLoading(rawLoading);
const syncedAt = ref('');
const syncing = ref(false);
const installingName = ref<string | null>(null);

// ── Preview / install modal ────────────────────────────────────────────────
const previewModal = ref(false);
const previewItem = ref<MarketplaceItem | null>(null);
const previewData = ref<PluginPreviewInfo | null>(null);
const previewLoading = ref(false);
const previewInstalling = ref(false);

const openPreview = async (item: MarketplaceItem) => {
  previewItem.value = item;
  previewData.value = null;
  previewLoading.value = true;
  previewModal.value = true;
  try {
    previewData.value = await pluginApi.preview(item.repo);
  } catch {
    // Preview failed — still show modal with basic info from marketplace
    previewData.value = {
      name: item.name,
      title: item.title,
      description: item.description,
      version: item.version,
      author: item.author,
      icon: item.icon,
      priority: 10,
      has_css: false,
      capabilities: {},
      settings: [],
      webhooks: [],
      pipelines: [],
    };
  } finally {
    previewLoading.value = false;
  }
};

const confirmInstall = async () => {
  if (!previewItem.value) return;
  previewInstalling.value = true;
  installingName.value = previewItem.value.name;
  try {
    const res = await pluginApi.install(previewItem.value.repo);
    installedIds.value.add(res.item.id);
    toast.add({ title: t('admin.plugins.install_success'), color: 'success' });
    previewModal.value = false;
  } catch (e: any) {
    toast.add({ title: e?.message, color: 'error' });
  } finally {
    previewInstalling.value = false;
    installingName.value = null;
  }
};

const search = ref('');
const filterType = ref('all');

// ── Installed IDs (for "already installed" badge) ──────────────────────────
const installedIds = ref<Set<string>>(new Set());

const loadInstalledIds = async () => {
  try {
    const res = await pluginApi.list();
    installedIds.value = new Set((res.items ?? []).map(p => p.id));
  } catch {}
};

// ── Type options ───────────────────────────────────────────────────────────
const typeOptions = computed(() => [
  { label: t('admin.plugins.market_type_all'), value: 'all' },
  { label: t('admin.plugins.market_type_hook'), value: 'hook' },
  { label: t('admin.plugins.market_type_integration'), value: 'integration' },
  { label: t('admin.plugins.market_type_theme'), value: 'theme' },
  { label: t('admin.plugins.market_type_editor'), value: 'editor' },
  { label: t('admin.plugins.market_type_analytics'), value: 'analytics' },
  { label: t('admin.plugins.market_type_moderation'), value: 'moderation' },
]);

// ── Fetch ──────────────────────────────────────────────────────────────────
const fetchMarketplace = async () => {
  rawLoading.value = true;
  try {
    const res = await pluginApi.marketplace(search.value || undefined, filterType.value !== 'all' ? filterType.value : undefined);
    marketItems.value = res.items ?? [];
    syncedAt.value = res.synced_at ?? '';
  } catch (e: any) {
    toast.add({ title: e?.message, color: 'error' });
  } finally {
    rawLoading.value = false;
  }
};

const doSync = async () => {
  syncing.value = true;
  try {
    const res = await pluginApi.syncMarketplace();
    toast.add({ title: t('admin.plugins.market_sync_done', { n: res.count }), color: 'success' });
    syncedAt.value = res.synced_at;
    await fetchMarketplace();
  } catch (e: any) {
    toast.add({ title: e?.message, color: 'error' });
  } finally {
    syncing.value = false;
  }
};

// kept for direct install (no preview), not used from UI currently
const installFromMarket = (item: MarketplaceItem) => openPreview(item);

const formatSyncedAt = (s: string) => s ? new Date(s).toLocaleString() : '';

// ── Watchers ───────────────────────────────────────────────────────────────
watchDebounced(search, () => fetchMarketplace(), { debounce: 300 });
watch(filterType, () => fetchMarketplace());

onMounted(() => {
  fetchMarketplace();
  loadInstalledIds();
});
</script>

<template>
  <AdminPageContainer>
    <!-- Preview / Install modal -->
    <UModal v-model:open="previewModal" :ui="{ content: 'max-w-md' }">
      <template #content>
        <div class="p-6">
          <h3 class="text-lg font-semibold text-highlighted mb-1">{{ $t('admin.plugins.preview_confirm_title') }}</h3>
          <p class="text-xs text-muted mb-4">{{ previewItem?.name }}</p>

          <!-- Loading skeleton -->
          <div v-if="previewLoading" class="space-y-3">
            <div class="flex items-start gap-3">
              <USkeleton class="h-12 w-12 rounded-md shrink-0" />
              <div class="flex-1 space-y-1.5">
                <USkeleton class="h-4 w-40" />
                <USkeleton class="h-3 w-24" />
                <USkeleton class="h-3 w-full" />
              </div>
            </div>
            <USkeleton class="h-28 w-full rounded-md" />
            <div class="grid grid-cols-3 gap-2">
              <USkeleton v-for="i in 3" :key="i" class="h-14 rounded-md" />
            </div>
          </div>

          <!-- Preview data -->
          <PluginPreviewPanel v-else-if="previewData" :info="previewData" />

          <!-- Type badge from marketplace if known -->
          <div v-if="previewItem?.type && !previewLoading" class="mt-3 flex items-center gap-2">
            <UBadge
              :label="$t(`admin.plugins.market_type_${previewItem.type}`)"
              color="neutral" variant="outline" size="sm" />
            <UBadge
              v-if="previewItem.is_official"
              :label="$t('admin.plugins.market_official')"
              leading-icon="i-tabler-rosette-discount-check"
              color="primary" variant="soft" size="sm" />
            <span v-if="previewItem.license" class="text-xs text-muted">{{ previewItem.license }}</span>
          </div>

          <div class="flex justify-end gap-2 mt-5">
            <UButton color="neutral" variant="ghost" @click="previewModal = false">{{ $t('common.cancel') }}</UButton>
            <UButton color="primary" leading-icon="i-tabler-download" :loading="previewInstalling" :disabled="previewLoading" @click="confirmInstall">
              {{ $t('admin.plugins.install') }}
            </UButton>
          </div>
        </div>
      </template>
    </UModal>

    <AdminPageHeader :title="$t('admin.plugins.market_title')" :subtitle="$t('admin.plugins.market_subtitle')">
      <template #actions>
        <div class="flex items-center gap-2">
          <span v-if="syncedAt" class="text-xs text-muted hidden sm:block">
            {{ $t('admin.plugins.market_synced_at', { time: formatSyncedAt(syncedAt) }) }}
          </span>
          <UButton color="neutral" variant="outline" icon="i-tabler-refresh" :loading="syncing" @click="doSync">
            {{ $t('admin.plugins.market_sync') }}
          </UButton>
        </div>
      </template>
    </AdminPageHeader>

    <AdminPageContent>
      <!-- Toolbar -->
      <div class="flex flex-wrap items-center gap-3 mb-4">
        <UInput
          v-model="search"
          :placeholder="$t('admin.plugins.search_placeholder')"
          leading-icon="i-tabler-search"
          class="w-56"
          size="sm">
          <template v-if="search" #trailing>
            <UButton icon="i-tabler-x" color="neutral" variant="ghost" size="xs" @click="search = ''" />
          </template>
        </UInput>
        <USelect v-model="filterType" :items="typeOptions" class="w-40" size="sm" />
      </div>

      <!-- Loading -->
      <div v-if="loading" class="space-y-3">
        <div v-for="i in 6" :key="i" class="flex items-center gap-4 p-4 bg-default border border-default rounded-md">
          <USkeleton class="h-12 w-12 rounded-md shrink-0" />
          <div class="flex-1 space-y-2">
            <USkeleton class="h-4 w-48" /><USkeleton class="h-3 w-full" /><USkeleton class="h-3 w-32" />
          </div>
          <USkeleton class="h-8 w-16 rounded-md shrink-0" />
        </div>
      </div>

      <!-- Empty -->
      <div v-else-if="marketItems.length === 0" class="flex flex-col items-center justify-center py-16">
        <UIcon name="i-tabler-building-store" class="size-16 text-muted mb-4" />
        <h3 class="text-lg font-medium text-highlighted mb-1">{{ $t('admin.plugins.market_empty') }}</h3>
        <p class="text-sm text-muted mb-4">{{ $t('admin.plugins.market_empty_desc') }}</p>
        <UButton color="neutral" variant="outline" icon="i-tabler-refresh" :loading="syncing" @click="doSync">
          {{ $t('admin.plugins.market_sync') }}
        </UButton>
      </div>

      <!-- List -->
      <div v-else class="space-y-3">
        <div
          v-for="item in marketItems" :key="item.name"
          class="flex items-center gap-4 p-4 bg-default border border-default rounded-md hover:shadow-sm transition-all">
          <div class="h-12 w-12 rounded-md bg-elevated flex items-center justify-center shrink-0">
            <UIcon :name="item.icon || 'i-tabler-plug'" class="size-6 text-primary" />
          </div>

          <div class="flex-1 min-w-0">
            <div class="flex items-start justify-between gap-4">
              <div class="flex-1 min-w-0">
                <div class="flex items-center gap-2 mb-1 flex-wrap">
                  <h3 class="text-sm font-semibold text-highlighted">{{ item.title }}</h3>
                  <UBadge :label="`v${item.version}`" color="neutral" variant="soft" size="sm" />
                  <UBadge
                    v-if="item.is_official"
                    :label="$t('admin.plugins.market_official')"
                    leading-icon="i-tabler-rosette-discount-check"
                    color="primary" variant="soft" size="sm" />
                  <UBadge
                    v-if="item.type"
                    :label="$t(`admin.plugins.market_type_${item.type}`)"
                    color="neutral" variant="outline" size="sm" />
                </div>
                <p class="text-xs text-muted mb-1.5 line-clamp-2">{{ item.description }}</p>
                <div class="flex items-center gap-3 flex-wrap">
                  <span class="text-xs text-muted">{{ $t("admin.plugins.author_label") }}{{ item.author }}</span>
                  <div v-if="item.tags?.length" class="flex gap-1 flex-wrap">
                    <span
                      v-for="tag in item.tags.slice(0, 4)" :key="tag"
                      class="text-xs px-1.5 py-0.5 bg-muted rounded text-muted">{{ tag }}</span>
                  </div>
                  <NuxtLink
                    v-if="item.homepage"
                    :to="item.homepage"
                    target="_blank"
                    class="text-xs text-primary hover:underline inline-flex items-center gap-0.5">
                    GitHub <UIcon name="i-tabler-external-link" class="size-3" />
                  </NuxtLink>
                </div>
              </div>

              <div class="shrink-0">
                <UButton
                  v-if="installedIds.has(item.name)"
                  size="sm" color="neutral" variant="soft" leading-icon="i-tabler-check" disabled>
                  {{ $t('admin.plugins.market_installed') }}
                </UButton>
                <UButton
                  v-else
                  size="sm" color="primary" variant="soft" leading-icon="i-tabler-download"
                  :loading="installingName === item.name"
                  @click="installFromMarket(item)">
                  {{ $t('admin.plugins.install') }}
                </UButton>
              </div>
            </div>
          </div>
        </div>
      </div>
    </AdminPageContent>
  </AdminPageContainer>
</template>

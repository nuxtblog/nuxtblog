<script setup lang="ts">
import type { MarketplaceItem, PluginItem, PluginPreviewInfo } from "~/composables/usePluginApi";

const { t } = useI18n();
const pluginApi = usePluginApi();
const optionApi = useOptionApi();
const toast = useToast();

// ── Proxy settings modal ──────────────────────────────────────────────────
const proxyModal = ref(false);
const proxyGithub = ref('');
const proxyHttp = ref('');
const proxySaving = ref(false);

const openProxySettings = async () => {
  proxyModal.value = true;
  try {
    const [gh, hp] = await Promise.all([
      optionApi.getOption('plugin_github_proxy').catch(() => ''),
      optionApi.getOption('plugin_http_proxy').catch(() => ''),
    ]);
    proxyGithub.value = (gh as string) || '';
    proxyHttp.value = (hp as string) || '';
  } catch {}
};

const saveProxySettings = async () => {
  proxySaving.value = true;
  try {
    await Promise.all([
      optionApi.setOption('plugin_github_proxy', proxyGithub.value.trim()),
      optionApi.setOption('plugin_http_proxy', proxyHttp.value.trim()),
    ]);
    toast.add({ title: t('admin.plugins.proxy_save_success'), color: 'success' });
    proxyModal.value = false;
  } catch (e: any) {
    toast.add({ title: e?.message ?? t('admin.plugins.proxy_save_failed'), color: 'error' });
  } finally {
    proxySaving.value = false;
  }
};

// ── State ──────────────────────────────────────────────────────────────────
const marketItems = ref<MarketplaceItem[]>([]);
const rawLoading = ref(true);
const loading = useMinLoading(rawLoading);
const syncedAt = ref('');
const syncing = ref(false);
const installingName = ref<string | null>(null);
const updatingName = ref<string | null>(null);

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
    const res = await pluginApi.install(previewItem.value.repo, previewItem.value.version);
    installedPlugins.value.set(res.item.id, res.item);
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
const filterRuntime = ref('all');

// ── Installed plugins (for status detection) ──────────────────────────────
const installedPlugins = ref<Map<string, PluginItem>>(new Map());

const loadInstalled = async () => {
  try {
    const res = await pluginApi.list();
    installedPlugins.value = new Map((res.items ?? []).map(p => [p.id, p]));
  } catch {}
};

/** Get installed plugin by marketplace name */
const getInstalled = (name: string) => installedPlugins.value.get(name);

/** Check if marketplace item has update available */
const hasUpdate = (item: MarketplaceItem) => {
  const installed = getInstalled(item.name);
  if (!installed) return false;
  return isNewerVersion(item.version, installed.version);
};

/** Update an installed plugin from marketplace */
const updateFromMarket = async (item: MarketplaceItem) => {
  const installed = getInstalled(item.name);
  if (!installed) return;
  updatingName.value = item.name;
  try {
    const res = await pluginApi.update(installed.id);
    installedPlugins.value.set(res.item.id, res.item);
    toast.add({ title: t('admin.plugins.update_success'), color: 'success' });
  } catch (e: any) {
    toast.add({ title: e?.message ?? t('admin.plugins.update_failed'), color: 'error' });
  } finally {
    updatingName.value = null;
  }
};

// ── Type options ───────────────────────────────────────────────────────────
const typeOptions = computed(() => [
  { label: t('admin.plugins.market_type_all'), value: 'all' },
  { label: t('admin.plugins.filter_type_go'), value: 'builtin' },
  { label: t('admin.plugins.filter_type_js'), value: 'js' },
  { label: t('admin.plugins.filter_type_yaml'), value: 'yaml' },
  { label: t('admin.plugins.filter_type_full'), value: 'full' },
]);

const runtimeOptions = computed(() => [
  { label: t('admin.plugins.market_runtime_all'), value: 'all' },
  { label: t('admin.plugins.market_runtime_compiled'), value: 'compiled' },
  { label: t('admin.plugins.market_runtime_interpreted'), value: 'interpreted' },
]);

// ── Filtered items ────────────────────────────────────────────────────────
const filtered = computed(() => {
  let result = marketItems.value;
  if (filterType.value !== 'all') result = result.filter(i => i.type === filterType.value);
  if (filterRuntime.value !== 'all') result = result.filter(i => i.runtime === filterRuntime.value);
  if (search.value) {
    const q = search.value.toLowerCase();
    result = result.filter(i =>
      i.title.toLowerCase().includes(q) ||
      i.description.toLowerCase().includes(q) ||
      i.name.toLowerCase().includes(q) ||
      i.tags?.some(tag => tag.toLowerCase().includes(q))
    );
  }
  return result;
});

// ── Fetch ──────────────────────────────────────────────────────────────────
const fetchMarketplace = async () => {
  rawLoading.value = true;
  try {
    const res = await pluginApi.marketplace();
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

const installFromMarket = (item: MarketplaceItem) => openPreview(item);

const formatSyncedAt = (s: string) => s ? new Date(s).toLocaleString() : '';

// ── Badge helpers ─────────────────────────────────────────────────────────
const typeBadgeColor = (type: string) => {
  switch (type) {
    case 'builtin': return 'info'
    case 'js': return 'warning'
    case 'yaml': return 'info'
    case 'full': return 'success'
    default: return 'neutral'
  }
};

const runtimeBadgeColor = (runtime: string) =>
  runtime === 'compiled' ? 'error' : 'success';

const runtimeIcon = (runtime: string) =>
  runtime === 'compiled' ? 'i-tabler-refresh' : 'i-tabler-bolt';

// ── Watchers ───────────────────────────────────────────────────────────────
watchDebounced(search, () => {}, { debounce: 300 });

onMounted(() => {
  fetchMarketplace();
  loadInstalled();
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

          <!-- Type / runtime / official badges -->
          <div v-if="previewItem && !previewLoading" class="mt-3 flex items-center gap-2 flex-wrap">
            <UBadge
              v-if="previewItem.type"
              :label="$t(`admin.plugins.type_${previewItem.type}`)"
              :color="typeBadgeColor(previewItem.type) as any"
              variant="soft" size="sm" />
            <UBadge
              v-if="previewItem.runtime"
              :label="$t(`admin.plugins.market_runtime_${previewItem.runtime}`)"
              :leading-icon="runtimeIcon(previewItem.runtime)"
              :color="runtimeBadgeColor(previewItem.runtime) as any"
              variant="outline" size="sm" />
            <UBadge
              v-if="previewItem.is_official"
              :label="$t('admin.plugins.market_official')"
              leading-icon="i-tabler-rosette-discount-check"
              color="success" variant="soft" size="sm" />
            <span v-if="previewItem.license" class="text-xs text-muted">{{ previewItem.license }}</span>
          </div>

          <div class="flex justify-end gap-2 mt-5">
            <UButton color="neutral" variant="ghost" @click="previewModal = false">{{ $t('common.cancel') }}</UButton>
            <UButton
              v-if="hasUpdate(previewItem!)"
              color="warning" leading-icon="i-tabler-arrow-up"
              :loading="previewInstalling" :disabled="previewLoading"
              @click="updateFromMarket(previewItem!); previewModal = false">
              {{ $t('admin.plugins.update_btn') }}
            </UButton>
            <UButton
              v-else-if="!getInstalled(previewItem?.name ?? '')"
              color="primary" leading-icon="i-tabler-download"
              :loading="previewInstalling" :disabled="previewLoading"
              @click="confirmInstall">
              {{ $t('admin.plugins.install') }}
            </UButton>
            <UButton v-else color="neutral" variant="soft" disabled leading-icon="i-tabler-check">
              {{ $t('admin.plugins.market_installed') }}
            </UButton>
          </div>
        </div>
      </template>
    </UModal>

    <!-- Proxy settings modal -->
    <UModal v-model:open="proxyModal" :ui="{ content: 'max-w-sm' }">
      <template #content>
        <div class="p-6">
          <h3 class="text-lg font-semibold text-highlighted mb-4">{{ $t('admin.plugins.proxy_settings_title') }}</h3>
          <div class="space-y-4">
            <div>
              <label class="text-sm font-medium text-highlighted mb-1 block">{{ $t('admin.plugins.proxy_github_label') }}</label>
              <UInput v-model="proxyGithub" placeholder="https://ghproxy.net" size="sm" />
              <p class="text-xs text-muted mt-1">{{ $t('admin.plugins.proxy_github_hint') }}</p>
            </div>
            <div>
              <label class="text-sm font-medium text-highlighted mb-1 block">{{ $t('admin.plugins.proxy_http_label') }}</label>
              <UInput v-model="proxyHttp" placeholder="http://127.0.0.1:7890" size="sm" />
              <p class="text-xs text-muted mt-1">{{ $t('admin.plugins.proxy_http_hint') }}</p>
            </div>
          </div>
          <div class="flex justify-end gap-2 mt-5">
            <UButton color="neutral" variant="ghost" @click="proxyModal = false">{{ $t('common.cancel') }}</UButton>
            <UButton color="primary" :loading="proxySaving" @click="saveProxySettings">{{ $t('common.save') }}</UButton>
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
          <UButton color="neutral" variant="ghost" icon="i-tabler-settings" @click="openProxySettings" />
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
        <USelect v-model="filterType" :items="typeOptions" class="w-32" size="sm" />
        <USelect v-model="filterRuntime" :items="runtimeOptions" class="w-36" size="sm" />
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
      <div v-else-if="filtered.length === 0" class="flex flex-col items-center justify-center py-16">
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
          v-for="item in filtered" :key="item.name"
          class="flex items-center gap-4 p-4 bg-default border border-default rounded-md hover:shadow-sm transition-all cursor-pointer"
          @click="openPreview(item)">
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
                    color="success" variant="soft" size="sm" />
                  <UBadge
                    v-if="item.type"
                    :label="$t(`admin.plugins.type_${item.type}`)"
                    :color="typeBadgeColor(item.type) as any"
                    variant="soft" size="sm" />
                  <UBadge
                    v-if="item.runtime"
                    :label="$t(`admin.plugins.market_runtime_${item.runtime}`)"
                    :leading-icon="runtimeIcon(item.runtime)"
                    :color="runtimeBadgeColor(item.runtime) as any"
                    variant="outline" size="sm" />
                  <!-- Update available badge -->
                  <UBadge
                    v-if="hasUpdate(item)"
                    :label="$t('admin.plugins.update_available', { version: item.version })"
                    color="warning" variant="soft" size="sm"
                    icon="i-tabler-arrow-up" />
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
                    class="text-xs text-primary hover:underline inline-flex items-center gap-0.5"
                    @click.stop>
                    GitHub <UIcon name="i-tabler-external-link" class="size-3" />
                  </NuxtLink>
                </div>
              </div>

              <div class="shrink-0" @click.stop>
                <!-- Has update -->
                <UButton
                  v-if="hasUpdate(item)"
                  size="sm" color="warning" variant="soft" leading-icon="i-tabler-arrow-up"
                  :loading="updatingName === item.name"
                  @click="updateFromMarket(item)">
                  {{ $t('admin.plugins.update_btn') }}
                </UButton>
                <!-- Already installed, up to date -->
                <UButton
                  v-else-if="getInstalled(item.name)"
                  size="sm" color="neutral" variant="soft" leading-icon="i-tabler-check" disabled>
                  {{ $t('admin.plugins.market_installed') }}
                </UButton>
                <!-- Not installed -->
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

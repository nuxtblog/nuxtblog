<script setup lang="ts">
import type { PluginItem, MarketplaceItem } from "~/composables/usePluginApi";
import { parseCapabilityBadges } from "~/composables/usePluginApi";
import { resolveTemplate, registerPluginI18n } from '~/composables/usePluginI18n'

const { t, locale } = useI18n();
const capBadges = (item: PluginItem) => parseCapabilityBadges(item.capabilities ?? '{}', t);
/** Resolve i18n template for plugin item text fields. */
const pt = (text: string, pluginId: string) => resolveTemplate(text, pluginId, locale.value);
const pluginApi = usePluginApi();
const toast = useToast();

// ── State ──────────────────────────────────────────────────────────────────
const items = ref<PluginItem[]>([]);
const rawLoading = ref(true);
const loading = useMinLoading(rawLoading);
const search = ref("");
const filterStatus = ref("all");
const filterType = ref("all");
const filterSource = ref("all");

// ── Update check state ───────────────────────────────────────────────────
const marketplaceItems = ref<MarketplaceItem[]>([]);
const updateMap = ref<Record<string, string>>({});

const selectedPlugins = ref<string[]>([]);
const batchPluginAction = ref<string | undefined>(undefined);

const isAllPluginsSelected = computed(() =>
  filtered.value.length > 0 && filtered.value.every(p => selectedPlugins.value.includes(p.id))
);
const isPluginsIndeterminate = computed(() =>
  selectedPlugins.value.length > 0 && !isAllPluginsSelected.value
);

function toggleSelectPlugin(id: string) {
  const idx = selectedPlugins.value.indexOf(id);
  if (idx > -1) selectedPlugins.value.splice(idx, 1);
  else selectedPlugins.value.push(id);
}

function toggleSelectAllPlugins() {
  if (isAllPluginsSelected.value) selectedPlugins.value = [];
  else selectedPlugins.value = filtered.value.map(p => p.id);
}

async function applyBatchPlugins() {
  if (!batchPluginAction.value || !selectedPlugins.value.length) return;
  let restartToastShown = false;
  try {
    if (batchPluginAction.value === 'enable') {
      await Promise.all(selectedPlugins.value.map(id => pluginApi.toggle(id, true)));
      items.value.forEach(p => { if (selectedPlugins.value.includes(p.id)) p.enabled = true; });
    } else if (batchPluginAction.value === 'disable') {
      await Promise.all(selectedPlugins.value.map(id => pluginApi.toggle(id, false)));
      items.value.forEach(p => { if (selectedPlugins.value.includes(p.id)) p.enabled = false; });
    } else if (batchPluginAction.value === 'uninstall') {
      const res = await pluginApi.batchUninstall(selectedPlugins.value);
      items.value = items.value.filter(p => !res.succeeded.includes(p.id));
      if (res.failed?.length) {
        toast.add({ title: t('admin.plugins.batch_uninstall_partial', { n: res.failed.length }), color: 'warning' });
      }
      if (res.need_restart) {
        toast.add({
          title: t('admin.plugins.uninstall_success'),
          description: t('admin.plugins.restart_required_uninstall'),
          color: 'warning',
          duration: 0,
        });
        restartToastShown = true;
      }
    } else if (batchPluginAction.value === 'update') {
      const res = await pluginApi.batchUpdate(selectedPlugins.value);
      if (res.succeeded?.length) await load();
      if (res.failed?.length) {
        toast.add({ title: t('admin.plugins.batch_update_partial', { n: res.failed.length }), color: 'warning' });
      }
    }
    if (!restartToastShown) {
      toast.add({ title: t('admin.plugins.batch_success'), color: 'success' });
    }
    selectedPlugins.value = [];
    batchPluginAction.value = undefined;
  } catch (e: any) {
    toast.add({ title: e?.message, color: 'error' });
  }
}

// ── Modals ──────────────────────────────────────────────────────────────────
const installModal = ref(false);
const settingsModal = ref(false);
const settingsPlugin = ref<PluginItem | null>(null);

const openSettings = (item: PluginItem) => {
  settingsPlugin.value = item;
  settingsModal.value = true;
};

const onInstalled = (item: PluginItem) => {
  items.value.unshift(item);
};

// ── Filters ────────────────────────────────────────────────────────────────
const statusFilterItems = computed(() => [
  { label: t("admin.plugins.filter_all"), value: "all" },
  { label: t("admin.plugins.filter_active"), value: "active" },
  { label: t("admin.plugins.filter_inactive"), value: "inactive" },
]);

const typeFilterItems = computed(() => [
  { label: t("admin.plugins.filter_type_all"), value: "all" },
  { label: t("admin.plugins.filter_type_go"), value: "builtin" },
  { label: t("admin.plugins.filter_type_js"), value: "js" },
  { label: t("admin.plugins.filter_type_yaml"), value: "yaml" },
]);

const sourceFilterItems = computed(() => [
  { label: t("admin.plugins.filter_source_all"), value: "all" },
  { label: t("admin.plugins.filter_source_builtin"), value: "builtin" },
  { label: t("admin.plugins.filter_source_external"), value: "external" },
]);

const filtered = computed(() => {
  let result = items.value;
  if (filterStatus.value === "active") result = result.filter(p => p.enabled);
  if (filterStatus.value === "inactive") result = result.filter(p => !p.enabled);
  if (filterType.value !== "all") result = result.filter(p => p.type === filterType.value);
  if (filterSource.value !== "all") result = result.filter(p => p.source === filterSource.value);
  if (search.value) {
    const q = search.value.toLowerCase();
    result = result.filter(
      p => p.title.toLowerCase().includes(q) || p.description.toLowerCase().includes(q),
    );
  }
  return result;
});

// ── Type badge helpers ───────────────────────────────────────────────────
const typeBadgeColor = (type: string) => {
  switch (type) {
    case 'builtin': return 'info'
    case 'js': return 'warning'
    case 'yaml': return 'info'
    case 'full': return 'success'
    default: return 'neutral'
  }
};
const typeBadgeLabel = (type: string) => t(`admin.plugins.type_${type}`) || type;

const hasUpdate = (item: PluginItem) => updateMap.value[item.id];
const updatableCount = computed(() => Object.keys(updateMap.value).length);

// ── Load ───────────────────────────────────────────────────────────────────
const load = async () => {
  rawLoading.value = true;
  try {
    const res = await pluginApi.list();
    items.value = res.items ?? [];
    // Register i18n messages for template resolution
    for (const item of items.value) {
      if (item.i18n) {
        try { registerPluginI18n(item.id, JSON.parse(item.i18n)) }
        catch { /* ignore */ }
      }
    }
  } catch (e: any) {
    toast.add({ title: e?.message, color: "error" });
  } finally {
    rawLoading.value = false;
  }
};

const checkUpdates = async () => {
  try {
    const res = await pluginApi.marketplace();
    marketplaceItems.value = res.items ?? [];
    const map: Record<string, string> = {};
    for (const mp of marketplaceItems.value) {
      const installed = items.value.find(p => p.id === mp.name);
      if (installed && installed.repo_url && mp.version &&
          isNewerVersion(mp.version, installed.version)) {
        map[installed.id] = mp.version;
      }
    }
    updateMap.value = map;
  } catch {
    // Silently fail
  }
};

onMounted(async () => {
  await load();
  checkUpdates();
});

watch([filterStatus, filterType, filterSource, search], () => {
  selectedPlugins.value = [];
  batchPluginAction.value = undefined;
});

// ── Toggle / Uninstall / Update ────────────────────────────────────────────
const togglingId = ref<string | null>(null);
const uninstallModal = ref(false);
const pendingUninstall = ref<PluginItem | null>(null);
const uninstalling = ref(false);
const updatingId = ref<string | null>(null);

const cascadePlugins = ref<string[]>([]);
const cascadeLoading = ref(false);

const impactHasDB = ref(false);
const impactMediaCats = ref<string[]>([]);

const disableModal = ref(false);
const pendingDisable = ref<PluginItem | null>(null);

const toggle = async (item: PluginItem) => {
  if (item.enabled) {
    try {
      cascadeLoading.value = true;
      const res = await pluginApi.unloadImpact(item.id);
      if (res.will_unload?.length > 0) {
        cascadePlugins.value = res.will_unload;
        pendingDisable.value = item;
        disableModal.value = true;
        return;
      }
    } catch {
    } finally {
      cascadeLoading.value = false;
    }
  }
  await doToggle(item);
};

const doToggle = async (item: PluginItem) => {
  togglingId.value = item.id;
  try {
    await pluginApi.toggle(item.id, !item.enabled);
    item.enabled = !item.enabled;
    if (!item.enabled && cascadePlugins.value.length > 0) {
      for (const depId of cascadePlugins.value) {
        const dep = items.value.find(p => p.id === depId);
        if (dep) dep.enabled = false;
      }
      cascadePlugins.value = [];
    }
  } catch (e: any) {
    toast.add({ title: e?.message, color: "error" });
  } finally {
    togglingId.value = null;
  }
};

const confirmDisable = async () => {
  if (!pendingDisable.value) return;
  disableModal.value = false;
  await doToggle(pendingDisable.value);
  pendingDisable.value = null;
};

const openUninstall = async (item: PluginItem) => {
  pendingUninstall.value = item;
  cascadePlugins.value = [];
  impactHasDB.value = false;
  impactMediaCats.value = [];
  cascadeLoading.value = true;
  uninstallModal.value = true;
  try {
    const res = await pluginApi.unloadImpact(item.id);
    cascadePlugins.value = res.will_unload ?? [];
    impactHasDB.value = res.has_db ?? false;
    impactMediaCats.value = res.media_cat_slugs ?? [];
  } catch {
  } finally {
    cascadeLoading.value = false;
  }
};

const updatePlugin = async (item: PluginItem) => {
  if (!item.repo_url) {
    toast.add({ title: t("admin.plugins.no_repo_url"), color: "warning" });
    return;
  }
  updatingId.value = item.id;
  try {
    const res = await pluginApi.update(item.id);
    const idx = items.value.findIndex(p => p.id === item.id);
    if (idx > -1) items.value[idx] = res.item;
    delete updateMap.value[item.id];
    toast.add({ title: t("admin.plugins.update_success"), color: "success" });
  } catch (e: any) {
    toast.add({ title: e?.message ?? t("admin.plugins.update_failed"), color: "error" });
  } finally {
    updatingId.value = null;
  }
};

const getPluginActions = (item: PluginItem) => {
  const actions: any[][] = [
    [
      {
        label: t("admin.plugins.settings_title"),
        icon: "i-tabler-settings",
        disabled: !item.enabled,
        onSelect: () => openSettings(item),
      },
      ...(item.repo_url ? [{
        label: t("admin.plugins.update_title"),
        icon: "i-tabler-refresh",
        onSelect: () => updatePlugin(item),
      }] : []),
    ],
  ];
  actions.push([
    {
      label: t("admin.plugins.uninstall_title"),
      icon: "i-tabler-trash",
      color: "error" as const,
      onSelect: () => openUninstall(item),
    },
  ]);
  return actions;
};

const confirmUninstall = async () => {
  if (!pendingUninstall.value) return;
  uninstalling.value = true;
  try {
    const res = await pluginApi.uninstall(pendingUninstall.value.id);
    items.value = items.value.filter(p => p.id !== pendingUninstall.value!.id);
    if (res?.need_restart) {
      toast.add({
        title: t("admin.plugins.uninstall_success"),
        description: t("admin.plugins.restart_required_uninstall"),
        color: "warning",
        duration: 0,
      });
    } else {
      toast.add({ title: t("admin.plugins.uninstall_success"), color: "success" });
    }
    uninstallModal.value = false;
  } catch (e: any) {
    toast.add({ title: e?.message, color: "error" });
  } finally {
    uninstalling.value = false;
  }
};
</script>

<template>
  <AdminPageContainer>
    <AdminPageHeader :title="$t('admin.plugins.title')" :subtitle="$t('admin.plugins.subtitle')">
      <template #actions>
        <UButton color="primary" icon="i-tabler-plus" @click="installModal = true">
          {{ $t("admin.plugins.install") }}
        </UButton>
      </template>
    </AdminPageHeader>

    <AdminPageContent>
      <!-- Install modal (extracted component) -->
      <PluginInstallModal v-model:open="installModal" @installed="onInstalled" />

      <!-- Settings modal (extracted component) -->
      <PluginSettingsModal v-model:open="settingsModal" :plugin="settingsPlugin" />

      <!-- Uninstall confirm modal -->
      <UModal v-model:open="uninstallModal" :ui="{ content: 'max-w-sm' }">
        <template #content>
          <div class="p-6">
            <h3 class="text-lg font-semibold text-highlighted mb-2">{{ $t("admin.plugins.uninstall_confirm_title") }}</h3>
            <p class="text-sm text-muted mb-4">{{ $t("admin.plugins.uninstall_confirm_desc", { name: pendingUninstall ? pt(pendingUninstall.title, pendingUninstall.id) : '' }) }}</p>
            <!-- Cascade warning -->
            <div v-if="cascadeLoading" class="mb-4">
              <USkeleton class="h-4 w-48" />
            </div>
            <div v-else-if="cascadePlugins.length > 0" class="mb-4 p-3 bg-warning-50 dark:bg-warning-950/30 border border-warning-200 dark:border-warning-800 rounded-md">
              <div class="flex items-start gap-2">
                <UIcon name="i-tabler-alert-triangle" class="size-4 text-warning-600 dark:text-warning-400 shrink-0 mt-0.5" />
                <div>
                  <p class="text-sm font-medium text-warning-800 dark:text-warning-200 mb-1">{{ $t("admin.plugins.cascade_warning") }}</p>
                  <ul class="text-xs text-warning-700 dark:text-warning-300 space-y-0.5">
                    <li v-for="pid in cascadePlugins" :key="pid" class="flex items-center gap-1">
                      <UIcon name="i-tabler-plug" class="size-3" />
                      {{ (() => { const p = items.find(x => x.id === pid); return p ? pt(p.title, p.id) : pid })() }}
                    </li>
                  </ul>
                </div>
              </div>
            </div>
            <!-- Resource deletion warning -->
            <div v-if="!cascadeLoading && (impactHasDB || impactMediaCats.length > 0)" class="mb-4 p-3 bg-error-50 dark:bg-error-950/30 border border-error-200 dark:border-error-800 rounded-md">
              <div class="flex items-start gap-2">
                <UIcon name="i-tabler-alert-triangle" class="size-4 text-error-600 dark:text-error-400 shrink-0 mt-0.5" />
                <div>
                  <p class="text-sm font-medium text-error-800 dark:text-error-200 mb-1">{{ $t("admin.plugins.uninstall_resource_warning") }}</p>
                  <ul class="text-xs text-error-700 dark:text-error-300 space-y-0.5">
                    <li v-if="impactHasDB" class="flex items-center gap-1">
                      <UIcon name="i-tabler-database" class="size-3" />
                      {{ $t("admin.plugins.uninstall_has_db") }}
                    </li>
                    <li v-if="impactMediaCats.length > 0" class="flex items-center gap-1">
                      <UIcon name="i-tabler-photo" class="size-3" />
                      {{ $t("admin.plugins.uninstall_has_media_cats", { cats: impactMediaCats.join(', ') }) }}
                    </li>
                  </ul>
                </div>
              </div>
            </div>
            <div class="flex justify-end gap-2">
              <UButton color="neutral" variant="ghost" @click="uninstallModal = false">{{ $t("common.cancel") }}</UButton>
              <UButton color="error" :loading="uninstalling" @click="confirmUninstall">{{ $t("admin.plugins.uninstall_title") }}</UButton>
            </div>
          </div>
        </template>
      </UModal>

      <!-- Disable cascade confirm modal -->
      <UModal v-model:open="disableModal" :ui="{ content: 'max-w-sm' }">
        <template #content>
          <div class="p-6">
            <h3 class="text-lg font-semibold text-highlighted mb-2">{{ $t("admin.plugins.disable_confirm_title") }}</h3>
            <p class="text-sm text-muted mb-4">{{ $t("admin.plugins.disable_confirm_desc", { name: pendingDisable ? pt(pendingDisable.title, pendingDisable.id) : '' }) }}</p>
            <div class="mb-4 p-3 bg-warning-50 dark:bg-warning-950/30 border border-warning-200 dark:border-warning-800 rounded-md">
              <div class="flex items-start gap-2">
                <UIcon name="i-tabler-alert-triangle" class="size-4 text-warning-600 dark:text-warning-400 shrink-0 mt-0.5" />
                <div>
                  <p class="text-sm font-medium text-warning-800 dark:text-warning-200 mb-1">{{ $t("admin.plugins.cascade_disable_warning") }}</p>
                  <ul class="text-xs text-warning-700 dark:text-warning-300 space-y-0.5">
                    <li v-for="pid in cascadePlugins" :key="pid" class="flex items-center gap-1">
                      <UIcon name="i-tabler-plug" class="size-3" />
                      {{ (() => { const p = items.find(x => x.id === pid); return p ? pt(p.title, p.id) : pid })() }}
                    </li>
                  </ul>
                </div>
              </div>
            </div>
            <div class="flex justify-end gap-2">
              <UButton color="neutral" variant="ghost" @click="disableModal = false; pendingDisable = null">{{ $t("common.cancel") }}</UButton>
              <UButton color="warning" @click="confirmDisable">{{ $t("admin.plugins.disable_confirm_btn") }}</UButton>
            </div>
          </div>
        </template>
      </UModal>

      <!-- Toolbar -->
      <div class="flex items-center gap-3 mb-4 flex-wrap">
        <UInput v-model="search" :placeholder="$t('admin.plugins.search_placeholder')" leading-icon="i-tabler-search" class="w-56" size="sm">
          <template v-if="search" #trailing>
            <UButton icon="i-tabler-x" color="neutral" variant="ghost" size="xs" @click="search = ''" />
          </template>
        </UInput>
        <USelect v-model="filterStatus" :items="statusFilterItems" class="w-28" size="sm" />
        <USelect v-model="filterType" :items="typeFilterItems" class="w-32" size="sm" />
        <USelect v-model="filterSource" :items="sourceFilterItems" class="w-32" size="sm" />
        <UBadge v-if="updatableCount > 0" :label="`${updatableCount} ${$t('admin.plugins.update_btn')}`" color="warning" variant="soft" size="sm" icon="i-tabler-arrow-up" />
      </div>

      <!-- Loading -->
      <div v-if="loading" class="space-y-3">
        <div v-for="i in 3" :key="i" class="flex items-center gap-4 p-4 bg-default border border-default rounded-md">
          <USkeleton class="h-12 w-12 rounded-md shrink-0" />
          <div class="flex-1 space-y-2">
            <USkeleton class="h-4 w-48" /><USkeleton class="h-3 w-full" /><USkeleton class="h-3 w-32" />
          </div>
          <USkeleton class="h-6 w-10 rounded-full shrink-0" />
        </div>
      </div>

      <!-- Empty -->
      <div v-else-if="filtered.length === 0" class="flex flex-col items-center justify-center py-16">
        <UIcon name="i-tabler-plug-x" class="size-16 text-muted mb-4" />
        <h3 class="text-lg font-medium text-highlighted mb-1">{{ $t("admin.plugins.no_plugins") }}</h3>
        <p class="text-sm text-muted mb-4">{{ $t("admin.plugins.no_plugins_desc") }}</p>
        <UButton color="primary" icon="i-tabler-plus" @click="installModal = true">{{ $t("admin.plugins.install") }}</UButton>
      </div>

      <!-- List -->
      <div v-else class="space-y-3">
        <div
          v-for="item in filtered" :key="item.id"
          class="group flex items-center gap-4 p-4 bg-default border border-default rounded-md hover:shadow-sm transition-all">
          <UCheckbox :model-value="selectedPlugins.includes(item.id)" @update:model-value="toggleSelectPlugin(item.id)" />
          <div class="h-12 w-12 rounded-md bg-elevated flex items-center justify-center shrink-0">
            <UIcon :name="item.icon || 'i-tabler-plug'" class="size-6 text-primary" />
          </div>
          <div class="flex-1 min-w-0">
            <div class="flex items-start justify-between gap-4">
              <div class="flex-1 min-w-0">
                <div class="flex items-center gap-2 mb-1 flex-wrap">
                  <h3 class="text-sm font-semibold text-highlighted">{{ pt(item.title, item.id) || item.id }}</h3>
                  <UBadge :label="`v${item.version}`" color="neutral" variant="soft" size="sm" />
                  <UBadge v-if="item.type" :label="typeBadgeLabel(item.type)" :color="typeBadgeColor(item.type) as any" variant="soft" size="sm" />
                  <UBadge v-if="item.source === 'builtin'" :label="$t('admin.plugins.source_builtin')" color="primary" variant="outline" size="sm" />
                  <UBadge v-else :label="$t('admin.plugins.source_external')" color="neutral" variant="outline" size="sm" />
                  <UBadge
                    v-if="hasUpdate(item)"
                    :label="$t('admin.plugins.update_available', { version: hasUpdate(item) })"
                    color="warning"
                    variant="soft"
                    size="sm"
                    icon="i-tabler-arrow-up" />
                  <UBadge
                    v-for="badge in capBadges(item)"
                    :key="badge.label"
                    :label="badge.label"
                    :color="badge.color as any"
                    variant="soft"
                    size="sm" />
                </div>
                <p class="text-xs text-muted mb-1">{{ pt(item.description, item.id) }}</p>
                <p class="text-xs text-muted">
                  {{ $t("admin.plugins.author_label") }}{{ item.author }}
                  <template v-if="item.repo_url">
                    <span class="mx-1">·</span>
                    <NuxtLink :to="item.repo_url" target="_blank" class="text-primary hover:underline inline-flex items-center gap-0.5">
                      GitHub <UIcon name="i-tabler-external-link" class="size-3" />
                    </NuxtLink>
                  </template>
                </p>
              </div>
              <div class="flex items-center gap-2 shrink-0">
                <UBadge
                  :label="item.enabled ? $t('admin.plugins.status_active') : $t('admin.plugins.status_inactive')"
                  :color="item.enabled ? 'success' : 'neutral'"
                  variant="soft" size="sm" />
                <UButton
                  v-if="hasUpdate(item)"
                  color="warning"
                  variant="soft"
                  size="xs"
                  icon="i-tabler-arrow-up"
                  :loading="updatingId === item.id"
                  @click="updatePlugin(item)">
                  {{ $t('admin.plugins.update_btn') }}
                </UButton>
                <USwitch :model-value="item.enabled" :loading="togglingId === item.id" @update:model-value="toggle(item)" />
                <UDropdownMenu :items="getPluginActions(item)" :popper="{ placement: 'bottom-end' }">
                  <UButton color="neutral" variant="ghost" icon="i-tabler-dots-vertical" square size="xs" class="opacity-0 group-hover:opacity-100 transition-opacity" />
                </UDropdownMenu>
              </div>
            </div>
          </div>
        </div>
      </div>
    </AdminPageContent>

    <AdminPageFooter>
      <template #left>
        <template v-if="filtered.length > 0">
          <UCheckbox :model-value="isAllPluginsSelected" :indeterminate="isPluginsIndeterminate" @update:model-value="toggleSelectAllPlugins" />
          <template v-if="selectedPlugins.length > 0">
            <span>{{ $t('common.selected_n', { n: selectedPlugins.length }) }}</span>
            <USeparator orientation="vertical" class="h-4" />
            <USelect
              v-model="batchPluginAction"
              :items="[
                { label: $t('admin.plugins.batch_enable'), value: 'enable' },
                { label: $t('admin.plugins.batch_disable'), value: 'disable' },
                ...( selectedPlugins.some(id => items.find(p => p.id === id)?.repo_url)
                  ? [{ label: $t('admin.plugins.batch_update'), value: 'update' }] : []),
                { label: $t('admin.plugins.batch_uninstall'), value: 'uninstall' },
              ]"
              :placeholder="$t('admin.posts.batch_action')"
              class="w-36" size="sm" />
            <UButton color="primary" variant="soft" size="sm" :disabled="!batchPluginAction" @click="applyBatchPlugins">{{ $t('common.apply') }}</UButton>
            <UButton color="neutral" variant="ghost" size="sm" @click="selectedPlugins = []; batchPluginAction = undefined">{{ $t('common.cancel') }}</UButton>
          </template>
          <span v-else class="text-xs">{{ $t('common.selectAll') }}</span>
        </template>
      </template>
    </AdminPageFooter>
  </AdminPageContainer>
</template>

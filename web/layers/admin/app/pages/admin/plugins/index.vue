<script setup lang="ts">
import type { PluginItem, PluginSettingField, PluginPreviewInfo, MarketplaceItem } from "~/composables/usePluginApi";
import { parseCapabilityBadges } from "~/composables/usePluginApi";

const { t } = useI18n();
const capBadges = (item: PluginItem) => parseCapabilityBadges(item.capabilities ?? '{}', t);
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
const updateMap = ref<Record<string, string>>({}); // plugin id → latest version

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

// ── Settings modal ─────────────────────────────────────────────────────────
const settingsModal = ref(false);
const settingsPlugin = ref<PluginItem | null>(null);
const settingsSchema = ref<PluginSettingField[]>([]);
const settingsValues = ref<Record<string, unknown>>({});
const settingsLoading = ref(false);
const settingsSaving = ref(false);

const openSettings = async (item: PluginItem) => {
  settingsPlugin.value = item;
  settingsModal.value = true;
  settingsLoading.value = true;
  try {
    const res = await pluginApi.getSettings(item.id);
    settingsSchema.value = res.schema ?? [];
    settingsValues.value = { ...(res.values ?? {}) };
    for (const field of settingsSchema.value) {
      if (!(field.key in settingsValues.value) && field.default !== undefined) {
        settingsValues.value[field.key] = field.default;
      }
    }
  } catch (e: any) {
    toast.add({ title: e?.message, color: "error" });
    settingsModal.value = false;
  } finally {
    settingsLoading.value = false;
  }
};

/** Group settings fields by their `group` property, preserving order. */
const settingsGroups = computed(() => {
  const groups: Array<{ name: string | undefined; fields: PluginSettingField[] }> = []
  const groupMap = new Map<string | undefined, PluginSettingField[]>()
  for (const field of settingsSchema.value) {
    const key = field.group || undefined
    if (!groupMap.has(key)) {
      const fields: PluginSettingField[] = []
      groupMap.set(key, fields)
      groups.push({ name: key, fields })
    }
    groupMap.get(key)!.push(field)
  }
  return groups
})

/** Evaluate a field's showIf condition against current settings values. */
function isFieldVisible(field: PluginSettingField): boolean {
  if (!field.showIf) return true
  try {
    // Support simple expressions like "key === true", "key === 'value'"
    const expr = field.showIf.trim()
    const match = expr.match(/^(\w+)\s*(===|!==|==|!=)\s*(.+)$/)
    if (!match) return true
    const [, key, op, rawVal] = match
    const actual = settingsValues.value[key]
    let expected: unknown = rawVal.trim()
    if (expected === 'true') expected = true
    else if (expected === 'false') expected = false
    else if (expected === 'null') expected = null
    else if (/^['"].*['"]$/.test(expected as string)) expected = (expected as string).slice(1, -1)
    else if (!isNaN(Number(expected))) expected = Number(expected)
    if (op === '===' || op === '==') return actual === expected
    if (op === '!==' || op === '!=') return actual !== expected
  } catch { /* show field on parse error */ }
  return true
}

const saveSettings = async () => {
  if (!settingsPlugin.value) return;
  settingsSaving.value = true;
  try {
    await pluginApi.updateSettings(settingsPlugin.value.id, settingsValues.value);
    toast.add({ title: t("admin.plugins.settings_save_success"), color: "success" });
    settingsModal.value = false;
  } catch (e: any) {
    toast.add({ title: e?.message, color: "error" });
  } finally {
    settingsSaving.value = false;
  }
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

/** Check if a plugin has an update available */
const hasUpdate = (item: PluginItem) => updateMap.value[item.id];

/** Number of plugins with available updates */
const updatableCount = computed(() =>
  Object.keys(updateMap.value).length
);

// ── Load ───────────────────────────────────────────────────────────────────
const load = async () => {
  rawLoading.value = true;
  try {
    const res = await pluginApi.list();
    items.value = res.items ?? [];
  } catch (e: any) {
    toast.add({ title: e?.message, color: "error" });
  } finally {
    rawLoading.value = false;
  }
};

/** Compare installed versions with marketplace to detect available updates */
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
    // Silently fail — update check is not critical
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

// ── Install modal ──────────────────────────────────────────────────────────
const installModal = ref(false);
const installTab = ref<"github" | "zip">("github");
const installUrl = ref("");
const zipFile = ref<File | null>(null);
const zipDragging = ref(false);
const fileInputRef = ref<HTMLInputElement | null>(null);
const installing = ref(false);

// Preview step state (GitHub tab only)
const previewStep = ref(false);
const previewData = ref<PluginPreviewInfo | null>(null);
const previewLoading = ref(false);

const openInstall = () => {
  installTab.value = "github";
  installUrl.value = "";
  zipFile.value = null;
  previewStep.value = false;
  previewData.value = null;
  installModal.value = true;
};

// Called when tab switches — reset preview state
watch(installTab, () => {
  previewStep.value = false;
  previewData.value = null;
});

const canInstall = computed(() => {
  if (installTab.value === "github") return installUrl.value.trim().length > 0;
  return zipFile.value !== null;
});

// Step 1 (GitHub): fetch preview then advance to step 2
const fetchPreview = async () => {
  if (!installUrl.value.trim()) return;
  previewLoading.value = true;
  try {
    previewData.value = await pluginApi.preview(installUrl.value.trim());
  } catch {
    // Preview failed — proceed anyway
    previewData.value = null;
  } finally {
    previewLoading.value = false;
    previewStep.value = true;
  }
};

const confirmInstall = async () => {
  if (!canInstall.value) return;
  installing.value = true;
  try {
    let item: PluginItem;
    if (installTab.value === "github") {
      const res = await pluginApi.install(installUrl.value.trim());
      item = res.item;
    } else {
      const res = await pluginApi.uploadZip(zipFile.value!);
      item = res.item;
    }
    items.value.unshift(item);
    if (item.need_restart) {
      toast.add({
        title: t("admin.plugins.install_success"),
        description: t("admin.plugins.restart_required_install"),
        color: "warning",
        duration: 0,
      });
    } else {
      toast.add({ title: t("admin.plugins.install_success"), color: "success" });
    }
    installModal.value = false;
  } catch (e: any) {
    toast.add({ title: e?.message, color: "error" });
  } finally {
    installing.value = false;
  }
};

const onFileInput = (e: Event) => {
  const f = (e.target as HTMLInputElement).files?.[0];
  if (f) zipFile.value = f;
};

const onDrop = (e: DragEvent) => {
  zipDragging.value = false;
  const f = e.dataTransfer?.files?.[0];
  if (f) zipFile.value = f;
};

// ── Toggle / Uninstall / Update ────────────────────────────────────────────
const togglingId = ref<string | null>(null);
const uninstallModal = ref(false);
const pendingUninstall = ref<PluginItem | null>(null);
const uninstalling = ref(false);
const updatingId = ref<string | null>(null);

// Cascade impact state
const cascadePlugins = ref<string[]>([]);
const cascadeLoading = ref(false);

// Disable confirmation modal
const disableModal = ref(false);
const pendingDisable = ref<PluginItem | null>(null);

const toggle = async (item: PluginItem) => {
  // If disabling, check for cascade impact first
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
      // If impact check fails, proceed without warning
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
    // If we disabled a plugin with cascade, remove dependents from active state in the UI
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
  cascadeLoading.value = true;
  uninstallModal.value = true;
  try {
    const res = await pluginApi.unloadImpact(item.id);
    cascadePlugins.value = res.will_unload ?? [];
  } catch {
    // Non-critical
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
        <UButton color="primary" icon="i-tabler-plus" @click="openInstall">
          {{ $t("admin.plugins.install") }}
        </UButton>
      </template>
    </AdminPageHeader>

    <AdminPageContent>
      <!-- Install modal -->
      <UModal v-model:open="installModal" :ui="{ content: 'max-w-md' }">
        <template #content>
          <div class="p-6">
            <h3 class="text-lg font-semibold text-highlighted mb-4">{{ $t("admin.plugins.install_title") }}</h3>

            <div class="flex gap-1 p-1 bg-muted rounded-md mb-4">
              <button
                v-for="tab in (['github', 'zip'] as const)" :key="tab"
                class="flex-1 flex items-center justify-center gap-1.5 py-1.5 px-3 rounded-md text-sm font-medium transition-colors"
                :class="installTab === tab ? 'bg-default text-highlighted shadow-xs' : 'text-muted hover:text-default'"
                @click="installTab = tab">
                <UIcon :name="tab === 'github' ? 'i-tabler-brand-github' : 'i-tabler-file-zip'" class="size-4" />
                {{ $t(`admin.plugins.install_tab_${tab}`) }}
              </button>
            </div>

            <div v-if="installTab === 'github'">
              <!-- Step 1: enter URL -->
              <template v-if="!previewStep">
                <UFormField :label="$t('admin.plugins.repo_url_label')">
                  <UInput v-model="installUrl" class="w-full" :placeholder="$t('admin.plugins.repo_url_placeholder')" @keydown.enter="fetchPreview" />
                </UFormField>
                <p class="text-xs text-muted mt-2">{{ $t("admin.plugins.repo_url_hint") }}</p>
              </template>
              <!-- Step 2: preview + confirm -->
              <template v-else>
                <p class="text-xs text-muted mb-3 font-mono truncate">{{ installUrl }}</p>
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
                </div>
                <PluginPreviewPanel v-else-if="previewData" :info="previewData" />
                <p v-else class="text-xs text-muted italic">{{ $t('admin.plugins.preview_no_caps') }}</p>
              </template>
            </div>
            <div v-else>
              <input ref="fileInputRef" type="file" accept=".zip" class="hidden" @change="onFileInput" />
              <div
                class="border-2 border-dashed rounded-md p-6 text-center cursor-pointer transition-colors"
                :class="zipDragging ? 'border-primary bg-primary/5' : 'border-default hover:border-primary/50'"
                @click="fileInputRef?.click()"
                @dragover.prevent="zipDragging = true"
                @dragleave.prevent="zipDragging = false"
                @drop.prevent="onDrop">
                <UIcon :name="zipFile ? 'i-tabler-file-check' : 'i-tabler-upload'" class="size-8 mx-auto mb-2" :class="zipFile ? 'text-success' : 'text-muted'" />
                <p v-if="zipFile" class="text-sm font-medium text-highlighted">{{ $t("admin.plugins.zip_selected", { name: zipFile.name }) }}</p>
                <p v-else class="text-sm text-muted">{{ $t("admin.plugins.zip_drop_label") }}</p>
              </div>
              <p class="text-xs text-muted mt-2">{{ $t("admin.plugins.zip_hint") }}</p>
            </div>

            <div class="flex justify-end gap-2 mt-6">
              <UButton color="neutral" variant="ghost" @click="installModal = false">{{ $t("common.cancel") }}</UButton>
              <!-- GitHub step 1: preview button -->
              <template v-if="installTab === 'github' && !previewStep">
                <UButton color="primary" :loading="previewLoading" :disabled="!canInstall" @click="fetchPreview">
                  {{ $t("admin.plugins.preview_confirm_title") }}
                </UButton>
              </template>
              <!-- GitHub step 2 or ZIP: back + install -->
              <template v-else>
                <UButton v-if="installTab === 'github'" color="neutral" variant="outline" @click="previewStep = false">{{ $t("common.back") }}</UButton>
                <UButton color="primary" leading-icon="i-tabler-download" :loading="installing" :disabled="!canInstall" @click="confirmInstall">{{ $t("admin.plugins.install") }}</UButton>
              </template>
            </div>
          </div>
        </template>
      </UModal>

      <!-- Settings modal -->
      <UModal v-model:open="settingsModal" :ui="{ content: 'max-w-md' }">
        <template #content>
          <div class="p-6">
            <h3 class="text-lg font-semibold text-highlighted mb-1">{{ $t("admin.plugins.settings_title") }}</h3>
            <p class="text-sm text-muted mb-5">{{ settingsPlugin?.title }}</p>

            <div v-if="settingsLoading" class="space-y-4 py-2">
              <div v-for="i in 3" :key="i" class="space-y-1.5">
                <USkeleton class="h-3 w-24" /><USkeleton class="h-9 w-full" />
              </div>
            </div>
            <p v-else-if="settingsSchema.length === 0" class="text-sm text-muted py-4 text-center">{{ $t("admin.plugins.settings_no_schema") }}</p>
            <div v-else class="space-y-5">
              <template v-for="(group, groupIdx) in settingsGroups" :key="group.name ?? '__default__'">
                <div v-if="group.name" :class="{ 'pt-3 border-t border-default': groupIdx > 0 }">
                  <p class="text-xs font-medium text-muted mb-3 uppercase tracking-wide">{{ group.name }}</p>
                </div>
                <template v-for="field in group.fields" :key="field.key">
                  <UFormField v-if="isFieldVisible(field)" :label="field.label || field.key" :description="field.description" :required="field.required">
                    <USwitch v-if="field.type === 'boolean'" :model-value="!!settingsValues[field.key]" @update:model-value="settingsValues[field.key] = $event" />
                    <USelect v-else-if="field.type === 'select'" :model-value="String(settingsValues[field.key] ?? '')" :items="(field.options ?? []).map(o => ({ label: o, value: o }))" class="w-full" @update:model-value="settingsValues[field.key] = $event" />
                    <UTextarea v-else-if="field.type === 'textarea'" :model-value="String(settingsValues[field.key] ?? '')" :placeholder="field.placeholder" class="w-full" :rows="3" @update:model-value="settingsValues[field.key] = $event" />
                    <UInput v-else-if="field.type === 'password'" type="password" :model-value="String(settingsValues[field.key] ?? '')" :placeholder="field.placeholder" class="w-full" @update:model-value="settingsValues[field.key] = $event" />
                    <UInput v-else-if="field.type === 'number'" type="number" :model-value="String(settingsValues[field.key] ?? '')" :placeholder="field.placeholder" class="w-full" @update:model-value="settingsValues[field.key] = Number($event)" />
                    <UInput v-else :model-value="String(settingsValues[field.key] ?? '')" :placeholder="field.placeholder" class="w-full" @update:model-value="settingsValues[field.key] = $event" />
                  </UFormField>
                </template>
              </template>
            </div>

            <!-- 3.5-F1: Capabilities display -->
            <div v-if="settingsPlugin && !settingsLoading" class="mt-5 pt-4 border-t border-default">
              <p class="text-xs font-medium text-muted mb-2 uppercase tracking-wide">{{ $t("admin.plugins.cap_section") }}</p>
              <div class="flex flex-wrap gap-1.5">
                <UBadge
                  v-for="badge in capBadges(settingsPlugin)"
                  :key="badge.label"
                  :label="badge.label"
                  :color="badge.color as any"
                  variant="soft"
                  size="sm" />
              </div>
            </div>

            <div class="flex justify-end gap-2 mt-6">
              <UButton color="neutral" variant="ghost" @click="settingsModal = false">{{ $t("common.cancel") }}</UButton>
              <UButton v-if="settingsSchema.length > 0" color="primary" :loading="settingsSaving" @click="saveSettings">{{ $t("common.save") }}</UButton>
            </div>
          </div>
        </template>
      </UModal>

      <!-- Uninstall confirm modal -->
      <UModal v-model:open="uninstallModal" :ui="{ content: 'max-w-sm' }">
        <template #content>
          <div class="p-6">
            <h3 class="text-lg font-semibold text-highlighted mb-2">{{ $t("admin.plugins.uninstall_confirm_title") }}</h3>
            <p class="text-sm text-muted mb-4">{{ $t("admin.plugins.uninstall_confirm_desc", { name: pendingUninstall?.title }) }}</p>
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
                      {{ items.find(p => p.id === pid)?.title || pid }}
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
            <p class="text-sm text-muted mb-4">{{ $t("admin.plugins.disable_confirm_desc", { name: pendingDisable?.title }) }}</p>
            <div class="mb-4 p-3 bg-warning-50 dark:bg-warning-950/30 border border-warning-200 dark:border-warning-800 rounded-md">
              <div class="flex items-start gap-2">
                <UIcon name="i-tabler-alert-triangle" class="size-4 text-warning-600 dark:text-warning-400 shrink-0 mt-0.5" />
                <div>
                  <p class="text-sm font-medium text-warning-800 dark:text-warning-200 mb-1">{{ $t("admin.plugins.cascade_disable_warning") }}</p>
                  <ul class="text-xs text-warning-700 dark:text-warning-300 space-y-0.5">
                    <li v-for="pid in cascadePlugins" :key="pid" class="flex items-center gap-1">
                      <UIcon name="i-tabler-plug" class="size-3" />
                      {{ items.find(p => p.id === pid)?.title || pid }}
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
        <UButton color="primary" icon="i-tabler-plus" @click="openInstall">{{ $t("admin.plugins.install") }}</UButton>
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
                  <h3 class="text-sm font-semibold text-highlighted">{{ item.title || item.id }}</h3>
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
                <p class="text-xs text-muted mb-1">{{ item.description }}</p>
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

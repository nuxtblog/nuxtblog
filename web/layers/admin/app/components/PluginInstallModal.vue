<script setup lang="ts">
import type { PluginItem, PluginPreviewInfo } from "~/composables/usePluginApi";

const props = defineProps<{
  open: boolean;
}>();

const emit = defineEmits<{
  (e: "update:open", val: boolean): void;
  (e: "installed", item: PluginItem): void;
}>();

const { t } = useI18n();
const pluginApi = usePluginApi();
const toast = useToast();

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

watch(() => props.open, (val) => {
  if (val) {
    installTab.value = "github";
    installUrl.value = "";
    zipFile.value = null;
    previewStep.value = false;
    previewData.value = null;
  }
});

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
    emit("installed", item);
    emit("update:open", false);
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
</script>

<template>
  <UModal :open="open" :ui="{ content: 'max-w-md' }" @update:open="emit('update:open', $event)">
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
          <UButton color="neutral" variant="ghost" @click="emit('update:open', false)">{{ $t("common.cancel") }}</UButton>
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
</template>

<template>
  <AdminPageContainer>
    <AdminPageHeader :title="$t('admin.settings.writing.title')" :subtitle="$t('admin.settings.writing.subtitle')">
      <template #actions>
        <UButton color="neutral" variant="outline" :disabled="isSaving" @click="loadSettings">{{ $t('common.reset') }}</UButton>
        <UButton color="primary" icon="i-tabler-device-floppy" :loading="isSaving" :disabled="isSaving" @click="saveSettings">
          {{ $t('common.save_changes') }}
        </UButton>
      </template>
    </AdminPageHeader>
    <AdminPageContent>
      <div v-if="isLoading" class="space-y-4">
        <UCard>
          <template #header><USkeleton class="h-5 w-40" /></template>
          <div class="space-y-4">
            <div v-for="i in 3" :key="i" class="space-y-2">
              <USkeleton class="h-4 w-24" />
              <USkeleton class="h-9 w-full rounded-md" />
            </div>
          </div>
        </UCard>
        <div class="flex justify-end gap-3">
          <USkeleton class="h-9 w-16 rounded-md" />
          <USkeleton class="h-9 w-24 rounded-md" />
        </div>
      </div>

      <template v-if="!isLoading">
        <UCard>
          <template #header>
            <h3 class="text-base font-semibold text-highlighted">{{ $t('admin.settings.writing.editor_settings') }}</h3>
          </template>
          <div class="space-y-4">
            <div>
              <label class="text-sm font-medium text-highlighted mb-2 block">{{ $t('admin.settings.writing.default_editor') }}</label>
              <div class="space-y-2">
                <label
                  v-for="editor in editors"
                  :key="editor.value"
                  class="flex items-start gap-3 p-3 bg-elevated/50 rounded-md cursor-pointer hover:bg-elevated">
                  <input v-model="form.defaultEditor" type="radio" :value="editor.value" class="w-4 h-4 text-primary mt-0.5" />
                  <div class="flex-1">
                    <div class="text-sm font-medium text-highlighted mb-1">{{ editor.label }}</div>
                    <div class="text-xs text-muted">{{ editor.description }}</div>
                  </div>
                </label>
              </div>
            </div>

            <div class="flex items-center justify-between pt-3 border-t border-default">
              <div>
                <h4 class="text-sm font-medium text-highlighted mb-1">{{ $t('admin.settings.writing.auto_save') }}</h4>
                <p class="text-xs text-muted">{{ $t('admin.settings.writing.auto_save_hint') }}</p>
              </div>
              <UCheckbox v-model="form.autoSave" />
            </div>

            <div v-if="form.autoSave" class="pt-3 border-t border-default">
              <UFormField :label="$t('admin.settings.writing.auto_save_interval')">
                <USelect
                  v-model="form.autoSaveInterval"
                  :items="[
                    { label: $t('admin.settings.writing.interval_30s'), value: 30 },
                    { label: $t('admin.settings.writing.interval_1m'), value: 60 },
                    { label: $t('admin.settings.writing.interval_2m'), value: 120 },
                    { label: $t('admin.settings.writing.interval_5m'), value: 300 },
                  ]"
                  class="w-full max-w-xs" />
              </UFormField>
            </div>
          </div>
        </UCard>

        <UCard class="mt-4">
          <template #header>
            <h3 class="text-base font-semibold text-highlighted">{{ $t('admin.settings.writing.code_highlight') }}</h3>
          </template>
          <div class="space-y-4">
            <UFormField :label="$t('admin.settings.writing.code_highlight_theme')" :hint="$t('admin.settings.writing.code_highlight_hint')">
              <USelect
                v-model="form.codeHighlightTheme"
                :items="highlightThemeOptions"
                class="w-full max-w-xs" />
            </UFormField>

            <div class="pt-3 border-t border-default">
              <UFormField :label="$t('admin.settings.writing.code_languages')" :hint="$t('admin.settings.writing.code_languages_hint')">
                <USelectMenu
                  v-model="form.codeLanguages"
                  :items="allHljsLanguages"
                  multiple
                  value-key="value"
                  label-key="label"
                  class="w-full max-w-lg"
                  :placeholder="$t('admin.settings.writing.code_languages_placeholder')" />
              </UFormField>
              <p v-if="!form.codeLanguages.length" class="text-xs text-muted mt-1">
                {{ $t('admin.settings.writing.code_languages_default') }}
              </p>
            </div>
          </div>
        </UCard>

      </template>
    </AdminPageContent>
  </AdminPageContainer>
</template>

<script setup lang="ts">
import { ref } from "vue";
import hljs from 'highlight.js';
import { HIGHLIGHT_THEMES } from '~/utils/highlightThemes';
import { formatLanguageLabel } from '~/utils/codeLanguages';

const { apiFetch } = useApiFetch();
const optionsStore = useOptionsStore();
const toast = useToast();
const { t } = useI18n();
const isSaving = ref(false);
const rawLoading = ref(true);
const isLoading = useMinLoading(rawLoading);

const form = ref({
  defaultEditor: "markdown",
  autoSave: true,
  autoSaveInterval: 60,
  codeHighlightTheme: "github",
  codeLanguages: [] as string[],
});

const highlightThemeOptions = HIGHLIGHT_THEMES.map(t => ({ label: t.label, value: t.id }));

// All languages supported by hljs, for the multi-select picker
const allHljsLanguages = hljs.listLanguages()
  .map(v => ({ label: formatLanguageLabel(v), value: v }))
  .sort((a, b) => a.label.localeCompare(b.label));

const editors = computed(() => [
  { value: "markdown",  label: t('admin.settings.writing.editor_markdown'), description: t('admin.settings.writing.editor_markdown_desc') },
  { value: "rich-text", label: t('admin.settings.writing.editor_richtext'), description: t('admin.settings.writing.editor_richtext_desc') },
]);

const loadSettings = async () => {
  try {
    const result = await apiFetch<{ options: Record<string, string> }>("/options/autoload");
    const opts = result.options ?? {};
    if (opts.default_editor !== undefined)     form.value.defaultEditor    = JSON.parse(opts.default_editor);
    if (opts.auto_save !== undefined)          form.value.autoSave         = JSON.parse(opts.auto_save) as boolean;
    if (opts.auto_save_interval !== undefined) form.value.autoSaveInterval = parseInt(JSON.parse(opts.auto_save_interval));
    if (opts.code_highlight_theme !== undefined) form.value.codeHighlightTheme = JSON.parse(opts.code_highlight_theme);
    if (opts.code_languages !== undefined) form.value.codeLanguages = JSON.parse(opts.code_languages) as string[];
  } catch (e) {
    console.error(e);
    toast.add({ title: t("admin.settings.writing.load_failed"), description: t("admin.settings.general.load_failed_desc"), color: "error", icon: "i-tabler-alert-circle" });
  } finally {
    rawLoading.value = false;
  }
};

const saveSettings = async () => {
  isSaving.value = true;
  try {
    const keyMap: Array<[string, unknown]> = [
      ["default_editor",     form.value.defaultEditor],
      ["auto_save",          form.value.autoSave],
      ["auto_save_interval", form.value.autoSaveInterval],
      ["code_highlight_theme", form.value.codeHighlightTheme],
      ["code_languages",       form.value.codeLanguages],
    ];
    await Promise.all(
      keyMap.map(([key, value]) =>
        apiFetch(`/options/${key}`, { method: "PUT", body: { value: JSON.stringify(value), autoload: 1 } })
      )
    );
    await optionsStore.reload();
    toast.add({ title: t("admin.settings.writing.saved"), description: t("admin.settings.writing.saved_desc"), color: "success", icon: "i-tabler-circle-check" });
  } catch (e) {
    console.error(e);
    toast.add({ title: t("common.save_failed"), description: t("admin.settings.general.save_failed_desc"), color: "error", icon: "i-tabler-alert-circle" });
  } finally {
    isSaving.value = false;
  }
};

await loadSettings();
</script>

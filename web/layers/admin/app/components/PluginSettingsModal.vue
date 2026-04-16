<script setup lang="ts">
import type { PluginItem, PluginSettingField } from "~/composables/usePluginApi";
import { parseCapabilityBadges } from "~/composables/usePluginApi";
import { resolveTemplate } from '~/composables/usePluginI18n'

const props = defineProps<{
  open: boolean;
  plugin: PluginItem | null;
}>();

const emit = defineEmits<{
  (e: "update:open", val: boolean): void;
}>();

const { t, locale } = useI18n();
const pluginApi = usePluginApi();
const toast = useToast();

const capBadges = (item: PluginItem) => parseCapabilityBadges(item.capabilities ?? '{}', t);

const settingsSchema = ref<PluginSettingField[]>([]);
const settingsValues = ref<Record<string, unknown>>({});
const settingsLoading = ref(false);
const settingsSaving = ref(false);

watch(() => props.open, async (val) => {
  if (!val || !props.plugin) return;
  settingsLoading.value = true;
  try {
    const res = await pluginApi.getSettings(props.plugin.id);
    settingsSchema.value = res.schema ?? [];
    settingsValues.value = { ...(res.values ?? {}) };
    for (const field of settingsSchema.value) {
      if (!(field.key in settingsValues.value) && field.default !== undefined) {
        settingsValues.value[field.key] = field.default;
      }
    }
  } catch (e: any) {
    toast.add({ title: e?.message, color: "error" });
    emit("update:open", false);
  } finally {
    settingsLoading.value = false;
  }
});

const saveSettings = async () => {
  if (!props.plugin) return;
  settingsSaving.value = true;
  try {
    await pluginApi.updateSettings(props.plugin.id, settingsValues.value);
    toast.add({ title: t("admin.plugins.settings_save_success"), color: "success" });
    emit("update:open", false);
  } catch (e: any) {
    toast.add({ title: e?.message, color: "error" });
  } finally {
    settingsSaving.value = false;
  }
};
</script>

<template>
  <UModal :open="open" :ui="{ content: 'max-w-md' }" @update:open="emit('update:open', $event)">
    <template #content>
      <div class="p-6">
        <h3 class="text-lg font-semibold text-highlighted mb-1">{{ $t("admin.plugins.settings_title") }}</h3>
        <p class="text-sm text-muted mb-5">{{ plugin ? resolveTemplate(plugin.title, plugin.id, locale) : '' }}</p>

        <div v-if="settingsLoading" class="space-y-4 py-2">
          <div v-for="i in 3" :key="i" class="space-y-1.5">
            <USkeleton class="h-3 w-24" /><USkeleton class="h-9 w-full" />
          </div>
        </div>
        <p v-else-if="settingsSchema.length === 0" class="text-sm text-muted py-4 text-center">{{ $t("admin.plugins.settings_no_schema") }}</p>
        <div v-else class="space-y-5">
          <PluginSettingFields
            :schema="settingsSchema"
            :model-value="settingsValues"
            :plugin-id="plugin?.id"
            @update:model-value="settingsValues = $event"
          />
        </div>

        <!-- Capabilities display -->
        <div v-if="plugin && !settingsLoading" class="mt-5 pt-4 border-t border-default">
          <p class="text-xs font-medium text-muted mb-2 uppercase tracking-wide">{{ $t("admin.plugins.cap_section") }}</p>
          <div class="flex flex-wrap gap-1.5">
            <UBadge
              v-for="badge in capBadges(plugin)"
              :key="badge.label"
              :label="badge.label"
              :color="badge.color as any"
              variant="soft"
              size="sm" />
          </div>
        </div>

        <div class="flex justify-end gap-2 mt-6">
          <UButton color="neutral" variant="ghost" @click="emit('update:open', false)">{{ $t("common.cancel") }}</UButton>
          <UButton v-if="settingsSchema.length > 0" color="primary" :loading="settingsSaving" @click="saveSettings">{{ $t("common.save") }}</UButton>
        </div>
      </div>
    </template>
  </UModal>
</template>

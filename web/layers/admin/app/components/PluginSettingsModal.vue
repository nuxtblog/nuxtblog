<script setup lang="ts">
import type { PluginItem, PluginSettingField } from "~/composables/usePluginApi";
import { parseCapabilityBadges } from "~/composables/usePluginApi";

const props = defineProps<{
  open: boolean;
  plugin: PluginItem | null;
}>();

const emit = defineEmits<{
  (e: "update:open", val: boolean): void;
}>();

const { t } = useI18n();
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

/** Group settings fields by their `group` property, preserving order. */
const settingsGroups = computed(() => {
  const groups: Array<{ name: string | undefined; fields: PluginSettingField[] }> = [];
  const groupMap = new Map<string | undefined, PluginSettingField[]>();
  for (const field of settingsSchema.value) {
    const key = field.group || undefined;
    if (!groupMap.has(key)) {
      const fields: PluginSettingField[] = [];
      groupMap.set(key, fields);
      groups.push({ name: key, fields });
    }
    groupMap.get(key)!.push(field);
  }
  return groups;
});

/** Evaluate a field's showIf condition against current settings values. */
function isFieldVisible(field: PluginSettingField): boolean {
  if (!field.showIf) return true;
  try {
    const expr = field.showIf.trim();
    const match = expr.match(/^(\w+)\s*(===|!==|==|!=)\s*(.+)$/);
    if (!match) return true;
    const [, key, op, rawVal] = match;
    const actual = settingsValues.value[key];
    let expected: unknown = rawVal.trim();
    if (expected === 'true') expected = true;
    else if (expected === 'false') expected = false;
    else if (expected === 'null') expected = null;
    else if (/^['"].*['"]$/.test(expected as string)) expected = (expected as string).slice(1, -1);
    else if (!isNaN(Number(expected))) expected = Number(expected);
    if (op === '===' || op === '==') return actual === expected;
    if (op === '!==' || op === '!=') return actual !== expected;
  } catch { /* show field on parse error */ }
  return true;
}

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
        <p class="text-sm text-muted mb-5">{{ plugin?.title }}</p>

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

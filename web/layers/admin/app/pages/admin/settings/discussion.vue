<template>
  <AdminPageContainer>
    <AdminPageHeader :title="$t('admin.settings.discussion.title')" :subtitle="$t('admin.settings.discussion.subtitle')">
      <template #actions>
        <UButton color="neutral" variant="outline" :disabled="isSaving" @click="loadSettings">{{ $t('common.reset') }}</UButton>
        <UButton color="primary" icon="i-tabler-device-floppy" :loading="isSaving" :disabled="isSaving" @click="saveSettings">
          {{ $t('common.save_changes') }}
        </UButton>
      </template>
    </AdminPageHeader>
    <AdminPageContent>
      <div v-if="isLoading" class="space-y-4">
        <UCard v-for="i in 2" :key="i">
          <template #header><USkeleton class="h-5 w-40" /></template>
          <div class="space-y-4">
            <div v-for="j in 3" :key="j" class="space-y-2">
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
            <h3 class="text-base font-semibold text-highlighted">{{ $t('admin.settings.discussion.default_post_settings') }}</h3>
          </template>
          <div class="flex items-center justify-between">
            <div>
              <h4 class="text-sm font-medium text-highlighted mb-1">{{ $t('admin.settings.discussion.allow_comments') }}</h4>
              <p class="text-xs text-muted">{{ $t('admin.settings.discussion.allow_comments_hint') }}</p>
            </div>
            <UCheckbox v-model="form.allowComments" />
          </div>
        </UCard>

        <UCard>
          <template #header>
            <h3 class="text-base font-semibold text-highlighted">{{ $t('admin.settings.discussion.moderation') }}</h3>
          </template>
          <div class="space-y-4">
            <div class="flex items-center justify-between">
              <div>
                <h4 class="text-sm font-medium text-highlighted mb-1">{{ $t('admin.settings.discussion.require_moderation') }}</h4>
                <p class="text-xs text-muted">{{ $t('admin.settings.discussion.require_moderation_desc') }}</p>
              </div>
              <UCheckbox v-model="form.requireModeration" />
            </div>

            <div class="flex items-center justify-between pt-3 border-t border-default">
              <div>
                <h4 class="text-sm font-medium text-highlighted mb-1">{{ $t('admin.settings.discussion.require_name_email') }}</h4>
                <p class="text-xs text-muted">{{ $t('admin.settings.discussion.require_name_email_desc') }}</p>
              </div>
              <UCheckbox v-model="form.requireNameEmail" />
            </div>

            <div class="pt-3 border-t border-default">
              <UFormField :label="$t('admin.settings.discussion.moderate_links')">
                <UInput v-model="form.moderateLinksCount" type="number" :min="0" :max="10" class="w-full max-w-xs" />
                <p class="text-xs text-muted mt-1">{{ $t('admin.settings.discussion.moderate_links_hint') }}</p>
              </UFormField>
            </div>

            <div class="pt-3 border-t border-default">
              <UFormField :label="$t('admin.settings.discussion.blocklist_label')">
                <UTextarea
                  v-model="form.commentBlocklist"
                  :placeholder="$t('admin.settings.discussion.blocklist_hint')"
                  :rows="4"
                  class="w-full resize-none" />
                <p class="text-xs text-muted mt-1">{{ $t('admin.settings.discussion.blocklist_hint') }}</p>
              </UFormField>
            </div>
          </div>
        </UCard>

      </template>
    </AdminPageContent>
  </AdminPageContainer>
</template>

<script setup lang="ts">
import { ref } from "vue";

const { apiFetch } = useApiFetch();
const toast = useToast();
const { t } = useI18n();
const isSaving = ref(false);
const rawLoading = ref(true);
const isLoading = useMinLoading(rawLoading);

const form = ref({
  allowComments: true,
  requireModeration: false,
  requireNameEmail: true,
  moderateLinksCount: 2,
  commentBlocklist: "",
});

const loadSettings = async () => {
  try {
    const result = await apiFetch<{ options: Record<string, string> }>("/options/autoload");
    const opts = result.options ?? {};
    if (opts.default_allow_comments !== undefined)     form.value.allowComments       = JSON.parse(opts.default_allow_comments) as boolean;
    if (opts.comment_moderation !== undefined)         form.value.requireModeration   = JSON.parse(opts.comment_moderation) as boolean;
    if (opts.comment_require_name_email !== undefined) form.value.requireNameEmail    = JSON.parse(opts.comment_require_name_email) as boolean;
    if (opts.comment_max_links !== undefined)          form.value.moderateLinksCount  = parseInt(JSON.parse(opts.comment_max_links));
    if (opts.comment_blacklist !== undefined)          form.value.commentBlocklist    = JSON.parse(opts.comment_blacklist);
  } catch (e) {
    console.error(e);
    toast.add({ title: t('common.load_failed'), description: t('common.cannot_read_settings'), color: "error", icon: "i-tabler-alert-circle" });
  } finally {
    rawLoading.value = false;
  }
};

const saveSettings = async () => {
  isSaving.value = true;
  try {
    const keyMap: Array<[string, unknown]> = [
      ["default_allow_comments",     form.value.allowComments],
      ["comment_moderation",         form.value.requireModeration],
      ["comment_require_name_email", form.value.requireNameEmail],
      ["comment_max_links",          form.value.moderateLinksCount],
      ["comment_blacklist",          form.value.commentBlocklist],
    ];
    await Promise.all(
      keyMap.map(([key, value]) =>
        apiFetch(`/options/${key}`, { method: "PUT", body: { value: JSON.stringify(value), autoload: 1 } })
      )
    );
    toast.add({ title: t('admin.settings.discussion.saved'), description: t('admin.settings.discussion.saved_desc'), color: "success", icon: "i-tabler-circle-check" });
  } catch (e) {
    console.error(e);
    toast.add({ title: t('common.save_failed'), description: t('common.settings_save_failed'), color: "error", icon: "i-tabler-alert-circle" });
  } finally {
    isSaving.value = false;
  }
};

await loadSettings();
</script>

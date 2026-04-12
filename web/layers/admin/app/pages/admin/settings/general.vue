<template>
  <AdminPageContainer>
    <AdminPageHeader :title="$t('admin.settings.general.title')" :subtitle="$t('admin.settings.general.subtitle')">
      <template #actions>
        <UButton color="neutral" variant="outline" :disabled="isSaving" @click="loadSettings">{{ $t('common.reset') }}</UButton>
        <UButton color="primary" icon="i-tabler-device-floppy" :loading="isSaving" :disabled="isSaving" @click="saveSettings">
          {{ $t('common.save_changes') }}
        </UButton>
      </template>
    </AdminPageHeader>
    <AdminPageContent>
      <!-- 骨架屏加载状态 -->
      <div v-if="isLoading" class="space-y-4">
        <UCard v-for="i in 4" :key="i">
          <template #header><USkeleton class="h-5 w-40" /></template>
          <div class="space-y-4">
            <div v-for="j in 3" :key="j" class="space-y-2">
              <USkeleton class="h-4 w-24" />
              <USkeleton class="h-9 w-full rounded-md" />
            </div>
          </div>
        </UCard>
      </div>

      <template v-if="!isLoading">
        <!-- 站点标识 -->
        <UCard>
          <template #header>
            <h3 class="text-base font-semibold text-highlighted">{{ $t('admin.settings.general.site_identity') }}</h3>
          </template>
          <div class="space-y-4">
            <UFormField :label="$t('admin.settings.general.site_title')" required>
              <UInput
                v-model="form.siteTitle"
                placeholder="Yozya Blog"
                class="w-full" />
              <p class="text-xs text-muted mt-1">{{ $t('admin.settings.general.site_title_hint') }}</p>
            </UFormField>
            <UFormField :label="$t('admin.settings.general.tagline')">
              <UInput
                v-model="form.tagline"
                :placeholder="$t('admin.settings.general.tagline_placeholder')"
                class="w-full" />
              <p class="text-xs text-muted mt-1">{{ $t('admin.settings.general.tagline_hint') }}</p>
            </UFormField>
            <UFormField :label="$t('admin.settings.general.default_author_bio')">
              <UInput
                v-model="form.defaultAuthorBio"
                :placeholder="$t('admin.settings.general.default_author_bio_placeholder')"
                class="w-full" />
              <p class="text-xs text-muted mt-1">{{ $t('admin.settings.general.default_author_bio_hint') }}</p>
            </UFormField>
          </div>
        </UCard>

        <!-- 默认媒体 -->
        <UCard>
          <template #header>
            <h3 class="text-base font-semibold text-highlighted">{{ $t('admin.settings.general.default_media') }}</h3>
          </template>
          <div class="space-y-4">
            <UFormField :label="$t('admin.settings.general.default_post_cover')">
              <UInput
                v-model="form.defaultPostCover"
                placeholder="/images/default-cover.svg"
                class="w-full" />
              <p class="text-xs text-muted mt-1">{{ $t('admin.settings.general.default_post_cover_hint') }}</p>
            </UFormField>
            <UFormField :label="$t('admin.settings.general.error_post_cover')">
              <UInput
                v-model="form.errorPostCover"
                placeholder="/images/default-cover.svg"
                class="w-full" />
              <p class="text-xs text-muted mt-1">{{ $t('admin.settings.general.error_post_cover_hint') }}</p>
            </UFormField>
            <UFormField :label="$t('admin.settings.general.default_user_bg')">
              <UInput
                v-model="form.defaultUserBg"
                placeholder="/images/default-user-bg.jpg"
                class="w-full" />
              <p class="text-xs text-muted mt-1">{{ $t('admin.settings.general.default_user_bg_hint') }}</p>
            </UFormField>
            <UFormField :label="$t('admin.settings.general.default_avatar')">
              <div class="space-y-2 w-full">
                <USelect
                  v-model="form.defaultAvatarType"
                  :items="[
                    { label: $t('admin.settings.general.avatar_initials'), value: 'initials' },
                    { label: $t('admin.settings.general.avatar_image'), value: 'image' },
                  ]"
                  class="w-full" />
                <UInput
                  v-if="form.defaultAvatarType === 'image'"
                  v-model="form.defaultAvatarUrl"
                  placeholder="/images/default-avatar.svg"
                  class="w-full" />
              </div>
              <p class="text-xs text-muted mt-1">{{ $t('admin.settings.general.default_avatar_hint') }}</p>
            </UFormField>
          </div>
        </UCard>

        <!-- 站点地址 -->
        <UCard>
          <template #header>
            <h3 class="text-base font-semibold text-highlighted">{{ $t('admin.settings.general.site_address') }}</h3>
          </template>
          <UFormField :label="$t('admin.settings.general.site_url')" required>
            <UInput
              v-model="form.siteUrl"
              type="url"
              placeholder="http://localhost:3000"
              class="w-full" />
            <p class="text-xs text-muted mt-1">{{ $t('admin.settings.general.site_url_hint') }}</p>
          </UFormField>
        </UCard>

        <!-- 管理设置 -->
        <UCard>
          <template #header>
            <h3 class="text-base font-semibold text-highlighted">{{ $t('admin.settings.general.admin_settings') }}</h3>
          </template>
          <div class="space-y-4">
            <UFormField :label="$t('admin.settings.general.admin_email')" required>
              <UInput
                v-model="form.adminEmail"
                type="email"
                placeholder="admin@example.com"
                class="w-full" />
              <p class="text-xs text-muted mt-1">{{ $t('admin.settings.general.admin_email_hint') }}</p>
            </UFormField>
            <UFormField :label="$t('admin.settings.general.site_language')">
              <USelect
                v-model="form.language"
                :items="languageItems"
                class="w-full" />
              <p class="text-xs text-muted mt-1">{{ $t('admin.settings.general.site_language_hint') }}</p>
            </UFormField>
          </div>
        </UCard>

        <!-- 页脚 -->
        <UCard>
          <template #header>
            <h3 class="text-base font-semibold text-highlighted">{{ $t('admin.settings.general.footer') }}</h3>
          </template>
          <div class="space-y-6">
            <UFormField :label="$t('admin.settings.general.footer_text')">
              <UInput
                v-model="form.footerText"
                placeholder="All rights reserved."
                class="w-full" />
              <p class="text-xs text-muted mt-1">{{ $t('admin.settings.general.footer_text_hint') }}</p>
            </UFormField>

            <USeparator />

            <!-- 页脚导航 -->
            <div class="space-y-2">
              <div class="flex items-center justify-between">
                <p class="text-sm font-medium text-highlighted">{{ $t('admin.settings.general.footer_nav') }}</p>
                <UButton
                  size="xs"
                  variant="soft"
                  color="primary"
                  leading-icon="i-tabler-plus"
                  @click="addFooterLink">
                  {{ $t('admin.settings.general.add_link') }}
                </UButton>
              </div>
              <div v-if="form.footerLinks.length" class="space-y-2">
                <div
                  v-for="(link, i) in form.footerLinks"
                  :key="i"
                  class="flex items-center gap-2">
                  <UInput
                    v-model="link.label"
                    :placeholder="$t('admin.settings.general.link_label_placeholder')"
                    size="sm"
                    class="w-28 shrink-0" />
                  <UInput
                    v-model="link.url"
                    placeholder="/about"
                    size="sm"
                    class="flex-1" />
                  <UButton
                    color="error"
                    variant="ghost"
                    icon="i-tabler-trash"
                    size="xs"
                    @click="form.footerLinks.splice(i, 1)" />
                </div>
              </div>
              <p v-else class="text-xs text-muted">{{ $t('admin.settings.general.no_links') }}</p>
            </div>

            <USeparator />

            <!-- 社交链接 -->
            <div class="space-y-2">
              <div class="flex items-center justify-between">
                <p class="text-sm font-medium text-highlighted">{{ $t('admin.settings.general.social_links') }}</p>
                <UButton
                  size="xs"
                  variant="soft"
                  color="primary"
                  leading-icon="i-tabler-plus"
                  @click="addSocialLink">
                  {{ $t('common.add') }}
                </UButton>
              </div>
              <div v-if="form.socialLinks.length" class="space-y-2">
                <div
                  v-for="(link, i) in form.socialLinks"
                  :key="i"
                  class="flex items-center gap-2">
                  <USelect
                    v-model="link.label"
                    :items="SOCIAL_PLATFORMS"
                    value-key="value"
                    label-key="label"
                    size="sm"
                    class="w-36 shrink-0" />
                  <UInput
                    v-if="link.label === '__custom__'"
                    v-model="link.customLabel"
                    :placeholder="$t('admin.settings.general.custom_name')"
                    size="sm"
                    class="w-28 shrink-0" />
                  <UInput
                    v-model="link.url"
                    placeholder="https://github.com/yourname"
                    size="sm"
                    class="flex-1" />
                  <UButton
                    color="error"
                    variant="ghost"
                    icon="i-tabler-trash"
                    size="xs"
                    @click="form.socialLinks.splice(i, 1)" />
                </div>
              </div>
              <p v-else class="text-xs text-muted">{{ $t('admin.settings.general.no_social_links') }}</p>
              <p class="text-xs text-muted">{{ $t('admin.settings.general.social_hint') }}</p>
            </div>
          </div>
        </UCard>

        <!-- 备案信息 -->
        <UCard>
          <template #header>
            <h3 class="text-base font-semibold text-highlighted">{{ $t('admin.settings.general.icp_info') }}</h3>
          </template>
          <div class="space-y-6">
            <div class="space-y-3">
              <p class="text-sm font-medium text-highlighted">{{ $t('admin.settings.general.icp_filing') }}</p>
              <UFormField :label="$t('admin.settings.general.filing_number')">
                <UInput
                  v-model="form.icpNumber"
                  :placeholder="$t('admin.settings.general.icp_number_placeholder')"
                  class="w-full" />
              </UFormField>
              <UFormField :label="$t('admin.settings.general.filing_url')">
                <UInput
                  v-model="form.icpUrl"
                  placeholder="https://beian.miit.gov.cn"
                  class="w-full" />
                <p class="text-xs text-muted mt-1">{{ $t('admin.settings.general.icp_url_hint') }}</p>
              </UFormField>
            </div>
            <USeparator />
            <div class="space-y-3">
              <p class="text-sm font-medium text-highlighted">{{ $t('admin.settings.general.police_filing') }}</p>
              <UFormField :label="$t('admin.settings.general.filing_number')">
                <UInput
                  v-model="form.policeNumber"
                  :placeholder="$t('admin.settings.general.police_number_placeholder')"
                  class="w-full" />
              </UFormField>
              <UFormField :label="$t('admin.settings.general.filing_url')">
                <UInput
                  v-model="form.policeUrl"
                  placeholder="https://www.beian.gov.cn/..."
                  class="w-full" />
                <p class="text-xs text-muted mt-1">{{ $t('admin.settings.general.police_url_hint') }}</p>
              </UFormField>
            </div>
          </div>
        </UCard>

      </template>
    </AdminPageContent>
  </AdminPageContainer>
</template>

<script setup lang="ts">
import type { LocaleCode } from '~/types/locale'

interface SimpleLink {
  label: string;
  url: string;
  customLabel?: string;
}

// SOCIAL_PLATFORMS is auto-imported from ~/utils/social

const { apiFetch } = useApiFetch();
const siteApi = useSiteApi();
const toast = useToast();
const optionsStore = useOptionsStore();
const { t } = useI18n();
const { $i18n } = useNuxtApp();

const languageItems = ref<Array<{ label: string; value: string }>>([]);
const { data: langsData } = await useAsyncData('site-languages', siteApi.getLanguages, { lazy: true });
watch(langsData, (data) => {
  if (data?.list) {
    languageItems.value = data.list.map((l) => ({ label: l.name, value: l.code }));
  }
}, { immediate: true });
const isSaving = ref(false);
const rawLoading = ref(true);
const isLoading = useMinLoading(rawLoading);

const form = ref({
  siteTitle: "",
  tagline: "",
  defaultAuthorBio: "",
  defaultPostCover: "",
  errorPostCover: "",
  defaultUserBg: "",
  siteUrl: "",
  adminEmail: "",
  language: "zh",
  footerText: "",
  icpNumber: "",
  icpUrl: "",
  policeNumber: "",
  policeUrl: "",
  footerLinks: [] as SimpleLink[],
  socialLinks: [] as SimpleLink[],
  defaultAvatarType: "initials",
  defaultAvatarUrl: "",
});

const addFooterLink = () => {
  form.value.footerLinks.push({ label: "", url: "" });
};
const addSocialLink = () => {
  form.value.socialLinks.push({ label: "GitHub", url: "" });
};
const loadSettings = async () => {
  try {
    const result = await apiFetch<{ options: Record<string, string> }>(
      "/options/autoload",
    );
    const opts = result.options ?? {};
    const p = (key: string) => {
      try {
        return JSON.parse(opts[key] ?? "null");
      } catch {
        return null;
      }
    };
    form.value.siteTitle = p("site_name") ?? "";
    form.value.tagline = p("site_description") ?? "";
    form.value.defaultAuthorBio = p("default_author_bio") ?? "";
    form.value.siteUrl = p("site_url") ?? "";
    form.value.adminEmail = p("admin_email") ?? "";
    form.value.language = p("site_language") ?? "zh";
    form.value.defaultPostCover = p("default_post_cover") ?? "";
    form.value.errorPostCover = p("error_post_cover") ?? "";
    form.value.footerText = p("footer_text") ?? "";
    form.value.icpNumber = p("icp_number") ?? "";
    form.value.icpUrl = p("icp_url") ?? "";
    form.value.policeNumber = p("police_number") ?? "";
    form.value.policeUrl = p("police_url") ?? "";
    form.value.footerLinks = p("footer_links") ?? [];
    form.value.socialLinks = p("social_links") ?? [];
    form.value.defaultUserBg = p("default_user_bg") ?? "";
    form.value.defaultAvatarType = p("default_avatar_type") ?? "initials";
    form.value.defaultAvatarUrl = p("default_avatar_url") ?? "";
  } catch (e) {
    toast.add({
      title: t("admin.settings.general.load_failed"),
      description: t("admin.settings.general.load_failed_desc"),
      color: "error",
    });
  } finally {
    rawLoading.value = false;
  }
};

const saveSettings = async () => {
  isSaving.value = true;
  try {
    // 保存 socialLinks 时把 __custom__ 替换为实际名称
    const socialLinks = form.value.socialLinks.map((l) => ({
      label: l.label === "__custom__" ? l.customLabel || t("admin.settings.general.default_link") : l.label,
      url: l.url,
    }));

    const keyMap: Array<[string, unknown]> = [
      ["site_name", form.value.siteTitle],
      ["site_description", form.value.tagline],
      ["default_author_bio", form.value.defaultAuthorBio],
      ["site_url", form.value.siteUrl],
      ["admin_email", form.value.adminEmail],
      ["site_language", form.value.language],
      ["default_post_cover", form.value.defaultPostCover],
      ["error_post_cover", form.value.errorPostCover],
      ["footer_text", form.value.footerText],
      ["icp_number", form.value.icpNumber],
      ["icp_url", form.value.icpUrl],
      ["police_number", form.value.policeNumber],
      ["police_url", form.value.policeUrl],
      ["footer_links", form.value.footerLinks],
      ["social_links", socialLinks],
      ["default_user_bg", form.value.defaultUserBg],
      ["default_avatar_type", form.value.defaultAvatarType],
      ["default_avatar_url", form.value.defaultAvatarUrl],
    ];
    await Promise.all(
      keyMap.map(([key, value]) =>
        apiFetch(`/options/${key}`, {
          method: "PUT",
          body: { value: JSON.stringify(value), autoload: 1 },
        }),
      ),
    );
    await $i18n.setLocale(form.value.language as LocaleCode)
    await optionsStore.reload()
    toast.add({
      title: t("admin.settings.general.saved"),
      description: t("admin.settings.general.saved_desc"),
      color: "success",
    });
  } catch (e) {
    toast.add({
      title: t("common.save_failed"),
      description: t("admin.settings.general.save_failed_desc"),
      color: "error",
    });
  } finally {
    isSaving.value = false;
  }
};

await loadSettings();
</script>

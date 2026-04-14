<template>
  <footer class="mt-auto border-t border-default bg-muted">
    <div class="container mx-auto px-4 py-10">
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-8 sm:gap-4">
        <!-- 左列：站点信息 -->
        <div class="space-y-3">
          <p class="font-semibold text-highlighted text-base">{{ siteName }}</p>
          <p v-if="tagline" class="text-sm text-muted leading-relaxed">
            {{ tagline }}
          </p>
          <p class="text-xs text-muted">
            &copy; {{ year }} {{ siteName
            }}{{ footerText ? ". " + footerText : "" }}
          </p>
        </div>

        <!-- 中列：页脚导航 -->
        <div v-if="footerLinks.length" class="space-y-2">
          <p
            class="text-xs font-semibold text-muted uppercase tracking-wider mb-3">
            {{ $t("site.footer.nav") }}
          </p>
          <a
            v-for="link in footerLinks"
            :key="link.label + link.url"
            :href="link.url"
            class="block text-sm text-muted hover:text-highlighted transition-colors"
            >{{ link.label }}</a
          >
        </div>
        <div v-else />

        <!-- 右列：社交 + 备案 -->
        <div class="space-y-4">
          <div v-if="socialLinks.length" class="space-y-2">
            <p
              class="text-xs font-semibold text-muted uppercase tracking-wider">
              {{ $t("site.footer.social") }}
            </p>
            <div class="flex flex-wrap gap-3">
              <UTooltip
                v-for="link in socialLinks"
                :key="link.label + link.url"
                :text="link.label">
                <a
                  :href="link.url"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="text-muted hover:text-highlighted transition-colors">
                  <UIcon
                    :name="getSocialIcon(link.label, link.url)"
                    class="size-5" />
                </a>
              </UTooltip>
            </div>
          </div>

          <!-- 备案 -->
          <div class="text-xs text-muted flex gap-2">
            <a
              v-if="icpNumber && icpUrl"
              rel="nofollow"
              target="_blank"
              :href="icpUrl"
              class="block hover:text-highlighted transition-colors"
              >{{ icpNumber }}</a
            >
            <span v-else-if="icpNumber">{{ icpNumber }}</span>

            <span v-if="icpNumber || policeNumber">・</span>
            <a
              v-if="policeNumber && policeUrl"
              rel="nofollow"
              target="_blank"
              :href="policeUrl"
              class="flex items-center gap-1 hover:text-highlighted transition-colors">
              <img src="/images/beian-ico.png" alt="" class="w-3.5 h-3.5" />
              {{ policeNumber }}
            </a>
            <span v-else-if="policeNumber" class="flex items-center gap-1">
              <img src="/images/beian-ico.png" alt="" class="w-3.5 h-3.5" />
              {{ policeNumber }}
            </span>
          </div>
        </div>
      </div>

      <ClientOnly>
        <div class="col-span-full">
          <ContributionSlot name="public:footer-extra" />
        </div>
      </ClientOnly>
    </div>
  </footer>
</template>

<script setup lang="ts">
interface SimpleLink {
  label: string;
  url: string;
}

const optionStore = useOptionsStore();

const year = new Date().getFullYear();
const siteName = computed(() => optionStore.get("site_name", "个人博客"));
const tagline = computed(() => optionStore.get("site_description", ""));
const footerText = computed(() => optionStore.get("footer_text", ""));
const icpNumber = computed(() => optionStore.get("icp_number", ""));
const icpUrl = computed(() => optionStore.get("icp_url", ""));
const policeNumber = computed(() => optionStore.get("police_number", ""));
const policeUrl = computed(() => optionStore.get("police_url", ""));

const footerLinks = computed(() =>
  optionStore.getJSON<SimpleLink[]>("footer_links", []),
);
const socialLinks = computed(() =>
  optionStore.getJSON<SimpleLink[]>("social_links", []),
);

// getSocialIcon is auto-imported from ~/utils/social
</script>

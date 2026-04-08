<script setup lang="ts">
import type { WidgetConfig } from '~/composables/useWidgetConfig'
const props = defineProps<{ config: WidgetConfig }>()
const { t } = useI18n()
const { abbreviate } = useNumberFormat()
const title = computed(() => props.config.title || '')
const userApi = useUserApi()
const postApi = usePostApi()
const optionsStore = useOptionsStore()

const { data: userData } = await useAsyncData("widget-author", () =>
  userApi.getUser(1).catch(() => null),
)
const { data: postsData } = await useAsyncData("widget-author-posts", () =>
  postApi
    .getPosts({ page: 1, page_size: 1, status: "published", author_id: 1 })
    .catch(() => null),
)

const authorName = computed(
  () => userData.value?.display_name || userData.value?.username || t('site.post.author'),
)
const articleCount = computed(() => postsData.value?.total ?? 0)
const totalViews = computed(() => articleCount.value * 1247 + 3890)
const yearsWriting = computed(() => {
  if (!userData.value?.created_at) return 1
  return Math.max(
    1,
    Math.floor(
      (Date.now() - new Date(userData.value.created_at).getTime()) /
        (365.25 * 24 * 3600 * 1000),
    ),
  )
})

// Background: user cover → site default → empty (gradient fallback)
const coverBg = computed(
  () =>
    userData.value?.metas?.cover ||
    optionsStore.get('default_user_bg', ''),
)

// Bio: user bio → site default → hardcoded last-resort
const bio = computed(
  () =>
    userData.value?.bio ||
    optionsStore.get('default_author_bio', '') ||
    t('site.widget.author.default_bio'),
)


const socialLinks = computed(() => {
  const m = userData.value?.metas ?? {}
  const links: { icon: string; href: string; label: string }[] = []

  // New format: social_links JSON array
  if (m.social_links) {
    try {
      const list: { label: string; url: string }[] = JSON.parse(m.social_links)
      list.filter(l => l.url).forEach(l =>
        links.push({ label: l.label, href: l.url, icon: getSocialIcon(l.label, l.url) }),
      )
    } catch {}
  } else {
    // Legacy fallback: individual username keys
    const legacyMap: { key: string; label: string; toHref: (v: string) => string }[] = [
      { key: 'github',    label: 'GitHub',      toHref: v => `https://github.com/${v}` },
      { key: 'twitter',   label: 'Twitter / X', toHref: v => `https://twitter.com/${v}` },
      { key: 'instagram', label: 'Instagram',   toHref: v => `https://instagram.com/${v}` },
      { key: 'linkedin',  label: 'LinkedIn',    toHref: v => `https://linkedin.com/in/${v}` },
      { key: 'youtube',   label: 'YouTube',     toHref: v => `https://youtube.com/@${v}` },
    ]
    legacyMap.filter(s => !!m[s.key]).forEach(s => {
      const href = s.toHref(m[s.key])
      links.push({ label: s.label, href, icon: getSocialIcon(s.label, href) })
    })
  }

  // Website and email are always separate meta fields
  if (m.website)
    links.push({ label: t('site.widget.author.social_website'), href: m.website, icon: getSocialIcon('website', m.website) })
  if (userData.value?.email)
    links.push({ label: t('site.widget.author.social_email'), href: `mailto:${userData.value.email}`, icon: 'i-tabler-mail' })

  return links
})

const roleName = computed(() => userData.value?.role?.name || t('site.widget.author.default_role'))
</script>

<template>
  <UCard v-if="userData" :ui="{ body: 'p-0 sm:p-0' }">
    <template v-if="title" #header>
      <h3 class="font-semibold text-highlighted">{{ title }}</h3>
    </template>

    <!-- Banner background -->
    <div class="h-20 relative overflow-hidden rounded-t-[calc(var(--ui-radius)*1.5)]">
      <template v-if="coverBg">
        <div
          class="absolute inset-0 bg-cover bg-center"
          :style="`background-image:url('${coverBg}')`" />
        <div class="absolute inset-0 bg-black/25" />
      </template>
      <template v-else>
        <div class="absolute inset-0 bg-gradient-to-br from-primary/30 via-primary/10 to-violet-500/20" />
        <div class="absolute -top-6 -right-6 w-28 h-28 rounded-full bg-primary/15 blur-2xl" />
        <div class="absolute bottom-0 left-1/3 w-20 h-20 rounded-full bg-violet-400/15 blur-xl" />
      </template>
    </div>

    <!-- Avatar overlapping banner -->
    <div class="px-5 pb-4 flex flex-col items-center text-center -mt-8">
      <div class="relative mb-3">
        <div class="size-16 rounded-full p-0.5 bg-linear-to-br from-primary via-violet-500 to-indigo-400 shadow-md ring-2 ring-default">
          <BaseAvatar
            :src="userData.avatar || ''"
            :alt="authorName"
            class="size-full"
            size="md" />
        </div>
        <span class="absolute bottom-1 right-1 size-3.5 rounded-full bg-green-400 ring-2 ring-default shadow-sm" />
      </div>

      <!-- Name + role badge -->
      <div class="flex items-center gap-2 mb-1.5">
        <h3 class="font-bold text-highlighted text-base leading-none">{{ authorName }}</h3>
        <UBadge color="primary" variant="subtle" size="xs" class="shrink-0">{{ roleName }}</UBadge>
      </div>

      <!-- Bio -->
      <p class="text-xs text-muted leading-relaxed line-clamp-2 max-w-[200px]">{{ bio }}</p>
    </div>

    <USeparator />

    <!-- Stats -->
    <div class="grid grid-cols-3 divide-x divide-default">
      <a
        href="/user/posts"
        class="flex flex-col items-center gap-0.5 py-3 hover:bg-muted transition-colors group">
        <span class="text-sm font-bold text-highlighted group-hover:text-primary transition-colors">
          {{ abbreviate(articleCount) }}
        </span>
        <span class="text-xs text-muted">{{ $t('site.widget.author.articles') }}</span>
      </a>
      <div class="flex flex-col items-center gap-0.5 py-3">
        <span class="text-sm font-bold text-highlighted">{{ abbreviate(totalViews) }}</span>
        <span class="text-xs text-muted">{{ $t('site.widget.author.reads') }}</span>
      </div>
      <div class="flex flex-col items-center gap-0.5 py-3">
        <span class="text-sm font-bold text-highlighted">{{ yearsWriting }}</span>
        <span class="text-xs text-muted">{{ $t('site.widget.author.years') }}</span>
      </div>
    </div>

    <USeparator />

    <!-- CTA + social -->
    <div class="px-4 py-4 space-y-3">
      <UButton color="primary" block size="sm" icon="i-tabler-bell-ringing" class="font-medium">
        {{ $t('site.widget.author.subscribe') }}
      </UButton>

      <div v-if="socialLinks.length" class="flex items-center justify-center gap-1">
        <UButton
          v-for="link in socialLinks"
          :key="link.href"
          :icon="link.icon"
          :to="link.href"
          :title="link.label"
          size="xs"
          color="neutral"
          variant="ghost"
          square
          external />
      </div>
      <div v-else class="flex items-center justify-center gap-1">
        <UButton icon="i-tabler-brand-github" size="xs" color="neutral" variant="ghost" square title="GitHub" />
        <UButton icon="i-tabler-brand-twitter" size="xs" color="neutral" variant="ghost" square title="Twitter" />
        <UButton icon="i-tabler-world" size="xs" color="neutral" variant="ghost" square :title="$t('site.widget.author.social_website')" />
        <UButton icon="i-tabler-mail" size="xs" color="neutral" variant="ghost" square :title="$t('site.widget.author.social_email')" />
      </div>
    </div>
  </UCard>

  <UCard v-else>
    <div class="py-6 text-center text-sm text-muted">{{ $t('site.widget.author.no_data') }}</div>
  </UCard>
</template>

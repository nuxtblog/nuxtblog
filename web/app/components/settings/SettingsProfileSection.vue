<template>
  <UCard>
    <template #header>
      <h2 class="font-semibold text-highlighted flex items-center gap-2">
        <UIcon name="i-tabler-user-circle" class="size-5 text-primary" />
        {{ $t('site.settings.profile_title') }}
      </h2>
    </template>

    <form class="space-y-5" @submit.prevent="saveProfile">
      <!-- 头像（只读，去个人主页改） -->
      <div class="flex items-center gap-4 pb-4 border-b border-default">
        <BaseAvatar
          :src="authStore.user?.avatar"
          :alt="authStore.user?.display_name || authStore.user?.username"
          size="xl" />
        <div>
          <p class="text-sm font-medium text-highlighted">{{ $t('site.settings.avatar_label') }}</p>
          <NuxtLink :to="`/user/${authStore.user?.id}`" class="text-xs text-primary hover:underline">
            {{ $t('site.settings.avatar_hint') }}
          </NuxtLink>
        </div>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <UFormField :label="$t('site.settings.display_name_label')" class="md:col-span-2">
          <UInput v-model="profile.display_name" :placeholder="$t('site.settings.display_name_placeholder')" class="w-full" />
        </UFormField>
        <UFormField :label="$t('site.settings.bio_label')" class="md:col-span-2">
          <UTextarea v-model="profile.bio" :placeholder="$t('site.settings.bio_placeholder')" :rows="3" class="w-full" />
        </UFormField>
        <UFormField :label="$t('site.settings.location_label')">
          <UInput v-model="profile.location" :placeholder="$t('site.settings.location_placeholder')" class="w-full" leading-icon="i-tabler-map-pin" />
        </UFormField>
        <UFormField :label="$t('site.settings.website_label')">
          <UInput v-model="profile.website" placeholder="https://yoursite.com" class="w-full" leading-icon="i-tabler-link" />
        </UFormField>
      </div>

      <!-- 社交链接 -->
      <div class="border-t border-default pt-5">
        <div class="flex items-center justify-between mb-3">
          <p class="text-sm font-medium text-highlighted flex items-center gap-2">
            <UIcon name="i-tabler-share" class="size-4 text-muted" />
            {{ $t('site.settings.social_section') }}
          </p>
          <UButton size="xs" variant="soft" color="primary" leading-icon="i-tabler-plus" @click="addSocialLink">
            {{ $t('site.settings.social_add_link') }}
          </UButton>
        </div>
        <div v-if="socialLinks.length" class="space-y-2">
          <div v-for="(link, i) in socialLinks" :key="i" class="flex items-center gap-2">
            <UIcon :name="getSocialIcon(link.label, link.url)" class="size-5 text-muted shrink-0" />
            <USelect v-model="link.label" :items="SOCIAL_PLATFORMS" value-key="value" label-key="label" size="sm" class="w-32 shrink-0" />
            <UInput
              v-if="link.label === '__custom__'"
              v-model="link.customLabel"
              :placeholder="$t('site.settings.social_custom_name')"
              size="sm"
              class="w-28 shrink-0" />
            <UInput v-model="link.url" :placeholder="$t('site.settings.social_url_placeholder')" size="sm" class="flex-1" />
            <UButton color="error" variant="ghost" icon="i-tabler-trash" size="xs" square @click="socialLinks.splice(i, 1)" />
          </div>
        </div>
        <p v-else class="text-xs text-muted">{{ $t('site.settings.social_no_links') }}</p>
      </div>

      <div class="flex items-center justify-between pt-2 border-t border-default">
        <p class="text-xs text-muted">
          {{ $t('site.settings.username_immutable', { username: authStore.user?.username }) }}
        </p>
        <UButton type="submit" color="primary" :loading="saving" icon="i-tabler-check">
          {{ $t('site.settings.save_profile') }}
        </UButton>
      </div>
    </form>
  </UCard>
</template>

<script setup lang="ts">
interface SocialLink { label: string; url: string; customLabel?: string }

const authStore = useAuthStore()
const userApi = useUserApi()
const toast = useToast()
const { t } = useI18n()

const profile = reactive({
  display_name: authStore.user?.display_name || '',
  bio: authStore.user?.bio || '',
  location: '',
  website: '',
})
const socialLinks = ref<SocialLink[]>([])
const saving = ref(false)
const loaded = ref(false)

function parseSocialLinks(metas: Record<string, string> | undefined): SocialLink[] {
  if (!metas) return []
  if (metas.social_links) {
    try { return JSON.parse(metas.social_links) } catch {}
  }
  const legacy: SocialLink[] = []
  const legacyMap: [string, string][] = [
    ['github', 'GitHub'], ['twitter', 'Twitter'], ['instagram', 'Instagram'],
    ['linkedin', 'LinkedIn'], ['youtube', 'YouTube'],
  ]
  for (const [key, label] of legacyMap) {
    if (metas[key]) legacy.push({ label, url: metas[key] })
  }
  return legacy
}

async function loadProfile() {
  if (loaded.value || !authStore.user?.id) return
  try {
    const u = await userApi.getUser(authStore.user.id)
    profile.display_name = u.display_name || ''
    profile.bio = u.bio || ''
    profile.location = u.metas?.location || ''
    profile.website = u.metas?.website || ''
    socialLinks.value = parseSocialLinks(u.metas)
    loaded.value = true
  } catch {}
}

onMounted(loadProfile)

function addSocialLink() {
  socialLinks.value.push({ label: 'GitHub', url: '' })
}

async function saveProfile() {
  if (!authStore.user?.id) return
  saving.value = true
  try {
    const links = socialLinks.value.map(l => ({
      label: l.label === '__custom__' ? (l.customLabel || 'Link') : l.label,
      url: l.url,
    }))
    await userApi.updateUser(authStore.user.id, {
      display_name: profile.display_name,
      bio: profile.bio,
      location: profile.location,
      website: profile.website,
      social_links: JSON.stringify(links),
    })
    authStore.setUser({ ...authStore.user!, display_name: profile.display_name, bio: profile.bio })
    toast.add({ title: t('site.settings.profile_saved'), color: 'success' })
  } catch (e: any) {
    toast.add({ title: t('site.settings.save_failed'), description: e.message, color: 'error' })
  } finally {
    saving.value = false
  }
}
</script>

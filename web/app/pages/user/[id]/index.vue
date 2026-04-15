<script setup lang="ts">
import type { UserResponse } from "~/types/api/user";
import type { PostListItemResponse } from "~/types/api/post";

const { containerClass } = useContainerWidth();
const { t } = useI18n();
const route = useRoute();
const router = useRouter();
const authStore = useAuthStore();
const optionsStore = useOptionsStore();
const userApi = useUserApi();
const postApi = usePostApi();
const reactionApi = useReactionApi();
const followApi = useFollowApi();
const toast = useToast();

const userId = computed(() => Number(route.params.id));
const isOwnProfile = computed(
  () => authStore.isLoggedIn && authStore.user?.id === userId.value,
);

// ── Tab 状态 ─────────────────────────────────────────────────────────────
const validTabs = ["home", "posts", "favorites", "following", "followers"];
const activeTab = ref(
  validTabs.includes(String(route.query.tab ?? ""))
    ? String(route.query.tab)
    : "home",
);
const tabs = computed(() => [
  { label: t("site.user.tab_home"), value: "home", icon: "i-tabler-home" },
  {
    label: t("site.user.tab_posts"),
    value: "posts",
    icon: "i-tabler-file-text",
  },
  ...(isOwnProfile.value
    ? [
        {
          label: t("site.user.tab_favorites"),
          value: "favorites",
          icon: "i-tabler-bookmark",
        },
      ]
    : []),
  {
    label: t("site.user.following_label"),
    value: "following",
    icon: "i-tabler-user-plus",
  },
  {
    label: t("site.user.followers_label"),
    value: "followers",
    icon: "i-tabler-users",
  },
]);

watch(activeTab, (tab) => {
  router.replace({ query: { ...route.query, tab } });
});

watch(
  () => route.query.tab,
  (tab) => {
    const next = validTabs.includes(String(tab ?? "")) ? String(tab) : "home";
    if (activeTab.value !== next) activeTab.value = next;
  },
);

// ── User & state & Posts ──────────────────────────────────────────────────────────
const user = ref<UserResponse | null>(null);
const rawLoading = ref(true);
const loading = useMinLoading(rawLoading);
const error = ref(false);

const posts = ref<PostListItemResponse[]>([]);
const postsTotal = ref(0);

// ── Bookmarks total (stats row only) ─────────────────────────────────────
const bookmarksTotal = ref(0);

// ── Follow state ──────────────────────────────────────────────────────────
const following = ref(false);
const followLoading = ref(false);
const followerCount = ref(0);
const followingCount = ref(0);

// ── Avatar / Cover previews (used in computed below) ─────────────────────
const avatarPreview = ref("");
const coverPreview = ref("");

// ── Computed ──────────────────────────────────────────────────────────────
const displayName = computed(
  () => user.value?.display_name || user.value?.username || "",
);
const avatar = computed(() => avatarPreview.value || user.value?.avatar || "");
const coverBg = computed(
  () =>
    coverPreview.value ||
    user.value?.metas?.cover ||
    optionsStore.get("default_user_bg", ""),
);
const defaultAuthorBio = computed(
  () => optionsStore.get("default_author_bio", "") as string,
);

const joinedDate = computed(() => {
  if (!user.value?.created_at) return "";
  return new Date(user.value.created_at).toLocaleDateString("zh-CN", {
    year: "numeric",
    month: "long",
  });
});
const socialLinks = computed(() => {
  const m = user.value?.metas ?? {};
  // New format: social_links is a JSON array of { label, url }
  if (m.social_links) {
    try {
      const list: { label: string; url: string }[] = JSON.parse(m.social_links);
      return list
        .filter((l) => l.url)
        .map((l, i) => ({
          key: String(i),
          label: l.label,
          href: l.url,
          icon: getSocialIcon(l.label, l.url),
        }));
    } catch {}
  }
  // Legacy fallback: individual username keys
  const legacyMap: {
    key: string;
    label: string;
    toHref: (v: string) => string;
  }[] = [
    {
      key: "github",
      label: "GitHub",
      toHref: (v) => `https://github.com/${v}`,
    },
    {
      key: "twitter",
      label: "Twitter / X",
      toHref: (v) => `https://twitter.com/${v}`,
    },
    {
      key: "instagram",
      label: "Instagram",
      toHref: (v) => `https://instagram.com/${v}`,
    },
    {
      key: "linkedin",
      label: "LinkedIn",
      toHref: (v) => `https://linkedin.com/in/${v}`,
    },
    {
      key: "youtube",
      label: "YouTube",
      toHref: (v) => `https://youtube.com/@${v}`,
    },
  ];
  return legacyMap
    .filter((s) => !!m[s.key])
    .map((s) => ({
      key: s.key,
      label: s.label,
      href: s.toHref(m[s.key]),
      icon: getSocialIcon(s.label, s.toHref(m[s.key])),
    }));
});

// ── Data loading ──────────────────────────────────────────────────────────
onMounted(async () => {
  if (!userId.value) {
    error.value = true;
    return;
  }
  rawLoading.value = true;
  error.value = false;
  try {
    const [userData, postsData] = await Promise.all([
      userApi.getUser(userId.value),
      postApi.getPosts({
        author_id: userId.value,
        status: "2",
        page: 1,
        page_size: 5,
      }),
    ]);
    user.value = userData;
    posts.value = postsData.data;
    postsTotal.value = postsData.total;
    if (isOwnProfile.value) {
      reactionApi
        .getBookmarks(1, 1)
        .then((r) => {
          bookmarksTotal.value = r.total;
        })
        .catch(() => {});
    }
    followApi
      .getStatus(userId.value)
      .then((s) => {
        followerCount.value = s.follower_count;
        followingCount.value = s.following_count;
        following.value = s.following;
      })
      .catch(() => {});
  } catch {
    error.value = true;
  } finally {
    rawLoading.value = false;
  }
});

async function toggleFollow() {
  if (!authStore.isLoggedIn) {
    router.push("/auth/login");
    return;
  }
  if (followLoading.value) return;
  followLoading.value = true;
  const wasFollowing = following.value;
  following.value = !wasFollowing;
  followerCount.value += wasFollowing ? -1 : 1;
  try {
    if (wasFollowing) {
      await followApi.unfollow(userId.value);
    } else {
      await followApi.follow(userId.value);
    }
  } catch {
    following.value = wasFollowing;
    followerCount.value += wasFollowing ? 1 : -1;
    toast.add({ title: t("site.user.op_failed"), color: "error" });
  } finally {
    followLoading.value = false;
  }
}

// ── Avatar updated callback ───────────────────────────────────────────────
function onAvatarUpdated(newAvatarId: number) {
  if (user.value) {
    user.value.avatar_id = newAvatarId;
  }
  if (authStore.user) {
    authStore.setUser({ ...authStore.user, avatar: avatarPreview.value });
  }
}

useHead(() => ({
  title: user.value
    ? t("site.user.home_title", { name: displayName.value })
    : t("site.user.default_title"),
}));
</script>

<template>
  <div class="min-h-screen pb-16">
    <!-- Cover -->
    <UserCoverCrop
      :model-value="coverBg"
      :user-id="userId"
      :editable="isOwnProfile"
      @update:model-value="coverPreview = $event" />

    <div :class="[containerClass, 'mx-auto px-4 md:px-6']">
      <!-- Loading skeleton -->
      <template v-if="loading">
        <div class="flex items-end justify-between -mt-14 mb-4">
          <USkeleton
            class="size-28 rounded-full ring-2 ring-white dark:ring-zinc-900 shrink-0" />
          <USkeleton class="h-8 w-24 rounded-md mb-1" />
        </div>
        <div
          class="rounded-md bg-default shadow-sm ring-1 ring-default p-5 mb-4 space-y-3">
          <USkeleton class="h-7 w-44" />
          <USkeleton class="h-4 w-24" />
          <USkeleton class="h-4 w-full" />
          <USkeleton class="h-4 w-3/4" />
          <div class="flex gap-3">
            <USkeleton class="h-4 w-24" />
            <USkeleton class="h-4 w-20" />
          </div>
        </div>
        <div
          class="grid grid-cols-4 rounded-md overflow-hidden ring-1 ring-default bg-default shadow-sm mb-6">
          <div
            v-for="i in 4"
            :key="i"
            class="flex flex-col items-center py-4 gap-2"
            :class="i !== 1 ? 'border-l border-default' : ''">
            <USkeleton class="h-8 w-12" />
            <USkeleton class="h-3 w-10" />
          </div>
        </div>
        <div class="flex gap-1 border-b border-default mb-6">
          <USkeleton v-for="i in 3" :key="i" class="h-9 w-16 rounded mb-1" />
        </div>
        <div class="space-y-3">
          <UCard v-for="i in 3" :key="i">
            <div class="flex items-center gap-3">
              <USkeleton class="size-12 rounded-md shrink-0" />
              <div class="flex-1 space-y-2">
                <USkeleton class="h-4 w-3/4" />
                <USkeleton class="h-3 w-1/2" />
              </div>
            </div>
          </UCard>
        </div>
      </template>

      <!-- Error -->
      <div v-else-if="error" class="text-center py-20">
        <UIcon
          name="i-tabler-user-off"
          class="size-14 text-muted mx-auto mb-4" />
        <p class="font-semibold text-highlighted mb-1">
          {{ $t("site.user.not_found") }}
        </p>
        <p class="text-sm text-muted mb-4">
          {{ $t("site.user.not_found_desc") }}
        </p>
        <UButton to="/" color="neutral" variant="outline" size="sm">{{
          $t("site.user.back_home")
        }}</UButton>
      </div>

      <!-- Content -->
      <template v-else-if="user">
        <!-- Plugin slot: user profile top -->
        <ClientOnly><ContributionSlot name="public:user-profile-top" /></ClientOnly>

        <!-- Avatar + actions row -->
        <div class="relative z-10 flex items-end justify-between -mt-14 mb-4">
          <UserAvatarCrop
            v-if="isOwnProfile"
            :model-value="avatar"
            :user-id="userId"
            :editable="true"
            :avatar-id="user.avatar_id"
            :alt="displayName"
            @update:model-value="avatarPreview = $event"
            @updated="onAvatarUpdated" />
          <BaseAvatar
            v-else
            :src="avatar"
            :alt="displayName"
            size="3xl"
            class="ring-2 ring-white dark:ring-zinc-900 shadow-xl shrink-0" />

          <UButton
            v-if="isOwnProfile"
            to="/user/settings"
            color="neutral"
            variant="outline"
            size="sm"
            icon="i-tabler-edit"
            class="mb-1">
            {{ $t("site.user.edit_profile") }}
          </UButton>
          <div v-else class="flex items-center gap-2 mb-1">
            <UButton
              v-if="authStore.isLoggedIn"
              :to="`/messages/${userId}`"
              color="neutral"
              variant="outline"
              size="sm"
              icon="i-tabler-message">
              {{ $t("site.widget.author.send_dm") }}
            </UButton>
            <UButton
              :color="following ? 'neutral' : 'primary'"
              :variant="following ? 'outline' : 'solid'"
              size="sm"
              :icon="following ? 'i-tabler-user-check' : 'i-tabler-user-plus'"
              :loading="followLoading"
              @click="toggleFollow">
              {{
                following ? $t("site.user.following") : $t("site.user.follow")
              }}
            </UButton>
            <ReportButton
              v-if="authStore.isLoggedIn && user"
              target-type="user"
              :target-id="user.id" />
          </div>
        </div>

        <!-- Profile info card -->
        <div
          class="rounded-md bg-default shadow-sm ring-1 ring-default p-5 mb-4">
          <h1 class="text-2xl font-bold text-highlighted leading-tight">
            {{ displayName }}
          </h1>
          <p class="text-sm text-muted mt-0.5 mb-3">@{{ user.username }}</p>

          <p v-if="user.bio" class="text-default leading-relaxed text-sm mb-4">
            {{ user.bio }}
          </p>
          <p v-else-if="isOwnProfile" class="text-muted italic text-sm mb-4">
            {{ $t("site.user.own_no_bio") }}
            <NuxtLink
              to="/user/settings"
              class="text-primary hover:underline"
              >{{ $t("site.user.own_no_bio_link") }}</NuxtLink
            >
            {{ $t("site.user.own_no_bio_suffix") }}
          </p>
          <p
            v-else-if="defaultAuthorBio"
            class="text-default leading-relaxed text-sm mb-4">
            {{ defaultAuthorBio }}
          </p>
          <p v-else class="text-muted italic text-sm mb-4">
            {{ $t("site.user.no_bio") }}
          </p>

          <div class="flex flex-wrap gap-x-5 gap-y-1.5 text-sm text-muted mb-4">
            <span v-if="joinedDate" class="flex items-center gap-1.5">
              <UIcon name="i-tabler-calendar" class="size-4 shrink-0" />
              {{ $t("site.user.join_date", { date: joinedDate }) }}
            </span>
            <span v-if="user.metas?.location" class="flex items-center gap-1.5">
              <UIcon name="i-tabler-map-pin" class="size-4 shrink-0" />
              {{ user.metas.location }}
            </span>
            <a
              v-if="user.metas?.website"
              :href="user.metas.website"
              target="_blank"
              rel="noopener noreferrer"
              class="flex items-center gap-1.5 text-primary hover:underline">
              <UIcon name="i-tabler-link" class="size-4 shrink-0" />
              {{ user.metas.website.replace(/^https?:\/\//, "") }}
            </a>
          </div>

          <div v-if="socialLinks.length > 0" class="flex gap-0.5 flex-wrap">
            <UTooltip
              v-for="link in socialLinks"
              :key="link.key"
              :text="link.label">
              <UButton
                :href="link.href"
                target="_blank"
                rel="noopener noreferrer"
                color="neutral"
                variant="ghost"
                size="sm"
                square
                :icon="link.icon" />
            </UTooltip>
          </div>
        </div>

        <!-- Stats row -->
        <div
          class="grid grid-cols-4 rounded-md overflow-hidden ring-1 ring-default bg-default shadow-sm mb-6">
          <div class="flex flex-col items-center py-4">
            <span class="text-2xl font-bold text-primary">{{
              postsTotal
            }}</span>
            <span class="text-xs text-muted mt-0.5">{{
              $t("site.user.posts_label")
            }}</span>
          </div>
          <button
            class="flex flex-col items-center py-4 border-l border-default hover:bg-elevated transition-colors"
            @click="activeTab = 'followers'">
            <span class="text-2xl font-bold text-primary">{{
              followerCount
            }}</span>
            <span class="text-xs text-muted mt-0.5">{{
              $t("site.user.followers_label")
            }}</span>
          </button>
          <button
            class="flex flex-col items-center py-4 border-l border-default hover:bg-elevated transition-colors"
            @click="activeTab = 'following'">
            <span class="text-2xl font-bold text-primary">{{
              followingCount
            }}</span>
            <span class="text-xs text-muted mt-0.5">{{
              $t("site.user.following_label")
            }}</span>
          </button>
          <div class="flex flex-col items-center py-4 border-l border-default">
            <template v-if="isOwnProfile">
              <span class="text-2xl font-bold text-primary">{{
                bookmarksTotal
              }}</span>
              <span class="text-xs text-muted mt-0.5">{{
                $t("site.user.favorites_label")
              }}</span>
            </template>
            <template v-else>
              <UBadge
                :label="user.role?.name || $t('site.user.default_role')"
                color="primary"
                variant="soft" />
              <span class="text-xs text-muted mt-1">{{
                $t("site.user.identity_label")
              }}</span>
            </template>
          </div>
        </div>
        <!-- Tabs -->
        <div class="flex gap-1 mb-2 p-1 bg-default rounded-md w-fit flex-wrap">
          <UButton
            v-for="tab in tabs"
            :key="tab.value"
            :icon="tab.icon"
            :variant="activeTab === tab.value ? 'solid' : 'ghost'"
            :color="activeTab === tab.value ? 'primary' : 'neutral'"
            size="sm"
            :class="activeTab === tab.value ? 'shadow ' : 'text-muted'"
            @click="activeTab = tab.value">
            {{ tab.label }}
          </UButton>
        </div>

        <UserTabHome
          v-if="activeTab === 'home'"
          :posts="posts"
          :loading="loading"
          @view-all="activeTab = 'posts'" />
        <UserTabPosts v-else-if="activeTab === 'posts'" :user-id="userId" />
        <UserTabFavorites
          v-else-if="activeTab === 'favorites' && isOwnProfile" />
        <UserTabFollowList
          v-else-if="activeTab === 'following'"
          :user-id="userId"
          type="following" />
        <UserTabFollowList
          v-else-if="activeTab === 'followers'"
          :user-id="userId"
          type="followers" />

        <!-- Plugin slot: user profile bottom -->
        <ClientOnly><ContributionSlot name="public:user-profile-bottom" /></ClientOnly>
      </template>
    </div>
  </div>
</template>

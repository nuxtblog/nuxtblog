<template>
  <AdminPageContainer>
    <AdminPageHeader
      :title="$t('admin.dashboard.title')"
      :subtitle="$t('admin.dashboard.welcome', { name: displayName })" />

    <AdminPageContent>
      <template v-if="pending">
        <div class="grid grid-cols-2 lg:grid-cols-4 gap-3 md:gap-4">
          <UCard v-for="i in 4" :key="i">
            <div class="flex items-center gap-4">
              <USkeleton class="size-10 rounded-md shrink-0" />
              <div class="flex-1 space-y-2">
                <USkeleton class="h-6 w-16" />
                <USkeleton class="h-3 w-24" />
              </div>
            </div>
          </UCard>
        </div>
        <div class="grid grid-cols-1 lg:grid-cols-3 gap-4 md:gap-6 mt-4">
          <div class="lg:col-span-2">
            <UCard>
              <template #header><USkeleton class="h-5 w-24" /></template>
              <div class="space-y-3">
                <div v-for="i in 5" :key="i" class="flex items-start gap-4">
                  <USkeleton class="h-16 w-24 rounded shrink-0" />
                  <div class="flex-1 space-y-2">
                    <USkeleton class="h-4 w-3/4" />
                    <USkeleton class="h-3 w-1/2" />
                  </div>
                </div>
              </div>
            </UCard>
          </div>
          <UCard>
            <template #header><USkeleton class="h-5 w-20" /></template>
            <div class="space-y-4">
              <div v-for="i in 4" :key="i" class="flex gap-3">
                <USkeleton class="size-8 rounded-full shrink-0" />
                <div class="flex-1 space-y-1.5">
                  <USkeleton class="h-3 w-32" />
                  <USkeleton class="h-3 w-full" />
                </div>
              </div>
            </div>
          </UCard>
        </div>
      </template>

      <template v-else>
        <!-- 统计卡片 -->
        <div class="grid grid-cols-2 lg:grid-cols-4 gap-3 md:gap-4">
          <StatCard
            icon="i-tabler-file-text"
            iconColor="text-primary"
            iconBg="bg-primary/10"
            :total="postStats.total"
            :label="$t('admin.dashboard.total_posts')">
            <template #description>
              {{ $t('admin.dashboard.published_draft', { published: postStats.published, draft: postStats.draft }) }}
            </template>
          </StatCard>

          <StatCard
            icon="i-tabler-message-circle"
            iconColor="text-info"
            iconBg="bg-info/10"
            :total="commentStats.total"
            :label="$t('admin.dashboard.total_comments')">
            <template #description>
              <span v-if="commentStats.pending > 0" class="text-warning">
                {{ $t('admin.dashboard.pending_review', { n: commentStats.pending }) }}
              </span>
              <span v-else class="text-success">{{ $t('admin.dashboard.no_pending') }}</span>
            </template>
          </StatCard>

          <StatCard
            icon="i-tabler-users"
            iconColor="text-success"
            iconBg="bg-success/10"
            :total="userStats.total"
            :label="$t('admin.dashboard.registered_users')">
            <template #description>{{ $t('admin.dashboard.active_users', { n: userStats.active }) }}</template>
          </StatCard>

          <StatCard
            icon="i-tabler-eye"
            iconColor="text-warning"
            iconBg="bg-warning/10"
            :total="abbreviateNumber(totalViews)"
            :label="$t('admin.dashboard.total_views')">
            <template #description>{{ $t('admin.dashboard.all_posts_views') }}</template>
          </StatCard>
        </div>

        <div class="grid grid-cols-1 lg:grid-cols-3 gap-4 md:gap-6">
          <!-- 左侧 -->
          <div class="lg:col-span-2 space-y-4 md:space-y-6">
            <!-- 快速操作 -->
            <UCard>
              <template #header>
                <h2 class="text-base font-semibold text-highlighted">
                  {{ $t('admin.dashboard.quick_actions') }}
                </h2>
              </template>
              <div class="grid grid-cols-2 md:grid-cols-4 gap-3">
                <QuickActionButton
                  :label="$t('admin.dashboard.write_post')"
                  class="cursor-pointer"
                  @action="router.push('/admin/posts/new')">
                  <template #icon><UIcon name="i-tabler-pencil" /></template>
                </QuickActionButton>
                <QuickActionButton
                  :label="$t('admin.dashboard.upload_media')"
                  class="cursor-pointer"
                  @action="router.push('/admin/media')">
                  <template #icon><UIcon name="i-tabler-photo" /></template>
                </QuickActionButton>
                <QuickActionButton
                  :label="$t('admin.dashboard.review_comments')"
                  class="cursor-pointer"
                  @action="router.push('/admin/comments')">
                  <template #icon><UIcon name="i-tabler-message" /></template>
                </QuickActionButton>
                <QuickActionButton
                  :label="$t('admin.dashboard.add_user')"
                  class="cursor-pointer"
                  @action="router.push('/admin/users/new')">
                  <template #icon><UIcon name="i-tabler-user-plus" /></template>
                </QuickActionButton>
              </div>
            </UCard>

            <!-- 最近文章 -->
            <UCard>
              <template #header>
                <div class="flex items-center justify-between">
                  <h2 class="text-base font-semibold text-highlighted">
                    {{ $t('admin.dashboard.recent_posts') }}
                  </h2>
                  <UButton
                    as="NuxtLink"
                    to="/admin/posts"
                    color="primary"
                    variant="link"
                    size="sm">
                    {{ $t('common.view_all') }}
                  </UButton>
                </div>
              </template>
              <div v-if="recentPosts.length" class="space-y-1">
                <div
                  v-for="post in recentPosts"
                  :key="post.id"
                  class="flex items-start gap-4 p-3 rounded-md hover:bg-elevated transition-colors group">
                  <div class="h-16 w-24 rounded overflow-hidden bg-muted shrink-0">
                    <BaseImg
                      :src="post.featured_img?.trim() || defaultCover"
                      :alt="post.title"
                      class="w-full h-full object-cover" />
                  </div>
                  <div class="flex-1 min-w-0">
                    <h3 class="font-medium text-highlighted truncate text-sm">
                      {{ post.title }}
                    </h3>
                    <div class="flex items-center gap-2 mt-1 text-xs text-muted">
                      <UBadge
                        :label="postStatusLabel(post.status)"
                        :color="postStatusColor(post.status)"
                        variant="soft"
                        size="xs" />
                      <span>{{ formatDate(post.updated_at) }}</span>
                      <span>·</span>
                      <span>{{ $t('common.views', { n: post.view_count }) }}</span>
                    </div>
                  </div>
                  <UButton
                    :to="`/admin/posts/edit/${post.id}`"
                    as="NuxtLink"
                    icon="i-tabler-pencil"
                    color="neutral"
                    variant="ghost"
                    size="xs"
                    square
                    class="shrink-0 opacity-0 group-hover:opacity-100 transition-opacity" />
                </div>
              </div>
              <div v-else class="py-8 text-center text-sm text-muted">
                {{ $t('admin.dashboard.no_posts') }}
              </div>
            </UCard>
          </div>

          <!-- 右侧 -->
          <div class="space-y-4 md:space-y-6">
            <!-- 最新评论 -->
            <UCard>
              <template #header>
                <div class="flex items-center justify-between">
                  <h2 class="text-base font-semibold text-highlighted">
                    {{ $t('admin.dashboard.recent_comments') }}
                  </h2>
                  <UButton
                    as="NuxtLink"
                    to="/admin/comments"
                    color="primary"
                    variant="link"
                    size="sm">
                    {{ $t('common.view_all') }}
                  </UButton>
                </div>
              </template>
              <div v-if="recentComments.length" class="space-y-4">
                <div
                  v-for="c in recentComments"
                  :key="c.comment_id"
                  class="flex gap-3">
                  <BaseAvatar
                    :src="c.author?.avatar"
                    :alt="c.author?.name || $t('common.anonymous')"
                    size="sm"
                    class="shrink-0" />
                  <div class="flex-1 min-w-0">
                    <div class="flex items-center gap-2 mb-0.5">
                      <span class="text-sm font-medium text-highlighted truncate">{{ c.author?.name || $t('common.anonymous') }}</span>
                      <UBadge
                        v-if="c.status === 'pending'"
                        :label="$t('common.pending')"
                        color="warning"
                        variant="soft"
                        size="xs" />
                    </div>
                    <p class="text-xs text-muted line-clamp-2">
                      {{ c.content }}
                    </p>
                    <p class="text-xs text-muted mt-1">
                      {{ formatDate(c.created_at) }}
                    </p>
                  </div>
                </div>
              </div>
              <div v-else class="py-8 text-center text-sm text-muted">
                {{ $t('admin.dashboard.no_comments') }}
              </div>
            </UCard>
          </div>
        </div>
      </template>
    </AdminPageContent>
  </AdminPageContainer>
</template>

<script setup lang="ts">
const router = useRouter();
const authStore = useAuthStore();
const { apiFetch } = useApiFetch();
const { defaultCover } = usePostCover();
const { t } = useI18n();

const displayName = computed(
  () => authStore.user?.display_name || authStore.user?.username || t('common.admin'),
);

// ── Types ───────────────────────────────────────────────────────────────────

interface PostListItem {
  id: number;
  slug: string;
  title: string;
  status: number;
  view_count: number;
  updated_at: string;
  featured_img?: { url: string };
}

interface PostListRes {
  total: number;
  data: PostListItem[];
}

interface CommentAdminItem {
  comment_id: number;
  content: string;
  status: string;
  created_at: string;
  author?: { name: string; avatar?: string };
}

interface CommentAdminRes {
  total: number;
  list: CommentAdminItem[];
}

interface PostStatsRes {
  total_posts: number;
  published_posts: number;
  draft_posts: number;
  private_posts: number;
  archived_posts: number;
  total_views: number;
  total_likes: number;
  total_comments: number;
}

interface UserListRes {
  total: number;
}

// ── Data fetching ────────────────────────────────────────────────────────────

const { data, pending } = await useAsyncData("dashboard", async () => {
  const [
    postStats,
    recentPostsRes,
    pendingComments,
    allComments,
    allUsers,
    activeUsers,
  ] = await Promise.all([
    apiFetch<PostStatsRes>("/posts/stats").catch(() => null),
    apiFetch<PostListRes>("/posts?page=1&page_size=5").catch(() => null),
    apiFetch<CommentAdminRes>("/admin/comments?page=1&size=1&status=pending").catch(() => null),
    apiFetch<CommentAdminRes>("/admin/comments?page=1&size=5").catch(() => null),
    apiFetch<UserListRes>("/users?page=1&size=1").catch(() => null),
    apiFetch<UserListRes>("/users?page=1&size=1&status=1").catch(() => null),
  ]);
  return { postStats, recentPostsRes, pendingComments, allComments, allUsers, activeUsers };
});

const postStats = computed(() => ({
  total: data.value?.postStats?.total_posts ?? 0,
  published: data.value?.postStats?.published_posts ?? 0,
  draft: data.value?.postStats?.draft_posts ?? 0,
}));

const commentStats = computed(() => ({
  total: data.value?.allComments?.total ?? 0,
  pending: data.value?.pendingComments?.total ?? 0,
}));

const userStats = computed(() => ({
  total: data.value?.allUsers?.total ?? 0,
  active: data.value?.activeUsers?.total ?? 0,
}));

const recentPosts = computed(() =>
  (data.value?.recentPostsRes?.data ?? []).map((p) => ({
    ...p,
    featured_img: p.featured_img?.url,
  })),
);

const recentComments = computed(() => data.value?.allComments?.list ?? []);

const totalViews = computed(() => data.value?.postStats?.total_views ?? 0);


// ── Helpers ──────────────────────────────────────────────────────────────────


const formatDate = (s: string): string => {
  const diff = Date.now() - new Date(s).getTime();
  const m = Math.floor(diff / 60000);
  if (m < 1) return t('common.just_now');
  if (m < 60) return t('common.minutes_ago', { n: m });
  const h = Math.floor(m / 60);
  if (h < 24) return t('common.hours_ago', { n: h });
  const d = Math.floor(h / 24);
  if (d < 30) return t('common.days_ago', { n: d });
  return new Date(s).toLocaleDateString(undefined, {
    month: "short",
    day: "numeric",
  });
};

const postStatusLabel = (status: number) =>
  ({
    1: t('admin.posts.status_draft'),
    2: t('admin.posts.status_published'),
    3: t('admin.posts.status_private'),
    4: t('admin.posts.status_archived'),
  })[status] ?? t('common.unknown');

type BadgeColor = "neutral" | "success" | "warning" | "error";
const postStatusColor = (status: number): BadgeColor =>
  (({ 1: "neutral", 2: "success", 3: "warning", 4: "neutral" })[
    status
  ] as BadgeColor) ?? "neutral";
</script>

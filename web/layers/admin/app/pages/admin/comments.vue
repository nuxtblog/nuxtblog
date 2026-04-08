<template>
  <AdminPageContainer>
    <AdminPageHeader :title="$t('admin.comments.title')" :subtitle="$t('admin.comments.subtitle')">
      <template #actions>
        <div class="flex items-center gap-4">
          <div class="text-center px-4 py-2 bg-default rounded-md border border-default">
            <div class="text-lg font-semibold text-highlighted">{{ stats.total }}</div>
            <div class="text-xs text-muted">{{ $t('admin.comments.total') }}</div>
          </div>
          <div class="text-center px-4 py-2 bg-warning/10 rounded-md border border-warning/20">
            <div class="text-lg font-semibold text-warning">{{ stats.pending }}</div>
            <div class="text-xs text-warning">{{ $t('admin.comments.pending') }}</div>
          </div>
        </div>
      </template>
    </AdminPageHeader>

    <AdminPageContent>
      <!-- 筛选工具栏 -->
      <div class="pb-4 border-b border-default">
        <div class="flex flex-col md:flex-row gap-3 mb-3">
          <!-- 状态 Tabs -->
          <div class="flex items-center gap-1 flex-1 overflow-x-auto border-b border-default pb-0">
            <button
              v-for="s in statusFilters"
              :key="s.value"
              class="px-3 py-2 text-sm font-medium rounded-t transition-colors whitespace-nowrap"
              :class="filterStatus === s.value ? 'text-primary border-b-2 border-primary' : 'text-muted hover:text-highlighted'"
              @click="filterStatus = s.value; currentPage = 1">
              {{ s.label }}
              <span v-if="s.count !== undefined" class="ml-1 text-xs text-muted">({{ s.count }})</span>
            </button>
          </div>
          <!-- 搜索 -->
          <UInput
            v-model="searchQuery"
            :placeholder="$t('admin.comments.search_placeholder')"
            leading-icon="i-tabler-search"
            class="w-56 shrink-0"
            size="sm" />
        </div>

      </div>

      <!-- 评论列表 -->
      <div class="mt-4">
        <!-- 加载骨架 -->
        <div v-if="loading" class="space-y-3">
          <div v-for="i in 5" :key="i" class="flex gap-4 p-4 border border-default rounded-md">
            <USkeleton class="size-4 rounded shrink-0 mt-1" />
            <USkeleton class="size-10 rounded-full shrink-0" />
            <div class="flex-1 space-y-2">
              <div class="flex gap-3">
                <USkeleton class="h-4 w-24" />
                <USkeleton class="h-3 w-32" />
              </div>
              <USkeleton class="h-3 w-full" />
              <USkeleton class="h-3 w-3/4" />
            </div>
          </div>
        </div>

        <!-- 评论列表 -->
        <div v-else-if="comments.length > 0" class="space-y-2">
          <div
            v-for="comment in comments"
            :key="comment.comment_id"
            class="group relative flex gap-3 px-4 py-3 rounded-md border border-default hover:bg-elevated transition-colors">

            <!-- 复选框 -->
            <UCheckbox
              class="shrink-0 mt-1"
              :model-value="selectedComments.includes(comment.comment_id)"
              @update:model-value="toggleSelect(comment.comment_id)" />

            <!-- 头像 -->
            <BaseAvatar
              :src="comment.author?.avatar"
              :alt="comment.author?.name"
              size="sm"
              class="shrink-0 mt-0.5" />

            <!-- 主体 -->
            <div class="flex-1 min-w-0">
              <!-- 作者行 -->
              <div class="flex items-center gap-1.5 flex-wrap mb-1">
                <span class="font-medium text-sm text-highlighted">{{ comment.author?.name || $t('common.anonymous_user') }}</span>
                <span class="text-xs text-muted">·</span>
                <span v-if="comment.author?.email" class="text-xs text-muted">{{ comment.author.email }}</span>
                <span class="text-xs text-muted">{{ formatDate(comment.created_at) }}</span>
              </div>

              <!-- 内容 -->
              <p class="text-sm text-highlighted whitespace-pre-wrap mb-2" v-html="renderComment(comment.content)" />

              <!-- 子评论 -->
              <div v-if="comment.children?.length" class="mb-2 space-y-2 pl-3 border-l-2 border-default">
                <div
                  v-for="child in comment.children"
                  :key="child.comment_id"
                  class="group/child flex gap-2">
                  <BaseAvatar
                    :src="child.author?.avatar"
                    :alt="child.author?.name"
                    size="xs"
                    class="shrink-0 mt-0.5" />
                  <div class="flex-1 min-w-0">
                    <div class="flex items-center gap-1.5 mb-0.5">
                      <span class="text-xs font-medium text-highlighted">{{ child.author?.name || $t('common.anonymous_user') }}</span>
                      <UBadge
                        :color="({ approved: 'success', pending: 'warning', spam: 'neutral', trash: 'error' } as any)[child.status] || 'neutral'"
                        variant="subtle"
                        size="xs">
                        {{ getStatusLabel(child.status) }}
                      </UBadge>
                      <span class="text-xs text-muted">{{ formatDate(child.created_at) }}</span>
                    </div>
                    <p class="text-xs text-highlighted mb-1" v-html="renderComment(child.content)" />
                    <!-- 子评论操作（hover 显示） -->
                    <div class="flex items-center gap-0.5 opacity-0 group-hover/child:opacity-100 transition-opacity">
                      <UButton color="neutral" variant="link" size="xs" class="h-auto px-0.5" @click="openEditModal(child)">{{ $t('admin.comments.edit') }}</UButton>
                      <span class="text-muted text-xs">·</span>
                      <UButton color="neutral" variant="link" size="xs" class="h-auto px-0.5" @click="openReplyModal(child, comment.post_id, comment.comment_id)">{{ $t('admin.comments.reply_action') }}</UButton>
                      <span class="text-muted text-xs">·</span>
                      <UDropdownMenu :items="[[
                        { label: $t('admin.comments.approve'), icon: 'i-tabler-circle-check', disabled: child.status === 'approved', onClick: () => updateStatus(child.comment_id, 'approved') },
                        { label: $t('admin.comments.batch_pending'), icon: 'i-tabler-clock', disabled: child.status === 'pending', onClick: () => updateStatus(child.comment_id, 'pending') },
                        { label: $t('admin.comments.mark_spam'), icon: 'i-tabler-ban', disabled: child.status === 'spam', onClick: () => updateStatus(child.comment_id, 'spam') },
                        { label: $t('admin.comments.move_trash'), icon: 'i-tabler-trash', disabled: child.status === 'trash', onClick: () => updateStatus(child.comment_id, 'trash') },
                      ]]">
                        <UButton color="neutral" variant="link" size="xs" class="h-auto px-0.5">{{ $t('admin.comments.moderate') }}</UButton>
                      </UDropdownMenu>
                      <span class="text-muted text-xs">·</span>
                      <UButton color="error" variant="link" size="xs" class="h-auto px-0.5" @click="deleteConfirmComment = child">{{ $t('admin.comments.delete') }}</UButton>
                    </div>
                  </div>
                </div>
              </div>

              <!-- 主评论操作（hover 显示） -->
              <div class="flex items-center gap-0.5 opacity-0 group-hover:opacity-100 transition-opacity">
                <UButton color="neutral" variant="link" size="xs" class="h-auto px-0.5" @click="openEditModal(comment)">{{ $t('admin.comments.edit') }}</UButton>
                <span class="text-muted text-xs">·</span>
                <UButton color="neutral" variant="link" size="xs" class="h-auto px-0.5" @click="openReplyModal(comment, comment.post_id)">{{ $t('admin.comments.reply_action') }}</UButton>
                <span class="text-muted text-xs">·</span>
                <UDropdownMenu :items="[[
                  { label: $t('admin.comments.approve'), icon: 'i-tabler-circle-check', disabled: comment.status === 'approved', onClick: () => updateStatus(comment.comment_id, 'approved') },
                  { label: $t('admin.comments.batch_pending'), icon: 'i-tabler-clock', disabled: comment.status === 'pending', onClick: () => updateStatus(comment.comment_id, 'pending') },
                  { label: $t('admin.comments.mark_spam'), icon: 'i-tabler-ban', disabled: comment.status === 'spam', onClick: () => updateStatus(comment.comment_id, 'spam') },
                  { label: $t('admin.comments.move_trash'), icon: 'i-tabler-trash', disabled: comment.status === 'trash', onClick: () => updateStatus(comment.comment_id, 'trash') },
                ]]">
                  <UButton color="neutral" variant="link" size="xs" class="h-auto px-0.5">{{ $t('admin.comments.moderate') }}</UButton>
                </UDropdownMenu>
                <span class="text-muted text-xs">·</span>
                <UButton color="error" variant="link" size="xs" class="h-auto px-0.5" @click="deleteConfirmComment = comment">{{ $t('admin.comments.delete') }}</UButton>
              </div>
            </div>

            <!-- 状态 badge + 3点菜单（右侧） -->
            <div class="flex items-center gap-3 shrink-0 self-start">
              <UBadge
                :color="({ approved: 'success', pending: 'warning', spam: 'neutral', trash: 'error' } as any)[comment.status] || 'neutral'"
                variant="soft"
                size="sm">
                {{ getStatusLabel(comment.status) }}
              </UBadge>
              <UDropdownMenu :items="[[
                { label: $t('admin.comments.edit_post'), icon: 'i-tabler-pencil', to: `/admin/posts/edit/${comment.post_id}` },
                { label: $t('admin.comments.view_post'), icon: 'i-tabler-external-link', onClick: () => openUrl(`/p/${comment.post_id}`) },
              ]]">
                <UButton
                  icon="i-tabler-dots-vertical"
                  color="neutral"
                  variant="ghost"
                  size="xs"
                  square
                  class="opacity-0 group-hover:opacity-100 transition-opacity" />
              </UDropdownMenu>
            </div>
          </div>
        </div>

        <!-- 空状态 -->
        <div v-else class="flex flex-col items-center justify-center py-16">
          <UIcon name="i-tabler-message-off" class="size-16 text-muted mb-4" />
          <h3 class="text-lg font-medium text-highlighted mb-1">{{ searchQuery ? $t('admin.comments.no_results') : $t('admin.comments.no_comments') }}</h3>
        </div>
      </div>

      <!-- 编辑评论 -->
      <UModal :open="!!editingComment" @update:open="val => { if (!val) editingComment = null }">
        <template #content>
          <div class="p-6">
            <h3 class="text-base font-semibold text-highlighted mb-4">{{ $t('admin.comments.edit_title') }}</h3>
            <form @submit.prevent="handleUpdateComment" class="space-y-4">
              <UFormField :label="$t('admin.comments.content_label')" required>
                <UTextarea v-model="editFormData.content" :rows="6" class="w-full" />
              </UFormField>
              <UFormField :label="$t('admin.comments.status_label')">
                <USelect
                  v-model="editFormData.status"
                  :items="[
                    { label: $t('admin.comments.status_pending_label'), value: 'pending' },
                    { label: $t('admin.comments.status_approved_label'), value: 'approved' },
                    { label: $t('admin.comments.status_spam_label'), value: 'spam' },
                    { label: $t('admin.comments.status_trash_label'), value: 'trash' },
                  ]"
                  class="w-full" />
              </UFormField>
              <div class="flex gap-3 justify-end">
                <UButton color="neutral" variant="outline" @click="editingComment = null">{{ $t('common.cancel') }}</UButton>
                <UButton type="submit" color="primary" :loading="updating">{{ $t('common.save') }}</UButton>
              </div>
            </form>
          </div>
        </template>
      </UModal>

      <!-- 删除确认 -->
      <UModal :open="!!deleteConfirmComment" @update:open="val => { if (!val) deleteConfirmComment = null }">
        <template #content>
          <div class="p-6">
            <h3 class="text-base font-semibold text-highlighted mb-2">{{ $t('admin.comments.delete_title') }}</h3>
            <p class="text-sm text-muted mb-4">
              {{ $t('admin.comments.delete_desc') }}
              <span v-if="deleteConfirmComment?.children?.length" class="block mt-1 text-warning">
                {{ $t('admin.comments.delete_with_replies', { n: deleteConfirmComment.children.length }) }}
              </span>
            </p>
            <div class="flex gap-3 justify-end">
              <UButton color="neutral" variant="outline" @click="deleteConfirmComment = null">{{ $t('common.cancel') }}</UButton>
              <UButton color="error" :loading="deleting" @click="confirmDelete">{{ $t('common.delete') }}</UButton>
            </div>
          </div>
        </template>
      </UModal>

      <!-- 回复评论 -->
      <UModal :open="!!replyingComment" @update:open="val => { if (!val) replyingComment = null }">
        <template #content>
          <div class="p-6">
            <h3 class="text-base font-semibold text-highlighted mb-1">{{ $t('admin.comments.reply_title') }}</h3>
            <p class="text-sm text-muted mb-4 line-clamp-2">
              {{ $t('admin.comments.reply_desc') }} <span class="font-medium text-highlighted">{{ replyingComment?.author?.name }}</span>：{{ replyingComment?.content }}
            </p>
            <form @submit.prevent="handleReply" class="space-y-4">
              <UFormField :label="$t('admin.comments.content_label')" required>
                <UTextarea v-model="replyContent" :rows="5" :placeholder="$t('admin.comments.reply_placeholder')" class="w-full" />
              </UFormField>
              <div class="flex gap-3 justify-end">
                <UButton color="neutral" variant="outline" @click="replyingComment = null">{{ $t('common.cancel') }}</UButton>
                <UButton type="submit" color="primary" :loading="replying" :disabled="!replyContent.trim()">{{ $t('common.send') }}</UButton>
              </div>
            </form>
          </div>
        </template>
      </UModal>
    </AdminPageContent>

    <AdminPageFooter>
      <template #left>
        <template v-if="comments.length > 0">
          <UCheckbox
            :model-value="comments.length > 0 && comments.every(c => selectedComments.includes(c.comment_id))"
            :indeterminate="selectedComments.length > 0 && !comments.every(c => selectedComments.includes(c.comment_id))"
            @update:model-value="val => { selectedComments = val ? comments.map(c => c.comment_id) : [] }" />
          <template v-if="selectedComments.length > 0">
            <span>{{ $t('admin.comments.selected_n', { n: selectedComments.length }) }}</span>
            <USeparator orientation="vertical" class="h-4" />
            <USelect
              v-model="batchAction"
              :items="[
                { label: $t('admin.comments.batch_approve'), value: 'approved' },
                { label: $t('admin.comments.batch_pending'), value: 'pending' },
                { label: $t('admin.comments.batch_spam'), value: 'spam' },
                { label: $t('admin.comments.batch_trash'), value: 'trash' },
              ]"
              :placeholder="$t('admin.comments.batch_ops')"
              class="w-36"
              size="sm" />
            <UButton color="primary" variant="soft" size="sm" :disabled="!batchAction" @click="handleBatchAction">{{ $t('admin.comments.batch_apply') }}</UButton>
            <UButton color="neutral" variant="ghost" size="sm" @click="selectedComments = []; batchAction = ''">{{ $t('common.cancel') }}</UButton>
          </template>
          <span v-else class="text-xs">{{ $t('common.selectAll') }}</span>
        </template>
      </template>
      <template #right>
        <UPagination v-if="!loading && totalComments > pageSize" v-model:page="currentPage" :total="totalComments" :items-per-page="pageSize" size="sm" />
      </template>
    </AdminPageFooter>
  </AdminPageContainer>
</template>

<script setup lang="ts">
interface CommentAuthor {
  id?: number;
  name: string;
  email?: string;
  avatar?: string;
  is_author?: boolean;
}

interface CommentItem {
  comment_id: number;
  post_id: number;
  content: string;
  parent_id?: number;
  status: string;
  ip_address?: string;
  user_agent?: string;
  created_at: string;
  updated_at: string;
  author?: CommentAuthor;
  children?: CommentItem[];
}

const { apiFetch } = useApiFetch();
const authStore = useAuthStore();
const toast = useToast();
const { renderComment } = useCommentRender();
const { t } = useI18n();

const comments = ref<CommentItem[]>([]);
const stats = ref({ total: 0, approved: 0, pending: 0, spam: 0, trash: 0 });
const statusFilters = computed(() => [
  { label: t("admin.comments.status_all"), value: "all", count: stats.value.total },
  { label: t("admin.comments.status_approved"), value: "approved", count: stats.value.approved },
  { label: t("admin.comments.status_pending"), value: "pending", count: stats.value.pending },
  { label: t("admin.comments.status_spam"), value: "spam", count: stats.value.spam },
  { label: t("admin.comments.status_trash"), value: "trash", count: stats.value.trash },
]);

const rawLoading = ref(true);
const loading = useMinLoading(rawLoading);
const updating = ref(false);
const deleting = ref(false);
const replying = ref(false);

const filterStatus = ref("all");
const searchQuery = ref("");
const currentPage = ref(1);
const pageSize = 20;
const totalComments = ref(0);

const selectedComments = ref<number[]>([]);
const batchAction = ref("");

const editingComment = ref<CommentItem | null>(null);
const editFormData = ref({ content: "", status: "pending" });
const deleteConfirmComment = ref<CommentItem | null>(null);
const replyingComment = ref<CommentItem | null>(null);
const replyParentId = ref<number | null>(null); // root parent id for flat threading
const replyContent = ref("");

// ── Data ──────────────────────────────────────────────────────────────────────

const fetchComments = async () => {
  try {
    const [main, allC, approvedC, pendingC, spamC, trashC] = await Promise.all([
      apiFetch<{ list: CommentItem[]; total: number }>("/admin/comments", {
        params: {
          page: currentPage.value,
          size: pageSize,
          status: filterStatus.value !== "all" ? filterStatus.value : undefined,
          keyword: searchQuery.value || undefined,
        },
      }),
      apiFetch<{ total: number }>("/admin/comments", { params: { size: 1 } }),
      apiFetch<{ total: number }>("/admin/comments", { params: { size: 1, status: "approved" } }),
      apiFetch<{ total: number }>("/admin/comments", { params: { size: 1, status: "pending" } }),
      apiFetch<{ total: number }>("/admin/comments", { params: { size: 1, status: "spam" } }),
      apiFetch<{ total: number }>("/admin/comments", { params: { size: 1, status: "trash" } }),
    ]);

    comments.value = main.list;
    totalComments.value = main.total;

    stats.value = {
      total: allC.total,
      approved: approvedC.total,
      pending: pendingC.total,
      spam: spamC.total,
      trash: trashC.total,
    };
  } finally {
    rawLoading.value = false;
  }
};

// ── Actions ───────────────────────────────────────────────────────────────────

const statusIntMap: Record<string, number> = { pending: 1, approved: 2, spam: 3, trash: 4 };

const updateStatus = async (id: number, status: string) => {
  try {
    await apiFetch(`/comments/${id}/status`, { method: "PUT", body: { status: statusIntMap[status] } });
    toast.add({ title: t("admin.comments.status_updated"), color: "success", icon: "i-tabler-circle-check" });
    await fetchComments();
  } catch {
    toast.add({ title: t("common.operation_failed"), color: "error", icon: "i-tabler-alert-circle" });
  }
};

const handleBatchAction = async () => {
  if (!batchAction.value || !selectedComments.value.length) return;
  try {
    await Promise.all(
      selectedComments.value.map((id) =>
        apiFetch(`/comments/${id}/status`, { method: "PUT", body: { status: statusIntMap[batchAction.value] } })
      )
    );
    toast.add({ title: t("admin.comments.batch_updated", { n: selectedComments.value.length }), color: "success", icon: "i-tabler-circle-check" });
    selectedComments.value = [];
    batchAction.value = "";
    await fetchComments();
  } catch {
    toast.add({ title: t("admin.comments.batch_failed"), color: "error", icon: "i-tabler-alert-circle" });
  }
};

const toggleSelect = (id: number) => {
  const idx = selectedComments.value.indexOf(id);
  if (idx > -1) selectedComments.value.splice(idx, 1);
  else selectedComments.value.push(id);
};

const openEditModal = (comment: CommentItem) => {
  editingComment.value = comment;
  editFormData.value = { content: comment.content, status: comment.status };
};

const handleUpdateComment = async () => {
  if (!editingComment.value) return;
  updating.value = true;
  try {
    const id = editingComment.value.comment_id;
    await apiFetch(`/comments/${id}`, { method: "PUT", body: { content: editFormData.value.content } });
    if (editFormData.value.status !== editingComment.value.status) {
      await apiFetch(`/comments/${id}/status`, { method: "PUT", body: { status: statusIntMap[editFormData.value.status] } });
    }
    toast.add({ title: t("admin.comments.updated"), color: "success", icon: "i-tabler-circle-check" });
    editingComment.value = null;
    await fetchComments();
  } catch {
    toast.add({ title: t("common.save_failed"), color: "error", icon: "i-tabler-alert-circle" });
  } finally {
    updating.value = false;
  }
};

const confirmDelete = async () => {
  if (!deleteConfirmComment.value) return;
  deleting.value = true;
  try {
    await apiFetch(`/comments/${deleteConfirmComment.value.comment_id}`, { method: "DELETE" });
    toast.add({ title: t("admin.comments.deleted"), color: "success", icon: "i-tabler-circle-check" });
    deleteConfirmComment.value = null;
    await fetchComments();
  } catch {
    toast.add({ title: t("common.delete_failed"), color: "error", icon: "i-tabler-alert-circle" });
  } finally {
    deleting.value = false;
  }
};

// rootParentId: when replying to a child, pass root comment id to keep thread flat (max 2 levels)
const openReplyModal = (comment: CommentItem, postId: number, rootParentId?: number) => {
  replyingComment.value = { ...comment, post_id: postId };
  replyParentId.value = rootParentId ?? comment.comment_id;
  // pre-fill @mention when replying to a child
  replyContent.value = rootParentId ? `@${comment.author?.name || t('common.anonymous_user')} ` : "";
};

const handleReply = async () => {
  if (!replyingComment.value || !replyContent.value.trim()) return;
  replying.value = true;
  try {
    await apiFetch("/comments", {
      method: "POST",
      body: {
        object_id: replyingComment.value.post_id,
        object_type: "post",
        parent_id: replyParentId.value,
        author_name: authStore.user?.display_name || authStore.user?.username || t('common.admin'),
        author_email: authStore.user?.email || "admin@example.com",
        content: replyContent.value.trim(),
      },
    });
    toast.add({ title: t("admin.comments.replied"), color: "success", icon: "i-tabler-circle-check" });
    replyingComment.value = null;
    await fetchComments();
  } catch {
    toast.add({ title: t("admin.comments.reply_failed"), color: "error", icon: "i-tabler-alert-circle" });
  } finally {
    replying.value = false;
  }
};

const openUrl = (url: string) => {
  if (import.meta.client) window.open(url, "_blank");
};

// ── Helpers ───────────────────────────────────────────────────────────────────

const getStatusLabel = (s: string) =>
  ({
    approved: t("admin.comments.label_approved"),
    pending: t("admin.comments.label_pending"),
    spam: t("admin.comments.label_spam"),
    trash: t("admin.comments.label_trash"),
  })[s] ?? s;

const formatDate = (s: string) => {
  const diff = Date.now() - new Date(s).getTime();
  const m = Math.floor(diff / 60000);
  if (m < 1) return t('common.just_now');
  if (m < 60) return t('common.minutes_ago', { n: m });
  const h = Math.floor(m / 60);
  if (h < 24) return t('common.hours_ago', { n: h });
  const d = Math.floor(h / 24);
  if (d < 30) return t('common.days_ago', { n: d });
  return new Date(s).toLocaleDateString(undefined, { month: "short", day: "numeric" });
};

// ── Watchers ──────────────────────────────────────────────────────────────────

let searchTimer: ReturnType<typeof setTimeout>;
watch(searchQuery, () => {
  clearTimeout(searchTimer);
  searchTimer = setTimeout(() => {
    currentPage.value = 1;
    selectedComments.value = [];
    batchAction.value = '';
    fetchComments();
  }, 400);
});

watch(filterStatus, () => {
  currentPage.value = 1;
  selectedComments.value = [];
  batchAction.value = '';
  fetchComments();
});
watch(currentPage, fetchComments);

onMounted(fetchComments);
</script>

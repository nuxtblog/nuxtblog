<template>
  <div v-if="!commentsOpen">
    <p class="text-sm text-muted">{{ $t('site.comment.closed') }}</p>
  </div>

  <div v-else>
    <!-- 标题 -->
    <h2 class="text-lg font-semibold text-highlighted mb-8">
      {{ total > 0 ? $t('site.comment.title_with_count', { n: total }) : $t('site.comment.title') }}
    </h2>

    <!-- 评论列表 -->
    <div v-if="loading" class="space-y-8 mb-10">
      <div v-for="i in 3" :key="i" class="flex gap-4">
        <USkeleton class="size-8 rounded-full shrink-0" />
        <div class="flex-1 space-y-2 pt-1">
          <USkeleton class="h-3 w-28" />
          <USkeleton class="h-3 w-full" />
          <USkeleton class="h-3 w-2/3" />
        </div>
      </div>
    </div>

    <div v-else-if="comments.length" class="space-y-8 mb-10">
      <div v-for="comment in comments" :key="comment.id" class="flex gap-4">
        <BaseAvatar
          :alt="comment.author_name"
          size="sm"
          class="shrink-0 mt-0.5 ring-2 ring-default" />

        <div class="flex-1 min-w-0">
          <div class="flex items-baseline gap-2 mb-1.5">
            <span class="text-sm font-medium text-highlighted">{{ comment.author_name }}</span>
            <span class="text-xs text-muted">{{ fmtDate(comment.created_at) }}</span>
          </div>
          <p
            class="text-sm text-default leading-relaxed whitespace-pre-wrap"
            v-html="renderComment(comment.content)" />
          <div class="mt-2 flex items-center gap-3">
            <button
              class="text-xs text-muted hover:text-primary transition-colors"
              @click="startReply(comment)">
              {{ $t('site.comment.reply') }}
            </button>
            <ReportButton v-if="authStore.isLoggedIn" target-type="comment" :target-id="comment.id" />
          </div>

          <!-- 子评论 -->
          <div v-if="comment.replies?.length" class="mt-6 space-y-6 pl-5 border-l border-default">
            <div v-for="reply in comment.replies" :key="reply.id" class="flex gap-3">
              <BaseAvatar
                :alt="reply.author_name"
                size="xs"
                class="shrink-0 mt-0.5 ring-1 ring-default" />
              <div class="flex-1 min-w-0">
                <div class="flex items-baseline gap-2 mb-1.5">
                  <span class="text-sm font-medium text-highlighted">{{ reply.author_name }}</span>
                  <span class="text-xs text-muted">{{ fmtDate(reply.created_at) }}</span>
                </div>
                <p
                  class="text-sm text-default leading-relaxed whitespace-pre-wrap"
                  v-html="renderComment(reply.content)" />
                <div class="mt-2 flex items-center gap-3">
                  <button
                    class="text-xs text-muted hover:text-primary transition-colors"
                    @click="startReply(reply, comment.id)">
                    {{ $t('site.comment.reply') }}
                  </button>
                  <ReportButton v-if="authStore.isLoggedIn" target-type="comment" :target-id="reply.id" />
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <p v-else class="text-sm text-muted mb-10">{{ $t('site.comment.no_comments') }}</p>

    <!-- 分页 -->
    <div v-if="total > pageSize" class="flex justify-center mb-10">
      <UPagination v-model:page="page" :total="total" :items-per-page="pageSize" size="sm" />
    </div>

    <!-- 提交表单（ClientOnly 避免 SSR/hydration value 不一致） -->
    <ClientOnly>
      <div class="space-y-4">
        <!-- 已登录身份栏 -->
        <div v-if="authStore.isLoggedIn" class="flex items-center gap-2 text-sm text-muted">
          <BaseAvatar
            :src="authStore.user?.avatar"
            :alt="authStore.user?.display_name || authStore.user?.username"
            size="xs" />
          <span>{{ $t('site.comment.posting_as', { name: authStore.user?.display_name }) }}</span>
        </div>

        <!-- 回复提示 -->
        <div v-if="replyTo" class="flex items-center gap-2 text-xs text-muted">
          <UIcon name="i-tabler-corner-down-right" class="size-3.5 shrink-0" />
          <span>{{ $t('site.comment.reply_to', { name: replyTo.author_name }) }}</span>
          <button class="ml-auto hover:text-highlighted transition-colors" @click="cancelReply">
            <UIcon name="i-tabler-x" class="size-3.5" />
          </button>
        </div>

        <!-- 游客姓名/邮箱 -->
        <div v-if="!authStore.isLoggedIn && requireNameEmail" class="grid grid-cols-2 gap-3">
          <UInput v-model="form.name" :placeholder="$t('site.comment.name_placeholder')" size="sm" />
          <UInput v-model="form.email" type="email" :placeholder="$t('site.comment.email_placeholder')" size="sm" />
        </div>

        <!-- 内容 -->
        <UTextarea
          ref="textareaRef"
          v-model="form.content"
          :placeholder="$t('site.comment.content_placeholder')"
          :rows="3"
          class="w-full resize-none" />

        <div class="flex items-center justify-between">
          <p v-if="moderation" class="text-xs text-muted">{{ $t('site.comment.moderation_hint') }}</p>
          <span v-else />
          <UButton
            color="primary"
            size="sm"
            :loading="submitting"
            :disabled="submitting || !form.content.trim()"
            @click="submitComment">
            {{ $t('site.comment.submit') }}
          </UButton>
        </div>
      </div>
    </ClientOnly>
  </div>
</template>

<script setup lang="ts">
interface CommentItem {
  id: number;
  author_name: string;
  content: string;
  created_at: string;
  replies?: CommentItem[];
}

const props = defineProps<{
  objectId: number;
  objectType?: string;
  commentMeta?: string; // "1" open | "0" closed | undefined = global default
}>();

const { t } = useI18n()
const { apiFetch } = useApiFetch();
const optionsStore = useOptionsStore();
const authStore = useAuthStore();
const { renderComment } = useCommentRender();
const toast = useToast();
const textareaRef = ref<HTMLElement | null>(null);

// ── Settings ─────────────────────────────────────────────────────────────────

const globalAllow = computed(() => optionsStore.getJSON<boolean>("default_allow_comments", true));
const moderation = computed(() => optionsStore.getJSON<boolean>("comment_moderation", false));
const requireNameEmail = computed(() => optionsStore.getJSON<boolean>("comment_require_name_email", true));

const commentsOpen = computed(() => {
  if (props.commentMeta === "0") return false;
  if (props.commentMeta === "1") return true;
  return globalAllow.value;
});

const objectType = computed(() => props.objectType ?? "post");

// ── Fetch ─────────────────────────────────────────────────────────────────────

const comments = ref<CommentItem[]>([]);
const total = ref(0);
const page = ref(1);
const pageSize = 20;
const loading = ref(false);

const fetchComments = async () => {
  if (!commentsOpen.value) return;
  loading.value = true;
  try {
    const res = await apiFetch<{ list: CommentItem[]; total: number }>("/comments", {
      params: { object_id: props.objectId, object_type: objectType.value, page: page.value, size: pageSize },
    });
    comments.value = res.list ?? [];
    total.value = res.total ?? 0;
  } finally {
    loading.value = false;
  }
};

watch(page, fetchComments);
onMounted(fetchComments);

// ── Reply ─────────────────────────────────────────────────────────────────────

const replyTo = ref<CommentItem | null>(null);
const replyParentId = ref<number | null>(null);

const startReply = (comment: CommentItem, rootId?: number) => {
  replyTo.value = comment;
  replyParentId.value = rootId ?? comment.id;
  form.value.content = `@${comment.author_name} `;
  nextTick(() => (textareaRef.value as any)?.$el?.querySelector("textarea")?.focus());
};

const cancelReply = () => {
  replyTo.value = null;
  replyParentId.value = null;
  form.value.content = "";
};

// ── Submit ────────────────────────────────────────────────────────────────────

const form = ref({ name: "", email: "", content: "" });
const submitting = ref(false);

const submitComment = async () => {
  if (!form.value.content.trim()) return;
  if (!authStore.isLoggedIn && requireNameEmail.value) {
    if (!form.value.name.trim() || !form.value.email.trim()) {
      toast.add({ title: t("site.comment.fill_name_email"), color: "warning", icon: "i-tabler-alert-circle" });
      return;
    }
  }

  const authorName = authStore.isLoggedIn
    ? (authStore.user?.display_name || authStore.user?.username || "用户")
    : (form.value.name.trim() || "匿名");
  const authorEmail = authStore.isLoggedIn
    ? (authStore.user?.email || "")
    : form.value.email.trim();

  submitting.value = true;
  try {
    await apiFetch("/comments", {
      method: "POST",
      body: {
        object_id: props.objectId,
        object_type: objectType.value,
        parent_id: replyParentId.value ?? undefined,
        author_name: authorName,
        author_email: authorEmail,
        content: form.value.content.trim(),
      },
    });
    toast.add({
      title: moderation.value ? t("site.comment.submitted_moderation") : t("site.comment.submitted"),
      color: "success",
      icon: "i-tabler-circle-check",
    });
    form.value.content = "";
    replyTo.value = null;
    replyParentId.value = null;
    if (!moderation.value) await fetchComments();
  } catch (err) {
    toast.add({
      title: t("site.comment.submit_failed"),
      description: (err as Error).message || "请稍后重试",
      color: "error",
      icon: "i-tabler-alert-circle",
    });
  } finally {
    submitting.value = false;
  }
};

// ── Helpers ───────────────────────────────────────────────────────────────────

const fmtDate = (s: string) => {
  const diff = Date.now() - new Date(s).getTime();
  const m = Math.floor(diff / 60000);
  if (m < 1) return t("site.activity.time_just_now");
  if (m < 60) return t("site.activity.time_minutes_ago", { n: m });
  const h = Math.floor(m / 60);
  if (h < 24) return t("site.activity.time_hours_ago", { n: h });
  const d = Math.floor(h / 24);
  if (d < 30) return t("site.activity.time_days_ago", { n: d });
  return new Date(s).toLocaleDateString("zh-CN", { month: "short", day: "numeric" });
};
</script>

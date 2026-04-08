<script setup lang="ts">
const props = defineProps<{
  postId: number
  likeCount?: number
  commentCount?: number
}>()

const { t } = useI18n()
const toast = useToast()
const authStore = useAuthStore()
const reactionApi = useReactionApi()
const router = useRouter()

// ── Visibility ────────────────────────────────────────────────────────────────
const visible = ref(false)
useEventListener('scroll', () => { visible.value = window.scrollY > 200 })

// ── State ─────────────────────────────────────────────────────────────────────
const liked = ref(false)
const bookmarked = ref(false)
const displayLikes = ref(props.likeCount ?? 0)
const likeAnim = ref(false)
const bookmarkAnim = ref(false)
const loadingLike = ref(false)
const loadingBookmark = ref(false)

const { apiFetch } = useApiFetch()

// Load like count + reaction status on mount
onMounted(async () => {
  // Fetch per-post stats for accurate like count
  try {
    const res = await apiFetch<{ stats?: { like_count?: number } }>(`/posts/${props.postId}`)
    if (res.stats?.like_count != null) displayLikes.value = res.stats.like_count
  } catch {}

  if (authStore.isLoggedIn) {
    try {
      const res = await reactionApi.getReaction(props.postId)
      liked.value = res.liked
      bookmarked.value = res.bookmarked
    } catch {}
  }
})

// ── Like ──────────────────────────────────────────────────────────────────────
const toggleLike = async () => {
  if (!authStore.isLoggedIn) {
    router.push(`/auth/login?redirect=${encodeURIComponent(useRoute().fullPath)}`)
    return
  }
  if (loadingLike.value) return
  loadingLike.value = true
  const wasLiked = liked.value
  // Optimistic update
  liked.value = !wasLiked
  displayLikes.value = liked.value ? displayLikes.value + 1 : Math.max(0, displayLikes.value - 1)
  likeAnim.value = true
  setTimeout(() => (likeAnim.value = false), 600)
  try {
    if (liked.value) {
      await reactionApi.likePost(props.postId)
    } else {
      await reactionApi.unlikePost(props.postId)
    }
  } catch {
    // Revert on error
    liked.value = wasLiked
    displayLikes.value = wasLiked ? displayLikes.value + 1 : Math.max(0, displayLikes.value - 1)
  } finally {
    loadingLike.value = false
  }
}

// ── Bookmark ──────────────────────────────────────────────────────────────────
const toggleBookmark = async () => {
  if (!authStore.isLoggedIn) {
    router.push(`/auth/login?redirect=${encodeURIComponent(useRoute().fullPath)}`)
    return
  }
  if (loadingBookmark.value) return
  loadingBookmark.value = true
  const wasBookmarked = bookmarked.value
  bookmarked.value = !wasBookmarked
  bookmarkAnim.value = true
  setTimeout(() => (bookmarkAnim.value = false), 600)
  try {
    if (bookmarked.value) {
      await reactionApi.bookmarkPost(props.postId)
    } else {
      await reactionApi.unbookmarkPost(props.postId)
    }
    toast.add({
      title: bookmarked.value ? t('site.post.bookmarked') : t('site.post.unbookmarked'),
      icon: bookmarked.value ? 'i-tabler-bookmark-filled' : 'i-tabler-bookmark-off',
      color: 'primary',
      duration: 1500,
    })
  } catch {
    bookmarked.value = wasBookmarked
  } finally {
    loadingBookmark.value = false
  }
}

// ── Share ─────────────────────────────────────────────────────────────────────
const copyLink = async () => {
  const url = window.location.href
  try {
    await navigator.clipboard.writeText(url)
  } catch {
    // Fallback for HTTP (non-secure context)
    const ta = document.createElement('textarea')
    ta.value = url
    ta.style.cssText = 'position:fixed;top:0;left:0;opacity:0'
    document.body.appendChild(ta)
    ta.focus()
    ta.select()
    document.execCommand('copy')
    document.body.removeChild(ta)
  }
  toast.add({ title: t('site.post.link_copied'), icon: 'i-tabler-check', color: 'success', duration: 1500 })
}

const shareToWeibo = () => {
  window.open(`https://service.weibo.com/share/share.php?url=${encodeURIComponent(window.location.href)}`, '_blank')
  toast.add({ title: t('site.post.share_weibo_done'), icon: 'i-tabler-brand-weibo', color: 'primary', duration: 1500 })
}

const shareToTwitter = () => {
  window.open(`https://twitter.com/intent/tweet?url=${encodeURIComponent(window.location.href)}`, '_blank')
  toast.add({ title: t('site.post.share_twitter_done'), icon: 'i-tabler-brand-twitter', color: 'primary', duration: 1500 })
}

const shareItems = computed(() => [[
  { label: t('site.post.copy_link'), icon: 'i-tabler-link', onClick: copyLink },
  { label: t('site.post.share_weibo'), icon: 'i-tabler-brand-weibo', onClick: shareToWeibo },
  { label: t('site.post.share_twitter'), icon: 'i-tabler-brand-twitter', onClick: shareToTwitter },
]])

const scrollToComments = () => {
  const el = document.getElementById('post-comments')
  if (el) el.scrollIntoView({ behavior: 'smooth', block: 'start' })
}

const fmtNum = (n: number) =>
  n >= 10000 ? (n / 10000).toFixed(1) + 'w' : n >= 1000 ? (n / 1000).toFixed(1) + 'k' : String(n)
</script>

<template>
  <Transition
    enter-active-class="transition duration-300 ease-out"
    enter-from-class="opacity-0 -translate-x-3 scale-95"
    enter-to-class="opacity-100 translate-x-0 scale-100"
    leave-active-class="transition duration-200 ease-in"
    leave-from-class="opacity-100 translate-x-0 scale-100"
    leave-to-class="opacity-0 -translate-x-3 scale-95"
  >
    <div
      v-if="visible"
      class="hidden lg:flex fixed left-5 top-1/2 -translate-y-1/2 z-40 flex-col items-center gap-1"
    >
      <div class="flex flex-col items-center gap-0.5 rounded-md bg-default/90 backdrop-blur-xl shadow-2xl shadow-black/10 ring-1 ring-default p-1.5">

        <!-- 点赞 -->
        <UTooltip :text="authStore.isLoggedIn ? (liked ? $t('site.post.unlike') : $t('site.post.like')) : $t('site.post.like_login')" side="right" :delay-duration="100">
          <button
            class="group relative flex flex-col items-center justify-center w-11 py-2 rounded-md transition-all duration-200"
            :class="liked ? 'bg-red-50 dark:bg-red-950/30' : 'hover:bg-red-50 dark:hover:bg-red-950/30'"
            @click="toggleLike"
          >
            <UIcon
              :name="liked ? 'i-tabler-heart-filled' : 'i-tabler-heart'"
              class="size-5 transition-all duration-200"
              :class="[liked ? 'text-red-500' : 'text-muted group-hover:text-red-500', likeAnim ? 'scale-125' : 'scale-100']"
            />
            <span class="text-[10px] font-medium mt-0.5 leading-none" :class="liked ? 'text-red-500' : 'text-muted'">
              {{ fmtNum(displayLikes) }}
            </span>
          </button>
        </UTooltip>

        <!-- 评论 -->
        <UTooltip :text="$t('site.post.comment')" side="right" :delay-duration="100">
          <button
            class="group flex flex-col items-center justify-center w-11 py-2 rounded-md hover:bg-primary/10 transition-all duration-200"
            @click="scrollToComments"
          >
            <UIcon name="i-tabler-message-circle" class="size-5 text-muted group-hover:text-primary transition-colors" />
            <span class="text-[10px] font-medium mt-0.5 leading-none text-muted">{{ fmtNum(commentCount ?? 0) }}</span>
          </button>
        </UTooltip>

        <!-- 收藏 -->
        <UTooltip :text="authStore.isLoggedIn ? (bookmarked ? $t('site.post.unbookmark') : $t('site.post.bookmark')) : $t('site.post.bookmark_login')" side="right" :delay-duration="100">
          <button
            class="group relative flex flex-col items-center justify-center w-11 py-2 rounded-md transition-all duration-200"
            :class="bookmarked ? 'bg-amber-50 dark:bg-amber-950/30' : 'hover:bg-amber-50 dark:hover:bg-amber-950/30'"
            @click="toggleBookmark"
          >
            <UIcon
              :name="bookmarked ? 'i-tabler-bookmark-filled' : 'i-tabler-bookmark'"
              class="size-5 transition-all duration-200"
              :class="[bookmarked ? 'text-amber-500' : 'text-muted group-hover:text-amber-500', bookmarkAnim ? 'scale-125' : 'scale-100']"
            />
            <span class="text-[10px] font-medium mt-0.5 leading-none text-muted">{{ bookmarked ? $t('site.post.bookmarked') : $t('site.post.bookmark') }}</span>
          </button>
        </UTooltip>

        <div class="w-6 h-px bg-border-default border border-default my-0.5 rounded-full" />

        <!-- 分享 -->
        <UDropdownMenu :items="shareItems" :ui="{ content: 'w-40' }">
          <UTooltip :text="$t('site.post.share')" side="right" :delay-duration="100">
            <button class="group flex flex-col items-center justify-center w-11 py-2 rounded-md hover:bg-primary/10 transition-all duration-200">
              <UIcon name="i-tabler-share-2" class="size-5 text-muted group-hover:text-primary transition-colors" />
              <span class="text-[10px] font-medium mt-0.5 leading-none text-muted">{{ $t('site.post.share') }}</span>
            </button>
          </UTooltip>
        </UDropdownMenu>

      </div>
    </div>
  </Transition>
</template>

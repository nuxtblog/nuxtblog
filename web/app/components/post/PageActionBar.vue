<script setup lang="ts">
const props = defineProps<{
  commentCount?: number
}>()

const { t } = useI18n()
const toast = useToast()

// ── Visibility ────────────────────────────────────────────────────────────────
const visible = ref(false)
useEventListener('scroll', () => { visible.value = window.scrollY > 200 })

// ── Share ─────────────────────────────────────────────────────────────────────
const copyLink = async () => {
  const url = window.location.href
  try {
    await navigator.clipboard.writeText(url)
  } catch {
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

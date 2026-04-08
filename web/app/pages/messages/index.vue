<script setup lang="ts">
import type { ConversationItem, MessageItem } from '~/composables/useMessageApi'
definePageMeta({ middleware: 'auth' })

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const messageApi = useMessageApi()
const userApi = useUserApi()
const toast = useToast()

// ── Conversations ─────────────────────────────────────────────────────────
const conversations = ref<ConversationItem[]>([])
const convLoading = ref(true)
const totalUnread = ref(0)

const loadConversations = async () => {
  try {
    const res = await messageApi.listConversations()
    conversations.value = res.items
    totalUnread.value = res.total_unread
  } finally {
    convLoading.value = false
  }
}

// ── Selected conversation ─────────────────────────────────────────────────
const selectedId = computed(() => {
  const w = route.query.with
  return w ? Number(w) : null
})

const otherUser = ref<{ id?: number; display_name?: string; username?: string; avatar?: string } | null>(null)
const messages = ref<MessageItem[]>([])
const hasMore = ref(false)
const msgLoading = ref(false)
const sending = ref(false)
const content = ref('')
const messagesEl = ref<HTMLElement | null>(null)

const displayName = computed(() =>
  otherUser.value?.display_name || otherUser.value?.username || t('site.message.user_default')
)

const openConversation = async (userId: number) => {
  if (selectedId.value === userId) return
  router.replace({ query: { with: userId } })
  otherUser.value = null
  messages.value = []
  hasMore.value = false
  msgLoading.value = true
  try {
    const [u, res] = await Promise.all([
      userApi.getUser(userId).catch(() => null),
      messageApi.listMessages(userId),
    ])
    otherUser.value = u
    messages.value = res.items
    hasMore.value = res.has_more
    await nextTick()
    scrollToBottom()
  } finally {
    msgLoading.value = false
  }
}

const loadMore = async () => {
  if (!messages.value.length || !selectedId.value) return
  const firstId = messages.value[0].id
  const prevH = messagesEl.value?.scrollHeight ?? 0
  const res = await messageApi.listMessages(selectedId.value, firstId)
  messages.value = [...res.items, ...messages.value]
  hasMore.value = res.has_more
  await nextTick()
  if (messagesEl.value) messagesEl.value.scrollTop = messagesEl.value.scrollHeight - prevH
}

const scrollToBottom = () => {
  if (messagesEl.value) messagesEl.value.scrollTop = messagesEl.value.scrollHeight
}

const send = async () => {
  const text = content.value.trim()
  if (!text || sending.value || !selectedId.value) return
  sending.value = true
  content.value = ''
  try {
    await messageApi.send(selectedId.value, text)
    const res = await messageApi.listMessages(selectedId.value)
    messages.value = res.items
    await nextTick()
    scrollToBottom()
    await loadConversations()
  } catch (e: any) {
    toast.add({ title: e?.message || t('site.message.send_failed'), color: 'error' })
    content.value = text
  } finally {
    sending.value = false
  }
}

// ── Helpers ───────────────────────────────────────────────────────────────
const formatConvTime = (s: string) => {
  if (!s) return ''
  const d = new Date(s)
  const now = new Date()
  if (d.toDateString() === now.toDateString())
    return d.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  if (now.getTime() - d.getTime() < 7 * 86400_000)
    return d.toLocaleDateString('zh-CN', { weekday: 'short' })
  return d.toLocaleDateString('zh-CN', { month: '2-digit', day: '2-digit' })
}

const formatMsgTime = (s: string) =>
  new Date(s).toLocaleString('zh-CN', { month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' })

const showTimestamp = (idx: number) => {
  if (idx === messages.value.length - 1) return true
  return messages.value[idx + 1]?.sender_id !== messages.value[idx]?.sender_id
}

const showOtherAvatar = (idx: number) => {
  if (messages.value[idx]?.sender_id === authStore.user?.id) return false
  if (idx === messages.value.length - 1) return true
  return messages.value[idx + 1]?.sender_id !== messages.value[idx]?.sender_id
}

const bubbleRounding = (idx: number) => {
  const msg = messages.value[idx]
  const isMine = msg.sender_id === authStore.user?.id
  const prevSame = idx > 0 && messages.value[idx - 1]?.sender_id === msg.sender_id
  const nextSame = idx < messages.value.length - 1 && messages.value[idx + 1]?.sender_id === msg.sender_id
  const base = 'rounded-md'
  if (isMine) {
    if (prevSame && nextSame) return `${base} rounded-r-md`
    if (prevSame) return `${base} rounded-tr-md`
    if (nextSame) return `${base} rounded-br-md`
    return base
  } else {
    if (prevSame && nextSame) return `${base} rounded-l-md`
    if (prevSame) return `${base} rounded-tl-md`
    if (nextSame) return `${base} rounded-bl-md`
    return base
  }
}

// ── Mobile: back from chat to list ────────────────────────────────────────
const backToList = () => router.replace({ query: {} })

// ── Init ──────────────────────────────────────────────────────────────────
onMounted(async () => {
  await loadConversations()
  if (selectedId.value) {
    openConversation(selectedId.value)
  }
})

useHead({ title: t('site.message.inbox_title') })
</script>

<template>
  <!-- Full-height split layout below the site header (h-16 = 4rem) -->
  <div class="flex overflow-hidden border-t border-default" style="height: calc(100dvh - 4rem)">

    <!-- ═══ LEFT PANEL — Conversation list ═══ -->
    <div
      class="flex flex-col w-full md:w-72 lg:w-80 shrink-0 border-r border-default bg-default"
      :class="selectedId ? 'hidden md:flex' : 'flex'">

      <!-- Panel header -->
      <div class="flex items-center justify-between px-4 h-14 border-b border-default shrink-0">
        <h1 class="font-bold text-highlighted flex items-center gap-2 text-[15px]">
          {{ $t('site.message.inbox_title') }}
          <span v-if="totalUnread > 0"
            class="inline-flex items-center justify-center min-w-[18px] h-[18px] px-1 rounded-full bg-error text-white text-[10px] font-bold">
            {{ totalUnread > 99 ? '99+' : totalUnread }}
          </span>
        </h1>
      </div>

      <!-- Conversation list (scrollable) -->
      <div class="flex-1 overflow-y-auto">

        <!-- Skeleton -->
        <div v-if="convLoading" class="px-2 pt-2 space-y-1">
          <div v-for="i in 6" :key="i" class="flex gap-3 items-center rounded-md px-3 py-3">
            <USkeleton class="size-10 rounded-full shrink-0" />
            <div class="flex-1 space-y-1.5">
              <div class="flex justify-between">
                <USkeleton class="h-3.5 w-24" />
                <USkeleton class="h-3 w-8" />
              </div>
              <USkeleton class="h-3 w-36" />
            </div>
          </div>
        </div>

        <!-- Empty -->
        <div v-else-if="conversations.length === 0"
          class="flex flex-col items-center justify-center gap-2 h-48 text-center px-4">
          <UIcon name="i-tabler-message-2-off" class="size-10 text-muted" />
          <p class="text-sm text-muted">{{ $t('site.message.empty') }}</p>
        </div>

        <!-- Items -->
        <div v-else class="px-2 pt-2 pb-4 space-y-0.5">
          <button
            v-for="conv in conversations"
            :key="conv.id"
            class="w-full flex items-center gap-3 rounded-md px-3 py-2.5 transition-colors text-left group"
            :class="selectedId === conv.other_user_id
              ? 'bg-primary/10'
              : 'hover:bg-elevated'"
            @click="openConversation(conv.other_user_id)">

            <!-- Avatar + unread dot -->
            <div class="relative shrink-0">
              <BaseAvatar :src="conv.other_avatar" :alt="conv.other_name" size="sm" />
              <span v-if="conv.unread_count > 0 && selectedId !== conv.other_user_id"
                class="absolute -top-0.5 -right-0.5 size-2.5 rounded-full bg-error ring-2 ring-default" />
            </div>

            <!-- Text -->
            <div class="flex-1 min-w-0">
              <div class="flex items-center justify-between gap-1 mb-0.5">
                <span class="text-[13px] truncate"
                  :class="selectedId === conv.other_user_id
                    ? 'font-semibold text-primary'
                    : conv.unread_count > 0
                      ? 'font-bold text-highlighted'
                      : 'font-medium text-highlighted'">
                  {{ conv.other_name }}
                </span>
                <span class="text-[11px] shrink-0"
                  :class="selectedId === conv.other_user_id ? 'text-primary/70' : 'text-muted'">
                  {{ formatConvTime(conv.last_msg_at) }}
                </span>
              </div>
              <div class="flex items-center justify-between gap-1">
                <p class="text-xs truncate"
                  :class="conv.unread_count > 0 && selectedId !== conv.other_user_id
                    ? 'text-default font-medium'
                    : 'text-muted'">
                  {{ conv.last_msg || $t('site.message.no_msg') }}
                </p>
                <span v-if="conv.unread_count > 0 && selectedId !== conv.other_user_id"
                  class="shrink-0 min-w-[18px] h-[18px] px-1 rounded-full bg-error text-white text-[10px] font-bold flex items-center justify-center">
                  {{ conv.unread_count > 99 ? '99+' : conv.unread_count }}
                </span>
              </div>
            </div>

          </button>
        </div>

      </div>
    </div>

    <!-- ═══ RIGHT PANEL — Chat area ═══ -->
    <div class="flex-1 flex flex-col overflow-hidden"
      :class="selectedId ? 'flex' : 'hidden md:flex'">

      <!-- No conversation selected (desktop empty state) -->
      <div v-if="!selectedId" class="flex-1 flex flex-col items-center justify-center gap-4 text-center px-8">
        <div class="size-20 rounded-full bg-primary/10 flex items-center justify-center">
          <UIcon name="i-tabler-message-2" class="size-10 text-primary" />
        </div>
        <div>
          <p class="font-semibold text-highlighted mb-1">{{ $t('site.message.inbox_title') }}</p>
          <p class="text-sm text-muted">{{ $t('site.message.select_hint', 'Select a conversation to start chatting') }}</p>
        </div>
      </div>

      <!-- Active chat -->
      <template v-else>

        <!-- Chat header -->
        <div class="flex items-center gap-3 px-4 h-14 border-b border-default bg-default shrink-0">
          <!-- Mobile back button -->
          <UButton
            class="md:hidden"
            color="neutral" variant="ghost" icon="i-tabler-arrow-left" square size="sm"
            @click="backToList" />

          <!-- User info (link to profile) -->
          <NuxtLink :to="`/user/${selectedId}`"
            class="flex items-center gap-2.5 flex-1 min-w-0 group">
            <BaseAvatar
              :src="otherUser?.avatar"
              :alt="displayName"
              size="sm" class="shrink-0" />
            <div class="min-w-0">
              <p class="text-sm font-semibold text-highlighted truncate group-hover:text-primary transition-colors">
                {{ displayName }}
              </p>
              <p v-if="otherUser?.username" class="text-[11px] text-muted">
                @{{ otherUser.username }}
              </p>
            </div>
          </NuxtLink>

          <UButton
            :to="`/user/${selectedId}`"
            color="neutral" variant="ghost" icon="i-tabler-user-circle" square size="sm"
            class="shrink-0" />
        </div>

        <!-- Messages (scrollable) -->
        <div
          ref="messagesEl"
          class="flex-1 overflow-y-auto px-4 py-4 space-y-0.5 bg-muted/30"
          style="overscroll-behavior: contain">

          <!-- Load more -->
          <div class="flex justify-center pb-2">
            <UButton
              v-if="hasMore"
              size="xs" color="neutral" variant="soft"
              icon="i-tabler-chevrons-up"
              @click="loadMore">
              {{ $t('site.message.load_earlier') }}
            </UButton>
          </div>

          <!-- Skeleton -->
          <div v-if="msgLoading" class="space-y-3 py-2">
            <div v-for="i in 8" :key="i" class="flex gap-2 items-end"
              :class="i % 3 === 0 ? 'justify-end' : 'justify-start'">
              <USkeleton v-if="i % 3 !== 0" class="size-8 rounded-full shrink-0" />
              <USkeleton class="h-10 rounded-md" :class="i % 3 === 0 ? 'w-32' : 'w-44'" />
            </div>
          </div>

          <!-- Empty -->
          <div v-else-if="messages.length === 0"
            class="flex flex-col items-center justify-center gap-3 h-40 text-center">
            <div class="size-14 rounded-full bg-primary/10 flex items-center justify-center">
              <UIcon name="i-tabler-message-circle" class="size-7 text-primary" />
            </div>
            <p class="text-sm text-muted">{{ $t('site.message.first_msg') }}</p>
          </div>

          <!-- Message bubbles -->
          <template v-else>
            <div
              v-for="(msg, idx) in messages"
              :key="msg.id"
              class="flex items-end gap-2"
              :class="[
                msg.sender_id === authStore.user?.id ? 'justify-end' : 'justify-start',
                showTimestamp(idx) ? 'mb-2' : 'mb-0.5',
              ]">

              <!-- Other user avatar (aligned to last bubble in group) -->
              <div v-if="msg.sender_id !== authStore.user?.id" class="w-8 shrink-0">
                <BaseAvatar
                  v-if="showOtherAvatar(idx)"
                  :src="otherUser?.avatar"
                  :alt="displayName"
                  size="xs" />
              </div>

              <!-- Bubble -->
              <div class="max-w-[65%] flex flex-col"
                :class="msg.sender_id === authStore.user?.id ? 'items-end' : 'items-start'">
                <div
                  class="px-3.5 py-2 text-sm leading-relaxed break-words whitespace-pre-wrap"
                  :class="[
                    bubbleRounding(idx),
                    msg.sender_id === authStore.user?.id
                      ? 'bg-primary text-white'
                      : 'bg-default text-highlighted shadow-xs ring-1 ring-default',
                  ]">
                  {{ msg.content }}
                </div>

                <!-- Timestamp + read receipt -->
                <div v-if="showTimestamp(idx)"
                  class="flex items-center gap-1 mt-1 px-1"
                  :class="msg.sender_id === authStore.user?.id ? 'justify-end' : 'justify-start'">
                  <span class="text-[10px] text-muted">{{ formatMsgTime(msg.created_at) }}</span>
                  <UIcon
                    v-if="msg.sender_id === authStore.user?.id"
                    :name="msg.is_read ? 'i-tabler-checks' : 'i-tabler-check'"
                    class="size-3.5"
                    :class="msg.is_read ? 'text-primary' : 'text-muted'" />
                </div>
              </div>

            </div>
          </template>
        </div>

        <!-- Input area -->
        <div class="border-t border-default bg-default px-4 py-3 shrink-0">
          <div class="flex items-end gap-2">
            <UTextarea
              v-model="content"
              :placeholder="$t('site.message.send_placeholder')"
              :rows="1"
              autoresize
              class="flex-1"
              @keydown.enter.exact.prevent="send"
              @keydown.enter.shift.exact="() => {}" />
            <UButton
              color="primary" icon="i-tabler-send-2" square
              :loading="sending" :disabled="!content.trim()"
              @click="send" />
          </div>
          <p class="text-[10px] text-muted mt-1.5 select-none">
            {{ $t('site.message.send_hint') }}
          </p>
        </div>

      </template>
    </div>

  </div>
</template>

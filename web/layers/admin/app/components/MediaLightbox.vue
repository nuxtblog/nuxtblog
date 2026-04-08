<template>
  <Teleport to="body">
    <Transition name="gallery-fade">
      <div
        v-if="open && media"
        class="fixed inset-0 z-[200] flex flex-col bg-black select-none"
        @click.self="emit('update:open', false)">

        <!-- 顶部工具栏 -->
        <div class="flex items-center justify-between px-4 py-3 bg-black/80 shrink-0">
          <span class="text-white/70 text-sm truncate max-w-xs">{{ media.filename }}</span>
          <div class="flex items-center gap-2">
            <span class="text-white/50 text-sm">
              {{ index + 1 }} / {{ images.length }}
              <template v-if="totalPages > 1"> · P{{ currentPage }}/{{ totalPages }}</template>
            </span>
            <UButton
              :to="media.cdn_url" target="_blank"
              color="neutral" variant="ghost" icon="i-tabler-external-link" square size="sm"
              class="text-white/70 hover:text-white hover:bg-white/10" />
            <UButton
              color="neutral" variant="ghost" icon="i-tabler-x" square size="sm"
              class="text-white/70 hover:text-white hover:bg-white/10"
              @click="emit('update:open', false)" />
          </div>
        </div>

        <!-- 主图区 -->
        <div class="flex-1 relative flex items-center justify-center overflow-hidden">
          <img
            :key="media.id"
            :src="media.cdn_url"
            :alt="media.alt_text"
            class="max-w-full max-h-full object-contain" />

          <!-- 左箭头 -->
          <button
            v-if="hasPrev"
            :disabled="pageLoading"
            class="absolute left-3 top-1/2 -translate-y-1/2 size-10 rounded-full bg-black/50 hover:bg-black/80 flex items-center justify-center text-white transition-colors disabled:opacity-40"
            @click.stop="emit('prev')">
            <UIcon name="i-tabler-chevron-left" class="size-6" />
          </button>

          <!-- 右箭头 -->
          <button
            v-if="hasNext"
            :disabled="pageLoading"
            class="absolute right-3 top-1/2 -translate-y-1/2 size-10 rounded-full bg-black/50 hover:bg-black/80 flex items-center justify-center text-white transition-colors disabled:opacity-40"
            @click.stop="emit('next')">
            <UIcon
              :name="pageLoading ? 'i-tabler-loader-2' : 'i-tabler-chevron-right'"
              :class="['size-6', pageLoading && 'animate-spin']" />
          </button>
        </div>

        <!-- 底部缩略图条 -->
        <div class="shrink-0 bg-black/80 px-4 py-2">
          <div ref="thumbStripEl" class="flex gap-2 overflow-x-auto scrollbar-hide" style="scroll-behavior: smooth">
            <button
              v-for="(img, i) in images"
              :key="img.id"
              class="shrink-0 size-14 rounded overflow-hidden border-2 transition-all"
              :class="i === index ? 'border-white opacity-100' : 'border-transparent opacity-50 hover:opacity-80'"
              @click.stop="emit('select-index', i)">
              <img :src="getThumbUrl(img)" :alt="img.alt_text" class="w-full h-full object-cover" />
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import type { MediaResponse } from '~/types/api/media'

const props = defineProps<{
  open: boolean
  media: MediaResponse | null
  images: MediaResponse[]
  index: number
  hasPrev: boolean
  hasNext: boolean
  pageLoading: boolean
  currentPage: number
  totalPages: number
}>()

const emit = defineEmits<{
  'update:open': [boolean]
  'prev': []
  'next': []
  'select-index': [number]
}>()

const thumbStripEl = ref<HTMLElement | null>(null)

const scrollThumbIntoView = (idx: number) => {
  nextTick(() => {
    const strip = thumbStripEl.value
    if (!strip) return
    const thumb = strip.children[idx] as HTMLElement | undefined
    thumb?.scrollIntoView({ block: 'nearest', inline: 'center', behavior: 'smooth' })
  })
}

// Auto-scroll thumbnail strip when index changes
watch(() => props.index, (idx) => scrollThumbIntoView(idx))
watch(() => props.open, (open) => { if (open) nextTick(() => scrollThumbIntoView(props.index)) })

// Keyboard navigation
const onKey = (e: KeyboardEvent) => {
  if (!props.open) return
  if (e.key === 'ArrowLeft') emit('prev')
  else if (e.key === 'ArrowRight') emit('next')
  else if (e.key === 'Escape') emit('update:open', false)
}
onMounted(() => window.addEventListener('keydown', onKey))
onUnmounted(() => window.removeEventListener('keydown', onKey))

// Thumbnail URL helper (simple: prefer thumbnail variant)
const getThumbUrl = (media: MediaResponse): string => {
  if (!media.variants) return media.cdn_url
  try {
    const v = JSON.parse(media.variants) as Record<string, string>
    return v.thumbnail || v.cover || v.content || media.cdn_url
  } catch {
    return media.cdn_url
  }
}
</script>

<style scoped>
.gallery-fade-enter-active,
.gallery-fade-leave-active { transition: opacity 0.15s ease; }
.gallery-fade-enter-from,
.gallery-fade-leave-to { opacity: 0; }
.scrollbar-hide { -ms-overflow-style: none; scrollbar-width: none; }
.scrollbar-hide::-webkit-scrollbar { display: none; }
</style>

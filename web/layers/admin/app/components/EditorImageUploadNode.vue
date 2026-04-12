<template>
  <NodeViewWrapper class="my-2">
    <div
      class="border-2 border-dashed border-default rounded-lg p-6 flex flex-col items-center justify-center gap-2 cursor-pointer transition-colors hover:border-primary hover:bg-elevated/50"
      :class="{ 'opacity-60 pointer-events-none': compressing }"
      @click="triggerFileInput"
      @dragover.prevent="dragover = true"
      @dragleave.prevent="dragover = false"
      @drop.prevent="handleDrop">
      <template v-if="compressing">
        <UIcon name="i-tabler-loader-2" class="size-8 text-muted animate-spin" />
        <span class="text-sm text-muted">{{ t('admin.posts.editor.image_compressing') }}</span>
      </template>
      <template v-else>
        <UIcon name="i-tabler-photo-up" class="size-8 text-muted" />
        <span class="text-sm text-muted">{{ t('admin.posts.editor.image_upload_hint') }}</span>
      </template>
    </div>
    <input
      ref="fileInput"
      type="file"
      accept="image/*"
      class="hidden"
      @change="handleFileSelect" />
  </NodeViewWrapper>
</template>

<script setup lang="ts">
import { NodeViewWrapper } from "@tiptap/vue-3"
import { compressImage } from "../utils/compressImage"

const props = defineProps<{
  editor: any
  node: any
  getPos: () => number
}>()

const { t } = useI18n()
const toast = useToast()
const fileInput = ref<HTMLInputElement | null>(null)
const compressing = ref(false)
const dragover = ref(false)

function triggerFileInput() {
  fileInput.value?.click()
}

async function handleFile(file: File) {
  if (!file.type.startsWith("image/")) return
  compressing.value = true
  try {
    const dataUrl = await compressImage(file)
    const pos = props.getPos()
    props.editor
      .chain()
      .focus()
      .deleteRange({ from: pos, to: pos + props.node.nodeSize })
      .setImage({ src: dataUrl, alt: file.name.replace(/\.[^.]+$/, "") })
      .run()
  } catch {
    toast.add({ title: t("admin.posts.editor.image_upload_failed"), color: "error" })
    compressing.value = false
  }
}

function handleFileSelect(event: Event) {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (file) handleFile(file)
}

function handleDrop(event: DragEvent) {
  dragover.value = false
  const file = event.dataTransfer?.files?.[0]
  if (file) handleFile(file)
}
</script>

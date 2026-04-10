<template>
  <AdminPageContainer>
    <AdminPageHeader
      :title="mode === 'create' ? t('admin.docs.editor.create_title') : t('admin.docs.editor.edit_title')"
      :subtitle="mode === 'create' ? t('admin.docs.editor.create_subtitle') : t('admin.docs.editor.edit_subtitle')">
      <template #actions>
        <div class="flex items-center gap-3">
          <UButton
            v-if="mode === 'edit'"
            color="neutral"
            variant="ghost"
            icon="i-tabler-history"
            @click="showRevisionModal = true">
            {{ t('admin.docs.editor.revisions') }}
          </UButton>
          <UButton color="neutral" variant="soft" :disabled="submitting" @click="handleSaveDraft">
            {{ t('admin.docs.editor.save_draft') }}
          </UButton>
          <UButton color="primary" :loading="submitting" @click="handlePublish">
            {{ mode === 'create' ? t('admin.docs.editor.publish') : t('common.save') }}
          </UButton>
        </div>
      </template>
    </AdminPageHeader>

    <AdminPageContent class="p-0 flex flex-col md:flex-row">
      <!-- 主内容区 -->
      <div class="flex-1 min-w-0 max-w-5xl mx-auto">
        <div class="flex-1 overflow-y-auto bg-default">

          <!-- 草稿恢复提示 -->
          <div v-if="showDraftRestore" class="px-8 sm:px-16 pt-8 pb-3">
            <UAlert
              icon="i-tabler-device-floppy"
              color="warning"
              variant="subtle"
              :title="t('admin.posts.editor.draft_found')">
              <template #description>
                <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-2">
                  <span class="text-sm text-gray-500">
                    {{ savedDraft?.savedAt ? t('admin.posts.editor.draft_saved_at', { time: new Date(savedDraft.savedAt).toLocaleString() }) : '' }}
                  </span>
                  <div class="flex gap-2">
                    <UButton size="xs" color="primary" variant="soft" @click="restoreDraft">{{ t('admin.posts.editor.restore_draft') }}</UButton>
                    <UButton size="xs" color="neutral" variant="ghost" @click="discardDraft">{{ t('admin.posts.editor.discard') }}</UButton>
                  </div>
                </div>
              </template>
            </UAlert>
          </div>

          <!-- 标题 -->
          <div class="px-8 sm:px-16 pt-8 pb-3">
            <input
              v-model="formData.title"
              type="text"
              :placeholder="t('admin.docs.editor.title_placeholder')"
              class="w-full text-3xl font-bold bg-transparent border-none outline-none placeholder:text-muted" />
          </div>

          <!-- 别名 -->
          <div class="px-8 sm:px-16 pb-4">
            <input
              v-model="formData.slug"
              type="text"
              :placeholder="t('admin.docs.editor.slug_placeholder')"
              class="w-full text-sm bg-transparent border-b border-default pb-2 outline-none placeholder:text-muted focus:border-primary transition-colors" />
          </div>

          <!-- 编辑器骨架屏 -->
          <div v-if="editorLoading" class="px-8 sm:px-16 py-4">
            <div class="flex items-center gap-2 border-b border-default pb-2 mb-6">
              <USkeleton v-for="i in 8" :key="i" class="h-7 w-7 rounded" />
              <USkeleton class="h-7 w-px mx-1" />
              <USkeleton v-for="i in 5" :key="'b' + i" class="h-7 w-7 rounded" />
            </div>
            <div class="space-y-4 min-h-[500px]">
              <USkeleton class="h-8 w-2/3" />
              <USkeleton class="h-4 w-full" />
              <USkeleton class="h-4 w-5/6" />
              <USkeleton class="h-4 w-4/5" />
              <USkeleton class="h-4 w-full" />
              <USkeleton class="h-4 w-3/4" />
              <div class="pt-4 space-y-3">
                <USkeleton class="h-6 w-48" />
                <USkeleton class="h-4 w-full" />
                <USkeleton class="h-4 w-5/6" />
                <USkeleton class="h-4 w-full" />
              </div>
            </div>
          </div>

          <!-- 编辑器 -->
          <UEditor
            v-else
            ref="editorRef"
            v-slot="{ editor, handlers }"
            v-model="formData.content"
            content-type="markdown"
            :placeholder="t('admin.docs.editor.content_placeholder')"
            :extensions="editorExtensions"
            :handlers="editorHandlers"
            :ui="{ base: 'px-8 sm:px-16 py-6' }"
            class="min-h-[500px] px-8 sm:px-16 pb-4"
            @ready="editorLoading = false">
            <UEditorToolbar
              :editor="editor"
              :items="toolbarItems"
              layout="fixed"
              class="border-b border-default sticky top-0 inset-x-0 px-4 py-1.5 z-10 bg-default/95 backdrop-blur overflow-x-auto" />

            <UEditorToolbar
              :editor="editor"
              :items="bubbleItems"
              class="z-50"
              layout="bubble"
              :should-show="({ editor: e, view, state }) => {
                const { selection } = state;
                return view.hasFocus() && !selection.empty && !e.isActive('image');
              }" />

            <UEditorDragHandle
              v-slot="{ ui, onClick }"
              :editor="editor"
              @node-change="selectedNode = $event">
              <UButton
                icon="i-tabler-plus"
                color="neutral"
                variant="ghost"
                size="sm"
                :class="ui.handle()"
                @click="(e) => {
                  e.stopPropagation();
                  const selected = onClick();
                  handlers.suggestion?.execute(editor, { pos: selected?.pos }).run();
                }" />
              <UDropdownMenu
                v-slot="{ open }"
                :modal="false"
                :items="dragHandleItems(editor)"
                :content="{ side: 'left' }"
                :ui="{ content: 'w-52', label: 'text-xs' }"
                @update:open="editor.chain().setMeta('lockDragHandle', $event).run()">
                <UButton
                  color="neutral"
                  variant="ghost"
                  active-variant="soft"
                  size="sm"
                  icon="i-tabler-grip-vertical"
                  :active="open"
                  :class="ui.handle()" />
              </UDropdownMenu>
            </UEditorDragHandle>

            <UEditorSuggestionMenu :editor="editor" :items="suggestionItems" :append-to="appendToBody" />
            <UEditorMentionMenu   :editor="editor" :items="[]"             :append-to="appendToBody" />
            <UEditorEmojiMenu     :editor="editor" :items="emojiItems"     :append-to="appendToBody" />
          </UEditor>

          <!-- 字数统计 -->
          <div class="px-8 sm:px-16 py-2 flex items-center gap-4 text-xs text-muted border-t border-default">
            <span>{{ t('admin.posts.editor.char_count', { n: charCount }) }}</span>
            <span>{{ t('admin.posts.editor.reading_minutes', { n: readingMinutes }) }}</span>
            <span v-if="autoSavedLabel" class="ml-auto flex items-center gap-1">
              <UIcon name="i-tabler-circle-check" class="size-3 text-success" />{{ autoSavedLabel }}
            </span>
          </div>

        </div>
      </div>

      <!-- 右侧边栏 -->
      <div class="md:w-80 shrink-0 border-l border-default bg-default overflow-y-auto">
        <div class="p-4 space-y-6">

          <!-- 所属合集 -->
          <div class="space-y-3">
            <h4 class="text-sm font-semibold text-highlighted">{{ t('admin.docs.editor.collection_label') }}</h4>
            <AdminSearchableSelect
              v-model="formData.collection_id"
              :items="collectionOptions"
              :placeholder="t('admin.docs.editor.collection_label')"
              :search-placeholder="t('common.search')"
              trigger-class="w-full justify-between"
              @update:model-value="onCollectionChange" />
          </div>

          <!-- 父文档 -->
          <div class="space-y-3">
            <h4 class="text-sm font-semibold text-highlighted">{{ t('admin.docs.editor.parent_label') }}</h4>
            <AdminSearchableSelect
              v-model="formData.parent_id"
              :items="parentDocOptions"
              :placeholder="t('admin.docs.editor.no_parent')"
              :search-placeholder="t('common.search')"
              trigger-class="w-full justify-between"
              :disabled="!formData.collection_id" />
          </div>

          <!-- 发布设置 -->
          <div class="space-y-3">
            <h4 class="text-sm font-semibold text-highlighted">{{ t('admin.docs.editor.publish_settings') }}</h4>

            <UFormField :label="t('common.status')">
              <USelect v-model="formData.status" :items="statusOptions" class="w-full" />
            </UFormField>

            <UFormField :label="t('admin.docs.editor.published_at')">
              <UInput
                v-model="publishedAtLocal"
                type="datetime-local"
                class="w-full" />
            </UFormField>

            <UFormField :label="t('admin.docs.editor.comment_status')">
              <div class="flex items-center gap-2">
                <USwitch v-model="commentStatusBool" />
                <span class="text-sm text-muted">{{ formData.comment_status === 1 ? t('admin.docs.editor.comments_open') : t('admin.docs.editor.comments_closed') }}</span>
              </div>
            </UFormField>

            <UFormField :label="t('admin.docs.editor.locale')">
              <UInput v-model="formData.locale" placeholder="zh" class="w-full" />
            </UFormField>

            <UFormField :label="t('admin.docs.editor.sort_order')">
              <UInput v-model.number="formData.sort_order" type="number" class="w-full" />
            </UFormField>
          </div>

          <!-- 摘要 -->
          <div class="space-y-3">
            <h4 class="text-sm font-semibold text-highlighted">{{ t('admin.posts.editor.excerpt') }}</h4>
            <UTextarea
              v-model="formData.excerpt"
              :rows="3"
              :placeholder="t('admin.posts.editor.excerpt_placeholder')"
              class="w-full" />
          </div>

          <!-- SEO 设置 -->
          <details class="space-y-3">
            <summary class="text-sm font-semibold text-highlighted cursor-pointer select-none flex items-center gap-2 py-1">
              <UIcon name="i-tabler-search" class="size-4" />
              {{ t('admin.docs.editor.seo_settings') }}
            </summary>
            <div class="space-y-3 pt-2">
              <UFormField label="Meta Title">
                <UInput v-model="seoData.meta_title" placeholder="留空使用文档标题" class="w-full" />
              </UFormField>
              <UFormField label="Meta Description">
                <UTextarea v-model="seoData.meta_desc" placeholder="页面描述摘要" :rows="3" class="w-full" />
              </UFormField>
              <UFormField label="OG Title">
                <UInput v-model="seoData.og_title" placeholder="社交分享标题" class="w-full" />
              </UFormField>
              <UFormField label="OG Image">
                <UInput v-model="seoData.og_image" placeholder="社交分享图片 URL" class="w-full" />
              </UFormField>
              <UFormField label="Canonical URL">
                <UInput v-model="seoData.canonical_url" placeholder="规范链接（可选）" class="w-full" />
              </UFormField>
              <UFormField label="Robots">
                <UInput v-model="seoData.robots" placeholder="index,follow" class="w-full" />
              </UFormField>
            </div>
          </details>

        </div>
      </div>
    </AdminPageContent>

    <!-- 版本历史弹窗 -->
    <UModal v-if="mode === 'edit'" v-model:open="showRevisionModal" :ui="{ content: 'max-w-2xl' }">
      <template #content>
        <div class="p-6">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-lg font-semibold text-highlighted">{{ t('admin.docs.editor.revisions') }}</h3>
            <UButton icon="i-tabler-x" color="neutral" variant="ghost" size="sm" square @click="showRevisionModal = false" />
          </div>

          <div v-if="revisions.length === 0" class="flex flex-col items-center justify-center py-12">
            <UIcon name="i-tabler-history" class="size-12 text-muted mb-2" />
            <p class="text-sm text-muted">暂无修订历史</p>
          </div>

          <div v-else class="space-y-2 max-h-96 overflow-y-auto">
            <div
              v-for="rev in revisions" :key="rev.id"
              class="flex items-center gap-3 p-3 border border-default rounded-lg group hover:bg-elevated transition-colors">
              <div class="flex-1 min-w-0">
                <p class="text-sm font-medium text-highlighted truncate">{{ rev.title || '(无标题)' }}</p>
                <div class="flex items-center gap-2 mt-0.5">
                  <span class="text-xs text-muted">{{ new Date(rev.created_at).toLocaleString('zh-CN') }}</span>
                  <span v-if="rev.rev_note" class="text-xs text-muted">· {{ rev.rev_note }}</span>
                </div>
              </div>
              <UButton
                size="xs"
                color="neutral"
                variant="outline"
                icon="i-tabler-restore"
                class="opacity-0 group-hover:opacity-100 transition-opacity"
                :loading="restoringRevisionId === rev.id"
                @click="handleRestoreRevision(rev)">
                {{ t('admin.docs.restore_revision') }}
              </UButton>
            </div>
          </div>
        </div>
      </template>
    </UModal>

    <!-- 图片上传（隐藏） -->
    <input
      ref="imageFileInput"
      type="file"
      accept="image/*"
      class="hidden"
      @change="handleImageFileSelect" />

  </AdminPageContainer>
</template>

<script setup lang="ts">
import type { Editor } from '@tiptap/vue-3'
import { Emoji, gitHubEmojis } from '@tiptap/extension-emoji'
import { TextAlign } from '@tiptap/extension-text-align'
import { Markdown } from 'tiptap-markdown'
import type { CreateDocRequest, UpdateDocRequest, DocRevisionItem } from '~/types/api/doc'

const { t } = useI18n()

// ── Props & Emits ─────────────────────────────────────────────────────────
export interface DocEditorInitialData {
  id?: number
  collectionId?: number
  parentId?: number | null
  title?: string
  slug?: string
  content?: string
  excerpt?: string
  status?: number
  commentStatus?: number
  locale?: string
  sortOrder?: number
  publishedAt?: string
  seo?: {
    meta_title?: string; meta_desc?: string; og_title?: string
    og_image?: string; canonical_url?: string; robots?: string
  }
}

const props = defineProps<{
  mode: 'create' | 'edit'
  initialData?: DocEditorInitialData
  submitting?: boolean
}>()

const emit = defineEmits<{ save: [payload: CreateDocRequest | UpdateDocRequest] }>()

// ── Editor extensions & emoji ─────────────────────────────────────────────
const editorExtensions = [
  Emoji,
  TextAlign.configure({ types: ['heading', 'paragraph'] }),
  Markdown.configure({ transformPastedText: true, transformCopiedText: true }),
]
const appendToBody = import.meta.client ? () => document.body : undefined
const emojiItems   = gitHubEmojis.filter((e) => !e.name.startsWith('regional_indicator_'))

// ── Toolbar config ────────────────────────────────────────────────────────
const { toolbarItems, bubbleItems, suggestionItems, selectedNode, dragHandleItems } =
  usePostEditorToolbar()

// ── Image upload ──────────────────────────────────────────────────────────
const toast          = useToast()
const mediaStore     = useMediaStore()
const imageFileInput = ref<HTMLInputElement | null>(null)
let   pendingEditor: Editor | null = null

const compressImage = (file: File, maxWidth = 1920, quality = 0.85): Promise<string> =>
  new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = (e) => {
      const img = new Image()
      img.onload = () => {
        const canvas = document.createElement('canvas')
        const scale  = Math.min(1, maxWidth / img.width)
        canvas.width  = Math.round(img.width  * scale)
        canvas.height = Math.round(img.height * scale)
        canvas.getContext('2d')!.drawImage(img, 0, 0, canvas.width, canvas.height)
        resolve(canvas.toDataURL('image/jpeg', quality))
      }
      img.onerror = reject
      img.src = e.target?.result as string
    }
    reader.onerror = reject
    reader.readAsDataURL(file)
  })

const editorHandlers = computed(() => ({
  image: {
    canExecute: (editor: Editor) => editor.isEditable,
    execute: (editor: Editor) => {
      pendingEditor = editor
      nextTick(() => imageFileInput.value?.click())
      return editor.chain()
    },
    isActive: (_editor: Editor) => false,
  },
}))

const handleImageFileSelect = async (event: Event) => {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (!file || !pendingEditor) return
  const editor = pendingEditor
  pendingEditor = null
  try {
    const dataUrl = await compressImage(file)
    editor.chain().focus().setImage({ src: dataUrl, alt: file.name.replace(/\.[^.]+$/, '') }).run()
  } catch {
    toast.add({ title: t('admin.posts.editor.image_upload_failed'), color: 'error' })
  } finally {
    if (imageFileInput.value) imageFileInput.value.value = ''
  }
}

const dataUrlToFile = (dataUrl: string, name: string): File => {
  const [header = "", b64 = ""] = dataUrl.split(",")
  const mime = header.match(/:(.*?);/)?.[1] ?? 'image/jpeg'
  const ext  = mime.split('/')[1] ?? 'jpg'
  const bin  = atob(b64)
  const arr  = new Uint8Array(bin.length)
  for (let i = 0; i < bin.length; i++) arr[i] = bin.charCodeAt(i)
  return new File([arr], name.includes('.') ? name : `${name}.${ext}`, { type: mime })
}

const uploadPendingImages = async (): Promise<void> => {
  const content = formData.value.content ?? ''
  const regex   = /!\[([^\]]*)\]\((data:[^)]+)\)/g
  const matches: { alt: string; src: string }[] = []
  let m: RegExpExecArray | null
  while ((m = regex.exec(content)) !== null) matches.push({ alt: m[1] ?? "", src: m[2] ?? "" })
  if (!matches.length) return

  toast.add({ title: t('admin.posts.editor.image_uploading'), color: 'neutral', duration: 0, id: 'img-upload' })
  try {
    const results = await Promise.all(
      matches.map(async ({ alt, src }) => {
        const name   = alt || `image-${Date.now()}`
        const result = await mediaStore.uploadMedia(dataUrlToFile(src, name), { title: name, category: 'doc' })
        return { src, cdnUrl: result?.cdn_url }
      }),
    )
    let updated = content
    for (const { src, cdnUrl } of results) {
      if (cdnUrl) updated = updated.replaceAll(src, cdnUrl)
    }
    formData.value.content = updated
    await nextTick()
  } finally {
    toast.remove('img-upload')
  }
}

// ── Form state ────────────────────────────────────────────────────────────
const docApi = useDocApi()
const { apiFetch } = useApiFetch()

const init = props.initialData

const DEFAULT_PUBLISHED_AT = () => new Date().toISOString().slice(0, 16)

const editorLoading = ref(true)

const publishedAtLocal = ref(
  init?.publishedAt ? new Date(init.publishedAt).toISOString().slice(0, 16) : DEFAULT_PUBLISHED_AT(),
)

const formData = ref<{
  collection_id: number | undefined
  parent_id: number | null | undefined
  title: string
  slug: string
  content: string
  excerpt: string
  status: number
  comment_status: number
  locale: string
  sort_order: number
}>({
  collection_id: init?.collectionId,
  parent_id:     init?.parentId ?? null,
  title:         init?.title    ?? '',
  slug:          init?.slug     ?? '',
  content:       init?.content  ?? '',
  excerpt:       init?.excerpt  ?? '',
  status:        init?.status   ?? 1,
  comment_status: init?.commentStatus ?? 1,
  locale:        init?.locale   ?? 'zh',
  sort_order:    init?.sortOrder ?? 0,
})

const seoData = ref({
  meta_title:    init?.seo?.meta_title    ?? '',
  meta_desc:     init?.seo?.meta_desc     ?? '',
  og_title:      init?.seo?.og_title      ?? '',
  og_image:      init?.seo?.og_image      ?? '',
  canonical_url: init?.seo?.canonical_url ?? '',
  robots:        init?.seo?.robots        ?? 'index,follow',
})

const commentStatusBool = computed({
  get: () => formData.value.comment_status === 1,
  set: (val: boolean) => { formData.value.comment_status = val ? 1 : 0 },
})

// ── Status/options ────────────────────────────────────────────────────────
const statusOptions = computed(() => [
  { label: t('admin.docs.editor.status_draft'),     value: 1 },
  { label: t('admin.docs.editor.status_published'), value: 2 },
  { label: t('admin.docs.editor.status_archived'),  value: 3 },
])

// ── Collections ───────────────────────────────────────────────────────────
const collections = ref<{ id: number; title: string; slug: string }[]>([])
const parentDocs  = ref<{ id: number; title: string }[]>([])

const collectionOptions = computed(() => [
  { label: t('admin.docs.all_collections'), value: undefined },
  ...collections.value.map(c => ({ label: c.title, value: c.id })),
])

const parentDocOptions = computed(() => [
  { label: t('admin.docs.editor.no_parent'), value: null },
  ...parentDocs.value
    .filter(d => d.id !== init?.id)
    .map(d => ({ label: d.title, value: d.id })),
])

async function fetchCollections() {
  try {
    const res = await docApi.getCollections({ page_size: 100, status: 2 })
    collections.value = res.data ?? []
  } catch {}
}

async function fetchParentDocs(collectionId: number) {
  try {
    const res = await docApi.getDocs({ collection_id: collectionId, page_size: 100 })
    parentDocs.value = res.data ?? []
  } catch {}
}

function onCollectionChange(val: number | undefined) {
  formData.value.parent_id = null
  if (val) fetchParentDocs(val)
  else parentDocs.value = []
}

// ── Auto slug from title ──────────────────────────────────────────────────
watch(() => formData.value.title, (val) => {
  if (!formData.value.slug) {
    formData.value.slug = val
      .toLowerCase()
      .trim()
      .replace(/[^\w\s-]/g, '')
      .replace(/\s+/g, '-')
      .replace(/-+/g, '-')
  }
})

// ── Unsaved changes ───────────────────────────────────────────────────────
const hasUnsavedChanges = ref(false)
watch([formData, seoData], () => { hasUnsavedChanges.value = true }, { deep: true })

const markSaved = () => {
  hasUnsavedChanges.value = false
  try { localStorage.removeItem(autoSaveKey.value) } catch {}
}

// ── Auto-save ─────────────────────────────────────────────────────────────
const autoSaveKey = computed(() =>
  props.mode === 'edit' && init?.id ? `blog:doc:draft:${init.id}` : 'blog:doc:draft:new',
)
const lastAutoSaved  = ref<Date | null>(null)
const autoSavedLabel = computed(() => {
  if (!lastAutoSaved.value) return ''
  return t('admin.posts.editor.auto_saved', {
    time: lastAutoSaved.value.toLocaleTimeString(undefined, { hour: '2-digit', minute: '2-digit' }),
  })
})
const showDraftRestore = ref(false)
const savedDraft = ref<{ title?: string; content?: string; excerpt?: string; savedAt?: string } | null>(null)

const doAutoSave = () => {
  if (!formData.value.title && !formData.value.content) return
  try {
    localStorage.setItem(autoSaveKey.value, JSON.stringify({
      title: formData.value.title, content: formData.value.content,
      excerpt: formData.value.excerpt, savedAt: new Date().toISOString(),
    }))
    lastAutoSaved.value = new Date()
  } catch {}
}

const restoreDraft = () => {
  if (!savedDraft.value) return
  if (savedDraft.value.title)   formData.value.title   = savedDraft.value.title
  if (savedDraft.value.content) formData.value.content = savedDraft.value.content
  if (savedDraft.value.excerpt) formData.value.excerpt = savedDraft.value.excerpt
  showDraftRestore.value = false
  localStorage.removeItem(autoSaveKey.value)
}

const discardDraft = () => {
  localStorage.removeItem(autoSaveKey.value)
  showDraftRestore.value = false
}

// ── Revision modal ────────────────────────────────────────────────────────
const showRevisionModal    = ref(false)
const revisions            = ref<DocRevisionItem[]>([])
const restoringRevisionId  = ref<number | null>(null)

async function loadRevisions() {
  if (!init?.id) return
  try {
    const res = await docApi.getRevisions(init.id)
    revisions.value = res.list ?? []
  } catch {}
}

async function handleRestoreRevision(rev: DocRevisionItem) {
  if (!init?.id) return
  restoringRevisionId.value = rev.id
  try {
    await docApi.restoreRevision(init.id, rev.id)
    toast.add({ title: t('admin.docs.restore_revision'), color: 'success' })
    showRevisionModal.value = false
    // Reload page to reflect restored content
    window.location.reload()
  } catch (e: any) {
    toast.add({ title: e?.message ?? t('common.operation_failed'), color: 'error' })
  } finally {
    restoringRevisionId.value = null
  }
}

// ── Word count ────────────────────────────────────────────────────────────
const charCount = computed(() =>
  (formData.value.content ?? '').replace(/[#*`\[\]()>_~\-|]/g, '').replace(/\s+/g, '').trim().length,
)
const readingMinutes = computed(() => Math.max(1, Math.ceil(charCount.value / 400)))

// ── Build & save ──────────────────────────────────────────────────────────
const isFormValid = computed(
  () => formData.value.title.trim() !== '' && !!formData.value.collection_id,
)

const buildPayload = (status: number): CreateDocRequest | UpdateDocRequest => {
  const base = {
    collection_id: formData.value.collection_id as number,
    parent_id:     formData.value.parent_id ?? undefined,
    title:         formData.value.title,
    slug:          formData.value.slug || formData.value.title,
    content:       formData.value.content,
    excerpt:       formData.value.excerpt || undefined,
    status:        status as 1 | 2 | 3,
    comment_status: formData.value.comment_status,
    locale:        formData.value.locale,
    sort_order:    formData.value.sort_order,
    published_at:  publishedAtLocal.value ? new Date(publishedAtLocal.value).toISOString() : undefined,
  }
  return base
}

const handleSaveDraft = async () => {
  if (!isFormValid.value) {
    toast.add({ title: t('admin.posts.editor.fill_required'), color: 'warning' })
    return
  }
  await uploadPendingImages()
  emit('save', buildPayload(1))
}

const handlePublish = async () => {
  if (!isFormValid.value) {
    toast.add({ title: t('admin.posts.editor.fill_required'), color: 'warning' })
    return
  }
  await uploadPendingImages()
  const status = props.mode === 'create' ? 2 : formData.value.status
  emit('save', buildPayload(status))
}

defineExpose({ isDirty: hasUnsavedChanges, getIsDirty: () => hasUnsavedChanges.value, markSaved, seoData })

// ── Lifecycle ─────────────────────────────────────────────────────────────
let autoSaveTimer: ReturnType<typeof setInterval>

onMounted(async () => {
  editorLoading.value = false
  await fetchCollections()

  if (init?.collectionId) {
    await fetchParentDocs(init.collectionId)
  }

  if (props.mode === 'edit' && init?.id) {
    loadRevisions()
  }

  if (!init?.content) {
    try {
      const saved = localStorage.getItem(autoSaveKey.value)
      if (saved) {
        const draft = JSON.parse(saved)
        if (draft.title || draft.content) { savedDraft.value = draft; showDraftRestore.value = true }
      }
    } catch {}
  }

  autoSaveTimer = setInterval(doAutoSave, 60_000)
})

onUnmounted(() => clearInterval(autoSaveTimer))
</script>

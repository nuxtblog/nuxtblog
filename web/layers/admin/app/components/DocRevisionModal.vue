<script setup lang="ts">
import type { DocRevisionItem } from "~/types/api/doc";

const props = defineProps<{
  open: boolean
  docId: number
}>()

const emit = defineEmits<{
  (e: 'update:open', val: boolean): void
  (e: 'restored'): void
}>()

const { t } = useI18n()
const toast = useToast()
const docApi = useDocApi()

const revisions = ref<DocRevisionItem[]>([])
const restoringRevisionId = ref<number | null>(null)

watch(() => props.open, async (val) => {
  if (!val) return
  try {
    const res = await docApi.getRevisions(props.docId)
    revisions.value = res.list ?? []
  } catch {}
})

const handleRestore = async (rev: DocRevisionItem) => {
  restoringRevisionId.value = rev.id
  try {
    await docApi.restoreRevision(props.docId, rev.id)
    toast.add({ title: t("admin.docs.restore_revision"), color: "success" })
    emit('update:open', false)
    emit('restored')
  } catch (e: any) {
    toast.add({ title: e?.message ?? t("common.operation_failed"), color: "error" })
  } finally {
    restoringRevisionId.value = null
  }
}
</script>

<template>
  <UModal :open="open" :ui="{ content: 'max-w-2xl' }" @update:open="emit('update:open', $event)">
    <template #content>
      <div class="p-6">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-lg font-semibold text-highlighted">{{ t("admin.docs.editor.revisions") }}</h3>
          <UButton icon="i-tabler-x" color="neutral" variant="ghost" size="sm" square @click="emit('update:open', false)" />
        </div>

        <div v-if="revisions.length === 0" class="flex flex-col items-center justify-center py-12">
          <UIcon name="i-tabler-history" class="size-12 text-muted mb-2" />
          <p class="text-sm text-muted">暂无修订历史</p>
        </div>

        <div v-else class="space-y-2 max-h-96 overflow-y-auto">
          <div
            v-for="rev in revisions"
            :key="rev.id"
            class="flex items-center gap-3 p-3 border border-default rounded-lg group hover:bg-elevated transition-colors">
            <div class="flex-1 min-w-0">
              <p class="text-sm font-medium text-highlighted truncate">{{ rev.title || "(无标题)" }}</p>
              <div class="flex items-center gap-2 mt-0.5">
                <span class="text-xs text-muted">{{ new Date(rev.created_at).toLocaleString("zh-CN") }}</span>
                <span v-if="rev.rev_note" class="text-xs text-muted">· {{ rev.rev_note }}</span>
              </div>
            </div>
            <UButton
              size="xs" color="neutral" variant="outline" icon="i-tabler-restore"
              class="opacity-0 group-hover:opacity-100 transition-opacity"
              :loading="restoringRevisionId === rev.id"
              @click="handleRestore(rev)">
              {{ t("admin.docs.restore_revision") }}
            </UButton>
          </div>
        </div>
      </div>
    </template>
  </UModal>
</template>

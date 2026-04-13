<script setup lang="ts">
import type { Editor } from '@tiptap/vue-3'

const props = defineProps<{
  editor: Editor
}>()

const open = ref(false)
const linkText = ref('')
const linkUrl = ref('')
const isEditing = ref(false)

function openForInsert() {
  const { from, to } = props.editor.state.selection
  linkText.value = from !== to ? props.editor.state.doc.textBetween(from, to, ' ') : ''
  linkUrl.value = ''
  isEditing.value = false
  open.value = true
  nextTick(() => {
    (document.querySelector('.link-popover-url') as HTMLInputElement)?.focus()
  })
}

function openForEdit() {
  const attrs = props.editor.getAttributes('link')
  linkUrl.value = attrs.href || ''
  // Get the linked text
  const { from, to } = props.editor.state.selection
  if (from === to) {
    // Cursor is inside link, expand selection to full mark
    props.editor.chain().focus().extendMarkRange('link').run()
  }
  const sel = props.editor.state.selection
  linkText.value = props.editor.state.doc.textBetween(sel.from, sel.to, ' ')
  isEditing.value = true
  open.value = true
  nextTick(() => {
    (document.querySelector('.link-popover-url') as HTMLInputElement)?.focus()
  })
}

function apply() {
  const url = linkUrl.value.trim()
  if (!url) {
    open.value = false
    return
  }
  const chain = props.editor.chain().focus()

  if (isEditing.value) {
    // Update existing link
    chain.extendMarkRange('link')
    const { from, to } = props.editor.state.selection
    const currentText = props.editor.state.doc.textBetween(from, to, ' ')
    if (linkText.value.trim() && linkText.value !== currentText) {
      chain.deleteRange({ from, to }).insertContent({
        type: 'text',
        text: linkText.value,
        marks: [{ type: 'link', attrs: { href: url } }],
      })
    } else {
      chain.setLink({ href: url })
    }
  } else {
    // Insert new link
    const { from, to } = props.editor.state.selection
    if (from === to) {
      // No selection — insert text with link
      chain.insertContent({
        type: 'text',
        text: linkText.value || url,
        marks: [{ type: 'link', attrs: { href: url } }],
      })
    } else {
      chain.setLink({ href: url })
    }
  }

  chain.run()
  open.value = false
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Enter') {
    e.preventDefault()
    apply()
  }
  if (e.key === 'Escape') {
    open.value = false
  }
}

defineExpose({ openForInsert, openForEdit })
</script>

<template>
  <UModal v-model:open="open" :ui="{ content: 'sm:max-w-md' }">
    <template #content>
      <div class="p-4 space-y-4">
        <h3 class="text-sm font-semibold text-highlighted">
          {{ isEditing ? $t('admin.posts.editor.edit_link') : $t('admin.posts.editor.insert_link') }}
        </h3>
        <UFormField :label="$t('admin.posts.editor.link_text')">
          <UInput
            v-model="linkText"
            :placeholder="$t('admin.posts.editor.link_text_placeholder')"
            class="w-full"
            @keydown="onKeydown" />
        </UFormField>
        <UFormField label="URL">
          <UInput
            v-model="linkUrl"
            class="link-popover-url w-full"
            placeholder="https://"
            @keydown="onKeydown" />
        </UFormField>
        <div class="flex justify-end gap-2">
          <UButton
            color="neutral"
            variant="ghost"
            size="sm"
            @click="open = false">
            {{ $t('common.cancel') }}
          </UButton>
          <UButton
            color="primary"
            size="sm"
            :disabled="!linkUrl.trim()"
            @click="apply">
            {{ isEditing ? $t('common.save') : $t('admin.posts.editor.insert') }}
          </UButton>
        </div>
      </div>
    </template>
  </UModal>
</template>

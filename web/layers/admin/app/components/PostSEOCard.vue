<script setup lang="ts">
interface SeoData {
  meta_title: string
  meta_desc: string
  og_title: string
  og_image: string
  canonical_url: string
  robots: string
}

defineProps<{
  title: string
  collapsed: boolean
}>()

const emit = defineEmits<{
  (e: 'toggle'): void
}>()

const seoData = defineModel<SeoData>('seo', { required: true })

const { t } = useI18n()
</script>

<template>
  <SidebarCard :title="title" :collapsed="collapsed" @toggle="emit('toggle')">
    <div class="space-y-3">
      <UFormField :label="t('admin.posts.editor.seo_title')">
        <UInput v-model="seoData.meta_title" :placeholder="t('admin.posts.editor.seo_title_placeholder')" class="w-full" />
      </UFormField>
      <UFormField :label="t('admin.posts.editor.seo_desc')">
        <UTextarea v-model="seoData.meta_desc" :rows="3" :placeholder="t('admin.posts.editor.seo_desc_placeholder')" class="w-full" />
      </UFormField>
      <UFormField :label="t('admin.posts.editor.og_title')">
        <UInput v-model="seoData.og_title" :placeholder="t('admin.posts.editor.og_title_placeholder')" class="w-full" />
      </UFormField>
      <UFormField :label="t('admin.posts.editor.og_image')">
        <UInput v-model="seoData.og_image" placeholder="https://..." class="w-full" />
      </UFormField>
      <UFormField label="Canonical URL">
        <UInput v-model="seoData.canonical_url" placeholder="https://..." class="w-full" />
      </UFormField>
      <UFormField label="Robots">
        <USelect
          v-model="seoData.robots"
          :items="[
            { label: 'index, follow', value: 'index,follow' },
            { label: 'noindex, follow', value: 'noindex,follow' },
            { label: 'index, nofollow', value: 'index,nofollow' },
            { label: 'noindex, nofollow', value: 'noindex,nofollow' },
          ]"
          class="w-full" />
      </UFormField>
    </div>
  </SidebarCard>
</template>

<script setup lang="ts">
interface SeoData {
  meta_title: string
  meta_desc: string
  og_title: string
  og_image: string
  canonical_url: string
  robots: string
}

const seoData = defineModel<SeoData>({ required: true })

const { t } = useI18n()
</script>

<template>
  <details class="space-y-3">
    <summary class="text-sm font-semibold text-highlighted cursor-pointer select-none flex items-center gap-2 py-1">
      <UIcon name="i-tabler-search" class="size-4" />
      {{ t("admin.docs.editor.seo_settings") }}
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
</template>

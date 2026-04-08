<template>
  <div class="min-h-screen flex flex-col">
    <BlogHeader />
    <main class="flex-1 flex items-center justify-center px-4 py-16">
      <div class="text-center max-w-md mx-auto">
        <!-- 错误码大字 -->
        <div class="text-8xl font-black text-primary/20 leading-none mb-6 select-none">
          {{ error?.statusCode ?? '?' }}
        </div>

        <!-- 图标 -->
        <div class="size-20 rounded-full bg-primary/10 flex items-center justify-center mx-auto mb-6">
          <UIcon
            :name="is404 ? 'i-tabler-map-search' : 'i-tabler-alert-triangle'"
            class="size-10 text-primary" />
        </div>

        <!-- 文字 -->
        <h1 class="text-2xl font-bold text-highlighted mb-2">
          {{ is404 ? $t('site.error.not_found_title') : $t('site.error.server_error_title') }}
        </h1>
        <p class="text-muted mb-8 leading-relaxed">
          {{ is404 ? $t('site.error.not_found_desc') : $t('site.error.server_error_desc') }}
        </p>

        <!-- 操作按钮 -->
        <div class="flex items-center justify-center gap-3">
          <UButton
            color="neutral"
            variant="outline"
            icon="i-tabler-arrow-left"
            @click="handleGoBack">
            {{ $t('site.error.go_back') }}
          </UButton>
          <UButton
            color="primary"
            icon="i-tabler-home"
            to="/">
            {{ $t('site.error.back_home') }}
          </UButton>
        </div>
      </div>
    </main>
    <BlogFooter />
  </div>
</template>

<script setup lang="ts">
import type { NuxtError } from '#app'

const props = defineProps<{ error: NuxtError | null }>()

const is404 = computed(() => props.error?.statusCode === 404)

const handleGoBack = () => {
  if (window.history.length > 1) {
    window.history.back()
  } else {
    clearError({ redirect: '/' })
  }
}
</script>

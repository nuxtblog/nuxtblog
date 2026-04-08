<template>
  <section class="w-full bg-gradient-to-r from-primary/10 via-secondary/10 to-transparent relative">
    <div>
      <div class="flex items-center justify-between py-6 px-6">
        <!-- 左侧：图标和文字 -->
        <div class="flex items-center space-x-4">
          <!-- 图标/徽章 -->
          <div
            class="inline-flex items-center justify-center size-16 rounded-md bg-gradient-to-br from-primary to-secondary text-white shadow-xl hover:scale-105 transition-transform duration-300"
          >
            <slot name="icon">
              <UIcon name="i-tabler-users" class="size-8" />
            </slot>
          </div>

          <!-- 标题和副标题 -->
          <div class="text-left">
            <h2 class="text-2xl font-bold text-highlighted mb-1">{{ title }}</h2>
            <p v-if="subtitle" class="text-muted text-sm">
              {{ subtitle }}
            </p>
          </div>
        </div>

        <!-- 右侧：操作按钮 -->
        <div class="flex items-center space-x-3">
          <slot name="actions" />
        </div>
      </div>

      <!-- 装饰线 -->
      <div class="h-1 bg-gradient-to-r w-full from-primary via-secondary to-transparent rounded-full" />
    </div>
  </section>
</template>

<script setup lang="ts">
export interface PageActionHeaderProps {
  title: string;
  subtitle?: string;
}
withDefaults(defineProps<PageActionHeaderProps>(), {
  title: "标题",
  subtitle: "副标题",
});
</script>

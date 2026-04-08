<script setup lang="ts">
const visible = ref(false);
const progress = ref(0);

const SHOW_THRESHOLD = 300;
const CIRCLE_RADIUS = 18;

const { t } = useI18n();

useEventListener(
  "scroll",
  useThrottleFn(() => {
    const scrolled = window.scrollY;
    const total = document.documentElement.scrollHeight - window.innerHeight;
    visible.value = scrolled > SHOW_THRESHOLD;
    progress.value = total > 0 ? Math.min(100, (scrolled / total) * 100) : 0;
  }),
);

const scrollToTop = () => window.scrollTo({ top: 0, behavior: "smooth" });

// SVG 圆环处理
const circumference = 2 * Math.PI * CIRCLE_RADIUS;
const dashOffset = computed(
  () => circumference - (progress.value / 100) * circumference,
);
</script>

<template>
  <Transition
    enter-active-class="transition duration-300 ease-out"
    enter-from-class="opacity-0 translate-y-4 scale-90"
    enter-to-class="opacity-100 translate-y-0 scale-100"
    leave-active-class="transition duration-200 ease-in"
    leave-from-class="opacity-100 translate-y-0 scale-100"
    leave-to-class="opacity-0 translate-y-4 scale-90">
    <button
      v-if="visible"
      class="fixed bottom-8 right-6 md:right-8 z-50 group size-11 flex items-center justify-center rounded-full bg-primary shadow-lg shadow-primary/30 hover:shadow-primary/50 hover:scale-110 active:scale-95 transition-all duration-200 cursor-pointer"
      :aria-label="t('common.back_to_top')"
      @click="scrollToTop">
      <!-- svg 圆环 -->
      <svg class="absolute inset-0 size-full -rotate-90" viewBox="0 0 44 44">
        <circle
          cx="22"
          cy="22"
          :r="CIRCLE_RADIUS"
          fill="none"
          stroke="white"
          stroke-opacity="0.2"
          stroke-width="2" />
        <circle
          cx="22"
          cy="22"
          :r="CIRCLE_RADIUS"
          fill="none"
          stroke="white"
          stroke-width="2"
          stroke-linecap="round"
          :stroke-dasharray="circumference"
          :stroke-dashoffset="dashOffset"
          class="transition-all duration-150" />
      </svg>
      <UIcon
        name="i-tabler-arrow-up"
        class="size-4 text-white relative z-10 group-hover:-translate-y-0.5 transition-transform duration-200" />
    </button>
  </Transition>
</template>

<template>
  <NuxtImg :src="currentSrc" v-bind="$attrs" />
</template>

<script setup lang="ts">
defineOptions({ inheritAttrs: false });

const props = defineProps<{ src: string }>();

const STATIC_FALLBACK = "/images/default-cover.svg";
const store = useOptionsStore();
const defaultCover = computed(() => store.get("default_post_cover", STATIC_FALLBACK));
const errorCover = computed(() => store.get("error_post_cover", STATIC_FALLBACK));

const currentSrc = ref(defaultCover.value || STATIC_FALLBACK);

const load = (src: string) => {
  const placeholder = defaultCover.value || STATIC_FALLBACK;

  // 如果传入的就是默认图，直接显示，不需要预加载
  if (!src || src === placeholder) {
    currentSrc.value = src || placeholder;
    return;
  }

  currentSrc.value = placeholder;

  const img = new Image();
  img.onload = () => { currentSrc.value = src; };
  img.onerror = () => { currentSrc.value = errorCover.value || STATIC_FALLBACK; };
  img.src = src;
};

onMounted(() => load(props.src));

watch(() => props.src, load);
</script>

<script setup lang="ts">
defineOptions({ inheritAttrs: false })

const props = defineProps<{
  src?: string | null
  alt?: string
}>()

const attrs = useAttrs()
const optionsStore = useOptionsStore()

const defaultAvatarUrl  = computed(() => optionsStore.get('default_avatar_url', '/images/default-avatar.svg'))
const defaultAvatarType = computed(() => optionsStore.get('default_avatar_type', 'initials'))

const effectiveSrc = computed(() => {
  if (props.src) return props.src
  if (defaultAvatarType.value === 'image' && defaultAvatarUrl.value) return defaultAvatarUrl.value
  return undefined
})

const fallbackText = computed(() => {
  const name = props.alt ?? ''
  return name.trim()[0]?.toUpperCase() || '?'
})
</script>

<template>
  <UAvatar
    v-bind="attrs"
    :src="effectiveSrc"
    :alt="alt"
    :text="fallbackText"
  />
</template>

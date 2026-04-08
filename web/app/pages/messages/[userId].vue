<script setup lang="ts">
// Redirect to the unified messages page with the conversation pre-selected
const route = useRoute()
const router = useRouter()

onMounted(() => {
  router.replace({ path: '/messages', query: { with: route.params.userId } })
})
</script>

<template>
  <div />
</template>

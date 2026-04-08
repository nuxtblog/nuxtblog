<script setup lang="ts">
definePageMeta({ middleware: 'auth' })

const authStore = useAuthStore()
await navigateTo(`/user/${authStore.user?.id}`, { replace: true })
</script>

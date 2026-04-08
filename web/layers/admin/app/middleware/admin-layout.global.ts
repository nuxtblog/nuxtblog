export default defineNuxtRouteMiddleware((to) => {
  if (!to.path.startsWith('/admin')) return

  const authStore = useAuthStore()

  if (!authStore.isLoggedIn) {
    return navigateTo(`/auth/login?redirect=${encodeURIComponent(to.fullPath)}`)
  }

  const { can } = usePermissions()
  if (!can('access_admin')) {
    return navigateTo('/')
  }

  // Only set admin layout after all checks pass — prevents layout leaking to
  // redirect destinations when the user lacks permission.
  if (to.meta.layout === undefined) {
    setPageLayout('admin')
  }
})

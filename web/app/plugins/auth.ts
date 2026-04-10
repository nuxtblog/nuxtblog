import type { LocaleCode } from '~/types/locale'

export default defineNuxtPlugin(async () => {
  const authStore = useAuthStore()
  const { $i18n } = useNuxtApp()
  if ((authStore.token || useCookie('blog_refresh_token').value) && !authStore.user) {
    await authStore.fetchMe()
    await authStore.syncLocale()
    if (authStore.user?.locale) {
      await $i18n.setLocale(authStore.user.locale as LocaleCode)
    }
  }
})

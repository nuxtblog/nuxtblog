import type { LocaleCode } from '~/types/locale'

export const useLocaleSwitch = () => {
  const { locale, setLocale } = useI18n()
  const authStore = useAuthStore()
  const { apiFetch } = useApiFetch()

  const switchLocale = async (code: string) => {
    await setLocale(code as LocaleCode)
    if (authStore.isLoggedIn && authStore.user) {
      await apiFetch(`/users/${authStore.user.id}`, {
        method: 'PUT',
        body: { locale: code },
      }).catch(() => {})
    }
  }

  return { locale, switchLocale }
}

import type { LocaleCode } from '~/types/locale'

/**
 * Load autoload options from the API once per full-page lifecycle.
 * On SSR the state is serialised into the Pinia payload and reused
 * by the client without a second network request.
 *
 * Also initialises the Nuxt i18n locale from site_language when the
 * user has no stored locale preference (i18n_locale cookie absent).
 */
export default defineNuxtPlugin(async () => {
  const store = useOptionsStore()
  await callOnce('options:load', () => store.load())

  const localeCookie = useCookie<string>('i18n_locale')
  if (!localeCookie.value) {
    const siteLanguage = store.get('site_language', 'zh')
    const { $i18n } = useNuxtApp()
    await $i18n.setLocale(siteLanguage as LocaleCode)
  }
})

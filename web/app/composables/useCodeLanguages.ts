import { DEFAULT_CODE_LANGUAGES, formatLanguageLabel } from '~/utils/codeLanguages'

/**
 * Returns the code language list for the editor code block dropdown.
 * Reads user config from optionsStore; falls back to DEFAULT_CODE_LANGUAGES.
 */
export function useCodeLanguages() {
  const optionsStore = useOptionsStore()

  const languages = computed(() => {
    const customLangs = optionsStore.getJSON<string[]>('code_languages', [])
    if (!customLangs.length) return DEFAULT_CODE_LANGUAGES

    return [
      { label: 'Plain', value: '' },
      ...customLangs.map(v => ({ label: formatLanguageLabel(v), value: v })),
    ]
  })

  return { languages }
}

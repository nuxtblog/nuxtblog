import { OPTIONS_SCHEMA, type OptionKey, type OptionValue } from '~/config/options'

/**
 * Type-safe accessor for site options.
 *
 * Unlike `useOptionsStore().get()` / `.getJSON()`, this composable infers the
 * return type automatically from the schema so no manual casting is needed.
 *
 * @example
 *   const { getOption } = useOption()
 *   const n       = getOption('posts_per_page')          // number
 *   const layout  = getOption('post_default_layout')     // string
 *   const widgets = getOption('blog_sidebar_widgets')    // WidgetConfig[]
 */
export const useOption = () => {
  const store = useOptionsStore()

  function getOption<K extends OptionKey>(key: K): OptionValue<K> {
    const schema = OPTIONS_SCHEMA[key]
    if (schema.type === 'boolean' || schema.type === 'json') {
      return store.getJSON(key, schema.default) as OptionValue<K>
    }
    if (schema.type === 'number') {
      const raw = store.get(key, String(schema.default))
      const parsed = parseFloat(raw)
      return (isNaN(parsed) ? schema.default : parsed) as OptionValue<K>
    }
    return store.get(key, schema.default as string) as OptionValue<K>
  }

  return { getOption, OPTIONS_SCHEMA }
}

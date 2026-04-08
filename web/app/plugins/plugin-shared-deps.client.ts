/**
 * Plugin Shared Dependencies — exposes Vue on window for plugin ESM imports.
 *
 * Plugins build with `external: ['vue']` and Rollup rewrites imports to
 * `/_shared/vue.mjs`, which re-exports from this window global.
 */
import * as Vue from 'vue'

export default defineNuxtPlugin(() => {
  ;(window as any).__nuxtblog_vue = Vue
})

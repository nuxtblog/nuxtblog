/**
 * Plugin Shared Dependencies — exposes Vue on window for plugin ESM imports.
 *
 * Plugins build with `external: ['vue']` and Rollup rewrites imports to
 * `/_shared/vue.mjs`, which re-exports from this window global.
 *
 * Also exposes a reactive theme tokens object (`window.__nuxtblog_theme`)
 * so plugins can programmatically access the current theme configuration
 * (e.g. primary color hex for SVG charts). Values are synced by theme.client.ts.
 */
import * as Vue from 'vue'
import { reactive } from 'vue'
import { DEFAULT_THEME, PRIMARY_COLORS, NEUTRAL_COLORS } from '~/config/theme'
import {
  // Existing (12)
  UButton, UCard, UBadge, UIcon, USkeleton, UTable,
  UInput, USelect, UTabs, UAlert, USwitch, USeparator,
  // Form (10)
  UTextarea, UCheckbox, UCheckboxGroup, URadioGroup,
  UFormField, UForm, UInputNumber, USelectMenu, USlider, UPinInput,
  // Overlays (7)
  UModal, UDrawer, USlideover, UPopover, UTooltip, UContextMenu, UDropdownMenu,
  // Data display (8)
  UAvatar, UAvatarGroup, UProgress, UPagination, UChip, UKbd, UCollapsible, UAccordion,
  // Other (2)
  ULink, UBreadcrumb,
  // Admin layout (4)
  AdminPageContainer, AdminPageHeader, AdminPageContent, AdminPageFooter,
} from '#components'

function findHex(colors: { name: string; hex: string }[], name: string): string {
  return colors.find(c => c.name === name)?.hex ?? ''
}

export default defineNuxtPlugin(() => {
  ;(window as any).__nuxtblog_vue = Vue
  ;(window as any).__nuxtblog_ui = {
    UButton, UCard, UBadge, UIcon, USkeleton, UTable,
    UInput, USelect, UTabs, UAlert, USwitch, USeparator,
    UTextarea, UCheckbox, UCheckboxGroup, URadioGroup,
    UFormField, UForm, UInputNumber, USelectMenu, USlider, UPinInput,
    UModal, UDrawer, USlideover, UPopover, UTooltip, UContextMenu, UDropdownMenu,
    UAvatar, UAvatarGroup, UProgress, UPagination, UChip, UKbd, UCollapsible, UAccordion,
    ULink, UBreadcrumb,
    AdminPageContainer, AdminPageHeader, AdminPageContent, AdminPageFooter,
  }
  // Legacy alias — kept for backward compatibility with plugins that import from @nuxtblog/admin
  ;(window as any).__nuxtblog_admin = {
    AdminPageContainer, AdminPageHeader, AdminPageContent, AdminPageFooter,
  }
  ;(window as any).__nuxtblog_theme = reactive({
    primary: DEFAULT_THEME.primary,
    neutral: DEFAULT_THEME.neutral,
    primaryHex: findHex(PRIMARY_COLORS, DEFAULT_THEME.primary),
    neutralHex: findHex(NEUTRAL_COLORS, DEFAULT_THEME.neutral),
    radius: DEFAULT_THEME.radius,
    colorMode: 'light' as 'light' | 'dark',
    font: DEFAULT_THEME.font,
    fontSize: DEFAULT_THEME.fontSize,
    /** Read a CSS variable's computed value from :root, e.g. getCssVar('--color-primary-300') */
    getCssVar(name: string): string {
      return getComputedStyle(document.documentElement).getPropertyValue(name).trim()
    },
  })
})

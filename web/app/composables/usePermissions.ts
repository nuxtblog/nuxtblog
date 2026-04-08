import { DEFAULT_ROLE_CAPABILITIES } from '~/config/permissions'
import type { Capability } from '~/config/permissions'

/**
 * Permission / Capability checker.
 *
 * Checks whether the current logged-in user has a given capability.
 * Role defaults are defined in ~/config/permissions.ts.
 * Overrides are stored in the options table under the key `role_capabilities`
 * and are loaded automatically via the options store.
 *
 * Example:
 *   const { can } = usePermissions()
 *   can('manage_options')   // true for admin, false for editor by default
 *   can('publish_posts')    // true for editor and admin
 */
export const usePermissions = () => {
  const authStore    = useAuthStore()
  const optionsStore = useOptionsStore()

  /**
   * Get the effective capability list for a role.
   * Role 4 (super admin) always gets all caps and cannot be overridden.
   * Other roles check options-stored overrides first, then fall back to defaults.
   */
  const capabilitiesForRole = (role: number): Capability[] => {
    if (role === 4) return DEFAULT_ROLE_CAPABILITIES[4] ?? []
    const overrides = optionsStore.getJSON<Record<string, Capability[]>>('role_capabilities', {})
    if (overrides[String(role)]?.length) {
      return overrides[String(role)]
    }
    return DEFAULT_ROLE_CAPABILITIES[role] ?? []
  }

  /**
   * Returns true if the current user has the given capability.
   * Unauthenticated users have no capabilities.
   */
  const can = (cap: Capability): boolean => {
    const role = authStore.user?.role
    if (!role) return false
    return capabilitiesForRole(role).includes(cap)
  }

  /**
   * Returns true if user has ALL of the given capabilities.
   */
  const canAll = (...caps: Capability[]): boolean => caps.every(can)

  /**
   * Returns true if user has ANY of the given capabilities.
   */
  const canAny = (...caps: Capability[]): boolean => caps.some(can)

  return { can, canAll, canAny }
}

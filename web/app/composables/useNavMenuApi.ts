import type { NavMenuItem, NavCustomMenu, NavMenuSlotKey } from '~/types/api/navMenu'

const CUSTOM_KEY = 'nav_custom_menus'

export function useNavMenuApi() {
  const { apiFetch } = useApiFetch()

  // ── Built-in slots ─────────────────────────────────────────────────────────

  async function loadSlot(slot: NavMenuSlotKey): Promise<NavMenuItem[]> {
    try {
      const res = await apiFetch<{ key: string; value: string }>(`/options/${slot}`)
      return JSON.parse(res.value) as NavMenuItem[]
    } catch {
      return []
    }
  }

  async function saveSlot(slot: NavMenuSlotKey, items: NavMenuItem[]): Promise<void> {
    await apiFetch(`/options/${slot}`, {
      method: 'PUT',
      body: { value: JSON.stringify(items), autoload: 1 },
    })
  }

  // ── Custom menus ───────────────────────────────────────────────────────────

  async function loadCustomMenus(): Promise<NavCustomMenu[]> {
    try {
      const res = await apiFetch<{ key: string; value: string }>(`/options/${CUSTOM_KEY}`)
      return JSON.parse(res.value) as NavCustomMenu[]
    } catch {
      return []
    }
  }

  async function saveCustomMenus(menus: NavCustomMenu[]): Promise<void> {
    await apiFetch(`/options/${CUSTOM_KEY}`, {
      method: 'PUT',
      body: { value: JSON.stringify(menus), autoload: 0 },
    })
  }

  return { loadSlot, saveSlot, loadCustomMenus, saveCustomMenus }
}

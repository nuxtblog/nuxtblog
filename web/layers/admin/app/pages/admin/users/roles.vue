<template>
  <AdminPageContainer>
    <AdminPageHeader :title="$t('admin.roles.title')" :subtitle="$t('admin.roles.subtitle')">
      <template #actions>
        <UButton color="primary" icon="i-tabler-plus" @click="openCreateModal">
          {{ $t('admin.roles.new_role') }}
        </UButton>
      </template>
    </AdminPageHeader>

    <AdminPageContent>
      <div class="flex flex-col lg:flex-row gap-6">

        <!-- ── 左侧：角色列表 ─────────────────────────────── -->
        <div class="lg:w-72 shrink-0 space-y-2">
          <div
            v-for="role in roles"
            :key="role.id"
            class="group flex items-center gap-3 p-3 rounded-md border cursor-pointer transition-all"
            :class="selectedRoleId === role.id
              ? 'border-primary bg-primary/5'
              : 'border-default bg-default hover:bg-elevated'"
            @click="selectedRoleId = role.id">
            <!-- 图标 -->
            <div
              class="size-9 rounded-md flex items-center justify-center shrink-0"
              :class="`bg-${role.color}-500/10`">
              <UIcon :name="role.icon" class="size-5" :class="`text-${role.color}-500`" />
            </div>
            <!-- 文字 -->
            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-2">
                <span class="text-sm font-medium text-highlighted truncate">{{ role.label }}</span>
                <UBadge v-if="role.immutable" :label="$t('admin.roles.immutable_badge')" color="neutral" variant="subtle" size="xs" />
              </div>
              <div class="text-xs text-muted mt-0.5">{{ $t('admin.roles.permissions_count', { n: capCount(role.id) }) }}</div>
            </div>
            <!-- 操作（非内置可删） -->
            <UDropdownMenu
              v-if="!role.immutable && role.id !== 1 && role.id !== 2"
              :items="[[
                { label: $t('admin.roles.edit_name'), icon: 'i-tabler-edit', onClick: () => openEditModal(role) },
                { label: $t('admin.roles.reset_default'), icon: 'i-tabler-refresh', onClick: () => resetRole(role.id) },
                { label: $t('admin.roles.delete_role'), icon: 'i-tabler-trash', color: 'error', onClick: () => openDeleteModal(role) },
              ]]"
              @click.stop>
              <UButton
                icon="i-tabler-dots-vertical"
                color="neutral" variant="ghost" size="xs" square
                class="opacity-0 group-hover:opacity-100 transition-opacity" />
            </UDropdownMenu>
          </div>

          <!-- 保存按钮 -->
          <div class="pt-2 border-t border-default space-y-2">
            <UButton
              color="primary" block icon="i-tabler-device-floppy"
              :loading="saving" @click="save">
              {{ $t('admin.roles.save_all') }}
            </UButton>
            <UButton
              color="neutral" variant="outline" block icon="i-tabler-refresh"
              :disabled="saving" @click="resetAll">
              {{ $t('admin.roles.reset_all') }}
            </UButton>
          </div>
        </div>

        <!-- ── 右侧：权限配置面板 ────────────────────────── -->
        <div class="flex-1 min-w-0">
          <template v-if="selectedRole">
            <!-- 面板头 -->
            <div class="flex items-center gap-3 mb-5 pb-4 border-b border-default">
              <div
                class="size-10 rounded-md flex items-center justify-center shrink-0"
                :class="`bg-${selectedRole.color}-500/10`">
                <UIcon :name="selectedRole.icon" class="size-6" :class="`text-${selectedRole.color}-500`" />
              </div>
              <div>
                <h3 class="font-semibold text-highlighted">{{ selectedRole.label }}</h3>
                <p class="text-sm text-muted">
                  {{ selectedRole.immutable ? $t('admin.roles.immutable_desc') : $t('admin.roles.enabled_n', { n: capCount(selectedRole.id) }) }}
                </p>
              </div>
              <div v-if="!selectedRole.immutable" class="ml-auto flex items-center gap-2">
                <UButton size="xs" color="neutral" variant="outline" icon="i-tabler-square-check" @click="selectAll(selectedRole.id)">{{ $t('admin.roles.select_all') }}</UButton>
                <UButton size="xs" color="neutral" variant="outline" icon="i-tabler-square" @click="clearAll(selectedRole.id)">{{ $t('admin.roles.clear_all') }}</UButton>
              </div>
            </div>

            <!-- 权限分组 -->
            <div class="space-y-5">
              <div v-for="(group, groupKey) in CAPABILITY_GROUPS" :key="groupKey">
                <!-- 组标题 + 组全选 -->
                <div class="flex items-center justify-between mb-2">
                  <span class="text-xs font-semibold text-muted uppercase tracking-wide">{{ group.label }}</span>
                  <UCheckbox
                    v-if="!selectedRole.immutable"
                    :model-value="isGroupAllSelected(selectedRole.id, group)"
                    :indeterminate="isGroupIndeterminate(selectedRole.id, group)"
                    size="xs"
                    @update:model-value="(val) => toggleGroup(selectedRole.id, group, val)" />
                </div>

                <!-- 权限项：每行一个 checkbox -->
                <div class="rounded-md border border-default divide-y divide-default overflow-hidden">
                  <label
                    v-for="(capLabel, cap) in group.caps"
                    :key="cap"
                    class="flex items-center gap-3 px-4 py-3 select-none"
                    :class="selectedRole.immutable
                      ? 'bg-muted/20 cursor-default opacity-60'
                      : 'hover:bg-elevated/50 transition-colors cursor-pointer'">
                    <UCheckbox
                      :model-value="hasCap(selectedRole.id, cap as Capability)"
                      :disabled="selectedRole.immutable"
                      @update:model-value="(val) => toggleCap(selectedRole.id, cap as Capability, val)" />
                    <div class="flex-1 min-w-0">
                      <div class="text-sm text-highlighted">{{ capLabel }}</div>
                      <div class="text-xs text-muted font-mono mt-0.5">{{ cap }}</div>
                    </div>
                  </label>
                </div>
              </div>
            </div>
          </template>

          <div v-else class="flex flex-col items-center justify-center h-48 text-muted">
            <UIcon name="i-tabler-arrow-left" class="size-8 mb-2" />
            <p class="text-sm">{{ $t('admin.roles.select_left') }}</p>
          </div>
        </div>
      </div>
    </AdminPageContent>

    <!-- 新建/编辑角色弹窗 -->
    <UModal v-model:open="showRoleModal" :title="editingRole ? $t('admin.roles.edit_modal_title') : $t('admin.roles.create_modal_title')">
      <template #content>
        <div class="p-6 space-y-4">
          <UFormField :label="$t('admin.roles.role_name_label')" required>
            <UInput v-model="roleForm.label" :placeholder="$t('admin.roles.role_name_placeholder')" class="w-full" />
          </UFormField>
          <UFormField :label="$t('admin.roles.base_role_label')">
            <USelect
              v-model="roleForm.baseRoleId"
              :items="inheritOptions"
              class="w-full" />
          </UFormField>
          <div class="flex justify-end gap-2 pt-2">
            <UButton color="neutral" variant="outline" @click="showRoleModal = false">{{ $t('common.cancel') }}</UButton>
            <UButton color="primary" :disabled="!roleForm.label.trim()" @click="confirmRoleModal">
              {{ editingRole ? $t('common.save') : $t('common.create') }}
            </UButton>
          </div>
        </div>
      </template>
    </UModal>

    <!-- 删除确认弹窗 -->
    <UModal v-model:open="showDeleteModal">
      <template #content>
        <div class="p-6">
          <div class="flex items-center gap-3 mb-4">
            <div class="size-10 rounded-full bg-error/10 flex items-center justify-center shrink-0">
              <UIcon name="i-tabler-alert-triangle" class="size-5 text-error" />
            </div>
            <div>
              <h3 class="font-semibold text-highlighted">{{ $t('admin.roles.delete_title', { name: deletingRole?.label }) }}</h3>
              <p class="text-sm text-muted mt-0.5">{{ $t('admin.roles.delete_desc') }}</p>
            </div>
          </div>
          <div class="flex justify-end gap-2 mt-6">
            <UButton color="neutral" variant="outline" @click="showDeleteModal = false">{{ $t('common.cancel') }}</UButton>
            <UButton color="error" @click="confirmDelete">{{ $t('common.delete') }}</UButton>
          </div>
        </div>
      </template>
    </UModal>
  </AdminPageContainer>
</template>

<script setup lang="ts">
import {
  CAPABILITY_GROUPS,
  DEFAULT_ROLE_CAPABILITIES,
  ROLE_LABELS,
  type Capability,
} from '~/config/permissions'

const toast = useToast()
const { t } = useI18n()
const optionApi = useOptionApi()
const saving = ref(false)

// ── Role definitions ───────────────────────────────────────────────────────────

interface RoleDef {
  id: number
  label: string
  color: string
  icon: string
  immutable: boolean
}

const BUILT_IN_ROLES: RoleDef[] = [
  { id: 1, label: ROLE_LABELS[1], color: 'neutral', icon: 'i-tabler-user',          immutable: false },
  { id: 2, label: ROLE_LABELS[2], color: 'primary', icon: 'i-tabler-pencil',        immutable: false },
  { id: 3, label: ROLE_LABELS[3], color: 'error',   icon: 'i-tabler-shield',        immutable: false },
  { id: 4, label: ROLE_LABELS[4], color: 'warning', icon: 'i-tabler-crown',         immutable: true  },
]

const customRoles = ref<RoleDef[]>([])
const roles = computed<RoleDef[]>(() => [...BUILT_IN_ROLES, ...customRoles.value])

const selectedRoleId = ref<number>(1)
const selectedRole = computed(() => roles.value.find(r => r.id === selectedRoleId.value) ?? null)

// ── Capability sets (role id string → Set) ─────────────────────────────────────

const ALL_CAPS = Object.values(CAPABILITY_GROUPS).flatMap(g => Object.keys(g.caps)) as Capability[]

const capSets = ref<Record<string, Set<Capability>>>({
  '1': new Set(DEFAULT_ROLE_CAPABILITIES[1]),
  '2': new Set(DEFAULT_ROLE_CAPABILITIES[2]),
  '3': new Set(DEFAULT_ROLE_CAPABILITIES[3]),
  '4': new Set(DEFAULT_ROLE_CAPABILITIES[4]),
})

// ── Load saved config ──────────────────────────────────────────────────────────

onMounted(async () => {
  try {
    // Load custom role definitions first so capSets keys are ready
    const defs = await optionApi.getOption('custom_role_defs')
    if (Array.isArray(defs)) {
      customRoles.value = defs as RoleDef[]
      const maxId = (defs as RoleDef[]).reduce((m, r) => Math.max(m, r.id), 9)
      nextCustomId = maxId + 1
    }

    const val = await optionApi.getOption('role_capabilities')
    if (val && typeof val === 'object') {
      const overrides = val as Record<string, string[]>
      for (const [roleId, caps] of Object.entries(overrides)) {
        if (roleId === '4') continue // super admin always all caps
        capSets.value[roleId] = new Set(caps as Capability[])
      }
    }
  } catch {}
})

// ── Helpers ────────────────────────────────────────────────────────────────────

const hasCap = (roleId: number, cap: Capability): boolean =>
  capSets.value[String(roleId)]?.has(cap) ?? false

const capCount = (roleId: number): number =>
  capSets.value[String(roleId)]?.size ?? 0

const toggleCap = (roleId: number, cap: Capability, enabled: boolean) => {
  const role = roles.value.find(r => r.id === roleId)
  if (role?.immutable) return
  // Replace the Set with a new instance so Vue detects the change.
  const prev = capSets.value[String(roleId)] ?? new Set<Capability>()
  const next = new Set(prev)
  if (enabled) next.add(cap)
  else next.delete(cap)
  capSets.value[String(roleId)] = next
}

const selectAll = (roleId: number) => {
  capSets.value[String(roleId)] = new Set(ALL_CAPS)
}
const clearAll = (roleId: number) => {
  capSets.value[String(roleId)] = new Set()
}

// 分组全选 helpers
const groupCaps = (group: (typeof CAPABILITY_GROUPS)[keyof typeof CAPABILITY_GROUPS]) =>
  Object.keys(group.caps) as Capability[]

const isGroupAllSelected = (roleId: number, group: any): boolean =>
  groupCaps(group).every(c => hasCap(roleId, c))

const isGroupIndeterminate = (roleId: number, group: any): boolean => {
  const caps = groupCaps(group)
  const selected = caps.filter(c => hasCap(roleId, c)).length
  return selected > 0 && selected < caps.length
}

const toggleGroup = (roleId: number, group: any, val: boolean) => {
  const prev = capSets.value[String(roleId)] ?? new Set<Capability>()
  const next = new Set(prev)
  groupCaps(group).forEach(c => val ? next.add(c) : next.delete(c))
  capSets.value[String(roleId)] = next
}

// ── Role modal (create / edit) ─────────────────────────────────────────────────

const showRoleModal = ref(false)
const editingRole = ref<RoleDef | null>(null)
const roleForm = ref({ label: '', baseRoleId: 2 })

const inheritOptions = computed(() => [
  { label: t('admin.roles.inherit_editor'), value: 2 },
  { label: t('admin.roles.inherit_admin'), value: 3 },
])

const openCreateModal = () => {
  editingRole.value = null
  roleForm.value = { label: '', baseRoleId: 2 }
  showRoleModal.value = true
}

const openEditModal = (role: RoleDef) => {
  editingRole.value = role
  roleForm.value = { label: role.label, baseRoleId: 2 }
  showRoleModal.value = true
}

let nextCustomId = 10 // will be updated after loading saved defs
const confirmRoleModal = () => {
  if (!roleForm.value.label.trim()) return
  if (editingRole.value) {
    // 编辑名称
    const r = customRoles.value.find(r => r.id === editingRole.value!.id)
    if (r) r.label = roleForm.value.label.trim()
  } else {
    // 新建：继承基础角色权限
    const id = nextCustomId++
    const baseCaps = capSets.value[String(roleForm.value.baseRoleId)]
    capSets.value[String(id)] = new Set(baseCaps)
    customRoles.value.push({
      id,
      label:      roleForm.value.label.trim(),
      color:      'success',
      icon:       'i-tabler-user-check',
      immutable:  false,
      baseRoleId: roleForm.value.baseRoleId, // retained for backend level resolution
    } as any)
    selectedRoleId.value = id
    toast.add({ title: t('admin.roles.role_created', { name: roleForm.value.label.trim() }), color: 'success' })
  }
  showRoleModal.value = false
}

// ── Delete role ────────────────────────────────────────────────────────────────

const showDeleteModal = ref(false)
const deletingRole = ref<RoleDef | null>(null)

const openDeleteModal = (role: RoleDef) => {
  deletingRole.value = role
  showDeleteModal.value = true
}

const confirmDelete = () => {
  if (!deletingRole.value) return
  const id = deletingRole.value.id
  customRoles.value = customRoles.value.filter(r => r.id !== id)
  delete capSets.value[String(id)]
  if (selectedRoleId.value === id) selectedRoleId.value = 1
  toast.add({ title: t('admin.roles.role_deleted'), color: 'neutral' })
  showDeleteModal.value = false
  deletingRole.value = null
}

// ── Reset ──────────────────────────────────────────────────────────────────────

const resetRole = (roleId: number) => {
  const defaults = DEFAULT_ROLE_CAPABILITIES[roleId]
  if (defaults) capSets.value[String(roleId)] = new Set(defaults)
  toast.add({ title: t('admin.roles.reset_role_success'), color: 'warning' })
}

const resetAll = () => {
  capSets.value = {
    '1': new Set(DEFAULT_ROLE_CAPABILITIES[1]),
    '2': new Set(DEFAULT_ROLE_CAPABILITIES[2]),
    '3': new Set(DEFAULT_ROLE_CAPABILITIES[3]),
    '4': new Set(DEFAULT_ROLE_CAPABILITIES[4]),
  }
  toast.add({ title: t('admin.roles.reset_all_success'), color: 'warning' })
}

// ── Save ───────────────────────────────────────────────────────────────────────

const save = async () => {
  saving.value = true
  try {
    const capPayload: Record<string, string[]> = {}
    for (const role of roles.value) {
      if (role.immutable) continue // 超级管理员不存储
      capPayload[String(role.id)] = [...(capSets.value[String(role.id)] ?? [])]
    }

    // Persist custom role definitions (id, label, baseRoleId) separately so the
    // backend can resolve the access level of custom role IDs (>4).
    const defsPayload = customRoles.value.map(r => ({
      id:         r.id,
      label:      r.label,
      baseRoleId: (r as any).baseRoleId ?? 2,
    }))

    await Promise.all([
      optionApi.setOption('role_capabilities', capPayload),
      optionApi.setOption('custom_role_defs',  defsPayload),
    ])
    toast.add({ title: t('admin.roles.saved'), color: 'success' })
  } catch (error: any) {
    toast.add({ title: t('common.save_failed'), description: error?.message, color: 'error' })
  } finally {
    saving.value = false
  }
}
</script>

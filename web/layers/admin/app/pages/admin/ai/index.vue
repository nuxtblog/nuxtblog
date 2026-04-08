<script setup lang="ts">
import type { AIConfig } from '~/composables/useAiApi'

definePageMeta({ layout: 'admin' })

const { t } = useI18n()
const toast = useToast()
const { listConfigs, createConfig, updateConfig, deleteConfig, activateConfig, testConfig } = useAiApi()

// ── State ─────────────────────────────────────────────────────────────────────

const rawLoading = ref(true)
const loading = useMinLoading(rawLoading)
const configs = ref<AIConfig[]>([])
const activeId = ref('')

const showModal = ref(false)
const isEditing = ref(false)
const submitting = ref(false)
const testing = ref<string | null>(null)
const activating = ref<string | null>(null)
const deletingId = ref<string | null>(null)
const showDeleteModal = ref(false)
const pendingDeleteId = ref('')

const form = ref({
  name: '',
  api_format: 'openai',
  label: '',
  api_key: '',
  model: '',
  base_url: '',
  timeout_ms: 30000,
})
const editingId = ref('')

// ── Quick-fill presets (cosmetic only, just fill defaults) ────────────────────

interface Preset {
  label: string
  icon: string
  api_format: 'openai' | 'claude'
  defaultModel: string
  modelOptions: string[]
}

const presets: Preset[] = [
  { label: 'OpenAI', icon: 'i-tabler-brand-openai', api_format: 'openai', defaultModel: 'gpt-4o-mini', modelOptions: ['gpt-4o-mini', 'gpt-4o', 'gpt-4-turbo', 'gpt-3.5-turbo'] },
  { label: 'Claude', icon: 'i-tabler-robot', api_format: 'claude', defaultModel: 'claude-3-5-sonnet-20241022', modelOptions: ['claude-3-5-sonnet-20241022', 'claude-3-5-haiku-20241022', 'claude-3-opus-20240229'] },
  { label: 'DeepSeek', icon: 'i-tabler-brain', api_format: 'openai', defaultModel: 'deepseek-chat', modelOptions: ['deepseek-chat', 'deepseek-reasoner'] },
  { label: 'Ollama', icon: 'i-tabler-server', api_format: 'openai', defaultModel: 'llama3.2', modelOptions: ['llama3.2', 'llama3.1', 'mistral', 'qwen2.5', 'gemma2'] },
  { label: '通义千问', icon: 'i-tabler-sparkles', api_format: 'openai', defaultModel: 'qwen-turbo', modelOptions: ['qwen-turbo', 'qwen-plus', 'qwen-max'] },
  { label: '智谱 GLM', icon: 'i-tabler-sparkles', api_format: 'openai', defaultModel: 'glm-4-flash', modelOptions: ['glm-4-flash', 'glm-4', 'glm-4-air'] },
]

const currentPreset = computed(() =>
  presets.find(p => p.label === form.value.label) ?? null,
)

const modelSuggestions = computed(() =>
  currentPreset.value?.modelOptions ?? [],
)

function applyPreset(p: Preset) {
  form.value.label = p.label
  form.value.api_format = p.api_format
  if (!isEditing.value || !form.value.model) {
    form.value.model = p.defaultModel
  }
}

// ── Data loading ──────────────────────────────────────────────────────────────

async function load() {
  rawLoading.value = true
  try {
    const res = await listConfigs()
    configs.value = res?.items ?? []
    activeId.value = res?.active_id ?? ''
  }
  catch (e: any) {
    toast.add({ title: e?.message ?? '加载失败', color: 'error' })
  }
  finally {
    rawLoading.value = false
  }
}

onMounted(load)

// ── Modal helpers ─────────────────────────────────────────────────────────────

function openCreate() {
  isEditing.value = false
  editingId.value = ''
  form.value = { name: '', api_format: 'openai', label: '', api_key: '', model: '', base_url: '', timeout_ms: 30000 }
  showModal.value = true
}

function openEdit(cfg: AIConfig) {
  isEditing.value = true
  editingId.value = cfg.id
  form.value = {
    name: cfg.name,
    api_format: cfg.api_format || 'openai',
    label: cfg.label ?? '',
    api_key: cfg.api_key,
    model: cfg.model,
    base_url: cfg.base_url,
    timeout_ms: cfg.timeout_ms,
  }
  showModal.value = true
}

async function handleSubmit() {
  if (!form.value.name || !form.value.model) {
    toast.add({ title: '请填写配置名称和模型名称', color: 'error' })
    return
  }
  submitting.value = true
  try {
    if (isEditing.value) {
      await updateConfig(editingId.value, form.value)
      toast.add({ title: '配置已更新', color: 'success' })
    }
    else {
      await createConfig(form.value)
      toast.add({ title: '配置已添加', color: 'success' })
    }
    showModal.value = false
    await load()
  }
  catch (e: any) {
    toast.add({ title: e?.message ?? '操作失败', color: 'error' })
  }
  finally {
    submitting.value = false
  }
}

// ── Actions ───────────────────────────────────────────────────────────────────

async function handleActivate(id: string) {
  activating.value = id
  try {
    await activateConfig(id)
    activeId.value = id
    configs.value.forEach(c => (c.is_active = c.id === id))
    toast.add({ title: '已切换为当前配置', color: 'success' })
  }
  catch (e: any) {
    toast.add({ title: e?.message ?? '切换失败', color: 'error' })
  }
  finally {
    activating.value = null
  }
}

async function handleTest(cfg: AIConfig) {
  testing.value = cfg.id
  try {
    const res = await testConfig(cfg.id)
    if (res?.ok) {
      toast.add({ title: res.message, color: 'success' })
    }
    else {
      toast.add({ title: res?.message ?? '测试失败', color: 'error' })
    }
  }
  catch (e: any) {
    toast.add({ title: e?.message ?? '连接失败', color: 'error' })
  }
  finally {
    testing.value = null
  }
}

function confirmDelete(id: string) {
  pendingDeleteId.value = id
  showDeleteModal.value = true
}

async function handleDelete() {
  deletingId.value = pendingDeleteId.value
  try {
    await deleteConfig(pendingDeleteId.value)
    toast.add({ title: '配置已删除', color: 'success' })
    showDeleteModal.value = false
    await load()
  }
  catch (e: any) {
    toast.add({ title: e?.message ?? '删除失败', color: 'error' })
  }
  finally {
    deletingId.value = null
  }
}

// ── Card display helpers ──────────────────────────────────────────────────────

function configIcon(cfg: AIConfig) {
  return cfg.api_format === 'claude' ? 'i-tabler-robot' : 'i-tabler-brand-openai'
}

function configFormatLabel(cfg: AIConfig) {
  return cfg.api_format === 'claude' ? 'Claude (Anthropic)' : 'OpenAI 兼容'
}
</script>

<template>
  <AdminPageContainer>
    <AdminPageHeader title="AI 配置" subtitle="管理 AI 服务配置，填写 Key、地址和模型即可开始使用">
      <template #actions>
        <UButton icon="i-tabler-plus" color="primary" @click="openCreate">
          添加配置
        </UButton>
      </template>
    </AdminPageHeader>

    <AdminPageContent>

      <!-- Loading -->
      <div v-if="loading" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        <UCard v-for="i in 3" :key="i">
          <div class="space-y-3">
            <div class="flex items-center gap-3">
              <USkeleton class="size-10 rounded-lg shrink-0" />
              <div class="flex-1 space-y-2">
                <USkeleton class="h-4 w-32" />
                <USkeleton class="h-3 w-20" />
              </div>
            </div>
            <USkeleton class="h-3 w-full" />
            <USkeleton class="h-3 w-3/4" />
            <div class="flex gap-2 pt-2">
              <USkeleton class="h-8 flex-1 rounded-md" />
              <USkeleton class="h-8 w-20 rounded-md" />
            </div>
          </div>
        </UCard>
      </div>

      <!-- Empty -->
      <div v-else-if="configs.length === 0" class="flex flex-col items-center justify-center py-20">
        <UIcon name="i-tabler-robot" class="size-16 text-muted mb-4" />
        <h3 class="text-lg font-medium text-highlighted mb-1">还没有 AI 配置</h3>
        <p class="text-sm text-muted mb-6">添加一个 AI 配置，填写接口地址、Key 和模型名称即可</p>
        <UButton icon="i-tabler-plus" color="primary" @click="openCreate">
          添加第一个配置
        </UButton>
      </div>

      <!-- Config cards grid -->
      <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        <div
          v-for="cfg in configs"
          :key="cfg.id"
          class="relative flex flex-col rounded-xl border-2 bg-default transition-all"
          :class="cfg.id === activeId
            ? 'border-primary shadow-md shadow-primary/10'
            : 'border-default hover:border-elevated hover:shadow-sm'">

          <!-- Active badge -->
          <div
            v-if="cfg.id === activeId"
            class="absolute -top-3 left-4 px-2.5 py-0.5 bg-primary text-white text-xs font-medium rounded-full shadow">
            当前使用
          </div>

          <div class="p-5 flex-1 flex flex-col gap-4">
            <!-- Header -->
            <div class="flex items-center gap-3">
              <div
                class="size-10 rounded-lg flex items-center justify-center shrink-0"
                :class="cfg.id === activeId ? 'bg-primary/10' : 'bg-elevated'">
                <UIcon
                  :name="configIcon(cfg)"
                  class="size-5"
                  :class="cfg.id === activeId ? 'text-primary' : 'text-muted'" />
              </div>
              <div class="flex-1 min-w-0">
                <p class="font-semibold text-highlighted truncate">{{ cfg.name }}</p>
                <p class="text-xs text-muted">{{ cfg.label || configFormatLabel(cfg) }}</p>
              </div>
              <UDropdownMenu
                :items="[
                  [{ label: '编辑', icon: 'i-tabler-pencil', onClick: () => openEdit(cfg) }],
                  [{ label: '删除', icon: 'i-tabler-trash', color: 'error' as const, onClick: () => confirmDelete(cfg.id) }],
                ]"
                :popper="{ placement: 'bottom-end' }">
                <UButton icon="i-tabler-dots-vertical" color="neutral" variant="ghost" size="xs" square />
              </UDropdownMenu>
            </div>

            <!-- Config details -->
            <div class="space-y-2 text-xs">
              <div class="flex items-center justify-between">
                <span class="text-muted">模型</span>
                <span class="font-mono text-highlighted bg-elevated px-2 py-0.5 rounded text-xs">{{ cfg.model }}</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-muted">格式</span>
                <UBadge :label="configFormatLabel(cfg)" :color="cfg.api_format === 'claude' ? 'warning' : 'primary'" variant="soft" size="xs" />
              </div>
              <div class="flex items-center justify-between">
                <span class="text-muted">API Key</span>
                <span class="font-mono text-highlighted">{{ cfg.api_key || '(无)' }}</span>
              </div>
              <div v-if="cfg.base_url" class="flex items-center justify-between">
                <span class="text-muted">Base URL</span>
                <span class="text-muted truncate max-w-36 text-right" :title="cfg.base_url">{{ cfg.base_url }}</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-muted">超时</span>
                <span class="text-muted">{{ cfg.timeout_ms / 1000 }}s</span>
              </div>
            </div>
          </div>

          <!-- Footer actions -->
          <div class="px-5 pb-5 flex items-center gap-2">
            <UButton
              variant="outline"
              color="neutral"
              size="sm"
              icon="i-tabler-plug"
              :loading="testing === cfg.id"
              class="flex-1"
              @click="handleTest(cfg)">
              测试连接
            </UButton>
            <UButton
              v-if="cfg.id !== activeId"
              size="sm"
              color="primary"
              variant="soft"
              icon="i-tabler-check"
              :loading="activating === cfg.id"
              @click="handleActivate(cfg.id)">
              启用
            </UButton>
            <UButton
              v-else
              size="sm"
              color="primary"
              variant="solid"
              icon="i-tabler-check"
              disabled>
              已启用
            </UButton>
          </div>
        </div>
      </div>

      <!-- Info box: available actions -->
      <div v-if="!loading" class="mt-8">
        <UCard>
          <template #header>
            <h2 class="text-base font-semibold text-highlighted">可用 AI 功能</h2>
          </template>
          <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
            <div
              v-for="feature in aiFeatures"
              :key="feature.id"
              class="flex items-start gap-3 p-3 rounded-lg bg-elevated/50">
              <div class="size-8 rounded-md bg-primary/10 flex items-center justify-center shrink-0">
                <UIcon :name="feature.icon" class="size-4 text-primary" />
              </div>
              <div>
                <p class="text-sm font-medium text-highlighted">{{ feature.name }}</p>
                <p class="text-xs text-muted mt-0.5">{{ feature.desc }}</p>
                <code class="text-xs font-mono text-muted mt-1 block">{{ feature.endpoint }}</code>
              </div>
            </div>
          </div>
        </UCard>
      </div>

    </AdminPageContent>

    <!-- ── Add / Edit Modal ──────────────────────────────────────────────────── -->
    <UModal v-model:open="showModal" :ui="{ width: 'sm:max-w-lg' }">
      <template #content>
        <div class="p-6 max-h-[90vh] overflow-y-auto">
          <h3 class="text-lg font-semibold text-highlighted mb-5">
            {{ isEditing ? '编辑 AI 配置' : '添加 AI 配置' }}
          </h3>

          <form class="space-y-4" @submit.prevent="handleSubmit">
            <!-- Name -->
            <UFormField label="配置名称" required>
              <UInput v-model="form.name" placeholder="例如：我的 AI 代理" class="w-full" />
            </UFormField>

            <!-- Quick-fill presets -->
            <div>
              <p class="text-sm text-muted mb-2">快速填写预设（可选）</p>
              <div class="flex flex-wrap gap-1.5">
                <button
                  v-for="p in presets"
                  :key="p.label"
                  type="button"
                  class="px-2.5 py-1 text-xs rounded-full border transition-colors"
                  :class="form.label === p.label
                    ? 'bg-primary text-white border-primary'
                    : 'border-default text-muted hover:border-primary/50 hover:text-highlighted'"
                  @click="applyPreset(p)">
                  {{ p.label }}
                </button>
              </div>
            </div>

            <!-- Label (display name) -->
            <UFormField label="服务商名称" hint="可选，显示在卡片上，如 OpenAI、My Proxy">
              <UInput v-model="form.label" placeholder="例如：My AI Proxy" class="w-full" />
            </UFormField>

            <!-- API Format toggle -->
            <UFormField label="API 格式" hint="绝大多数服务（含代理）选 OpenAI 兼容；直连 Anthropic 官方选 Claude">
              <div class="flex gap-2">
                <button
                  type="button"
                  class="flex-1 py-2 text-sm rounded-lg border-2 transition-colors"
                  :class="form.api_format === 'openai'
                    ? 'border-primary bg-primary/5 text-primary font-medium'
                    : 'border-default text-muted hover:border-primary/30'"
                  @click="form.api_format = 'openai'">
                  OpenAI 兼容
                </button>
                <button
                  type="button"
                  class="flex-1 py-2 text-sm rounded-lg border-2 transition-colors"
                  :class="form.api_format === 'claude'
                    ? 'border-warning bg-warning/5 text-warning font-medium'
                    : 'border-default text-muted hover:border-warning/30'"
                  @click="form.api_format = 'claude'">
                  Claude (Anthropic)
                </button>
              </div>
            </UFormField>

            <!-- Base URL -->
            <UFormField label="Base URL" required hint="API 根地址，例如 https://ai.example.com 或 http://localhost:11434">
              <UInput
                v-model="form.base_url"
                placeholder="https://your-ai-api.com"
                class="w-full" />
            </UFormField>

            <!-- API Key -->
            <UFormField label="API Key" hint="无需认证的本地服务可留空">
              <UInput
                v-model="form.api_key"
                type="password"
                placeholder="sk-..."
                class="w-full" />
            </UFormField>

            <!-- Model -->
            <UFormField label="模型名称" required hint="直接输入模型 ID，或点击下方建议">
              <UInput
                v-model="form.model"
                placeholder="例如 gpt-4o-mini、deepseek-chat..."
                class="w-full" />
              <div v-if="modelSuggestions.length > 0" class="flex flex-wrap gap-1.5 mt-2">
                <button
                  v-for="m in modelSuggestions"
                  :key="m"
                  type="button"
                  class="px-2 py-0.5 text-xs rounded border border-default text-muted hover:border-primary/50 hover:text-highlighted transition-colors font-mono"
                  @click="form.model = m">
                  {{ m }}
                </button>
              </div>
            </UFormField>

            <!-- Timeout -->
            <UFormField label="超时（毫秒）" hint="AI 生成时间较长时可适当调大">
              <UInput
                v-model.number="form.timeout_ms"
                type="number"
                :min="5000"
                :max="120000"
                :step="5000"
                class="w-full" />
            </UFormField>

            <div class="flex gap-3 pt-2 justify-end">
              <UButton color="neutral" variant="outline" type="button" @click="showModal = false">
                取消
              </UButton>
              <UButton type="submit" color="primary" :loading="submitting">
                {{ isEditing ? '保存修改' : '添加配置' }}
              </UButton>
            </div>
          </form>
        </div>
      </template>
    </UModal>

    <!-- ── Delete confirm modal ──────────────────────────────────────────────── -->
    <UModal v-model:open="showDeleteModal">
      <template #content>
        <div class="p-6">
          <div class="flex items-center gap-3 mb-4">
            <div class="size-10 rounded-full bg-error/10 flex items-center justify-center shrink-0">
              <UIcon name="i-tabler-alert-triangle" class="size-5 text-error" />
            </div>
            <div>
              <h3 class="font-semibold text-highlighted">删除 AI 配置</h3>
              <p class="text-sm text-muted mt-0.5">此操作不可撤销，确定要删除这个配置吗？</p>
            </div>
          </div>
          <div class="flex justify-end gap-2 mt-6">
            <UButton color="neutral" variant="outline" @click="showDeleteModal = false">取消</UButton>
            <UButton color="error" :loading="deletingId === pendingDeleteId" @click="handleDelete">
              删除
            </UButton>
          </div>
        </div>
      </template>
    </UModal>

  </AdminPageContainer>
</template>

<script lang="ts">
const aiFeatures = [
  { id: 'polish', name: '内容润色', desc: '对文章内容进行润色和改写', icon: 'i-tabler-sparkles', endpoint: 'POST /api/v1/ai/polish' },
  { id: 'summarize', name: '自动摘要', desc: '生成文章摘要/Excerpt', icon: 'i-tabler-list-check', endpoint: 'POST /api/v1/ai/summarize' },
  { id: 'suggest-tags', name: '推荐标签', desc: '根据内容自动提取标签', icon: 'i-tabler-tags', endpoint: 'POST /api/v1/ai/suggest-tags' },
  { id: 'from-url', name: '从 URL 生成', desc: '抓取网页并改写为博客文章', icon: 'i-tabler-link', endpoint: 'POST /api/v1/ai/from-url' },
  { id: 'translate', name: '内容翻译', desc: '支持中英日韩等多语言翻译', icon: 'i-tabler-language', endpoint: 'POST /api/v1/ai/translate' },
  { id: 'plugin', name: '插件 SDK', desc: '插件可通过 nuxtblog.ai 调用', icon: 'i-tabler-plug', endpoint: 'nuxtblog.ai.polish(ctx)' },
]
</script>

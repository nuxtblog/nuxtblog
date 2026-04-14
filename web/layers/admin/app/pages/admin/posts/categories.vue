<template>
  <AdminPageContainer>
    <AdminPageHeader
      :title="$t('admin.posts.categories.title')"
      :subtitle="$t('admin.posts.categories.subtitle')">
      <template #actions>
        <UButton icon="i-tabler-plus" color="primary" @click="openCreateModal">
          {{ $t("admin.posts.categories.new_category") }}
        </UButton>
      </template>
    </AdminPageHeader>

    <AdminPageContent>
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <!-- 左侧：新建/编辑表单 -->
        <div class="lg:col-span-1">
          <UCard>
            <template #header>
              <h2 class="text-base font-semibold text-highlighted">
                {{
                  isEditing
                    ? $t("admin.posts.categories.edit_category")
                    : $t("admin.posts.categories.add_category")
                }}
              </h2>
            </template>

            <form @submit.prevent="handleSubmit" class="space-y-4">
              <UFormField
                :label="$t('admin.posts.categories.name_label')"
                required>
                <UInput
                  v-model="formData.name"
                  required
                  maxlength="100"
                  :placeholder="$t('admin.posts.categories.name_placeholder')"
                  class="w-full" />
              </UFormField>

              <UFormField :label="$t('admin.posts.categories.slug_label')">
                <UInput
                  v-model="formData.slug"
                  maxlength="100"
                  :placeholder="$t('admin.posts.categories.slug_auto')"
                  class="w-full" />
              </UFormField>

              <ParentCategorySelector
                v-model.number="formData.parent_id"
                :exclude-id="editingTerm?.term_id" />

              <UFormField :label="$t('admin.posts.categories.desc_label')">
                <UTextarea
                  v-model="formData.description"
                  :rows="4"
                  maxlength="255"
                  :placeholder="$t('admin.posts.categories.desc_placeholder')"
                  class="w-full" />
                <template #hint>
                  {{ formData.description?.length || 0 }} / 255
                </template>
              </UFormField>

              <div class="flex gap-3 pt-2">
                <UButton
                  type="submit"
                  color="primary"
                  :loading="submitting"
                  class="flex-1">
                  {{
                    isEditing
                      ? $t("admin.posts.categories.update_btn")
                      : $t("admin.posts.categories.add_btn")
                  }}
                </UButton>
                <UButton
                  v-if="isEditing"
                  type="button"
                  color="neutral"
                  variant="outline"
                  @click="cancelEdit">
                  {{ $t("common.cancel") }}
                </UButton>
              </div>
            </form>
          </UCard>
        </div>

        <!-- 右侧：分类列表/树 -->
        <div class="lg:col-span-2">
          <UCard :ui="{ body: 'p-0' }">
            <!-- 工具栏 -->
            <div
              class="p-4 border-b border-default flex items-center justify-between gap-4">
              <UInput
                v-model="searchQuery"
                :placeholder="$t('admin.posts.categories.search_placeholder')"
                leading-icon="i-tabler-search"
                class="flex-1"
                size="sm" />
              <ViewModeSwitcher
                v-model="viewMode"
                :modes="[
                  {
                    value: 'tree',
                    title: $t('admin.posts.categories.tree_view'),
                  },
                  {
                    value: 'list',
                    title: $t('admin.posts.categories.list_view'),
                  },
                ]" />
            </div>

            <ClientOnly>
              <!-- 加载状态 -->
              <div v-if="loading" class="p-4 space-y-3">
                <div
                  v-for="i in 5"
                  :key="i"
                  class="flex items-center gap-4 p-4 bg-default border border-default rounded-md">
                  <USkeleton class="size-4 rounded" />
                  <div class="flex-1 space-y-1.5">
                    <div class="flex items-center gap-3">
                      <USkeleton
                        :class="`h-4 w-${i % 2 === 0 ? '32' : '24'}`" />
                      <USkeleton class="h-3 w-16" />
                      <USkeleton class="h-5 w-10 rounded-full" />
                    </div>
                    <USkeleton class="h-3 w-1/2" />
                  </div>
                </div>
              </div>

              <!-- 树状视图 -->
              <div
                v-else-if="viewMode === 'tree' && filteredTree.length > 0"
                class="p-4">
                <CategoryNode
                  v-model="visibleTree"
                  :parent-id="0"
                  @edit="openEditModal"
                  @delete="handleDelete" />
              </div>

              <!-- 列表视图 -->
              <div
                v-else-if="viewMode === 'list' && filteredList.length > 0"
                class="space-y-3 p-4">
                <div
                  v-for="category in filteredList"
                  :key="category.term_id"
                  class="flex items-center group gap-4 p-4 bg-default border border-default rounded-md hover:shadow-sm transition-all">
                  <div class="flex-1 min-w-0">
                    <div class="flex items-center gap-3">
                      <h3 class="text-sm font-semibold text-highlighted">
                        {{ category.name }}
                      </h3>
                      <span class="text-xs text-muted font-mono"
                        >({{ category.slug }})</span
                      >
                      <UBadge
                        :label="
                          $t('admin.posts.categories.post_count', {
                            n: category.count,
                          })
                        "
                        size="sm"
                        color="primary"
                        variant="soft" />
                    </div>
                    <p
                      v-if="category.description"
                      class="text-sm text-muted mt-2 line-clamp-1">
                      {{ category.description }}
                    </p>
                    <div
                      v-if="getParentName(category)"
                      class="flex items-center gap-1.5 mt-2 text-xs text-muted">
                      <UIcon name="i-tabler-arrow-right" class="size-3.5" />
                      <span
                        >{{ $t("admin.posts.categories.parent_label") }}
                        <strong>{{ getParentName(category) }}</strong></span
                      >
                    </div>
                  </div>

                  <div
                    class="shrink-0 opacity-0 group-hover:opacity-100 transition-opacity">
                    <UDropdownMenu
                      :items="getCategoryActions(category)"
                      :popper="{ placement: 'bottom-end' }">
                      <UButton
                        color="neutral"
                        variant="ghost"
                        icon="i-tabler-dots-vertical"
                        square
                        size="xs" />
                    </UDropdownMenu>
                  </div>
                </div>
              </div>

              <!-- 空状态 -->
              <div
                v-else
                class="flex flex-col items-center justify-center py-16">
                <UIcon
                  name="i-tabler-folder-open"
                  class="size-16 text-muted mb-4" />
                <h3 class="text-lg font-medium text-highlighted mb-1">
                  {{ $t("admin.posts.categories.no_categories") }}
                </h3>
              </div>

              <template #fallback>
                <div class="p-4 space-y-3">
                  <div
                    v-for="i in 5"
                    :key="i"
                    class="flex items-center gap-4 p-4 bg-default border border-default rounded-md">
                    <USkeleton class="size-4 rounded" />
                    <div class="flex-1 space-y-1.5">
                      <div class="flex items-center gap-3">
                        <USkeleton class="h-4 w-28" />
                        <USkeleton class="h-3 w-16" />
                        <USkeleton class="h-5 w-10 rounded-full" />
                      </div>
                      <USkeleton class="h-3 w-1/2" />
                    </div>
                  </div>
                </div>
              </template>
            </ClientOnly>
          </UCard>

          <!-- 统计信息 -->
          <ClientOnly>
            <UCard v-if="!loading && categories.length > 0" class="mt-4">
              <div class="grid grid-cols-3 gap-4 text-center">
                <div>
                  <div class="text-2xl font-semibold text-highlighted">
                    {{ categories.length }}
                  </div>
                  <div class="text-sm text-muted">
                    {{ $t("admin.posts.categories.total_categories") }}
                  </div>
                </div>
                <div>
                  <div class="text-2xl font-semibold text-highlighted">
                    {{ categories.filter((t) => !t.parent_id).length }}
                  </div>
                  <div class="text-sm text-muted">
                    {{ $t("admin.posts.categories.top_level") }}
                  </div>
                </div>
                <div>
                  <div class="text-2xl font-semibold text-highlighted">
                    {{ categories.reduce((sum, t) => sum + t.count, 0) }}
                  </div>
                  <div class="text-sm text-muted">
                    {{ $t("admin.posts.categories.related_posts") }}
                  </div>
                </div>
              </div>
            </UCard>
          </ClientOnly>
        </div>
      </div>

      <!-- 删除确认弹窗 -->
      <UModal v-model:open="showDeleteModal">
        <template #content>
          <div class="p-6">
            <div class="flex items-center gap-3 mb-4">
              <div
                class="size-10 rounded-full bg-error/10 flex items-center justify-center shrink-0">
                <UIcon
                  name="i-tabler-alert-triangle"
                  class="size-5 text-error" />
              </div>
              <div>
                <h3 class="font-semibold text-highlighted">
                  {{ $t("admin.posts.categories.confirm_delete") }}
                </h3>
                <p class="text-sm text-muted mt-0.5">
                  {{
                    $t("admin.posts.categories.delete_confirm", {
                      name: deleteConfirmTerm?.name,
                    })
                  }}
                </p>
              </div>
            </div>
            <div class="flex justify-end gap-2 mt-6">
              <UButton
                color="neutral"
                variant="outline"
                @click="showDeleteModal = false">
                {{ $t("common.cancel") }}
              </UButton>
              <UButton color="error" :loading="deleting" @click="confirmDelete">
                {{ $t("admin.posts.categories.confirm_delete") }}
              </UButton>
            </div>
          </div>
        </template>
      </UModal>
    </AdminPageContent>
  </AdminPageContainer>
</template>

<script setup lang="ts">
import type { TermDetailResponse } from "~/types/api/term";
import type {
  TermDetailWithChildren,
  TermWithLevel,
} from "~/types/models/term";

const categoryStore = useCategoryStore();
const { t } = useI18n();

const searchQuery = ref("");
const isEditing = ref(false);
const editingTerm = ref<TermDetailResponse | null>(null);
const submitting = ref(false);
const deleting = ref(false);
const deleteConfirmTerm = ref<TermDetailResponse | null>(null);
const showDeleteModal = ref(false);
const viewMode = ref<"tree" | "list">("tree");

const { getParentName } = categoryStore;
const { categories, loading: storeLoading } = storeToRefs(categoryStore);
const loading = useMinLoading(storeLoading);

type TermFormData = {
  name?: string;
  slug?: string;
  taxonomy?: "category" | "tag";
  description?: string;
  parent_id?: number | null;
};

const formData = ref<TermFormData>({
  name: "",
  slug: "",
  taxonomy: "category",
  description: "",
  parent_id: undefined,
});

const visibleTree = ref<TermDetailWithChildren[]>([]);

watch(
  () => categoryStore.categoryTree,
  (tree) => {
    visibleTree.value = JSON.parse(JSON.stringify(tree));
  },
  { immediate: true, deep: true },
);

const filteredList = computed(() => {
  if (!searchQuery.value) return categories.value;
  const query = searchQuery.value.toLowerCase();
  return categories.value.filter(
    (t) =>
      t.name.toLowerCase().includes(query) ||
      t.slug.toLowerCase().includes(query) ||
      t.description.toLowerCase().includes(query),
  );
});

const filteredTree = computed(() => {
  if (!searchQuery.value) return visibleTree.value;
  const query = searchQuery.value.toLowerCase();
  const ids = categoryStore.categories
    .filter(
      (t) =>
        t.name.toLowerCase().includes(query) ||
        t.slug.toLowerCase().includes(query),
    )
    .map((t) => t.term_id);

  type TreeNode = TermDetailWithChildren;
  const filterTree = (list: TreeNode[]): TreeNode[] =>
    list
      .map((node) => ({
        ...node,
        children: node.children ? filterTree(node.children) : [],
      }))
      .filter(
        (n) => ids.includes(n.term_id) || (n.children && n.children.length > 0),
      );

  return filterTree(visibleTree.value);
});

const openCreateModal = () => {
  isEditing.value = false;
  editingTerm.value = null;
  formData.value = { name: "", slug: "", description: "", parent_id: null };
};

const openEditModal = (term: TermDetailResponse) => {
  isEditing.value = true;
  editingTerm.value = term;
  formData.value = {
    name: term.name,
    slug: term.slug,
    description: term.description,
  };
};

const cancelEdit = () => {
  isEditing.value = false;
  editingTerm.value = null;
  formData.value = { name: "", slug: "", description: "" };
};

const toast = useToast();

const toSlug = (s: string) =>
  s
    .toLowerCase()
    .replace(/\s+/g, "-")
    .replace(/[^a-z0-9-]+/g, "")
    .replace(/^-+|-+$/g, "");

const handleSubmit = async () => {
  if (!formData.value.name?.trim()) return;
  submitting.value = true;
  try {
    const slug = toSlug(formData.value.slug || formData.value.name!);
    if (!slug) {
      toast.add({
        title: t("admin.posts.categories.invalid_slug"),
        description: t("admin.posts.categories.invalid_slug_desc"),
        color: "error",
      });
      return;
    }
    if (isEditing.value && editingTerm.value) {
      await categoryStore.updateCategory(editingTerm.value.term_id, {
        name: formData.value.name,
        slug,
        description: formData.value.description,
        parent_id: formData.value.parent_id,
      });
      toast.add({
        title: t("admin.posts.categories.updated"),
        color: "success",
      });
    } else {
      await categoryStore.addNewCategory({
        name: formData.value.name!,
        slug,
        description: formData.value.description,
        parent_id: formData.value.parent_id || undefined,
      });
      toast.add({
        title: t("admin.posts.categories.created"),
        color: "success",
      });
    }
    cancelEdit();
  } catch (err: any) {
    toast.add({
      title: t("admin.posts.categories.submit_failed"),
      description: err?.message || t("common.retry"),
      color: "error",
    });
  } finally {
    submitting.value = false;
  }
};

const handleDelete = (term: TermDetailResponse) => {
  deleteConfirmTerm.value = term;
  showDeleteModal.value = true;
};

const getCategoryActions = (category: TermDetailResponse) => [
  [
    {
      label: t("common.edit"),
      icon: "i-tabler-pencil",
      onSelect: () => openEditModal(category),
    },
  ],
  [
    {
      label: t("common.delete"),
      icon: "i-tabler-trash",
      color: "error" as const,
      onSelect: () => handleDelete(category),
    },
  ],
];

const confirmDelete = async () => {
  if (!deleteConfirmTerm.value) return;
  deleting.value = true;
  try {
    await categoryStore.deleteCategory(deleteConfirmTerm.value.term_id);
    showDeleteModal.value = false;
    deleteConfirmTerm.value = null;
    toast.add({ title: t("admin.posts.categories.deleted"), color: "success" });
  } catch (err: any) {
    toast.add({
      title: t("admin.posts.categories.delete_failed"),
      description: err?.message || t("common.retry"),
      color: "error",
    });
  } finally {
    deleting.value = false;
  }
};

onMounted(async () => {
  if (!categoryStore.isLoaded) {
    await categoryStore.loadCategories();
  }
});
</script>

<template>
  <AdminPageContainer>
    <AdminPageHeader :title="$t('admin.posts.tags.title')" :subtitle="$t('admin.posts.tags.subtitle')">
      <template #actions>
        <UButton icon="i-tabler-plus" color="primary" @click="openCreateModal">
          {{ $t('admin.posts.tags.new_tag') }}
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
                {{ isEditing ? $t('admin.posts.tags.edit_tag') : $t('admin.posts.tags.add_tag') }}
              </h2>
            </template>

            <form @submit.prevent="handleSubmit" class="space-y-4">
              <UFormField :label="$t('admin.posts.tags.name_label')" :hint="$t('admin.posts.tags.name_hint')" required>
                <UInput
                  v-model="formData.name"
                  required
                  maxlength="100"
                  :placeholder="$t('admin.posts.tags.name_placeholder')"
                  class="w-full"
                />
              </UFormField>

              <UFormField :label="$t('admin.posts.tags.slug_label')" :hint="$t('admin.posts.tags.slug_hint')">
                <UInput
                  v-model="formData.slug"
                  maxlength="100"
                  :placeholder="$t('admin.posts.tags.slug_label')"
                  class="w-full"
                />
              </UFormField>

              <UFormField :label="$t('admin.posts.tags.desc_label')">
                <UTextarea
                  v-model="formData.description"
                  :rows="4"
                  maxlength="255"
                  :placeholder="$t('admin.posts.tags.desc_placeholder')"
                  class="w-full"
                />
                <template #hint>{{ formData.description?.length || 0 }} / 255</template>
              </UFormField>

              <div class="flex gap-3 pt-2">
                <UButton type="submit" color="primary" :loading="submitting" class="flex-1">
                  {{ isEditing ? $t('admin.posts.tags.update_btn') : $t('admin.posts.tags.add_btn') }}
                </UButton>
                <UButton v-if="isEditing" type="button" color="neutral" variant="outline" @click="cancelEdit">
                  {{ $t('common.cancel') }}
                </UButton>
              </div>
            </form>
          </UCard>
        </div>

        <!-- 右侧：标签列表 -->
        <div class="lg:col-span-2">
          <UCard :ui="{ body: 'p-0' }">
            <div class="p-4 border-b border-default">
              <UInput
                v-model="searchQuery"
                :placeholder="$t('admin.posts.tags.search_placeholder')"
                leading-icon="i-tabler-search"
                :trailing-icon="searchQuery ? 'i-tabler-x' : undefined"
                @click:trailing="searchQuery = ''"
                class="w-full"
                size="sm"
              />
            </div>

            <!-- 加载状态 -->
            <div v-if="loading" class="space-y-3 p-4">
              <div v-for="i in 6" :key="i" class="flex items-center gap-4 p-4 bg-default border border-default rounded-md">
                <div class="flex-1 space-y-2">
                  <div class="flex items-center gap-3">
                    <USkeleton class="h-4 w-24" />
                    <USkeleton class="h-4 w-16" />
                    <USkeleton class="h-5 w-14 rounded-full" />
                  </div>
                  <USkeleton class="h-3 w-2/3" />
                </div>
                <div class="flex gap-1">
                  <USkeleton class="size-6 rounded" />
                  <USkeleton class="size-6 rounded" />
                </div>
              </div>
            </div>

            <!-- 标签列表 -->
            <div v-else-if="filteredTags.length > 0" class="space-y-3 p-4">
              <div
                v-for="tag in filteredTags"
                :key="tag.term_id"
                class="flex items-center group gap-4 p-4 bg-default border border-default rounded-md hover:shadow-sm transition-all"
              >
                  <div class="flex-1 min-w-0">
                    <div class="flex items-center gap-3 flex-wrap">
                      <h3 class="text-sm font-medium text-highlighted">{{ tag.name }}</h3>
                      <span class="text-xs text-muted">({{ tag.slug }})</span>
                      <UBadge :label="$t('admin.posts.tags.post_count', { n: tag.count })" color="neutral" variant="soft" size="xs" />
                    </div>
                    <p v-if="tag.description" class="text-xs text-muted mt-1 line-clamp-2">
                      {{ tag.description }}
                    </p>
                  </div>

                  <div class="shrink-0 opacity-0 group-hover:opacity-100 transition-opacity">
                    <UDropdownMenu
                      :items="getTagActions(tag)"
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
            <div v-else class="flex flex-col items-center justify-center py-16">
              <UIcon name="i-tabler-tag-off" class="size-16 text-muted mb-4" />
              <h3 class="text-lg font-medium text-highlighted mb-1">{{ $t('admin.posts.tags.no_tags') }}</h3>
              <p class="text-sm text-muted">
                {{ searchQuery ? $t('admin.posts.tags.no_results') : $t('admin.posts.tags.start_create') }}
              </p>
            </div>
          </UCard>

          <!-- 统计信息 -->
          <UCard v-if="!loading && tags.length > 0" class="mt-4">
            <div class="grid grid-cols-2 gap-4 text-center">
              <div>
                <div class="text-2xl font-semibold text-highlighted">{{ tags.length }}</div>
                <div class="text-sm text-muted">{{ $t('admin.posts.tags.total_tags') }}</div>
              </div>
              <div>
                <div class="text-2xl font-semibold text-highlighted">{{ totalPosts }}</div>
                <div class="text-sm text-muted">{{ $t('admin.posts.tags.related_posts') }}</div>
              </div>
            </div>
          </UCard>
        </div>
      </div>

      <!-- 删除确认弹窗 -->
      <UModal v-model:open="showDeleteModal">
        <template #content>
          <div class="p-6">
            <div class="flex items-center gap-3 mb-4">
              <div class="size-10 rounded-full bg-error/10 flex items-center justify-center shrink-0">
                <UIcon name="i-tabler-alert-triangle" class="size-5 text-error" />
              </div>
              <div>
                <h3 class="font-semibold text-highlighted">{{ $t('admin.posts.tags.confirm_delete') }}</h3>
                <p class="text-sm text-muted mt-0.5">
                  {{ $t('admin.posts.tags.delete_confirm', { name: deleteConfirmTag?.name }) }}
                </p>
              </div>
            </div>
            <div v-if="deleteConfirmTag && deleteConfirmTag.count > 0" class="flex items-center gap-2 text-sm text-warning mb-4">
              <UIcon name="i-tabler-alert-circle" class="size-4" />
              {{ $t('admin.posts.tags.delete_warning', { n: deleteConfirmTag.count }) }}
            </div>
            <div class="flex justify-end gap-2 mt-6">
              <UButton color="neutral" variant="outline" @click="showDeleteModal = false">{{ $t('common.cancel') }}</UButton>
              <UButton color="error" :loading="deleting" @click="confirmDelete">{{ $t('admin.posts.tags.confirm_delete') }}</UButton>
            </div>
          </div>
        </template>
      </UModal>
    </AdminPageContent>
  </AdminPageContainer>
</template>

<script setup lang="ts">
import type { TermDetailResponse } from "~/types/api/term";

const tagStore = useTagStore();
const { t } = useI18n();

const searchQuery = ref("");
const isEditing = ref(false);
const editingTag = ref<TermDetailResponse | null>(null);
const submitting = ref(false);
const deleteConfirmTag = ref<TermDetailResponse | null>(null);
const showDeleteModal = ref(false);
const deleting = ref(false);

const { loading: storeLoading, tags } = storeToRefs(tagStore);
const loading = useMinLoading(storeLoading);

const formData = ref({ name: "", slug: "", description: "" });

const filteredTags = computed(() => {
  if (!searchQuery.value) return tags.value;
  const q = searchQuery.value.toLowerCase();
  return tags.value.filter(
    (t) =>
      t.name.toLowerCase().includes(q) ||
      t.slug.toLowerCase().includes(q) ||
      (t.description && t.description.toLowerCase().includes(q))
  );
});

const totalPosts = computed(() => tags.value.reduce((sum, tag) => sum + (tag.count || 0), 0));

const openCreateModal = () => {
  isEditing.value = false;
  editingTag.value = null;
  formData.value = { name: "", slug: "", description: "" };
};

const openEditModal = (tag: TermDetailResponse) => {
  isEditing.value = true;
  editingTag.value = tag;
  formData.value = { name: tag.name, slug: tag.slug, description: tag.description || "" };
};

const cancelEdit = () => {
  isEditing.value = false;
  editingTag.value = null;
  formData.value = { name: "", slug: "", description: "" };
};

const toast = useToast();

// Sanitize to backend slug format: lowercase, a-z0-9 and hyphens only
const toSlug = (s: string) =>
  s.toLowerCase().replace(/\s+/g, "-").replace(/[^a-z0-9-]+/g, "").replace(/^-+|-+$/g, "");

const handleSubmit = async () => {
  if (!formData.value.name.trim()) return;
  submitting.value = true;
  try {
    const slug = toSlug(formData.value.slug || formData.value.name);
    if (!slug) {
      toast.add({ title: t("admin.posts.tags.invalid_slug"), description: t("admin.posts.tags.invalid_slug_desc"), color: "error" });
      return;
    }
    if (isEditing.value && editingTag.value) {
      await tagStore.updateTag(editingTag.value.term_id, { name: formData.value.name, slug, description: formData.value.description });
      toast.add({ title: t("admin.posts.tags.updated"), color: "success" });
    } else {
      await tagStore.addNewTag({ name: formData.value.name, slug, description: formData.value.description });
      toast.add({ title: t("admin.posts.tags.created"), color: "success" });
    }
    cancelEdit();
  } catch (err: any) {
    toast.add({ title: t("admin.posts.tags.submit_failed"), description: err?.message || t("common.retry"), color: "error" });
  } finally {
    submitting.value = false;
  }
};

const handleDelete = (tag: TermDetailResponse) => {
  deleteConfirmTag.value = tag;
  showDeleteModal.value = true;
};

const getTagActions = (tag: TermDetailResponse) => [
  [
    {
      label: t("common.edit"),
      icon: "i-tabler-pencil",
      onSelect: () => openEditModal(tag),
    },
  ],
  [
    {
      label: t("common.delete"),
      icon: "i-tabler-trash",
      color: "error" as const,
      onSelect: () => handleDelete(tag),
    },
  ],
];

const confirmDelete = async () => {
  if (!deleteConfirmTag.value) return;
  deleting.value = true;
  try {
    await tagStore.deleteTag(deleteConfirmTag.value.term_id);
    showDeleteModal.value = false;
    deleteConfirmTag.value = null;
    toast.add({ title: t("admin.posts.tags.deleted"), color: "success" });
  } catch (err: any) {
    toast.add({ title: t("admin.posts.tags.delete_failed"), description: err?.message || t("common.retry"), color: "error" });
  } finally {
    deleting.value = false;
  }
};

onMounted(() => {
  if (!tagStore.isLoaded) {
    tagStore.loadTags();
  }
});
</script>

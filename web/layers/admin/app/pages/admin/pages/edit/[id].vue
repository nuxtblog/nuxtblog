<template>
  <div v-if="pending" class="flex items-center justify-center min-h-screen">
    <UIcon name="i-tabler-loader-2" class="size-8 text-muted animate-spin" />
  </div>
  <div v-else-if="!initialData">
    <AdminPageContainer>
      <AdminPageHeader title="页面不存在" subtitle="未找到该页面" />
    </AdminPageContainer>
  </div>
  <AdminPostEditor
    v-else
    mode="edit"
    :simple="true"
    :initial-data="initialData"
    :submitting="submitting"
    @save="handleSave" />
</template>

<script setup lang="ts">
import type { CreatePostRequest, UpdatePostRequest } from "~/types/api/post";
import type { PostEditorInitialData } from "~/components/AdminPostEditor.vue";

const route = useRoute();
const toast = useToast();
const postApi = usePostApi();
const postStore = usePostStore();

const submitting = ref(false);

const statusStrToInt: Record<string, number> = {
  draft: 1,
  published: 2,
  private: 3,
  archived: 4,
};
const commentStatusStrToInt: Record<string, number> = { open: 1, closed: 0 };

const { data: initialData, pending } = await useAsyncData(
  `page-edit-${route.params.id}`,
  async (): Promise<PostEditorInitialData | null> => {
    const post = await postApi
      .getPost(Number(route.params.id))
      .catch(() => null);
    if (!post) return null;

    return {
      id: post.id,
      title: post.title,
      slug: post.slug,
      content: post.content,
      excerpt: post.excerpt,
      status: statusStrToInt[post.status] ?? 1,
      postType: 2, // always page
      commentStatus: commentStatusStrToInt[post.comment_status] ?? 1,
      featuredImgId: post.featured_img?.id,
      featuredImgUrl: post.featured_img?.url,
      publishedAt: post.published_at,
      categoryTaxonomyIds: [],
      selectedTagObjects: [],
      isBanner: false,
      isFeatured: false,
      customMetaFields: [],
    };
  },
);

const handleSave = async (payload: CreatePostRequest | UpdatePostRequest) => {
  if (!initialData.value?.id) return;
  submitting.value = true;
  try {
    await postStore.updatePost(
      initialData.value.id,
      payload as UpdatePostRequest,
    );
    toast.add({ title: "页面已更新", color: "success" });
  } catch (error: any) {
    const msg = error?.data?.message || error?.message || "未知错误";
    toast.add({ title: "更新失败", description: msg, color: "error" });
  } finally {
    submitting.value = false;
  }
};
</script>

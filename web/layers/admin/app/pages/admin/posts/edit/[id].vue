<template>
  <div v-if="pending" class="flex items-center justify-center min-h-screen">
    <UIcon name="i-tabler-loader-2" class="size-8 text-muted animate-spin" />
  </div>
  <div v-else-if="!initialData">
    <AdminPageContainer>
      <AdminPageHeader title="文章不存在" subtitle="未找到该文章" />
    </AdminPageContainer>
  </div>
  <AdminPostEditor
    v-else
    ref="editorRef"
    mode="edit"
    :initial-data="initialData"
    :submitting="submitting"
    @save="handleSave" />
</template>

<script setup lang="ts">
import type { CreatePostRequest, UpdatePostRequest } from "~/types/api/post";
import type { PostEditorInitialData } from "~/components/AdminPostEditor.vue";
import type { TermDetailResponse } from "~/types/api/term";

const route = useRoute();
const toast = useToast();
const { t } = useI18n();
const postApi = usePostApi();
const postStore = usePostStore();
const { apiFetch } = useApiFetch();

const submitting = ref(false);

const editorRef = useTemplateRef<{
  isDirty: Ref<boolean>;
  getIsDirty: () => boolean;
  markSaved: () => void;
  seoData: Ref<Record<string, string>>;
}>("editorRef");

// 离开页面未保存警告
onBeforeRouteLeave(() => {
  if (editorRef.value?.getIsDirty()) {
    return window.confirm("有未保存的内容，确定要离开吗？");
  }
});

const handleBeforeUnload = (e: BeforeUnloadEvent) => {
  if (editorRef.value?.getIsDirty()) {
    e.preventDefault();
  }
};
onMounted(() => window.addEventListener("beforeunload", handleBeforeUnload));
onUnmounted(() =>
  window.removeEventListener("beforeunload", handleBeforeUnload),
);

const statusStrToInt: Record<string, number> = {
  draft: 1,
  published: 2,
  private: 3,
  archived: 4,
};
const postTypeStrToInt: Record<string, number> = {
  post: 1,
  page: 2,
  custom: 3,
};
const commentStatusStrToInt: Record<string, number> = {
  open: 1,
  closed: 0,
};

interface TaxonomyItem {
  id: number;
  term_id: number;
  taxonomy: string;
  description: string;
  parent_id?: number;
  post_count: number;
  extra: string;
  term?: { id: number; name: string; slug: string };
}

const { data: initialData, pending } = await useAsyncData(
  `post-edit-${route.params.id}`,
  async (): Promise<PostEditorInitialData | null> => {
    const post = await postApi
      .getPost(Number(route.params.id))
      .catch(() => null);
    if (!post) return null;

    const taxonomyRes = await apiFetch<{ list: TaxonomyItem[] }>(
      "/object-taxonomies",
      { params: { object_id: post.id, object_type: "post" } },
    ).catch(() => ({ list: [] as TaxonomyItem[] }));

    const taxonomies = taxonomyRes?.list ?? [];
    const categoryTaxonomyIds = taxonomies
      .filter((t) => t.taxonomy === "category")
      .map((t) => t.id);
    const selectedTagObjects: TermDetailResponse[] = taxonomies
      .filter((t) => t.taxonomy === "tag")
      .map((t) => ({
        term_id: t.term_id,
        name: t.term?.name ?? "",
        slug: t.term?.slug ?? "",
        taxonomy: t.taxonomy as "tag",
        count: t.post_count,
        description: t.description,
        term_taxonomy_id: t.id,
        parent_id: t.parent_id,
      }));

    const metas = post.metas ?? {};
    const isBanner = metas.is_banner === "1";
    const isFeatured = metas.is_featured === "1";
    const customMetaFields = Object.entries(metas)
      .filter(
        ([k]) =>
          ![
            "is_banner",
            "is_featured",
            "post_layout",
            "post_sidebar",
            "post_downloads",
          ].includes(k),
      )
      .map(([key, value]) => ({ key, value }));

    // SEO — included in admin get-by-id response
    const postAny = post as any;
    const seo = postAny.seo ?? {};

    return {
      id: post.id,
      title: post.title,
      slug: post.slug,
      content: post.content,
      excerpt: post.excerpt,
      status: statusStrToInt[post.status] ?? 1,
      postType: postTypeStrToInt[post.post_type] ?? 1,
      commentStatus: commentStatusStrToInt[post.comment_status] ?? 1,
      featuredImgId: post.featured_img?.id,
      featuredImgUrl: post.featured_img?.url,
      publishedAt: post.published_at,
      authorId: (post as any).author_id ?? post.author?.id,
      categoryTaxonomyIds,
      selectedTagObjects,
      isBanner,
      isFeatured,
      customMetaFields,
      seo: {
        meta_title: seo.meta_title ?? "",
        meta_desc: seo.meta_desc ?? "",
        og_title: seo.og_title ?? "",
        og_image: seo.og_image ?? "",
        canonical_url: seo.canonical_url ?? "",
        robots: seo.robots ?? "index,follow",
      },
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

    // 同步保存 SEO 数据
    const seo = editorRef.value?.seoData?.value;
    if (seo) {
      await apiFetch(`/posts/${initialData.value.id}/seo`, {
        method: "PUT",
        body: seo,
      }).catch(() => {});
    }

    editorRef.value?.markSaved();
    const isScheduled = payload.status === 1 && !!(payload as any).published_at &&
      new Date((payload as any).published_at).getTime() > Date.now()
    toast.add({
      title: isScheduled ? t('admin.posts.editor.scheduled_saved') : t('admin.posts.editor.save'),
      icon: isScheduled ? 'i-tabler-clock' : 'i-tabler-circle-check',
      color: 'success',
    });
  } catch (error: any) {
    const msg = error?.data?.message || error?.message || "未知错误";
    toast.add({ title: "更新失败", description: msg, color: "error" });
  } finally {
    submitting.value = false;
  }
};
</script>

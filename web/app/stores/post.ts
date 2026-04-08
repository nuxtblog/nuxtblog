import type {
  PostResponse,
  PostListItemResponse,
  PaginationResponse,
} from "~/types/api/post";

export const usePostStore = defineStore("post", () => {
  const posts = ref<PostListItemResponse[]>([]);
  const currentPost = ref<PostResponse | null>(null);
  const loading = ref(false);
  const detailLoading = ref(false);
  const error = ref<string | null>(null);
  const pagination = ref<{
    page: number;
    pageSize: number;
    total: number;
    totalPages: number;
  }>({
    page: 1,
    pageSize: 10,
    total: 0,
    totalPages: 0,
  });

  const postApi = usePostApi();

  /**
   * 获取文章列表
   */
  const loadPosts = async (page = 1, pageSize = 10, params?: any) => {
    loading.value = true;
    error.value = null;
    try {
      const res = await postApi.getPosts({
        page,
        pageSize,
        ...params,
      });
      posts.value = res.data;
      pagination.value = {
        page: res.page,
        pageSize: res.page_size,
        total: res.total,
        totalPages: res.total_pages,
      };
      return res;
    } catch (err) {
      console.error("获取文章列表失败:", err);
      error.value = "获取文章列表失败";
      return null;
    } finally {
      loading.value = false;
    }
  };

  /**
   * 根据 ID 获取文章详情
   */
  const getPostById = async (id: number): Promise<PostResponse | null> => {
    detailLoading.value = true;
    error.value = null;
    try {
      const post = await postApi.getPost(id);
      currentPost.value = post;
      return post;
    } catch (err) {
      console.error("获取文章详情失败:", err);
      error.value = "获取文章详情失败";
      return null;
    } finally {
      detailLoading.value = false;
    }
  };

  /**
   * 根据 slug 获取文章详情
   */
  const getPostBySlug = async (slug: string): Promise<PostResponse | null> => {
    detailLoading.value = true;
    error.value = null;
    try {
      const post = await postApi.getPostBySlug(slug);
      currentPost.value = post;
      return post;
    } catch (err) {
      console.error("获取文章详情失败:", err);
      error.value = "获取文章详情失败";
      return null;
    } finally {
      detailLoading.value = false;
    }
  };

  /**
   * 创建文章
   */
  const addPost = async (data: any): Promise<PostResponse | null> => {
    try {
      const createdPost = await postApi.createPost(data);
      posts.value.unshift(createdPost as any);
      return createdPost;
    } catch (err) {
      console.error("创建文章失败:", err);
      throw err;
    }
  };

  /**
   * 更新文章
   */
  const updatePost = async (
    postId: number,
    data: any
  ): Promise<PostResponse | null> => {
    try {
      const updated = await postApi.updatePost(postId, data);
      const index = posts.value.findIndex((p) => p.id === postId);
      if (index !== -1) {
        posts.value[index] = updated as any;
      }
      if (currentPost.value?.id === postId) {
        currentPost.value = updated;
      }
      return updated;
    } catch (err) {
      console.error("更新文章失败:", err);
      throw err;
    }
  };

  /**
   * 删除文章
   */
  const deletePost = async (postId: number) => {
    try {
      await postApi.deletePost(postId);
      posts.value = posts.value.filter((p) => p.id !== postId);
      if (currentPost.value?.id === postId) {
        currentPost.value = null;
      }
    } catch (err) {
      console.error("删除文章失败:", err);
      throw err;
    }
  };

  /**
   * 增加文章浏览量
   */
  const incrementView = async (postId: number) => {
    try {
      await postApi.incrementViewCount(postId);
      if (currentPost.value?.id === postId && currentPost.value) {
        currentPost.value.view_count = (currentPost.value.view_count || 0) + 1;
      }
    } catch (err) {
      console.error("增加浏览量失败:", err);
    }
  };

  /**
   * 验证文章密码
   */
  const verifyPassword = async (
    postId: number,
    password: string
  ): Promise<boolean> => {
    const res = await postApi.verifyPostPassword(postId, { password });
    return res.valid;
  };

  /**
   * 清空当前文章
   */
  const clearCurrentPost = () => {
    currentPost.value = null;
  };

  /**
   * 是否已加载列表
   */
  const isLoaded = computed(() => posts.value.length > 0);

  return {
    posts,
    currentPost,
    loading,
    detailLoading,
    error,
    pagination,
    isLoaded,
    loadPosts,
    getPostById,
    getPostBySlug,
    addPost,
    updatePost,
    deletePost,
    incrementView,
    verifyPassword,
    clearCurrentPost,
  };
});

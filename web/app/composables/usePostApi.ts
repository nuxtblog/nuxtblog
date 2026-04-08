// composables/usePostApi.ts
import type {
  CreatePostRequest,
  UpdatePostRequest,
  PostQueryRequest,
  PostResponse,
  PostListItemResponse,
  PaginationResponse,
  VerifyPasswordRequest,
} from "~/types/api/post";

export const usePostApi = () => {
  const { apiFetch } = useApiFetch();

  /**
   * 获取文章列表
   */
  const getPosts = async (params?: PostQueryRequest) => {
    return await apiFetch<PaginationResponse<PostListItemResponse>>("/posts", {
      method: "GET",
      params,
    });
  };

  /**
   * 获取单个文章详情
   */
  const getPost = async (id: number) => {
    return await apiFetch<PostResponse>(`/posts/${id}`, {
      method: "GET",
    });
  };

  /**
   * 通过 slug 获取文章
   */
  const getPostBySlug = async (slug: string) => {
    return await apiFetch<PostResponse>(`/posts/slug/${slug}`, {
      method: "GET",
    });
  };

  /**
   * 创建文章
   */
  const createPost = async (data: CreatePostRequest) => {
    return await apiFetch<PostResponse>("/posts", {
      method: "POST",
      body: data,
    });
  };

  /**
   * 更新文章
   */
  const updatePost = async (id: number, data: UpdatePostRequest) => {
    return await apiFetch<PostResponse>(`/posts/${id}`, {
      method: "PUT",
      body: data,
    });
  };

  /**
   * 删除文章
   */
  const deletePost = async (id: number) => {
    return await apiFetch<void>(`/posts/${id}`, {
      method: "DELETE",
    });
  };

  /**
   * 验证文章密码
   */
  const verifyPostPassword = async (
    id: number,
    data: VerifyPasswordRequest
  ) => {
    return await apiFetch<{ valid: boolean }>(
      `/posts/${id}/verify-password`,
      {
        method: "POST",
        body: data,
      }
    );
  };

  /**
   * 增加文章浏览量
   */
  const incrementViewCount = async (id: number) => {
    return await apiFetch<void>(`/posts/${id}/view`, {
      method: "POST",
    });
  };

  /**
   * 批量更新文章字段
   */
  const batchUpdatePosts = async (body: {
    ids: number[]
    featured_img_id?: number | null
    status?: number
    term_taxonomy_ids?: number[]
    author_id?: number
  }) => {
    return await apiFetch<{ affected: number }>('/posts/batch', {
      method: 'PATCH',
      body,
    })
  }

  return {
    getPosts,
    getPost,
    getPostBySlug,
    createPost,
    updatePost,
    deletePost,
    verifyPostPassword,
    incrementViewCount,
    batchUpdatePosts,
  };
};

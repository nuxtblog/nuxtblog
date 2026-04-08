import type {
  PostMetaRequest,
  PostMetaResponse,
  BatchPostMetaRequest,
} from "~/types/api/postMeta";

export const usePostMetaApi = () => {
  const config = useRuntimeConfig();
  const baseURL = config.public.apiBase || "/api";

  /**
   * 创建单个元数据
   */
  const createMeta = async (postId: number, data: PostMetaRequest) => {
    return await $fetch<PostMetaResponse>(`${baseURL}/posts/${postId}/meta`, {
      method: "POST",
      body: data,
    });
  };

  /**
   * 批量创建元数据
   */
  const batchCreateMeta = async (
    postId: number,
    data: BatchPostMetaRequest
  ) => {
    return await $fetch<void>(`${baseURL}/posts/${postId}/meta/batch`, {
      method: "POST",
      body: data,
    });
  };

  /**
   * 更新元数据
   */
  const updateMeta = async (
    postId: number,
    metaKey: string,
    metaValue: string
  ) => {
    return await $fetch<PostMetaResponse>(
      `${baseURL}/posts/${postId}/meta/${encodeURIComponent(metaKey)}`,
      {
        method: "PUT",
        body: { meta_value: metaValue },
      }
    );
  };

  /**
   * 删除元数据
   */
  const deleteMeta = async (postId: number, metaKey: string) => {
    return await $fetch<void>(
      `${baseURL}/posts/${postId}/meta/${encodeURIComponent(metaKey)}`,
      { method: "DELETE" }
    );
  };

  /**
   * 获取单个元数据
   */
  const getMeta = async (postId: number, metaKey: string) => {
    return await $fetch<PostMetaResponse>(
      `${baseURL}/posts/${postId}/meta/${encodeURIComponent(metaKey)}`
    );
  };

  /**
   * 获取文章的所有元数据
   */
  const getMetasByPostId = async (postId: number) => {
    return await $fetch<PostMetaResponse[]>(`${baseURL}/posts/${postId}/meta`);
  };

  /**
   * 更新或创建元数据（upsert）
   */
  const updateOrCreateMeta = async (postId: number, data: PostMetaRequest) => {
    return await $fetch<PostMetaResponse>(
      `${baseURL}/posts/${postId}/meta/upsert`,
      { method: "POST", body: data }
    );
  };

  return {
    createMeta,
    batchCreateMeta,
    updateMeta,
    deleteMeta,
    getMeta,
    getMetasByPostId,
    updateOrCreateMeta,
  };
};

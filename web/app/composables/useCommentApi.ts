// composables/useCommentApi.ts
import type {
  CreateCommentRequest,
  UpdateCommentRequest,
  CommentQueryRequest,
  CommentResponse,
  CommentListResponse,
  CommentTreeResponse,
  CommentStatsResponse,
  BatchCommentUpdateStatusRequest,
} from "~/types/api/comment";

export const useCommentApi = () => {
  const { apiFetch } = useApiFetch();

  /**
   * 创建评论
   */
  const createComment = async (data: CreateCommentRequest) => {
    return await apiFetch<CommentResponse>('/comments', {
      method: "POST",
      body: data,
    });
  };

  /**
   * 更新评论
   */
  const updateComment = async (id: number, data: UpdateCommentRequest) => {
    return await apiFetch<CommentResponse>(`/comments/${id}`, {
      method: "PUT",
      body: data,
    });
  };

  /**
   * 删除评论
   */
  const deleteComment = async (id: number) => {
    return await apiFetch<void>(`/comments/${id}`, {
      method: "DELETE",
    });
  };

  /**
   * 获取评论详情
   */
  const getCommentById = async (id: number) => {
    return await apiFetch<CommentResponse>(`/comments/${id}`, {
      method: "GET",
    });
  };

  /**
   * 获取评论列表（分页、查询）
   */
  const getComments = async (params?: CommentQueryRequest) => {
    return await apiFetch<{ list: CommentListResponse[]; total: number; page: number; page_size: number }>('/comments', {
      method: "GET",
      params,
    });
  };

  /**
   * 获取评论树形结构
   */
  const getCommentTree = async (
    post_id: number,
    status?: string,
    page?: number,
    page_size?: number
  ) => {
    return await apiFetch<CommentTreeResponse>('/comments/tree', {
      method: "GET",
      params: {
        post_id,
        status,
        page,
        page_size,
      },
    });
  };

  /**
   * 获取评论统计
   */
  const getCommentStats = async (post_id?: number) => {
    return await apiFetch<CommentStatsResponse>('/comments/stats', {
      method: "GET",
      params: { post_id },
    });
  };

  /**
   * 批量更新状态
   */
  const batchUpdateStatus = async (data: BatchCommentUpdateStatusRequest) => {
    return await apiFetch<void>('/comments/batch/status', {
      method: "PUT",
      body: data,
    });
  };

  /**
   * 批量删除
   */
  const batchDelete = async (ids: number[]) => {
    return await apiFetch<void>('/comments/batch', {
      method: "DELETE",
      body: { ids },
    });
  };

  /**
   * 批准评论
   */
  const approveComment = async (id: number) => {
    return await apiFetch(`/comments/${id}/approve`, {
      method: "POST",
    });
  };

  /**
   * 拒绝评论
   */
  const rejectComment = async (id: number) => {
    return await apiFetch(`/comments/${id}/reject`, {
      method: "POST",
    });
  };

  /**
   * 标记为垃圾评论
   */
  const markAsSpam = async (id: number) => {
    return await apiFetch(`/comments/${id}/spam`, {
      method: "POST",
    });
  };

  return {
    createComment,
    updateComment,
    deleteComment,
    getCommentById,
    getComments,
    getCommentTree,
    getCommentStats,
    batchUpdateStatus,
    batchDelete,
    approveComment,
    rejectComment,
    markAsSpam,
  };
};

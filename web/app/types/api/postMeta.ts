// 单个元数据请求
export interface PostMetaRequest {
  /** 元数据 key，最大长度 191 */
  meta_key: string;
  /** 元数据值 */
  meta_value: string;
}

// 批量元数据请求
export interface BatchPostMetaRequest {
  /** 元数据列表 */
  metas: PostMetaRequest[];
}

// 元数据响应
export interface PostMetaResponse {
  /** 元数据ID */
  meta_id: number;
  /** 所属文章ID */
  post_id: number;
  /** 元数据 key */
  meta_key: string;
  /** 元数据值 */
  meta_value: string;
}

// 批量更新文章状态请求
export interface BatchPostUpdateStatusRequest {
  /** 文章ID列表 */
  ids: number[];
  /** 状态 */
  status: "draft" | "published" | "private";
}

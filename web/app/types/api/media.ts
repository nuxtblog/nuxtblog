// 媒体分类（5 个内置类别，与后端 consts.BuiltinMediaCategories 保持同步）
// user   = 用户头像及个人封面（更新时会删除旧图）
// post   = 文章封面及正文内嵌图片
// doc    = 文档内嵌图片及附件
// moment = 动态附图/视频
// banner = 站点横幅 / Hero 图
export type MediaCategory = 'user' | 'post' | 'doc' | 'moment' | 'banner'

// 上传媒体请求
export interface UploadMediaRequest {
  alt_text?: string;
  title?: string;
  category?: MediaCategory;
}

// 更新媒体信息
export interface UpdateMediaRequest {
  alt_text?: string;
  title?: string;
}

// 查询媒体请求
export interface MediaQueryRequest {
  page?: number;
  size?: number;
  mime_type?: string;
  uploader_id?: number;
  category?: MediaCategory | '';
  storage_type?: number;
}

// 媒体响应 (matches backend MediaItem)
export interface MediaResponse {
  id: number;
  uploader_id: number;
  storage_type: number; // 1=local 2=S3 3=OSS 4=COS
  cdn_url: string;
  filename: string;
  mime_type: string;
  file_size: number;
  width?: number;
  height?: number;
  duration?: number;
  alt_text: string;
  title: string;
  category: MediaCategory;
  variants: string;
  created_at: string;
}

export interface MediaListResponse {
  list: MediaResponse[];
  total: number;
  page: number;
  size: number;
}

export interface MediaTypeStat {
  type: string;
  name: string;
  description: string;
  icon: string;
  count: number;
}

// 媒体分类 — 动态字符串，不再硬编码（内置: avatar, cover, post, doc, moment, banner）
export type MediaCategory = string

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
  storage_type: number; // 1=local 2=S3 3=OSS 4=COS 5=external
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

// ── Extension Groups & Format Policies ────────────────────────────────────────

export interface ExtensionGroup {
  name: string;
  label_zh: string;
  label_en: string;
  extensions: string[];
  max_size_mb: number;
}

export interface FormatPolicy {
  name: string;
  label_zh: string;
  label_en: string;
  is_system: boolean;
  groups: string[];
}

// types/post.ts
export type PostType = 1 | 2 | 3; // 1=post 2=page 3=custom
export type PostStatus = 1 | 2 | 3 | 4; // 1=draft 2=published 3=private 4=archived
export type CommentStatus = 0 | 1; // 0 =closed 1=open
export type OrderDirection = "asc" | "desc";
export type OrderBy =
  | "created_at"
  | "published_at"
  | "view_count"
  | "comment_count";

export interface CreatePostRequest {
  post_type?: PostType;
  title: string;
  slug?: string;
  content?: string;
  excerpt?: string;
  featured_img_id?: number;
  status?: PostStatus;
  published_at?: string;
  password?: string;
  comment_status?: CommentStatus;
  author_id: number;
  term_taxonomy_ids?: number[];
  metas?: Record<string, string>;
}

export interface UpdatePostRequest {
  post_type?: PostType;
  title?: string;
  slug?: string;
  content?: string;
  excerpt?: string;
  featured_img_id?: number;
  status?: PostStatus;
  published_at?: string;
  password?: string;
  comment_status?: CommentStatus;
  term_taxonomy_ids?: number[];
  metas?: Record<string, string>;
}

export interface PostQueryRequest {
  page?: number;
  page_size?: number;
  post_type?: string;
  status?: string;
  author_id?: number;
  term_taxonomy_id?: number;
  taxonomy?: string;
  search?: string;
  order_by?: OrderBy;
  order?: OrderDirection;
  start_date?: string;
  end_date?: string;
  meta_key?: string;
  meta_value?: string;
  sort_by?: string;
  published_after?: string;
  include_category_ids?: string;
  exclude_category_ids?: string;
}

export interface AuthorResponse {
  id: number;
  username: string;
  nickname: string;
  avatar?: string;
}

export interface MediaResponse {
  id: number;
  url: string;
  title?: string;
  mime_type?: string;
}

export interface TermResponse {
  id: number;
  term_id: number;
  name: string;
  slug: string;
  taxonomy: string;
  count: number;
  description: string;
  term_taxonomy_id: number;
}

export interface PostResponse {
  id: number;
  post_type: string;
  title: string;
  slug: string;
  content: string;
  excerpt: string;
  featured_img?: MediaResponse;
  status: string;
  published_at?: string;
  modified_at?: string;
  has_password: boolean;
  comment_status: string;
  comment_count: number;
  view_count: number;
  stats?: { view_count?: number; like_count?: number; comment_count?: number; share_count?: number };
  author?: AuthorResponse;
  terms?: TermResponse[];
  metas?: Record<string, string>;
  created_at: string;
  updated_at: string;
}

export interface PostListItemResponse {
  id: number;
  post_type: string;
  title: string;
  slug: string;
  excerpt: string;
  featured_img?: MediaResponse;
  status: string;
  published_at?: string;
  comment_count: number;
  view_count: number;
  author?: AuthorResponse;
  created_at: string;
  updated_at: string;
  metas?: Record<string, string>;
}

export interface PaginationResponse<T = any> {
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
  data: T[];
}

export interface VerifyPasswordRequest {
  password: string;
}

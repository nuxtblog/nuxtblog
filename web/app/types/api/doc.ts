export type DocStatus = 1 | 2 | 3  // 1=draft 2=published 3=archived
export type CollectionStatus = 1 | 2  // 1=draft 2=published

export interface DocCollectionItem {
  id: number
  slug: string
  title: string
  description: string
  cover_img_id: number | null
  author_id: number
  status: CollectionStatus
  locale: string
  sort_order: number
  created_at: string
  updated_at: string
  doc_count?: number
}

export interface DocItem {
  id: number
  collection_id: number
  parent_id: number | null
  sort_order: number
  status: DocStatus
  title: string
  slug: string
  excerpt: string
  author_id: number
  comment_status: number
  locale: string
  published_at: string | null
  created_at: string
  updated_at: string
  stats?: DocStatsItem
}

export interface DocDetailItem extends DocItem {
  content: string
  seo?: DocSeoItem
}

export interface DocStatsItem {
  view_count: number
  like_count: number
  comment_count: number
}

export interface DocSeoItem {
  meta_title: string
  meta_desc: string
  og_title: string
  og_image: string
  canonical_url: string
  robots: string
  structured_data: string
}

export interface DocRevisionItem {
  id: number
  doc_id: number
  author_id: number
  title: string
  rev_note: string
  created_at: string
}

// Request types
export interface CreateCollectionRequest {
  slug: string
  title: string
  description?: string
  cover_img_id?: number
  status: CollectionStatus
  locale?: string
  sort_order?: number
}

export interface UpdateCollectionRequest {
  slug?: string
  title?: string
  description?: string
  cover_img_id?: number
  status?: CollectionStatus
  locale?: string
  sort_order?: number
}

export interface CollectionQueryRequest {
  page?: number
  page_size?: number
  status?: number
  locale?: string
  author_id?: number
}

export interface CreateDocRequest {
  collection_id: number
  parent_id?: number
  title: string
  slug: string
  content?: string
  excerpt?: string
  status: DocStatus
  comment_status?: number
  locale?: string
  sort_order?: number
  published_at?: string
}

export interface UpdateDocRequest {
  collection_id?: number
  parent_id?: number
  title?: string
  slug?: string
  content?: string
  excerpt?: string
  status?: DocStatus
  comment_status?: number
  locale?: string
  sort_order?: number
  published_at?: string
}

export interface DocQueryRequest {
  page?: number
  page_size?: number
  collection_id?: number
  parent_id?: number
  status?: number
  author_id?: number
  locale?: string
  keyword?: string
}

export interface PaginatedCollections {
  data: DocCollectionItem[]
  total: number
  page: number
  page_size: number
  total_pages: number
}

export interface PaginatedDocs {
  data: DocItem[]
  total: number
  page: number
  page_size: number
  total_pages: number
}

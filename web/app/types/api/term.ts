// 创建分类/标签请求
export interface CreateTermRequest {
  /** 名称，必填，最大长度 100 */
  name: string;
  /** 别名，可选，最大长度 100 */
  slug?: string;
  /** 类型，必填，category 或 tag */
  taxonomy: "category" | "tag";
  /** 描述，可选，最大长度 255 */
  description?: string;
  /** 父级ID，可选 */
  parent_id?: number;
}

// 更新分类/标签请求
export interface UpdateTermRequest {
  /** 名称，可选 */
  name?: string;
  /** 别名，可选 */
  slug?: string;
  /** 描述，可选 */
  description?: string;
  /** 父级ID，可选 */
  parent_id?: number | null;
}

// 查询分类/标签请求
export interface TermQueryRequest {
  /** 类型，可选，category 或 tag */
  taxonomy?: "category" | "tag";
  /** 父级ID，可选 */
  parent_id?: number;
  /** 搜索关键字，可选 */
  search?: string;
}

// 分类/标签详情响应
export interface TermDetailResponse {
  /** 分类/标签ID */
  term_id: number;
  /** 名称 */
  name: string;
  /** 别名 */
  slug: string;
  /** 类型 category | tag */
  taxonomy: "category" | "tag" | string;
  /** 描述 */
  description: string;
  /** 父级ID，可选 */
  parent_id?: number;
  /** 关联文章数量 */
  count: number;
  /** 术语分类ID */
  term_taxonomy_id: number;
}

// 列表请求（分页）
export interface ListTermRequest {
  /** 类型，可选，category | tag */
  taxonomy?: "category" | "tag";
  /** 父级ID，可选 */
  parent_id?: number;
  /** 搜索关键字，可选 */
  search?: string;

  /** 页码，可选，默认 1 */
  page?: number;
  /** 每页数量，可选，默认 10 或 20 */
  page_size?: number;
}

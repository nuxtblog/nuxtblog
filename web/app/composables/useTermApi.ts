// composables/useTermApi.ts
import type {
  CreateTermRequest,
  UpdateTermRequest,
  TermDetailResponse,
} from "~/types/api/term";

// Backend taxonomy list item shape (before mapping)
interface TaxonomyApiItem {
  id: number;
  term_id: number;
  taxonomy: string;
  description: string;
  parent_id: number;
  post_count: number;
  extra?: string;
  term: {
    id: number;
    name: string;
    slug: string;
    created_at?: string;
  };
  children?: TaxonomyApiItem[];
}

interface TaxonomyListData {
  list: TaxonomyApiItem[];
  total: number;
}

function mapTaxonomyItem(item: TaxonomyApiItem): TermDetailResponse {
  return {
    term_id: item.term_id,
    name: item.term.name,
    slug: item.term.slug,
    taxonomy: item.taxonomy as "category" | "tag",
    description: item.description || "",
    parent_id: item.parent_id || undefined,
    count: item.post_count || 0,
    term_taxonomy_id: item.id,
  }
}

export const useTermApi = () => {
  const { apiFetch } = useApiFetch()

  /**
   * 创建术语（分类/标签）- two-step: POST /terms then POST /taxonomies
   */
  const createTerm = async (data: CreateTermRequest): Promise<TermDetailResponse> => {
    const slug = data.slug || data.name.toLowerCase().replace(/\s+/g, '-').replace(/[^\w-]+/g, '')
    // Step 1: create the term (name + slug)
    const termRes = await apiFetch<{ id: number }>('/terms', {
      method: 'POST',
      body: { name: data.name, slug },
    })
    // Step 2: create the taxonomy entry
    const taxRes = await apiFetch<{ id: number }>('/taxonomies', {
      method: 'POST',
      body: {
        term_id: termRes.id,
        taxonomy: data.taxonomy,
        description: data.description || '',
        parent_id: data.parent_id,
      },
    })
    return {
      term_id: termRes.id,
      term_taxonomy_id: taxRes.id,
      name: data.name,
      slug,
      taxonomy: data.taxonomy,
      description: data.description || '',
      parent_id: data.parent_id,
      count: 0,
    }
  };

  /**
   * 更新术语 - name/slug → PUT /terms/{termId}, description/parent → PUT /taxonomies/{termTaxonomyId}
   */
  const updateTerm = async (termId: number, termTaxonomyId: number, data: UpdateTermRequest) => {
    const promises: Promise<unknown>[] = []
    if (data.name !== undefined || data.slug !== undefined) {
      const body: Record<string, unknown> = {}
      if (data.name !== undefined) body.name = data.name
      if (data.slug !== undefined) body.slug = data.slug
      promises.push(apiFetch<void>(`/terms/${termId}`, { method: 'PUT', body }))
    }
    if (data.description !== undefined || data.parent_id !== undefined) {
      const body: Record<string, unknown> = {}
      if (data.description !== undefined) body.description = data.description
      if (data.parent_id !== undefined) body.parent_id = data.parent_id
      promises.push(apiFetch<void>(`/taxonomies/${termTaxonomyId}`, { method: 'PUT', body }))
    }
    await Promise.all(promises)
  };

  /**
   * 删除术语 - DELETE /terms/{termId} cascades to taxonomies + object bindings
   */
  const deleteTerm = async (termId: number) => {
    return await apiFetch<void>(`/terms/${termId}`, {
      method: 'DELETE',
    });
  };

  /**
   * 获取单个术语详情
   */
  const getTermById = async (termId: number) => {
    const item = await apiFetch<TaxonomyApiItem>(`/taxonomies/${termId}`);
    return mapTaxonomyItem(item)
  };

  /**
   * 获取术语列表
   */
  const getTerms = async (params?: {
    name?: string;
    taxonomy?: string;
    parent?: number;
  }) => {
    const data = await apiFetch<TaxonomyListData>('/taxonomies', {
      method: "GET",
      params,
    });
    return (data.list || []).map(mapTaxonomyItem);
  };

  return {
    createTerm,
    updateTerm,
    deleteTerm,
    getTermById,
    getTerms,
  };
};

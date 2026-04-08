import type {
  DocCollectionItem, DocItem, DocDetailItem, DocRevisionItem,
  CreateCollectionRequest, UpdateCollectionRequest, CollectionQueryRequest,
  CreateDocRequest, UpdateDocRequest, DocQueryRequest,
  PaginatedCollections, PaginatedDocs,
} from '~/types/api/doc'

export const useDocApi = () => {
  const { apiFetch } = useApiFetch()

  // Collections
  const getCollections = (params?: CollectionQueryRequest) =>
    apiFetch<PaginatedCollections>('/doc-collections', { params })

  const getCollection = (id: number) =>
    apiFetch<DocCollectionItem>(`/doc-collections/${id}`)

  const createCollection = (data: CreateCollectionRequest) =>
    apiFetch<{ id: number }>('/doc-collections', { method: 'POST', body: data })

  const updateCollection = (id: number, data: UpdateCollectionRequest) =>
    apiFetch<void>(`/doc-collections/${id}`, { method: 'PUT', body: data })

  const deleteCollection = (id: number) =>
    apiFetch<void>(`/doc-collections/${id}`, { method: 'DELETE' })

  // Docs
  const getDocs = (params?: DocQueryRequest) =>
    apiFetch<PaginatedDocs>('/docs', { params })

  const getDoc = (id: number) =>
    apiFetch<DocDetailItem>(`/docs/${id}`)

  const getDocBySlug = (slug: string) =>
    apiFetch<DocDetailItem>(`/docs/slug/${slug}`)

  const createDoc = (data: CreateDocRequest) =>
    apiFetch<{ id: number }>('/docs', { method: 'POST', body: data })

  const updateDoc = (id: number, data: UpdateDocRequest) =>
    apiFetch<void>(`/docs/${id}`, { method: 'PUT', body: data })

  const deleteDoc = (id: number) =>
    apiFetch<void>(`/docs/${id}`, { method: 'DELETE' })

  const getRevisions = (id: number, page = 1, size = 10) =>
    apiFetch<{ list: DocRevisionItem[]; total: number }>(`/docs/${id}/revisions`, { params: { page, size } })

  const restoreRevision = (docId: number, revisionId: number) =>
    apiFetch<void>(`/docs/${docId}/revisions/${revisionId}/restore`, { method: 'POST' })

  const incrementView = (id: number) =>
    apiFetch<void>(`/docs/${id}/view`, { method: 'POST' })

  return {
    getCollections, getCollection, createCollection, updateCollection, deleteCollection,
    getDocs, getDoc, getDocBySlug, createDoc, updateDoc, deleteDoc,
    getRevisions, restoreRevision, incrementView,
  }
}

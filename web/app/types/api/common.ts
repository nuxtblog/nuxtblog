export interface PaginationResponse<T> {
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
  data: T;
}

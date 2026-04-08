import type {
  UploadMediaRequest,
  UpdateMediaRequest,
  MediaQueryRequest,
  MediaResponse,
  MediaListResponse,
} from "~/types/api/media";

export const useMediaApi = () => {
  const { apiFetch } = useApiFetch();

  const upload = async (
    file: File,
    data?: UploadMediaRequest
  ): Promise<MediaResponse> => {
    const formData = new FormData();
    formData.append("file", file);
    if (data?.alt_text) formData.append("alt_text", data.alt_text);
    if (data?.title) formData.append("title", data.title);
    if (data?.category) formData.append("category", data.category);

    return apiFetch<MediaResponse>("/medias/upload", {
      method: "POST",
      body: formData,
    });
  };

  const list = async (query?: MediaQueryRequest): Promise<MediaListResponse> => {
    return apiFetch<MediaListResponse>("/medias", {
      method: "GET",
      params: {
        page: query?.page ?? 1,
        size: query?.size ?? 20,
        ...(query?.mime_type && { mime_type: query.mime_type }),
        ...(query?.uploader_id && { uploader_id: query.uploader_id }),
        ...(query?.category && { category: query.category }),
      },
    });
  };

  const getByID = async (id: number): Promise<MediaResponse> => {
    return apiFetch<MediaResponse>(`/medias/${id}`);
  };

  const update = async (
    id: number,
    data: UpdateMediaRequest
  ): Promise<void> => {
    return apiFetch<void>(`/medias/${id}`, {
      method: "PUT",
      body: data,
    });
  };

  const link = async (data: {
    url: string;
    title?: string;
    alt_text?: string;
    category?: MediaCategory;
  }): Promise<MediaResponse> => {
    return apiFetch<MediaResponse>("/medias/link", {
      method: "POST",
      body: data,
    });
  };

  const delete_ = async (id: number): Promise<void> => {
    return apiFetch<void>(`/medias/${id}`, {
      method: "DELETE",
    });
  };

  const localize = async (id: number): Promise<MediaResponse> => {
    return apiFetch<MediaResponse>(`/medias/${id}/localize`, {
      method: "POST",
    });
  };

  return {
    upload,
    link,
    list,
    getByID,
    update,
    delete: delete_,
    localize,
  };
};

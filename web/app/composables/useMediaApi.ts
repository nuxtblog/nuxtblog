import type {
  MediaCategory,
  UploadMediaRequest,
  UpdateMediaRequest,
  MediaQueryRequest,
  MediaResponse,
  MediaListResponse,
  ExtensionGroup,
  FormatPolicy,
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
        ...(query?.storage_type != null && { storage_type: query.storage_type }),
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

  const getExtensionGroups = async (): Promise<{ list: ExtensionGroup[] }> => {
    return apiFetch<{ list: ExtensionGroup[] }>("/admin/media/extension-groups");
  };

  const getFormatPolicies = async (): Promise<{ list: FormatPolicy[] }> => {
    return apiFetch<{ list: FormatPolicy[] }>("/admin/media/format-policies");
  };

  const saveExtensionGroups = async (groups: ExtensionGroup[]): Promise<void> => {
    return apiFetch<void>("/admin/media/extension-groups", {
      method: "PUT",
      body: { groups },
    });
  };

  const createFormatPolicy = async (data: Omit<FormatPolicy, "is_system">): Promise<void> => {
    return apiFetch<void>("/admin/media/format-policies", {
      method: "POST",
      body: data,
    });
  };

  const updateFormatPolicy = async (name: string, data: Partial<FormatPolicy>): Promise<void> => {
    return apiFetch<void>(`/admin/media/format-policies/${name}`, {
      method: "PUT",
      body: data,
    });
  };

  const deleteFormatPolicy = async (name: string): Promise<void> => {
    return apiFetch<void>(`/admin/media/format-policies/${name}`, {
      method: "DELETE",
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
    getExtensionGroups,
    getFormatPolicies,
    saveExtensionGroups,
    createFormatPolicy,
    updateFormatPolicy,
    deleteFormatPolicy,
  };
};

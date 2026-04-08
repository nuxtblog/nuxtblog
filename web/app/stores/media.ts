import type {
  MediaResponse,
  MediaQueryRequest,
  UpdateMediaRequest,
  UploadMediaRequest,
} from "~/types/api/media";

export const useMediaStore = defineStore("media", () => {
  const medias = ref<MediaResponse[]>([]);
  const currentMedia = ref<MediaResponse | null>(null);
  const loading = ref(false);
  const uploading = ref(false);
  const error = ref<string | null>(null);
  const pagination = ref({
    page: 1,
    size: 20,
    total: 0,
  });

  const mediaApi = useMediaApi();

  const fetchMedias = async (query?: MediaQueryRequest) => {
    loading.value = true;
    error.value = null;
    try {
      const result = await mediaApi.list(query);
      medias.value = result.list;
      pagination.value = {
        page: result.page,
        size: result.size,
        total: result.total,
      };
    } catch (err) {
      error.value = err instanceof Error ? err.message : "获取媒体列表失败";
      console.error("Failed to fetch medias:", err);
    } finally {
      loading.value = false;
    }
  };

  const fetchMediaDetail = async (id: number) => {
    error.value = null;
    try {
      currentMedia.value = await mediaApi.getByID(id);
    } catch (err) {
      error.value = err instanceof Error ? err.message : "获取媒体详情失败";
      console.error("Failed to fetch media detail:", err);
    }
  };

  const uploadMedia = async (
    file: File,
    data?: UploadMediaRequest
  ): Promise<MediaResponse | null> => {
    uploading.value = true;
    error.value = null;
    try {
      const media = await mediaApi.upload(file, data);
      medias.value.unshift(media);
      pagination.value.total += 1;
      return media;
    } catch (err) {
      error.value = err instanceof Error ? err.message : "上传媒体失败";
      console.error("Failed to upload media:", err);
      return null;
    } finally {
      uploading.value = false;
    }
  };

  const updateMedia = async (
    id: number,
    data: UpdateMediaRequest
  ): Promise<void> => {
    error.value = null;
    await mediaApi.update(id, data);
    const updated = await mediaApi.getByID(id);
    const index = medias.value.findIndex((m) => m.id === id);
    if (index > -1) medias.value[index] = updated;
    if (currentMedia.value?.id === id) currentMedia.value = updated;
  };

  const deleteMedia = async (id: number): Promise<void> => {
    error.value = null;
    await mediaApi.delete(id);
    medias.value = medias.value.filter((m) => m.id !== id);
    if (currentMedia.value?.id === id) currentMedia.value = null;
    pagination.value.total -= 1;
  };

  const batchDeleteMedias = async (ids: number[]): Promise<void> => {
    error.value = null;
    await Promise.all(ids.map((id) => mediaApi.delete(id)));
    medias.value = medias.value.filter((m) => !ids.includes(m.id));
    if (currentMedia.value && ids.includes(currentMedia.value.id)) {
      currentMedia.value = null;
    }
    pagination.value.total -= ids.length;
  };

  const clearCurrentMedia = () => { currentMedia.value = null; };
  const clearError = () => { error.value = null; };

  return {
    medias,
    currentMedia,
    loading,
    uploading,
    error,
    pagination,

    fetchMedias,
    fetchMediaDetail,
    uploadMedia,
    updateMedia,
    deleteMedia,
    batchDeleteMedias,
    clearCurrentMedia,
    clearError,
  };
});

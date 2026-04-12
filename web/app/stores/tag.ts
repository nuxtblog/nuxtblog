import type { TermDetailResponse } from "~/types/api/term";

export const useTagStore = defineStore("tag", () => {
  const tags = ref<TermDetailResponse[]>([]);
  const loading = ref(false);
  const error = ref<string | null>(null);

  const termApi = useTermApi();

  /**
   * 加载全部标签
   */
  const loadTags = async () => {
    loading.value = true;
    error.value = null;

    try {
      const terms = await termApi.getTerms({ taxonomy: "tag" });

      // 排序：按名称升序
      tags.value = terms.sort((a, b) => a.name.localeCompare(b.name));

      return tags.value;
    } catch (err) {
      console.error("加载标签失败:", err);
      error.value = "加载标签失败";
      return [];
    } finally {
      loading.value = false;
    }
  };

  /**
   * 根据 ID 获取标签（优先本地，没有再请求）
   */
  const getTagById = async (id: number): Promise<TermDetailResponse | null> => {
    const local = tags.value.find((t) => t.term_id === id);
    if (local) return local;

    try {
      const tag = await termApi.getTermById(id);
      tags.value.push(tag);
      return tag;
    } catch (err) {
      console.error("获取标签失败:", err);
      return null;
    }
  };

  /**
   * 新增标签
   */
  const addNewTag = async (data: { name: string; slug?: string; description?: string }): Promise<TermDetailResponse> => {
    const trimmed = data.name.trim();
    if (!trimmed) {
      throw new Error("标签名称不能为空");
    }

    // 检查是否已存在
    const existing = tags.value.find(
      (t) => t.name.toLowerCase() === trimmed.toLowerCase()
    );
    if (existing) return existing;

    const raw = data.slug || trimmed.toLowerCase().replace(/\s+/g, "-").replace(/[^a-z0-9-]+/g, "");
    const slug = raw || `tag-${Date.now().toString(36)}`;

    try {
      const createdTag = await termApi.createTerm({
        name: trimmed,
        slug,
        taxonomy: "tag",
        description: data.description,
      });

      tags.value.push(createdTag);
      tags.value.sort((a, b) => a.name.localeCompare(b.name));

      return createdTag;
    } catch (err) {
      console.error("创建标签失败:", err);
      throw err;
    }
  };

  /**
   * 更新标签
   */
  const updateTag = async (
    tagId: number,
    data: { name?: string; slug?: string; description?: string }
  ) => {
    try {
      const tag = tags.value.find((t) => t.term_id === tagId);
      if (!tag) throw new Error("Tag not found");
      await termApi.updateTerm(tagId, tag.term_taxonomy_id, data);

      const index = tags.value.findIndex((t) => t.term_id === tagId);
      if (index !== -1) {
        tags.value[index] = { ...tags.value[index], ...data };
        tags.value.sort((a, b) => a.name.localeCompare(b.name));
      }

      return tags.value[index];
    } catch (err) {
      console.error("更新标签失败:", err);
      throw err;
    }
  };

  /**
   * 删除标签
   */
  const deleteTag = async (tagId: number) => {
    try {
      await termApi.deleteTerm(tagId);
      tags.value = tags.value.filter((t) => t.term_id !== tagId);
    } catch (err) {
      console.error("删除标签失败:", err);
      throw err;
    }
  };

  /**
   * 是否已加载
   */
  const isLoaded = computed(() => tags.value.length > 0);

  return {
    tags,
    loading,
    error,
    isLoaded,

    loadTags,
    getTagById,
    addNewTag,
    updateTag,
    deleteTag,
  };
});

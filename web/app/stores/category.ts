import { ref, computed } from "vue";
import type { TermDetailResponse, UpdateTermRequest } from "~/types/api/term";
import type { TermWithLevel } from "~/types/models/term";

export const useCategoryStore = defineStore("category", () => {
  const categories = ref<TermDetailResponse[]>([]); // 平铺列表
  const loading = ref(false);
  const error = ref<string | null>(null);

  const termApi = useTermApi();

  /**
   * 加载分类列表
   */
  const loadCategories = async () => {
    loading.value = true;
    error.value = null;

    try {
      const list = await termApi.getTerms({ taxonomy: "category" });
      categories.value = list.sort((a, b) => a.name.localeCompare(b.name));
      return categories.value;
    } catch (err) {
      console.error("加载分类失败:", err);
      error.value = "加载分类失败";
      return [];
    } finally {
      loading.value = false;
    }
  };

  /**
   * 根据 ID 获取分类（优先本地）
   */
  const getCategoryById = async (
    id: number,
  ): Promise<TermDetailResponse | null> => {
    const local = categories.value.find((c) => c.term_id === id);
    if (local) return local;

    try {
      const category = await termApi.getTermById(id);
      categories.value.push(category);
      categories.value.sort((a, b) => a.name.localeCompare(b.name));
      return category;
    } catch (err) {
      console.error("获取分类失败:", err);
      return null;
    }
  };

  /**
   * 新增分类
   */
  const addNewCategory = async (data: {
    name: string;
    slug?: string;
    description?: string;
    parent_id?: number;
  }) => {
    const trimmed = data.name.trim();
    if (!trimmed) throw new Error("分类名称不能为空");

    // 是否已存在同名分类（同层级）
    const existing = categories.value.find(
      (c) =>
        c.name.toLowerCase() === trimmed.toLowerCase() &&
        c.parent_id === data.parent_id,
    );
    if (existing) return existing;

    const slug = data.slug || trimmed.toLowerCase().replace(/\s+/g, "-");

    try {
      const created = await termApi.createTerm({
        name: trimmed,
        slug,
        taxonomy: "category",
        description: data.description,
        parent_id: data.parent_id,
      });
      categories.value.push(created);
      categories.value.sort((a, b) => a.name.localeCompare(b.name));
      console.log("创建分类成功");

      return created;
    } catch (err) {
      console.error("创建分类失败:", err);
      throw err;
    }
  };

  /**
   * 更新分类
   */
  const updateCategory = async (
    categoryId: number,
    data: UpdateTermRequest,
  ) => {
    try {
      const category = categories.value.find((c) => c.term_id === categoryId);
      if (!category) throw new Error("Category not found");
      await termApi.updateTerm(categoryId, category.term_taxonomy_id, data);
      const index = categories.value.findIndex((c) => c.term_id === categoryId);
      if (index !== -1) {
        categories.value[index] = {
          ...categories.value[index],
          ...data,
          parent_id:
            data.parent_id !== undefined
              ? (data.parent_id ?? undefined)
              : categories.value[index].parent_id,
        };
      }
      categories.value.sort((a, b) => a.name.localeCompare(b.name));
      console.log("更新分类成功");
      return categories.value[index];
    } catch (err) {
      console.error("更新分类失败:", err);
      throw err;
    }
  };
  /**
   * 更新父级ID
   */
  const updateCategoryParent = async (
    categoryId: number,
    newParentId: number | undefined | null,
  ) => {
    try {
      const category = await getCategoryById(categoryId);
      if (!category) {
        throw new Error(`Category ${categoryId} not found`);
      }

      // 防止把分类移动到自身下面，但 null 是合法的（表示顶级）
      if (newParentId !== null && categoryId === newParentId) {
        throw new Error("Category cannot be its own parent");
      }

      // 无需更新
      if (category.parent_id === newParentId) {
        console.log("Parent ID is already the target value");
        return;
      }

      await updateCategory(categoryId, { parent_id: newParentId });
    } catch (err) {
      console.error("Failed to update category parent:", err);
      throw err;
    }
  };

  /**
   * 删除分类
   */
  const deleteCategory = async (categoryId: number) => {
    try {
      await termApi.deleteTerm(categoryId);
      categories.value = categories.value.filter(
        (c) => c.term_id !== categoryId,
      );
    } catch (err) {
      console.error("删除分类失败:", err);
      throw err;
    }
  };

  /**
   * 树结构（computed，从平铺列表生成）
   */
  const categoryTree = computed(() => {
    const map = new Map<
      number,
      TermDetailResponse & { children: TermDetailResponse[] }
    >();
    const roots: (TermDetailResponse & { children: TermDetailResponse[] })[] =
      [];

    // 初始化 map
    categories.value.forEach((cat) => {
      map.set(cat.term_id, { ...cat, children: [] });
    });

    // 构建树
    map.forEach((cat) => {
      if (cat.parent_id && map.has(cat.parent_id)) {
        map.get(cat.parent_id)!.children.push(cat);
      } else {
        roots.push(cat);
      }
    });

    return roots;
  });

  /**
   * 是否已加载
   */
  const isLoaded = computed(() => categories.value.length > 0);

  /**
   * 获取父类名称
   */
  const getParentName = (category: TermDetailResponse) => {
    if (!category.parent_id) return null;
    const parent = categories.value.find(
      (c) => c.term_id === category.parent_id,
    );
    return parent?.name ?? null;
  };
  /**
   * 获取子类数量
   */
  const getChildrenCount = (term_id: number) => {
    const collectChildren = (parentId: number): number => {
      const directChildren = categories.value.filter(
        (c) => c.parent_id === parentId,
      );
      let count = directChildren.length;
      for (const child of directChildren) {
        count += collectChildren(child.term_id);
      }
      return count;
    };

    return collectChildren(term_id);
  };

  /**
   * 获取平铺的分类列表，用于选择器
   * @param excludeId 要排除的分类 ID（用于编辑时排除自身及其子分类）
   * @returns 平铺后的分类列表，包含层级信息
   */
  const getFlattenedParents = (excludeId?: number): TermWithLevel[] => {
    const flatList: TermWithLevel[] = [];

    if (!categoryTree.value || categoryTree.value.length === 0) {
      return flatList;
    }

    // 标记要排除的 ID 及其所有子分类
    const excludeIds = new Set<number>();
    const markExcludeIds = (
      items: (TermDetailResponse & { children?: TermDetailResponse[] })[],
    ) => {
      items.forEach((item) => {
        if (item.term_id === excludeId) {
          excludeIds.add(item.term_id);
          if (item.children) {
            markExcludeIds(item.children);
          }
        } else if (item.children) {
          markExcludeIds(item.children);
        }
      });
    };

    if (excludeId) {
      markExcludeIds(categoryTree.value);
    }

    // 递归平铺树结构
    const addToList = (
      items: (TermDetailResponse & { children?: TermDetailResponse[] })[],
      level = 0,
    ) => {
      items.forEach((item) => {
        if (!excludeIds.has(item.term_id)) {
          flatList.push({ ...item, level } as TermWithLevel);
          if (item.children && item.children.length > 0) {
            addToList(item.children, level + 1);
          }
        }
      });
    };

    addToList(categoryTree.value);

    return flatList;
  };

  // 加载的时候要load一下
  loadCategories();

  return {
    categories,
    categoryTree,
    loading,
    error,
    isLoaded,
    loadCategories,
    getCategoryById,
    addNewCategory,
    updateCategory,
    deleteCategory,
    updateCategoryParent,

    getParentName,
    getChildrenCount,
    getFlattenedParents,
  };
});

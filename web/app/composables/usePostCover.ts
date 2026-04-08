const STATIC_FALLBACK = "/images/default-cover.svg";

/**
 * Returns cover image helpers sourced from the options store.
 *
 * - `defaultCover` — used as `src` when a post has no featured image
 * - `errorCover`   — used in `@error` when the image URL is broken / 404
 * - `onImgError`   — ready-made event handler for `@error`
 */
export const usePostCover = () => {
  const store = useOptionsStore();

  const defaultCover = computed(() =>
    store.get("default_post_cover", STATIC_FALLBACK),
  );

  const errorCover = computed(() =>
    store.get("error_post_cover", STATIC_FALLBACK),
  );

  const onImgError = (e: string | Event) => {
    if (typeof e === "string") return;
    const img = e.target as HTMLImageElement;

    if (!img) return;

    // 防止无限循环
    if (img.dataset.errorHandled) return;

    img.dataset.errorHandled = "true";
    img.src = errorCover.value || STATIC_FALLBACK;
  };

  return { defaultCover, errorCover, onImgError };
};

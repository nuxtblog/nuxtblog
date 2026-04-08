/**
 * Ensures a loading ref stays true for at least `minMs` milliseconds.
 * Prevents skeleton flickering when data loads too fast.
 */
export const useMinLoading = (
  loadingRef: Ref<boolean>,
  minMs = 300,
): Readonly<Ref<boolean>> => {
  const display = ref(loadingRef.value);
  let timer: ReturnType<typeof setTimeout> | null = null;
  let startTime = loadingRef.value ? Date.now() : 0;

  watch(loadingRef, (val) => {
    if (val) {
      if (timer) {
        clearTimeout(timer);
        timer = null;
      }
      startTime = Date.now();
      display.value = true;
    } else {
      const remaining = minMs - (Date.now() - startTime);
      if (remaining > 0) {
        timer = setTimeout(() => {
          display.value = false;
        }, remaining);
      } else {
        display.value = false;
      }
    }
  });

  onUnmounted(() => {
    if (timer) clearTimeout(timer);
  });

  return readonly(display);
};

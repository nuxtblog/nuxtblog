const WIDTH_MAP: Record<string, string> = {
  '5xl': 'max-w-5xl',
  '6xl': 'max-w-6xl',
  '7xl': 'max-w-7xl',
}

export const useContainerWidth = () => {
  const optionsStore = useOptionsStore()
  const containerClass = computed(
    () => WIDTH_MAP[optionsStore.get('site_container_width', '7xl')] ?? 'max-w-7xl',
  )
  return { containerClass }
}

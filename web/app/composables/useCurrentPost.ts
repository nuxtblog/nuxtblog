export interface TocItem {
  id: string
  text: string
  level: number
}

export interface DownloadItem {
  name: string
  url: string
  size?: string
  desc?: string
}

/**
 * Shared state for the currently viewed post.
 * Populated by posts/[slug].vue, consumed by sidebar widgets (Toc, Downloads).
 */
export const useCurrentPost = () => {
  const toc = useState<TocItem[]>('post-toc', () => [])
  const downloads = useState<DownloadItem[]>('post-downloads', () => [])

  const clearPost = () => {
    toc.value = []
    downloads.value = []
  }

  return { toc, downloads, clearPost }
}

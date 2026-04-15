import type { InjectionKey, Ref } from 'vue'

export interface SidebarContext {
  /** 页面类型 */
  pageType: 'homepage' | 'blog' | 'post' | 'page' | 'other'
  /** 当前文章/页面摘要（仅 post/page 页面有值） */
  post?: {
    id: number
    slug: string
    title: string
    excerpt?: string
    featuredImg?: string
    author?: { id: number; nickname: string; avatar?: string }
    terms?: Array<{ id: number; name: string; taxonomy: string; slug: string }>
    publishedAt?: string
    viewCount?: number
    commentCount?: number
    likeCount?: number
  }
}

export const SIDEBAR_CONTEXT_KEY: InjectionKey<Ref<SidebarContext>> = Symbol('sidebarContext')

/** 供内置小部件使用的 composable */
export function useSidebarContext(): Ref<SidebarContext> {
  return inject(SIDEBAR_CONTEXT_KEY, ref({ pageType: 'other' }))
}

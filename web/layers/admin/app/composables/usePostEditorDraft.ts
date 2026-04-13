import type { CreatePostRequest } from "~/types/api/post";
import { useEditorDraft } from './useEditorDraft'

export function usePostEditorDraft(
  formData: Ref<CreatePostRequest>,
  opts: { mode: "create" | "edit"; postId?: number; hasInitialContent: boolean },
) {
  return useEditorDraft(formData, {
    ...opts,
    keyPrefix: 'blog:draft',
    entityId: opts.postId,
  })
}

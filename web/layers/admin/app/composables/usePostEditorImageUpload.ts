import type { CreatePostRequest } from "~/types/api/post";
import { useEditorImageUpload } from './useEditorImageUpload'

export function usePostEditorImageUpload(
  formData: Ref<CreatePostRequest>,
  editorRef: Ref<any>,
) {
  return useEditorImageUpload(formData, editorRef, { imageCategory: 'post' })
}

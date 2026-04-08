import type { PostCard } from "~/types/models/post";
import type { PostListItemResponse } from "~/types/api/post";

const FALLBACK_COVER = '/images/default-cover.svg'

export function transformPostList(
  posts: PostListItemResponse[],
  defaultCover = FALLBACK_COVER
): PostCard[] {
  return posts.map((post) => ({
    id: post.id,
    slug: post.slug,
    title: post.title,
    excerpt: post.excerpt,
    cover: post.featured_img?.url || defaultCover,
    author: post.author,
    date: post.updated_at,
    comments: post.comment_count,
    views: post.view_count,

    category: "",
    authorAvatar: post.author?.avatar || "",
    authorName: post.author?.nickname || "Anonymous",
  }));
}

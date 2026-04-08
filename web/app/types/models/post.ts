export interface PostCard {
  id: number;
  slug: string;
  title: string;
  excerpt: string;
  cover?: string;
  category: string;
  date: string;
  views: number;
  comments: number;

  authorName: string;
  authorAvatar: string;
}

<script setup lang="ts">
const route = useRoute();
const id = Number(route.params.id);

if (!id || isNaN(id)) {
  throw createError({ statusCode: 404, message: "文章不存在" });
}

const { apiFetch } = useApiFetch();

const post = await apiFetch<{ slug: string; post_type: number }>(
  `/posts/${id}`,
  { params: { fields: "slug,post_type" } }
).catch(() => null);

if (!post?.slug) {
  throw createError({ statusCode: 404, message: "文章不存在" });
}

const target = post.post_type === 2 ? `/pages/${post.slug}` : `/posts/${post.slug}`;

await navigateTo(target, { redirectCode: 301 });
</script>

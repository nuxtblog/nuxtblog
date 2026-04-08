// middleware/admin-redirect.ts
import type { RouteLocationNormalized } from "vue-router";

export default defineNuxtRouteMiddleware((to: RouteLocationNormalized) => {
  // 如果访问 /admin 就直接跳转到 /admin/dashboard
  if (to.path === "/admin") {
    return navigateTo("/admin/dashboard");
  }
});

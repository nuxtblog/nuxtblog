/**
 * 管理后台侧边栏状态
 * - desktop: 侧边栏内联显示，可折叠为图标模式
 * - mobile: 侧边栏作为抽屉层叠覆盖，点击遮罩关闭
 */
const mobileOpen = ref(false);
const desktopCollapsed = ref(false);

export function useAdminSidebar() {
  const toggleMobile = () => {
    mobileOpen.value = !mobileOpen.value;
  };

  const closeMobile = () => {
    mobileOpen.value = false;
  };

  const toggleDesktop = () => {
    desktopCollapsed.value = !desktopCollapsed.value;
  };

  return {
    mobileOpen: readonly(mobileOpen),
    desktopCollapsed: readonly(desktopCollapsed),
    toggleMobile,
    closeMobile,
    toggleDesktop,
  };
}

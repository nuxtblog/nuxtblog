/**
 * Contribution Slot Registry
 *
 * Central registry of all named slots where plugins can inject UI.
 * Using these constants instead of raw strings gives compile-time safety.
 */
export const CONTRIBUTION_SLOTS = {
  // ── Admin ───────────────────────────────────────────────────────────────
  ADMIN_SIDEBAR_NAV: 'admin:sidebar-nav',
  ADMIN_DASHBOARD_STATS: 'admin:dashboard-stats',
  ADMIN_DASHBOARD_BOTTOM: 'admin:dashboard-bottom',
  ADMIN_POST_EDITOR_FOOTER: 'admin:post-editor-footer',

  // ── Post Editor ─────────────────────────────────────────────────────────
  POST_EDITOR_TOOLBAR: 'post-editor:toolbar',
  POST_EDITOR_CONTEXT: 'post-editor:context',

  // ── Post List ───────────────────────────────────────────────────────────
  POST_LIST_ROW_ACTION: 'post-list:row-action',

  // ── Public: Layout ──────────────────────────────────────────────────────
  PUBLIC_HEADER_ACTIONS: 'public:header-actions',
  PUBLIC_FOOTER_EXTRA: 'public:footer-extra',
  PUBLIC_SIDEBAR_TOP: 'public:sidebar-top',
  PUBLIC_SIDEBAR_BOTTOM: 'public:sidebar-bottom',
  PUBLIC_SIDEBAR_WIDGET: 'public:sidebar-widget',

  // ── Public: Home ────────────────────────────────────────────────────────
  PUBLIC_HOME_TOP: 'public:home-top',
  PUBLIC_HOME_BOTTOM: 'public:home-bottom',

  // ── Public: Posts list ──────────────────────────────────────────────────
  PUBLIC_POSTS_TOP: 'public:posts-top',
  PUBLIC_POSTS_BOTTOM: 'public:posts-bottom',

  // ── Public: Single Post ─────────────────────────────────────────────────
  PUBLIC_POST_BEFORE_TITLE: 'public:post-before-title',
  PUBLIC_POST_AFTER_TITLE: 'public:post-after-title',
  PUBLIC_POST_BEFORE_CONTENT: 'public:post-before-content',
  PUBLIC_POST_AFTER_CONTENT: 'public:post-after-content',
  PUBLIC_POST_BEFORE_COMMENTS: 'public:post-before-comments',
  PUBLIC_POST_AFTER_COMMENTS: 'public:post-after-comments',

  // ── Public: Categories ──────────────────────────────────────────────────
  PUBLIC_CATEGORIES_TOP: 'public:categories-top',
  PUBLIC_CATEGORIES_BOTTOM: 'public:categories-bottom',
  PUBLIC_CATEGORY_TOP: 'public:category-top',
  PUBLIC_CATEGORY_BOTTOM: 'public:category-bottom',

  // ── Public: Tags ────────────────────────────────────────────────────────
  PUBLIC_TAGS_TOP: 'public:tags-top',
  PUBLIC_TAGS_BOTTOM: 'public:tags-bottom',

  // ── Public: Archive ─────────────────────────────────────────────────────
  PUBLIC_ARCHIVE_TOP: 'public:archive-top',
  PUBLIC_ARCHIVE_BOTTOM: 'public:archive-bottom',

  // ── Public: Docs ────────────────────────────────────────────────────────
  PUBLIC_DOCS_TOP: 'public:docs-top',
  PUBLIC_DOCS_BOTTOM: 'public:docs-bottom',

  // ── Public: Pages ───────────────────────────────────────────────────────
  PUBLIC_PAGE_BEFORE_CONTENT: 'public:page-before-content',
  PUBLIC_PAGE_AFTER_CONTENT: 'public:page-after-content',

  // ── Public: User Profile ────────────────────────────────────────────────
  PUBLIC_USER_PROFILE_TOP: 'public:user-profile-top',
  PUBLIC_USER_PROFILE_BOTTOM: 'public:user-profile-bottom',

  // ── Public: User Menu ────────────────────────────────────────────────
  PUBLIC_USER_MENU: 'public:user-menu',

  // ── Public: Floating Toolbar ──────────────────────────────────────────
  PUBLIC_FLOATING_TOOLBAR: 'public:floating-toolbar',

  // ── Public: Post Actions ──────────────────────────────────────────────
  PUBLIC_POST_ACTIONS: 'public:post-actions',
} as const

export type ContributionSlotName = typeof CONTRIBUTION_SLOTS[keyof typeof CONTRIBUTION_SLOTS]

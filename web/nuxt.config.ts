// https://nuxt.com/docs/api/configuration/nuxt-config

export default defineNuxtConfig({
  compatibilityDate: "2025-11-15",
  devtools: { enabled: false },
  css: ["~/assets/styles/main.css"],

  app: {
    head: {
      link: [{ rel: "icon", type: "image/svg+xml", href: "/icon.svg" }],
    },
  },

  runtimeConfig: {
    public: {
      apiBase:
        process.env.NUXT_PUBLIC_API_BASE || "http://localhost:9000/api/v1",
    },
  },

  vite: {
    optimizeDeps: {
      include: [
        "@nuxt/ui > prosemirror-state",
        "@nuxt/ui > prosemirror-transform",
        "@nuxt/ui > prosemirror-model",
        "@nuxt/ui > prosemirror-view",
        "@nuxt/ui > prosemirror-gapcursor",
        "@tiptap/extension-emoji",
        "@tiptap/extension-text-align",
        "dayjs",
        "dayjs/plugin/relativeTime",
        "dayjs/locale/zh-cn",
        "dayjs/locale/en",
        "zod",
        "tiptap-markdown",
        "vue-draggable-plus",
        "marked",
        "highlight.js",
        "dompurify",
      ],
    },
    server: {
      proxy: {
        // 只代理业务 API，避免拦截 Nuxt 内部路由（如 /api/_nuxt_icon/）
        "/api/v1": {
          target: "http://localhost:9000",
          changeOrigin: true,
        },
        // 插件自定义路由（/api/plugin/xxx/...）
        "/api/plugin": {
          target: "http://localhost:9000",
          changeOrigin: true,
        },
        // 插件共享依赖 shim（/_shared/vue.mjs 等，用于 public.mjs）
        "/_shared": {
          target: "http://localhost:9000",
          changeOrigin: true,
        },
        // 插件静态资源（admin.mjs, public.js 等）
        "/api/plugins": {
          target: "http://localhost:9000",
          changeOrigin: true,
        },
        "/uploads": {
          target: "http://localhost:9000",
          changeOrigin: true,
        },
      },
    },
  },

  extends: ["./layers/admin"],
  imports: {
    dirs: ["modules/editor/composables"],
  },

  modules: [
    "@nuxtjs/i18n",
    "@nuxt/ui",
    "@pinia/nuxt",
    "@nuxt/image",
    "@nuxtjs/robots",
    "@nuxtjs/sitemap",
    "@vee-validate/nuxt",
    "@vueuse/nuxt",
  ],

  i18n: {
    locales: [
      { code: "zh", name: "中文", file: "zh.ts" },
      { code: "en", name: "English", file: "en.ts" },
    ],
    langDir: "locales/",
    defaultLocale: "zh",
    strategy: "no_prefix",
    detectBrowserLanguage: {
      useCookie: true,
      cookieKey: "i18n_locale",
      redirectOn: "root",
    },
    experimental: {
      typedOptionsAndMessages: "default",
    },
  },

  // 使用本地图标包；客户端 scan 打包主应用用到的图标，
  // 插件运行时用到的图标通过 Nuxt 内置 /api/_nuxt_icon/ 按需加载
  icon: {
    serverBundle: {
      collections: ["tabler"],
    },
    clientBundle: {
      scan: true,
    },
  },

  // 禁用 Google Fonts（网络受限环境）
  fonts: {
    providers: {
      google: false,
      googleicons: false,
    },
  },
});

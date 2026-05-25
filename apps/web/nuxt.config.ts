// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  devtools: { enabled: true },
  compatibilityDate: '2025-05-21',
  typescript: { strict: true },

  /** HTML5 History — 禁止 hash 路由（#） */
  router: {
    options: {
      hashMode: false,
    },
  },

  build: {
    transpile: ['@crm/ui-kit'],
  },

  app: {
    head: {
      script: [
        {
          key: 'theme-init',
          innerHTML:
            "(function(){try{var m=document.cookie.match(/crm-theme=([^;]+)/);document.documentElement.setAttribute('data-theme',m&&m[1]==='v2'?'v2':'v1')}catch(e){document.documentElement.setAttribute('data-theme','v1')}})();",
          type: 'text/javascript',
          tagPosition: 'head',
        },
      ],
      link: [
        { rel: 'preconnect', href: 'https://fonts.googleapis.com' },
        { rel: 'preconnect', href: 'https://fonts.gstatic.com', crossorigin: '' },
        {
          rel: 'stylesheet',
          href: 'https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@400;500;600;700&display=swap',
        },
      ],
    },
  },

  /** @nuxt/ui 会自行 install @nuxtjs/tailwindcss，勿重复注册否则 Tailwind 配置被覆盖、布局类失效 */
  modules: ['@crm/ui-kit/nuxt', '@pinia/nuxt', '@nuxtjs/i18n', '@nuxt/ui'],

  ui: {
    primary: 'violet',
  },

  tailwindcss: {
    configPath: './tailwind.config.ts',
    viewer: false,
  },

  components: [
    { path: '~/components/common', pathPrefix: false },
    { path: '~/components/feature/admin', prefix: 'Admin', pathPrefix: false },
    { path: '~/components/feature/auth', prefix: 'Login', pathPrefix: false },
    { path: '~/components/feature/app', prefix: 'App', pathPrefix: false },
    { path: '~/components/feature/crm', prefix: 'Crm', pathPrefix: false },
    { path: '~/components/feature/leads', prefix: 'Leads', pathPrefix: false },
    { path: '~/components/feature/accounts', prefix: 'Accounts', pathPrefix: false },
    { path: '~/components/feature/contacts', prefix: 'Contacts', pathPrefix: false },
    { path: '~/components/feature/deals', prefix: 'Deals', pathPrefix: false },
    { path: '~/components/feature/dashboard', prefix: 'Dashboard', pathPrefix: false },
    { path: '~/components/feature/ai', prefix: 'Ai', pathPrefix: false },
    { path: '~/components/feature/design', pathPrefix: false },
  ],

  i18n: {
    locales: [
      { code: 'zh', iso: 'zh-CN', file: 'zh-CN.json' },
      { code: 'en', iso: 'en-US', file: 'en-US.json' },
    ],
    defaultLocale: 'zh',
    lazy: true,
    langDir: 'locales',
    restructureDir: false,
    bundle: {
      optimizeTranslationDirective: false,
    },
  },

  runtimeConfig: {
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:8080',
      /** 仅显式 true 时用 mock；默认走真实 /api/leads */
      useLeadsMock: process.env.NUXT_PUBLIC_USE_LEADS_MOCK === 'true',
      /** 仅显式 true 时用 mock；默认走真实 API */
      useDealsMock: process.env.NUXT_PUBLIC_USE_DEALS_MOCK === 'true',
      /** 仅显式 true 时用 mock；默认走真实 /api/dashboard/* */
      useDashboardMock: process.env.NUXT_PUBLIC_USE_DASHBOARD_MOCK === 'true',
    },
  },
})

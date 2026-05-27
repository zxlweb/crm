// https://nuxt.com/docs/api/configuration/nuxt-config
import { resolvePublicApiBase } from './utils/resolve-public-api-base'

export default defineNuxtConfig({
  devtools: { enabled: true },
  compatibilityDate: '2025-05-21',
  typescript: { strict: true },

  /** 局域网可通过本机 IP 访问（如 http://10.x.x.x:3000） */
  devServer: {
    host: '0.0.0.0',
    port: 3000,
  },

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
    { path: '~/components/feature/settings', pathPrefix: false },
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
      apiBase: resolvePublicApiBase(),
      /** 仅显式 true 时用 mock；默认走真实 /api/leads */
      useLeadsMock: process.env.NUXT_PUBLIC_USE_LEADS_MOCK === 'true',
      /** 仅显式 true 时用 mock；默认走真实 API */
      useDealsMock: process.env.NUXT_PUBLIC_USE_DEALS_MOCK === 'true',
      /** 仅显式 true 时用 mock；默认走真实 /api/dashboard/* */
      useDashboardMock: process.env.NUXT_PUBLIC_USE_DASHBOARD_MOCK === 'true',
      /** Phase 4：settings / audit / custom-fields；默认真实 API */
      useSettingsMock: process.env.NUXT_PUBLIC_USE_SETTINGS_MOCK === 'true',
      /** Phase 4：super-admin 健康度/套餐/TOP；默认真实 API */
      useAdminInsightsMock: process.env.NUXT_PUBLIC_USE_ADMIN_INSIGHTS_MOCK === 'true',
    },
  },
})

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  devtools: { enabled: true },
  compatibilityDate: '2025-05-21',
  typescript: { strict: true },

  css: ['~/assets/css/design-system.css'],

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

  modules: ['@nuxtjs/tailwindcss', '@pinia/nuxt', '@nuxtjs/i18n'],

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
    },
  },
})

#!/bin/bash

# ==================== CRM Frontend (Nuxt 3) 目录结构创建脚本 ====================
echo "🚀 开始创建 Nuxt 3 Frontend 项目目录结构..."

# 创建主目录
mkdir -p frontend
cd frontend

# 创建 Nuxt 3 标准 + 推荐的 CRM 目录结构
mkdir -p app/{assets,components,composables,features,middleware,pages,stores,utils}
mkdir -p app/components/{base,ui,feature,layout}
mkdir -p app/features/{auth,dashboard,leads,contacts,companies,deals,settings,rbac}
mkdir -p composables
mkdir -p types
mkdir -p i18n/locales
mkdir -p public
mkdir -p server/api

# 创建核心文件
touch nuxt.config.ts
touch tailwind.config.ts
touch tsconfig.json
touch .env
touch README.md
touch app/app.vue
touch app/pages/index.vue

# 创建重要 composables
touch composables/useApi.ts
touch composables/useAuth.ts
touch composables/useTenant.ts
touch composables/usePermission.ts
touch composables/useI18n.ts

# 创建 i18n 文件
touch i18n/locales/zh-CN.json
tt defineNuxtConfig({
  devtools: { enabled: true },
  typescript: { strict: true },
  
  modules: [
    '@nuxtjs/tailwindcss',
    '@pinia/nuxt',
    '@nuxtjs/i18n',
    '@nuxt/icon'
  ],

  i18n: {
    locales: [
      { code: 'zh', iso: 'zh-CN', file: 'zh-CN.json' },
      { code: 'en', iso: 'en-US', file: 'en-US.json' }
    ],
    defaultLocale: 'zh',
    lazy: true,
  },

  runtimeConfig: {
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:8080'
    }
  }
})
EOL

echo '
# CRM Frontend (Nuxt 3)

企业级多租户 CRM 系统前端

技术栈：
- Nuxt 3 + Vue 3
- TypeScript
- Tailwind CSS
- Pinia
- Zod
- Nuxt I18n
' > README.md

echo "✅ Nuxt 3 前端目录结构创建完成！"

# 显示目录树
echo "📁 生成的目录结构："
tree -L 3

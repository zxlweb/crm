# CRM Frontend (Nuxt 3)

企业级多租户 CRM 前端脚手架。

## 技术栈

- Nuxt 3 + Vue 3 + TypeScript (strict)
- Tailwind CSS
- Pinia
- Nuxt I18n（中/英）
- Zod（表单校验，Phase 1+ 使用）

## 快速开始

```bash
cd frontend
cp .env.example .env
npm install
npm run dev
```

## 核心 Composables

| 文件 | 说明 |
|------|------|
| `composables/use-api.ts` | API 请求，自动携带 JWT + `X-Tenant-ID` |
| `composables/use-auth.ts` | Token 状态管理 |
| `composables/use-tenant.ts` | 当前租户上下文 |
| `composables/use-permission.ts` | 权限检查 `can(resource, action)` |

## 权限组件

```vue
<PermissionGuard resource="leads" action="view">
  <LeadsList />
</PermissionGuard>
```

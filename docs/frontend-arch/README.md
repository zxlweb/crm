# 前端架构文档

Nuxt 3 前端实现约定，与 `frontend/` 代码目录对应。

## 当前状态

脚手架已就绪（Phase 0），详细规范随 Phase 1 补充。

## 技术栈

| 技术 | 用途 |
|------|------|
| Nuxt 3 | 框架、路由、SSR |
| TypeScript | 严格模式 |
| Tailwind CSS | 样式 |
| Pinia | 全局状态 |
| Nuxt I18n | 中/英 |
| Zod | 表单与 API 校验（Phase 1+） |

## 目录结构

```
frontend/
├── app.vue
├── pages/
├── components/          # base / ui / layout / feature
├── composables/
│   ├── use-api.ts       # API + JWT + X-Tenant-ID
│   ├── use-auth.ts
│   ├── use-tenant.ts
│   └── use-permission.ts
├── i18n/locales/        # zh-CN.json, en-US.json
└── nuxt.config.ts
```

## 核心约定（草案）

| 主题 | 约定 |
|------|------|
| 多租户 | `useTenant` 管理当前租户；`useApi` 自动带 `X-Tenant-ID` |
| 权限 | `usePermission().can(resource, action)` + `<PermissionGuard>` |
| i18n | key 使用英文；默认 `zh` |
| 命名 | 组件 PascalCase；composables `useXxx`；文件 kebab-case |

## 文档索引

| 文档 | 说明 |
|------|------|
| [ux-design-guidelines.md](./ux-design-guidelines.md) | **UX 设计交互规范**（原则、布局、组件、动效、无障碍） |
| [design-system.md](./design-system.md) | 设计系统 Token、双主题、图表组件 |
| [01-directory-structure.md](./01-directory-structure.md) | 目录结构 |

## 待补充文档

- [ ] `02-state-and-api.md`
- [ ] `03-rbac-frontend.md`

## 相关文档

- [API 设计](../api/00-api-design.md)
- [多租户架构](../architecture/01-multi-tenancy.md)
- [RBAC 设计](../architecture/02-rbac-design.md)

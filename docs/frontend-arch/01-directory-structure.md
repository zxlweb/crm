# 前端目录结构

**版本**：v1.0 · Monorepo

```
crm/
├── apps/web/               # @crm/web — Nuxt 应用
│   ├── pages/、layouts/、feature/、common/
│   ├── composables/        # useApi、useTheme（provide ui-kit）
│   └── nuxt.config.ts      # modules: @crm/ui-kit/nuxt
└── packages/ui-kit/        # @crm/ui-kit — 可发布组件库（见 02-ui-kit-modules.md）
```

## 应用内组件（`apps/web/components/`）

```
components/
├── common/permission-guard.vue
└── feature/
    ├── admin/
    └── auth/
```

Chart / UiThemeToggle 由 **`@crm/ui-kit`** 通过 Nuxt 模块注册。

## 约定

- 业务逻辑优先放 `composables` 或 `utils`（可测试）
- API 调用统一经 `useApi`
- 权限判断使用 `usePermission().can()` 或 `<PermissionGuard>`
- Monorepo 与组件模型：[02-ui-kit-modules.md](./02-ui-kit-modules.md)（双包、依赖、bridge、排障）
- 合并前：`pnpm --filter @crm/web lint:layers`；组件库：`pnpm --filter @crm/ui-kit test`


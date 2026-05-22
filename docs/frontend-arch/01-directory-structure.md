# 前端目录结构

**版本**：v0.1 · Phase 0

```
frontend/
├── app.vue                 # 根布局
├── pages/                  # 路由页面
├── components/             # 通用与业务组件
│   └── permission-guard.vue
├── composables/            # useApi / useAuth / useTenant / usePermission
├── utils/                  # 可单测纯函数（如 permissions.ts）
├── locales/                # i18n 文案
├── e2e/                    # （仓库根 e2e/）Playwright
├── nuxt.config.ts
└── vitest.config.ts
```

## 约定

- 业务逻辑优先放 `composables` 或 `utils`（可测试）
- API 调用统一经 `useApi`
- 权限判断使用 `usePermission().can()` 或 `<PermissionGuard>`

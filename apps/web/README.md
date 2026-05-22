# @crm/web

CRM Nuxt 3 应用（monorepo `apps/web`）。

## 依赖

- [`@crm/ui-kit`](../../packages/ui-kit) — 设计系统 CSS、Chart 组件、`UiThemeToggle`

## 开发

在**仓库根目录**：

```bash
pnpm install
pnpm dev
```

或：

```bash
pnpm --filter @crm/web dev
```

## 测试

```bash
pnpm --filter @crm/web test
pnpm --filter @crm/web lint:layers
```

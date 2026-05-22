# 前端架构文档

Nuxt 3 前端实现约定；应用在 `apps/web/`，组件库在 `packages/ui-kit/`。

## 命名规范

| 类型 | 格式 | 示例 |
|------|------|------|
| 索引 | `README.md`（无编号） | 本文件 |
| 正文 | `NN-{主题}.md` | `02-ui-kit-modules.md` |
| HTML | 与 md 同名（`make docs-html` 生成） | — |

`NN` 为两位序号，表示**推荐阅读顺序**。新增文档占用下一空闲序号并更新下表。

## 技术栈

| 技术 | 用途 |
|------|------|
| Nuxt 3 | 框架、路由、SSR |
| TypeScript | 严格模式 |
| Tailwind CSS | 样式 |
| Pinia | 全局状态 |
| Nuxt I18n | 中/英 |
| Zod | 表单与 API 校验（Phase 1+） |

## Monorepo 目录

```
apps/web/                 # @crm/web
packages/ui-kit/          # @crm/ui-kit
```

根目录：`pnpm install` → `pnpm dev`

## 文档索引（按阅读顺序）

| 序号 | 文档 | 说明 |
|------|------|------|
| 01 | [01-directory-structure.md](./01-directory-structure.md) | 应用目录速览 |
| 02 | [02-ui-kit-modules.md](./02-ui-kit-modules.md) | **ui-kit 模块、双包模型、依赖、主题 bridge、排障** |
| 03 | [03-design-system.md](./03-design-system.md) | Token、主题、Chart/Card API |
| 04 | [04-ux-design-guidelines.md](./04-ux-design-guidelines.md) | UX、动效、检查清单 |
| 05 | [05-component-scenarios.md](./05-component-scenarios.md) | Card/Chart 场景登记 |

## 待补充（预留序号）

| 序号 | 计划文档 | 说明 |
|------|----------|------|
| 06 | `06-state-and-api.md` | 前端状态与 API 封装 |
| 07 | `07-rbac-frontend.md` | 权限 UI 与路由守卫 |

## 相关文档

- [MVP 任务清单](../tasks/00-mvp-task-breakdown.md)
- [API 设计](../api/00-api-design.md)
- [多租户](../architecture/01-multi-tenancy.md) · [RBAC](../architecture/02-rbac-design.md)

# 架构决策记录（ADR）

ADR 记录**已接受的重要技术决策**及其背景，供后续开发与 AI 协作时快速理解「为什么」。

## 原则

- **不可改写历史**：已 Accepted 的 ADR 正文不修改，仅可标记 `Superseded`  
- **一事一议**：一个决策一篇 ADR，避免大而全  
- **及时记录**：代码落地后 1–2 天内补充，避免事后补回忆  

## 编号与命名

```
docs/architecture/adr/
├── README.md
├── 0000-template.md          # 复制此模板新建
├── 0001-shared-db-multi-tenancy.md
├── 0002-casbin-rbac.md
└── ...
```

- 编号：四位递增 `0001`, `0002`, …  
- 文件名：`NNNN-简短英文-kebab-case.md`  

## 状态定义

| 状态 | 含义 |
|------|------|
| Proposed | 提案中，未实施 |
| Accepted | 已采纳并实施 |
| Deprecated | 仍有效但计划替换 |
| Superseded | 已被新 ADR 取代（文内注明新编号） |

## 索引

| ADR | 标题 | 状态 | 日期 |
|-----|------|------|------|
| [0001](./0001-shared-db-multi-tenancy.md) | 共享数据库 + tenant_id 多租户 | Accepted | 2026-05-21 |
| [0002](./0002-casbin-rbac.md) | Casbin RBAC + resource:action | Accepted | 2026-05-21 |
| [0003](./0003-goose-migrations.md) | Goose 作为数据库迁移工具 | Accepted | 2026-05-21 |

## 何时写 ADR

- 数据库/中间件/框架选型  
- 多租户、权限、安全相关取舍  
- 前后端边界、API 风格变更  
- 明确拒绝过的方案（避免重复讨论）  

## 不必写 ADR

- 纯 UI 样式、文案调整  
- 无争议的 CRUD 字段增删  
- 临时调试手段  

新建 ADR：复制 [0000-template.md](./0000-template.md)。

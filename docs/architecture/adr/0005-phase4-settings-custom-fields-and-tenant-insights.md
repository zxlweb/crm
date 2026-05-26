# ADR-0005: Phase 4 Settings、自定义字段与租户健康度统一架构

**状态**：Accepted  
**日期**：2026-05-26  
**决策者**：系统架构师

## 背景（Context）

Phase 4 目标从“业务功能扩展”转为“平台可运营收尾”。现状问题：

1. 租户配置分散在 `tenants.config` 与前端常量，变更链路不可审计；
2. 自定义字段能力缺失，行业适配只能改代码；
3. 审计数据可查但难以按管理视角消费；
4. Super Admin 缺少跨租户健康度视图，不利于运营决策。

若继续沿用零散接口，会导致 FE/BE 反复改 schema、RBAC 映射不稳定、QA 用例难固化。

## 决策（Decision）

1. **Settings 统一入口**
   - 采用 `/api/settings/tenant` + `/api/settings/features`，集中管理租户可配置项。
   - 业务开关与 quota 等低频配置保留 JSONB 扩展能力。

2. **自定义字段采用“元数据驱动”**
   - 新增 `custom_fields` 元数据模型，MVP 支持 `text/select/date`。
   - 业务实体值存储沿用现有实体扩展机制（按服务层映射），避免 Phase 4 直接引入复杂 EAV。

3. **审计聚合以现有 `audit_logs` 为单一事实源**
   - 统计 API 不新增冗余明细表。
   - 导出能力先同步 CSV，后续再演进异步导出。

4. **Super Admin 运营统计走独立权限域**
   - 引入 `admin_tenant_insights:view`，路由集中在 `/api/super-admin/stats/*`。
   - 普通租户上下文与跨租户上下文严格隔离。

5. **前端保持双包模型**
   - 新增 `ChartRadar` 于 `@crm/ui-kit`，业务页面仅消费公开组件。
   - 业务逻辑留在 composables，页面负责编排。

## 备选方案（Options Considered）

| 方案 | 优点 | 缺点 |
|------|------|------|
| A. 所有配置继续写 `tenants.config`（无独立接口） | 开发快 | 边界不清、审计弱、前端耦合 |
| B. EAV 通用字段平台一次到位 | 扩展最强 | 超出 MVP 复杂度，测试成本高 |
| **C. Settings + 元数据字段 + 审计聚合（采纳）** | 成本可控、可交付、便于并行三轨 | 后续高级字段能力需再扩展 |

## 后果（Consequences）

### 正面

- Phase 4 BE/FE/QA 可按单一契约并行推进；
- 配置与字段变更均可审计，满足收尾合规要求；
- Super Admin 获得可解释的健康度与套餐分布视图。

### 负面 / 风险

- `custom_fields` 仅支持 3 种类型，复杂行业诉求需后置；
- 审计实时聚合在高数据量下可能出现查询压力，需监控并评估物化视图。

### 后续行动

- [ ] 迁移新增资源动作：`settings:*`、`custom_fields:*`、`audit:export`、`admin_tenant_insights:view`
- [ ] FE 补 `ChartRadar` Story 与 VTU
- [ ] QA 增加“设置变更 -> 审计可见 -> 越权阻断”E2E 主链路

## 关联文档

- PRD: [phase-4-system-settings-close-prd.md](../../prd/phase-4-system-settings-close-prd.md)
- API: [phase-4-system-settings-api.md](../../api/phase-4-system-settings-api.md)
- Architecture: [05-phase-4-settings-close-architecture.md](../05-phase-4-settings-close-architecture.md)
- Tasks: [00-mvp-task-breakdown.md](../../tasks/00-mvp-task-breakdown.md) Phase 4

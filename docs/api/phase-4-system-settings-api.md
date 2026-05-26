# Phase 4 API 契约：系统设置（Settings）与收尾

**版本**：v1.0  
**日期**：2026-05-26  
**状态**：Accepted（Implementation 开工前评审基线）  
**输入**：[phase-4-system-settings-close-prd.md](../prd/phase-4-system-settings-close-prd.md) · [00-mvp-task-breakdown.md](../tasks/00-mvp-task-breakdown.md) §Phase 4  
**架构**：[05-phase-4-settings-close-architecture.md](../architecture/05-phase-4-settings-close-architecture.md)  
**索引**：[00-api-design.md](./00-api-design.md)

> HTML：`phase-4-system-settings-api.html`（`make docs-html` 生成）

---

## 1. 范围与原则

| 类别 | Phase 4 交付 |
|------|--------------|
| **Settings (L1)** | 租户配置读取/更新、业务开关 |
| **Custom Fields (L1)** | Accounts/Contacts/Leads/Deals 自定义字段（文本/枚举/日期） |
| **Audit Stats (L1)** | 审计操作统计（Bar / Donut）+ 导出契约 |
| **Super Admin Insights (L1)** | 租户健康度 Radar、套餐分布 Donut、租户 TOP Bar |
| **不在 MVP** | 字段表达式规则、跨租户智能诊断、复杂审批流 |

**通用约定**（继承 [00-api-design.md](./00-api-design.md)）：

- Header：`Authorization`、`X-Tenant-ID`
- 响应：`{ code, message, data, pagination? }`
- 列表默认 `page=1`、`page_size=20`（最大 100）
- 写操作写 `audit_logs`
- **禁止** body 传 `tenant_id`

---

## 2. 枚举与共享类型

### 2.1 设置开关 `settings_key`

| key | 类型 | 默认 | 说明 |
|-----|------|------|------|
| `default_locale` | string | `zh-CN` | 默认语言，`zh-CN`/`en-US` |
| `timezone` | string | `Asia/Shanghai` | 租户时区 |
| `ai_preview_enabled` | boolean | `false` | Phase 2 Preview 能力总开关 |
| `lead_import_mode` | string | `manual_review` | `manual_review` / `auto_merge` |
| `sales_quota.amount` | number | `0` | Dashboard 配额 |
| `sales_quota.currency` | string | `CNY` | 配额币种 |

### 2.2 自定义字段 `field_type`

| code | 中文 | 校验 |
|------|------|------|
| `text` | 文本 | `max_length <= 255` |
| `select` | 枚举 | `options` 非空，value 唯一 |
| `date` | 日期 | `YYYY-MM-DD` |

### 2.3 自定义字段 `entity_type`

`account` | `contact` | `lead` | `deal`

### 2.4 健康度维度 `health_dimension`

`activity` | `config_completeness` | `audit_risk` | `data_freshness` | `feature_adoption`

---

## 3. Settings（租户配置）

### 3.1 资源字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `tenant_id` | uuid | 当前租户 |
| `tenant_name` | string | 租户名称 |
| `default_locale` | string | 默认语言 |
| `timezone` | string | 时区 |
| `business_switches` | object | 业务开关集合 |
| `sales_quota` | object | `{ amount, currency, period }` |
| `updated_at` | timestamptz | 最近更新时间 |
| `updated_by` | uuid | 最近更新人 |

### 3.2 接口

**权限**：`settings:view`、`settings:update`

| 方法 | 路径 | 权限 | 说明 |
|------|------|------|------|
| GET | `/api/settings/tenant` | view | 获取当前租户配置 |
| PATCH | `/api/settings/tenant` | update | 局部更新配置 |
| GET | `/api/settings/features` | view | 获取开关字典与当前值 |

**PATCH `/api/settings/tenant` 请求示例**：

```json
{
  "tenant_name": "Acme China",
  "default_locale": "en-US",
  "timezone": "Asia/Shanghai",
  "business_switches": {
    "ai_preview_enabled": true,
    "lead_import_mode": "manual_review"
  },
  "sales_quota": {
    "amount": 5000000,
    "currency": "CNY",
    "period": "2026-05"
  }
}
```

---

## 4. Custom Fields（自定义字段）

### 4.1 字段模型

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `id` | uuid | 出 | |
| `entity_type` | enum | 是 | `account/contact/lead/deal` |
| `field_key` | string | 是 | 租户内 + 实体内唯一，snake_case |
| `field_label` | object | 是 | `{ "zh-CN": "...", "en-US": "..." }` |
| `field_type` | enum | 是 | `text/select/date` |
| `required` | boolean | 否 | 默认 false |
| `options` | object[] | 否 | `select` 必填 |
| `default_value` | any | 否 | 按类型校验 |
| `display_order` | int | 否 | 默认 100 |
| `is_active` | boolean | 否 | 默认 true |
| `created_at/updated_at` | timestamptz | 出 | |

### 4.2 接口

**权限**：`custom_fields:view`、`custom_fields:update`

| 方法 | 路径 | 权限 | 说明 |
|------|------|------|------|
| GET | `/api/settings/custom-fields` | view | 列表（支持 `entity_type`） |
| POST | `/api/settings/custom-fields` | update | 创建 |
| PATCH | `/api/settings/custom-fields/:id` | update | 更新 |
| DELETE | `/api/settings/custom-fields/:id` | update | 逻辑删除（`is_active=false`） |

**POST 请求示例**：

```json
{
  "entity_type": "lead",
  "field_key": "industry_segment",
  "field_label": { "zh-CN": "行业子类", "en-US": "Industry Segment" },
  "field_type": "select",
  "required": false,
  "options": [
    { "value": "saas", "label": { "zh-CN": "SaaS", "en-US": "SaaS" } },
    { "value": "manufacturing", "label": { "zh-CN": "制造业", "en-US": "Manufacturing" } }
  ],
  "display_order": 30
}
```

---

## 5. Audit Stats（审计统计）

### 5.1 统计接口

**权限**：`audit:view`

| 方法 | 路径 | 图表 | Query |
|------|------|------|-------|
| GET | `/api/audit/stats/by-action` | `ChartDonut` | `from`, `to`, `module`, `actor_role` |
| GET | `/api/audit/stats/trend` | `ChartBar` | `from`, `to`, `granularity=day|week`, `action` |
| GET | `/api/audit/stats/top-actors` | `ChartBar` | `from`, `to`, `limit` |

**`by-action` 响应示例**：

```json
{
  "items": [
    { "action": "settings.update", "count": 18 },
    { "action": "custom_field.create", "count": 7 },
    { "action": "rbac.role.assign", "count": 5 }
  ],
  "total": 30
}
```

### 5.2 导出接口

**GET** `/api/audit/export`  
**权限**：`audit:export`  
**Query**：`from`, `to`, `module`, `actor_id`, `format=csv`

成功返回文件流；越权或跨租户返回 `403/404`。

---

## 6. Super Admin（跨租户运营）

> 此组接口仅需 `Authorization` + `is_super_admin=true`，不要求 `X-Tenant-ID`。

### 6.1 租户健康度

**GET** `/api/super-admin/stats/tenant-health`

**响应 data**：

```json
{
  "dimensions": ["activity", "config_completeness", "audit_risk", "data_freshness", "feature_adoption"],
  "items": [
    {
      "tenant_id": "uuid",
      "tenant_name": "Acme",
      "scores": {
        "activity": 78,
        "config_completeness": 85,
        "audit_risk": 42,
        "data_freshness": 80,
        "feature_adoption": 66
      },
      "overall_score": 70
    }
  ]
}
```

### 6.2 套餐分布 / 租户 TOP

| 方法 | 路径 | 图表 | Query |
|------|------|------|-------|
| GET | `/api/super-admin/stats/plan-distribution` | `ChartDonut` | `from`, `to` |
| GET | `/api/super-admin/stats/top-tenants` | `ChartBar` | `metric=activity|revenue|risk`, `limit=10` |

---

## 7. RBAC 动作与路由映射

| 资源 | action | 路由前缀 |
|------|--------|----------|
| `settings` | `view`, `update` | `/api/settings/tenant`、`/api/settings/features` |
| `custom_fields` | `view`, `update` | `/api/settings/custom-fields*` |
| `audit` | `view`, `export` | `/api/audit/stats/*`、`/api/audit/export` |
| `admin_tenant_insights` | `view` | `/api/super-admin/stats/*` |

---

## 8. 审计动作（新增）

| action | 触发 |
|--------|------|
| `settings.update` | PATCH `/api/settings/tenant` |
| `custom_field.create` | POST `/api/settings/custom-fields` |
| `custom_field.update` | PATCH `/api/settings/custom-fields/:id` |
| `custom_field.delete` | DELETE `/api/settings/custom-fields/:id` |
| `audit.export` | GET `/api/audit/export` |
| `admin.tenant_health.view` | GET `/api/super-admin/stats/tenant-health` |

---

## 9. 错误码补充

| message | HTTP | 场景 |
|---------|------|------|
| `setting_key_invalid` | 400 | 不支持的配置键 |
| `custom_field_key_conflict` | 409 | 字段 key 冲突 |
| `custom_field_type_invalid` | 400 | 类型与 default/options 不匹配 |
| `audit_export_rate_limited` | 429 | 导出限流 |
| `admin_scope_denied` | 403 | 非 super admin 访问跨租户统计 |

---

## 10. 实现顺序（供三轨领任务）

| 顺序 | 任务 ID | 交付 |
|------|---------|------|
| 1 | 4.1 | `settings/tenant` + `custom-fields` CRUD 契约 |
| 2 | 4.2 | `audit/stats/*` + `audit/export` |
| 3 | 4.3 | `super-admin/stats/tenant-health` + Radar 数据 |
| 4 | 4.4 | `plan-distribution` + `top-tenants` |
| 5 | 4.5 | i18n key 清单、Swagger 对齐、部署文档验收 |

---

## 11. 修订记录

| 日期 | 说明 |
|------|------|
| 2026-05-26 | v1.0：Phase 4 架构师 2a 首版，覆盖 settings/custom-fields/audit/admin-insights 契约 |

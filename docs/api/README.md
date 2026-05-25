# API 设计文档

接口契约与 [PRD](../prd/) 同步，编码前定稿。

## 文件索引

| 文件 | 范围 | 状态 |
|------|------|------|
| [00-api-design.md](./00-api-design.md) | MVP 全量 API 索引 | v0.1 |
| [phase-2-crm-ai.md](./phase-2-crm-ai.md) | Phase 2 客户关系 + 情绪旅程 + AI Preview | v1.0 Accepted |
| [phase-3-deals-dashboard-api.md](./phase-3-deals-dashboard-api.md) | Phase 3 商机 + 仪表盘 | v1.0 Accepted |

## 命名规范

- 总索引：`00-api-design.md`
- **Phase 迭代契约**（与 PRD 同 slug）：`phase-N-{slug}-api.md`（如 `phase-2-crm-ai.md`、`phase-3-deals-dashboard-api.md`）
- **单模块详约**（可选拆分）：`{模块}-api.md`（如 `auth-api.md`）

## 单接口描述模板

```markdown
### POST /api/auth/login

**描述**：用户登录  
**认证**：否  

**请求体**
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|

**响应 data**
| 字段 | 类型 | 说明 |
|------|------|------|

**错误码**：400 / 401 / 500
```

## 统一约定

- 响应格式：`{ code, message, data, pagination? }`（见 00-api-design）  
- 业务接口需 Header：`Authorization: Bearer <token>`、`X-Tenant-ID: <uuid>`  
- 权限标识：`resource:action`（与 Casbin 一致）

# 产品需求文档（PRD）

## 文件索引

| 文件 | 模块 | 状态 |
|------|------|------|
| [00-crm-overview.md](./00-crm-overview.md) | MVP 总览、角色、功能范围 | v0.2 |
| [phase-2-relationship-crm-prd.md](./phase-2-relationship-crm-prd.md) | Phase 2 关系经营 + AI Preview + **§11.5 销售工作台** | v0.4 |

## 命名规范

- 总览：`00-{产品名}-overview.md`
- 模块：`{模块}-prd.md`（如 `auth-prd.md`、`leads-prd.md`）

## 模块 PRD 必备章节

1. 背景与目标用户  
2. 用户故事 / 验收标准  
3. 功能列表（含 Out of Scope）  
4. 非功能需求（性能、安全、i18n）  
5. 依赖与里程碑（对应 `docs/tasks/`）

## 协作

- 新模块开发前：先 PRD → 再 [api/](../api/)  
- HTML 预览：运行 `make docs-html` 后打开同名 `.html`

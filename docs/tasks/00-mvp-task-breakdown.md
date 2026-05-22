### **任务拆分清单**
#### **MVP 总体开发阶段**

**Phase 0: 基础架构（已完成 · v2.1 收尾，无 CI）**
- [x] 完善数据库核心 Schema + 迁移脚本（00001–00003，含 dev seed）
- [x] Backend 项目初始化（Gin + GORM + Casbin）
- [x] Frontend 项目初始化（Nuxt 3 + TypeScript + Tailwind + Pinia + i18n）
- [x] 多租户中间件 + Tenant Context（含单元/集成测试）
- [x] RBAC 基础框架（Casbin 集成 + 路由权限映射测试）
- [x] 自动化测试：Backend `make test`、Frontend Vitest、E2E 冒烟（`e2e/`）
- [x] QA + Code Review 留档 → [qa/phase-0-qa.md](../qa/phase-0-qa.md) · [reviews/phase-0-review.md](../reviews/phase-0-review.md) · [meeting-notes/phase-0-notes.md](../meeting-notes/phase-0-notes.md)

**Phase 1: 认证与权限系统（建议优先完成）**
- [x] 用户登录 / Refresh Token（注册待后续）
- [ ] 用户注册
- [x] 多租户登录与切换
- [x] Super Admin 管理后台入口
- [x] RBAC 完整实现（角色、权限、分配）
- [x] 前端权限控制组件（usePermission、Permission 组件）
- [ ] 审计日志基础记录

**Phase 2: 客户与线索模块**
- [ ] 公司 (Accounts) CRUD + 搜索
- [ ] 联系人 (Contacts) CRUD + 关联公司
- [ ] 线索 (Leads) CRUD + 状态流转
- [ ] 线索导入（Excel）
- [ ] 线索分配功能
- [ ] 跟进记录（Activity）基础

**Phase 3: 商机与仪表盘**
- [ ] 商机 (Deals) CRUD + Pipeline 看板
- [ ] 销售阶段管理
- [ ] 仪表盘（关键指标 + 漏斗图）

**Phase 4: 系统设置与收尾**
- [ ] 租户配置管理
- [ ] 自定义字段支持
- [ ] i18n 完整接入（中英）
- [x] 基础错误处理与统一响应格式
- [ ] Swagger API 文档
- [x] 权限与多租户集成测试（Phase 0 `scope_test` + QA）
- [ ] 部署文档

#### **MVP 迭代开发阶段**

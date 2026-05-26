# Phase 4 PRD：系统设置与收尾

**产品名称**：EnterpriseFlow CRM  
**版本**：Phase 4 v0.3  
**日期**：2026-05-26  
**状态**：实现中（含成员/部门 v0.3 增补）  
**MVP 总览**：[00-crm-overview.md](./00-crm-overview.md)  
**前置依赖**：Phase 0–3（认证、RBAC、客户与线索、商机与仪表盘）  
**关联任务**：[00-mvp-task-breakdown.md](../tasks/00-mvp-task-breakdown.md) § Phase 4  
**架构产出**：[phase-4-system-settings-api.md](../api/phase-4-system-settings-api.md) · [05-phase-4-settings-close-architecture.md](../architecture/05-phase-4-settings-close-architecture.md) · [ADR-0005](../architecture/adr/0005-phase4-settings-custom-fields-and-tenant-insights.md)

---

## 0. 执行摘要（给决策者）

Phase 4 不再新增核心销售业务链路，目标是把 MVP 从“可演示”推进到“可交付、可运营、可审计”：  

| 收尾主题 | 业务价值 |
|----------|----------|
| 租户配置与自定义字段 | 让不同企业在不改代码的前提下完成本地化配置 |
| **角色与权限（系统设置内）** | 管理员可在 UI 为岗位/角色勾选功能权限，无需改 Casbin 策略文件 |
| 审计报表与健康度图表 | 让管理员/老板能看见“系统是否被正确使用” |
| i18n 全量落地 | 支撑双语团队统一使用体验 |
| Swagger 与部署文档 | 降低交接成本，支撑实施与后续迭代 |

**结论**：Phase 4 的验收标准是“平台运营能力闭环”与“文档闭环”，而不是新增销售功能。其中 **角色权限配置** 是租户管理员日常最高频的运营动作之一，必须在 `/settings` 内一站式完成。

---

## 1. 背景与业务目标

| 目标 ID | 背景问题 | Phase 4 目标 | 可度量结果 |
|---------|----------|--------------|------------|
| G1 | 目前租户配置能力分散，部分参数靠代码/seed | 建立统一“系统设置”入口（租户配置 + 业务开关） | 管理员可在 UI/API 完成关键配置，无需改代码 |
| G2 | 业务实体字段扩展能力不足 | 提供自定义字段机制（先覆盖 Accounts/Contacts/Leads/Deals） | 新增字段可在列表/详情可见并可检索（按约定范围） |
| G3 | 审计数据可查但不易读 | 提供审计统计图与租户健康度可视化 | 管理员 1 分钟内识别高风险租户或异常操作 |
| G4 | 文档链路缺口影响交接 | 完成 Swagger、部署文档、阶段收尾文档 | 新成员可按文档完成本地部署与接口调试 |
| G5 | Phase 1 已有 RBAC API，但缺少面向管理员的配置入口 | 在 **系统设置** 提供角色列表 + 权限树勾选 | 管理员 10 分钟内完成「新建角色 → 勾选权限 → 成员生效」 |
| G6 | 经理/管理员需按 **部门** 看团队，但无成员与部门对照 | **设置 → 成员** 列表 + 部门列；部门写入 `user_tenants.department` | 管理员可核对「谁属于哪个事业群」；与 Phase 3 `data_scope` 联动 |

---

## 2. 用户角色 (Persona) 与场景

| 角色 | 核心场景 | 关键诉求 |
|------|----------|----------|
| Tenant Admin（管理员） | 配置租户偏好、字段模板、**角色权限**、语言 | 不改代码即可完成配置，且有权限边界 |
| Sales Manager（销售经理） | 查看团队使用与审计趋势 | 快速识别异常使用与执行风险 |
| Super Admin（平台管理员） | 查看跨租户健康度、套餐与活跃情况 | 可做平台级运营判断 |
| Viewer（只读） | 查看设置与报表（受限） | 可读不可改，防误操作 |

| 场景 ID | 场景描述 | 验收结果 |
|---------|----------|----------|
| S1 | 管理员在 `settings` 配置租户信息、默认语言、业务开关 | 保存后即时生效（按配置类型定义是否需要刷新） |
| S2 | 管理员给 Leads 新增一个“行业子类”字段并用于录入 | 新字段在表单展示、详情可见、审计可追踪 |
| S3 | 经理在审计页查看最近 7 天关键操作分布 | 图表可筛选、可按角色/模块查看 |
| S4 | Super Admin 在后台查看租户健康度与套餐分布 | 可识别低活跃/高风险租户并导出复盘 |
| S5 | 管理员在 `settings` →「角色与权限」为「销售经理」勾选仪表盘与线索权限 | 保存后该角色用户刷新即生效；写审计 `rbac.role.permissions` |
| S6 | 管理员新建「客服专员」角色并仅授予联系人只读 | 无 `leads:update` 的用户无法编辑线索；越权 API 返回 403 |

### 2.1 参考交互（竞品对标，非 MVP 全量）

对标截图中的「组织架构 + 权限管理」布局，本 Phase **Must** 交付的是 **系统设置内的角色权限树**；**Won't** 在本 Phase 交付完整「部门/岗位组织树」（见 §7 MoSCoW）。

| 参考能力 | MVP（Phase 4） | 后续 |
|----------|----------------|------|
| 左侧选组织岗位 | 左侧 **角色列表**（租户内 `roles`） | 部门/岗位组织树（Phase 5+） |
| 右侧多级权限勾选 | **模块 → 资源 → action** 三级树 | 更细粒度「数据可见性」子项（如「可查看团队排名金额」） |
| 用户数展示 | 角色卡片展示 `user_count` | **设置 → 成员** 分配角色（`PUT /api/rbac/members/:id/roles`）已交付 |
| 部门列 | — | 成员表展示 `department`（一级部门名）；编辑部门 UI Phase 5+ |

---

## 3. 功能需求（User Story + Acceptance Criteria）

### 3.1 模块映射：任务清单 ↔ 本 PRD

| 任务 ID | PRD 章节 | 交付摘要 |
|---------|----------|----------|
| 4.1 | §3.2、§3.5、§4 | 租户配置 + 自定义字段 + **角色与权限** |
| 4.2 | §3.3、§4 | 审计日志统计图（Bar/Donut） |
| 4.3 | §3.4、§4 | Admin 租户健康度（Radar） |
| 4.4 | §3.4 | Admin 套餐分布（Donut） |
| 4.5 | §5、§6、§7 | i18n 完整接入 + Swagger + 部署文档 |

### 3.2 系统设置（租户配置 + 自定义字段）

| Story ID | 作为 | 我要 | 以便 | Acceptance Criteria |
|----------|------|------|------|---------------------|
| ST-01 | Tenant Admin | 配置租户基础信息（名称、时区、默认语言、业务开关） | 适配企业管理习惯 | 仅 `tenant_admin` 可改；保存写审计；跨租户不可见 |
| ST-02 | Tenant Admin | 管理模块级开关（如 AI Preview、导入策略） | 控制功能灰度与演示边界 | 开关变更有审计；前端按开关显隐 |
| ST-03 | Tenant Admin | 为 Accounts/Contacts/Leads/Deals 增加自定义字段 | 贴合不同行业字段需求 | 支持文本/枚举/日期（MVP）；字段名唯一；字段变更可追踪 |
| ST-04 | Sales | 在业务表单中填写自定义字段 | 避免在备注里写结构化数据 | 表单可见且可校验；详情与列表按配置展示 |

### 3.5 角色与权限（系统设置 · Must）

> **能力边界**：复用 Phase 1 已定义的 RBAC 契约（见 `docs/api/00-api-design.md` §RBAC），PM 不新增 OpenAPI 路径；Phase 4 交付 **设置页产品化** 与 **权限字典分组展示**。

#### 3.5.1 信息架构

| 元素 | 说明 |
|------|------|
| 入口 | `/settings`，Tab「角色与权限」（`?tab=roles`） |
| 布局 | 左：角色列表；右：当前角色的权限树 + 保存 |
| 可见性 | `rbac:view` 可进入只读；`rbac:manage` 可新建角色与保存勾选 |
| 权限树层级 | **模块**（如客户与商机、商业智能、系统设置）→ **资源**（如 `leads`、`dashboard`）→ **动作**（`view`/`create`/`update`/`delete`/`export`/`manage`） |

**模块分组（与实现对齐）**：

| 模块 key | 中文 | 包含资源（示例） |
|----------|------|------------------|
| `crm` | 客户与商机 | `leads`, `accounts`, `contacts`, `deals` |
| `insights` | 商业智能 | `dashboard` |
| `settings` | 系统设置 | `settings`, `custom_fields`, `audit` |
| `platform` | 平台运营 | `admin_tenant_insights`（仅 super_admin 角色可见项） |
| `system` | 系统管理 | `rbac` |

#### 3.5.2 User Story

| Story ID | 作为 | 我要 | 以便 | Acceptance Criteria |
|----------|------|------|------|---------------------|
| RB-01 | Tenant Admin | 查看本租户全部角色及成员数 | 了解权限分布 | `GET /api/rbac/roles`；列表含 `name`、`user_count`、`permission_ids` 数量 |
| RB-02 | Tenant Admin | 新建自定义角色（名称、描述） | 适配岗位差异 | `POST /api/rbac/roles`；`is_system=false`；重名返回业务错误 |
| RB-03 | Tenant Admin | 为选中角色勾选/取消权限 | 控制功能访问 | 模块/资源级 **全选/半选/取消** 联动子级；保存 `POST .../roles/:id/permissions` |
| RB-04 | Tenant Admin | 保存后成员权限即时生效 | 无需重新部署 | 被分配该角色的用户下次请求或刷新 `my-permissions` 后菜单/按钮与 API 一致 |
| RB-05 | Tenant Admin | 系统预置角色有保护提示 | 避免误删核心能力 | `is_system=true` 角色展示「系统」标记；**不可删除**（MVP 可不支持删角色，仅禁删系统角色） |
| RB-06 | Sales Manager | 只读查看某角色已有权限 | 协同制定权限策略 | `rbac:view` 时树只读，无保存按钮，展示 `roleReadonlyHint` |
| RB-07 | 系统 | 权限变更写审计 | 满足合规复盘 | 审计含 `role_id`、变更前后 `permission_ids` 摘要（或 diff 条数） |

#### 3.5.3 权限树交互细则

| 规则 | 要求 |
|------|------|
| 全选模块 | 勾选模块 checkbox → 该模块下全部 `permission_id` 加入草稿 |
| 半选态 | 模块或资源下仅部分 action 勾选时，父级 checkbox 为 **indeterminate** |
| 脏数据提示 | `draftPermIds` 与保存前 baseline 不一致时显示「保存权限」按钮 |
| 保存失败 | 展示错误文案，不静默丢失草稿 |
| 字典来源 | `GET /api/rbac/permission-items`；禁止前端硬编码权限 UUID |
| 标签 i18n | 模块/资源用 `permModule.*`、`permResource.*`；action 展示 `resource:action` + 可选 `description` |
| 生效范围 | 仅本租户；禁止勾选其他租户权限项 |

#### 3.5.4 与业务模块的映射（验收抽检）

| 勾选权限 | 预期 UI/API 行为 |
|----------|------------------|
| 取消 `leads:view` | 侧栏无 Leads；直接访问 `/leads` 或 API → 403 |
| 授予 `audit:view` 无 `export` | 可进 `/settings/audit` 看图，无导出按钮 |
| 授予 `settings:update` | 可改租户配置 Tab；Viewer 角色无此权限 |
| 授予 `rbac:manage` | 可保存权限树；仅 `rbac:view` 只读 |

#### 3.5.5 数据范围（与 Phase 3 联动）

| Story ID | 作为 | 我要 | AC | 状态 |
|----------|------|------|-----|------|
| RB-08 | Tenant Admin | 为角色配置 `data_scope`（`self` / `department` / `all`） | 保存后统计范围联动 | **Should（未做 UI）**；MVP 由 **角色名 + `rbac:manage`** 推导，见 Phase 3 PRD §5.3 |
| RB-09 | Tenant Admin | 在成员列表查看每人 **部门** 与 **角色** | `GET /api/rbac/members` 含 `department`；设置页「成员」Tab | **Must（已交付）** |
| RB-10 | Tenant Admin | 为成员分配/调整角色 | `PUT /api/rbac/members/:id/roles`；Modal 勾选角色 | **Must（已交付）** |

**MVP 数据范围判定（无需 RB-08 UI）**：

| 角色（示例） | `data_scope` | 说明 |
|--------------|--------------|------|
| 租户管理员（`rbac:manage`） | `all` | 全租户 CRM + 仪表盘 **部门排行** |
| 销售经理 | `department` | 同 `user_tenants.department` 成员的数据 + **成员排行** |
| 销售代表 / 只读 | `self` | 仅本人 |

部门字段由迁移/运营 seed（如小西演示 `00020_user_tenant_department.sql`）；**暂不支持**在设置页编辑部门（Phase 5+ 或钉钉同步）。

### 3.6 成员管理（系统设置 · Must）

| 元素 | 说明 |
|------|------|
| 入口 | `/settings`，Tab「成员」（`?tab=members`） |
| 可见性 | `rbac:view` 可查看；`rbac:manage` 可编辑成员角色 |
| 列 | 成员（姓名/邮箱）、**部门**、角色、加入时间 |
| API | `GET /api/rbac/members`、`PUT /api/rbac/members/:id/roles` |

| Story ID | 作为 | 我要 | AC |
|----------|------|------|-----|
| MB-01 | Tenant Admin | 查看本租户全部成员及角色 | 列表含 `roles[]`、`department` |
| MB-02 | Tenant Admin | 为成员勾选角色并保存 | 保存后该用户 `my-permissions` 与菜单一致；写审计 |
| MB-03 | Tenant Admin | 按姓名/邮箱搜索成员 | 客户端或 API 筛选；空态文案 |
| MB-04 | Sales Manager | （间接）仅见本部门 CRM 数据 | 依赖 MB-01 维护的 `department` + 角色「销售经理」；见 [phase-3-deals-dashboard-prd.md](./phase-3-deals-dashboard-prd.md) §4.4.4 |

### 3.3 审计与设置报表

| Story ID | 作为 | 我要 | Acceptance Criteria |
|----------|------|------|---------------------|
| AU-01 | Tenant Admin | 查看按操作类型聚合的审计图（Donut/Bar） | 支持按时间、模块、用户角色筛选 |
| AU-02 | Sales Manager | 识别异常操作波动（如删除、权限变更） | 图表可切时间范围；空态/错误态可识别 |
| AU-03 | Tenant Admin | 导出审计摘要用于复盘 | 至少支持 CSV 导出（MVP）或保留导出契约 |

### 3.4 Super Admin 运营视图

| Story ID | 作为 | 我要 | Acceptance Criteria |
|----------|------|------|---------------------|
| AD-01 | Super Admin | 查看租户健康度雷达图（活跃、配置完整度、审计风险、数据新鲜度） | `ChartRadar` 接统计 API，维度定义固定且可解释 |
| AD-02 | Super Admin | 查看套餐分布与头部租户贡献 | `ChartDonut` + `ChartBar` 组合展示，支持时间过滤 |
| AD-03 | Super Admin | 按租户 drill-down 查看详情 | 可跳转租户详情页或详情弹层（按现有 IA） |

---

## 4. 多租户与 RBAC 要求

| 维度 | 要求 |
|------|------|
| 多租户隔离 | 所有 settings/custom-fields/audit 查询强制使用 JWT 上下文租户；禁止请求体透传 tenant_id |
| 跨租户访问 | 跨租户资源访问返回 `404`（防枚举） |
| 权限动作建议 | `settings:view/update`、`custom_fields:view/update`、`audit:view/export`、`admin_tenant_insights:view`、**`rbac:view/manage`** |
| 角色权限页 | 进入 Tab 需 `rbac:view`；新建角色、保存勾选需 `rbac:manage` |
| 最小权限原则 | Viewer 默认只读；字段变更、开关变更、**角色权限变更**必须显式授权 |
| 系统角色 | `is_system=true` 的角色禁止删除；敏感权限变更建议二次确认（Could） |
| 审计要求 | 配置变更、字段变更、**角色权限变更**、导出动作均写审计，含 before/after 摘要 |
| 数据范围 | 经理仅看本租户；Super Admin 才可看跨租户汇总 |

---

## 5. i18n 与本地化需求

| 范围 | 要求 | 验收 |
|------|------|------|
| 页面文案 | 设置页、审计页、Admin 统计页、**角色与权限 Tab** 全量中英双语 key | 不允许硬编码中文/英文常量 |
| 权限树标签 | `permModule.*`、`permResource.*`、`roleManagement*` 等 | 切换语言后模块名与提示一致 |
| 自定义字段 | 字段显示名支持中英文标签（MVP 可先主标签 + 备用标签） | 切语言后标签同步变化 |
| 时间与数字 | 时间按租户时区展示；金额/数字按 locale 格式化 | `zh-CN` 与 `en-US` 样例截图可对照 |
| 错误与空态 | API 错误、权限提示、空态文案均国际化 | QA 冒烟覆盖关键错误态 |

---

## 6. 非功能需求（性能、安全、审计）

| 类别 | 要求 |
|------|------|
| 性能 | 设置页与审计页首屏 < 1.5s（单租户 200 用户规模）；统计 API P95 < 300ms |
| 安全 | 配置与字段变更必须鉴权；防 IDOR；导出接口限频 |
| 审计 | 关键动作 100% 留痕（配置、字段、权限、导出） |
| 可用性 | 图表 API 失败时单图降级，不阻塞整页；提供清晰错误态 |
| 测试 | BE：单测+集成测（租户/RBAC/审计）；FE：Vitest（`use-rbac`、`RolePermissionManager`）；QA：E2E（**改角色权限 → 越权/放行**） |
| 文档交付 | Swagger 与部署文档可独立指导联调和部署 |

---

## 7. 优先级 (MoSCoW)

| 级别 | Phase 4 范围 |
|------|--------------|
| Must | 租户配置管理（基础信息、开关） |
| Must | **系统设置 → 角色与权限**：角色列表 + 权限树勾选 + 保存（`rbac:view/manage`） |
| Must | **系统设置 → 成员**：成员列表 + 部门列 + 分配角色（`rbac:view/manage`） |
| Must | 自定义字段（文本/枚举/日期）基础能力 |
| Must | 审计统计图（`ChartDonut` / `ChartBar`） |
| Must | i18n 全量接入（设置、审计、Admin） |
| Must | Swagger API 文档 + 部署文档 |
| Should | Super Admin 健康度 `ChartRadar` |
| Should | 审计摘要导出能力 |
| Should | 角色级 `data_scope` 配置 UI（RB-08） |
| Should | 成员 **部门** 在线编辑 / 钉钉部门同步 |
| Could | 健康度阈值告警与推荐动作；权限变更二次确认 Modal |
| Won't（MVP） | **多级组织架构树**（部门/岗位层级 UI，对标截图左侧）；按品牌/代理商等细粒度数据权限 |
| **已交付（跨 Phase）** | 一级 `department` 字段 + Phase 3 数据范围/排行（见 Phase 3 PRD v0.2） |
| Won't（MVP） | 可视化字段设计器高级能力（联动规则、表达式）、跨租户智能诊断 |

---

## 8. 成功指标

| 指标 | 目标值 |
|------|--------|
| 配置可用性 | 管理员可在 5 分钟内完成“租户配置 + 开关调整 + 生效验证” |
| 权限配置效率 | 新建角色并勾选一组 CRM 只读权限 < 10 分钟；保存后抽检 1 个成员 API 行为正确 |
| 字段扩展效率 | 新增一个业务字段从创建到页面可用 < 10 分钟 |
| 审计可读性 | 管理员可在 1 分钟内定位近 7 天高频风险操作 |
| i18n 完整度 | Phase 4 范围页面双语覆盖率 100% |
| 文档可交接性 | 新成员按部署文档可在 30 分钟内完成本地启动与 Swagger 调试 |
| 质量门禁 | 与 Phase 4 相关自动化测试通过率 100% |

---

## 9. 风险与依赖

| 风险 | 影响 | 缓解策略 |
|------|------|----------|
| 自定义字段范围过大导致延期 | 影响 Phase 4 关闭节奏 | MVP 仅支持文本/枚举/日期，复杂类型后置 |
| 审计统计口径不统一 | 图表失真，管理判断错误 | 先冻结统计口径与时间窗，再开发图表 |
| i18n 收尾遗漏 | 上线后语言混杂 | 设立 i18n 清单与 CI 检查（禁硬编码） |
| 权限树与 Casbin 不同步 | 勾选保存后仍 403 | 保存走统一 `assignRolePermissions`；集成测覆盖 `rbac:manage` 与 `leads:view` 撤销 |
| 误改系统角色 | 租户不可用 | 系统角色禁删；关键权限变更审计 + 可选二次确认 |
| 部署文档与真实环境偏差 | 交接失败 | 用“从零环境复跑”验证文档可执行性 |

| 依赖项 | 责任方 | 说明 |
|--------|--------|------|
| Phase 4 API 契约与字段模型 | Architect | 落地到 `docs/api/` 与必要 ADR |
| RBAC 字典与路由 | BE（Phase 1 已有） | `permission-items` 含 Phase 4 新资源；`route.go` 与 PRD 资源表一致 |
| ChartRadar 组件与主题一致性 | FE / UI-kit | 遵循双包模型，先补 ui-kit 再嵌页面 |
| 审计与权限集成测试 | BE / QA | 覆盖多租户隔离、越权、导出场景 |
| 文档闭环（QA/Review/Notes） | PM / QA / Dev | 按 `development-workflow` 第 6 阶段执行 |

---

## 10. 前端切面（2b，Implementation 门禁）

### 10.1 路由与角色边界

| 路由 | 页面职责 | 角色边界 |
|------|----------|----------|
| `/settings` | 租户配置 + 自定义字段 + **角色与权限 Tab** | `settings/custom_fields` 编辑需对应 action；`rbac:view` 可见权限 Tab |
| `/settings?tab=roles` | 角色列表 + 权限树 | `rbac:manage` 可保存；`rbac:view` 只读 |
| `/settings?tab=members` | 成员列表 + 部门 + 分配角色 | `rbac:manage` 可保存角色；`rbac:view` 只读 |
| `/settings/audit` | 审计统计与导出 | `tenant_admin/manager` 查看，导出需 `audit:export` |
| `/admin` | 平台健康度、套餐分布、TOP 租户 | 仅 `super_admin` |

### 10.2 信息流（API -> composable -> feature -> ui-kit）

| 场景 | composable | feature 组件 | ui-kit 组件 |
|------|------------|--------------|-------------|
| 租户配置 | `use-settings` | `feature/settings/TenantSettingsForm*` | 表单组件 |
| 自定义字段 | `use-custom-fields` | `feature/settings/CustomField*` | 表格/弹窗组件 |
| 角色与权限 | `use-rbac` | `feature/settings/RolePermissionManager` | 原生 checkbox 树（业务层） |
| 成员管理 | `use-rbac` | `feature/settings/TenantMembersManager` | `UiModal`、表格 |
| 权限分组 | — | — | `utils/rbac-permission-groups`（`buildPermissionModules`） |
| 审计图 | `use-audit-stats` | `feature/settings/AuditStats*` | `ChartDonut`、`ChartBar` |
| 平台运营图 | `use-admin-tenant-insights` | `feature/admin/TenantInsights*` | `ChartRadar`、`ChartDonut`、`ChartBar` |

### 10.3 前端实现约束

- 页面层仅做编排，不承载业务校验；校验、权限预判断放在 composable。
- `Pinia` 仅用于跨页共享的租户配置快照，页内状态保持局部。
- 图表仅通过 `@crm/ui-kit` 暴露层引用，禁止 deep import。
- 所有新增页面必须提供稳定 `data-testid`，供 QA Phase 4 E2E 复用。

**角色与权限 `data-testid`（QA 必用）**：

| testid | 含义 |
|--------|------|
| `settings-tab-roles` | 角色 Tab |
| `role-permission-manager` | 管理器根节点 |
| `role-list` / `role-item-{id}` | 角色列表 |
| `role-create-btn` / `role-create-modal` | 新建角色 |
| `permission-tree` | 权限树容器 |
| `perm-{resource}-{action}` | 单个 action 勾选 |
| `role-perm-save` | 保存权限 |
| `settings-tab-members` | 成员 Tab |
| `tenant-members-manager` | 成员管理器根节点 |
| `member-row-{id}` | 成员行 |
| `members-roles-dialog` | 分配角色弹窗 |

### 10.4 QA 验收路径（角色权限）

1. `tenant-admin@demo.com` → `/settings?tab=roles` → 新建角色 → 仅勾选 `contacts:view` → 保存  
2. 将测试用户分配该角色（seed 或 `users/:id/roles`）→ 登录 → 可访问 Contacts、不可编辑 Leads  
3. `viewer@demo.com` 进入角色 Tab → 只读，无 `role-perm-save`  
4. 审计日志可见权限变更记录（action 以 BE 实现为准，如 `rbac.role.permissions`）

### 10.5 QA 验收路径（成员与部门）

1. `tenant-admin` → `/settings?tab=members` → 可见部门列（如「灵狐数据」）与各成员角色  
2. 为 `linghu@xiaoxi.com` 确认角色含「销售经理」、部门为「灵狐数据」  
3. 该用户登录 → Dashboard `data_scope=department`、团队排行为本部门成员（Phase 3 E2E）  
4. `moye@xiaoxi.com`（销售代表）→ 无团队排行、仅本人商机  

---

## 11. 修订记录

| 日期 | 说明 |
|------|------|
| 2026-05-26 | PM v0.1 初稿：系统设置、自定义字段、审计报表、i18n、Swagger、部署文档收尾范围定义 |
| 2026-05-26 | Architect 2a/2b 评审：API 契约、架构基线、ADR-0005 已落地，可进入三轨实现 |
| 2026-05-26 | PM v0.2：补充 §3.5 系统设置内「角色与权限」；对标组织架构截图界定 MVP/Won't；RBAC 验收与 E2E testid |
| 2026-05-26 | PM v0.3：§3.6 成员与部门（MB-01–04）；RB-08/09/10 与 Phase 3 数据范围对齐；MoSCoW 区分一级 department vs 组织树 |

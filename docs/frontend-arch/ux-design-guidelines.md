# UX 设计交互规范

EnterpriseFlow B2B CRM 前端用户体验与交互约定。视觉 Token 与主题实现见 [design-system.md](./design-system.md)。

| 文档 | 职责 |
|------|------|
| 本文档 | 体验原则、布局、交互、反馈、动效、无障碍 |
| [design-system.md](./design-system.md) | 颜色 Token、主题切换、图表组件 API |
| [01-directory-structure.md](./01-directory-structure.md) | 目录与组件分层 |

**在线预览**：`/design`（主题）、`/charts`（图表）

---

## 1. 设计目标

面向 **企业销售与运营** 用户，界面应：

- **专业可信**：B2B SaaS 气质，信息密度适中，避免娱乐化装饰
- **高效清晰**：关键指标一眼可读，列表与表单操作路径短
- **反馈及时**：加载、成功、失败、无权限均有明确状态
- **多租户安全**：租户边界在 UI 上可感知，避免误操作

---

## 2. 设计原则

### 2.1 一致性（Consistency）

- 同一类操作使用相同控件、文案与反馈方式
- 间距、圆角、字号遵循统一尺度（见 §4）
- 颜色仅通过设计 Token（`ds-*` / `--ds-*`），禁止页面内随意写死色值
- 图标风格统一：线性图标，`stroke-width` 1.75–2，与侧栏/按钮对齐

### 2.2 层次清晰（Hierarchy）

| 层级 | 用途 | 示例类名 |
|------|------|----------|
| L1 页面标题 | 当前模块主标题 | `text-ds-fg-heading` 大字号加粗 |
| L2 区块标题 | 卡片/表格标题 | `font-semibold text-ds-fg-heading` |
| L3 辅助说明 | 副标题、表头说明 | `text-sm text-ds-fg-muted` |
| L4 数据强调 | KPI、金额 | 大号数字 + 可选趋势色 |

### 2.3 反馈及时（Feedback）

所有用户操作须在 **300ms 内** 给出可见反馈（含禁用态、Loading、Toast/行内错误）。

| 场景 | 要求 |
|------|------|
| 按钮点击 | 立即 `disabled` + loading 文案（如「保存中…」） |
| 列表/详情加载 | 骨架屏或居中 Spinner，禁止空白无提示 |
| 表单提交失败 | 字段下或表单顶部错误区，`role="alert"` |
| 成功操作 | 行内状态更新或轻提示，避免仅控制台 log |
| 无权限 | 明确文案 + 引导（非空白页） |

### 2.4 容错与防误操作（Error Prevention）

- 破坏性操作（删除、停用租户）必须二次确认
- 主按钮与危险按钮视觉区分：主操作用 `ds-btn-primary`，危险用 `text-ds-danger` / 红色浅底
- 表单离开未保存：后续迭代建议拦截（Phase 2+）

---

## 3. 视觉与品牌

### 3.1 双主题

| 主题 | 场景 | 特征 |
|------|------|------|
| **V1**（默认） | 日常办公、明亮环境 | 浅紫 SaaS、白卡片、浅紫轨道 |
| **V2** | 大屏/夜间偏好 | 深灰 `#0A0A0A`、紫色图表强调（Optrixx 风） |

- 用户通过 `UiThemeToggle` 或 `/design` 切换，偏好写入 Cookie `crm-theme`
- 新页面必须 **双主题可读**，禁止只适配 V1

### 3.2 主色板（5 色）

EnterpriseFlow 以 **紫色为品牌主色**（非蓝色）。实现时优先用 CSS 变量 / Tailwind `ds-*`，下表给出 V1 参考色值（V2 同 Token、色值略亮，见 `design-system.css`）。

| 角色 | 参考色 (V1) | Token | 典型场景 |
|------|-------------|-------|----------|
| **Primary** | `#7c3aed` (Violet-600) | `--ds-brand`、`--ds-fg-brand`、`bg-ds-brand`、`text-ds-fg-brand` | 主按钮、文字链接、侧栏激活、图表主线、输入框聚焦环、品牌渐变区 |
| **Success** | `#059669` (Emerald-600) | `--ds-success`、`bg-ds-success-subtle`、`text-ds-fg-success` | 保存成功、启用/活跃状态、正向 KPI 趋势、完成类 Toast、表格「活跃」圆点 |
| **Warning** | `#ca8a04` (Amber-600) | *待接入* `--ds-warning`（规划） | 待审批、即将到期、配额将满、非阻断性提示条、需关注但不紧急的标签 |
| **Danger** | `#b91c1c` (Red-700) | `--ds-danger`、`bg-ds-danger-subtle`、`text-ds-fg-danger` | 删除/停用确认、表单校验错误、API 失败提示、破坏性主操作、失败状态 |
| **Neutral** | `#64748b` (Slate-500) | `--ds-fg-muted`、`--ds-fg`、`--ds-border`、`bg-ds-bg-muted` | 正文辅助说明、表头/轴标签、分割线、占位符、禁用态文案、次要按钮 |

#### 各角色使用说明

**Primary（品牌紫）**

- 每个页面仅 **1 个** 主 CTA（如「保存」「登录」），避免满屏紫按钮
- 链接：`text-ds-fg-brand` + hover 加深，不必下划线
- 禁止：大面积纯紫背景（除登录品牌区、洞察卡片等已定义模块）；**不以蓝色 `#2563eb` 作主色**

**Success（成功绿）**

- 语义为「已完成 / 正常 / 正向」，不用于 mere 信息提示
- 背景用浅底 `bg-ds-success-subtle`，文字/图标用 `text-ds-success`
- 图表正向趋势可与 Primary 并存，趋势数字优先 Success

**Warning（警告黄）**

- 语义为「需留意、可继续操作」，阻断性低于 Danger
- 用于：审批中、同步延迟、试用即将结束、字段未填但不影响提交草稿
- UI：浅黄底 + 深黄字/图标；**勿与 Danger 混用**于同一表单行
- 落地前在 `design-system.css` 增加 `--ds-warning` / `--ds-bg-warning-subtle` 后统一引用

**Danger（危险红）**

- 必须配文案说明后果；删除/停用需二次确认
- 错误信息靠近字段或表单顶部，`role="alert"`
- 列表「停用」等操作：可用紫系次要按钮，**仅不可逆删除**用 Danger 强调

**Neutral（中性灰）**

- 承载 80% 界面信息：说明、标签、边框、未选中导航
- 标题用 `--ds-fg-heading`，正文用 `--ds-fg`，弱化用 `--ds-fg-muted`
- 图表对比线、次要系列用 `--ds-chart-secondary`（中性灰紫）

#### 色板与主题关系

```
Primary   ████  #7c3aed  →  V2: #9333ea
Success   ████  #059669  →  V2: #34d399
Warning   ████  #ca8a04  →  两主题共用（待 Token 化）
Danger    ████  #b91c1c  →  V2: #f87171
Neutral   ████  #64748b  →  V2: #71717a（文案）；边框为半透明白
```

完整变量表见 [design-system.md](./design-system.md)。

### 3.3 字体

- 家族：**Plus Jakarta Sans**（见 `nuxt.config.ts`）
- 正文：`text-sm`（14px）为主
- 标题：`text-lg` ~ `text-3xl`，`font-semibold` / `font-bold`
- 数字/KPI：可加 `tracking-tight`、等宽不必强制

### 3.4 圆角与阴影

| 元素 | 圆角 | 阴影 |
|------|------|------|
| 卡片 | `rounded-2xl` | `shadow-sm` 或 `shadow-ds-brand`（主 CTA） |
| 输入框 / 按钮 | `rounded-xl` | 主按钮 `ds-brand-shadow` |
| 标签 / Badge | `rounded-full` | 无或极轻 |

---

## 4. 布局与间距

### 4.1 栅格

- 后台主布局：**侧栏 + 顶栏 + 内容区**（`layouts/admin.vue`）
- 内容区内边距：`p-6 lg:p-8`
- 卡片间距：垂直 `space-y-6`，网格 `gap-4` / `gap-6`
- 移动端：侧栏隐藏，顶栏保留主题切换与关键操作

### 4.2 间距尺度（Tailwind）

优先使用：`2 / 3 / 4 / 5 / 6 / 8`（即 8px 倍数）。表单项之间 `space-y-5`，表单项内 label 与 input `mb-1.5`。

### 4.3 内容宽度

- 表单单列：`max-w-sm` ~ `max-w-md`
- 仪表盘/列表：`max-w-6xl` 或全宽表格 `overflow-x-auto`
- 营销/选择页（`/design`）：`max-w-5xl` 居中

---

## 5. 组件交互规范

### 5.1 按钮

| 类型 | 样式 | 交互 |
|------|------|------|
| 主按钮 | `.ds-btn-primary` | `cursor-pointer`；`disabled:opacity-50`；hover 变色 |
| 次按钮 | `border border-ds-border` + 文字色 | 同左 |
| 文字按钮 | `text-ds-fg-muted hover:text-ds-fg-brand` | 用于次要导航 |
| 危险 | 红字/红底浅 | 必须带确认 |

- 所有可点击元素：**必须** `cursor-pointer`（项目 UX 硬性要求）
- 触摸目标最小高度约 **40px**（`py-2.5` ~ `py-3`）

### 5.2 表单

- 标签：`text-sm font-medium text-ds-fg`
- 输入：`.ds-input` 或等价 Token 类
- 占位符：`placeholder:text-ds-fg-muted` 风格
- 聚焦：边框 `--ds-border-focus` + 浅紫光晕（见 design-system.css）
- 校验：Zod + 行内错误，错误色 `text-ds-danger`，容器 `bg-ds-danger-subtle`

### 5.3 导航

- 侧栏项：默认 `text-ds-fg-nav`，激活 `.ds-nav-active` 或 `bg-ds-bg-muted` + `text-ds-fg-nav-active`
- 当前路由与 `NuxtLink` `active-class` 保持一致
- 面包屑/顶栏副标题：`text-xs text-ds-fg-brand`

### 5.4 表格

- 表头：`text-xs uppercase tracking-wide text-ds-fg-brand`，背景 `bg-ds-bg-muted`
- 行 hover：`hover:bg-ds-bg-muted`
- 状态列：圆点 + 文字（活跃/停用），成功/中性 Token
- 空状态：居中 `text-ds-fg-muted`，说明 + 可选主操作
- loading：骨架屏或居中 Spinner
- 排序：图标 + 文字（升序/降序）
- 筛选：下拉菜单 + 文字
- 导出：按钮 + 文字

### 5.5 卡片（KPI / 图表容器）

- 外壳：`.ds-card` + `rounded-2xl p-5`
- Hover：轻微抬升 + 阴影变化
- 指标卡：标题 muted + 数字 heading + 可选趋势（`text-ds-success`）
- 图表：使用 `ChartShell` 包裹，标题区与图表区分隔 `border-b border-ds-border-muted`

---

## 6. 数据可视化 UX

图表统一走 `components/chart/*` + `useChartTheme()`，**不要**在业务页手写 ECharts option。

### 6.1 图表选择

| 数据意图 | 组件 |
|----------|------|
| 趋势、对比（当前 vs 历史） | `ChartLine` |
| 品类占比、排名 | `ChartBar`（`horizontal`） |
| 时间序列对比 | `ChartBar`（纵向） |
| 转化漏斗 | `ChartFunnel` |

### 6.2 视觉约定

- **折线**：平滑曲线；主系列紫色 + 面积浅渐变；对比系列灰色虚线
- **柱状**：紫渐变柱体 + 浅紫轨道；忌过度动效（扫光已移除）
- **漏斗**：自上而下递减，标签置于块内，tooltip 显示名称与数值
- **坐标轴**：尽量无网格或极浅网格；轴标签 `text-ds-fg-muted` 量级
- **Tooltip**：圆角卡片，跟随 `--ds-chart-tooltip-*`

### 6.3 加载与空数据

- 图表区加载：`ClientOnly` fallback 文案「加载中…」
- 无数据：卡片内说明文案，不显示空白坐标系

---

## 7. 动效与过渡

### 7.1 时长与缓动

| 类型 | 时长 | 缓动 |
|------|------|------|
| 悬停、颜色 | 150–200ms | `ease-in-out` |
| 面板、模态 | 250–300ms | `ease-out` |
| 图表入场 | 600–900ms | `cubicOut`（组件内已配置） |

推荐类名：`transition-colors duration-200`、`transition-all duration-300 ease-in-out`。

### 7.2 必须有过渡的场景

- 页面区块出现（可选 `fade` + 轻微 `translate-y`）
- 模态框 / Drawer：淡入 + 缩放或滑入
- 列表项增删：高度 + 透明度（避免生硬跳变）
- 按钮：hover/active，避免瞬间跳色
- Toast 提示：从顶部滑入
- Skeleton Loading：脉冲动画

### 7.3 减少动效（无障碍）

```css
@media (prefers-reduced-motion: reduce) {
  /* 关闭非必要动画，保留必要状态色变化 */
}
```

图表组件在系统开启「减少动态效果」时，应关闭入场动画（`animated={false}`）。

### 7.4 禁止事项

- 全屏闪烁、高频循环扫光（影响性能与观感）
- 每帧刷新 ECharts 整表 option（会导致柱长异常）
- 超过 400ms 的阻塞式动画且无加载提示

---

## 8. 国际化（i18n）

- 文案 key：**英文**（如 `loginTitle`），展示由 `zh-CN` / `en-US` 翻译
- 禁止在组件内硬编码中文（除注释）
- 数字、货币、日期：使用 `Intl` 或统一 formatter，图表 `yFormatter` 与表格一致
- 邮箱等含 `@` 的文案在 JSON 中使用 `admin{'@'}demo.com` 避免 vue-i18n 解析错误

---

## 9. 权限与多租户体验

### 9.1 多租户

- 登录后展示可切换租户列表（`useTenant`）
- API 请求带 `X-Tenant-ID`；切换租户后刷新权限与列表
- UI 上避免跨租户数据混在同一视图（Super Admin 除外）
- 权限隐藏使用 v-permission + 平滑淡出
- 无权限时显示友好提示而非空白




### 9.2 RBAC

- 无权限：使用 `PermissionGuard` 或路由中间件拦截，文案 `noPermission`
- 菜单项：无权限不展示或禁用并 tooltip 说明
- 操作按钮：提交前可用 `can(resource, action)` 隐藏或禁用

### 9.3 Super Admin

- 与普通租户后台视觉一致，使用同一 `admin` 布局与 Token
- 危险操作（停用租户）突出确认与结果反馈

---

## 10. 无障碍（Accessibility）

- 图标按钮提供 `aria-hidden`（装饰）或 `aria-label`（仅图标可点击）
- 表单错误 `role="alert"`
- 对比度：正文与背景满足 WCAG AA（V1/V2 均需检查 muted 文字）
- 键盘：主流程可 Tab 聚焦，焦点环使用 `focus:ring-2 focus:ring-ds-brand/40`
- 图表：核心数据不仅依赖颜色，需标签/数值/tooltip 辅助

---

## 11. 文案语气（Tone of Voice）

- **简洁**：按钮 2–4 字为主（「保存」「停用」「导出」）
- **专业**：避免口语化；错误信息说明原因 + 可操作建议
- **一致**：同一动作全局同一动词（如统一用「停用」而非「禁用/关闭」混用）

---

## 12. 开发检查清单

新页面或组件合并前自检：

- [ ] 颜色使用 `ds-*` / `--ds-*`，无硬编码主题色
- [ ] V1、V2 下文字与边框可读
- [ ] 可点击元素有 `cursor-pointer` 与 hover/focus
- [ ] 异步操作有 loading / 错误 / 空状态
- [ ] 文案走 i18n key
- [ ] 图表使用 `ChartShell` + `ChartLine` / `ChartBar` / `ChartFunnel`
- [ ] 动效时长 ≤400ms，且考虑 `prefers-reduced-motion`
- [ ] 权限与租户边界符合 §9

---

## 13. 参考实现

| 页面/组件 | 说明 |
|-----------|------|
| `pages/login.vue` | 双栏登录、主题、表单反馈 |
| `pages/admin/index.vue` | KPI、图表、表格、租户操作 |
| `pages/charts/index.vue` | 图表组件陈列 |
| `pages/design.vue` | 主题选择与 Token 预览 |
| `components/ui/theme-toggle.vue` | 主题切换控件 |
| `components/chart/*` | 图表族 |
| `assets/css/design-system.css` | Token 源文件 |

---

## 14. 修订记录

| 日期 | 说明 |
|------|------|
| 2025-05 | 初版：对齐 V1/V2 设计系统、图表组件与前端架构师 UX 约定 |

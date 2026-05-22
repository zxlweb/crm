# 设计系统（Design System）

EnterpriseFlow 使用 **三层 Token 逻辑** + **CSS 变量** + **Tailwind 语义类**，由 `@crm/ui-kit` 统一提供，应用在 `apps/web` 消费。

> 体验原则：[04-ux-design-guidelines.md](./04-ux-design-guidelines.md) · Token 预览：`/design` · **场景手册**：[05-component-scenarios.md](./05-component-scenarios.md)（Chart + Card）

---

## 1. 三层逻辑（必读）

```
┌─────────────────────────────────────────────────────────┐
│  Component  组件层                                       │
│  .ds-card / .ds-input / .ds-btn-primary                 │
│  --ds-input-height / --ds-button-height-md              │
├─────────────────────────────────────────────────────────┤
│  Semantic  语义层（随 data-theme="v1"|"v2" 切换）          │
│  --ds-bg / --ds-fg / --ds-brand / --ds-chart-*          │
│  --ds-shadow-sm|md|lg（明暗主题不同阴影）                  │
├─────────────────────────────────────────────────────────┤
│  Foundation  基础层（:root，与主题无关）                   │
│  --ds-font-* / --ds-text-* / --ds-space-*               │
│  --ds-radius-* / --ds-duration-* / --ds-z-*             │
└─────────────────────────────────────────────────────────┘
```

| 层级 | 定义位置 | 是否随 V1/V2 变化 | 典型用途 |
|------|----------|-------------------|----------|
| **Foundation** | `:root` | 否 | 字号、间距、圆角、动效、z-index |
| **Semantic** | `:root` + `[data-theme]` | **是** | 背景、文字、边框、品牌、状态、图表色、阴影 |
| **Component** | CSS 工具类 + 尺寸 token | 部分（继承 semantic 色） | 按钮、输入、卡片、导航 |

**规则**：页面与 feature 组件只使用 `ds-*` / `var(--ds-*)`，禁止硬编码 `#7c3aed`、`rounded-xl`（应 `rounded-ds-xl`）等脱离 token 的值。

---

## 2. 主题（Semantic 色彩）

| ID | 名称 | 说明 |
|----|------|------|
| `v1` | 浅色紫 SaaS | **默认**，白底卡片、浅紫渐变 |
| `v2` | 深色 Optrixx | 深灰底、紫色图表强调 |

- Cookie：`crm-theme`（365 天）
- HTML：`document.documentElement[data-theme="v1"|"v2"]`
- 切换：`useTheme()` / `<UiThemeToggle />`；桥接与排障见 [02-ui-kit-modules.md](./02-ui-kit-modules.md#41-主题桥接必须)
- 预览：`/design`，或 `?theme=v2`

---

## 3. Foundation Token 一览

### 3.1 字体 Typography

| CSS 变量 | 说明 | Tailwind |
|----------|------|----------|
| `--ds-font-sans` | 主字体 Plus Jakarta Sans | `font-sans` |
| `--ds-font-mono` | 等宽 | `font-mono` |
| `--ds-text-xs` … `--ds-text-4xl` | 字号 | `text-ds-xs` … `text-ds-4xl` |
| `--ds-leading-*` | 行高（与字号配对） | 随 `text-ds-*` |
| `--ds-font-normal` … `--ds-font-bold` | 字重 400–700 | `font-ds-medium` 等 |

工具类：`.ds-text-body`、`.ds-text-caption`、`.ds-text-heading`

### 3.2 间距 Spacing（4px 基准）

| Token | 值 |
|-------|-----|
| `--ds-space-1` | 4px |
| `--ds-space-2` | 8px |
| `--ds-space-4` | 16px |
| `--ds-space-6` | 24px |
| `--ds-space-8` | 32px |

Tailwind：`p-ds-4`、`gap-ds-3`、`m-ds-6` 等。

### 3.3 圆角 Radius

| Token | 值 | Tailwind |
|-------|-----|----------|
| `--ds-radius-sm` | 6px | `rounded-ds-sm` |
| `--ds-radius-lg` | 12px | `rounded-ds-lg` |
| `--ds-radius-xl` | 16px | `rounded-ds-xl` |
| `--ds-radius-2xl` | 20px | `rounded-ds-2xl` |
| `--ds-radius-full` | 药丸形 | `rounded-ds-full` |

### 3.4 阴影 Shadow（随主题）

| Token | 用途 |
|-------|------|
| `--ds-shadow-sm` | 卡片默认 |
| `--ds-shadow-md` | 浮层、下拉 |
| `--ds-shadow-lg` | 模态 |
| `--ds-brand-shadow` | 主按钮 |
| `--ds-shadow-focus` | 聚焦环 |

Tailwind：`shadow-ds-sm`、`shadow-ds-brand`、`shadow-ds-focus`

### 3.5 动效 Motion

| Token | 值 | Tailwind |
|-------|-----|----------|
| `--ds-duration-fast` | 150ms | `duration-ds-fast` |
| `--ds-duration-normal` | 200ms | `duration-ds-normal` |
| `--ds-duration-slow` | 300ms | `duration-ds-slow` |
| `--ds-ease-default` | 标准曲线 | `ease-ds-default` |

工具类：`.ds-transition`（颜色/背景/边框/阴影过渡）

### 3.6 层级 Z-index

`--ds-z-dropdown` (1000) → `--ds-z-toast` (1080)  
Tailwind：`z-ds-modal`、`z-ds-dropdown` 等。

### 3.7 布局 Layout

| Token | 说明 |
|-------|------|
| `--ds-sidebar-width` | 侧栏宽度 260px |
| `--ds-topbar-height` | 顶栏 64px |
| `--ds-content-max` | 内容区最大宽 1280px |
| `--ds-page-px` / `--ds-page-py` | 页面内边距 |

---

## 4. Semantic Token 分组

| 分组 | 前缀示例 | Tailwind 示例 |
|------|----------|---------------|
| Surface | `--ds-bg`, `--ds-bg-elevated`, `--ds-bg-muted` | `bg-ds-bg`, `bg-ds-bg-elevated` |
| Foreground | `--ds-fg`, `--ds-fg-heading`, `--ds-fg-muted` | `text-ds-fg-heading` |
| Border | `--ds-border`, `--ds-border-focus` | `border-ds-border` |
| Brand | `--ds-brand`, `--ds-brand-hover` | `bg-ds-brand` |
| Status | `--ds-success`, `--ds-danger` | `text-ds-success` |
| Chart | `--ds-chart-line-end`, `--ds-chart-grid` | 图表内部读取 |

完整列表见 `packages/ui-kit/src/tokens/catalog.ts` 或 `/design` 预览。

---

## 5. Card 组件（与 Chart 并列）

场景登记见 [05-component-scenarios.md](./05-component-scenarios.md) Part A；**新增风格只加场景行 + `Card*` 组件**，不另开文档。

| 组件 | 场景 id | 用途 |
|------|---------|------|
| **`CardMetric`** | `dashboard` | 看板 KPI：左圆标 + 主数值 + 底栏同比 |
| **`CardShell`** | `content` | 标题 + 内容区（表/表单/详情） |

### 5.1 `CardMetric`（`dashboard` 场景）

| Prop | 类型 | 说明 |
|------|------|------|
| `label` | `string` | 指标名 |
| `value` | `string \| number` | 主数值（自动千分位） |
| `compareLabel` | `string?` | 底栏左 |
| `trend` | `string?` | 底栏右，如 `12%` |
| `trendDirection` | `'up' \| 'down' \| 'flat'` | 绿↑ / 红↓ |
| `iconTone` | `'info' \| 'calendar' \| 'accent' \| 'brand' \| 'neutral'` | 圆标底色 |

```vue
<CardMetric
  label="租户总数"
  :value="12"
  compare-label="累计"
  trend="+2"
  trend-direction="up"
  icon-tone="brand"
>
  <template #icon><!-- SVG --></template>
</CardMetric>
```

### 5.2 `CardShell`（`content` 场景）

| Prop | 说明 |
|------|------|
| `title` / `subtitle` | 可选标题区 |
| `noPadding` | 内容区是否去掉 `p-ds-5` |
| 默认 slot | 表格、表单等 |

### 5.3 类型

```ts
import { CARD_SCENARIOS, type CardScenarioId } from '@crm/ui-kit'
// CARD_SCENARIOS.dashboard.component === 'CardMetric'
```

---

## 6. Component Token 与工具类

| 尺寸 Token | 值 | 用途 |
|------------|-----|------|
| `--ds-input-height` | 40px | 输入框 `h-ds-input` |
| `--ds-button-height-md` | 40px | 按钮 `h-ds-btn-md` |
| `--ds-icon-md` | 20px | 图标 |

| 工具类 | 说明 |
|--------|------|
| `.ds-card` | 标准卡片（surface + 圆角 + 阴影） |
| `.ds-input` | 表单输入 |
| `.ds-btn-primary` | 主按钮 |
| `.ds-nav-active` | 侧栏激活 |
| `.ds-focus-ring` | 统一 focus-visible |
| `.ds-panel-brand` | 登录品牌区背景 |

---

## 7. 文件结构（Monorepo）

```
packages/ui-kit/
├── src/styles/design-system.css   # 全部 --ds-* 变量 + .ds-* 工具类
├── src/tokens/
│   ├── theme.ts                   # ThemeId、applyThemeToDocument
│   ├── foundation.ts              # 尺度元数据（文档/工具）
│   ├── catalog.ts                 # DS_TOKEN_GROUPS 目录
│   └── index.ts
├── tailwind.preset.ts             # Tailwind 扩展（apps/web 引用）

apps/web/
├── tailwind.config.ts             # presets: [ui-kit/tailwind.preset]
├── composables/use-theme.ts
├── plugins/00-ui-kit-theme.ts
└── pages/design.vue               # Token 全量预览
```

---

## 8. 使用方式

### 切换主题

```vue
<UiThemeToggle />
```

```ts
const { id, isDark, setTheme } = useTheme()
setTheme('v2')
```

### Tailwind 示例（推荐）

```vue
<article class="rounded-ds-xl border border-ds-border bg-ds-bg-elevated p-ds-6 shadow-ds-sm">
  <h2 class="text-ds-2xl font-ds-semibold text-ds-fg-heading">标题</h2>
  <p class="mt-ds-2 text-ds-sm text-ds-fg-muted">说明文字</p>
  <button class="ds-btn-primary ds-transition mt-ds-4 h-ds-btn-md rounded-ds-lg px-ds-4 text-ds-sm">
    提交
  </button>
</article>
```

### 读取 Token 目录（工具/文档）

```ts
import { DS_TOKEN_GROUPS, DS_TOKEN_LAYERS } from '@crm/ui-kit/tokens'
```

### 新增页面约定

1. 色彩、阴影、圆角、间距、字号均使用 `ds-*` token
2. 图表颜色读 `--ds-chart-*`（`useChartTheme`）
3. 条件插画可用 `isDark`：`v-if="!isDark"` / `v-else`
4. 动效使用 `ds-transition` + `duration-ds-*`，避免随意 `duration-300`

---

## 9. 图表组件（ECharts）

基于 `echarts` + `vue-echarts`，颜色从 `--ds-chart-*` 读取，随 V1/V2 自动切换。

| 组件 | 用途 |
|------|------|
| `ChartShell` | 卡片外壳：标题、主指标、图例 |
| `ChartLine` | 折线图 |
| `ChartBar` | 柱状图（横/纵） |
| `ChartFunnel` | 漏斗图 |

展示页：`/charts`  
场景手册：[05-component-scenarios.md](./05-component-scenarios.md)；Phase 交付勾选：[../tasks/00-mvp-task-breakdown.md](../tasks/00-mvp-task-breakdown.md)。

---

## 10. 路由

| 路径 | 说明 |
|------|------|
| `/login` | 统一登录 |
| `/admin` | Super Admin |
| `/design` | **全量 Token 预览**（三层目录） |
| `/charts` | Chart 组件展示 |
| `/cards` | Card 场景展示（`dashboard` / `content`） |

---

## 11. 修订记录

| 日期 | 说明 |
|------|------|
| 2025-05 | 初版：双主题色彩 |
| 2026-05-22 | 补齐 Foundation/Semantic/Component 三层；typography/spacing/radius/shadow/motion/z-index；tailwind.preset；/design 目录 |
| 2026-05-22 | `CardMetric` / `CardShell`；场景表 `CARD_SCENARIOS` |

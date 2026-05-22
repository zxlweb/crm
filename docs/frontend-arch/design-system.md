# 设计系统（Design System）

EnterpriseFlow 前端使用 **CSS 变量 + Tailwind 语义类** 实现双主题，V1 为默认，V2 为可切换深色主题。

> 体验与交互原则见 [ux-design-guidelines.md](./ux-design-guidelines.md)。

## 主题

| ID | 名称 | 说明 |
|----|------|------|
| `v1` | 浅色紫 SaaS | **默认**，白底卡片、浅紫渐变 |
| `v2` | 深色 Optrixx | 深灰底 `#0A0A0A`、紫色图表强调 |

- Cookie：`crm-theme`（365 天）
- HTML 属性：`document.documentElement[data-theme="v1"|"v2"]`
- 预览：`/design`，或 URL `?theme=v2`

## 文件结构

```
frontend/
├── design-system/tokens.ts      # 类型与元数据
├── assets/css/design-system.css # CSS 变量定义
├── composables/use-theme.ts     # setTheme / toggleTheme
├── plugins/theme.client.ts      # 客户端同步 data-theme
├── components/ui/theme-toggle.vue
```

## 使用方式

### 切换主题

```vue
<UiThemeToggle />
```

```ts
const { id, isDark, setTheme, toggleTheme } = useTheme()
setTheme('v2')
```

### Tailwind 语义色

| 类名 | 用途 |
|------|------|
| `bg-ds-bg` | 页面背景 |
| `bg-ds-bg-elevated` | 卡片、表单区 |
| `bg-ds-bg-muted` | 次级背景、表头 |
| `text-ds-fg-heading` | 标题 |
| `text-ds-fg-muted` | 辅助文字 |
| `bg-ds-brand` / `text-ds-on-brand` | 主按钮 |
| `border-ds-border` | 边框 |

### 工具类（CSS）

| 类名 | 用途 |
|------|------|
| `.ds-card` | 标准卡片 |
| `.ds-input` | 输入框 |
| `.ds-btn-primary` | 主按钮 |
| `.ds-panel-brand` | 登录左侧品牌区 |
| `.ds-nav-active` | 侧栏激活项 |

### 新增页面约定

1. 禁止硬编码 `violet-600`、`#0D0D0D` 等主题色，使用 `ds-*` token
2. 图表 SVG 读取 `--ds-chart-*` 变量（参考 `AdminOverviewChart`）
3. 条件插画可用 `isDark`：`v-if="!isDark"` / `v-else`

## 图表组件（ECharts）

基于 `echarts` + `vue-echarts`，颜色从 `--ds-chart-*` 读取，随 V1/V2 自动切换。

| 组件 | 用途 |
|------|------|
| `ChartShell` | 卡片外壳：标题、主指标、图例 |
| `ChartLine` | 折线图（平滑曲线、面积、对比虚线、发光） |
| `ChartBar` | 柱状图（横向/纵向；紫色渐变、轨道背景、错峰入场） |
| `ChartFunnel` | 漏斗图（销售管道） |

展示页：`/charts`

```vue
<ChartShell title="Sales Overview" metric="$88,692" metric-label="average sales">
  <ChartLine
    :categories="days"
    :series="[
      { name: 'Current', data: current, primary: true },
      { name: 'Before', data: before, compare: true },
    ]"
    :y-formatter="(v) => `$${v}k`"
  />
</ChartShell>
```

## 路由

| 路径 | 说明 |
|------|------|
| `/login` | 统一登录（随主题变化） |
| `/admin` | 统一 Super Admin |
| `/design` | 主题选择与 token 预览 |
| `/charts` | 图表组件展示（折线/柱/漏斗） |
| `/v2/login` | 重定向至 `/login` 并设 V2 |
| `/v2/admin` | 重定向至 `/admin` 并设 V2 |

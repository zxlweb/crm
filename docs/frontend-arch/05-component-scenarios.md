# CRM 组件场景手册（Chart + Card）

**版本**：v1.1  
**原则**：新增固定设计风格 = 新增 **场景 id** + **`Card*` / `Chart*` 组件** + 在本表补一行；**不单独开文档**。

**关联**：[03-design-system.md](./03-design-system.md) · [../tasks/00-mvp-task-breakdown.md](../tasks/00-mvp-task-breakdown.md) · [prd/00-crm-overview.md](../prd/00-crm-overview.md)

**展示页**：`/charts`（图表）· `/cards`（卡片场景）· `/design`（主题 / `Ui*` 原子）

---

## 0. 如何扩展（Chart / Card 通用）

| 步骤 | Card 示例 | Chart 示例 |
|------|-----------|------------|
| 1. 定场景 id | `dashboard` | `pipeline`（已有 funnel） |
| 2. 实现组件 | `packages/ui-kit/.../card/card-metric.vue` | `chart-line.vue` |
| 3. 注册 Nuxt | `module.ts` → `prefix: 'Card'` | `prefix: 'Chart'`；Chart 还需 `runtime/echarts.client.ts` |
| 4. **更新案例页** | `apps/web/pages/cards/index.vue` 新增 section | `apps/web/pages/charts/index.vue` 新增 section |
| 5. **i18n** | `chartsSection*` / `cardsSection*` 等 key（zh + en） | 同上 |
| 6. 登记场景 | 下表 Part A，状态 ✅ | 下表 Part B，状态 ✅ |
| 7. 业务接入 | Admin/Dashboard 顶栏 | 报表区 |

### 0.1 案例页工作流（每次加组件必做）

新增或**首次对外使用**的 `@crm/ui-kit` 组件（`Chart*` / `Card*` / `Ui*`）时，**同一 PR / 同一 FE 子任务**内必须完成案例页，禁止「只合组件、案例页后补」。

**原则**：组件与案例页 **同 PR、同子任务** 交付；案例页使用 fixture，不绑真实 API。

| 组件类型 | 路由 | 案例页文件 |
|----------|------|------------|
| `Chart*` | `/charts` | `apps/web/pages/charts/index.vue` |
| `Card*` | `/cards` | `apps/web/pages/cards/index.vue` |
| `Ui*` | `/design` | `apps/web/pages/design.vue` |

`apps/web/components/feature/*` 的**可复用**展示：优先在对应 Chart/Card 案例页加「业务用法」小节；仅该域使用的组件可在业务页内嵌 demo，但仍须在本文件场景表登记。

**标准步骤（与 §0 表步骤 4–7 对应）**

1. **实现组件** — `packages/ui-kit/src/components/ui/...`
2. **注册** — `src/index.ts` 导出；Chart 类还需 `runtime/echarts.client.ts` 注册 ECharts 类型
3. **更新案例页** — 上表对应 `index.vue` 增加独立 `<section>`（`ChartShell` / 栅格与现有风格一致），**fixture 数据**，不依赖真实 API
4. **i18n** — `locales/zh-CN.json` 与 `en-US.json` 同步 section 标题/说明（英文 key）
5. **场景登记** — 下表 Part A/B 状态 ✅，执行 `make docs-html`
6. **自测** — 本地打开案例路由，切换 V1/V2 与明暗主题；有 E2E 时跑 `e2e/tests/charts-theme.spec.ts`

**案例页 section 模板（Chart）**

```vue
<section>
  <h2 class="mb-4 text-sm font-semibold uppercase tracking-wider text-ds-fg-muted">
    {{ $t('chartsSectionDonut') }}
  </h2>
  <ChartShell :title="$t('chartsDonutTitle')" :height="320">
    <ChartDonut :items="donutItems" :height="300" />
  </ChartShell>
</section>
```

**完成定义（DoD）**

- [ ] 案例页能单独展示该组件的典型 props（含空数据 / 单条数据若常见）
- [ ] 中英文文案齐全
- [ ] 本文件场景表已更新
- [ ] `pnpm --filter @crm/ui-kit test`（若改 ui-kit）

类型注册：`packages/ui-kit/src/card/types.ts` → `CARD_SCENARIOS`

---

## Part A · Card 场景

与 **Chart** 同级，组件前缀 **`Card*`**。

### A.1 场景注册表

| 场景 id | 组件 | 用途 | 状态 |
|---------|------|------|------|
| **`dashboard`** | **`CardMetric`** | 看板 / Admin / Dashboard **KPI 行**（左圆标 + 底栏同比） | ✅ |
| **`content`** | **`CardShell`** | 表格、表单、详情**内容容器**（标题 + slot） | ✅ |

> 新风格示例：`CardCompare` + 场景 `compare` — 只加上表一行 + 一个 vue 文件，无需新 md。

### A.2 `dashboard` 场景 · `CardMetric`（设计偏好）

布局（**禁止**改成右上大图标 + 趋势贴标题的旧样式）：

```
┌────────────────────────────────────┐
│ (○)  指标名 · 小字灰                │
│      1,235 · text-ds-3xl           │
│  同比增长              12% ↓       │
└────────────────────────────────────┘
```

| 规则 | 说明 |
|------|------|
| 栅格 | `grid gap-ds-4 sm:grid-cols-2 xl:grid-cols-4` |
| 趋势 | `up` → 绿 ↑；`down` → 红 ↓ |
| `iconTone` | 同行相邻卡用不同 tone（info / calendar / accent / brand / neutral） |

```vue
<CardMetric
  label="昨日产量(万吨)"
  :value="22"
  compare-label="同比增长"
  trend="12%"
  trend-direction="down"
  icon-tone="info"
>
  <template #icon><!-- SVG h-5 w-5 --></template>
</CardMetric>
```

### A.3 `content` 场景 · `CardShell`

```vue
<CardShell title="租户列表" subtitle="最近更新">
  <table>...</table>
</CardShell>
```

### A.4 模块 × Card

| 模块 | 场景 | 组件 |
|------|------|------|
| Admin / Dashboard 顶栏 KPI | `dashboard` | `CardMetric` |
| 列表/表格外框 | `content` | `CardShell` |
| 图表外框 | — | `ChartShell`（与 Card 分工） |

---

## Part B · Chart 场景

组件前缀 **`Chart*`**。下列章节沿用原图表选型内容。

### B.1 选型原则

| 原则 | 说明 |
|------|------|
| **回答一个问题** | 每张图只服务一个业务问题 |
| **匹配数据形态** | 时间序列 → 折线；转化 → 漏斗；占比 → 环/饼 |
| **设计系统一致** | `--ds-chart-*`，`useChartTheme` |

### B.2 总览：模块 × 图表

| CRM 模块 | 典型场景 | 推荐组件 | ui-kit 状态 |
|----------|----------|----------|-------------|
| **仪表盘** | 营收趋势 | `ChartLine` | ✅ |
| **仪表盘** | KPI 数字 | `CardMetric`（**Card**） | ✅ |
| **仪表盘** | 管道 | `ChartFunnel` | ✅ |
| **线索** | 来源占比 | `ChartDonut` | ✅ |
| **跟进 Activity** | 类型分布 TOP | `ChartBar` | ✅ |
| **商机** | 销售管道 | `ChartFunnel` | ✅ |
| **仪表盘** | KPI 迷你趋势 | `ChartSparkline` | ✅ |
| **仪表盘** | 配额完成率 | `ChartGauge` | ✅ |
| **团队** | 业绩排行 | `ChartBar` | ✅ |

（完整矩阵见历史版本各模块行，与 Part B 图表专项一致。）

### B.3 图表类型速查

折线、柱状、漏斗、环形、雷达、热力、桑基等 — 见 [03-design-system.md](./03-design-system.md) §9。

### B.4 与 MVP Phase 同步（Chart）

| MVP Phase | 新增 Chart | Card 并行 |
|-----------|------------|-----------|
| Phase 1 | `ChartLine` on `/admin` | `CardMetric` on `/admin` ✅ |
| Phase 2 | `ChartDonut` | `CardMetric` on Leads 报表 |
| Phase 3 | `ChartSparkline`、`ChartGauge` | Dashboard 顶栏 `CardMetric` 行 |

Phase 排期与验收勾选见 [../tasks/00-mvp-task-breakdown.md](../tasks/00-mvp-task-breakdown.md)。

---

## 修订记录

| 日期 | 说明 |
|------|------|
| 2026-05-22 | 初版图表场景（`05-crm-chart-scenarios`） |
| 2026-05-22 | 合并为 Chart + Card；`CardMetric` / `dashboard` 场景；废止独立 07 文档 |
| 2026-05-22 | 统一编号为 `05-component-scenarios.md` |
| 2026-05-24 | §0.1 案例页工作流；`ChartDonut` 案例页与 Leads 报表接入 |
| 2026-05-24 | Activity 摘要 `ChartBar`：线索详情 + `/charts` 案例区块 |

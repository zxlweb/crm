<template>
  <ClientOnly>
    <VChart
      class="w-full"
      :style="{ height: `${height}px` }"
      :option="option"
      autoresize
    />
    <template #fallback>
      <div
        class="flex items-center justify-center text-sm text-ds-fg-muted"
        :style="{ height: `${height}px` }"
      >
        {{ loadingText }}
      </div>
    </template>
  </ClientOnly>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { HeatmapSeriesOption } from 'echarts/charts'
import type { ChartHeatmapPoint } from '../../../chart/types'
import { useChartTheme } from '../../../chart/use-chart-theme'

const props = withDefaults(
  defineProps<{
    points: ChartHeatmapPoint[]
    rows?: string[]
    columns?: string[]
    /** 按列语义着色（如情绪维）；设置后隐藏连续 visualMap */
    columnColors?: string[]
    height?: number
    loadingText?: string
    valueLabel?: string
  }>(),
  { height: 260, loadingText: 'Loading…', valueLabel: 'Count' },
)

const { colors, baseGrid } = useChartTheme()

const rowLabels = computed(() => {
  if (props.rows?.length) return props.rows
  return [...new Set(props.points.map((p) => p.row))]
})

const columnLabels = computed(() => {
  if (props.columns?.length) return props.columns
  return [...new Set(props.points.map((p) => p.column))]
})

const maxValue = computed(() =>
  Math.max(...props.points.map((p) => p.value), 1),
)

const pointMeta = computed(() => {
  const map = new Map<string, string>()
  for (const p of props.points) {
    map.set(`${p.column}\0${p.row}`, p.meta ?? '')
  }
  return map
})

const heatmapData = computed(() =>
  props.points.map((p) => {
    const x = columnLabels.value.indexOf(p.column)
    const y = rowLabels.value.indexOf(p.row)
    return [x, y, p.value] as [number, number, number]
  }),
)

function hexToRgb(hex: string) {
  const h = hex.replace('#', '')
  const n = h.length === 3
    ? h.split('').map((c) => parseInt(c + c, 16))
    : [h.slice(0, 2), h.slice(2, 4), h.slice(4, 6)].map((c) => parseInt(c, 16))
  return { r: n[0], g: n[1], b: n[2] }
}

function cellColor(colIndex: number, value: number) {
  const base = props.columnColors?.[colIndex]
  if (!base) return undefined
  const { r, g, b } = hexToRgb(base)
  const t = 0.22 + 0.78 * (value / maxValue.value)
  return `rgba(${r}, ${g}, ${b}, ${t.toFixed(3)})`
}

const option = computed(() => {
  const c = colors.value
  const rows = rowLabels.value
  const cols = columnLabels.value

  const series: HeatmapSeriesOption = {
    type: 'heatmap',
    data: heatmapData.value,
    label: {
      show: true,
      color: c.axis,
      fontSize: 11,
      fontWeight: 600,
      formatter: (params) => {
        const val = Array.isArray(params.value) ? params.value[2] : params.value
        return val === 0 ? '' : String(val)
      },
    },
    itemStyle: {
      borderColor: 'rgba(255,255,255,0.85)',
      borderWidth: 2,
      borderRadius: 6,
      ...(props.columnColors
        ? {
            color: (params: { value?: [number, number, number]; data?: [number, number, number] }) => {
              const tuple = params.value ?? params.data
              if (!tuple) return colors.value.primary
              const [x, , val] = tuple
              return cellColor(x, val) ?? colors.value.primary
            },
          }
        : {}),
    },
    emphasis: {
      itemStyle: {
        shadowBlur: 12,
        shadowColor: c.glow,
      },
    },
  }

  return {
    animationDuration: 500,
    grid: baseGrid({ top: 8, right: 16, bottom: 52, left: 76 }),
    tooltip: {
      position: 'top',
      backgroundColor: c.tooltipBg,
      borderColor: c.tooltipBorder,
      borderWidth: 1,
      padding: [10, 14],
      textStyle: { color: c.tooltipFg, fontSize: 12 },
      extraCssText: 'border-radius: 12px; box-shadow: 0 8px 24px rgba(124,58,237,0.15);',
      formatter: (params: unknown) => {
        const item = (Array.isArray(params) ? params[0] : params) as {
          value?: [number, number, number]
        }
        const tuple = item.value
        if (!tuple) return ''
        const [x, y, val] = tuple
        const col = cols[x] ?? ''
        const row = rows[y] ?? ''
        const meta = pointMeta.value.get(`${col}\0${row}`)
        const metaLine = meta ? `<div style="opacity:0.75;font-size:11px;margin-top:4px">${meta}</div>` : ''
        return `<div style="opacity:0.7;font-size:11px">${row} · ${col}</div><b>${props.valueLabel}: ${val}</b>${metaLine}`
      },
    },
    xAxis: {
      type: 'category',
      data: cols,
      splitArea: { show: true, areaStyle: { color: ['transparent', 'rgba(124,58,237,0.03)'] } },
      axisLine: { show: false },
      axisTick: { show: false },
      axisLabel: { color: c.axis, fontSize: 11, margin: 8 },
    },
    yAxis: {
      type: 'category',
      data: rows,
      splitArea: { show: true, areaStyle: { color: ['transparent', 'rgba(124,58,237,0.03)'] } },
      axisLine: { show: false },
      axisTick: { show: false },
      axisLabel: { color: c.axis, fontSize: 11, margin: 8 },
    },
    visualMap: props.columnColors
      ? undefined
      : {
          min: 0,
          max: maxValue.value,
          calculable: false,
          orient: 'horizontal',
          left: 'center',
          bottom: 4,
          itemWidth: 120,
          itemHeight: 8,
          textStyle: { color: c.axis, fontSize: 10 },
          inRange: {
            color: ['#f5f3ff', '#ddd6fe', '#a78bfa', '#7c3aed'],
          },
        },
    series: [series],
  }
})
</script>

export interface ChartLegendItem {
  label: string
  muted?: boolean
  dashed?: boolean
  color?: string
}

/** 图表系列数据 */
export interface ChartSeries {
  name: string
  data: number[]
  /** 对比线（灰色虚线，如「30 days before」） */
  compare?: boolean
  /** 主系列：面积渐变 + 发光 */
  primary?: boolean
}

export interface ChartFunnelItem {
  name: string
  value: number
}

export interface ChartBarItem {
  name: string
  value: number
}

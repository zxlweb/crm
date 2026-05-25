import type { TagIconPath } from './tag-icons'

/** 将 SVG path d 转为 ECharts path:// 片段（命令 + 逗号分隔坐标） */
export function svgPathDToEchart(d: string): string {
  const tokens: string[] = []
  const re = /([MmLlHhVvCcSsQqTtAaZz])|(-?\d*\.?\d+(?:e[-+]?\d+)?)/g
  let match: RegExpExecArray | null
  while ((match = re.exec(d)) !== null) {
    tokens.push(match[1] ?? match[2]!)
  }

  let out = ''
  let i = 0
  while (i < tokens.length) {
    const cmd = tokens[i]!
    if (!/^[A-Za-z]$/.test(cmd)) {
      i += 1
      continue
    }
    out += cmd
    i += 1
    const nums: string[] = []
    while (i < tokens.length && !/^[A-Za-z]$/.test(tokens[i]!)) {
      nums.push(tokens[i]!)
      i += 1
    }
    if (nums.length) out += nums.join(',')
  }
  return out
}

function circleToEchart(cx: number, cy: number, r: number): string {
  return `M${cx},${cy - r}a${r},${r},0,1,1,0,${2 * r}a${r},${r},0,1,1,0,-${2 * r}`
}

function shapeToEchartFragment(shape: TagIconPath): string {
  switch (shape.tag) {
    case 'circle':
      return circleToEchart(shape.cx ?? 0, shape.cy ?? 0, shape.r ?? 0)
    case 'line':
      return `M${shape.x1},${shape.y1}L${shape.x2},${shape.y2}`
    case 'path':
      return shape.d ? svgPathDToEchart(shape.d) : ''
    case 'rect': {
      const x = shape.x ?? 0
      const y = shape.y ?? 0
      const w = shape.width ?? 0
      const h = shape.height ?? 0
      const rx = shape.rx ?? 0
      if (rx > 0) {
        return `M${x + rx},${y}h${w - 2 * rx}a${rx},${rx},0,0,1,${rx},${rx}v${h - 2 * rx}a${rx},${rx},0,0,1,-${rx},${rx}h-${w - 2 * rx}a${rx},${rx},0,0,1,-${rx},-${rx}v-${h - 2 * rx}a${rx},${rx},0,0,1,${rx},-${rx}z`
      }
      return `M${x},${y}h${w}v${h}h-${w}z`
    }
    case 'polyline':
      return shape.points
        ? `M${shape.points.replace(/\s+/g, ' L').replace(/,/g, ',')}`
        : ''
    default:
      return ''
  }
}

/** 与 UiTagIcon 同源路径，供 ECharts symbol 使用 */
export function tagIconPathsToEchartPath(shapes: TagIconPath[]): string {
  return `path://${shapes.map(shapeToEchartFragment).join('')}`
}

import type { ChartHeatmapPoint } from '@crm/ui-kit'
import type { TeamHeatmapCell } from '~/types/dashboard'
import type { RelationshipHealth } from '~/types/lead'

export type HeatmapDimension = 'health' | 'emotion'

const HEALTH_ORDER: RelationshipHealth[] = ['high', 'medium', 'low']
const EMOTION_ORDER = ['positive', 'neutral', 'negative'] as const

export function listHeatmapMembers(cells: TeamHeatmapCell[]): string[] {
  const seen = new Set<string>()
  const rows: string[] = []
  for (const cell of cells) {
    if (!seen.has(cell.memberId)) {
      seen.add(cell.memberId)
      rows.push(cell.memberName)
    }
  }
  return rows
}

export function buildHeatmapPoints(
  cells: TeamHeatmapCell[],
  dimension: HeatmapDimension,
  label: (key: string) => string,
): ChartHeatmapPoint[] {
  if (dimension === 'health') {
    return cells.map((cell) => ({
      row: cell.memberName,
      column: label(`relationshipHealth.${cell.health}`),
      value: cell.count,
      meta: label(`dashboardHeatmapEmotion.${cell.emotion}`),
    }))
  }

  return cells.map((cell) => ({
    row: cell.memberName,
    column: label(`dashboardHeatmapEmotionShort.${cell.emotion}`),
    value: cell.count,
    meta: label(`relationshipHealth.${cell.health}`),
  }))
}

export function heatmapColumns(
  dimension: HeatmapDimension,
  label: (key: string) => string,
): string[] {
  if (dimension === 'health') {
    return HEALTH_ORDER.map((h) => label(`relationshipHealth.${h}`))
  }
  return EMOTION_ORDER.map((e) => label(`dashboardHeatmapEmotionShort.${e}`))
}

/** 情绪维：积极紫 / 中性灰 / 消极琥珀 */
export const EMOTION_HEATMAP_COLUMN_COLORS = ['#7c3aed', '#a1a1aa', '#f97316']
